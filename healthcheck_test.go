package s3client_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	s3client "github.com/ONSdigital/dp-s3"
	"github.com/ONSdigital/dp-s3/mock"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	ExistingBucket   = "csv-exported"
	InexistentBucket = "thisBucketDoesNotExist"
	ExpectedRegion   = "eu-west-1"
	UnexpectedRegion = "us-west-1"
	InexistentRegion = "atlantida-north-1"
)

// msgWrongRegion is the message returned when we try to get a bucket with the wrong region
func msgWrongRegion(region, bucketName string) string {
	return fmt.Sprintf("Get https://s3-%s.amazonaws.com/%s/: 301 response missing Location header", region, bucketName)
}

// msgInexistentRegion is the message returned when we try to get a bucket with an inexistent region
func msgInexistentRegion(bucketName string) string {
	return fmt.Sprintf("Get /%s/: unsupported protocol scheme \"\"", bucketName)
}

// bucketExists is the mock function for requests for existing buckets
func bucketExists(bucketName, path string) (io.ReadCloser, error) {
	return ioutil.NopCloser(strings.NewReader("MockReadCloser")), nil
}

// bucketDoesNotExist is the mock function for requests with inexistent regions
func bucketDoesNotExist(bucketName, path string) (io.ReadCloser, error) {
	errBucket := s3client.ErrBucketDoesNotExist{BucketName: bucketName}
	return nil, &errBucket
}

// bucketWrongRegion is the mock function for requests with wrong region for bucket
func bucketWrongRegion(bucketName, path string) (io.ReadCloser, error) {
	return nil, errors.New(msgWrongRegion(UnexpectedRegion, bucketName))
}

// bucketInexistentRegion is the mock function for requests with inexistent region
func bucketInexistentRegion(bucketName, path string) (io.ReadCloser, error) {
	return nil, errors.New(msgInexistentRegion(bucketName))
}

func TestBucketOk(t *testing.T) {
	Convey("Given that S3 client is available, bucket exists and it was created in the same region as the S3 client config", t, func() {

		// Create S3Client with mock AmzClient
		s3AmzCliMock := &mock.AmzClientMock{
			GetBucketReaderFunc: bucketExists,
		}
		s3Cli := s3client.NewWithClient(s3AmzCliMock, ExistingBucket)

		// mock CheckState for test validation
		mockCheckState := mock.CheckStateMock{
			UpdateFunc: func(status, message string, statusCode int) error {
				return nil
			},
		}

		Convey("Checker updates the CheckState to a successful state", func() {
			s3Cli.Checker(context.Background(), &mockCheckState)
			So(len(s3AmzCliMock.GetBucketReaderCalls()), ShouldEqual, 1)
			updateCalls := mockCheckState.UpdateCalls()
			So(len(updateCalls), ShouldEqual, 1)
			So(updateCalls[0].Status, ShouldEqual, health.StatusOK)
			So(updateCalls[0].Message, ShouldEqual, s3client.MsgHealthy)
			So(updateCalls[0].StatusCode, ShouldEqual, 0)
		})
	})
}

func TestBucketDoesNotExist(t *testing.T) {
	Convey("Given that S3 client is available and bucket does not exist", t, func() {

		// Create S3Client with mock AmzClient
		s3AmzCliMock := &mock.AmzClientMock{
			GetBucketReaderFunc: bucketDoesNotExist,
		}
		s3Cli := s3client.NewWithClient(s3AmzCliMock, InexistentBucket)

		// mock CheckState for test validation
		mockCheckState := mock.CheckStateMock{
			UpdateFunc: func(status, message string, statusCode int) error {
				return nil
			},
		}

		Convey("Checker updates the CheckState to a critical state with the relevant error message", func() {
			s3Cli.Checker(context.Background(), &mockCheckState)
			So(len(s3AmzCliMock.GetBucketReaderCalls()), ShouldEqual, 1)
			updateCalls := mockCheckState.UpdateCalls()
			expectedErr := s3client.ErrBucketDoesNotExist{BucketName: InexistentBucket}
			So(len(updateCalls), ShouldEqual, 1)
			So(updateCalls[0].Status, ShouldEqual, health.StatusCritical)
			So(updateCalls[0].Message, ShouldEqual, expectedErr.Error())
			So(updateCalls[0].StatusCode, ShouldEqual, 0)
		})
	})
}

func TestBucketUnexpectedRegion(t *testing.T) {
	Convey("Given that S3 client is available and bucket was created in a different region than the S3 client config", t, func() {

		// Create S3Client with mock AmzClient
		s3AmzCliMock := &mock.AmzClientMock{
			GetBucketReaderFunc: bucketWrongRegion,
		}
		s3Cli := s3client.NewWithClient(s3AmzCliMock, ExistingBucket)

		// mock CheckState for test validation
		mockCheckState := mock.CheckStateMock{
			UpdateFunc: func(status, message string, statusCode int) error {
				return nil
			},
		}

		Convey("Checker updates the CheckState to a critical state with the relevant error message", func() {
			s3Cli.Checker(context.Background(), &mockCheckState)
			So(len(s3AmzCliMock.GetBucketReaderCalls()), ShouldEqual, 1)
			updateCalls := mockCheckState.UpdateCalls()
			expectedErr := errors.New(msgWrongRegion(UnexpectedRegion, ExistingBucket))
			So(len(updateCalls), ShouldEqual, 1)
			So(updateCalls[0].Status, ShouldEqual, health.StatusCritical)
			So(updateCalls[0].Message, ShouldEqual, expectedErr.Error())
			So(updateCalls[0].StatusCode, ShouldEqual, 0)
		})
	})
}

func TestBucketInexistentRegion(t *testing.T) {
	Convey("Given that S3 client is available, bucket exists, but S3 is configured with an inexistent region", t, func() {

		// Create S3Client with mock AmzClient
		s3AmzCliMock := &mock.AmzClientMock{
			GetBucketReaderFunc: bucketInexistentRegion,
		}
		s3Cli := s3client.NewWithClient(s3AmzCliMock, ExistingBucket)

		// mock CheckState for test validation
		mockCheckState := mock.CheckStateMock{
			UpdateFunc: func(status, message string, statusCode int) error {
				return nil
			},
		}

		Convey("Checker updates the CheckState to a critical state with the relevant error message", func() {
			s3Cli.Checker(context.Background(), &mockCheckState)
			So(len(s3AmzCliMock.GetBucketReaderCalls()), ShouldEqual, 1)
			updateCalls := mockCheckState.UpdateCalls()
			expectedErr := errors.New(msgInexistentRegion(ExistingBucket))
			So(len(updateCalls), ShouldEqual, 1)
			So(updateCalls[0].Status, ShouldEqual, health.StatusCritical)
			So(updateCalls[0].Message, ShouldEqual, expectedErr.Error())
			So(updateCalls[0].StatusCode, ShouldEqual, 0)
		})
	})
}
