// https://github.com/topfreegames/mystack-logger
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/topfreegames/mystack-logger/logger"
	"github.com/topfreegames/mystack-logger/storage"
)

//LogsHandler handler
type LogsHandler struct {
	App            *App
	storageAdapter storage.Adapter
	logger         *logrus.Logger
}

// NewLogsHandler ctor
func NewLogsHandler(
	app *App,
	storageAdapter storage.Adapter,
	logger *logrus.Logger,
) *LogsHandler {
	return &LogsHandler{
		App:            app,
		storageAdapter: storageAdapter,
		logger:         logger,
	}
}

func getVars(r *http.Request) (string, string) {
	vars := mux.Vars(r)
	app := vars["app"]
	splited := strings.Split(r.URL.Path, "/")
	if len(app) == 0 {
		app = splited[3]
	}

	email := emailFromCtx(r.Context())
	user := userFromEmail(email)

	return app, user
}

func userFromEmail(email string) string {
	user := strings.Split(email, "@")[0]
	user = strings.Replace(user, ".", "-", -1)
	return user
}

type flushWriter struct {
	f http.Flusher
	w io.Writer
}

func (fw *flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	if fw.f != nil {
		fw.f.Flush()
	}
	return
}

//ServeHTTP method
func (l *LogsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	lines := 50
	follow := false
	var err error
	if args, ok := query["lines"]; ok {
		if len(args) > 0 {
			if lines, err = strconv.Atoi(args[0]); err != nil {
				lines = 50
			}
		}
	}

	if args, ok := query["follow"]; ok {
		if len(args) > 0 {
			if args[0] == "true" {
				follow = true
			}
		}
	}

	app, user := getVars(r)
	log := l.logger.WithFields(logrus.Fields{
		"user":   user,
		"app":    app,
		"source": "api/logs.go",
		"lines":  lines,
	})
	log.Debug("getting logs")
	logs, err := l.storageAdapter.Read(fmt.Sprintf("%s-%s", app, user), lines)
	if err != nil {
		if strings.HasPrefix(err.Error(), "Could not find logs for") {
			w.WriteHeader(http.StatusNoContent)
		} else {
			log.Error(err)
			l.App.HandleError(w, http.StatusBadRequest, "error getting app logs", err)
		}
		return
	}
	if follow {
		// hack fox nginx
		w.Header().Add("X-Accel-Buffering", "no")
	}
	fw := flushWriter{w: w}
	if f, ok := w.(http.Flusher); ok {
		fw.f = f
	}
	for _, line := range logs {
		fw.Write([]byte(fmt.Sprintf("%s\n", strings.TrimSuffix(line, "\n"))))
	}

	if follow {
		followerChan := make(chan []byte)
		follower := logger.NewLogFollower(followerChan)
		f := l.App.Collector.AddFollower(follower)
		closedChan := w.(http.CloseNotifier).CloseNotify()
		close := false
		for !close {
			select {
			case msg := <-followerChan:
				message := new(logger.Message)
				if err := json.Unmarshal(msg, message); err != nil {
					log.Error(err)
				} else {
					if message.Kubernetes.Labels["app"] == app {
						fw.Write([]byte(fmt.Sprintf("%s\n", strings.TrimSuffix(logger.BuildApplicationLogMessage(message), "\n"))))
					}
				}
			case <-closedChan:
				l.App.Collector.RemoveFollower(f)
				log.Debug("exiting log streaming...")
				close = true
			}
		}
	}
	log.Debug("logs done.")
}
