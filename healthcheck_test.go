package s3client_test

import (
	"context"
	"errors"
	"testing"

	health "github.com/ONSdigital/dp-healthcheck/v2/healthcheck"
	s3client "github.com/ONSdigital/dp-s3/v2"
	"github.com/ONSdigital/dp-s3/v2/mock"
	"github.com/aws/aws-sdk-go/service/s3"
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
var msgWrongRegion = "BucketRegionError"

// msgInexistentRegion is the message returned when we try to get a bucket with an inexistent region
var msgInexistentRegion = "RequestError"

// msgBucketNotFound is the message returned when we try to get a bucket that does not exist
var msgBucketNotFound = "BucketNotFound"

// bucketExists is the mock function for requests for existing buckets
func bucketExists(*s3.HeadBucketInput) (*s3.HeadBucketOutput, error) {
	return &s3.HeadBucketOutput{}, nil
}

// bucketDoesNotExist is the mock function for requests with inexistent regions
func bucketDoesNotExist(*s3.HeadBucketInput) (*s3.HeadBucketOutput, error) {
	return &s3.HeadBucketOutput{}, errors.New(msgBucketNotFound)
}

// bucketWrongRegion is the mock function for requests with wrong region for bucket
func bucketWrongRegion(*s3.HeadBucketInput) (*s3.HeadBucketOutput, error) {
	return &s3.HeadBucketOutput{}, errors.New(msgWrongRegion)
}

// bucketInexistentRegion is the mock function for requests with inexistent region
func bucketInexistentRegion(*s3.HeadBucketInput) (*s3.HeadBucketOutput, error) {
	return &s3.HeadBucketOutput{}, errors.New(msgInexistentRegion)
}

func TestBucketOk(t *testing.T) {
	Convey("Given that S3 client is available, bucket exists and it was created in the same region as the S3 client config", t, func() {

		// Create S3Client with SDK Mock for existing bucket
		sdkMock := &mock.S3SDKClientMock{
			HeadBucketFunc: bucketExists,
		}
		s3Cli := s3client.InstantiateClient(sdkMock, nil, ExistingBucket, ExpectedRegion, nil)

		// CheckState for test validation
		checkState := health.NewCheckState(s3client.ServiceName)

		Convey("Checker updates the CheckState to a successful state", func() {
			s3Cli.Checker(context.Background(), checkState)
			So(len(sdkMock.HeadBucketCalls()), ShouldEqual, 1)
			So(checkState.Status(), ShouldEqual, health.StatusOK)
			So(checkState.Message(), ShouldEqual, s3client.MsgHealthy)
			So(checkState.StatusCode(), ShouldEqual, 0)
		})
	})
}

func TestBucketDoesNotExist(t *testing.T) {
	Convey("Given that S3 client is available and bucket does not exist", t, func() {

		// Create S3Client with SDK Mock for inexistent bucket
		sdkMock := &mock.S3SDKClientMock{
			HeadBucketFunc: bucketDoesNotExist,
		}
		s3Cli := s3client.InstantiateClient(sdkMock, nil, InexistentBucket, ExpectedRegion, nil)

		// CheckState for test validation
		checkState := health.NewCheckState(s3client.ServiceName)

		Convey("Checker updates the CheckState to a critical state with the relevant error message", func() {
			s3Cli.Checker(context.Background(), checkState)
			So(len(sdkMock.HeadBucketCalls()), ShouldEqual, 1)
			So(checkState.Status(), ShouldEqual, health.StatusCritical)
			So(checkState.Message(), ShouldEqual, msgBucketNotFound)
			So(checkState.StatusCode(), ShouldEqual, 0)
		})
	})
}

func TestBucketUnexpectedRegion(t *testing.T) {
	Convey("Given that S3 client is available and bucket was created in a different region than the S3 client config", t, func() {

		// Create S3Client with SDK Mock for unexpected region for bucket
		sdkMock := &mock.S3SDKClientMock{
			HeadBucketFunc: bucketWrongRegion,
		}
		s3Cli := s3client.InstantiateClient(sdkMock, nil, ExistingBucket, UnexpectedRegion, nil)

		// CheckState for test validation
		checkState := health.NewCheckState(s3client.ServiceName)

		Convey("Checker updates the CheckState to a critical state with the relevant error message", func() {
			s3Cli.Checker(context.Background(), checkState)
			So(len(sdkMock.HeadBucketCalls()), ShouldEqual, 1)
			So(checkState.Status(), ShouldEqual, health.StatusCritical)
			So(checkState.Message(), ShouldEqual, msgWrongRegion)
			So(checkState.StatusCode(), ShouldEqual, 0)
		})
	})
}

func TestBucketInexistentRegion(t *testing.T) {
	Convey("Given that S3 client is available, bucket exists, but S3 is configured with an inexistent region", t, func() {

		// Create S3Client with SDK Mock for inexistent region
		sdkMock := &mock.S3SDKClientMock{
			HeadBucketFunc: bucketInexistentRegion,
		}
		s3Cli := s3client.InstantiateClient(sdkMock, nil, ExistingBucket, InexistentRegion, nil)

		// CheckState for test validation
		checkState := health.NewCheckState(s3client.ServiceName)

		Convey("Checker updates the CheckState to a critical state with the relevant error message", func() {
			s3Cli.Checker(context.Background(), checkState)
			So(len(sdkMock.HeadBucketCalls()), ShouldEqual, 1)
			expectedErr := errors.New(msgInexistentRegion)
			So(checkState.Status(), ShouldEqual, health.StatusCritical)
			So(checkState.Message(), ShouldEqual, expectedErr.Error())
			So(checkState.StatusCode(), ShouldEqual, 0)
		})
	})
}
