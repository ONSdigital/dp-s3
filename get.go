// file: get.go
//
// Contains methods to get objects or metadata from S3,
// with or without a used-defined psk for encryption,
// passing the object key or a full path in a specific aws allowed style.
//
// Requires "s3:GetObject" action allowed by IAM policy for objects inside the bucket,
// as defined by `read-{bucketName}-bucket` policies in dp-setup
package s3

import (
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// GetFromS3URL returns an io.ReadCloser instance and the content length (size in bytes) for the given S3 URL,
// in the format specified by URLStyle.
// More information about s3 URL styles: https://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html
// If the URL defines a region (if provided) or bucket different from the one configured in this client, an error will be returned.
//
// The caller is responsible for closing the returned ReadCloser.
// For example, it may be closed in a defer statement: defer r.Close()
func (cli *Client) GetFromS3URL(rawURL string, style URLStyle) (io.ReadCloser, *int64, error) {
	return cli.doGetFromS3URL(rawURL, style, nil)

}

// GetFromS3URLWithPSK returns an io.ReadCloser instance and the content length (size in bytes) for the given S3 URL,
// in the format specified by URLStyle, using the provided PSK for encryption.
// More information about s3 URL styles: https://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html
// If the URL defines a region (if provided) or bucket different from the one configured in this client, an error will be returned.
//
// The caller is responsible for closing the returned ReadCloser.
// For example, it may be closed in a defer statement: defer r.Close()
func (cli *Client) GetFromS3URLWithPSK(rawURL string, style URLStyle, psk []byte) (io.ReadCloser, *int64, error) {
	return cli.doGetFromS3URL(rawURL, style, psk)
}

func (cli *Client) doGetFromS3URL(rawURL string, style URLStyle, psk []byte) (io.ReadCloser, *int64, error) {
	logData := log.Data{
		"raw_url":   rawURL,
		"url_style": style.String(),
	}

	// Parse URL with the provided format style
	s3Url, err := ParseURL(rawURL, style)
	if err != nil {
		return nil, nil, NewError(fmt.Errorf("error parsing url: %w", err), logData)
	}

	// Validate that URL and client bucket names match
	if s3Url.BucketName != cli.bucketName {
		logData["bucket_name"] = cli.bucketName
		return nil, nil, NewUnexpectedBucketError(errors.New("unexpected bucket name in url"), logData)
	}

	// Validate that URL and client regions match, if URL provides one
	if len(s3Url.Region) > 0 && s3Url.Region != cli.region {
		logData["region"] = cli.region
		return nil, nil, NewUnexpectedRegionError(errors.New("unexpected aws region in url"), logData)
	}

	if psk == nil {
		return cli.Get(s3Url.Key)
	}
	return cli.GetWithPSK(s3Url.Key, psk)
}

// Get returns an io.ReadCloser instance for the given path (inside the bucket configured for this client)
// and the content length (size in bytes).
// They 'key' parameter refers to the path for the file under the bucket.
//
// The caller is responsible for closing the returned ReadCloser.
// For example, it may be closed in a defer statement: defer r.Close()
func (cli *Client) Get(key string) (io.ReadCloser, *int64, error) {

	input := &s3.GetObjectInput{
		Bucket: aws.String(cli.bucketName),
		Key:    aws.String(key),
	}

	result, err := cli.sdkClient.GetObject(input)
	if err != nil {
		return nil, nil, NewError(fmt.Errorf("error getting object from s3: %w", err), log.Data{
			"bucket_name": cli.bucketName,
			"s3_key":      key, // key is the s3 filename with path (it's not a cryptographic key)
			"user_psk":    false,
		})
	}

	return result.Body, result.ContentLength, nil
}

// GetWithPSK returns an io.ReadCloser instance for the given path (inside the bucket configured for this client)
// and the content length (size in bytes). It uses the provided PSK for encryption.
// The 'key' parameter refers to the path for the file under the bucket.
//
// The caller is responsible for closing the returned ReadCloser.
// For example, it may be closed in a defer statement: defer r.Close()
func (cli *Client) GetWithPSK(key string, psk []byte) (io.ReadCloser, *int64, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(cli.bucketName),
		Key:    aws.String(key),
	}

	result, err := cli.cryptoClient.GetObjectWithPSK(input, psk)
	if err != nil {
		return nil, nil, NewError(fmt.Errorf("error getting object from s3: %w", err), log.Data{
			"bucket_name": cli.bucketName,
			"s3_key":      key, // key is the s3 filename with path (it's not a cryptographic key)
			"user_psk":    true,
		})
	}

	return result.Body, result.ContentLength, nil
}

// Head returns a HeadObjectOutput containing an object metadata obtained from ah HTTP HEAD call
func (cli *Client) Head(key string) (*s3.HeadObjectOutput, error) {
	result, err := cli.sdkClient.HeadObject(&s3.HeadObjectInput{
		Bucket: &cli.bucketName,
		Key:    &key,
	})
	if err != nil {
		return nil, NewError(
			fmt.Errorf("error trying to obtain s3 object metadata with HeadObject call: %w", err),
			log.Data{
				"bucket_name": cli.bucketName,
				"s3_key":      key, // key is the s3 filename with path (it's not a cryptographic key)
			},
		)
	}
	return result, nil
}

func (cli *Client) FileExists(key string) (bool, error) {
	_, err := cli.sdkClient.HeadObject(&s3.HeadObjectInput{
		Bucket: &cli.bucketName,
		Key:    &key,
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound": // s3.ErrCodeNoSuchKey does not work, aws is missing this error code so we hardwire a string
				return false, nil
			default:
				return false, aerr
			}
		}

		return false, err
	}

	return true, nil
}

func (cli *Client) GetBucketPolicy(BucketName string) (*s3.GetBucketPolicyOutput, error) {
	result, err := cli.sdkClient.GetBucketPolicy(&s3.GetBucketPolicyInput{
		Bucket: aws.String(BucketName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound": // s3.ErrCodeNoSuchKey does not work, aws is missing this error code so we hardwire a string
				return nil, nil
			default:
				return nil, aerr
			}
		}

		return nil, err
	}
	return result, nil
}

func (cli *Client) PutBucketPolicy(BucketName string, policy string) (*s3.PutBucketPolicyOutput, error) {

	result, err := cli.sdkClient.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(BucketName),
		Policy: aws.String(string(policy)),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound": // s3.ErrCodeNoSuchKey does not work, aws is missing this error code so we hardwire a string
				return nil, nil
			default:
				return nil, aerr
			}
		}

		return nil, err
	}
	return result, nil
}

func (cli *Client) ListObjects(BucketName string) (*s3.ListObjectsOutput, error) {

	result, err := cli.sdkClient.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(BucketName),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound": // s3.ErrCodeNoSuchKey does not work, aws is missing this error code so we hardwire a string
				return nil, nil
			default:
				return nil, aerr
			}
		}

		return nil, err
	}
	return result, nil
}
