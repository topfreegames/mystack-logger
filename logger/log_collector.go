// mystack/mystack-logger
// https://github.com/topfreegames/mystack/mystack-controller
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package logger

import (
	"errors"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/spf13/viper"
	"github.com/topfreegames/mystack-logger/storage"
)

// LogCollector is the the module that will grab logs // from nsq and write them into redis
type LogCollector struct {
	config          *viper.Viper
	storageAdapter  storage.Adapter
	nsqConsumer     *nsq.Consumer
	nsqHandlerCount int
	nsqLogsTopic    string
	nsqdURL         string
	redisURL        string
	Listening       bool
	stopTimeout     time.Duration
}

// NewLogCollector instantiates a new LogCollector
func NewLogCollector(storageAdapter storage.Adapter, config *viper.Viper) *LogCollector {
	l := &LogCollector{
		config:         config,
		storageAdapter: storageAdapter,
	}
	l.configure()
	l.configureNsqConsumer()
	return l
}

func (l *LogCollector) configure() {
	l.loadConfigurationDefaults()
	l.nsqdURL = l.config.GetString("nsqd.url")
	l.redisURL = l.config.GetString("nsqd.url")
	l.nsqLogsTopic = l.config.GetString("nsqd.logs-topic")
	l.nsqHandlerCount = l.config.GetInt("nsqd.handler-count")
	l.stopTimeout = time.Duration(l.config.GetInt("stop-timeout")) * time.Second
}

func (l *LogCollector) loadConfigurationDefaults() {
	l.config.SetDefault("nsqd.url", "localhost:4155")
	l.config.SetDefault("nsqd.logs-topic", "logs")
	l.config.SetDefault("nsqd.handler-count", 30)
	l.config.SetDefault("stop-timeout", 30)
}

func (l *LogCollector) configureNsqConsumer() error {
	cfg := nsq.NewConfig()
	c, err := nsq.NewConsumer(l.nsqLogsTopic, "consume", cfg)
	if err != nil {
		return err
	}
	l.nsqConsumer = c
	return nil
}

// Start starts logcollector
func (l *LogCollector) Start() error {
	if !l.Listening {
		l.Listening = true
		l.nsqConsumer.AddConcurrentHandlers(nsq.HandlerFunc(func(msg *nsq.Message) error {
			if err := Handle(msg.Body, l.storageAdapter); err != nil {
				msg.Requeue(-1)
				return err
			}
			return nil
		}), l.nsqHandlerCount)
		err := l.nsqConsumer.ConnectToNSQD(l.nsqdURL)
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop stops the collector
func (l *LogCollector) Stop() error {
	l.nsqConsumer.Stop()
	tmr := time.NewTimer(l.stopTimeout)
	defer tmr.Stop()
	select {
	case <-tmr.C:
		return errors.New("timeout waiting for logcollector to stop")
	case <-l.nsqConsumer.StopChan:
		return nil
	}
}

// Stopped is the Aggregator interface implementation
func (l *LogCollector) Stopped() <-chan error {
	retCh := make(chan error)
	go func() {
		<-l.nsqConsumer.StopChan
		retCh <- nil
	}()
	return retCh
}
