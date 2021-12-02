// file: client.go
//
// Contains the Client struct definition and constructors,
// as well as getters to read some private fields like bucketName or session.
//
// If multiple clients are required, it is advised to reuse the same AWS session.
package s3

import (
	"fmt"
	"sync"

	"github.com/ONSdigital/dp-s3/v2/crypto"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Client client with SDK client, CryptoClient and BucketName
type Client struct {
	sdkClient      S3SDKClient
	cryptoClient   S3CryptoClient
	sdkUploader    S3SDKUploader
	cryptoUploader S3CryptoUploader
	bucketName     string
	region         string
	mutexUploadID  *sync.Mutex
	session        *session.Session
}

// NewClient creates a new S3 Client configured for the given region and bucket name.
// Note: This function will create a new session, if you already have a session, please use NewUploaderWithSession instead
// Any error establishing the AWS session will be returned
func NewClient(region string, bucketName string) (*Client, error) {
	s, err := session.NewSession(&aws.Config{Region: &region})
	if err != nil {
		return nil, NewError(
			fmt.Errorf("error creating session: %w", err),
			log.Data{
				"region":      region,
				"bucket_name": bucketName,
			},
		)
	}
	return NewClientWithSession(bucketName, s), nil
}

// NewClientWithSession creates a new S3 Client configured for the given bucket name, using the provided session and region within it.
func NewClientWithSession(bucketName string, s *session.Session) *Client {
	// Get region for the Session config
	region := s.Config.Region

	// Create AWS-SDK-S3 client with the session
	sdkClient := s3.New(s)

	// Create an AWS-SDK-S3 Uploader.
	sdkUploader := s3manager.NewUploader(s)

	// Create crypto client, which allows user to provide a psk
	cryptoClient := crypto.New(s, &crypto.Config{HasUserDefinedPSK: true})

	// Create crypto uploader, which allows user to provide a psk
	cryptoUploader := crypto.NewUploader(s, &crypto.Config{HasUserDefinedPSK: true})

	return InstantiateClient(sdkClient, cryptoClient, sdkUploader, cryptoUploader, bucketName, *region, s)
}

// InstantiateClient creates a new instance of S3 struct with the provided clients, bucket and region.
func InstantiateClient(sdkClient S3SDKClient, cryptoClient S3CryptoClient, sdkUploader S3SDKUploader, cryptoUploader S3CryptoUploader, bucketName, region string, s *session.Session) *Client {
	return &Client{
		sdkClient:      sdkClient,
		cryptoClient:   cryptoClient,
		sdkUploader:    sdkUploader,
		cryptoUploader: cryptoUploader,
		bucketName:     bucketName,
		region:         region,
		mutexUploadID:  &sync.Mutex{},
		session:        s,
	}
}

// Session returns the Session of this client
func (cli *Client) Session() *session.Session {
	return cli.session
}

// BucketName returns the bucket name used by this S3 client
func (cli *Client) BucketName() string {
	return cli.bucketName
}
