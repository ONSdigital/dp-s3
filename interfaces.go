package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//go:generate moq -out ./mock/s3-sdk.go -pkg mock . S3SDKClient
//go:generate moq -out ./mock/s3-crypto.go -pkg mock . S3CryptoClient
//go:generate moq -out ./mock/s3-uploader.go -pkg mock . S3SDKUploader
//go:generate moq -out ./mock/s3-crypto-uploader.go -pkg mock . S3CryptoUploader

// S3SDKClient represents the sdk client with methods required by dp-s3 client
type S3SDKClient interface {
	ListMultipartUploads(ctx context.Context, in *s3.ListMultipartUploadsInput, optFns ...func(*s3.Options)) (*s3.ListMultipartUploadsOutput, error)
	ListParts(ctx context.Context, in *s3.ListPartsInput, optFns ...func(*s3.Options)) (*s3.ListPartsOutput, error)
	CompleteMultipartUpload(ctx context.Context, in *s3.CompleteMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CompleteMultipartUploadOutput, error)
	CreateMultipartUpload(ctx context.Context, in *s3.CreateMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CreateMultipartUploadOutput, error)
	UploadPart(ctx context.Context, in *s3.UploadPartInput, optFns ...func(*s3.Options)) (*s3.UploadPartOutput, error)
	HeadBucket(ctx context.Context, in *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error)
	HeadObject(ctx context.Context, in *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	GetObject(ctx context.Context, in *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	GetBucketPolicy(ctx context.Context, in *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error)
	PutBucketPolicy(ctx context.Context, in *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error)
	ListObjects(ctx context.Context, in *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error)
}

// S3CryptoClient represents the cryptoclient with methods required to upload parts with encryption
type S3CryptoClient interface {
	UploadPartWithPSK(ctx context.Context, in *s3.UploadPartInput, psk []byte) (out *s3.UploadPartOutput, err error)
	GetObjectWithPSK(ctx context.Context, in *s3.GetObjectInput, psk []byte) (out *s3.GetObjectOutput, err error)
	PutObjectWithPSK(ctx context.Context, in *s3.PutObjectInput, psk []byte) (out *s3.PutObjectOutput, err error)
}

// S3SDKUploader represents the sdk uploader with methods required by dp-s3 client
type S3SDKUploader interface {
	Upload(ctx context.Context, in *s3.PutObjectInput, options ...func(*manager.Uploader)) (out *manager.UploadOutput, err error)
}

// S3CryptoUploader represents the s3crypto Uploader with methods required to upload parts with encryption
type S3CryptoUploader interface {
	UploadWithPSK(ctx context.Context, in *s3.PutObjectInput, psk []byte) (out *manager.UploadOutput, err error)
}
