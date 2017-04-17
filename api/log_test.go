// https://github.com/topfreegames/mystack-logger
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	. "github.com/topfreegames/mystack-logger/api"
)

var _ = Describe("Log Handler", func() {
	var request *http.Request
	var recorder *httptest.ResponseRecorder
	var logger *logrus.Logger
	var log *LogsHandler

	BeforeEach(func() {
		// Record HTTP responses.
		recorder = httptest.NewRecorder()

		logger = logrus.New()
		log = NewLogsHandler(app, storageAdapter, logger)
	})

	Describe("GET /logs/users/{user}/apps/{app}", func() {
		Context("when all services healthy", func() {
			It("returns a status code of 204 if no logs for app", func() {
				request, _ = http.NewRequest("GET", "/logs/users/testuser/apps/testapp", nil)
				log.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(204))
			})

			It("returns last logs if exists", func() {
				storageAdapter.Start()
				for i := 0; i < 5; i++ {
					err := storageAdapter.Write("testapp2-testuser", fmt.Sprintf("message %d", i))
					Expect(err).NotTo(HaveOccurred())
				}
				Eventually(func() int {
					messages, _ := storageAdapter.Read("testapp2-testuser", 8)
					return len(messages)
				}, 10).Should(Equal(5))

				request, _ = http.NewRequest("GET", "/logs/users/testuser/apps/testapp2", nil)
				log.ServeHTTP(recorder, request)
				fmt.Printf(recorder.Body.String())
				Expect(recorder.Body.String()).To(Equal(`message 0
message 1
message 2
message 3
message 4
`))
				err := storageAdapter.Destroy("testapp2-testuser")
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns only x last logs if exists", func() {
				storageAdapter.Start()
				for i := 0; i < 5; i++ {
					err := storageAdapter.Write("testapp2-testuser", fmt.Sprintf("message %d", i))
					Expect(err).NotTo(HaveOccurred())
				}
				Eventually(func() int {
					messages, _ := storageAdapter.Read("testapp2-testuser", 2)
					return len(messages)
				}, 10).Should(Equal(2))

				request, _ = http.NewRequest("GET", "/logs/users/testuser/apps/testapp2?lines=2", nil)
				log.ServeHTTP(recorder, request)
				fmt.Printf(recorder.Body.String())
				Expect(recorder.Body.String()).To(Equal(`message 3
message 4
`))
				err := storageAdapter.Destroy("testapp2-testuser")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
