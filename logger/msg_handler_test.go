// mystack/mystack-logger
// https://github.com/topfreegames/mystack-controller
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

var _ = Describe("Message Handler", func() {

	var (
		validAppMessage = `{"log": "test message", "stream": "stderr", "time": "2016-10-18T20:29:38+00:00", "docker": {"container_id": "containerId"}, "kubernetes": {"namespace_name": "foo", "pod_id": "podId", "pod_name": "foo-web-845861952-nzf60", "container_name": "foo-web", "labels": {"app": "foo",
		"heritage": "mystack", "mystack/owner":"testuser", "type": "web", "version": "v2"}, "host": "host"}}`

		badjson = `{"log":}`
	)

	It("should handle and persist valid message", func() {
		storageAdapter.Start()
		defer storageAdapter.Destroy("testapp2-testuser")
		err := logger.Handle([]byte(validAppMessage), storageAdapter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(func() int {
			messages, _ := storageAdapter.Read("foo-testuser", 1)
			return len(messages)
		}, 10).Should(Equal(1))
	})

	It("should not handle invalid json", func() {
		storageAdapter.Start()
		defer storageAdapter.Destroy("testapp2-testuser")
		err := logger.Handle([]byte(badjson), storageAdapter)
		Expect(err).To(HaveOccurred())
	})

})
