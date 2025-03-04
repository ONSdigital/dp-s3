// file: upload.go
//
// Contains methods to efficiently upload files to S3
// by using the high level SDK manager uploader methods,
// which automatically split large objects in chunks and uploads them concurrently.
//
// Requires "s3:PutObject" action allowed by IAM policy for the bucket,
// as defined by `write-{bucketName}-bucket` policies in dp-setup
package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// PutWithPSK uploads the provided contents to the key in the bucket configured for this client, using the provided PSK.
// The 'key' parameter refers to the path for the file under the bucket.
func (cli *Client) PutWithPSK(ctx context.Context, key *string, reader *bytes.Reader, psk []byte) error {
	input := &s3.PutObjectInput{
		Body:   reader,
		Key:    key,
		Bucket: &cli.bucketName,
	}

	if _, err := cli.cryptoClient.PutObjectWithPSK(ctx, input, psk); err != nil {
		return NewError(fmt.Errorf("error putting object to s3: %w", err), log.Data{
			"bucket_name": cli.bucketName,
			"s3_key":      key, // key is the s3 filename with path (it's not a cryptographic key)
			"user_psk":    true,
		})
	}
	return nil
}

// Upload uploads a file to S3 using the AWS Manager, which will automatically split up large objects and upload them concurrently.
func (cli *Client) Upload(ctx context.Context, input *s3.PutObjectInput, options ...func(*manager.Uploader)) (*manager.UploadOutput, error) {
	logData, err := cli.ValidateUploadInput(input)
	if err != nil {
		return nil, NewError(
			fmt.Errorf("validation error for Upload: %w", err),
			logData,
		)
	}

	output, err := cli.sdkUploader.Upload(ctx, input, options...)
	if err != nil {
		return nil, NewError(
			fmt.Errorf("failed to upload: %w", err),
			logData,
		)
	}
	return output, nil
}

// UploadWithPSK uploads a file to S3 using cryptoclient, which allows you to encrypt the file with a given psk.
func (cli *Client) UploadWithPSK(ctx context.Context, input *s3.PutObjectInput, psk []byte) (*manager.UploadOutput, error) {
	logData, err := cli.ValidateUploadInput(input)
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

	output, err := cli.cryptoUploader.UploadWithPSK(ctx, input, psk)
	if err != nil {
		return nil, NewError(
			fmt.Errorf("failed to upload with psk: %w", err),
			logData,
		)
	}
	return output, nil
}

// ValidateUploadInput checks the upload input and returns an error
// if there is a bucket override mismatch or s3 key is not provided
func (cli *Client) ValidateUploadInput(input *s3.PutObjectInput) (log.Data, error) {
	logData := log.Data{
		"bucket_name": cli.bucketName,
	}

	if input == nil {
		return logData, errors.New("nil input provided")
	}

	if input.Key == nil || len(*input.Key) == 0 {
		return logData, errors.New("nil or empty s3 key provided in input")
	}
	logData["s3_key"] = *input.Key // key is the s3 filename with path (it's not a cryptographic key)

	if input.Bucket == nil {
		input.Bucket = &cli.bucketName
	} else if *input.Bucket != cli.bucketName {
		logData["input_bucket_name"] = *input.Bucket
		return logData, errors.New("unexpected bucket name provided in upload input")
	}

	return logData, nil
}
