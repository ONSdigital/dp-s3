package s3client_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"testing"

	s3client "github.com/ONSdigital/dp-s3"
	"github.com/ONSdigital/dp-s3/mock"
	"github.com/aws/aws-sdk-go/service/s3"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGet(t *testing.T) {

	Convey("Given an S3 client configured with a bucket and region", t, func() {

		payload := []byte("test data")
		bucket := "myBucket"
		objKey := "my/object/key"
		region := "eu-north-1"

		sdkMock := &mock.S3SDKClientMock{
			GetObjectFunc: func(input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{
					Body: ioutil.NopCloser(bytes.NewReader(payload)),
				}, nil
			},
		}

		s3Cli := s3client.InstantiateClient(sdkMock, nil, bucket, region, nil)

		Convey("Get returns an io.Reader with the expected payload", func() {
			ret, err := s3Cli.Get(objKey)
			So(err, ShouldBeNil)
			buf := new(bytes.Buffer)
			buf.ReadFrom(ret)
			ret.Close()
			So(buf.Bytes(), ShouldResemble, payload)
			So(len(sdkMock.GetObjectCalls()), ShouldEqual, 1)
			So(sdkMock.GetObjectCalls()[0].In1, ShouldResemble, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &objKey,
			})
		})

		Convey("GetFromS3URL called with a valid global URL returns an io.Reader with the expected payload", func() {
			validGlobalURL := fmt.Sprintf("s3://%s/%s", bucket, objKey)
			ret, err := s3Cli.GetFromS3URL(validGlobalURL, s3client.StyleAliasVirtualHosted)
			So(err, ShouldBeNil)
			buf := new(bytes.Buffer)
			buf.ReadFrom(ret)
			ret.Close()
			So(buf.Bytes(), ShouldResemble, payload)
			So(len(sdkMock.GetObjectCalls()), ShouldEqual, 1)
			So(sdkMock.GetObjectCalls()[0].In1, ShouldResemble, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &objKey,
			})
		})

		Convey("GetFromS3URL called with a valid regional URL returns an io.Reader with the expected payload", func() {
			validRegionalURL := fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s", region, bucket, objKey)
			ret, err := s3Cli.GetFromS3URL(validRegionalURL, s3client.StylePath)
			So(err, ShouldBeNil)
			buf := new(bytes.Buffer)
			buf.ReadFrom(ret)
			ret.Close()
			So(buf.Bytes(), ShouldResemble, payload)
			So(len(sdkMock.GetObjectCalls()), ShouldEqual, 1)
			So(sdkMock.GetObjectCalls()[0].In1, ShouldResemble, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &objKey,
			})
		})

		Convey("GetFromS3URL called with a valid global URL with the wrong bucket returns ErrUnexpectedBucket", func() {
			wrongBucketGlobalURL := fmt.Sprintf("s3://%s/%s", "wrongBucket", objKey)
			_, err := s3Cli.GetFromS3URL(wrongBucketGlobalURL, s3client.StyleAliasVirtualHosted)
			So(err, ShouldResemble, &s3client.ErrUnexpectedBucket{ExpectedBucketName: bucket, BucketName: "wrongBucket"})
			So(len(sdkMock.GetObjectCalls()), ShouldEqual, 0)
		})

		Convey("GetFromS3URL called with a valid regional URL with the wrong region returns ErrUnexpectedBucket", func() {
			wrongRegionRegionalURL := fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s", "wrongRegion", bucket, objKey)
			_, err := s3Cli.GetFromS3URL(wrongRegionRegionalURL, s3client.StylePath)
			So(err, ShouldResemble, &s3client.ErrUnexpectedRegion{ExpectedRegion: region, Region: "wrongRegion"})
			So(len(sdkMock.GetObjectCalls()), ShouldEqual, 0)
		})

		Convey("GetFromS3URL called with a malformed URL returns error", func() {
			malformedURL := "This%Url%Is%Malformed"
			_, err := s3Cli.GetFromS3URL(malformedURL, s3client.StyleAliasVirtualHosted)
			So(err, ShouldResemble, &url.Error{Op: "parse", URL: malformedURL, Err: url.EscapeError("%Ur")})
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

		cryptoMock := &mock.S3CryptoClientMock{
			GetObjectWithPSKFunc: func(input *s3.GetObjectInput, inPsk []byte) (*s3.GetObjectOutput, error) {
				return &s3.GetObjectOutput{
					Body: ioutil.NopCloser(bytes.NewReader(payload)),
				}, nil
			},
		}

		s3Cli := s3client.InstantiateClient(nil, cryptoMock, bucket, region, nil)

		Convey("GetWithPSK returns an io.Reader with the expected payload", func() {
			ret, err := s3Cli.GetWithPSK(objKey, psk)
			So(err, ShouldBeNil)
			buf := new(bytes.Buffer)
			buf.ReadFrom(ret)
			ret.Close()
			So(buf.Bytes(), ShouldResemble, payload)
			So(len(cryptoMock.GetObjectWithPSKCalls()), ShouldEqual, 1)
			So(cryptoMock.GetObjectWithPSKCalls()[0].In2, ShouldResemble, psk)
			So(cryptoMock.GetObjectWithPSKCalls()[0].In1, ShouldResemble, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &objKey,
			})
		})

		Convey("GetFromS3URLWithPSK called with a valid global URL returns an io.Reader with the expected payload", func() {
			validGlobalURL := fmt.Sprintf("s3://%s/%s", bucket, objKey)
			ret, err := s3Cli.GetFromS3URLWithPSK(validGlobalURL, s3client.StyleAliasVirtualHosted, psk)
			So(err, ShouldBeNil)
			buf := new(bytes.Buffer)
			buf.ReadFrom(ret)
			ret.Close()
			So(buf.Bytes(), ShouldResemble, payload)
			So(len(cryptoMock.GetObjectWithPSKCalls()), ShouldEqual, 1)
			So(cryptoMock.GetObjectWithPSKCalls()[0].In2, ShouldResemble, psk)
			So(cryptoMock.GetObjectWithPSKCalls()[0].In1, ShouldResemble, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &objKey,
			})
		})
	})
}

