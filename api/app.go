// https://github.com/topfreegames/mystack-logger
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/topfreegames/mystack-logger/errors"
	"github.com/topfreegames/mystack-logger/logger"
	"github.com/topfreegames/mystack-logger/storage"
)

//App is our API application
type App struct {
	Address        string
	Config         *viper.Viper
	storageAdapter storage.Adapter
	Logger         *logrus.Logger
	Router         *mux.Router
	Server         *http.Server
	EmailDomain    []string
	Unsecure       bool
	Collector      *logger.LogCollector
}

//NewApp is the app constructor
func NewApp(
	host string, port int,
	config *viper.Viper,
	logger *logrus.Logger,
	storageAdapter storage.Adapter,
	collector *logger.LogCollector,
	unsecure bool) (*App, error) {
	a := &App{
		Config:         config,
		Address:        fmt.Sprintf("%s:%d", host, port),
		Logger:         logger,
		storageAdapter: storageAdapter,
		EmailDomain:    config.GetStringSlice("oauth.acceptedDomains"),
		Unsecure:       unsecure,
		Collector:      collector,
	}
	err := a.configureApp()
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) getRouter() *mux.Router {
	r := mux.NewRouter()
	//TODO fix healthcheck
	r.Handle("/healthcheck", Chain(
		&HealthcheckHandler{
			App:            a,
			storageAdapter: a.storageAdapter,
		},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
	)).Methods("GET").Name("healthcheck")

	r.Handle("/logs/apps/{app}", Chain(
		&LogsHandler{
			App:            a,
			storageAdapter: a.storageAdapter,
			logger:         a.Logger,
		},
		&LoggingMiddleware{App: a},
		&VersionMiddleware{},
		&AccessMiddleware{App: a},
	)).Methods("GET").Name("logs")

	return r
}

func (a *App) configureApp() error {
	a.configureServer()
	return nil
}

func (a *App) configureServer() {
	a.Router = a.getRouter()
	a.Server = &http.Server{
		Addr:         a.Address,
		Handler:      a.Router,
		WriteTimeout: 20 * time.Minute,
	}
}

//HandleError writes an error response with message and status
func (a *App) HandleError(w http.ResponseWriter, status int, msg string, err interface{}) {
	w.WriteHeader(status)
	var sErr errors.SerializableError
	val, ok := err.(errors.SerializableError)
	if ok {
		sErr = val
	} else {
		sErr = errors.NewGenericError(msg, err.(error))
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(sErr.Serialize())
}

//ListenAndServe requests
func (a *App) ListenAndServe() (io.Closer, error) {
	listener, err := net.Listen("tcp", a.Address)
	if err != nil {
		return nil, err
	}

	err = a.Server.Serve(listener)
	if err != nil {
		listener.Close()
		return nil, err
	}

	return listener, nil
}
