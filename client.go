// file: client.go
//
// Contains the Client struct definition and constructors,
// as well as getters to read some private fields like bucketName or cfg.
//
// If multiple clients are required, it is advised to reuse the same AWS config.
package s3

import (
	"context"
	"fmt"
	"sync"

	"github.com/ONSdigital/dp-s3/v2/crypto"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Client: client with sdkClient, cryptoClient, sdkUploader, cryptoUploader, bucketName, region, mutexUploadID and cfg
type Client struct {
	sdkClient      S3SDKClient
	cryptoClient   S3CryptoClient
	sdkUploader    S3SDKUploader
	cryptoUploader S3CryptoUploader
	bucketName     string
	region         string
	mutexUploadID  *sync.Mutex
	cfg            aws.Config
}

// NewClient creates a new S3 Client configured for the given region and bucket name.
// Note: This function will create a new config, if you already have a config, please use NewUploader instead
// Any error establishing the AWS config will be returned
func NewClient(ctx context.Context, region string, bucketName string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, NewError(
			fmt.Errorf("error creating config: %w", err),
			log.Data{
				"region":      region,
				"bucket_name": bucketName,
			},
		)
	}
	return NewClientWithConfig(bucketName, cfg), nil
}

// NewClientWithCredentials creates a new S3 Client configured for the given region and bucket name with creds.
// Note: This function will create a new config, if you already have a config, please use NewUploader instead
// Any error establishing the AWS config will be returned
func NewClientWithCredentials(ctx context.Context, region string, bucketName string, awsAccessKey string, awsSecretKey string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKey, awsSecretKey, "")),
	)
	if err != nil {
		return nil, NewError(
			fmt.Errorf("error creating config: %w", err),
			log.Data{
				"region":      region,
				"bucket_name": bucketName,
			},
		)
	}
	return NewClientWithConfig(bucketName, cfg), nil
}

// NewClientWithConfig creates a new S3 Client configured for the given bucket name, using the provided config and region within it.
func NewClientWithConfig(bucketName string, cfg aws.Config, optFns ...func(*s3.Options)) *Client {
	// Get region for the Config
	region := cfg.Region

	// Create AWS-SDK-S3 client with the config
	sdkClient := s3.NewFromConfig(cfg, optFns...)

	// Create an AWS-SDK-S3 Uploader.
	sdkUploader := manager.NewUploader(sdkClient)

	// Create crypto client, which allows user to provide a psk
	cryptoClient := crypto.New(cfg, &crypto.Config{HasUserDefinedPSK: true})

	// Create crypto uploader, which allows user to provide a psk
	cryptoUploader := crypto.NewUploader(cfg, &crypto.Config{HasUserDefinedPSK: true})

	return InstantiateClient(sdkClient, cryptoClient, sdkUploader, cryptoUploader, bucketName, region, cfg)
}

// InstantiateClient creates a new instance of S3 struct with the provided clients, bucket and region.
func InstantiateClient(sdkClient S3SDKClient, cryptoClient S3CryptoClient, sdkUploader S3SDKUploader, cryptoUploader S3CryptoUploader, bucketName, region string, cfg aws.Config) *Client {
	return &Client{
		sdkClient:      sdkClient,
		cryptoClient:   cryptoClient,
		sdkUploader:    sdkUploader,
		cryptoUploader: cryptoUploader,
		bucketName:     bucketName,
		region:         region,
		mutexUploadID:  &sync.Mutex{},
		cfg:            cfg,
	}
}

// Config returns the Config of this client
func (cli *Client) Config() aws.Config {
	return cli.cfg
}

// BucketName returns the bucket name used by this S3 client
func (cli *Client) BucketName() string {
	return cli.bucketName
}
