package s3client

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/ONSdigital/log.go/log"

	"github.com/ONSdigital/s3crypto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3 client with SDK client, CryptoClient and BucketName
type S3 struct {
	sdkClient    S3SDKClient
	cryptoClient S3CryptoClient
	bucketName   string
	region       string
}

// Resumable represents a resumable upload
type Resumable struct {
	UploadKey   string
	Type        string
	ChunkNumber int
	TotalChunks int
	FileName    string
}

// New creates a new S3 Client configured for the given region and bucket name.
// If HasUserDefinedPSK is true, it will also have a crypto client.
func New(region string, bucketName string, HasUserDefinedPSK bool) (*S3, error) {

	// Create AWS session with the provided region
	sess, err := session.NewSession(&aws.Config{Region: &region})
	if err != nil {
		return nil, err
	}

	// Create AWS-SDK-S3 client with the session
	sdkClient := s3.New(sess)

	// If we require crypto client (HasUserDefinedPSK), create it.
	var cryptoClient S3CryptoClient
	if HasUserDefinedPSK {
		cryptoConfig := &s3crypto.Config{HasUserDefinedPSK: true}
		s3cryptoClient := s3crypto.New(sess, cryptoConfig)

		cryptoClient = s3cryptoClient
	}

	cli := Instantiate(sdkClient, cryptoClient, bucketName, region)
	return cli, nil
}

// Instantiate creates a new instance of S3 struct with the provided clients, bucket and region
func Instantiate(sdkClient S3SDKClient, cryptoClient S3CryptoClient, bucketName, region string) *S3 {
	return &S3{
		sdkClient:    sdkClient,
		cryptoClient: cryptoClient,
		bucketName:   bucketName,
		region:       region,
	}
}

// Upload handles the uploading a file to AWS S3, into the bucket configured for this client
func (cli *S3) Upload(ctx context.Context, resum Resumable, b []byte) error {
	return cli.UploadWithPsk(ctx, resum, b, nil)
}

// UploadWithPsk handles the uploading a file to AWS S3, into the bucket configured for this client, using a user-defined psk
func (cli *S3) UploadWithPsk(ctx context.Context, resum Resumable, b []byte, psk []byte) error {

	listMultiInput := &s3.ListMultipartUploadsInput{
		Bucket: &cli.bucketName,
	}

	listMultiOutput, err := cli.sdkClient.ListMultipartUploads(listMultiInput)
	if err != nil {
		log.Event(ctx, "error fetching multipart list", log.Error(err))
		return err
	}

	var uploadID string
	for _, upload := range listMultiOutput.Uploads {
		if *upload.Key == resum.UploadKey {
			uploadID = *upload.UploadId
		}
	}

	if len(uploadID) == 0 {

		createMultiInput := &s3.CreateMultipartUploadInput{
			Bucket:      &cli.bucketName,
			Key:         &resum.UploadKey,
			ContentType: &resum.Type,
		}

		createMultiOutput, err := cli.sdkClient.CreateMultipartUpload(createMultiInput)
		if err != nil {
			log.Event(ctx, "error creating multipart upload", log.Error(err))
			return err
		}

		uploadID = *createMultiOutput.UploadId
	}

	rs := bytes.NewReader(b)

	n := int64(resum.ChunkNumber)

	uploadPartInput := &s3.UploadPartInput{
		UploadId:   &uploadID,
		Bucket:     &cli.bucketName,
		Key:        &resum.UploadKey,
		Body:       rs,
		PartNumber: &n,
	}

	if psk != nil {
		_, err = cli.cryptoClient.UploadPartWithPSK(uploadPartInput, psk)
		if err != nil {
			log.Event(ctx, "error uploading with psk", log.Error(err))
			return err
		}
	} else {
		_, err = cli.sdkClient.UploadPart(uploadPartInput)
		if err != nil {
			log.Event(ctx, "error uploading part", log.Error(err))
			return err
		}
	}

	log.Event(ctx, "chunk accepted", log.Data{
		"chunk_number": resum.ChunkNumber,
		"max_chunks":   resum.TotalChunks,
		"file_name":    resum.FileName,
	})

	input := &s3.ListPartsInput{
		Key:      &resum.UploadKey,
		Bucket:   &cli.bucketName,
		UploadId: &uploadID,
	}

	output, err := cli.sdkClient.ListParts(input)
	if err != nil {
		log.Event(ctx, "error listing parts", log.Error(err))
		return err
	}

	parts := output.Parts
	if len(parts) == resum.TotalChunks {
		return cli.completeUpload(ctx, uploadID, resum.UploadKey, resum.TotalChunks, parts)
	}

	return nil
}

