// https://github.com/topfreegames/mystack-logger
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/topfreegames/mystack-logger/api"
	"github.com/topfreegames/mystack-logger/metadata"
	"github.com/topfreegames/mystack-logger/storage"
)

var _ = Describe("Healthcheck Handler", func() {
	var request *http.Request
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		// Record HTTP responses.
		recorder = httptest.NewRecorder()
	})

	Describe("GET /healthcheck", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/healthcheck", nil)
		})

		Context("when all services healthy", func() {
			It("returns a status code of 200", func() {
				app.Router.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})

			It("returns working string", func() {
				app.Router.ServeHTTP(recorder, request)
				Expect(recorder.Body.String()).To(Equal(`{"healthy": true}`))
			})

			It("returns the version as a header", func() {
				app.Router.ServeHTTP(recorder, request)
				Expect(recorder.Header().Get("x-mystack-version")).To(Equal(metadata.Version))
			})

			It("returns status code of 500 if redis is unavailable", func() {
				config := viper.New()
				config.Set("redis.url", "redis://:@localhost:11111")
				logger := logrus.New()
				storageAdapter, _ := storage.NewRedisStorageAdapter(config)
				app, _ := api.NewApp("localhost", 8686, config, logger, storageAdapter, collector, false)
				app.Router.ServeHTTP(recorder, request)
				var obj map[string]interface{}
				err := json.Unmarshal([]byte(recorder.Body.String()), &obj)
				Expect(err).NotTo(HaveOccurred())
				Expect(obj["code"]).To(Equal("MST-001"))
				Expect(obj["error"]).To(Equal("error performing healthcheck"))
				Expect(obj["description"]).Should(SatisfyAny(
					Equal("dial tcp [::1]:11111: getsockopt: connection refused"),
					Equal("dial tcp 127.0.0.1:11111: getsockopt: connection refused"),
				))
			})
		})
	})
})
