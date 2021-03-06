// https://github.com/topfreegames/mystack-logger
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/topfreegames/mystack-logger/api"
	"github.com/topfreegames/mystack-logger/logger"
	"github.com/topfreegames/mystack-logger/storage"

	"testing"
)

var config *viper.Viper
var app *api.App
var storageAdapter storage.Adapter
var log *logrus.Logger
var collector *logger.LogCollector

func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MyStack Logger - Api Suite")
}

var _ = BeforeSuite(func() {
	config = viper.New()
	config.Set("log-buffer-size", 10)
	config.Set("redis.pipeline-timeout", 1)
	log = logrus.New()
	storageAdapter, _ = storage.NewRedisStorageAdapter(config)
	collector = logger.NewLogCollector(storageAdapter, config)
	app, _ = api.NewApp("localhost", 8686, config, log, storageAdapter, collector, false)
})
