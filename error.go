package s3

import (
	"github.com/ONSdigital/log.go/v2/log"
)

// Error is the s3 package's error type
type Error struct {
	err     error
	logData map[string]interface{}
}

// NewError creates a new Error
func NewError(err error, logData map[string]interface{}) *Error {
	return &Error{
		err:     err,
		logData: logData,
	}
}

// Error implements the Go standard error interface
func (e *Error) Error() string {
	if e.err == nil {
		return "nil"
	}
	return e.err.Error()
}

// LogData implements the DataLogger interface which allows you extract
// embedded log.Data from an error
func (e *Error) LogData() map[string]interface{} {
	if e.logData == nil {
		return log.Data{}
	}
	return e.logData
}

// Unwrap returns the wrapped error
func (e *Error) Unwrap() error {
	return e.err
}