func TestPutWithPSK(t *testing.T) {

	Convey("Given an S3 client configured with a bucket, region and psk", t, func() {

		psk := []byte("test psk")
		payload := []byte("test data")
		bucket := "myBucket"
		objKey := "my/object/key"
		region := "eu-north-1"

		payloadReader := bytes.NewReader(payload)

		cryptoMock := &mock.S3CryptoClientMock{
			PutObjectWithPSKFunc: func(in1 *s3.PutObjectInput, in2 []byte) (*s3.PutObjectOutput, error) {
				return &s3.PutObjectOutput{}, nil
			},
		}

		s3Cli := s3client.InstantiateClient(nil, cryptoMock, bucket, region, nil)

		Convey("PutWithPSK calls the expected cryptoClient with provided key, reader and client-configured bucket", func() {
			err := s3Cli.PutWithPSK(&objKey, payloadReader, psk)
			So(err, ShouldBeNil)
			So(len(cryptoMock.PutObjectWithPSKCalls()), ShouldEqual, 1)
			So(cryptoMock.PutObjectWithPSKCalls()[0].In1, ShouldResemble, &s3.PutObjectInput{
				Bucket: &bucket,
				Key:    &objKey,
				Body:   payloadReader,
			})
			So(cryptoMock.PutObjectWithPSKCalls()[0].In2, ShouldResemble, psk)
		})
	})
}

