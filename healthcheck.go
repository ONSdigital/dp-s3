package s3client

import (
	"context"
	"fmt"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
)

//go:generate moq -out ./mock/check_state.go -pkg mock . CheckState

// CheckState interface corresponds to the healthcheck CheckState structure
type CheckState interface {
	Update(status, message string, statusCode int) error
}

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

// Checker performs a check health of S3 and returns it inside a Check structure
func (s3 *S3) Checker(ctx context.Context, state CheckState) error {
	reader, err := s3.Get("s3://" + s3.BucketName)
	if err != nil {
		// errBucket := ErrBucketDoesNotExist{BucketName: s3.BucketName}
		state.Update(health.StatusCritical, err.Error(), 0)
		return nil
	}
	reader.Close()
	state.Update(health.StatusOK, MsgHealthy, 0)
	return nil
}
