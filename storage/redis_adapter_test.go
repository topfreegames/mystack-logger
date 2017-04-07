// mystack/mystack-logger
// https://github.com/topfreegames/mystack-controller
//
// licensed under the mit license:
// http://www.opensource.org/licenses/mit-license
// copyright © 2017 top free games <backend@tfgco.com>

package storage_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	. "github.com/topfreegames/mystack-logger/storage"
)

var _ = Describe("Redis adapter test", func() {
	It("TestRedisReadFromNonExistingApp", func() {
		a, err := NewRedisStorageAdapter(config)
		Expect(err).NotTo(HaveOccurred())
		messages, err := a.Read(app, 10)
		Expect(messages).To(BeNil())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(fmt.Sprintf("Could not find logs for '%s'", app)))
	})

	It("TestRedisWithBadBufferSizes", func() {
		for _, size := range []int{-1, 0} {
			config := viper.New()
			config.Set("log-buffer-size", size)
			a, err := NewRedisStorageAdapter(config)
			Expect(a).To(BeNil())
			Expect(err).To(HaveOccurred())
		}
	})

	It("TestRedisLogs", func() {
		config := viper.New()
		config.Set("log-buffer-size", 10)
		a, err := NewRedisStorageAdapter(config)
		Expect(err).NotTo(HaveOccurred())
		a.Start()
		defer a.Stop()
		for i := 0; i < 5; i++ {
			err := a.Write(app, fmt.Sprintf("message %d", i))
			Expect(err).NotTo(HaveOccurred())
		}
		Eventually(func() int {
			messages, _ := a.Read(app, 8)
			return len(messages)
		}, 10).Should(Equal(5))
		messages, err := a.Read(app, 3)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(messages)).To(Equal(3))
		for i := 0; i < 3; i++ {
			expectedMessage := fmt.Sprintf("message %d", i+2)
			Expect(messages[i]).To(Equal(expectedMessage))
		}
		for i := 5; i < 11; i++ {
			err = a.Write(app, fmt.Sprintf("message %d", i))
			Expect(err).NotTo(HaveOccurred())
		}
		time.Sleep(time.Second * 2)
		messages, err = a.Read(app, 20)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(messages)).To(Equal(10))
		for i := 0; i < 10; i++ {
			expectedMessage := fmt.Sprintf("message %d", i+1)
			Expect(messages[i]).To(Equal(expectedMessage))
		}
	})

	It("TestRedisDestroy", func() {
		a, err := NewRedisStorageAdapter(config)
		Expect(err).NotTo(HaveOccurred())
		a.Start()
		defer a.Stop()
		err = a.Write(app, "Hello, log!")
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() bool {
			exists, _ := a.(*RedisAdapter).RedisClient.Exists(app).Result()
			return exists
		}, 10).Should(BeTrue())
		err = a.Destroy(app)
		Expect(err).NotTo(HaveOccurred())
		exists, err := a.(*RedisAdapter).RedisClient.Exists(app).Result()
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeFalse())
	})
})
