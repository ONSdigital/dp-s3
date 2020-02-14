package s3client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/ONSdigital/log.go/log"

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
// If HasUserDefinedPSK is true, it will also have a crypto client.
// Note: This function will create a new session, if you already have a session, please use NewUploaderWithSession instead
func NewClient(region string, bucketName string, hasUserDefinedPSK bool) (*S3, error) {
	s, err := session.NewSession(&aws.Config{Region: &region})
	if err != nil {
		return nil, err
	}
	return NewClientWithSession(region, bucketName, hasUserDefinedPSK, s), nil
}

// NewClientWithSession creates a new S3 Client configured for the given region and bucket name, using the provided session.
// If HasUserDefinedPSK is true, it will also have a crypto client.
func NewClientWithSession(region string, bucketName string, hasUserDefinedPSK bool, s *session.Session) *S3 {

	// Create AWS-SDK-S3 client with the session
	var sdkClient S3SDKClient = s3.New(s)

	// If we require crypto client (HasUserDefinedPSK), create it.
	var cryptoClient S3CryptoClient
	if hasUserDefinedPSK {
		cryptoClient = s3crypto.New(s, &s3crypto.Config{HasUserDefinedPSK: true})
	}

	return InstantiateClient(sdkClient, cryptoClient, bucketName, region, s)
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

	// Get UploadID or create it if it does not exist (atomically)
	uploadID, err := cli.doGetOrCreateMultipartUpload(ctx, req)
	if err != nil {
		return err
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
		return err
	}

	log.Event(ctx, "chunk accepted", log.Data{
		"chunk_number": req.ChunkNumber,
		"max_chunks":   req.TotalChunks,
		"file_name":    req.FileName,
	})

	// List parts so that we can validate if the upload operation is complete
	output, err := cli.sdkClient.ListParts(
		&s3.ListPartsInput{
			Key:      &req.UploadKey,
			Bucket:   &cli.bucketName,
			UploadId: &uploadID,
		},
	)
	if err != nil {
		log.Event(ctx, "error listing parts", log.Error(err))
		return err
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
		log.Event(ctx, "error fetching multipart list", log.Error(err))
		return "", err
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
		log.Event(ctx, "error creating multipart upload", log.Error(err))
		return "", err
	}
	return *createMultiOutput.UploadId, nil
}

// doUploadPart performs the upload using the sdkClient if no psk is provided, or the cryptoClient if psk is provided
func (cli *S3) doUploadPart(ctx context.Context, input *s3.UploadPartInput, psk []byte) (out *s3.UploadPartOutput, err error) {

	if psk != nil {
		// Upload Part with PSK
		out, err = cli.cryptoClient.UploadPartWithPSK(input, psk)
		if err != nil {
			log.Event(ctx, "error uploading with psk", log.Error(err))
		}
		return
	}

	// Upload part without user-defined PSK
	out, err = cli.sdkClient.UploadPart(input)
	if err != nil {
		log.Event(ctx, "error uploading part", log.Error(err))
	}
	return
}

// CheckPartUploaded returns true only if the chunk corresponding to the provided chunkNumber has been uploaded.
// If all the chunks have been uploaded, we complete the upload operation.
func (cli *S3) CheckPartUploaded(ctx context.Context, req *UploadPartRequest) (bool, error) {

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
		if *upload.Key == req.UploadKey {
			uploadID = *upload.UploadId
			break
		}
	}
	if len(uploadID) == 0 {
		log.Event(ctx, "chunk number not uploaded", log.Data{"chunk_number": req.ChunkNumber, "file_name": req.FileName})
		return false, &ErrNotUploaded{UploadKey: req.UploadKey}
	}

	input := &s3.ListPartsInput{
		Key:      &req.UploadKey,
		Bucket:   &cli.bucketName,
		UploadId: &uploadID,
	}

	output, err := cli.sdkClient.ListParts(input)
	if err != nil {
		log.Event(ctx, "chunk number verification error", log.Error(err), log.Data{"chunk_number": req.ChunkNumber, "file_name": req.FileName})
		return false, &ErrListParts{err.Error()}
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
			log.Event(ctx, "chunk number already uploaded", log.Data{"chunk_number": req.ChunkNumber, "file_name": req.FileName, "identifier": req.UploadKey})
			return true, nil
		}
	}

	log.Event(ctx, "chunk number failed to upload", log.Data{"chunk_number": req.ChunkNumber, "file_name": req.FileName})
	return false, &ErrChunkNumberNotFound{req.ChunkNumber}
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

		log.Event(ctx, "going to attempt to complete", log.Data{"complete": completeInput})

		_, err := cli.sdkClient.CompleteMultipartUpload(completeInput)
		if err != nil {
			log.Event(ctx, "error completing upload", log.Error(err))
			return err
		}
	}
	return nil
}

// GetPathStyleURL returns an https S3 URL from the provided path and the bucket and region
// configured for the client in path-style. Note: this format is deprecated by Amazon
// https://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html
func (cli *S3) GetPathStyleURL(path string) string {
	url := "https://s3-%s.amazonaws.com/%s/%s"
	return fmt.Sprintf(url, cli.region, cli.bucketName, path)
}

// GetFromS3URL returns an io.ReadCloser instance for the given S3 virtual-hosted-style URL (s3://bucketName/objectKey)
// If the URL defines a bucket different from the one configured in this client, an error will be returned.
func (cli *S3) GetFromS3URL(rawURL string) (io.ReadCloser, error) {

	// Use the S3 URL implementation as the S3 drivers don't seem to handle fully qualified URLs that include the
	// bucket name.
	url, err := NewURL(rawURL)
	if err != nil {
		return nil, err
	}

	// Validate that bucket defined by URL matches the bucket of this client
	if url.BucketName() != cli.bucketName {
		return nil, &ErrUnexpectedBucket{
			ExpectedBucketName: cli.bucketName, BucketName: url.BucketName()}
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
