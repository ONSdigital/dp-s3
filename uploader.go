package s3client

import (
	"errors"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/ONSdigital/s3crypto"
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
	cryptoUploader := s3crypto.NewUploader(s, &s3crypto.Config{HasUserDefinedPSK: true})

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

// UploadWithPSK uploads a file to S3 using cryptoclient, which allows you to encrypt the file with a given psk
func (u *Uploader) UploadWithPSK(input *s3manager.UploadInput, psk []byte) (*s3manager.UploadOutput, error) {
	logData := log.Data{
		"bucket_name": u.bucketName,
	}

	// param validation
	if input == nil {
		return nil, NewError(errors.New("nil input provided to UploadWithPSK"), logData)
	}
	logData["s3_key"] = input.Key // key is the s3 filename with path (it's not a cryptographic key)
	if psk == nil || len(psk) == 0 {
		return nil, NewError(errors.New("nil or empty psk provided to UploadWithPSK"), logData)
	}

	// Check that the requested Bucket is the correct one, or assign it if nil
	if err := u.validateRequestBucket(input); err != nil {
		return nil, err
	}

	// Perform the Upload with SPK
	output, err := u.cryptoUploader.UploadWithPSK(input, psk)
	if err != nil {
		return nil, NewError(err, logData)
	}
	return output, nil
}

// Upload uploads a file to S3 using the AWS s3Manager, which will automatically split up large objects and upload them concurrently
func (u *Uploader) Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	logData := log.Data{
		"bucket_name": u.bucketName,
	}

	// param validation
	if input == nil {
		return nil, NewError(errors.New("nil input provided to UploadWithPSK"), logData)
	}
	logData["s3_key"] = input.Key // key is the s3 filename with path (it's not a cryptographic key)

	// Check that the requested Bucket is the correct one, or assign it if nil
	if err := u.validateRequestBucket(input); err != nil {
		return nil, err
	}

	// Perform the Upload using the AWS SDK Uploader
	return u.sdkUploader.Upload(input, options...)
}

// validateRequestBucket checks that the requested bucket matches the bucket configured in this client, or assigns it if it is nil
func (u *Uploader) validateRequestBucket(input *s3manager.UploadInput) error {
	if input.Bucket == nil {
		input.Bucket = &u.bucketName
	} else if *input.Bucket != u.bucketName {
		return NewError(
			errors.New("unexpected bucket name provided in upload input"),
			log.Data{
				"client_bucket_name": u.bucketName,
				"input_bucket_name":  input.Bucket,
			})
	}
	return nil
}
