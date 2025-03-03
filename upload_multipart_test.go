package s3_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	dps3 "github.com/ONSdigital/dp-s3/v2"
	"github.com/ONSdigital/dp-s3/v2/mock"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUploadPart(t *testing.T) {

	Convey("Given an S3 client with the intention of performing a multi-part upload", t, func() {

		bucket := ExistingBucket
		payload := []byte("test data")
		testUploadId := "testUploadId"
		expectedPart := int32(1)
		testKey := "testKey"

		Convey("An error listing multipart uploads results in Upload failing with said error", func() {

			// Create S3 client with SDK Mock which fails to ListMultipartUploads
			listMultipartUploadsErr := errors.New("ListMultipartUploads failed")
			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(ctx context.Context, in1 *s3.ListMultipartUploadsInput, optFns ...func(*s3.Options)) (*s3.ListMultipartUploadsOutput, error) {
					return nil, listMultipartUploadsErr
				},
			}

			// Instantiate and call Upload
			cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, ExpectedRegion, aws.Config{})
			_, err := cli.UploadPart(context.Background(), &dps3.UploadPartRequest{
				UploadKey:   testKey,
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 1,
				FileName:    "helloworld",
			}, payload)

			// Validate
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldResemble, fmt.Errorf("error fetching multipart list: %w", listMultipartUploadsErr).Error())
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In.Bucket, ShouldResemble, ExistingBucket)
		})

		Convey("If the upload S3 object key can be found in the list of multipart upload, Upload will use it", func() {

			// Create S3 client with SDK Mock with empty list of Multipart uploads
			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(ctx context.Context, in1 *s3.ListMultipartUploadsInput, optFns ...func(*s3.Options)) (*s3.ListMultipartUploadsOutput, error) {
					return createUploads(testUploadId, testKey), nil
				},
				UploadPartFunc: func(ctx context.Context, in1 *s3.UploadPartInput, optFns ...func(*s3.Options)) (*s3.UploadPartOutput, error) {
					return &s3.UploadPartOutput{ETag: aws.String(`"1234567890"`)}, nil
				},
				ListPartsFunc: func(ctx context.Context, in1 *s3.ListPartsInput, optFns ...func(*s3.Options)) (*s3.ListPartsOutput, error) {
					return createListPartsOutput(&expectedPart), nil
				},
			}

			// Instantiate and call Upload
			cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, ExpectedRegion, aws.Config{})
			response, err := cli.UploadPart(context.Background(), &dps3.UploadPartRequest{
				UploadKey:   testKey,
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 2,
				FileName:    "helloworld",
			}, payload)

			// Validate
			So(err, ShouldBeNil)
			So(response.Etag, ShouldEqual, `"1234567890"`)
			So(response.AllPartsUploaded, ShouldBeFalse)
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In.Bucket, ShouldResemble, bucket)
			So(len(sdkMock.ListPartsCalls()), ShouldEqual, 1)
			So(len(sdkMock.UploadPartCalls()), ShouldEqual, 1)
			So(*sdkMock.UploadPartCalls()[0].In.UploadId, ShouldEqual, testUploadId)
			So(*sdkMock.UploadPartCalls()[0].In.Bucket, ShouldEqual, bucket)
			So(*sdkMock.UploadPartCalls()[0].In.Key, ShouldEqual, testKey)
		})

		Convey("If the upload S3 object key cannot be found in the list of multipart uploads, Upload will create a new one, "+
			"and don't complete it if some chunks have not been uploaded yet", func() {

			// Create S3 client with SDK Mock with empty list of Multipart uploads
			testUploadId := "testUploadId"
			testKey := "testKey"
			expectedPart := int32(1)
			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(ctx context.Context, in1 *s3.ListMultipartUploadsInput, optFns ...func(*s3.Options)) (*s3.ListMultipartUploadsOutput, error) {
					return &s3.ListMultipartUploadsOutput{}, nil
				},
				CreateMultipartUploadFunc: func(ctx context.Context, in1 *s3.CreateMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CreateMultipartUploadOutput, error) {
					return &s3.CreateMultipartUploadOutput{UploadId: &testUploadId}, nil
				},
				UploadPartFunc: func(ctx context.Context, in1 *s3.UploadPartInput, optFns ...func(*s3.Options)) (*s3.UploadPartOutput, error) {
					return &s3.UploadPartOutput{ETag: aws.String(`"1234567890"`)}, nil
				},
				ListPartsFunc: func(ctx context.Context, in1 *s3.ListPartsInput, optFns ...func(*s3.Options)) (*s3.ListPartsOutput, error) {
					return createListPartsOutput(&expectedPart), nil
				},
			}

			// Instantiate and call Upload
			s3Cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, ExpectedRegion, aws.Config{})
			response, err := s3Cli.UploadPart(context.Background(), &dps3.UploadPartRequest{
				UploadKey:   testKey,
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 2,
				FileName:    "helloworld",
			}, payload)

			// Validate
			So(err, ShouldBeNil)
			So(response.Etag, ShouldEqual, `"1234567890"`)
			So(response.AllPartsUploaded, ShouldBeFalse)
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In.Bucket, ShouldResemble, bucket)
			So(len(sdkMock.CreateMultipartUploadCalls()), ShouldEqual, 1)
			So(len(sdkMock.UploadPartCalls()), ShouldEqual, 1)
			So(*sdkMock.UploadPartCalls()[0].In.UploadId, ShouldEqual, testUploadId)
			So(*sdkMock.UploadPartCalls()[0].In.Bucket, ShouldEqual, bucket)
			So(*sdkMock.UploadPartCalls()[0].In.Key, ShouldEqual, testKey)
			So(len(sdkMock.ListPartsCalls()), ShouldEqual, 1)
		})

		Convey("If the upload S3 object key cannot be found in the list of multipart uploads, Upload will create a new one, "+
			"and complete it if all chunks have been uploaded", func() {

			// Create S3 client with SDK Mock with empty list of Multipart uploads
			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(ctx context.Context, in1 *s3.ListMultipartUploadsInput, optFns ...func(*s3.Options)) (*s3.ListMultipartUploadsOutput, error) {
					return &s3.ListMultipartUploadsOutput{}, nil
				},
				CreateMultipartUploadFunc: func(ctx context.Context, in1 *s3.CreateMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CreateMultipartUploadOutput, error) {
					return &s3.CreateMultipartUploadOutput{UploadId: &testUploadId}, nil
				},
				UploadPartFunc: func(ctx context.Context, in1 *s3.UploadPartInput, optFns ...func(*s3.Options)) (*s3.UploadPartOutput, error) {
					return &s3.UploadPartOutput{ETag: aws.String(`"1234567890"`)}, nil
				},
				ListPartsFunc: func(ctx context.Context, in1 *s3.ListPartsInput, optFns ...func(*s3.Options)) (*s3.ListPartsOutput, error) {
					return createListPartsOutput(&expectedPart), nil
				},
				CompleteMultipartUploadFunc: func(ctx context.Context, input *s3.CompleteMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CompleteMultipartUploadOutput, error) {
					return &s3.CompleteMultipartUploadOutput{}, nil
				},
			}

			// Instantiate and call Upload
			s3Cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, ExpectedRegion, aws.Config{})
			response, err := s3Cli.UploadPart(context.Background(), &dps3.UploadPartRequest{
				UploadKey:   testKey,
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 1,
				FileName:    "helloworld",
			}, payload)

			// Validate
			So(err, ShouldBeNil)
			So(response.Etag, ShouldEqual, `"1234567890"`)
			So(response.AllPartsUploaded, ShouldBeTrue)
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In.Bucket, ShouldResemble, bucket)
			So(len(sdkMock.CreateMultipartUploadCalls()), ShouldEqual, 1)
			So(len(sdkMock.UploadPartCalls()), ShouldEqual, 1)
			So(*sdkMock.UploadPartCalls()[0].In.UploadId, ShouldEqual, testUploadId)
			So(*sdkMock.UploadPartCalls()[0].In.Bucket, ShouldEqual, bucket)
			So(*sdkMock.UploadPartCalls()[0].In.Key, ShouldEqual, testKey)
			So(len(sdkMock.ListPartsCalls()), ShouldEqual, 1)
			So(len(sdkMock.CompleteMultipartUploadCalls()), ShouldEqual, 1)
		})

		Convey("UploadWithPsk performs an upload with the provided PSK", func() {

			psk := []byte("test psk")

			// Create S3 client with SDK Mock with empty list of Multipart uploads
			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(ctx context.Context, in1 *s3.ListMultipartUploadsInput, optFns ...func(*s3.Options)) (*s3.ListMultipartUploadsOutput, error) {
					return &s3.ListMultipartUploadsOutput{}, nil
				},
				CreateMultipartUploadFunc: func(ctx context.Context, in1 *s3.CreateMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CreateMultipartUploadOutput, error) {
					return &s3.CreateMultipartUploadOutput{UploadId: &testUploadId}, nil
				},
				ListPartsFunc: func(ctx context.Context, in1 *s3.ListPartsInput, optFns ...func(*s3.Options)) (*s3.ListPartsOutput, error) {
					return createListPartsOutput(&expectedPart), nil
				},
			}

			cryptoMock := &mock.S3CryptoClientMock{
				UploadPartWithPSKFunc: func(ctx context.Context, in1 *s3.UploadPartInput, in2 []byte) (*s3.UploadPartOutput, error) {
					return &s3.UploadPartOutput{ETag: aws.String(`"1234567890"`)}, nil
				},
			}

			// Instantiate and call UploadWithPsk
			s3Cli := dps3.InstantiateClient(sdkMock, cryptoMock, nil, nil, bucket, ExpectedRegion, aws.Config{})
			response, err := s3Cli.UploadPartWithPsk(context.Background(), &dps3.UploadPartRequest{
				UploadKey:   testKey,
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 2,
				FileName:    "helloworld",
			}, payload, psk)

			// Validate
			So(err, ShouldBeNil)
			So(response.Etag, ShouldEqual, `"1234567890"`)
			So(response.AllPartsUploaded, ShouldBeFalse)
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In.Bucket, ShouldResemble, bucket)
			So(len(sdkMock.CreateMultipartUploadCalls()), ShouldEqual, 1)
			So(len(cryptoMock.UploadPartWithPSKCalls()), ShouldEqual, 1)
			So(len(sdkMock.ListPartsCalls()), ShouldEqual, 1)
		})

		Convey("UploadWithPsk performs an upload with the provided PSK - all parts uploaded", func() {
			psk := []byte("test psk")

			// Create S3 client with SDK Mock with empty list of Multipart uploads
			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(ctx context.Context, in1 *s3.ListMultipartUploadsInput, optFns ...func(*s3.Options)) (*s3.ListMultipartUploadsOutput, error) {
					return &s3.ListMultipartUploadsOutput{}, nil
				},
				CreateMultipartUploadFunc: func(ctx context.Context, in1 *s3.CreateMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CreateMultipartUploadOutput, error) {
					return &s3.CreateMultipartUploadOutput{UploadId: &testUploadId}, nil
				},
				ListPartsFunc: func(ctx context.Context, in1 *s3.ListPartsInput, optFns ...func(*s3.Options)) (*s3.ListPartsOutput, error) {
					return createListPartsOutput(&expectedPart), nil
				},
				CompleteMultipartUploadFunc: func(ctx context.Context, input *s3.CompleteMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CompleteMultipartUploadOutput, error) {
					return &s3.CompleteMultipartUploadOutput{}, nil
				},
			}

			cryptoMock := &mock.S3CryptoClientMock{
				UploadPartWithPSKFunc: func(ctx context.Context, in1 *s3.UploadPartInput, in2 []byte) (*s3.UploadPartOutput, error) {
					return &s3.UploadPartOutput{ETag: aws.String(`"1234567890"`)}, nil
				},
			}

			// Instantiate and call Upload
			s3Cli := dps3.InstantiateClient(sdkMock, cryptoMock, nil, nil, bucket, ExpectedRegion, aws.Config{})
			response, err := s3Cli.UploadPartWithPsk(context.Background(), &dps3.UploadPartRequest{
				UploadKey:   testKey,
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 1,
				FileName:    "helloworld",
			}, payload, psk)

			// Validate
			So(err, ShouldBeNil)
			So(response.Etag, ShouldEqual, `"1234567890"`)
			So(response.AllPartsUploaded, ShouldBeTrue)
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In.Bucket, ShouldResemble, bucket)
			So(len(sdkMock.CreateMultipartUploadCalls()), ShouldEqual, 1)
			So(len(sdkMock.UploadPartCalls()), ShouldEqual, 0)
			So(len(sdkMock.ListPartsCalls()), ShouldEqual, 1)
			So(len(sdkMock.CompleteMultipartUploadCalls()), ShouldEqual, 1)
		})
	})
}

