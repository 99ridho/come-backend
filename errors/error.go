package errors

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message string `json:"error_message"`
	Code    int    `json:"error_code"`
}

func NewError(message string, code int) *Error {
	return &Error{
		Message: message,
		Code:    code,
	}
}

func NewErrorWithStatusCode(code int) *Error {
	return &Error{
		Message: http.StatusText(code),
		Code:    code,
	}
}

func (e *Error) WriteTo(w http.ResponseWriter) {
	w.WriteHeader(e.Code)
	json.NewEncoder(w).Encode(e)
}

func (e *Error) Error() string {
	return e.Message
}
