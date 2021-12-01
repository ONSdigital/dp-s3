// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/ONSdigital/dp-s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"sync"
)

// Ensure, that S3CryptoUploaderMock does implement s3.S3CryptoUploader.
// If this is not the case, regenerate this file with moq.
var _ s3.S3CryptoUploader = &S3CryptoUploaderMock{}

// S3CryptoUploaderMock is a mock implementation of s3.S3CryptoUploader.
//
// 	func TestSomethingThatUsesS3CryptoUploader(t *testing.T) {
//
// 		// make and configure a mocked s3.S3CryptoUploader
// 		mockedS3CryptoUploader := &S3CryptoUploaderMock{
// 			UploadWithPSKFunc: func(ctx context.Context, in *s3manager.UploadInput, psk []byte) (*s3manager.UploadOutput, error) {
// 				panic("mock out the UploadWithPSK method")
// 			},
// 		}
//
// 		// use mockedS3CryptoUploader in code that requires s3.S3CryptoUploader
// 		// and then make assertions.
//
// 	}
type S3CryptoUploaderMock struct {
	// UploadWithPSKFunc mocks the UploadWithPSK method.
	UploadWithPSKFunc func(ctx context.Context, in *s3manager.UploadInput, psk []byte) (*s3manager.UploadOutput, error)

	// calls tracks calls to the methods.
	calls struct {
		// UploadWithPSK holds details about calls to the UploadWithPSK method.
		UploadWithPSK []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *s3manager.UploadInput
			// Psk is the psk argument value.
			Psk []byte
		}
	}
	lockUploadWithPSK sync.RWMutex
}

// UploadWithPSK calls UploadWithPSKFunc.
func (mock *S3CryptoUploaderMock) UploadWithPSK(ctx context.Context, in *s3manager.UploadInput, psk []byte) (*s3manager.UploadOutput, error) {
	if mock.UploadWithPSKFunc == nil {
		panic("S3CryptoUploaderMock.UploadWithPSKFunc: method is nil but S3CryptoUploader.UploadWithPSK was just called")
	}
	callInfo := struct {
		Ctx context.Context
		In  *s3manager.UploadInput
		Psk []byte
	}{
		Ctx: ctx,
		In:  in,
		Psk: psk,
	}
	mock.lockUploadWithPSK.Lock()
	mock.calls.UploadWithPSK = append(mock.calls.UploadWithPSK, callInfo)
	mock.lockUploadWithPSK.Unlock()
	return mock.UploadWithPSKFunc(ctx, in, psk)
}

// UploadWithPSKCalls gets all the calls that were made to UploadWithPSK.
// Check the length with:
//     len(mockedS3CryptoUploader.UploadWithPSKCalls())
func (mock *S3CryptoUploaderMock) UploadWithPSKCalls() []struct {
	Ctx context.Context
	In  *s3manager.UploadInput
	Psk []byte
} {
	var calls []struct {
		Ctx context.Context
		In  *s3manager.UploadInput
		Psk []byte
	}
	mock.lockUploadWithPSK.RLock()
	calls = mock.calls.UploadWithPSK
	mock.lockUploadWithPSK.RUnlock()
	return calls
}