func TestCheckUpload(t *testing.T) {

	Convey("Given an S3 client with the intention of checking if a chunk has been uploaded in a multipart upload", t, func() {

		bucket := ExistingBucket

		filename := "helloworld"
		Convey("An error listing multipart uploads results in CheckUplaod failing with said error", func() {

			// Create S3 client with SDK Mock which fails to ListMultipartUploads
			listMultipartUploadsErr := errors.New("ListMultipartUploads failed")
			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(ctx context.Context, in1 *s3.ListMultipartUploadsInput, optFns ...func(*s3.Options)) (*s3.ListMultipartUploadsOutput, error) {
					return nil, listMultipartUploadsErr
				},
			}

			// Instantiate and call CheckUpload
			cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, ExpectedRegion, aws.Config{})
			ok, err := cli.CheckPartUploaded(context.Background(), &dps3.UploadPartRequest{
				UploadKey:   "12345",
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 1,
				FileName:    filename,
			})

			// Validate
			So(ok, ShouldBeFalse)
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In.Bucket, ShouldResemble, ExistingBucket)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldResemble, fmt.Errorf("error fetching multipart upload list: %w", listMultipartUploadsErr).Error())
		})

		Convey("If the upload S3 object key cannot be found in the list of multipart uploads, then the CheckUpload will fail with a ErrNotUploaded error", func() {

			// Create S3 client with SDK Mock with empty list of Multipart uploads
			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(ctx context.Context, in1 *s3.ListMultipartUploadsInput, optFns ...func(*s3.Options)) (*s3.ListMultipartUploadsOutput, error) {
					return &s3.ListMultipartUploadsOutput{}, nil
				},
			}

			// Instantiate and call CheckUpload
			cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, ExpectedRegion, aws.Config{})
			ok, err := cli.CheckPartUploaded(context.Background(), &dps3.UploadPartRequest{
				UploadKey:   "12345",
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 1,
				FileName:    filename,
			})

			// Validate
			So(ok, ShouldBeFalse)
			_, isNewNotUploaded := err.(*dps3.ErrNotUploaded)
			So(isNewNotUploaded, ShouldBeTrue)

			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In.Bucket, ShouldResemble, bucket)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldResemble, "s3 key not uploaded")
		})

		Convey("An error listing parts for a particular multipart upload results in ErrListParts error", func() {

			// Create S3 client with SDK Mock which fails to ListParts for a valid multipart upload
			skdListPartsErr := errors.New("ListMultipartUploads failed")

			expectedKey := "12345"
			expectedUploadID := "myID"

			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(ctx context.Context, in1 *s3.ListMultipartUploadsInput, optFns ...func(*s3.Options)) (*s3.ListMultipartUploadsOutput, error) {
					return createMultipartUploads(in1.Bucket, &expectedKey, &expectedUploadID), nil
				},
				ListPartsFunc: func(ctx context.Context, in1 *s3.ListPartsInput, optFns ...func(*s3.Options)) (*s3.ListPartsOutput, error) {
					return nil, skdListPartsErr
				},
			}

			// Instantiate and call CheckUpload
			cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, ExpectedRegion, aws.Config{})
			ok, err := cli.CheckPartUploaded(context.Background(), &dps3.UploadPartRequest{
				UploadKey:   expectedKey,
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 1,
				FileName:    filename,
			})

			// Validate
			So(ok, ShouldBeFalse)
			_, isNewErrListParts := err.(*dps3.ErrListParts)
			So(isNewErrListParts, ShouldBeTrue)

			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In.Bucket, ShouldResemble, bucket)
			So(len(sdkMock.ListPartsCalls()), ShouldEqual, 1)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldResemble, fmt.Errorf("list parts failed: %w", skdListPartsErr).Error())
		})

		Convey("If the chunk has been uploaded but the multipart upload is not completed yet, then the function should return true", func() {

			expectedKey := "12345"
			expectedUploadID := "myID"
			expectedPart := int32(1)

			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(ctx context.Context, in1 *s3.ListMultipartUploadsInput, optFns ...func(*s3.Options)) (*s3.ListMultipartUploadsOutput, error) {
					return createMultipartUploads(in1.Bucket, &expectedKey, &expectedUploadID), nil
				},
				ListPartsFunc: func(ctx context.Context, in1 *s3.ListPartsInput, optFns ...func(*s3.Options)) (*s3.ListPartsOutput, error) {
					return createListPartsOutput(&expectedPart), nil
				},
			}

			// Instantiate and call CheckUpload
			cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, ExpectedRegion, aws.Config{})
			ok, err := cli.CheckPartUploaded(context.Background(), &dps3.UploadPartRequest{
				UploadKey:   expectedKey,
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 10,
				FileName:    filename,
			})

			// Validate
			So(ok, ShouldBeTrue)
			So(err, ShouldBeNil)
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In.Bucket, ShouldResemble, bucket)
			So(len(sdkMock.ListPartsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListPartsCalls()[0].In.Key, ShouldEqual, expectedKey)
			So(*sdkMock.ListPartsCalls()[0].In.Bucket, ShouldEqual, bucket)
			So(*sdkMock.ListPartsCalls()[0].In.UploadId, ShouldEqual, expectedUploadID)
		})

		Convey("Provided chunk not being found in the list of parts results in ErrChunkNumberNotFound being returned", func() {

			expectedKey := "12345"
			expectedUploadID := "myID"
			unexpectedPart := int32(3)

			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(ctx context.Context, in1 *s3.ListMultipartUploadsInput, optFns ...func(*s3.Options)) (*s3.ListMultipartUploadsOutput, error) {
					return createMultipartUploads(in1.Bucket, &expectedKey, &expectedUploadID), nil
				},
				ListPartsFunc: func(ctx context.Context, in1 *s3.ListPartsInput, optFns ...func(*s3.Options)) (*s3.ListPartsOutput, error) {
					return createListPartsOutput(&unexpectedPart), nil
				},
			}

			// Instantiate and call CheckUpload
			cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, ExpectedRegion, aws.Config{})
			ok, err := cli.CheckPartUploaded(context.Background(), &dps3.UploadPartRequest{
				UploadKey:   expectedKey,
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 10,
				FileName:    filename,
			})

			// Validate
			So(ok, ShouldBeFalse)

			So(err.Error(), ShouldEqual, "chunk number not found")
			_, isNewChunkNumberNotFound := err.(*dps3.ErrChunkNumberNotFound)
			So(isNewChunkNumberNotFound, ShouldBeTrue)
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In.Bucket, ShouldResemble, bucket)
			So(len(sdkMock.ListPartsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListPartsCalls()[0].In.Key, ShouldEqual, expectedKey)
			So(*sdkMock.ListPartsCalls()[0].In.Bucket, ShouldEqual, bucket)
			So(*sdkMock.ListPartsCalls()[0].In.UploadId, ShouldEqual, expectedUploadID)
		})

		Convey("Provided chunk being successfully uploaded as part of a completed multipart upload results in the function completing the upload and returning true", func() {

			expectedKey := "12345"
			expectedUploadID := "myID"
			expectedPart := int32(1)

			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(ctx context.Context, in1 *s3.ListMultipartUploadsInput, optFns ...func(*s3.Options)) (*s3.ListMultipartUploadsOutput, error) {
					return createMultipartUploads(in1.Bucket, &expectedKey, &expectedUploadID), nil
				},
				ListPartsFunc: func(ctx context.Context, in1 *s3.ListPartsInput, optFns ...func(*s3.Options)) (*s3.ListPartsOutput, error) {
					return createListPartsOutput(&expectedPart), nil
				},
				CompleteMultipartUploadFunc: func(ctx context.Context, input *s3.CompleteMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CompleteMultipartUploadOutput, error) {
					return &s3.CompleteMultipartUploadOutput{}, nil
				},
			}

			// Instantiate and call CheckUpload
			cli := dps3.InstantiateClient(sdkMock, nil, nil, nil, bucket, ExpectedRegion, aws.Config{})
			ok, err := cli.CheckPartUploaded(context.Background(), &dps3.UploadPartRequest{
				UploadKey:   expectedKey,
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 1,
				FileName:    filename,
			})

			// Validate
			So(ok, ShouldBeTrue)
			So(err, ShouldBeNil)
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In.Bucket, ShouldResemble, bucket)
			So(len(sdkMock.ListPartsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListPartsCalls()[0].In.Key, ShouldEqual, expectedKey)
			So(*sdkMock.ListPartsCalls()[0].In.Bucket, ShouldEqual, bucket)
			So(*sdkMock.ListPartsCalls()[0].In.UploadId, ShouldEqual, expectedUploadID)
			So(len(sdkMock.CompleteMultipartUploadCalls()), ShouldEqual, 1)
		})

	})
}

func createMultipartUploads(bucket, key, id *string) *s3.ListMultipartUploadsOutput {
	uploads := make([]types.MultipartUpload, 0, 1)
	uploads = append(uploads, types.MultipartUpload{
		Key:      key,
		UploadId: id,
	})
	return &s3.ListMultipartUploadsOutput{
		Bucket:  bucket,
		Uploads: uploads,
	}
}

// createListPartsOutput returns a ListPartsOutput with a single part, corresponding to the provided partNumber
func createListPartsOutput(partNumber *int32) *s3.ListPartsOutput {
	parts := make([]types.Part, 0, 1)
	parts = append(parts, types.Part{
		PartNumber: partNumber,
	})
	return &s3.ListPartsOutput{
		Parts: parts,
	}
}

// createUploads returns a ListMultipartUploadsOutput with a single upload, with the provided uploadID
func createUploads(uploadID, key string) *s3.ListMultipartUploadsOutput {
	uploads := make([]types.MultipartUpload, 0, 1)
	uploads = append(uploads, types.MultipartUpload{
		UploadId: &uploadID,
		Key:      &key,
	})
	return &s3.ListMultipartUploadsOutput{
		Uploads: uploads,
	}
}
