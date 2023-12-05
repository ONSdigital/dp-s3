package s3_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	dps3 "github.com/ONSdigital/dp-s3/v2"
	"github.com/ONSdigital/dp-s3/v2/mock"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	testS3Key  = "test/s3Key.csv"
	testBucket = "my-bucket"
)

func TestUpload(t *testing.T) {
	Convey("Given a client configured with a successful SDK uploader", t, func() {
		sdkUploaderMock := &mock.S3SDKUploaderMock{
			UploadFunc: func(in *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
				return nil, nil
			},
		}
		cli := dps3.InstantiateClient(nil, nil, sdkUploaderMock, nil, testBucket, ExpectedRegion, nil)

		Convey("Calling Upload with a valid s3 key results in sdk Upload being called as expected", func() {
			_, err := cli.Upload(&s3manager.UploadInput{Key: &testS3Key})
			So(err, ShouldBeNil)
			So(len(sdkUploaderMock.UploadCalls()), ShouldEqual, 1)
			So(*sdkUploaderMock.UploadCalls()[0].In.Bucket, ShouldEqual, testBucket)
		})

		Convey("Calling Upload with nil input returns the expected error", func() {
			_, err := cli.Upload(nil)
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
		errUploader := errors.New("failed to upload file")
		sdkUploaderMock := &mock.S3SDKUploaderMock{
			UploadFunc: func(in *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
				return nil, errUploader
			},
		}
		cli := dps3.InstantiateClient(nil, nil, sdkUploaderMock, nil, testBucket, ExpectedRegion, nil)

		Convey("Calling Upload with a valid s3 key results in the expected error being returned", func() {
			_, err := cli.Upload(&s3manager.UploadInput{Key: &testS3Key})
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

func TestUploadWithContext(t *testing.T) {
	ctx := context.Background()

	Convey("Given a client configured with a successful SDK uploader", t, func() {
		sdkUploaderMock := &mock.S3SDKUploaderMock{
			UploadWithContextFunc: func(ctx context.Context, in *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
				return nil, nil
			},
		}
		cli := dps3.InstantiateClient(nil, nil, sdkUploaderMock, nil, testBucket, ExpectedRegion, nil)

		Convey("Calling UploadWithContext with a valid s3 key and context results in sdk UploadWithContext being called as expected", func() {
			_, err := cli.UploadWithContext(ctx, &s3manager.UploadInput{Key: &testS3Key})
			So(err, ShouldBeNil)
			So(len(sdkUploaderMock.UploadWithContextCalls()), ShouldEqual, 1)
			So(*sdkUploaderMock.UploadWithContextCalls()[0].In.Bucket, ShouldEqual, testBucket)
			So(sdkUploaderMock.UploadWithContextCalls()[0].Ctx, ShouldResemble, ctx)
		})

		Convey("Calling UploadWithContext with nil context returns the expected error", func() {
			_, err := cli.UploadWithContext(nil, &s3manager.UploadInput{Key: &testS3Key})
			So(err, ShouldResemble, dps3.NewError(
				errors.New("nil context provided to UploadWithContext"),
				log.Data{
					"bucket_name": testBucket,
					"s3_key":      testS3Key,
				},
			))
			So(len(sdkUploaderMock.UploadWithContextCalls()), ShouldEqual, 0)
		})

		Convey("Calling UploadWithContext with nil input returns the expected error", func() {
			_, err := cli.UploadWithContext(ctx, nil)
			So(err, ShouldResemble, dps3.NewError(
				fmt.Errorf("validation error for UploadWithContext: %w",
					errors.New("nil input provided"),
				),
				log.Data{
					"bucket_name": testBucket,
				},
			))
			So(len(sdkUploaderMock.UploadWithContextCalls()), ShouldEqual, 0)
		})
	})

	Convey("Given a client configured with an SDK uploader that fails to upload", t, func() {
		errUploader := errors.New("failed to upload file")
		sdkUploaderMock := &mock.S3SDKUploaderMock{
			UploadWithContextFunc: func(ctx context.Context, in *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
				return nil, errUploader
			},
		}
		cli := dps3.InstantiateClient(nil, nil, sdkUploaderMock, nil, testBucket, ExpectedRegion, nil)

		Convey("Calling UploadWithContext with a valid s3 key and context results in the expected error being returned", func() {
			_, err := cli.UploadWithContext(ctx, &s3manager.UploadInput{Key: &testS3Key})
			So(err, ShouldResemble, dps3.NewError(
				fmt.Errorf("failed to upload with context: %w", errUploader),
				log.Data{
					"bucket_name": testBucket,
					"s3_key":      testS3Key,
				},
			))
			So(len(sdkUploaderMock.UploadWithContextCalls()), ShouldEqual, 1)
			So(*sdkUploaderMock.UploadWithContextCalls()[0].In.Bucket, ShouldEqual, testBucket)
			So(sdkUploaderMock.UploadWithContextCalls()[0].Ctx, ShouldResemble, ctx)
		})
	})
}

func TestUploadWithPSK(t *testing.T) {
	Convey("Given a client configured with a user-defined psk uploader", t, func() {
		psk := []byte("test psk")
		cryptoUploaderMock := &mock.S3CryptoUploaderMock{
			UploadWithPSKFunc: func(ctx context.Context, in *s3manager.UploadInput, psk []byte) (*s3manager.UploadOutput, error) {
				return &s3manager.UploadOutput{}, nil
			},
		}
		cli := dps3.InstantiateClient(nil, nil, nil, cryptoUploaderMock, testBucket, ExpectedRegion, nil)

		Convey("Calling UploadWithPSK with a valid s3 key results in crypto UploadWithPSK being called as expected", func() {
			_, err := cli.UploadWithPSK(&s3manager.UploadInput{Key: &testS3Key}, psk)
			So(err, ShouldBeNil)
			So(len(cryptoUploaderMock.UploadWithPSKCalls()), ShouldEqual, 1)
			So(*cryptoUploaderMock.UploadWithPSKCalls()[0].In.Bucket, ShouldEqual, testBucket)
			So(cryptoUploaderMock.UploadWithPSKCalls()[0].Ctx, ShouldEqual, nil)
		})

		Convey("Calling UploadWithPSK with nil psk returns the expected error", func() {
			_, err := cli.UploadWithPSK(&s3manager.UploadInput{Key: &testS3Key}, nil)
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
			_, err := cli.UploadWithPSK(nil, psk)
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
		errCryptoUploader := errors.New("failed to upload file")
		psk := []byte("test psk")
		cryptoUploaderMock := &mock.S3CryptoUploaderMock{
			UploadWithPSKFunc: func(ctx context.Context, in *s3manager.UploadInput, psk []byte) (*s3manager.UploadOutput, error) {
				return nil, errCryptoUploader
			},
		}
		cli := dps3.InstantiateClient(nil, nil, nil, cryptoUploaderMock, testBucket, ExpectedRegion, nil)

		Convey("Calling UploadWithPSK with a valid s3 key and context results in the expected error being returned", func() {
			_, err := cli.UploadWithPSK(&s3manager.UploadInput{Key: &testS3Key}, psk)
			So(err, ShouldResemble, dps3.NewError(
				fmt.Errorf("failed to upload with psk: %w", errCryptoUploader),
				log.Data{
					"bucket_name": testBucket,
					"s3_key":      testS3Key,
				},
			))
			So(len(cryptoUploaderMock.UploadWithPSKCalls()), ShouldEqual, 1)
			So(*cryptoUploaderMock.UploadWithPSKCalls()[0].In.Bucket, ShouldEqual, testBucket)
			So(cryptoUploaderMock.UploadWithPSKCalls()[0].Ctx, ShouldEqual, nil)
		})
	})
}

func TestUploadWithPSKAndContext(t *testing.T) {
	ctx := context.Background()

	Convey("Given a client configured with a user-defined psk uploader", t, func() {
		psk := []byte("test psk")
		cryptoUploaderMock := &mock.S3CryptoUploaderMock{
			UploadWithPSKFunc: func(ctx context.Context, in *s3manager.UploadInput, psk []byte) (*s3manager.UploadOutput, error) {
				return &s3manager.UploadOutput{}, nil
			},
		}
		cli := dps3.InstantiateClient(nil, nil, nil, cryptoUploaderMock, testBucket, ExpectedRegion, nil)

		Convey("Calling UploadWithPSKAndContext with a valid s3 key and context results in crypto UploadWithPSK being called as expected", func() {
			_, err := cli.UploadWithPSKAndContext(ctx, &s3manager.UploadInput{Key: &testS3Key}, psk)
			So(err, ShouldBeNil)
			So(len(cryptoUploaderMock.UploadWithPSKCalls()), ShouldEqual, 1)
			So(*cryptoUploaderMock.UploadWithPSKCalls()[0].In.Bucket, ShouldEqual, testBucket)
			So(cryptoUploaderMock.UploadWithPSKCalls()[0].Ctx, ShouldResemble, ctx)
		})

		Convey("Calling UploadWithPSKAndContext with nil context returns the expected error", func() {
			_, err := cli.UploadWithPSKAndContext(nil, &s3manager.UploadInput{Key: &testS3Key}, psk)
			So(err, ShouldResemble, dps3.NewError(
				errors.New("nil context provided to UploadWithPSKAndContext"),
				log.Data{
					"bucket_name": testBucket,
					"s3_key":      testS3Key,
				},
			))
			So(len(cryptoUploaderMock.UploadWithPSKCalls()), ShouldEqual, 0)
		})

		Convey("Calling UploadWithPSKAndContext with nil psk returns the expected error", func() {
			_, err := cli.UploadWithPSKAndContext(ctx, &s3manager.UploadInput{Key: &testS3Key}, nil)
			So(err, ShouldResemble, dps3.NewError(
				errors.New("nil or empty psk provided to UploadWithPSKAndContext"),
				log.Data{
					"bucket_name": testBucket,
					"s3_key":      testS3Key,
				},
			))
			So(len(cryptoUploaderMock.UploadWithPSKCalls()), ShouldEqual, 0)
		})

		Convey("Calling UploadWithPSKAndContext with nil input returns the expected error", func() {
			_, err := cli.UploadWithPSKAndContext(ctx, nil, psk)
			So(err, ShouldResemble, dps3.NewError(
				fmt.Errorf("validation error for UploadWithPSKAndContext: %w",
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
		errCryptoUploader := errors.New("failed to upload file")
		psk := []byte("test psk")
		cryptoUploaderMock := &mock.S3CryptoUploaderMock{
			UploadWithPSKFunc: func(ctx context.Context, in *s3manager.UploadInput, psk []byte) (*s3manager.UploadOutput, error) {
				return nil, errCryptoUploader
			},
		}
		cli := dps3.InstantiateClient(nil, nil, nil, cryptoUploaderMock, testBucket, ExpectedRegion, nil)

		Convey("Calling UploadWithPSKAndContext with a valid s3 key and context results in the expected error being returned", func() {
			_, err := cli.UploadWithPSKAndContext(ctx, &s3manager.UploadInput{Key: &testS3Key}, psk)
			So(err, ShouldResemble, dps3.NewError(
				fmt.Errorf("failed to upload with psk: %w", errCryptoUploader),
				log.Data{
					"bucket_name": testBucket,
					"s3_key":      testS3Key,
				},
			))
			So(len(cryptoUploaderMock.UploadWithPSKCalls()), ShouldEqual, 1)
			So(*cryptoUploaderMock.UploadWithPSKCalls()[0].In.Bucket, ShouldEqual, testBucket)
			So(cryptoUploaderMock.UploadWithPSKCalls()[0].Ctx, ShouldResemble, ctx)
		})
	})
}

func TestValidateInput(t *testing.T) {
	cli := dps3.InstantiateClient(nil, nil, nil, nil, testBucket, "", nil)

	Convey("validating an input with only an s3 key is successful with the expected logData", t, func() {
		logData, err := cli.ValidateUploadInput(&s3manager.UploadInput{
			Key: &testS3Key,
		})

		So(err, ShouldBeNil)
		So(logData, ShouldResemble, log.Data{
			"bucket_name": testBucket,
			"s3_key":      testS3Key,
		})
	})

	Convey("validating an input with an s3 key and the expected bucket is successful with the expected logData", t, func() {
		logData, err := cli.ValidateUploadInput(&s3manager.UploadInput{
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
		logData, err := cli.ValidateUploadInput(&s3manager.UploadInput{})

		So(err, ShouldResemble, errors.New("nil or empty s3 key provided in input"))
		So(logData, ShouldResemble, log.Data{
			"bucket_name": testBucket,
		})
	})

	Convey("validating an input without s3 key fails with the expected error and logData", t, func() {
		emptyKey := ""
		logData, err := cli.ValidateUploadInput(&s3manager.UploadInput{
			Key: &emptyKey,
		})
		So(err, ShouldResemble, errors.New("nil or empty s3 key provided in input"))
		So(logData, ShouldResemble, log.Data{
			"bucket_name": testBucket,
		})
	})

	Convey("validating an input with an s3 key and an unexpected bucket fails with the expected error and logData", t, func() {
		otherBucket := "otherBucket"
		logData, err := cli.ValidateUploadInput(&s3manager.UploadInput{
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
