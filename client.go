package s3client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/ONSdigital/log.go/v2/log"

	"github.com/ONSdigital/s3crypto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3 client with SDK client, CryptoClient and BucketName
type S3 struct {
	sdkClient     S3SDKClient
	cryptoClient  S3CryptoClient
	bucketName    string
	region        string
	mutexUploadID *sync.Mutex
	session       *session.Session
}

// UploadPartRequest represents a part upload request
type UploadPartRequest struct {
	UploadKey   string
	Type        string
	ChunkNumber int64
	TotalChunks int
	FileName    string
}

// NewClient creates a new S3 Client configured for the given region and bucket name.
// Note: This function will create a new session, if you already have a session, please use NewUploaderWithSession instead
// Any error establishing the AWS session will be returned
func NewClient(region string, bucketName string) (*S3, error) {
	s, err := session.NewSession(&aws.Config{Region: &region})
	if err != nil {
		return nil, NewError(
			fmt.Errorf("error creating session: %w", err),
			log.Data{
				"region":      region,
				"bucket_name": bucketName,
			},
		)
	}
	return NewClientWithSession(bucketName, s), nil
}

// NewClientWithSession creates a new S3 Client configured for the given bucket name, using the provided session and region within it.
func NewClientWithSession(bucketName string, s *session.Session) *S3 {
	// Get region for the Session config
	region := s.Config.Region

	// Create AWS-SDK-S3 client with the session
	sdkClient := s3.New(s)

	// Create crypto client, which allows user to provide a psk
	cryptoClient := s3crypto.New(s, &s3crypto.Config{HasUserDefinedPSK: true})

	return InstantiateClient(sdkClient, cryptoClient, bucketName, *region, s)
}

// InstantiateClient creates a new instance of S3 struct with the provided clients, bucket and region
func InstantiateClient(sdkClient S3SDKClient, cryptoClient S3CryptoClient, bucketName, region string, s *session.Session) *S3 {
	return &S3{
		sdkClient:     sdkClient,
		cryptoClient:  cryptoClient,
		bucketName:    bucketName,
		region:        region,
		mutexUploadID: &sync.Mutex{},
		session:       s,
	}
}

// Session returns the Session of this client
func (cli *S3) Session() *session.Session {
	return cli.session
}

// BucketName returns the bucket name used by this S3 client
func (cli *S3) BucketName() string {
	return cli.bucketName
}

// UploadPart handles the uploading a file to AWS S3, into the bucket configured for this client
func (cli *S3) UploadPart(ctx context.Context, req *UploadPartRequest, payload []byte) error {
	return cli.UploadPartWithPsk(ctx, req, payload, nil)
}

// UploadPartWithPsk handles the uploading a file to AWS S3, into the bucket configured for this client, using a user-defined psk
func (cli *S3) UploadPartWithPsk(ctx context.Context, req *UploadPartRequest, payload []byte, psk []byte) error {
	logData := log.Data{
		"chunk_number": req.ChunkNumber,
		"max_chunks":   req.TotalChunks,
		"file_name":    req.FileName,
		"bucket_name":  cli.bucketName,
		"user_psk":     psk != nil,
	}

	// Get UploadID or create it if it does not exist (atomically)
	uploadID, err := cli.doGetOrCreateMultipartUpload(ctx, req)
	if err != nil {
		return NewError(err, logData)
	}

	// Do the upload against AWS
	_, err = cli.doUploadPart(
		ctx,
		&s3.UploadPartInput{
			UploadId:   &uploadID,
			Bucket:     &cli.bucketName,
			Key:        &req.UploadKey,
			Body:       bytes.NewReader(payload),
			PartNumber: &req.ChunkNumber,
		},
		psk,
	)
	if err != nil {
		return NewError(err, logData)
	}

	log.Info(ctx, "chunk accepted", logData)

	// List parts so that we can validate if the upload operation is complete
	output, err := cli.sdkClient.ListParts(
		&s3.ListPartsInput{
			Key:      &req.UploadKey,
			Bucket:   &cli.bucketName,
			UploadId: &uploadID,
		},
	)
	if err != nil {
		return NewError(fmt.Errorf("error listing parts: %w", err), logData)
	}

	// If all parts have been uploaded, we call completeUpload
	parts := output.Parts
	if len(parts) == req.TotalChunks {
		return cli.completeUpload(ctx, uploadID, req, parts)
	}

	// Otherwise we don't need to perform any other operation.
	return nil
}

