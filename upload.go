// file: upload.go
//
// Contains methods to efficiently upload files to S3
// by using the high level SDK s3manager uploader methods,
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
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// PutWithPSK uploads the provided contents to the key in the bucket configured for this client, using the provided PSK.
// The 'key' parameter refers to the path for the file under the bucket.
func (cli *Client) PutWithPSK(key *string, reader *bytes.Reader, psk []byte) error {
	input := &s3.PutObjectInput{
		Body:   reader,
		Key:    key,
		Bucket: &cli.bucketName,
	}

	if _, err := cli.cryptoClient.PutObjectWithPSK(input, psk); err != nil {
		return NewError(fmt.Errorf("error putting object to s3: %w", err), log.Data{
			"bucket_name": cli.bucketName,
			"s3_key":      key, // key is the s3 filename with path (it's not a cryptographic key)
			"user_psk":    true,
		})
	}
	return nil
}

// Upload uploads a file to S3 using the AWS s3Manager, which will automatically split up large objects and upload them concurrently.
func (cli *Client) Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	logData, err := cli.ValidateUploadInput(input)
	if err != nil {
		return nil, NewError(
			fmt.Errorf("validation error for Upload: %w", err),
			logData,
		)
	}

	output, err := cli.sdkUploader.Upload(input, options...)
	if err != nil {
		return nil, NewError(
			fmt.Errorf("failed to upload: %w", err),
			logData,
		)
	}
	return output, nil
}

// UploadWithContext uploads a file to S3 using the AWS s3Manager with context, which will automatically split up large objects and upload them concurrently.
// The provided context may be used to abort the operation.
func (cli *Client) UploadWithContext(ctx context.Context, input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	logData, err := cli.ValidateUploadInput(input)
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

	output, err := cli.sdkUploader.UploadWithContext(ctx, input, options...)
	if err != nil {
		return nil, NewError(
			fmt.Errorf("failed to upload with context: %w", err),
			logData,
		)
	}
	return output, nil
}

// UploadWithPSK uploads a file to S3 using cryptoclient, which allows you to encrypt the file with a given psk.
func (cli *Client) UploadWithPSK(input *s3manager.UploadInput, psk []byte) (*s3manager.UploadOutput, error) {
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

	output, err := cli.cryptoUploader.UploadWithPSK(nil, input, psk)
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
func (cli *Client) UploadWithPSKAndContext(ctx context.Context, input *s3manager.UploadInput, psk []byte, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	logData, err := cli.ValidateUploadInput(input)
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
func (cli *Client) ValidateUploadInput(input *s3manager.UploadInput) (log.Data, error) {
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
