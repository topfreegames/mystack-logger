// https://github.com/topfreegames/mystack-logger
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package logger_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"github.com/topfreegames/mystack-logger/storage"

	"testing"
)

var config *viper.Viper
var storageAdapter storage.Adapter
var app string

func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MyStack Logger - Log Suite")
}

var _ = BeforeSuite(func() {
	config = viper.New()
	storageAdapter, _ = storage.NewRedisStorageAdapter(config)
	app = "test-app"
})
