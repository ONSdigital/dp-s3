package s3client

import (
	"io"
	"time"

	"github.com/goamz/goamz/aws"
	goamzs3 "github.com/goamz/goamz/s3"
)

// S3 provides AWS S3 functions that support fully qualified URL's using s3 client from goamz s3 package, which implements AmzClient interface
type S3 struct {
	cli        AmzClient
	BucketName string
}

// New returns a new AWS specific file.Provider instance configured for the given region.
func New(region string, bucketName string) (*S3, error) {

	// AWS credentials gathered from the env.
	auth, err := aws.GetAuth("", "", "", time.Time{})
	if err != nil {
		return nil, err
	}

	// Initialise amazon client with internal client.
	amzCli := &AmzClientImpl{
		goamzs3.New(auth, aws.Regions[region]),
	}

	// Create S3 with the created amzClient.
	return NewWithClient(amzCli, bucketName), nil
}

// NewWithClient returns a new S3 structure for the provided AmzClient instance.
func NewWithClient(client AmzClient, bucketName string) *S3 {
	return &S3{
		cli:        client,
		BucketName: bucketName,
	}
}

// Get returns an io.ReadCloser instance for the given fully qualified S3 URL.
// Note that this function will mutate the bucket name of the S3 object if it's different than the existing one.
func (s3 *S3) Get(rawURL string) (io.ReadCloser, error) {

	// Use the S3 URL implementation as the S3 drivers don't seem to handle fully qualified URLs that include the
	// bucket name.
	url, err := NewURL(rawURL)
	if err != nil {
		return nil, err
	}

	// Get the Bucket Reader.
	reader, err := s3.cli.GetBucketReader(url.BucketName(), url.Path())
	if err != nil {
		return nil, err
	}

	// If getting the reader was successful, set the bucket name.
	s3.BucketName = url.BucketName()
	return reader, nil
}
