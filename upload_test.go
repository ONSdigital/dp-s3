package s3_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	dps3 "github.com/ONSdigital/dp-s3/v2"
	"github.com/ONSdigital/dp-s3/v2/mock"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	testS3Key  = "test/s3Key.csv"
	testBucket = "my-bucket"
)

func TestUpload(t *testing.T) {
	Convey("Given a client configured with a successful SDK uploader", t, func() {
		ctx := context.Background()

		sdkUploaderMock := &mock.S3SDKUploaderMock{
			UploadFunc: func(ctx context.Context, in *s3.PutObjectInput, options ...func(*manager.Uploader)) (*manager.UploadOutput, error) {
				return nil, nil
			},
		}
		cli := dps3.InstantiateClient(nil, nil, sdkUploaderMock, nil, testBucket, ExpectedRegion, aws.Config{})

		Convey("Calling Upload with a valid s3 key results in sdk Upload being called as expected", func() {
			_, err := cli.Upload(ctx, &s3.PutObjectInput{Key: &testS3Key})
			So(err, ShouldBeNil)
			So(len(sdkUploaderMock.UploadCalls()), ShouldEqual, 1)
			So(*sdkUploaderMock.UploadCalls()[0].In.Bucket, ShouldEqual, testBucket)
		})

		Convey("Calling Upload with nil input returns the expected error", func() {
			_, err := cli.Upload(ctx, nil)
			So(err, ShouldResemble, dps3.NewError(
				fmt.Errorf("validation error for Upload: %w",
					errors.New("nil input provided"),
				),
				log.Data{
					"bucket_name": testBucket,
				},
			))
			So(len(sdkUploaderMock.UploadCalls()), ShouldEqual, 0)
		})
	})

	Convey("Given a client configured with an SDK uploader that fails to upload", t, func() {
		ctx := context.Background()

		errUploader := errors.New("failed to upload file")
		sdkUploaderMock := &mock.S3SDKUploaderMock{
			UploadFunc: func(ctx context.Context, in *s3.PutObjectInput, options ...func(*manager.Uploader)) (*manager.UploadOutput, error) {
				return nil, errUploader
			},
		}
		cli := dps3.InstantiateClient(nil, nil, sdkUploaderMock, nil, testBucket, ExpectedRegion, aws.Config{})

		Convey("Calling Upload with a valid s3 key results in the expected error being returned", func() {
			_, err := cli.Upload(ctx, &s3.PutObjectInput{Key: &testS3Key})
			So(err, ShouldResemble, dps3.NewError(
				fmt.Errorf("failed to upload: %w", errUploader),
				log.Data{
					"bucket_name": testBucket,
					"s3_key":      testS3Key,
				},
			))
			So(len(sdkUploaderMock.UploadCalls()), ShouldEqual, 1)
			So(*sdkUploaderMock.UploadCalls()[0].In.Bucket, ShouldEqual, testBucket)
		})
	})
}

