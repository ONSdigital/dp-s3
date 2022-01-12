package s3

import (
	"github.com/ONSdigital/log.go/v2/log"
)

// S3Error is the s3 package's error type
type S3Error struct {
	err     error
	logData map[string]interface{}
}

// NewError creates a new S3Error
func NewError(err error, logData map[string]interface{}) *S3Error {
	return &S3Error{
		err:     err,
		logData: logData,
	}
}

// S3Error implements the Go standard error interface
func (e *S3Error) Error() string {
	if e.err == nil {
		return "nil"
	}
	return e.err.Error()
}

// LogData implements the DataLogger interface which allows you extract
// embedded log.Data from an
//error
func (e *S3Error) LogData() map[string]interface{} {
	if e.logData == nil {
		return log.Data{}
	}
	return e.logData
}

// Unwrap returns the wrapped error
func (e *S3Error) Unwrap() error {
	return e.err
}

// ErrUnexpectedRegion if a request tried to access an unexpected region
type ErrUnexpectedRegion struct {
	S3Error
}

func NewUnexpectedRegionError(err error, logData map[string]interface{}) *ErrUnexpectedRegion {
	return &ErrUnexpectedRegion{
		S3Error: S3Error{
			err:     err,
			logData: logData,
		},
	}
}

// ErrUnexpectedBucket if a request tried to access an unexpected bucket
type ErrUnexpectedBucket struct {
	S3Error
}

func NewUnexpectedBucketError(err error, logData map[string]interface{}) *ErrUnexpectedBucket {
	return &ErrUnexpectedBucket{
		S3Error: S3Error{
			err:     err,
			logData: logData,
		},
	}
}

// ErrNotUploaded if an s3Key could not be found in ListMultipartUploads
type ErrNotUploaded struct {
	S3Error
}

func NewErrNotUploaded(err error, logData map[string]interface{}) *ErrNotUploaded {
	return &ErrNotUploaded{
		S3Error: S3Error{
			err:     err,
			logData: logData,
		},
	}
}

// ErrListParts represents an error returned by S3 ListParts
type ErrListParts struct {
	S3Error
}

func NewListPartsError(err error, logData map[string]interface{}) *ErrListParts {
	return &ErrListParts{
		S3Error: S3Error{
			err:     err,
			logData: logData,
		},
	}
}

// ErrChunkNumberNotFound if a chunk number could not be found in an existing multipart upload.
type ErrChunkNumberNotFound struct {
	S3Error
}

func NewChunkNumberNotFound(err error, logData map[string]interface{}) *ErrChunkNumberNotFound {
	return &ErrChunkNumberNotFound{
		S3Error: S3Error{
			err:     err,
			logData: logData,
		},
	}
}
