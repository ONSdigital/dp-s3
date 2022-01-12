package s3

import (
	"fmt"
	"github.com/ONSdigital/log.go/v2/log"
)

type DPError interface {
	error
}

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
func (e Error) Error() string {
	if e.err == nil {
		return "nil"
	}
	return e.err.Error()
}

// LogData implements the DataLogger interface which allows you extract
// embedded log.Data from an error
func (e Error) LogData() map[string]interface{} {
	if e.logData == nil {
		return log.Data{}
	}
	return e.logData
}

// Unwrap returns the wrapped error
func (e Error) Unwrap() error {
	return e.err
}

func NewUnexpectedRegionError(err error, logData map[string]interface{}) ErrUnexpectedRegion {
	return ErrUnexpectedRegion{
		Error: Error{
			err:     err,
			logData: logData,
		},
	}
}

// ErrUnexpectedRegion if a request tried to access an unexpected region
type ErrUnexpectedRegion struct {
	Error
	//Region         string
	//ExpectedRegion string
}

//// Error returns the error message with the requested and expected regions
//func (e *ErrUnexpectedRegion) Error() string {
//	return fmt.Sprintf("Unexpected region: %s. This S3 client is configured with region %s",
//		e.Region, e.ExpectedRegion)
//}

// ErrUnexpectedBucket if a request tried to access an unexpected bucket
type ErrUnexpectedBucket struct {
	BucketName         string
	ExpectedBucketName string
}

// Error returns the error message with the requested and expected bucket names.
func (e *ErrUnexpectedBucket) Error() string {
	return fmt.Sprintf("Unexpected bucket: %s. This S3 client is configured with bucket %s",
		e.BucketName, e.ExpectedBucketName)
}

// ErrBucketNotFound if the bucket configured for this client does not exist
type ErrBucketNotFound struct {
	BucketName string
}

// Error returns the error message with the bucket name
func (e *ErrBucketNotFound) Error() string {
	return fmt.Sprintf("Bucket %s not found", e.BucketName)
}

// ErrNotUploaded if an s3Key could not be found in ListMultipartUploads
type ErrNotUploaded struct {
	UploadKey string
}

// Error returns the error message with the chunk number that could not be found
func (e *ErrNotUploaded) Error() string {
	return fmt.Sprintf("%s not uploaded", e.UploadKey)
}

// ErrListParts represents an error returned by S3 ListParts
type ErrListParts struct {
	Msg string
}

// Error returns the underlaying error message returned by ListParts
func (e *ErrListParts) Error() string {
	return fmt.Sprintf("ListParts error: %s", e.Msg)
}

// ErrChunkNumberNotFound if a chunk number could not be found in an existing multipart upload.
type ErrChunkNumberNotFound struct {
	ChunkNumber int64
}

// Error returns the error message containing the chunk num er that could not be found
func (e *ErrChunkNumberNotFound) Error() string {
	return fmt.Sprintf("Chunk number %d not found", e.ChunkNumber)
}

// ErrInvalidUploader is the error returned when the user tries to execute an operation with the the wrong type of Uploader
type ErrInvalidUploader struct {
	ExpectCrypto bool
}

func (e *ErrInvalidUploader) Error() string {
	if e.ExpectCrypto {
		return fmt.Sprintf("Expected Crypto Uploader, but uploader was initialised with only AWS SDK Uploader")
	}
	return fmt.Sprintf("Expected AWS SDK Uploader, but uploader was initialised with only Crypto Uploader")
}
