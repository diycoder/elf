// Package errors provides a way to return detailed information
// The error is normally JSON encoded.
package errors

import (
	"encoding/json"
	"fmt"
)

type Error struct {
	Code   int    `json:"code,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// New generates a custom error.
func New(code int, detail string) error {
	return &Error{
		Code:   code,
		Detail: detail,
	}
}

// Parse tries to parse a JSON string into an error. If that
// fails, it will set the given string as the error detail.
func Parse(err string) *Error {
	e := new(Error)
	errr := json.Unmarshal([]byte(err), e)
	if errr != nil {
		e.Detail = err
	}
	return e
}

// BadRequest generates a 400 error.
func BadRequest(format string, a ...interface{}) error {
	return &Error{
		Code:   400,
		Detail: fmt.Sprintf(format, a...),
	}
}

// Unauthorized generates a 401 error.
func Unauthorized(format string, a ...interface{}) error {
	return &Error{
		Code:   401,
		Detail: fmt.Sprintf(format, a...),
	}
}

// Forbidden generates a 403 error.
func Forbidden(format string, a ...interface{}) error {
	return &Error{
		Code:   403,
		Detail: fmt.Sprintf(format, a...),
	}
}

// NotFound generates a 404 error.
func NotFound(format string, a ...interface{}) error {
	return &Error{
		Code:   404,
		Detail: fmt.Sprintf(format, a...),
	}
}

// MethodNotAllowed generates a 405 error.
func MethodNotAllowed(format string, a ...interface{}) error {
	return &Error{
		Code:   405,
		Detail: fmt.Sprintf(format, a...),
	}
}

// Timeout generates a 408 error.
func Timeout(format string, a ...interface{}) error {
	return &Error{
		Code:   408,
		Detail: fmt.Sprintf(format, a...),
	}
}

// Conflict generates a 409 error.
func Conflict(format string, a ...interface{}) error {
	return &Error{
		Code:   409,
		Detail: fmt.Sprintf(format, a...),
	}
}

// InternalServerError generates a 500 error.
func InternalServerError(format string, a ...interface{}) error {
	return &Error{
		Code:   500,
		Detail: fmt.Sprintf(format, a...),
	}
}

// Equal tries to compare errors
func Equal(err1 error, err2 error) bool {
	verr1, ok1 := err1.(*Error)
	verr2, ok2 := err2.(*Error)

	if ok1 != ok2 {
		return false
	}

	if !ok1 {
		return err1 == err2
	}

	if verr1.Code != verr2.Code {
		return false
	}

	return true
}

// FromError try to convert go error to *Error
func FromError(err error) *Error {
	if verr, ok := err.(*Error); ok && verr != nil {
		return verr
	}

	return Parse(err.Error())
}
