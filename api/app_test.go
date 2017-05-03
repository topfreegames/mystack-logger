// https://github.com/topfreegames/mystack-logger
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/topfreegames/mystack-logger/api"
)

var _ = Describe("App", func() {

	Describe("NewApp", func() {
		It("should return new app", func() {
			application, err := api.NewApp("0.0.0.0", 8686, config, log, storageAdapter, collector, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(application).NotTo(BeNil())
			Expect(application.Address).NotTo(Equal(""))
			Expect(application.Logger).NotTo(BeNil())
			Expect(application.Router).NotTo(BeNil())
			Expect(application.Server).NotTo(BeNil())
			Expect(application.Config).To(Equal(config))
		})
	})
})
