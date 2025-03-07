//go:build integration
// +build integration

package s3_test

// TODO: move back to top of file

import (
	"context"
	"testing"

	dps3 "github.com/ONSdigital/dp-s3/v3"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	localstackHost = "http://localhost:4566"
	payload        = "TESTING"
	bucket         = "testing"
	file           = "test.csv"
	fileType       = "text/csv"
)

func TestMultipartUploadIntegrationTest(t *testing.T) {
	Convey("Multipart Upload", t, func() {
		cfg, err := config.LoadDefaultConfig(context.Background(),
			config.WithRegion("eu-west-1"),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
		)
		So(err, ShouldBeNil)

		dpClient := dps3.NewClientWithConfig(bucket, cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(localstackHost)
			o.UsePathStyle = true
		})

		awsClient := s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(localstackHost)
			o.UsePathStyle = true
		})

		Convey("When Uploading all parts of a file", func() {
			r, err := dpClient.UploadPart(context.Background(), &dps3.UploadPartRequest{
				UploadKey:   file,
				Type:        fileType,
				ChunkNumber: 1,
				TotalChunks: 1,
				FileName:    file,
			}, []byte(payload))

			Convey("Then response should contain etag & confirm all parts uploaded", func() {
				So(err, ShouldBeNil)
				So(r.AllPartsUploaded, ShouldBeTrue)
				So(r.Etag, ShouldEqual, `"907953dcbd01ad68db1f19be286936f4"`) // ETag should always be quoted https://datatracker.ietf.org/doc/html/rfc2616#section-14.19
			})

			Convey("And should have created the file in S3", func() {
				_, err = awsClient.HeadObject(context.Background(), &s3.HeadObjectInput{
					Bucket: aws.String(bucket),
					Key:    aws.String(file),
				})

				So(err, ShouldBeNil)
			})

			Convey("And the file content in S3 should match given payload", func() {
				buf := manager.NewWriteAtBuffer([]byte{})
				dl := manager.NewDownloader(awsClient)
				_, err = dl.Download(context.Background(), buf, &s3.GetObjectInput{
					Bucket: aws.String(bucket),
					Key:    aws.String(file),
				})

				So(err, ShouldBeNil)
				So(string(buf.Bytes()), ShouldEqual, payload)
			})
		})

		Convey("When uploading parts under 5mb", func() {
			dpClient.UploadPart(context.Background(), &dps3.UploadPartRequest{
				UploadKey:   file,
				Type:        fileType,
				ChunkNumber: 1,
				TotalChunks: 2,
				FileName:    file,
			}, []byte(payload))

			_, err := dpClient.UploadPart(context.Background(), &dps3.UploadPartRequest{
				UploadKey:   file,
				Type:        fileType,
				ChunkNumber: 2,
				TotalChunks: 2,
				FileName:    file,
			}, []byte(payload))

			Convey("Then it should return chunk too small error", func() {
				So(err.Error(), ShouldContainSubstring, "EntityTooSmall")
			})
		})
	})
}
