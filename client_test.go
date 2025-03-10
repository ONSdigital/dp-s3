package s3_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	dps3 "github.com/ONSdigital/dp-s3/v3"
	"github.com/ONSdigital/dp-s3/v3/mock"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPutWithPSK(t *testing.T) {
	Convey("Given an S3 client configured with a bucket, region and psk", t, func() {
		ctx := context.Background()

		psk := []byte("test psk")
		payload := []byte("test data")
		bucket := "myBucket"
		objKey := "my/object/key"
		region := "eu-north-1"

		payloadReader := bytes.NewReader(payload)

		cryptoMock := &mock.S3CryptoClientMock{
			PutObjectWithPSKFunc: func(ctx context.Context, in1 *s3.PutObjectInput, in2 []byte) (*s3.PutObjectOutput, error) {
				return &s3.PutObjectOutput{}, nil
			},
		}

		cli := dps3.InstantiateClient(nil, cryptoMock, nil, nil, bucket, region, aws.Config{})

		Convey("PutWithPSK calls the expected cryptoClient with provided key, reader and client-configured bucket", func() {
			err := cli.PutWithPSK(ctx, &objKey, payloadReader, psk)
			So(err, ShouldBeNil)
			So(len(cryptoMock.PutObjectWithPSKCalls()), ShouldEqual, 1)
			So(cryptoMock.PutObjectWithPSKCalls()[0].In, ShouldResemble, &s3.PutObjectInput{
				Bucket: &bucket,
				Key:    &objKey,
				Body:   payloadReader,
			})
			So(cryptoMock.PutObjectWithPSKCalls()[0].Psk, ShouldResemble, psk)
		})
	})
}

// readBytes reads the bytes from the provided ReadCloser and asserts that there is no error
func readBytes(ret io.ReadCloser) []byte {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(ret)
	So(err, ShouldBeNil)
	err = ret.Close()
	So(err, ShouldBeNil)
	return buf.Bytes()
}
