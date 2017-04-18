// https://github.com/topfreegames/mystack-logger
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package errors

import "encoding/json"

//AccessError happens when an unidentified error occurs
type AccessError struct {
	Message     string
	SourceError error
}

//NewAccessError ctor
func NewAccessError(message string, err error) *AccessError {
	return &AccessError{
		Message:     message,
		SourceError: err,
	}
}

func (e *AccessError) Error() string {
	return e.SourceError.Error()
}

//Serialize returns the error serialized
func (e *AccessError) Serialize() []byte {
	g, _ := json.Marshal(map[string]interface{}{
		"code":        "OFF-002",
		"error":       e.Message,
		"description": e.SourceError.Error(),
	})

	return g
}
