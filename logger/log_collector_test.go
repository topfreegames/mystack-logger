// https://github.com/topfreegames/mystack-logger
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package logger_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/topfreegames/mystack-logger/logger"
)

var _ = Describe("Log Collector", func() {

	var collector logger.Aggregator

	BeforeEach(func() {
		collector = logger.NewLogCollector(storageAdapter, config)
	})

	Context("Start", func() {
		It("should start", func() {
			defer collector.Stop()
			Expect(collector.(*logger.LogCollector).Listening).To(Equal(false))
			collector.Start()
			Expect(collector.(*logger.LogCollector).Listening).To(Equal(true))
		})
	})
})
