/*
File copied form s3crypto repository
Original repo: https://github.com/ONSdigital/s3crypto
*/
package crypto

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	encryptionKeyHeader = "Pskencrypted"

	maxChunkSize = 5 * 1024 * 1024
)

// ErrNoPrivateKey is returned when an attempt is made to access a method that requires a private key when it has not been provided
var ErrNoPrivateKey = errors.New("you have not provided a private key and therefore do not have permission to complete this action")

// ErrNoMetadataPSK is returned when the file you are trying to download is not encrypted
var ErrNoMetadataPSK = errors.New("no encrypted key found for this file, you are trying to download a file which is not encrypted")

// Config represents the configuration items for the
// CryptoClient
type Config struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey

	HasUserDefinedPSK  bool
	MultipartChunkSize int
}

// CryptoClient provides a wrapper to the aws-sdk-go-v2 S3
// object
type CryptoClient struct {
	s3Client *s3.Client

	privKey           *rsa.PrivateKey
	publicKey         *rsa.PublicKey
	hasUserDefinedPSK bool
	chunkSize         int
}

type cryptoReader struct {
	s3Reader io.ReadCloser

	psk       []byte
	chunkSize int

	currChunk []byte
}

func (r *cryptoReader) Read(b []byte) (int, error) {
	if r.chunkSize == 0 {
		r.chunkSize = maxChunkSize
	}

	if len(r.currChunk) == 0 {
		p := make([]byte, r.chunkSize)

		n, err := io.ReadFull(r.s3Reader, p)
		if err != nil && err != io.ErrUnexpectedEOF {
			return n, err
		}

		unencryptedChunk, err := decryptObjectContent(r.psk, io.NopCloser(bytes.NewReader(p[:n])))
		if err != nil {
			return 0, err
		}

		r.currChunk = unencryptedChunk
	}

	var n int
	if len(r.currChunk) >= len(b) {
		copy(b, r.currChunk[:len(b)])
		n = len(b)
		r.currChunk = r.currChunk[len(b):]
	} else {
		copy(b, r.currChunk)
		n = len(r.currChunk)
		r.currChunk = nil
	}

	return n, nil
}

type encryptoReader struct {
	s3Reader io.Reader

	psk       []byte
	chunkSize int
	lastChunk bool

	currChunk []byte
}

func (r *encryptoReader) Read(b []byte) (int, error) {
	if r.lastChunk && len(r.currChunk) == 0 {
		return 0, io.EOF
	}

	if r.chunkSize == 0 {
		r.chunkSize = maxChunkSize
	}

	if len(r.currChunk) == 0 {
		p := make([]byte, r.chunkSize)

		n, err := io.ReadFull(r.s3Reader, p)
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				r.lastChunk = true
			} else {
				return n, err
			}
		}

		unencryptedChunk, err := encryptObjectContent(r.psk, bytes.NewReader(p[:n]))
		if err != nil {
			return 0, err
		}

		r.currChunk = unencryptedChunk
	}

	var n int
	if len(r.currChunk) >= len(b) {
		copy(b, r.currChunk[:len(b)])
		n = len(b)
		r.currChunk = r.currChunk[len(b):]
	} else {
		copy(b, r.currChunk)
		n = len(r.currChunk)
		r.currChunk = nil
	}

	return n, nil
}

func (r *cryptoReader) Close() error {
	return r.s3Reader.Close()
}

// Uploader provides a wrapper to the aws-sdk-go-v2 manager uploader
// for encryption
type Uploader struct {
	*CryptoClient

	s3uploader manager.Uploader
}

// New supports the creation of an Encryption supported client
// with a given aws config and rsa Private Key.
func New(awsConfig aws.Config, cfg *Config, optFns ...func(*s3.Options)) *CryptoClient {
	if cfg.MultipartChunkSize == 0 {
		cfg.MultipartChunkSize = maxChunkSize
	}
	cc := &CryptoClient{s3.NewFromConfig(awsConfig, optFns...), cfg.PrivateKey, cfg.PublicKey, cfg.HasUserDefinedPSK, cfg.MultipartChunkSize}

	if cc.privKey != nil {
		cc.publicKey = &cc.privKey.PublicKey
	}

	return cc
}

// NewUploader creates a new instance of the crypto Uploader
func NewUploader(awsConfig aws.Config, cfg *Config, optFns ...func(*s3.Options)) *Uploader {
	cc := &CryptoClient{s3.NewFromConfig(awsConfig, optFns...), cfg.PrivateKey, cfg.PublicKey, cfg.HasUserDefinedPSK, cfg.MultipartChunkSize}

	if cc.privKey != nil {
		cc.publicKey = &cc.privKey.PublicKey
	}

	return &Uploader{
		CryptoClient: cc,

		s3uploader: *manager.NewUploader(s3.NewFromConfig(awsConfig, optFns...)),
	}
}

