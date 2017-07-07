package s3

import (
	"io"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"time"
)

// New returns a new AWS specific file.Provider instance configured for the given region.
func New(awsRegion string) (*S3, error) {

	// AWS credentials gathered from the env.
	auth, err := aws.GetAuth("", "", "", time.Time{})
	if err != nil {
		return nil, err
	}

	s3 := s3.New(auth, aws.Regions[awsRegion])

	return &S3{
		s3: s3,
	}, nil
}

// S3 provides AWS S3 functions that support fully qualified URL's.
type S3 struct {
	s3 *s3.S3
}

// Get a io.ReadCloser instance for the given fully qualified S3 URL.
func (s3 *S3) Get(rawURL string) (io.ReadCloser, error) {

	// Use the S3 URL implementation as the S3 drivers don't seem to handle fully qualified URLs that include the
	// bucket name.
	url, err := NewURL(rawURL)
	if err != nil {
		return nil, err
	}

	bucket := s3.s3.Bucket(url.BucketName())
	reader, error := bucket.GetReader(url.Path())
	if err != nil {
		return nil, error
	}

	return reader, nil
}