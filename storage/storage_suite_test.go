// https://github.com/topfreegames/mystack-logger
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package storage_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"

	"testing"
)

var config *viper.Viper
var app string

func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MyStack Logger - Storage Suite")
}

var _ = BeforeSuite(func() {
	config = viper.New()
	app = "test-app"
})
