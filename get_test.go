package s3_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"

	dps3 "github.com/ONSdigital/dp-s3/v3"
	"github.com/ONSdigital/dp-s3/v3/mock"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGet(t *testing.T) {
	Convey("Given an S3 client configured with a bucket and region", t, func() {
		ctx := context.Background()

		payload := []byte("test data")
		bucket := "myBucket"
		objKey := "my/object/key"
		region := "eu-north-1"
		contentLen := int64(123)

		sdkMock := &mock.S3SDKClientMock{
			GetObjectFunc: func(ctx context.Context, input *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{
					Body:          io.NopCloser(bytes.NewReader(payload)),
					ContentLength: &contentLen,
				}, nil
			},
		}

		cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, region, aws.Config{})

		Convey("Get returns an io.Reader with the expected payload", func() {
			ret, cLen, err := cli.Get(ctx, objKey)
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
			ret, cLen, err := cli.GetFromS3URL(ctx, validGlobalURL, dps3.AliasVirtualHostedStyle)
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
			ret, cLen, err := cli.GetFromS3URL(ctx, validRegionalURL, dps3.PathStyle)
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
			_, _, err := cli.GetFromS3URL(ctx, wrongBucketGlobalURL, dps3.AliasVirtualHostedStyle)
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
			_, _, err := cli.GetFromS3URL(ctx, wrongRegionRegionalURL, dps3.PathStyle)
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
			_, _, err := cli.GetFromS3URL(ctx, malformedURL, dps3.AliasVirtualHostedStyle)
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
		ctx := context.Background()

		psk := []byte("test psk")
		payload := []byte("test data")
		bucket := "myBucket"
		objKey := "my/object/key"
		region := "eu-north-1"
		contentLen := int64(123)

		cryptoMock := &mock.S3CryptoClientMock{
			GetObjectWithPSKFunc: func(ctx context.Context, input *s3.GetObjectInput, inPsk []byte) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{
					Body:          io.NopCloser(bytes.NewReader(payload)),
					ContentLength: &contentLen,
				}, nil
			},
		}

		cli := dps3.InstantiateClient(nil, cryptoMock, nil, nil, bucket, region, aws.Config{})

		Convey("GetWithPSK returns an io.Reader with the expected payload", func() {
			ret, cLen, err := cli.GetWithPSK(ctx, objKey, psk)
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
			ret, cLen, err := cli.GetFromS3URLWithPSK(ctx, validGlobalURL, dps3.AliasVirtualHostedStyle, psk)
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

func TestFileExists(t *testing.T) {
	ctx := context.Background()

	bucket := "myBucket"
	region := "eu-north-1"
	contentLen := int64(123)
	objKey := "my/object/key"

	Convey("Given an S3 client that returns a valid HeadObject response", t, func() {
		sdkMock := &mock.S3SDKClientMock{
			HeadObjectFunc: func(ctx context.Context, in *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{
					ContentLength: &contentLen,
				}, nil
			},
		}
		cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, region, aws.Config{})

		Convey("When the file exists", func() {
			exists, err := cli.FileExists(ctx, objKey)
			So(err, ShouldBeNil)
			So(exists, ShouldBeTrue)
		})
	})

	Convey("Given an S3 client that returns a Not Found Error", t, func() {
		sdkMock := &mock.S3SDKClientMock{
			HeadObjectFunc: func(ctx context.Context, in *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
				return nil, &types.NotFound{}
			},
		}
		cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, region, aws.Config{})

		Convey("When the file exists", func() {
			exists, err := cli.FileExists(ctx, objKey)
			So(err, ShouldBeNil)
			So(exists, ShouldBeFalse)
		})
	})

	Convey("Given an S3 client that returns a unexpected AWS Error", t, func() {
		sdkMock := &mock.S3SDKClientMock{
			HeadObjectFunc: func(ctx context.Context, in *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
				return nil, &smithy.GenericAPIError{}
			},
		}
		cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, region, aws.Config{})

		Convey("When the file exists", func() {
			_, err := cli.FileExists(ctx, objKey)
			So(err, ShouldBeError)
		})
	})

	Convey("Given an S3 client that returns a unexpected Error", t, func() {
		sdkMock := &mock.S3SDKClientMock{
			HeadObjectFunc: func(ctx context.Context, in *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
				return nil, errors.New("very broken")
			},
		}
		cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, region, aws.Config{})

		Convey("When the file exists", func() {
			_, err := cli.FileExists(ctx, objKey)
			So(err, ShouldBeError)
		})
	})
}

func TestHead(t *testing.T) {
	ctx := context.Background()

	bucket := "myBucket"
	region := "eu-north-1"
	contentLen := int64(123)
	objKey := "my/object/key"

	Convey("Given an S3 client that returns a valid HeadObject response", t, func() {
		sdkMock := &mock.S3SDKClientMock{
			HeadObjectFunc: func(ctx context.Context, in *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
				return &s3.HeadObjectOutput{
					ContentLength: &contentLen,
				}, nil
			},
		}
		cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, region, aws.Config{})

		Convey("Head returns the expected output returned by the sdk client without error", func() {
			out, err := cli.Head(ctx, objKey)
			So(err, ShouldBeNil)
			So(*out.ContentLength, ShouldEqual, contentLen)
		})
	})

	Convey("Given an S3 client that returns an error on a HeadObject request", t, func() {
		errHead := errors.New("headObject error")
		sdkMock := &mock.S3SDKClientMock{
			HeadObjectFunc: func(ctx context.Context, in *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
				return nil, errHead
			},
		}
		s3Cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, region, aws.Config{})

		Convey("Head returns the expected error", func() {
			_, err := s3Cli.Head(ctx, objKey)
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

func TestGetBucketPolicy(t *testing.T) {
	ctx := context.Background()

	bucket := "myBucket"
	region := "eu-north-1"
	policy := "policy"

	expectedReturn := &s3.GetBucketPolicyOutput{
		Policy: &policy,
	}

	Convey("Given an S3 client that returns a valid BucketPolicy response", t, func() {
		sdkMock := &mock.S3SDKClientMock{
			GetBucketPolicyFunc: func(ctx context.Context, in *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
				return &s3.GetBucketPolicyOutput{
					Policy: &policy,
				}, nil
			},
		}
		cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, region, aws.Config{})

		Convey("GetBucketPolicy returns the expected output returned by the sdk client without error", func() {
			out, err := cli.GetBucketPolicy(ctx, bucket)
			So(err, ShouldBeNil)
			So(out, ShouldResemble, expectedReturn)
		})
	})

	Convey("Given an S3 client that returns an error on a BucketPolicy request", t, func() {
		errPolicy := errors.New("BucketPolicy error")
		sdkMock := &mock.S3SDKClientMock{
			GetBucketPolicyFunc: func(ctx context.Context, in *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
				return nil, errPolicy
			},
		}
		s3Cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, region, aws.Config{})

		Convey("BucketPolicy returns the expected error", func() {
			_, err := s3Cli.GetBucketPolicy(ctx, bucket)
			So(err, ShouldResemble, errPolicy)
		})
	})
	Convey("Given an S3 client that returns an aws error on a BucketPolicy request", t, func() {
		sdkMock := &mock.S3SDKClientMock{
			GetBucketPolicyFunc: func(ctx context.Context, in *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
				return nil, &types.NotFound{}
			},
		}
		s3Cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, region, aws.Config{})

		Convey("BucketPolicy returns the expected error", func() {
			_, err := s3Cli.GetBucketPolicy(ctx, bucket)
			So(err, ShouldBeNil)
		})
	})
}

func TestPutBucketPolicy(t *testing.T) {
	ctx := context.Background()

	bucket := "myBucket"
	region := "eu-north-1"
	policy := "policy"
	expectedReturn := &s3.PutBucketPolicyOutput{}

	Convey("Given an S3 client that returns a valid BucketPolicy response", t, func() {
		sdkMock := &mock.S3SDKClientMock{
			PutBucketPolicyFunc: func(ctx context.Context, in *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error) {
				return &s3.PutBucketPolicyOutput{}, nil
			},
		}
		cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, region, aws.Config{})

		Convey("putBucketPolicy returns the expected output returned by the sdk client without error", func() {
			out, err := cli.PutBucketPolicy(ctx, bucket, policy)
			So(err, ShouldBeNil)
			So(out, ShouldResemble, expectedReturn)
		})
	})

	Convey("Given an S3 client that returns an error on a BucketPolicy request", t, func() {
		sdkMock := &mock.S3SDKClientMock{
			PutBucketPolicyFunc: func(ctx context.Context, in *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error) {
				return nil, &types.NotFound{}
			},
		}
		s3Cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, region, aws.Config{})

		Convey("BucketPolicy returns the expected error", func() {
			_, err := s3Cli.PutBucketPolicy(ctx, bucket, policy)
			So(err, ShouldBeNil)
		})
	})
	Convey("Given an S3 client that returns an aws error on a BucketPolicy request", t, func() {
		sdkMock := &mock.S3SDKClientMock{
			PutBucketPolicyFunc: func(ctx context.Context, in *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error) {
				return nil, &types.NotFound{}
			},
		}
		s3Cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, region, aws.Config{})

		Convey("BucketPolicy returns the expected error", func() {
			_, err := s3Cli.PutBucketPolicy(ctx, bucket, policy)
			So(err, ShouldBeNil)

		})
	})
}

