package s3client

import (
	"context"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

// ServiceName S3
const ServiceName = "S3"

const codeNotFound = "NotFound"

// MsgHealthy is the message in the Check structure when S3 is healthy
const MsgHealthy = "S3 is healthy"

// Checker validates that the S3 bucket exists, and updates the provided CheckState accordingly
func (cli *S3) Checker(ctx context.Context, state *health.CheckState) error {
	err := cli.ValidateBucket()
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// AWS error
			cli.handleAWSErr(awsErr, state)
			return nil
		}
		// Generic error
		state.Update(health.StatusCritical, err.Error(), 0)
		return err
	}
	// Success
	state.Update(health.StatusOK, MsgHealthy, 0)
	return nil
}

// handleAWSErr updates the provided CheckState with a Critical state and a message according to the provided AWS error.
// For inexistent buckets, a relevant error message will be generated, for any other error we use the AWS Code (consice string).
func (cli *S3) handleAWSErr(err awserr.Error, state *health.CheckState) {
	code := err.Code()
	switch code {
	case codeNotFound:
		// Bucket not found
		errBucket := ErrBucketNotFound{BucketName: cli.bucketName}
		state.Update(health.StatusCritical, errBucket.Error(), 0)
	default:
		// Other AWS error
		state.Update(health.StatusCritical, err.Code(), 0)
	}
}
