//+build integration

package s3_test

import (
	"context"
	"testing"

	dps3 "github.com/ONSdigital/dp-s3/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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
		s, err := session.NewSession(&aws.Config{
			Endpoint:         aws.String(localstackHost),
			Region:           aws.String("eu-west-1"),
			S3ForcePathStyle: aws.Bool(true),
			Credentials:      credentials.NewStaticCredentials("test", "test", ""),
		})
		So(err, ShouldBeNil)

		dpClient := dps3.NewClientWithSession(bucket, s)
		awsClient := s3.New(s)

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
				_, err = awsClient.HeadObject(&awss3.HeadObjectInput{
					Bucket: aws.String(bucket),
					Key:    aws.String(file),
				})

				So(err, ShouldBeNil)
			})

			Convey("And the file content in S3 should match given payload", func() {
				buf := aws.WriteAtBuffer{}
				dl := s3manager.NewDownloaderWithClient(awsClient)
				_, err = dl.Download(&buf, &awss3.GetObjectInput{
					Bucket: aws.String(bucket),
					Key:    aws.String(file),
				})

				So(err, ShouldBeNil)
				So(string(buf.Bytes()), ShouldEqual, payload)
			})
		})
	})
}
