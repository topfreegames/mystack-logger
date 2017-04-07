// https://github.com/topfreegames/mystack-logger
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"net/http"

	"github.com/topfreegames/mystack-logger/storage"
)

//HealthcheckHandler handler
type HealthcheckHandler struct {
	App            *App
	storageAdapter storage.Adapter
}

//ServeHTTP method
func (h *HealthcheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l := loggerFromContext(r.Context())
	l.Debug("Performing healthcheck...")
	_, err := h.storageAdapter.Healthcheck()
	if err != nil {
		l.Error(err)
		h.App.HandleError(w, http.StatusBadRequest, "error performing healthcheck", err)
		return
	}
	Write(w, http.StatusOK, `{"healthy": true}`)
	l.Debug("Healthcheck done.")
}
