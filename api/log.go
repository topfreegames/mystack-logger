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
	vars := mux.Vars(r)
	app := vars["app"]
	user := vars["user"]
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