func TestUploadWithPSK(t *testing.T) {
	Convey("Given a client configured with a user-defined psk uploader", t, func() {
		ctx := context.Background()

		psk := []byte("test psk")
		cryptoUploaderMock := &mock.S3CryptoUploaderMock{
			UploadWithPSKFunc: func(ctx context.Context, in *s3.PutObjectInput, psk []byte) (*manager.UploadOutput, error) {
				return &manager.UploadOutput{}, nil
			},
		}
		cli := dps3.InstantiateClient(nil, nil, nil, cryptoUploaderMock, testBucket, ExpectedRegion, aws.Config{})

		Convey("Calling UploadWithPSK with a valid s3 key results in crypto UploadWithPSK being called as expected", func() {
			_, err := cli.UploadWithPSK(ctx, &s3.PutObjectInput{Key: &testS3Key}, psk)
			So(err, ShouldBeNil)
			So(len(cryptoUploaderMock.UploadWithPSKCalls()), ShouldEqual, 1)
			So(*cryptoUploaderMock.UploadWithPSKCalls()[0].In.Bucket, ShouldEqual, testBucket)
			So(cryptoUploaderMock.UploadWithPSKCalls()[0].Ctx, ShouldNotBeNil)
		})

		Convey("Calling UploadWithPSK with nil psk returns the expected error", func() {
			_, err := cli.UploadWithPSK(ctx, &s3.PutObjectInput{Key: &testS3Key}, nil)
			So(err, ShouldResemble, dps3.NewError(
				errors.New("nil or empty psk provided to UploadWithPSK"),
				log.Data{
					"bucket_name": testBucket,
					"s3_key":      testS3Key,
				},
			))
			So(len(cryptoUploaderMock.UploadWithPSKCalls()), ShouldEqual, 0)
		})

		Convey("Calling UploadWithPSK with nil input returns the expected error", func() {
			_, err := cli.UploadWithPSK(ctx, nil, psk)
			So(err, ShouldResemble, dps3.NewError(
				fmt.Errorf("validation error for UploadWithPSK: %w",
					errors.New("nil input provided"),
				),
				log.Data{
					"bucket_name": testBucket,
				},
			))
			So(len(cryptoUploaderMock.UploadWithPSKCalls()), ShouldEqual, 0)
		})
	})

	Convey("Given a client configured with a user-defined psk uploader that fails to upload", t, func() {
		ctx := context.Background()

		errCryptoUploader := errors.New("failed to upload file")
		psk := []byte("test psk")
		cryptoUploaderMock := &mock.S3CryptoUploaderMock{
			UploadWithPSKFunc: func(ctx context.Context, in *s3.PutObjectInput, psk []byte) (*manager.UploadOutput, error) {
				return nil, errCryptoUploader
			},
		}
		cli := dps3.InstantiateClient(nil, nil, nil, cryptoUploaderMock, testBucket, ExpectedRegion, aws.Config{})

		Convey("Calling UploadWithPSK with a valid s3 key and context results in the expected error being returned", func() {
			_, err := cli.UploadWithPSK(ctx, &s3.PutObjectInput{Key: &testS3Key}, psk)
			So(err, ShouldResemble, dps3.NewError(
				fmt.Errorf("failed to upload with psk: %w", errCryptoUploader),
				log.Data{
					"bucket_name": testBucket,
					"s3_key":      testS3Key,
				},
			))
			So(len(cryptoUploaderMock.UploadWithPSKCalls()), ShouldEqual, 1)
			So(*cryptoUploaderMock.UploadWithPSKCalls()[0].In.Bucket, ShouldEqual, testBucket)
			So(cryptoUploaderMock.UploadWithPSKCalls()[0].Ctx, ShouldNotBeNil)
		})
	})
}

func TestValidateInput(t *testing.T) {
	cli := dps3.InstantiateClient(nil, nil, nil, nil, testBucket, "", aws.Config{})

	Convey("validating an input with only an s3 key is successful with the expected logData", t, func() {
		logData, err := cli.ValidateUploadInput(&s3.PutObjectInput{
			Key: &testS3Key,
		})

		So(err, ShouldBeNil)
		So(logData, ShouldResemble, log.Data{
			"bucket_name": testBucket,
			"s3_key":      testS3Key,
		})
	})

	Convey("validating an input with an s3 key and the expected bucket is successful with the expected logData", t, func() {
		logData, err := cli.ValidateUploadInput(&s3.PutObjectInput{
			Key:    &testS3Key,
			Bucket: &testBucket,
		})

		So(err, ShouldBeNil)
		So(logData, ShouldResemble, log.Data{
			"bucket_name": testBucket,
			"s3_key":      testS3Key,
		})
	})

	Convey("validating a nil input fails with the expected error and logData", t, func() {
		logData, err := cli.ValidateUploadInput(nil)

		So(err, ShouldResemble, errors.New("nil input provided"))
		So(logData, ShouldResemble, log.Data{
			"bucket_name": testBucket,
		})
	})

	Convey("validating an input without s3 key fails with the expected error and logData", t, func() {
		logData, err := cli.ValidateUploadInput(&s3.PutObjectInput{})

		So(err, ShouldResemble, errors.New("nil or empty s3 key provided in input"))
		So(logData, ShouldResemble, log.Data{
			"bucket_name": testBucket,
		})
	})

	Convey("validating an input without s3 key fails with the expected error and logData", t, func() {
		emptyKey := ""
		logData, err := cli.ValidateUploadInput(&s3.PutObjectInput{
			Key: &emptyKey,
		})
		So(err, ShouldResemble, errors.New("nil or empty s3 key provided in input"))
		So(logData, ShouldResemble, log.Data{
			"bucket_name": testBucket,
		})
	})

	Convey("validating an input with an s3 key and an unexpected bucket fails with the expected error and logData", t, func() {
		otherBucket := "otherBucket"
		logData, err := cli.ValidateUploadInput(&s3.PutObjectInput{
			Key:    &testS3Key,
			Bucket: &otherBucket,
		})

		So(err, ShouldResemble, errors.New("unexpected bucket name provided in upload input"))
		So(logData, ShouldResemble, log.Data{
			"bucket_name":       testBucket,
			"input_bucket_name": otherBucket,
			"s3_key":            testS3Key,
		})
	})
}
