package s3client

import "errors"

// ErrChunkNumberNotUploaded if upload key could not be found
var ErrChunkNumberNotUploaded = errors.New("Chunk number not uploaded")

// ErrChunkNumberNotFound if ListParts failed for a particular key, bucket and uploadID
var ErrChunkNumberNotFound = errors.New("Chunk number not found")

// ErrUnexpectedBucket if a request tried to access an unexpected bucket
var ErrUnexpectedBucket = errors.New("Unexpected bucket")
