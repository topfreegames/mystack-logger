// mystack/mystack-logger
// https://github.com/topfreegames/mystack-controller
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	r "gopkg.in/redis.v5"
)

type message struct {
	app         string
	messageBody string
}

func newMessage(app string, messageBody string) *message {
	return &message{
		app:         app,
		messageBody: messageBody,
	}
}

type messagePipeliner struct {
	bufferSize    int
	messageCount  int
	pipeline      *r.Pipeline
	timeoutTicker *time.Ticker
	queuedApps    map[string]bool
	errCh         chan error
}

func newMessagePipeliner(bufferSize int, redisClient *r.Client, timeout time.Duration, errCh chan error) *messagePipeliner {
	return &messagePipeliner{
		bufferSize:    bufferSize,
		pipeline:      redisClient.Pipeline(),
		timeoutTicker: time.NewTicker(timeout),
		queuedApps:    map[string]bool{},
		errCh:         errCh,
	}
}

func (mp *messagePipeliner) addMessage(message *message) {
	if err := mp.pipeline.RPush(message.app, message.messageBody).Err(); err == nil {
		mp.queuedApps[message.app] = true
		mp.messageCount++
	} else {
		mp.errCh <- fmt.Errorf("Error adding rpush to %s to the pipeline: %s", message.app, err)
	}
}

func (mp messagePipeliner) execPipeline() {
	for app := range mp.queuedApps {
		if err := mp.pipeline.LTrim(app, int64(-1*mp.bufferSize), -1).Err(); err != nil {
			mp.errCh <- fmt.Errorf("Error adding ltrim of %s to the pipeline: %s", app, err)
		}
	}
	if _, err := mp.pipeline.Exec(); err != nil {
		mp.errCh <- fmt.Errorf("Error executing pipeline: %s", err)
	}
}

type redisAdapter struct {
	started         bool
	bufferSize      int
	redisClient     *r.Client
	redisURL        string
	messageChannel  chan *message
	stopCh          chan struct{}
	config          *viper.Viper
	pipelineTimeout time.Duration
	pipelineLength  int
}

// NewRedisStorageAdapter returns a pointer to a new instance of a redis-based storage.Adapter.
func NewRedisStorageAdapter(config *viper.Viper) (Adapter, error) {
	rsa := &redisAdapter{
		messageChannel: make(chan *message),
		stopCh:         make(chan struct{}),
		config:         config,
	}
	err := rsa.configure()
	if err != nil {
		return nil, err
	}
	err = rsa.configureRedis()
	if err != nil {
		return nil, err
	}
	return rsa, nil
}

func (a *redisAdapter) configureRedis() error {
	redisCfg, err := r.ParseURL(a.config.GetString("redis.url"))
	if err != nil {
		return err
	}
	redisClient := r.NewClient(redisCfg)
	a.redisClient = redisClient
	return nil
}

func (a *redisAdapter) configure() error {
	a.loadConfigurationDefaults()
	a.redisURL = a.config.GetString("redis.url")
	a.pipelineTimeout = time.Duration(a.config.GetInt("redis.pipeline-timeout")) * time.Second
	a.pipelineLength = a.config.GetInt("redis.pipeline-length")
	a.bufferSize = a.config.GetInt("log-buffer-size")
	if a.bufferSize <= 0 {
		return fmt.Errorf("Invalid buffer size: %d", a.bufferSize)
	}
	return nil
}

func (a *redisAdapter) loadConfigurationDefaults() {
	a.config.SetDefault("redis.url", "localhost:6399")
	a.config.SetDefault("redis.pipeline-timeout", 1)
	a.config.SetDefault("redis.pipeline-length", 20)
	a.config.SetDefault("log-buffer-size", 1000)
}

// Start the storage adapter. Invocations of this function are not concurrency safe and multiple
// serialized invocations have no effect.
func (a *redisAdapter) Start() {
	if !a.started {
		a.started = true
		errCh := make(chan error)
		mp := newMessagePipeliner(a.bufferSize, a.redisClient, a.pipelineTimeout, errCh)
		go func() {
			defer mp.pipeline.Close()
			for {
				select {
				case err := <-errCh:
					log.Println(err)
				case <-a.stopCh:
					return
				case message := <-a.messageChannel:
					mp.addMessage(message)
					if mp.messageCount == a.pipelineLength {
						go mp.execPipeline()
					}
				case <-mp.timeoutTicker.C:
					go mp.execPipeline()
				}
			}
		}()
	}
}

// Write adds a log message to to an app-specific list in redis using ring-buffer-like semantics
func (a *redisAdapter) Write(app string, messageBody string) error {
	a.messageChannel <- newMessage(app, messageBody)
	return nil
}

// Read retrieves a specified number of log lines from an app-specific list in redis
func (a *redisAdapter) Read(app string, lines int) ([]string, error) {
	stringSliceCmd := a.redisClient.LRange(app, int64(-1*lines), -1)
	result, err := stringSliceCmd.Result()
	if err != nil {
		return nil, err
	}
	if len(result) > 0 {
		return result, nil
	}
	return nil, fmt.Errorf("Could not find logs for '%s'", app)
}

// Destroy deletes an app-specific list from redis
func (a *redisAdapter) Destroy(app string) error {
	if err := a.redisClient.Del(app).Err(); err != nil {
		return err
	}
	return nil
}

// Reopen the storage adapter-- in the case of this implementation, a no-op
func (a *redisAdapter) Reopen() error {
	return nil
}

// Stop the storage adapter. Additional writes may not be performed after stopping.
func (a *redisAdapter) Stop() {
	close(a.stopCh)
}
