package s3_test

import (
	"context"
	"errors"
	"testing"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	dps3 "github.com/ONSdigital/dp-s3/v3"
	"github.com/ONSdigital/dp-s3/v3/mock"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
func bucketExists(ctx context.Context, input *s3.HeadBucketInput, opts ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
	return &s3.HeadBucketOutput{}, nil
}

// bucketDoesNotExist is the mock function for requests with inexistent regions
func bucketDoesNotExist(ctx context.Context, input *s3.HeadBucketInput, opts ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
	return &s3.HeadBucketOutput{}, errors.New(msgBucketNotFound)
}

// bucketWrongRegion is the mock function for requests with wrong region for bucket
func bucketWrongRegion(ctx context.Context, input *s3.HeadBucketInput, opts ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
	return &s3.HeadBucketOutput{}, errors.New(msgWrongRegion)
}

// bucketInexistentRegion is the mock function for requests with inexistent region
func bucketInexistentRegion(ctx context.Context, input *s3.HeadBucketInput, opts ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
	return &s3.HeadBucketOutput{}, errors.New(msgInexistentRegion)
}

func TestBucketOk(t *testing.T) {
	Convey("Given that S3 client is available, bucket exists and it was created in the same region as the S3 client config", t, func() {
		// Create S3 client with SDK Mock for existing bucket
		sdkMock := &mock.S3SDKClientMock{
			HeadBucketFunc: bucketExists,
		}
		cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, ExistingBucket, ExpectedRegion, aws.Config{})

		// CheckState for test validation
		checkState := health.NewCheckState(dps3.ServiceName)

		Convey("Checker updates the CheckState to a successful state", func() {
			cli.Checker(context.Background(), checkState)
			So(len(sdkMock.HeadBucketCalls()), ShouldEqual, 1)
			So(checkState.Status(), ShouldEqual, health.StatusOK)
			So(checkState.Message(), ShouldEqual, dps3.MsgHealthy)
			So(checkState.StatusCode(), ShouldEqual, 0)
		})
	})
}

func TestBucketDoesNotExist(t *testing.T) {
	Convey("Given that S3 client is available and bucket does not exist", t, func() {
		// Create S3 client with SDK Mock for inexistent bucket
		sdkMock := &mock.S3SDKClientMock{
			HeadBucketFunc: bucketDoesNotExist,
		}
		cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, InexistentBucket, ExpectedRegion, aws.Config{})

		// CheckState for test validation
		checkState := health.NewCheckState(dps3.ServiceName)

		Convey("Checker updates the CheckState to a critical state with the relevant error message", func() {
			cli.Checker(context.Background(), checkState)
			So(len(sdkMock.HeadBucketCalls()), ShouldEqual, 1)
			So(checkState.Status(), ShouldEqual, health.StatusCritical)
			So(checkState.Message(), ShouldEqual, msgBucketNotFound)
			So(checkState.StatusCode(), ShouldEqual, 0)
		})
	})
}

func TestBucketUnexpectedRegion(t *testing.T) {
	Convey("Given that S3 client is available and bucket was created in a different region than the S3 client config", t, func() {
		// Create S3 client with SDK Mock for unexpected region for bucket
		sdkMock := &mock.S3SDKClientMock{
			HeadBucketFunc: bucketWrongRegion,
		}
		cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, ExistingBucket, UnexpectedRegion, aws.Config{})

		// CheckState for test validation
		checkState := health.NewCheckState(dps3.ServiceName)

		Convey("Checker updates the CheckState to a critical state with the relevant error message", func() {
			cli.Checker(context.Background(), checkState)
			So(len(sdkMock.HeadBucketCalls()), ShouldEqual, 1)
			So(checkState.Status(), ShouldEqual, health.StatusCritical)
			So(checkState.Message(), ShouldEqual, msgWrongRegion)
			So(checkState.StatusCode(), ShouldEqual, 0)
		})
	})
}

func TestBucketInexistentRegion(t *testing.T) {
	Convey("Given that S3 client is available, bucket exists, but S3 is configured with an inexistent region", t, func() {
		// Create S3 client with SDK Mock for inexistent region
		sdkMock := &mock.S3SDKClientMock{
			HeadBucketFunc: bucketInexistentRegion,
		}
		cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, ExistingBucket, InexistentRegion, aws.Config{})

		// CheckState for test validation
		checkState := health.NewCheckState(dps3.ServiceName)

		Convey("Checker updates the CheckState to a critical state with the relevant error message", func() {
			cli.Checker(context.Background(), checkState)
			So(len(sdkMock.HeadBucketCalls()), ShouldEqual, 1)
			expectedErr := errors.New(msgInexistentRegion)
			So(checkState.Status(), ShouldEqual, health.StatusCritical)
			So(checkState.Message(), ShouldEqual, expectedErr.Error())
			So(checkState.StatusCode(), ShouldEqual, 0)
		})
	})
}
