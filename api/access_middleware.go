// https://github.com/topfreegames/mystack-logger
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//AccessMiddleware guarantees that the user is logged
type AccessMiddleware struct {
	App  *App
	next http.Handler
}

const emailKey = contextKey("emailKey")

//NewContextWithEmail save email on context
func NewContextWithEmail(ctx context.Context, email string) context.Context {
	c := context.WithValue(ctx, emailKey, email)
	return c
}

func emailFromCtx(ctx context.Context) string {
	email := ctx.Value(emailKey)
	if email == nil {
		return ""
	}
	return email.(string)
}

//ServeHTTP methods
func (m *AccessMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger := loggerFromContext(r.Context())
	logger.Info("Checking access token")
	if m.App.Unsecure {
		logger.Debug("unsecure mode detected, defaulting user to testuser.example.com")
		ctx := NewContextWithEmail(r.Context(), "testuser")
		m.next.ServeHTTP(w, r.WithContext(ctx))
	} else {
		logger := loggerFromContext(r.Context())
		logger.Info("Checking access token")

		accessToken := r.Header.Get("Authorization")
		accessToken = strings.TrimPrefix(accessToken, "Bearer ")

		url := fmt.Sprintf("http://mystack-controller:8080/users?token=%s", accessToken)
		resp, err := http.Get(url)
		if err != nil {
			m.App.HandleError(w, http.StatusInternalServerError, "access error", err)
			return
		}

		defer resp.Body.Close()
		bts, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			m.App.HandleError(w, resp.StatusCode, "access error", fmt.Errorf(string(bts)))
			return
		}

		obj := make(map[string]string)
		err = json.Unmarshal(bts, &obj)
		if err != nil {
			m.App.HandleError(w, http.StatusUnauthorized, "access error", err)
			return
		}

		ctx := NewContextWithEmail(r.Context(), obj["email"])

		logger.Info("Access token checked")
		m.next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (m *AccessMiddleware) verifyEmailDomain(email string) bool {
	for _, domain := range m.App.EmailDomain {
		if strings.HasSuffix(email, domain) {
			return true
		}
	}
	return false
}

//SetNext handler
func (m *AccessMiddleware) SetNext(next http.Handler) {
	m.next = next
}
