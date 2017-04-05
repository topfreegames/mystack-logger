// mystack/mystack-logger
// https://github.com/topfreegames/mystack-controller
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package logger

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/topfreegames/mystack-logger/storage"
)

const (
	podPattern = `(\w.*)-(\w.*)-(\w.*)-(\w.*)`
	timeFormat = "2006-01-02T15:04:05-07:00"
)

var (
	podRegex = regexp.MustCompile(podPattern)
)

func handle(rawMessage []byte, storageAdapter storage.Adapter) error {
	message := new(Message)
	if err := json.Unmarshal(rawMessage, message); err != nil {
		return err
	}
	labels := message.Kubernetes.Labels
	storageAdapter.Write(fmt.Sprintf("%s-%s", labels["app"], labels["mystack/owner"]), buildApplicationLogMessage(message))
	return nil
}

func buildApplicationLogMessage(message *Message) string {
	p := podRegex.FindStringSubmatch(message.Kubernetes.PodName)
	tag := fmt.Sprintf(
		"%s.%s",
		message.Kubernetes.Labels["type"],
		message.Kubernetes.Labels["version"])
	if len(p) > 0 {
		tag = fmt.Sprintf("%s.%s", tag, p[len(p)-1])
	}
	return fmt.Sprintf("%s %s: %s",
		message.Time.Format(timeFormat),
		message.Kubernetes.Labels["app"],
		message.Log)
}