// doGetOrCreateMultipartUpload atomically gets the UploadId for the specified bucket
// and S3 object key, and if it does not find it, it creates it.
func (cli *S3) doGetOrCreateMultipartUpload(ctx context.Context, req *UploadPartRequest) (uploadID string, err error) {

	cli.mutexUploadID.Lock()
	defer cli.mutexUploadID.Unlock()

	// List existing Multipart uploads for our bucket
	listMultiOutput, err := cli.sdkClient.ListMultipartUploads(
		&s3.ListMultipartUploadsInput{
			Bucket: &cli.bucketName,
		})
	if err != nil {
		return "", fmt.Errorf("error fetching multipart list: %w", err)
	}

	// Try to find a multipart upload for the same s3 object that we want
	for _, upload := range listMultiOutput.Uploads {
		if *upload.Key == req.UploadKey {
			uploadID = *upload.UploadId
			return uploadID, nil
		}
	}

	// If we didn't find the Multipart upload, create it
	createMultiOutput, err := cli.sdkClient.CreateMultipartUpload(
		&s3.CreateMultipartUploadInput{
			Bucket:      &cli.bucketName,
			Key:         &req.UploadKey,
			ContentType: &req.Type,
		})
	if err != nil {
		return "", fmt.Errorf("error creating multipart upload: %w", err)
	}
	return *createMultiOutput.UploadId, nil
}

// doUploadPart performs the upload using the sdkClient if no psk is provided, or the cryptoClient if psk is provided
func (cli *S3) doUploadPart(ctx context.Context, input *s3.UploadPartInput, psk []byte) (out *s3.UploadPartOutput, err error) {
	if psk != nil {
		// Upload Part with PSK
		out, err = cli.cryptoClient.UploadPartWithPSK(input, psk)
		if err != nil {
			return nil, fmt.Errorf("error uploading part with psk: %w", err)
		}
		return
	}

	// Upload part without user-defined PSK
	out, err = cli.sdkClient.UploadPart(input)
	if err != nil {
		return nil, fmt.Errorf("error uploading part with psk: %w", err)
	}
	return
}

// CheckPartUploaded returns true only if the chunk corresponding to the provided chunkNumber has been uploaded.
// If all the chunks have been uploaded, we complete the upload operation.
func (cli *S3) CheckPartUploaded(ctx context.Context, req *UploadPartRequest) (bool, error) {
	logData := log.Data{
		"chunk_number": req.ChunkNumber,
		"max_chunks":   req.TotalChunks,
		"file_name":    req.FileName,
		"bucket_name":  cli.bucketName,
		"identifier":   req.UploadKey,
	}

	listMultiInput := &s3.ListMultipartUploadsInput{
		Bucket: &cli.bucketName,
	}

	listMultiOutput, err := cli.sdkClient.ListMultipartUploads(listMultiInput)
	if err != nil {
		return false, NewError(fmt.Errorf("error fetching multipart upload list: %w", err), logData)
	}

	var uploadID string
	for _, upload := range listMultiOutput.Uploads {
		if *upload.Key == req.UploadKey {
			uploadID = *upload.UploadId
			break
		}
	}
	if len(uploadID) == 0 {
		return false, NewError(errors.New("s3 key not uploaded"), logData)
	}

	input := &s3.ListPartsInput{
		Key:      &req.UploadKey,
		Bucket:   &cli.bucketName,
		UploadId: &uploadID,
	}

	output, err := cli.sdkClient.ListParts(input)
	if err != nil {
		return false, NewError(fmt.Errorf("chunk number verification error: %w", err), logData)
	}

	// TODO: If there are more than 1000 parts, they will be paginated, so we would need to call ListParts again with the provided Marker until we have all of them.
	// Reference: https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#S3.ListParts

	parts := output.Parts
	if len(parts) == req.TotalChunks {
		if err = cli.completeUpload(ctx, uploadID, req, parts); err != nil {
			return false, err
		}
		return true, nil
	}

	for _, part := range parts {
		if *part.PartNumber == int64(req.ChunkNumber) {
			log.Info(ctx, "chunk already uploaded", logData)
			return true, nil
		}
	}

	return false, NewError(errors.New("chunk number not found"), logData)
}