func TestUploadPart(t *testing.T) {

	Convey("Given an S3 client with the intention of performing a multi-part upload", t, func() {

		bucket := ExistingBucket
		payload := []byte("test data")
		testUploadId := "testUploadId"
		expectedPart := int64(1)
		testKey := "testKey"

		Convey("An error listing multipart uploads results in Upload failing with said error", func() {

			// Create S3Client with SDK Mock which fails to ListMultipartUploads
			listMultipartUploadsErr := errors.New("ListMultipartUploads failed")
			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(in1 *s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error) {
					return nil, listMultipartUploadsErr
				},
			}

			// Instantiate and call Upload
			s3Cli := s3client.InstantiateClient(sdkMock, nil, bucket, ExpectedRegion, nil)
			err := s3Cli.UploadPart(context.Background(), &s3client.UploadPartRequest{
				UploadKey:   testKey,
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 1,
				FileName:    "helloworld",
			}, payload)

			// Validate
			So(err, ShouldResemble, listMultipartUploadsErr)
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In1.Bucket, ShouldResemble, ExistingBucket)
		})

		Convey("If the upload S3 object key can be found in the list of multipart upload, Upload will use it", func() {

			// Create S3Client with SDK Mock with empty list of Multipart uploads
			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(in1 *s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error) {
					return createUploads(testUploadId, testKey), nil
				},
				UploadPartFunc: func(in1 *s3.UploadPartInput) (*s3.UploadPartOutput, error) {
					return &s3.UploadPartOutput{}, nil
				},
				ListPartsFunc: func(in1 *s3.ListPartsInput) (*s3.ListPartsOutput, error) {
					return createListPartsOutput(&expectedPart), nil
				},
			}

			// Instantiate and call Upload
			s3Cli := s3client.InstantiateClient(sdkMock, nil, bucket, ExpectedRegion, nil)
			err := s3Cli.UploadPart(context.Background(), &s3client.UploadPartRequest{
				UploadKey:   testKey,
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 2,
				FileName:    "helloworld",
			}, payload)

			// Validate
			So(err, ShouldBeNil)
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In1.Bucket, ShouldResemble, bucket)
			So(len(sdkMock.ListPartsCalls()), ShouldEqual, 1)
			So(len(sdkMock.UploadPartCalls()), ShouldEqual, 1)
			So(*sdkMock.UploadPartCalls()[0].In1.UploadId, ShouldEqual, testUploadId)
			So(*sdkMock.UploadPartCalls()[0].In1.Bucket, ShouldEqual, bucket)
			So(*sdkMock.UploadPartCalls()[0].In1.Key, ShouldEqual, testKey)
		})

		Convey("If the upload S3 object key cannot be found in the list of multipart uploads, Upload will create a new one, "+
			"and don't complete it if some chunks have not been uploaded yet", func() {

			// Create S3Client with SDK Mock with empty list of Multipart uploads
			testUploadId := "testUploadId"
			testKey := "testKey"
			expectedPart := int64(1)
			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(in1 *s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error) {
					return &s3.ListMultipartUploadsOutput{}, nil
				},
				CreateMultipartUploadFunc: func(in1 *s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error) {
					return &s3.CreateMultipartUploadOutput{UploadId: &testUploadId}, nil
				},
				UploadPartFunc: func(in1 *s3.UploadPartInput) (*s3.UploadPartOutput, error) {
					return &s3.UploadPartOutput{}, nil
				},
				ListPartsFunc: func(in1 *s3.ListPartsInput) (*s3.ListPartsOutput, error) {
					return createListPartsOutput(&expectedPart), nil
				},
			}

			// Instantiate and call Upload
			s3Cli := s3client.InstantiateClient(sdkMock, nil, bucket, ExpectedRegion, nil)
			err := s3Cli.UploadPart(context.Background(), &s3client.UploadPartRequest{
				UploadKey:   testKey,
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 2,
				FileName:    "helloworld",
			}, payload)

			// Validate
			So(err, ShouldBeNil)
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In1.Bucket, ShouldResemble, bucket)
			So(len(sdkMock.CreateMultipartUploadCalls()), ShouldEqual, 1)
			So(len(sdkMock.UploadPartCalls()), ShouldEqual, 1)
			So(*sdkMock.UploadPartCalls()[0].In1.UploadId, ShouldEqual, testUploadId)
			So(*sdkMock.UploadPartCalls()[0].In1.Bucket, ShouldEqual, bucket)
			So(*sdkMock.UploadPartCalls()[0].In1.Key, ShouldEqual, testKey)
			So(len(sdkMock.ListPartsCalls()), ShouldEqual, 1)
		})

		Convey("If the upload S3 object key cannot be found in the list of multipart uploads, Upload will create a new one, "+
			"and complete it if all chunks have been uploaded", func() {

			// Create S3Client with SDK Mock with empty list of Multipart uploads
			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(in1 *s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error) {
					return &s3.ListMultipartUploadsOutput{}, nil
				},
				CreateMultipartUploadFunc: func(in1 *s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error) {
					return &s3.CreateMultipartUploadOutput{UploadId: &testUploadId}, nil
				},
				UploadPartFunc: func(in1 *s3.UploadPartInput) (*s3.UploadPartOutput, error) {
					return &s3.UploadPartOutput{}, nil
				},
				ListPartsFunc: func(in1 *s3.ListPartsInput) (*s3.ListPartsOutput, error) {
					return createListPartsOutput(&expectedPart), nil
				},
				CompleteMultipartUploadFunc: func(input *s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error) {
					return &s3.CompleteMultipartUploadOutput{}, nil
				},
			}

			// Instantiate and call Upload
			s3Cli := s3client.InstantiateClient(sdkMock, nil, bucket, ExpectedRegion, nil)
			err := s3Cli.UploadPart(context.Background(), &s3client.UploadPartRequest{
				UploadKey:   testKey,
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 1,
				FileName:    "helloworld",
			}, payload)

			// Validate
			So(err, ShouldBeNil)
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In1.Bucket, ShouldResemble, bucket)
			So(len(sdkMock.CreateMultipartUploadCalls()), ShouldEqual, 1)
			So(len(sdkMock.UploadPartCalls()), ShouldEqual, 1)
			So(*sdkMock.UploadPartCalls()[0].In1.UploadId, ShouldEqual, testUploadId)
			So(*sdkMock.UploadPartCalls()[0].In1.Bucket, ShouldEqual, bucket)
			So(*sdkMock.UploadPartCalls()[0].In1.Key, ShouldEqual, testKey)
			So(len(sdkMock.ListPartsCalls()), ShouldEqual, 1)
			So(len(sdkMock.CompleteMultipartUploadCalls()), ShouldEqual, 1)
		})

		Convey("UploadWithPsk performs an upload with the provided PSK", func() {

			psk := []byte("test psk")

			// Create S3Client with SDK Mock with empty list of Multipart uploads
			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(in1 *s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error) {
					return &s3.ListMultipartUploadsOutput{}, nil
				},
				CreateMultipartUploadFunc: func(in1 *s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error) {
					return &s3.CreateMultipartUploadOutput{UploadId: &testUploadId}, nil
				},
				ListPartsFunc: func(in1 *s3.ListPartsInput) (*s3.ListPartsOutput, error) {
					return createListPartsOutput(&expectedPart), nil
				},
			}

			cryptoMock := &mock.S3CryptoClientMock{
				UploadPartWithPSKFunc: func(in1 *s3.UploadPartInput, in2 []byte) (*s3.UploadPartOutput, error) {
					return &s3.UploadPartOutput{}, nil
				},
			}

			// Instantiate and call UploadWithPsk
			s3Cli := s3client.InstantiateClient(sdkMock, cryptoMock, bucket, ExpectedRegion, nil)
			err := s3Cli.UploadPartWithPsk(context.Background(), &s3client.UploadPartRequest{
				UploadKey:   testKey,
				Type:        "text/plain",
				ChunkNumber: 1,
				TotalChunks: 2,
				FileName:    "helloworld",
			}, payload, psk)

			// Validate
			So(err, ShouldBeNil)
			So(len(sdkMock.ListMultipartUploadsCalls()), ShouldEqual, 1)
			So(*sdkMock.ListMultipartUploadsCalls()[0].In1.Bucket, ShouldResemble, bucket)
			So(len(sdkMock.CreateMultipartUploadCalls()), ShouldEqual, 1)
			So(len(cryptoMock.UploadPartWithPSKCalls()), ShouldEqual, 1)
			So(len(sdkMock.ListPartsCalls()), ShouldEqual, 1)
		})
	})
}

