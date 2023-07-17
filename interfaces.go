package s3

import (
	"context"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//go:generate moq -out ./mock/s3-sdk.go -pkg mock . S3SDKClient
//go:generate moq -out ./mock/s3-crypto.go -pkg mock . S3CryptoClient
//go:generate moq -out ./mock/s3-uploader.go -pkg mock . S3SDKUploader
//go:generate moq -out ./mock/s3-crypto-uploader.go -pkg mock . S3CryptoUploader

// S3SDKClient represents the sdk client with methods required by dp-s3 client
type S3SDKClient interface {
	ListMultipartUploads(in *s3.ListMultipartUploadsInput) (out *s3.ListMultipartUploadsOutput, err error)
	ListParts(in *s3.ListPartsInput) (out *s3.ListPartsOutput, err error)
	CompleteMultipartUpload(in *s3.CompleteMultipartUploadInput) (out *s3.CompleteMultipartUploadOutput, err error)
	CreateMultipartUpload(in *s3.CreateMultipartUploadInput) (out *s3.CreateMultipartUploadOutput, err error)
	UploadPart(in *s3.UploadPartInput) (out *s3.UploadPartOutput, err error)
	HeadBucket(in *s3.HeadBucketInput) (out *s3.HeadBucketOutput, err error)
	HeadObject(in *s3.HeadObjectInput) (out *s3.HeadObjectOutput, err error)
	GetObject(in *s3.GetObjectInput) (out *s3.GetObjectOutput, err error)
	GetBucketPolicy(in *s3.GetBucketPolicyInput) (out *s3.GetBucketPolicyOutput, err error)
}

// S3CryptoClient represents the cryptoclient with methods required to upload parts with encryption
type S3CryptoClient interface {
	UploadPartWithPSK(in *s3.UploadPartInput, psk []byte) (out *s3.UploadPartOutput, err error)
	GetObjectWithPSK(in *s3.GetObjectInput, psk []byte) (out *s3.GetObjectOutput, err error)
	PutObjectWithPSK(in *s3.PutObjectInput, psk []byte) (out *s3.PutObjectOutput, err error)
}

// S3SDKUploader represents the sdk uploader with methods required by dp-s3 client
type S3SDKUploader interface {
	Upload(in *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (out *s3manager.UploadOutput, err error)
	UploadWithContext(ctx context.Context, in *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (out *s3manager.UploadOutput, err error)
}

// S3CryptoUploader represents the s3crypto Uploader with methods required to upload parts with encryption
type S3CryptoUploader interface {
	UploadWithPSK(ctx context.Context, in *s3manager.UploadInput, psk []byte) (out *s3manager.UploadOutput, err error)
}
