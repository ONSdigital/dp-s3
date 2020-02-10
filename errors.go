package s3client

import (
	"fmt"
)

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
	ChunkNumber int
}

// Error returns the error message containing the chunk num er that could not be found
func (e *ErrChunkNumberNotFound) Error() string {
	return fmt.Sprintf("Chunk number %d not found", e.ChunkNumber)
}