func TestListObjects(t *testing.T) {
	ctx := context.Background()

	bucket := "myBucket"
	region := "eu-north-1"
	expectedReturn := &s3.ListObjectsOutput{}

	Convey("Given an S3 client that returns a valid ListObjects response", t, func() {
		sdkMock := &mock.S3SDKClientMock{
			ListObjectsFunc: func(ctx context.Context, in *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return &s3.ListObjectsOutput{}, nil
			},
		}
		cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, region, aws.Config{})

		Convey("ListObjects returns the expected output returned by the sdk client without error", func() {
			out, err := cli.ListObjects(ctx, bucket)
			So(err, ShouldBeNil)
			So(out, ShouldResemble, expectedReturn)
		})
	})

	Convey("Given an S3 client that returns an non aws error on a ListObjects request", t, func() {
		errBucket := errors.New("NoSuchBucket")
		sdkMock := &mock.S3SDKClientMock{
			ListObjectsFunc: func(ctx context.Context, in *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return nil, errBucket
			},
		}
		s3Cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, region, aws.Config{})

		Convey("BucketPolicy returns the expected error", func() {
			_, err := s3Cli.ListObjects(ctx, bucket)
			So(err, ShouldResemble, errBucket)
		})
	})

	Convey("Given an S3 client that returns an aws error on a ListObjects request", t, func() {
		sdkMock := &mock.S3SDKClientMock{
			ListObjectsFunc: func(ctx context.Context, in *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
				return nil, &types.NotFound{}
			},
		}
		s3Cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, region, aws.Config{})

		Convey("BucketPolicy returns the expected error", func() {
			_, err := s3Cli.ListObjects(ctx, bucket)
			var notFoundErr *types.NotFound
			So(errors.As(err, &notFoundErr), ShouldBeTrue)
		})
	})
}
