// file: healthcheck.go
//
// Contains methods to get the health state of an S3 client from S3,
// by checking that the bucket exists in the provided region.
//
// Requires "s3:ListBucket" action allowed by IAM policy for the bucket,
// as defined by `check-{bucketName}-bucket` policies in dp-setup
package s3

import (
	"context"
	"errors"
	"fmt"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

// ServiceName S3
const ServiceName = "S3"

// MsgHealthy is the message in the Check structure when S3 is healthy
const MsgHealthy = "S3 is healthy"

// Checker validates that the S3 bucket exists, and updates the provided CheckState accordingly.
// Any error during the state update will be returned
func (cli *Client) Checker(ctx context.Context, state *health.CheckState) error {
	if err := cli.ValidateBucket(ctx); err != nil {
		return cli.handleAWSErr(err, state)
	}
	// Success
	return state.Update(health.StatusOK, MsgHealthy, 0)
}

// handleAWSErr updates the provided CheckState with a Critical state and a message according to the provided AWS error.
// For inexistent buckets, a relevant error message will be generated, for any other error we use the AWS Code (consice string).
// Any error during the state update will be returned
func (cli *Client) handleAWSErr(err error, state *health.CheckState) error {
	var bucketNotFoundErr *types.NoSuchBucket
	if errors.As(err, &bucketNotFoundErr) {
		// Bucket not found
		return state.Update(health.StatusCritical, fmt.Sprintf("Bucket not found: %s", cli.bucketName), 0)
	}

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		// Other AWS service error
		return state.Update(health.StatusCritical, apiErr.ErrorMessage(), 0)
	}

	// Generic error
	return state.Update(health.StatusCritical, err.Error(), 0)
}

// ValidateBucket checks that the bucket exists and returns an error if it
// does not exist or there was some other error trying to get this information.
func (cli *Client) ValidateBucket(ctx context.Context) error {
	_, err := cli.sdkClient.HeadBucket(
		ctx,
		&s3.HeadBucketInput{
			Bucket: aws.String(cli.bucketName),
		},
	)
	return err
}
