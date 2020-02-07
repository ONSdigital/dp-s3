package s3client_test

import (
	"context"
	"errors"
	"testing"

	s3client "github.com/ONSdigital/dp-s3"
	"github.com/ONSdigital/dp-s3/mock"
	"github.com/aws/aws-sdk-go/service/s3"
	. "github.com/smartystreets/goconvey/convey"
)

// Return a ListMultipartUploadsOutput with a single upload for the provided bucket, key and id
func createMultipartUploads(bucket, key, id *string) *s3.ListMultipartUploadsOutput {
	uploads := make([]*s3.MultipartUpload, 0, 1)
	uploads = append(uploads, &s3.MultipartUpload{
		Key:      key,
		UploadId: id,
	})
	return &s3.ListMultipartUploadsOutput{
		Bucket:  bucket,
		Uploads: uploads,
	}
}

// Returns a ListPartsOutput with a single part, corresponding to the provided partNumber
func createListPartsOutput(partNumber *int64) *s3.ListPartsOutput {
	parts := make([]*s3.Part, 0, 1)
	parts = append(parts, &s3.Part{
		PartNumber: partNumber,
	})
	return &s3.ListPartsOutput{
		Parts: parts,
	}
}

func TestCheckUpload(t *testing.T) {

	Convey("Given an S3 client", t, func() {

		bucket := ExistingBucket

		Convey("An error listing multipart uploads results in CheckUplaod failing with said error", func() {

			// Create S3Client with SDK Mock which fails to ListMultipartUploads
			listMultipartUploadsErr := errors.New("ListMultipartUploads failed")
			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(in1 *s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error) {
					return nil, listMultipartUploadsErr
				},
			}

			// Instantiate and call CheckUpload
			s3Cli := s3client.Instantiate(sdkMock, nil, bucket, ExpectedRegion)
			ok, err := s3Cli.CheckUploaded(context.Background(), &s3client.UploadRequest{
				UploadKey:   "12345",
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 1,
				FileName:    "helloworld",
			})

			// Validate
			So(ok, ShouldBeFalse)
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In1.Bucket, ShouldResemble, ExistingBucket)
			So(err, ShouldResemble, listMultipartUploadsErr)
		})

		Convey("If the upload key cannot be found in the list of multipart uploads, then the CheckUpload will fail with a ErrNotUploaded error", func() {

			// Create S3Client with SDK Mock with empty list of Multipart uploads
			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(in1 *s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error) {
					return &s3.ListMultipartUploadsOutput{}, nil
				},
			}

			// Instantiate and call CheckUpload
			s3Cli := s3client.Instantiate(sdkMock, nil, bucket, ExpectedRegion)
			ok, err := s3Cli.CheckUploaded(context.Background(), &s3client.UploadRequest{
				UploadKey:   "12345",
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 1,
				FileName:    "helloworld",
			})

			// Validate
			So(ok, ShouldBeFalse)
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In1.Bucket, ShouldResemble, bucket)
			So(err, ShouldResemble, &s3client.ErrNotUploaded{UploadKey: "12345"})
		})

		Convey("An error listing parts for a particular multipart upload results in ErrListParts error", func() {

			// Create S3Client with SDK Mock which fails to ListParts for a valid multipart upload
			skdListPartsErr := errors.New("ListMultipartUploads failed")

			expectedKey := "12345"
			expectedUploadID := "myID"

			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(in1 *s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error) {
					return createMultipartUploads(in1.Bucket, &expectedKey, &expectedUploadID), nil
				},
				ListPartsFunc: func(in1 *s3.ListPartsInput) (*s3.ListPartsOutput, error) {
					return nil, skdListPartsErr
				},
			}

			// Instantiate and call CheckUpload
			s3Cli := s3client.Instantiate(sdkMock, nil, bucket, ExpectedRegion)
			ok, err := s3Cli.CheckUploaded(context.Background(), &s3client.UploadRequest{
				UploadKey:   expectedKey,
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 1,
				FileName:    "helloworld",
			})

			// Validate
			So(ok, ShouldBeFalse)
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In1.Bucket, ShouldResemble, bucket)
			So(len(sdkMock.ListPartsCalls()), ShouldEqual, 1)
			So(err, ShouldResemble, &s3client.ErrListParts{Msg: skdListPartsErr.Error()})
		})

		Convey("Found part in incomplete upload", func() {

			expectedKey := "12345"
			expectedUploadID := "myID"
			expectedPart := int64(1)

			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(in1 *s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error) {
					return createMultipartUploads(in1.Bucket, &expectedKey, &expectedUploadID), nil
				},
				ListPartsFunc: func(in1 *s3.ListPartsInput) (*s3.ListPartsOutput, error) {
					return createListPartsOutput(&expectedPart), nil
				},
			}

			// Instantiate and call CheckUpload
			s3Cli := s3client.Instantiate(sdkMock, nil, bucket, ExpectedRegion)
			ok, err := s3Cli.CheckUploaded(context.Background(), &s3client.UploadRequest{
				UploadKey:   expectedKey,
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 10,
				FileName:    "helloworld",
			})

			// Validate
			So(ok, ShouldBeTrue)
			So(err, ShouldBeNil)
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In1.Bucket, ShouldResemble, bucket)
			So(len(sdkMock.ListPartsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListPartsCalls()[0].In1.Key, ShouldEqual, expectedKey)
			So(*sdkMock.ListPartsCalls()[0].In1.Bucket, ShouldEqual, bucket)
			So(*sdkMock.ListPartsCalls()[0].In1.UploadId, ShouldEqual, expectedUploadID)
		})

		Convey("Part not found", func() {

			expectedKey := "12345"
			expectedUploadID := "myID"
			unexpectedPart := int64(3)

			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(in1 *s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error) {
					return createMultipartUploads(in1.Bucket, &expectedKey, &expectedUploadID), nil
				},
				ListPartsFunc: func(in1 *s3.ListPartsInput) (*s3.ListPartsOutput, error) {
					return createListPartsOutput(&unexpectedPart), nil
				},
			}

			// Instantiate and call CheckUpload
			s3Cli := s3client.Instantiate(sdkMock, nil, bucket, ExpectedRegion)
			ok, err := s3Cli.CheckUploaded(context.Background(), &s3client.UploadRequest{
				UploadKey:   expectedKey,
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 10,
				FileName:    "helloworld",
			})

			// Validate
			So(ok, ShouldBeFalse)
			So(err, ShouldResemble, &s3client.ErrChunkNumberNotFound{1})
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In1.Bucket, ShouldResemble, bucket)
			So(len(sdkMock.ListPartsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListPartsCalls()[0].In1.Key, ShouldEqual, expectedKey)
			So(*sdkMock.ListPartsCalls()[0].In1.Bucket, ShouldEqual, bucket)
			So(*sdkMock.ListPartsCalls()[0].In1.UploadId, ShouldEqual, expectedUploadID)
		})

		Convey("Found part in comple upload results in CompleteMultipartUpload", func() {

			expectedKey := "12345"
			expectedUploadID := "myID"
			expectedPart := int64(1)

			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(in1 *s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error) {
					return createMultipartUploads(in1.Bucket, &expectedKey, &expectedUploadID), nil
				},
				ListPartsFunc: func(in1 *s3.ListPartsInput) (*s3.ListPartsOutput, error) {
					return createListPartsOutput(&expectedPart), nil
				},
				CompleteMultipartUploadFunc: func(input *s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error) {
					return &s3.CompleteMultipartUploadOutput{}, nil
				},
			}

			// Instantiate and call CheckUpload
			s3Cli := s3client.Instantiate(sdkMock, nil, bucket, ExpectedRegion)
			ok, err := s3Cli.CheckUploaded(context.Background(), &s3client.UploadRequest{
				UploadKey:   expectedKey,
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 1,
				FileName:    "helloworld",
			})

			// Validate
			So(ok, ShouldBeTrue)
			So(err, ShouldBeNil)
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In1.Bucket, ShouldResemble, bucket)
			So(len(sdkMock.ListPartsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListPartsCalls()[0].In1.Key, ShouldEqual, expectedKey)
			So(*sdkMock.ListPartsCalls()[0].In1.Bucket, ShouldEqual, bucket)
			So(*sdkMock.ListPartsCalls()[0].In1.UploadId, ShouldEqual, expectedUploadID)
			So(len(sdkMock.CompleteMultipartUploadCalls()), ShouldEqual, 1)
		})

	})
}

func TestUpload(t *testing.T) {
	// TODO implement
}

func TestGet(t *testing.T) {

	Convey("Given an S3 client configured with a bucket and region", t, func() {

		sdkMock := &mock.S3SDKClientMock{
			GetObjectFunc: func(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{}, nil
			},
		}

		s3Cli := s3client.Instantiate(sdkMock, nil, "bucket", "eu-north-1")

		Convey("getURL returns the correct fully qualified URL with the bucket, region and requested path", func() {
			_, err := s3Cli.Get("objectKey")
			So(err, ShouldBeNil)
		})

	})
}

func TestGetUrl(t *testing.T) {

	Convey("Given an S3 client configured with a bucket and region", t, func() {
		s3Cli := s3client.Instantiate(nil, nil, "bucket", "eu-north-1")

		Convey("getURL returns the correct fully qualified URL with the bucket, region and requested path", func() {
			url := s3Cli.GetURL("path")
			So(url, ShouldEqual, "https://s3-eu-north-1.amazonaws.com/bucket/path")
		})

	})
}