// CheckUploaded check uploaded. Returns true only if the chunk corresponding to the provided chunkNumber has been uploaded.
// If the upload is finished, we complete it.
func (cli *S3) CheckUploaded(ctx context.Context, uploadKey string, chunkNumber, totalChunks int) (bool, error) {

	listMultiInput := &s3.ListMultipartUploadsInput{
		Bucket: &cli.bucketName,
	}

	listMultiOutput, err := cli.sdkClient.ListMultipartUploads(listMultiInput)
	if err != nil {
		log.Event(ctx, "error fetching multipart upload list", log.Error(err))
		return false, err
	}

	var uploadID string
	for _, upload := range listMultiOutput.Uploads {
		if *upload.Key == uploadKey {
			uploadID = *upload.UploadId
		}
	}
	if len(uploadID) == 0 {
		return false, ErrChunkNumberNotUploaded
	}

	input := &s3.ListPartsInput{
		Key:      &uploadKey,
		Bucket:   &cli.bucketName,
		UploadId: &uploadID,
	}

	output, err := cli.sdkClient.ListParts(input)
	if err != nil {
		log.Event(ctx, "error listing parts", log.Error(err))
		return false, ErrChunkNumberNotFound
	}

	parts := output.Parts
	if len(parts) == totalChunks {
		if err = cli.completeUpload(ctx, uploadID, uploadKey, totalChunks, parts); err != nil {
			return false, err
		}
		return true, nil
	}

	for _, part := range parts {
		if *part.PartNumber == int64(chunkNumber) {
			return true, nil
		}
	}

	return false, ErrChunkNumberNotFound
}

// completeUpload if all parts have been uploaded, we complete the multipart upload.
func (cli *S3) completeUpload(ctx context.Context, uploadID, uploadKey string, totalChunks int, parts []*s3.Part) error {
	var completedParts []*s3.CompletedPart

	for _, part := range parts {
		completedParts = append(completedParts, &s3.CompletedPart{
			PartNumber: part.PartNumber,
			ETag:       part.ETag,
		})
	}

	if len(completedParts) == totalChunks {
		completeInput := &s3.CompleteMultipartUploadInput{
			Key:      &uploadKey,
			UploadId: &uploadID,
			MultipartUpload: &s3.CompletedMultipartUpload{
				Parts: completedParts,
			},
			Bucket: &cli.bucketName,
		}

		log.Event(ctx, "going to attempt to complete", log.Data{"complete": completeInput})

		_, err := cli.sdkClient.CompleteMultipartUpload(completeInput)
		if err != nil {
			log.Event(ctx, "error completing upload", log.Error(err))
			return err
		}
	}

	return nil
}

// GetURL returns a full S3 URL from the provided path and the bucket and region configured for the client.
func (cli *S3) GetURL(path string) string {
	url := "https://s3-%s.amazonaws.com/%s/%s"
	return fmt.Sprintf(url, cli.region, cli.bucketName, path)
}

// GetFromURL returns an io.ReadCloser instance for the given fully qualified S3 URL.
// If the URL defines a bucket different from the one configured in this client, an error will be returned.
func (cli *S3) GetFromURL(rawURL string) (io.ReadCloser, error) {

	// Use the S3 URL implementation as the S3 drivers don't seem to handle fully qualified URLs that include the
	// bucket name.
	url, err := NewURL(rawURL)
	if err != nil {
		return nil, err
	}

	// Validate that bucket defined by URL matches the bucket of this client
	if url.BucketName() != cli.bucketName {
		return nil, ErrUnexpectedBucket
	}

	return cli.Get(url.Path())
}

// Get returns an io.ReadCloser instance for the given path (inside the bucket configured for this client).
// They 'key' parameter refers to the path for the file under the bucket.
func (cli *S3) Get(key string) (io.ReadCloser, error) {

	input := &s3.GetObjectInput{
		Bucket: aws.String(cli.bucketName),
		Key:    aws.String(key),
	}

	result, err := cli.sdkClient.GetObject(input)
	if err != nil {
		return nil, err
	}

	return result.Body, nil
}

// ValidateBucket checks that the bucket exists and returns an error if it
// does not exist or there was some other error trying to get this information.
func (cli *S3) ValidateBucket() error {

	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(cli.bucketName),
		MaxKeys: aws.Int64(1),
	}

	_, err := cli.sdkClient.ListObjectsV2(input)
	return err
}
