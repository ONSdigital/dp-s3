package s3client

import (
	"github.com/ONSdigital/s3crypto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Uploader with SDK uploader client, CryptoClient uploader client, region and BucketName
type Uploader struct {
	sdkUploader    S3SDKUploader
	cryptoUploader S3CryptoUploader
	bucketName     string
	region         string
	session        *session.Session
}

// NewUploader creates a new Uploader configured for the given region and bucket name.
// If hasUserDefinedPSK is true, it will be a crypto uploader client, otherwise it will be a vanilla S3 SDK Uploader.
// Note: This function will create a new session, if you already have a session, please use NewUploaderWithSession instead
func NewUploader(region string, bucketName string, hasUserDefinedPSK bool) (*Uploader, error) {
	s, err := session.NewSession(&aws.Config{Region: &region})
	if err != nil {
		return nil, err
	}
	return NewUploaderWithSession(region, bucketName, hasUserDefinedPSK, s), nil
}

// NewUploaderWithSession creates a new Uploader configured for the given region and bucket name, using the provided session.
// If hasUserDefinedPSK is true, it will be a crypto uploader client, otherwise it will be a vanilla S3 SDK Uploader.
func NewUploaderWithSession(region string, bucketName string, hasUserDefinedPSK bool, s *session.Session) *Uploader {

	// If we require crypto client (HasUserDefinedPSK), create it.
	if hasUserDefinedPSK {
		cryptoUploader := s3crypto.NewUploader(s, &s3crypto.Config{HasUserDefinedPSK: true})
		return InstantiateUploader(nil, cryptoUploader, bucketName, region, s)
	}

	// Otherwise create a vanilla AWS-SDK-S3 Uploader.
	sdkUploader := s3manager.NewUploader(s)
	return InstantiateUploader(sdkUploader, nil, bucketName, region, s)
}

// InstantiateUploader creates a new instance of Uploader struct with the provided clients, bucket and region
func InstantiateUploader(sdkUploader S3SDKUploader, cryptoUploader S3CryptoUploader, bucketName, region string, s *session.Session) *Uploader {
	return &Uploader{
		sdkUploader:    sdkUploader,
		cryptoUploader: cryptoUploader,
		bucketName:     bucketName,
		region:         region,
		session:        s,
	}
}

// Session returns the Session of this uploader
func (u *Uploader) Session() *session.Session {
	return u.session
}

// UploadWithPSK uploads a file to S3 using cryptoclient, which allows you to encrypt the file with a given psk
func (u *Uploader) UploadWithPSK(input *s3manager.UploadInput, psk []byte) (*s3manager.UploadOutput, error) {
	if u.cryptoUploader == nil {
		return nil, &ErrInvalidUploader{expectCrypto: true}
	}
	return u.cryptoUploader.UploadWithPSK(input, psk)
}

// Upload uploads a file to S3 using the AWS s3Manager, which will automatically split up large objects and upload them concurrently
func (u *Uploader) Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	if u.cryptoUploader == nil {
		return nil, &ErrInvalidUploader{expectCrypto: false}
	}
	return u.sdkUploader.Upload(input, options...)
}
