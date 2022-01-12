package s3_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"testing"

	dps3 "github.com/ONSdigital/dp-s3/v2"
	"github.com/ONSdigital/dp-s3/v2/mock"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go/service/s3"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGet(t *testing.T) {

	Convey("Given an S3 client configured with a bucket and region", t, func() {

		payload := []byte("test data")
		bucket := "myBucket"
		objKey := "my/object/key"
		region := "eu-north-1"
		contentLen := int64(123)

		sdkMock := &mock.S3SDKClientMock{
			GetObjectFunc: func(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{
					Body:          ioutil.NopCloser(bytes.NewReader(payload)),
					ContentLength: &contentLen,
				}, nil
			},
		}

		cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, region, nil)

		Convey("Get returns an io.Reader with the expected payload", func() {
			ret, cLen, err := cli.Get(objKey)
			So(err, ShouldBeNil)
			So(*cLen, ShouldEqual, contentLen)
			b := readBytes(ret)
			So(b, ShouldResemble, payload)
			So(len(sdkMock.GetObjectCalls()), ShouldEqual, 1)
			So(sdkMock.GetObjectCalls()[0].In, ShouldResemble, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &objKey,
			})
		})

		Convey("GetFromS3URL called with a valid global URL returns an io.Reader with the expected payload", func() {
			validGlobalURL := fmt.Sprintf("s3://%s/%s", bucket, objKey)
			ret, cLen, err := cli.GetFromS3URL(validGlobalURL, dps3.AliasVirtualHostedStyle)
			So(err, ShouldBeNil)
			So(*cLen, ShouldEqual, contentLen)
			b := readBytes(ret)
			So(b, ShouldResemble, payload)
			So(len(sdkMock.GetObjectCalls()), ShouldEqual, 1)
			So(sdkMock.GetObjectCalls()[0].In, ShouldResemble, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &objKey,
			})
		})

		Convey("GetFromS3URL called with a valid regional URL returns an io.Reader with the expected payload", func() {
			validRegionalURL := fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s", region, bucket, objKey)
			ret, cLen, err := cli.GetFromS3URL(validRegionalURL, dps3.PathStyle)
			So(err, ShouldBeNil)
			So(*cLen, ShouldEqual, contentLen)
			b := readBytes(ret)
			So(b, ShouldResemble, payload)
			So(len(sdkMock.GetObjectCalls()), ShouldEqual, 1)
			So(sdkMock.GetObjectCalls()[0].In, ShouldResemble, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &objKey,
			})
		})

		Convey("GetFromS3URL called with a valid global URL with the wrong bucket returns ErrUnexpectedBucket", func() {
			wrongBucketGlobalURL := fmt.Sprintf("s3://%s/%s", "wrongBucket", objKey)
			_, _, err := cli.GetFromS3URL(wrongBucketGlobalURL, dps3.AliasVirtualHostedStyle)
			So(err, ShouldResemble, dps3.NewUnexpectedBucketError(
				errors.New("unexpected bucket name in url"),
				log.Data{"bucket_name": bucket,
					"raw_url":   wrongBucketGlobalURL,
					"url_style": "AliasVirtualHosted",
				},
			))
			So(len(sdkMock.GetObjectCalls()), ShouldEqual, 0)
		})

		Convey("GetFromS3URL called with a valid regional URL with the wrong region returns ErrUnexpectedBucket", func() {
			wrongRegionRegionalURL := fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s", "wrongRegion", bucket, objKey)
			_, _, err := cli.GetFromS3URL(wrongRegionRegionalURL, dps3.PathStyle)
			So(err, ShouldResemble, dps3.NewUnexpectedRegionError(
				errors.New("unexpected aws region in url"),
				log.Data{"region": region,
					"raw_url":   wrongRegionRegionalURL,
					"url_style": "Path",
				},
			))
			So(len(sdkMock.GetObjectCalls()), ShouldEqual, 0)
		})

		Convey("GetFromS3URL called with a malformed URL returns error", func() {
			malformedURL := "This%Url%Is%Malformed"
			_, _, err := cli.GetFromS3URL(malformedURL, dps3.AliasVirtualHostedStyle)
			So(err, ShouldResemble, dps3.NewError(
				fmt.Errorf("error parsing url: %w",
					&url.Error{Op: "parse", URL: malformedURL, Err: url.EscapeError("%Ur")},
				),
				log.Data{
					"raw_url":   malformedURL,
					"url_style": "AliasVirtualHosted"},
			))
			So(len(sdkMock.GetObjectCalls()), ShouldEqual, 0)
		})
	})
}