// completeUpload if all parts have been uploaded, we complete the multipart upload.
func (cli *S3) completeUpload(ctx context.Context, uploadID string, req *UploadPartRequest, parts []*s3.Part) error {
	var completedParts []*s3.CompletedPart

	for _, part := range parts {
		completedParts = append(completedParts, &s3.CompletedPart{
			PartNumber: part.PartNumber,
			ETag:       part.ETag,
		})
	}

	if len(completedParts) == req.TotalChunks {
		completeInput := &s3.CompleteMultipartUploadInput{
			Key:      &req.UploadKey,
			UploadId: &uploadID,
			MultipartUpload: &s3.CompletedMultipartUpload{
				Parts: completedParts,
			},
			Bucket: &cli.bucketName,
		}

		log.Info(ctx, "attemtping to complete multipart upload", log.Data{"complete": completeInput})

		_, err := cli.sdkClient.CompleteMultipartUpload(completeInput)
		if err != nil {
			return NewError(fmt.Errorf("error completing multipart upload: %w", err), log.Data{"complete": completeInput})
		}
	}
	return nil
}

// GetFromS3URL returns an io.ReadCloser instance and the content length (size in bytes) for the given S3 URL,
// in the format specified by URLStyle.
// More information about s3 URL styles: https://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html
// If the URL defines a region (if provided) or bucket different from the one configured in this client, an error will be returned.
func (cli *S3) GetFromS3URL(rawURL string, style URLStyle) (io.ReadCloser, *int64, error) {
	return cli.doGetFromS3URL(rawURL, style, nil)

}

// GetFromS3URLWithPSK returns an io.ReadCloser instance and the content length (size in bytes) for the given S3 URL,
// in the format specified by URLStyle, using the provided PSK for encryption.
// More information about s3 URL styles: https://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html
// If the URL defines a region (if provided) or bucket different from the one configured in this client, an error will be returned.
func (cli *S3) GetFromS3URLWithPSK(rawURL string, style URLStyle, psk []byte) (io.ReadCloser, *int64, error) {
	return cli.doGetFromS3URL(rawURL, style, psk)
}

func (cli *S3) doGetFromS3URL(rawURL string, style URLStyle, psk []byte) (io.ReadCloser, *int64, error) {
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
		return nil, nil, NewError(errors.New("unexpected bucket name in url"), logData)
	}

	// Validate that URL and client regions match, if URL provides one
	if len(s3Url.Region) > 0 && s3Url.Region != cli.region {
		logData["region"] = cli.region
		return nil, nil, NewError(errors.New("unexpected aws region in url"), logData)
	}

	if psk == nil {
		return cli.Get(s3Url.Key)
	}
	return cli.GetWithPSK(s3Url.Key, psk)
}

// Get returns an io.ReadCloser instance for the given path (inside the bucket configured for this client)
// and the content length (size in bytes).
// They 'key' parameter refers to the path for the file under the bucket.
func (cli *S3) Get(key string) (io.ReadCloser, *int64, error) {

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
func (cli *S3) GetWithPSK(key string, psk []byte) (io.ReadCloser, *int64, error) {

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

// PutWithPSK uploads the provided contents to the key in the bucket configured for this client, using the provided PSK.
// The 'key' parameter refers to the path for the file under the bucket.
func (cli *S3) PutWithPSK(key *string, reader *bytes.Reader, psk []byte) error {

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

// ValidateBucket checks that the bucket exists and returns an error if it
// does not exist or there was some other error trying to get this information.
func (cli *S3) ValidateBucket() error {
	_, err := cli.sdkClient.HeadBucket(
		&s3.HeadBucketInput{
			Bucket: aws.String(cli.bucketName),
		},
	)
	return err
}
