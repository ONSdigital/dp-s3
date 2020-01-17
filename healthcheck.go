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
	currentTime := time.Now().UTC()
	s3.Check.LastChecked = &currentTime
	if err != nil {
		s3.Check.LastFailure = &currentTime
		s3.Check.Status = health.StatusCritical
		s3.Check.Message = err.Error()
		return s3.Check, err
	}
	reader.Close()
	s3.Check.LastSuccess = &currentTime
	s3.Check.Status = health.StatusOK
	s3.Check.Message = MsgHealthy
	return s3.Check, nil
}
