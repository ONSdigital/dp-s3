package s3client

import (
	"context"
	"fmt"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ServiceName S3
const ServiceName = "S3"

// MsgHealthy is the message in the Check structure when S3 is healthy
const MsgHealthy = "S3 is healthy"

// ErrBucketDoesNotExist is an Error to handle failures getting the S3 bucket
type ErrBucketDoesNotExist struct {
	BucketName string
}

// Error returns the error message with the bucket name
func (e *ErrBucketDoesNotExist) Error() string {
	return fmt.Sprintf("Bucket %s does not exist", e.BucketName)
}

// Checker validates that the S3 bucket exists, and updates the provided CheckState accordingly
func (cli *S3) Checker(ctx context.Context, state *health.CheckState) error {
	err := cli.ValidateBucket()
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case s3.ErrCodeNoSuchBucket:
				// Bucket does not exist
				errBucket := ErrBucketDoesNotExist{BucketName: cli.bucketName}
				state.Update(health.StatusCritical, errBucket.Error(), 0)
				return nil
			default:
				// Other AWS error
				state.Update(health.StatusCritical, awsErr.Code(), 0)
				return nil
			}
		}
		// Generic error
		state.Update(health.StatusCritical, err.Error(), 0)
		return err
	}
	// Success
	state.Update(health.StatusOK, MsgHealthy, 0)
	return nil
}
