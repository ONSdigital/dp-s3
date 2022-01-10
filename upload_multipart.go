// file: upload_multipart.go
//
// Contains methods to upload files to S3 in chunks
// by using the low level SDK methods that give the caller control over
// the multipart uploading process.
//
// Requires "s3:PutObject", "s3:GetObject" and "s3:AbortMultipartUpload" actions allowed by IAM policy for the bucket,
// as defined by `multipart-{bucketName}-bucket` policies in dp-setup
package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/aws/aws-sdk-go/service/s3"
)

// UploadPartRequest represents a part upload request
type UploadPartRequest struct {
	UploadKey   string
	Type        string
	ChunkNumber int64
	TotalChunks int
	FileName    string
}

type MultipartUploadResponse struct {
	Etag string
	AllPartsUploaded bool
}

// UploadPart handles the uploading a file to AWS S3, into the bucket configured for this client
func (cli *Client) UploadPart(ctx context.Context, req *UploadPartRequest, payload []byte) (MultipartUploadResponse, error) {
	return cli.UploadPartWithPsk(ctx, req, payload, nil)
}

// UploadPartWithPsk handles the uploading a file to AWS S3, into the bucket configured for this client, using a user-defined psk
func (cli *Client) UploadPartWithPsk(ctx context.Context, req *UploadPartRequest, payload []byte, psk []byte) (MultipartUploadResponse, error) {
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
		return MultipartUploadResponse{}, NewError(err, logData)
	}

	// Do the upload against AWS
	uploadPartOutput, err := cli.doUploadPart(
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
		return MultipartUploadResponse{}, NewError(err, logData)
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
		return MultipartUploadResponse{}, NewError(fmt.Errorf("error listing parts: %w", err), logData)
	}

	// If all parts have been uploaded, we call completeUpload
	parts := output.Parts
	if len(parts) == req.TotalChunks {
		return MultipartUploadResponse{
			Etag:             *uploadPartOutput.ETag,
			AllPartsUploaded: true,
		}, cli.completeUpload(ctx, uploadID, req, parts)
	}

	// Otherwise we don't need to perform any other operation.
	return MultipartUploadResponse{
		Etag:             *uploadPartOutput.ETag,
		AllPartsUploaded: false,
	}, nil
}

// doGetOrCreateMultipartUpload atomically gets the UploadId for the specified bucket
// and S3 object key, and if it does not find it, it creates it.
// The uploadID is returned. If an error happens, it will be wrapped and returned.
func (cli *Client) doGetOrCreateMultipartUpload(ctx context.Context, req *UploadPartRequest) (string, error) {
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
			return *upload.UploadId, nil
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
// The UploadPartOutput is returned. If an error happens, it will be wrapped and returned.
func (cli *Client) doUploadPart(ctx context.Context, input *s3.UploadPartInput, psk []byte) (*s3.UploadPartOutput, error) {
	if psk != nil {
		// Upload Part with PSK
		out, err := cli.cryptoClient.UploadPartWithPSK(input, psk)
		if err != nil {
			return nil, fmt.Errorf("error uploading part with psk: %w", err)
		}
		return out, nil
	}

	// Upload part without user-defined PSK
	out, err := cli.sdkClient.UploadPart(input)
	if err != nil {
		return nil, fmt.Errorf("error uploading part: %w", err)
	}
	return out, err
}

// CheckPartUploaded returns true only if the chunk corresponding to the provided chunkNumber has been uploaded.
// If all the chunks have been uploaded, we complete the upload operation.
// A boolean value which indicates if the call uploaded the last part is returned. If an error happens, it will be wrapped and returned.
func (cli *Client) CheckPartUploaded(ctx context.Context, req *UploadPartRequest) (bool, error) {
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
		return false, NewError(fmt.Errorf("list parts failed: %w", err), logData)
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
func (cli *Client) completeUpload(ctx context.Context, uploadID string, req *UploadPartRequest, parts []*s3.Part) error {
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

		log.Info(ctx, "attempting to complete multipart upload", log.Data{"complete": completeInput})

		_, err := cli.sdkClient.CompleteMultipartUpload(completeInput)
		if err != nil {
			return NewError(fmt.Errorf("error completing multipart upload: %w", err), log.Data{"complete": completeInput})
		}
	}
	return nil
}
