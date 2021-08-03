package s3client

import (
	"context"
	"fmt"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

// ServiceName S3
const ServiceName = "S3"

const codeNotFound = "NotFound"

// MsgHealthy is the message in the Check structure when S3 is healthy
const MsgHealthy = "S3 is healthy"

// Checker validates that the S3 bucket exists, and updates the provided CheckState accordingly.
// Any error during the state update will be returned
func (cli *S3) Checker(ctx context.Context, state *health.CheckState) error {
	if err := cli.ValidateBucket(); err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// AWS error
			return cli.handleAWSErr(awsErr, state)
		}
		// Generic error
		return state.Update(health.StatusCritical, err.Error(), 0)
	}
	// Success
	return state.Update(health.StatusOK, MsgHealthy, 0)
}

// handleAWSErr updates the provided CheckState with a Critical state and a message according to the provided AWS error.
// For inexistent buckets, a relevant error message will be generated, for any other error we use the AWS Code (consice string).
// Any error during the state update will be returned
func (cli *S3) handleAWSErr(err awserr.Error, state *health.CheckState) error {
	code := err.Code()
	switch code {
	case codeNotFound:
		// Bucket not found
		return state.Update(health.StatusCritical, fmt.Sprintf("Bucket not found: %s", cli.bucketName), 0)
	default:
		// Other AWS error
		return state.Update(health.StatusCritical, err.Code(), 0)
	}
}
