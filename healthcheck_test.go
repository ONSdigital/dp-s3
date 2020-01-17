package s3client_test

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"

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
	return nil, s3client.ErrBucketDoesNotExist
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

		s3AmzCliMock := &mock.AmzClientMock{
			GetBucketReaderFunc: bucketExists,
		}
		s3Cli := s3client.NewWithClient(s3AmzCliMock)

		Convey("Checker returns a successful Check structure", func() {
			validateSuccessfulCheck(s3Cli, ExistingBucket)
			So(len(s3AmzCliMock.GetBucketReaderCalls()), ShouldEqual, 1)
		})
	})
}

func TestBucketDoesNotExist(t *testing.T) {
	Convey("Given that S3 client is available and bucket does not exist", t, func() {

		s3AmzCliMock := &mock.AmzClientMock{
			GetBucketReaderFunc: bucketDoesNotExist,
		}
		s3Cli := s3client.NewWithClient(s3AmzCliMock)

		Convey("Checker returns a critical Check structure with the relevant error message", func() {
			_, err := validateCriticalCheck(s3Cli, InexistentBucket, 500, s3client.ErrBucketDoesNotExist.Error())
			So(err.Error(), ShouldEqual, s3client.ErrBucketDoesNotExist.Error())
			So(len(s3AmzCliMock.GetBucketReaderCalls()), ShouldEqual, 1)
		})
	})
}

func TestBucketUnexpectedRegion(t *testing.T) {
	Convey("Given that S3 client is available and bucket was created in a different region than the S3 client config", t, func() {

		s3AmzCliMock := &mock.AmzClientMock{
			GetBucketReaderFunc: bucketWrongRegion,
		}
		s3Cli := s3client.NewWithClient(s3AmzCliMock)

		Convey("Checker returns a critical Check structure with the relevant error message", func() {
			msg := msgWrongRegion(UnexpectedRegion, ExistingBucket)
			_, err := validateCriticalCheck(s3Cli, ExistingBucket, 500, msg)
			So(err.Error(), ShouldEqual, msg)
			So(len(s3AmzCliMock.GetBucketReaderCalls()), ShouldEqual, 1)
		})
	})
}

func TestBucketInexistentRegion(t *testing.T) {
	Convey("Given that S3 client is available, bucket exists, but S3 is configured with an inexistent region", t, func() {

		s3AmzCliMock := &mock.AmzClientMock{
			GetBucketReaderFunc: bucketInexistentRegion,
		}
		s3Cli := s3client.NewWithClient(s3AmzCliMock)

		Convey("Checker returns a critical Check structure with the relevant error message", func() {
			msg := msgInexistentRegion(ExistingBucket)
			_, err := validateCriticalCheck(s3Cli, ExistingBucket, 500, msg)
			So(err.Error(), ShouldEqual, msg)
			So(len(s3AmzCliMock.GetBucketReaderCalls()), ShouldEqual, 1)
		})
	})
}

func validateSuccessfulCheck(cli *s3client.S3, bucketName string) (check *health.Check) {
	t0 := time.Now().UTC()
	check, err := cli.Checker(nil, bucketName)
	t1 := time.Now().UTC()
	So(err, ShouldBeNil)
	So(check.Name, ShouldEqual, s3client.ServiceName)
	So(check.Status, ShouldEqual, health.StatusOK)
	So(check.StatusCode, ShouldEqual, 200)
	So(check.Message, ShouldEqual, s3client.MsgHealthy)
	So(check.LastChecked, ShouldHappenOnOrBetween, t0, t1)
	So(check.LastSuccess, ShouldHappenOnOrBetween, t0, t1)
	So(check.LastFailure, ShouldHappenBefore, t0)
	return check
}

func validateWarningCheck(cli *s3client.S3, bucketName string, expectedCode int, expectedMessage string) (check *health.Check, err error) {
	t0 := time.Now().UTC()
	check, err = cli.Checker(nil, bucketName)
	t1 := time.Now().UTC()
	So(check.Name, ShouldEqual, s3client.ServiceName)
	So(check.Status, ShouldEqual, health.StatusWarning)
	So(check.StatusCode, ShouldEqual, expectedCode)
	So(check.Message, ShouldEqual, expectedMessage)
	So(check.LastChecked, ShouldHappenOnOrBetween, t0, t1)
	So(check.LastSuccess, ShouldHappenBefore, t0)
	So(check.LastFailure, ShouldHappenOnOrBetween, t0, t1)
	return check, err
}

func validateCriticalCheck(cli *s3client.S3, bucketName string, expectedCode int, expectedMessage string) (check *health.Check, err error) {
	t0 := time.Now().UTC()
	check, err = cli.Checker(nil, bucketName)
	t1 := time.Now().UTC()
	So(check.Name, ShouldEqual, s3client.ServiceName)
	So(check.Status, ShouldEqual, health.StatusCritical)
	So(check.StatusCode, ShouldEqual, expectedCode)
	So(check.Message, ShouldEqual, expectedMessage)
	So(check.LastChecked, ShouldHappenOnOrBetween, t0, t1)
	So(check.LastSuccess, ShouldHappenBefore, t0)
	So(check.LastFailure, ShouldHappenOnOrBetween, t0, t1)
	return check, err
}
