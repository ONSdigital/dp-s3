package s3client

import (
	"io"
	"time"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/goamz/goamz/aws"
	goamzs3 "github.com/goamz/goamz/s3"
)

// S3 provides AWS S3 functions that support fully qualified URL's using s3 client from goamz s3 package, which implements AmzClient interface
type S3 struct {
	cli   AmzClient
	Check *health.Check
}

// New returns a new AWS specific file.Provider instance configured for the given region.
func New(region string) (*S3, error) {

	// AWS credentials gathered from the env.
	auth, err := aws.GetAuth("", "", "", time.Time{})
	if err != nil {
		return nil, err
	}

	// Initialize amazon client with internal client
	amzCli := &AmzClientImpl{
		goamzs3.New(auth, aws.Regions[region]),
	}

	// ceate S3 with the created amzClient
	return NewWithClient(amzCli), nil
}

// NewWithClient returns a new S3 structure for the provided AmzClient instance.
func NewWithClient(client AmzClient) *S3 {

	// Initial Check struct
	check := &health.Check{Name: ServiceName}

	// Create S3 with amzClient and check
	return &S3{
		cli:   client,
		Check: check,
	}
}

// Get returns an io.ReadCloser instance for the given fully qualified S3 URL.
func (s3 *S3) Get(rawURL string) (io.ReadCloser, error) {

	// Use the S3 URL implementation as the S3 drivers don't seem to handle fully qualified URLs that include the
	// bucket name.
	url, err := NewURL(rawURL)
	if err != nil {
		return nil, err
	}

	reader, err := s3.cli.GetBucketReader(url.BucketName(), url.Path())
	if err != nil {
		return nil, err
	}

	return reader, nil
}
