package s3client

import (
	"context"
	"errors"
	"fmt"

	"github.com/ONSdigital/dp-s3/crypto"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Uploader extends S3 client in order to perform Upload operations easily using the aws s3manager package.
// It allows Uploads with or without user-provided PSK for encryption.
type Uploader struct {
	*S3
	sdkUploader    S3SDKUploader
	cryptoUploader S3CryptoUploader
}

// NewUploader creates a new Uploader configured for the given region and bucket name.
// Note: This function will create a new AWS session, if you already have a valid session, please use NewUploaderWithSession instead
// If an error occurs while establishing the AWS session, it will be returned
func NewUploader(region string, bucketName string) (*Uploader, error) {
	s, err := session.NewSession(&aws.Config{Region: &region})
	if err != nil {
		return nil, err
	}
	return NewUploaderWithSession(bucketName, s), nil
}

// NewUploaderWithSession creates a new Uploader configured for the given bucket name, using the provided session and regions defined by it.
func NewUploaderWithSession(bucketName string, s *session.Session) *Uploader {
	// Create base S3 client using the provided session
	s3Client := NewClientWithSession(bucketName, s)

	// Create an AWS-SDK-S3 Uploader.
	sdkUploader := s3manager.NewUploader(s)

	// Create crypto client, which allows user to provide a psk
	cryptoUploader := crypto.NewUploader(s, &crypto.Config{HasUserDefinedPSK: true})

	return InstantiateUploader(s3Client, sdkUploader, cryptoUploader)
}

// InstantiateUploader creates a new instance of Uploader struct with the provided clients, bucket and region
func InstantiateUploader(s3Client *S3, sdkUploader S3SDKUploader, cryptoUploader S3CryptoUploader) *Uploader {
	return &Uploader{s3Client, sdkUploader, cryptoUploader}
}

// Session returns the Session of this uploader
func (u *Uploader) Session() *session.Session {
	return u.session
}

// Upload uploads a file to S3 using the AWS s3Manager, which will automatically split up large objects and upload them concurrently.
func (u *Uploader) Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	logData, err := u.ValidateInput(input)
	if err != nil {
		return nil, NewError(
			fmt.Errorf("validation error for Upload: %w", err),
			logData,
		)
	}

	output, err := u.sdkUploader.Upload(input, options...)
	if err != nil {
		return nil, NewError(
			fmt.Errorf("failed to upload: %w", err),
			logData,
		)
	}
	return output, nil
}

// Upload uploads a file to S3 using the AWS s3Manager with context, which will automatically split up large objects and upload them concurrently.
// The provided context may be used to abort the operation.
func (u *Uploader) UploadWithContext(ctx context.Context, input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	logData, err := u.ValidateInput(input)
	if err != nil {
		return nil, NewError(
			fmt.Errorf("validation error for UploadWithContext: %w", err),
			logData,
		)
	}
	if ctx == nil {
		return nil, NewError(
			errors.New("nil context provided to UploadWithContext"),
			logData,
		)
	}

	output, err := u.sdkUploader.UploadWithContext(ctx, input, options...)
	if err != nil {
		return nil, NewError(
			fmt.Errorf("failed to upload with context: %w", err),
			logData,
		)
	}
	return output, nil
}

// UploadWithPSK uploads a file to S3 using cryptoclient, which allows you to encrypt the file with a given psk.
func (u *Uploader) UploadWithPSK(input *s3manager.UploadInput, psk []byte) (*s3manager.UploadOutput, error) {
	logData, err := u.ValidateInput(input)
	if err != nil {
		return nil, NewError(
			fmt.Errorf("validation error for UploadWithPSK: %w", err),
			logData,
		)
	}
	if len(psk) == 0 {
		return nil, NewError(
			errors.New("nil or empty psk provided to UploadWithPSK"),
			logData,
		)
	}

	output, err := u.cryptoUploader.UploadWithPSK(nil, input, psk)
	if err != nil {
		return nil, NewError(
			fmt.Errorf("failed to upload with psk: %w", err),
			logData,
		)
	}
	return output, nil
}

// UploadWithPSKAndContext uploads a file to S3 using cryptoclient, which allows you to encrypt the file with a given psk.
// The provided context may be used to abort the operation.
func (u *Uploader) UploadWithPSKAndContext(ctx context.Context, input *s3manager.UploadInput, psk []byte, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	logData, err := u.ValidateInput(input)
	if err != nil {
		return nil, NewError(
			fmt.Errorf("validation error for UploadWithPSKAndContext: %w", err),
			logData,
		)
	}
	if ctx == nil {
		return nil, NewError(
			errors.New("nil context provided to UploadWithPSKAndContext"),
			logData,
		)
	}
	if len(psk) == 0 {
		return nil, NewError(
			errors.New("nil or empty psk provided to UploadWithPSKAndContext"),
			logData,
		)
	}

	output, err := u.cryptoUploader.UploadWithPSK(ctx, input, psk)
	if err != nil {
		return nil, NewError(
			fmt.Errorf("failed to upload with psk: %w", err),
			logData,
		)
	}
	return output, nil
}

// ValidateInput checks the input and returns an error
// if there is a bucket override mismatch or s3 key is not provided
func (u *Uploader) ValidateInput(input *s3manager.UploadInput) (log.Data, error) {
	logData := log.Data{
		"bucket_name": u.bucketName,
	}

	if input == nil {
		return logData, errors.New("nil input provided")
	}

	if input.Key == nil || len(*input.Key) == 0 {
		return logData, errors.New("nil or empty s3 key provided in input")
	}
	logData["s3_key"] = *input.Key // key is the s3 filename with path (it's not a cryptographic key)

	if input.Bucket == nil {
		input.Bucket = &u.bucketName
	} else if *input.Bucket != u.bucketName {
		logData["input_bucket_name"] = *input.Bucket
		return logData, errors.New("unexpected bucket name provided in upload input")
	}

	return logData, nil
}
