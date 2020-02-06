package s3client

import (
	"github.com/aws/aws-sdk-go/service/s3"
)

//go:generate moq -out ./mock/s3-sdk.go -pkg mock . S3SDKClient
//go:generate moq -out ./mock/s3-crypto.go -pkg mock . S3CryptoClient

// S3SDKClient represents the client with methods required to upload a multipart upload to s3
type S3SDKClient interface {
	ListMultipartUploads(*s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error)
	ListParts(*s3.ListPartsInput) (*s3.ListPartsOutput, error)
	CompleteMultipartUpload(input *s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error)
	CreateMultipartUpload(*s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error)
	UploadPart(*s3.UploadPartInput) (*s3.UploadPartOutput, error)
	ListObjectsV2(input *s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error)
	GetObject(input *s3.GetObjectInput) (*s3.GetObjectOutput, error)
}

// S3CryptoClient represents the cryptoclient with methods required to upload parts with encryption
type S3CryptoClient interface {
	UploadPartWithPSK(*s3.UploadPartInput, []byte) (*s3.UploadPartOutput, error)
}