func TestGetWithPSK(t *testing.T) {
	Convey("Given an S3 client configured with a bucket, region and psk", t, func() {

		psk := []byte("test psk")
		payload := []byte("test data")
		bucket := "myBucket"
		objKey := "my/object/key"
		region := "eu-north-1"
		contentLen := int64(123)

		cryptoMock := &mock.S3CryptoClientMock{
			GetObjectWithPSKFunc: func(input *s3.GetObjectInput, inPsk []byte) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{
					Body:          ioutil.NopCloser(bytes.NewReader(payload)),
					ContentLength: &contentLen,
				}, nil
			},
		}

		cli := dps3.InstantiateClient(nil, cryptoMock, nil, nil, bucket, region, nil)

		Convey("GetWithPSK returns an io.Reader with the expected payload", func() {
			ret, cLen, err := cli.GetWithPSK(objKey, psk)
			So(err, ShouldBeNil)
			So(*cLen, ShouldEqual, contentLen)
			b := readBytes(ret)
			So(b, ShouldResemble, payload)
			So(len(cryptoMock.GetObjectWithPSKCalls()), ShouldEqual, 1)
			So(cryptoMock.GetObjectWithPSKCalls()[0].Psk, ShouldResemble, psk)
			So(cryptoMock.GetObjectWithPSKCalls()[0].In, ShouldResemble, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &objKey,
			})
		})

		Convey("GetFromS3URLWithPSK called with a valid global URL returns an io.Reader with the expected payload", func() {
			validGlobalURL := fmt.Sprintf("s3://%s/%s", bucket, objKey)
			ret, cLen, err := cli.GetFromS3URLWithPSK(validGlobalURL, dps3.AliasVirtualHostedStyle, psk)
			So(err, ShouldBeNil)
			So(*cLen, ShouldEqual, contentLen)
			b := readBytes(ret)
			So(b, ShouldResemble, payload)
			So(len(cryptoMock.GetObjectWithPSKCalls()), ShouldEqual, 1)
			So(cryptoMock.GetObjectWithPSKCalls()[0].Psk, ShouldResemble, psk)
			So(cryptoMock.GetObjectWithPSKCalls()[0].In, ShouldResemble, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &objKey,
			})
		})
	})
}

func TestHead(t *testing.T) {
	bucket := "myBucket"
	region := "eu-north-1"
	contentLen := int64(123)
	objKey := "my/object/key"

	Convey("Given an S3 client that returns a valid HeadObject response", t, func() {
		sdkMock := &mock.S3SDKClientMock{
			HeadObjectFunc: func(in *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{
					ContentLength: &contentLen,
				}, nil
			},
		}
		cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, region, nil)

		Convey("Head returns the expected output returned by the sdk client without error", func() {
			out, err := cli.Head(objKey)
			So(err, ShouldBeNil)
			So(*out.ContentLength, ShouldEqual, contentLen)
		})
	})

	Convey("Given an S3 client that returns an error on a HeadObject request", t, func() {
		errHead := errors.New("headObject error")
		sdkMock := &mock.S3SDKClientMock{
			HeadObjectFunc: func(in *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
				return nil, errHead
			},
		}
		s3Cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, region, nil)

		Convey("Head returns the expected error", func() {
			_, err := s3Cli.Head(objKey)
			So(err, ShouldResemble, dps3.NewError(
				fmt.Errorf("error trying to obtain s3 object metadata with HeadObject call: %w", errHead),
				log.Data{
					"bucket_name": bucket,
					"s3_key":      objKey,
				},
			))
		})
	})
}