// CreateMultipartUploadRequest wraps the SDK method by creating a PSK which
// is encrypted using the public key and stored as metadata against the completed
// object, as well as temporarily being stored as its own object while the Multipart
// upload is being updated.
func (c *CryptoClient) CreateMultipartUpload(ctx context.Context, input *s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error) {
	if !c.hasUserDefinedPSK {
		psk := createPSK()

		ekStr, err := c.encryptKey(psk)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt PSK: %w", err)
		}

		if input.Metadata == nil {
			input.Metadata = make(map[string]string)
		}
		input.Metadata[encryptionKeyHeader] = ekStr

		if err := c.storeEncryptedKey(ctx, input, ekStr); err != nil {
			return nil, fmt.Errorf("failed to store encrypted PSK: %w", err)
		}
	}

	out, err := c.s3Client.CreateMultipartUpload(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("CreateMultipartUpload failed: %w", err)
	}

	return out, nil
}

// UploadPartRequest wraps the SDK method by retrieving the encrypted PSK from the temporary
// object, decrypting the PSK using the private key, before stream encoding the content
// for the particular part
func (c *CryptoClient) UploadPartRequest(ctx context.Context, input *s3.UploadPartInput) (*s3.UploadPartOutput, error) {
	ekStr, err := c.getEncryptedKey(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt PSK: %w", err)
	}

	psk, err := c.decryptKey(ekStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt PSK: %w", err)
	}

	encryptedContent, err := encryptObjectContent(psk, input.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt content: %w", err)
	}

	input.Body = bytes.NewReader(encryptedContent)

	out, err := c.s3Client.UploadPart(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("UploadPart failed: %w", err)
	}

	return out, nil
}

// UploadPartRequestWithPSK wraps the SDK method encrypting the part contents with a user defined
// PSK
func (c *CryptoClient) UploadPartWithPSK(ctx context.Context, input *s3.UploadPartInput, psk []byte) (*s3.UploadPartOutput, error) {
	encryptedContent, err := encryptObjectContent(psk, input.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt content: %w", err)
	}

	input.Body = bytes.NewReader(encryptedContent)

	out, err := c.s3Client.UploadPart(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("UploadPart failed: %w", err)
	}

	return out, nil
}

// PutObjectRequest wraps the SDK method by creating a PSK, encrypting it using the public key,
// and encrypting the object content using the PSK
func (c *CryptoClient) PutObject(ctx context.Context, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	psk := createPSK()

	ekStr, err := c.encryptKey(psk)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt PSK: %w", err)
	}

	if input.Metadata == nil {
		input.Metadata = make(map[string]string)
	}
	input.Metadata[encryptionKeyHeader] = ekStr

	encryptedContent, err := encryptObjectContent(psk, input.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt content: %w", err)
	}

	input.Body = bytes.NewReader(encryptedContent)

	out, err := c.s3Client.PutObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("PutObject failed: %w", err)
	}

	return out, nil
}

// PutObjectRequestWithPSK wraps the SDK method by encrypting the object content with a user defined PSK
func (c *CryptoClient) PutObjectWithPSK(ctx context.Context, input *s3.PutObjectInput, psk []byte) (*s3.PutObjectOutput, error) {
	encryptedContent, err := encryptObjectContent(psk, input.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt content: %w", err)
	}

	input.Body = bytes.NewReader(encryptedContent)

	out, err := c.s3Client.PutObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("PutObject failed: %w", err)
	}

	return out, nil
}

// GetObjectRequest wraps the SDK method by retrieving the encrypted PSK from the object metadata.
// The PSK is then decrypted, and is then used to decrypt the content of the object.
func (c *CryptoClient) GetObject(ctx context.Context, input *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	out, err := c.s3Client.GetObject(ctx, input)

	ekStr := out.Metadata[encryptionKeyHeader]

	if ekStr == "" {
		return nil, ErrNoMetadataPSK
	}

	psk, err := c.decryptKey(ekStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt PSK: %w", err)
	}

	var content io.Reader
	if c.chunkSize > 0 {
		content, err = decryptObjectContentChunks(c.chunkSize, psk, out.Body)
	} else {
		content, err = decryptObjectContentChunks(maxChunkSize, psk, out.Body)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt content: %w", err)
	}

	out.Body = io.NopCloser(content)

	return out, nil
}

// GetObjectRequestWithPSK wraps the SDK method by decrypting the retrieved object content with the given PSK
func (c *CryptoClient) GetObjectWithPSK(ctx context.Context, input *s3.GetObjectInput, psk []byte) (*s3.GetObjectOutput, error) {
	out, err := c.s3Client.GetObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("GetObject failed: %w", err)
	}

	out.Body = &cryptoReader{
		s3Reader:  out.Body,
		psk:       psk,
		chunkSize: c.chunkSize,
	}

	return out, nil
}

