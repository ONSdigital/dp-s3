package s3client

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//go:generate moq -out ./mock/s3-sdk.go -pkg mock . S3SDKClient
//go:generate moq -out ./mock/s3-crypto.go -pkg mock . S3CryptoClient
//go:generate moq -out ./mock/s3-uploader.go -pkg mock . S3SDKUploader
//go:generate moq -out ./mock/s3-crypto-uploader.go -pkg mock . S3CryptoUploader

// S3SDKClient represents the sdk client with methods required by dp-s3 client
type S3SDKClient interface {
	ListMultipartUploads(*s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error)
	ListParts(*s3.ListPartsInput) (*s3.ListPartsOutput, error)
	CompleteMultipartUpload(*s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error)
	CreateMultipartUpload(*s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error)
	UploadPart(*s3.UploadPartInput) (*s3.UploadPartOutput, error)
	ListObjectsV2(*s3.ListObjectsV2Input) (*s3.ListObjectsV2Output, error)
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
}

// S3CryptoClient represents the cryptoclient with methods required to upload parts with encryption
type S3CryptoClient interface {
	UploadPartWithPSK(*s3.UploadPartInput, []byte) (*s3.UploadPartOutput, error)
	GetObjectWithPSK(*s3.GetObjectInput, []byte) (*s3.GetObjectOutput, error)
}

// S3SDKUploader represents the sdk uploader with methods required by dp-s3 client
type S3SDKUploader interface {
	Upload(*s3manager.UploadInput, ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

// S3CryptoUploader represents the s3crypto Uploader with methods required to upload parts with encryption
type S3CryptoUploader interface {
	UploadWithPSK(*s3manager.UploadInput, []byte) (*s3manager.UploadOutput, error)
}