func TestCheckUpload(t *testing.T) {

	Convey("Given an S3 client with the intention of checking if a chunk has been uploaded in a multipart upload", t, func() {

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
			s3Cli := s3client.InstantiateClient(sdkMock, nil, bucket, ExpectedRegion, nil)
			ok, err := s3Cli.CheckPartUploaded(context.Background(), &s3client.UploadPartRequest{
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

		Convey("If the upload S3 object key cannot be found in the list of multipart uploads, then the CheckUpload will fail with a ErrNotUploaded error", func() {

			// Create S3Client with SDK Mock with empty list of Multipart uploads
			sdkMock := &mock.S3SDKClientMock{
				ListMultipartUploadsFunc: func(in1 *s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error) {
					return &s3.ListMultipartUploadsOutput{}, nil
				},
			}

			// Instantiate and call CheckUpload
			s3Cli := s3client.InstantiateClient(sdkMock, nil, bucket, ExpectedRegion, nil)
			ok, err := s3Cli.CheckPartUploaded(context.Background(), &s3client.UploadPartRequest{
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
			s3Cli := s3client.InstantiateClient(sdkMock, nil, bucket, ExpectedRegion, nil)
			ok, err := s3Cli.CheckPartUploaded(context.Background(), &s3client.UploadPartRequest{
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

		Convey("If the chunk has been uploaded but the multipart upload is not completed yet, then the function should return true", func() {

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
			s3Cli := s3client.InstantiateClient(sdkMock, nil, bucket, ExpectedRegion, nil)
			ok, err := s3Cli.CheckPartUploaded(context.Background(), &s3client.UploadPartRequest{
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

		Convey("Provided chunk not being found in the list of parts results in ErrChunkNumberNotFound being returned", func() {

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
			s3Cli := s3client.InstantiateClient(sdkMock, nil, bucket, ExpectedRegion, nil)
			ok, err := s3Cli.CheckPartUploaded(context.Background(), &s3client.UploadPartRequest{
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

		Convey("Provided chunk being successfully uploaded as part of a completed multipart upload results in the function completing the upload and returning true", func() {

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
			s3Cli := s3client.InstantiateClient(sdkMock, nil, bucket, ExpectedRegion, nil)
			ok, err := s3Cli.CheckPartUploaded(context.Background(), &s3client.UploadPartRequest{
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

// createListPartsOutput returns a ListPartsOutput with a single part, corresponding to the provided partNumber
func createListPartsOutput(partNumber *int64) *s3.ListPartsOutput {
	parts := make([]*s3.Part, 0, 1)
	parts = append(parts, &s3.Part{
		PartNumber: partNumber,
	})
	return &s3.ListPartsOutput{
		Parts: parts,
	}
}

// createUploads returns a ListMultipartUploadsOutput with a single upload, with the provided uploadID
func createUploads(uploadID, key string) *s3.ListMultipartUploadsOutput {
	uploads := make([]*s3.MultipartUpload, 0, 1)
	uploads = append(uploads, &s3.MultipartUpload{
		UploadId: &uploadID,
		Key:      &key,
	})
	return &s3.ListMultipartUploadsOutput{
		Uploads: uploads,
	}
}
