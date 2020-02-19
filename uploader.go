package s3client

import (
	"github.com/ONSdigital/s3crypto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Uploader extends S3 client, with an additional SDK uploader client or CryptoClient uploader client in order to perform Upload operations.
type Uploader struct {
	*S3
	sdkUploader    S3SDKUploader
	cryptoUploader S3CryptoUploader
}

// NewUploader creates a new Uploader configured for the given region and bucket name.
// If hasUserDefinedPSK is true, it will be a crypto uploader client, otherwise it will be a vanilla S3 SDK Uploader.
// Note: This function will create a new AWS session, if you already have a valid session, please use NewUploaderWithSession instead
func NewUploader(region string, bucketName string, hasUserDefinedPSK bool) (*Uploader, error) {
	s, err := session.NewSession(&aws.Config{Region: &region})
	if err != nil {
		return nil, err
	}
	return NewUploaderWithSession(bucketName, hasUserDefinedPSK, s), nil
}

// NewUploaderWithSession creates a new Uploader configured for the given bucket name, using the provided session and region within it.
// If hasUserDefinedPSK is true, it will be a crypto uploader client, otherwise it will be a vanilla S3 SDK Uploader.
func NewUploaderWithSession(bucketName string, hasUserDefinedPSK bool, s *session.Session) *Uploader {

	// Create base S3 client using the provided session
	s3Client := NewClientWithSession(bucketName, hasUserDefinedPSK, s)

	// If we require crypto client (HasUserDefinedPSK), create it.
	if hasUserDefinedPSK {
		cryptoUploader := s3crypto.NewUploader(s, &s3crypto.Config{HasUserDefinedPSK: true})
		return InstantiateUploader(s3Client, nil, cryptoUploader)
	}

	// Otherwise create a vanilla AWS-SDK-S3 Uploader.
	sdkUploader := s3manager.NewUploader(s)
	return InstantiateUploader(s3Client, sdkUploader, nil)
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

	// We need CryptoUploader to perform this operation
	if u.cryptoUploader == nil {
		return nil, &ErrInvalidUploader{ExpectCrypto: true}
	}

	// Check that the requested Bucket is the correct one, or assign it if nil
	if err := u.validateRequestBucket(input); err != nil {
		return nil, err
	}

	// Perform the Upload with SPK
	return u.cryptoUploader.UploadWithPSK(input, psk)
}

// Upload uploads a file to S3 using the AWS s3Manager, which will automatically split up large objects and upload them concurrently
func (u *Uploader) Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {

	// We need SDKUploader to perform this operation
	if u.sdkUploader == nil {
		return nil, &ErrInvalidUploader{ExpectCrypto: false}
	}

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
		return &ErrUnexpectedBucket{
			BucketName:         *input.Bucket,
			ExpectedBucketName: u.bucketName,
		}
	}
	return nil
}