func (c *CryptoClient) storeEncryptedKey(ctx context.Context, input *s3.CreateMultipartUploadInput, key string) error {
	keyFileName := *input.Key + ".key"

	objectInput := &s3.PutObjectInput{
		Body:   strings.NewReader(key),
		Bucket: input.Bucket,
		Key:    &keyFileName,
	}

	_, err := c.s3Client.PutObject(ctx, objectInput)
	if err != nil {
		return fmt.Errorf("failed to store encrypted PSK: %w", err)
	}

	return nil
}

func (c *CryptoClient) getEncryptedKey(ctx context.Context, input *s3.UploadPartInput) (string, error) {
	keyFileName := *input.Key + ".key"

	objectInput := &s3.GetObjectInput{
		Bucket: input.Bucket,
		Key:    &keyFileName,
	}

	objectOutput, err := c.s3Client.GetObject(ctx, objectInput)
	if err != nil {
		return "", err
	}

	key, err := io.ReadAll(objectOutput.Body)
	if err != nil {
		return "", err
	}

	return string(key), nil
}

func (c *CryptoClient) removeEncryptedKey(ctx context.Context, input *s3.CompleteMultipartUploadInput) error {
	keyFileName := *input.Key + ".key"

	objectInput := &s3.DeleteObjectInput{
		Bucket: input.Bucket,
		Key:    &keyFileName,
	}

	_, err := c.s3Client.DeleteObject(ctx, objectInput)
	return err
}

func (c *CryptoClient) encryptKey(psk []byte) (string, error) {
	hash := sha1.New()
	encryptedKey, err := rsa.EncryptOAEP(hash, rand.Reader, c.publicKey, psk, []byte(""))
	if err != nil {
		return "", err
	}

	ekStr := hex.EncodeToString(encryptedKey)
	return ekStr, nil
}

func (c *CryptoClient) decryptKey(encryptedKeyHex string) ([]byte, error) {
	if c.privKey == nil {
		return nil, ErrNoPrivateKey
	}

	encryptedKey, err := hex.DecodeString(encryptedKeyHex)
	if err != nil {
		return nil, err
	}

	hash := sha1.New()
	return rsa.DecryptOAEP(hash, rand.Reader, c.privKey, encryptedKey, []byte(""))
}

// Upload provides a wrapper for the sdk method with encryption
func (u *Uploader) Upload(ctx context.Context, input *s3.PutObjectInput) (output *manager.UploadOutput, err error) {
	psk := createPSK()

	ekStr, err := u.CryptoClient.encryptKey(psk)
	if err != nil {
		return
	}

	input.Metadata = make(map[string]string)
	input.Metadata[encryptionKeyHeader] = ekStr

	encryptedContent, err := encryptObjectContent(psk, input.Body)
	if err != nil {
		return
	}

	input.Body = bytes.NewReader(encryptedContent)

	return u.s3uploader.Upload(ctx, input)
}

// UploadWithPSK allows you to encrypt the file with a given psk
func (u *Uploader) UploadWithPSK(ctx context.Context, input *s3.PutObjectInput, psk []byte) (output *manager.UploadOutput, err error) {
	input.Body = &encryptoReader{
		psk:      psk,
		s3Reader: input.Body,
	}

	return u.s3uploader.Upload(ctx, input)
}

func encryptObjectContent(psk []byte, b io.Reader) ([]byte, error) {
	unencryptedBytes, err := io.ReadAll(b)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(psk)
	if err != nil {
		return nil, err
	}

	encryptedBytes := make([]byte, len(unencryptedBytes))

	stream := cipher.NewCFBEncrypter(block, psk)

	stream.XORKeyStream(encryptedBytes, unencryptedBytes)

	return encryptedBytes, nil
}

func decryptObjectContentChunks(size int, psk []byte, r io.ReadCloser) (io.Reader, error) {

	p := make([]byte, size)

	buf := &bytes.Buffer{}
	for {
		n, err := io.ReadFull(r, p)
		if err != nil && err != io.ErrUnexpectedEOF {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		unencryptedChunk, err := decryptObjectContent(psk, io.NopCloser(bytes.NewReader(p[:n])))
		if err != nil {
			return nil, err
		}

		_, err = buf.Write(unencryptedChunk)
		if err != nil {
			return nil, err
		}
	}

	return buf, nil
}

func decryptObjectContent(psk []byte, b io.ReadCloser) ([]byte, error) {
	encryptedBytes, err := io.ReadAll(b)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(psk)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCFBDecrypter(block, psk)

	unencryptedBytes := make([]byte, len(encryptedBytes))
	stream.XORKeyStream(unencryptedBytes, encryptedBytes)

	return unencryptedBytes, nil
}

func createPSK() []byte {
	key := make([]byte, 16)
	rand.Read(key)

	return key
}
