// https://github.com/topfreegames/mystack-logger
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/topfreegames/mystack-logger/storage"
)

//LogsHandler handler
type LogsHandler struct {
	App            *App
	storageAdapter storage.Adapter
	logger         *logrus.Logger
}

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

//ServeHTTP method
func (l *LogsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	lines := 50
	var err error
	if args, ok := query["lines"]; ok {
		if len(args) > 0 {
			if lines, err = strconv.Atoi(args[0]); err != nil {
				lines = 50
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
	for _, line := range logs {
		fmt.Fprintf(w, "%s\n", strings.TrimSuffix(line, "\n"))
	}
	log.Debug("logs done.")
}
