package s3client

import (
	"io"

	goamzs3 "github.com/goamz/goamz/s3"
)

//go:generate moq -out ./mock/goamz-s3.go -pkg mock . AmzClient

// AmzClient is an interface representing an s3 client wrapper for goamz s3 package
type AmzClient interface {
	GetBucketReader(bucketName string, path string) (io.ReadCloser, error)
}

// AmzClientImpl implements the AmzClient interface wrapping the real S3 calls to nested structs (e.g Bucket)
type AmzClientImpl struct {
	s3 *goamzs3.S3
}

// GetBucketReader returns io.ReadCloser for a bucket name and path, using the Bucket struct in goamz
func (amzCli *AmzClientImpl) GetBucketReader(bucketName string, path string) (io.ReadCloser, error) {
	return amzCli.s3.Bucket(bucketName).GetReader(path)
}
