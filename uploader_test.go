package s3client_test

import (
	"errors"
	"testing"

	s3client "github.com/ONSdigital/dp-s3/v2"
	"github.com/ONSdigital/dp-s3/v2/mock"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	. "github.com/smartystreets/goconvey/convey"
)

func uploadOK(in1 *s3manager.UploadInput, in2 ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	return &s3manager.UploadOutput{}, nil
}

func uploadWithPskOk(in1 *s3manager.UploadInput, in2 []byte) (*s3manager.UploadOutput, error) {
	return &s3manager.UploadOutput{}, nil
}

func TestUpload(t *testing.T) {

	Convey("Given an Uploader configured without user-defined psk", t, func() {
		sdkMock := &mock.S3SDKClientMock{}
		s3Cli := s3client.InstantiateClient(sdkMock, nil, ExistingBucket, ExpectedRegion, nil)
		sdkUploaderMock := &mock.S3SDKUploaderMock{
			UploadFunc: uploadOK,
		}
		uploader := s3client.InstantiateUploader(s3Cli, sdkUploaderMock, nil)

		Convey("Upload with no bucket in parameter uploads the file to the configured S3 bucket using AWS SDK", func() {
			_, err := uploader.Upload(&s3manager.UploadInput{})
			So(err, ShouldBeNil)
			So(len(sdkUploaderMock.UploadCalls()), ShouldEqual, 1)
			So(*sdkUploaderMock.UploadCalls()[0].In1.Bucket, ShouldEqual, ExistingBucket)
		})

		Convey("Upload with expected Bucket in parameter uploads the file to the configured S3 bucket using AWS SDK", func() {
			validBucket := ExistingBucket
			_, err := uploader.Upload(&s3manager.UploadInput{
				Bucket: &validBucket,
			})
			So(err, ShouldBeNil)
			So(len(sdkUploaderMock.UploadCalls()), ShouldEqual, 1)
			So(*sdkUploaderMock.UploadCalls()[0].In1.Bucket, ShouldEqual, ExistingBucket)
		})

		Convey("Tying to upload a file to the wrong S3 bucket results in ErrUnexpectedBucket error", func() {
			wrongBucket := "someBucket"
			_, err := uploader.Upload(&s3manager.UploadInput{
				Bucket: &wrongBucket,
			})
			So(err, ShouldResemble, s3client.NewError(
				errors.New("unexpected bucket name provided in upload input"),
				log.Data{
					"client_bucket_name": "csv-exported",
					"input_bucket_name":  "someBucket",
				},
			))
			So(len(sdkUploaderMock.UploadCalls()), ShouldEqual, 0)
		})
	})
}

func TestUploadWithPSK(t *testing.T) {
	s3Key := "my/s3/key"

	Convey("Given an Uploader configured with user-defined psk", t, func() {
		psk := []byte("test psk")
		sdkMock := &mock.S3SDKClientMock{}
		s3Cli := s3client.InstantiateClient(sdkMock, nil, ExistingBucket, ExpectedRegion, nil)
		cryptoUploaderMock := &mock.S3CryptoUploaderMock{
			UploadWithPSKFunc: uploadWithPskOk,
		}
		uploader := s3client.InstantiateUploader(s3Cli, nil, cryptoUploaderMock)

		Convey("UploadWithPSK with no bucket in parameter uploads the file to the configured S3 bucket using Crypto Uploader", func() {
			_, err := uploader.UploadWithPSK(&s3manager.UploadInput{Key: &s3Key}, psk)
			So(err, ShouldBeNil)
			So(len(cryptoUploaderMock.UploadWithPSKCalls()), ShouldEqual, 1)
			So(*cryptoUploaderMock.UploadWithPSKCalls()[0].In1.Bucket, ShouldEqual, ExistingBucket)
		})

		Convey("Upload with expected Bucket in parameter uploads the file to the configured S3 bucket using Crypto Uploader", func() {
			validBucket := ExistingBucket
			_, err := uploader.UploadWithPSK(&s3manager.UploadInput{
				Key:    &s3Key,
				Bucket: &validBucket,
			}, psk)
			So(err, ShouldBeNil)
			So(len(cryptoUploaderMock.UploadWithPSKCalls()), ShouldEqual, 1)
			So(*cryptoUploaderMock.UploadWithPSKCalls()[0].In1.Bucket, ShouldEqual, ExistingBucket)
		})

		Convey("Tying to upload a file to the wrong S3 bucket results in ErrUnexpectedBucket error", func() {
			wrongBucket := "someBucket"
			_, err := uploader.UploadWithPSK(&s3manager.UploadInput{
				Key:    &s3Key,
				Bucket: &wrongBucket,
			}, psk)
			So(err, ShouldResemble, s3client.NewError(
				errors.New("unexpected bucket name provided in upload input"),
				log.Data{
					"client_bucket_name": "csv-exported",
					"input_bucket_name":  "someBucket",
				},
			))
			So(len(cryptoUploaderMock.UploadWithPSKCalls()), ShouldEqual, 0)
		})
	})
}
