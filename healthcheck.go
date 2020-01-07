package s3client

import (
	"context"
	"errors"
	"time"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
)

// ServiceName S3
const ServiceName = "S3"

// MsgHealthy is the message in the Check structure when S3 is healthy
const MsgHealthy = "S3 is healthy"

// Error definitions
var (
	ErrBucketDoesNotExist = errors.New("The specified bucket does not exist")
)

// minTime is the oldest time for Check structure.
var minTime = time.Unix(0, 0)

// Checker performs a check health of S3 and returns it inside a Check structure
func (s3 *S3) Checker(ctx *context.Context, bucketName string) (*health.Check, error) {
	reader, err := s3.Get("s3://" + bucketName)
	if err != nil {
		return getCheck(ctx, health.StatusCritical, err.Error()), err
	}
	reader.Close()
	return getCheck(ctx, health.StatusOK, MsgHealthy), nil
}

// getCheck creates a Check structure and populates it according to status and message provided
func getCheck(ctx *context.Context, status, message string) *health.Check {

	currentTime := time.Now().UTC()

	check := &health.Check{
		Name:        ServiceName,
		Status:      status,
		Message:     message,
		LastChecked: currentTime,
		LastSuccess: minTime,
		LastFailure: minTime,
	}

	switch status {
	case health.StatusOK:
		check.StatusCode = 200
		check.LastSuccess = currentTime
	case health.StatusWarning:
		check.StatusCode = 429
		check.LastFailure = currentTime
	default:
		check.StatusCode = 500
		check.Status = health.StatusCritical
		check.LastFailure = currentTime
	}

	return check
}
