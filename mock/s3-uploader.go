// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/ONSdigital/dp-s3/v2"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"sync"
)

// Ensure, that S3SDKUploaderMock does implement s3.S3SDKUploader.
// If this is not the case, regenerate this file with moq.
var _ s3.S3SDKUploader = &S3SDKUploaderMock{}

// S3SDKUploaderMock is a mock implementation of s3.S3SDKUploader.
//
// 	func TestSomethingThatUsesS3SDKUploader(t *testing.T) {
//
// 		// make and configure a mocked s3.S3SDKUploader
// 		mockedS3SDKUploader := &S3SDKUploaderMock{
// 			UploadFunc: func(in *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
// 				panic("mock out the Upload method")
// 			},
// 			UploadWithContextFunc: func(ctx context.Context, in *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
// 				panic("mock out the UploadWithContext method")
// 			},
// 		}
//
// 		// use mockedS3SDKUploader in code that requires s3.S3SDKUploader
// 		// and then make assertions.
//
// 	}
type S3SDKUploaderMock struct {
	// UploadFunc mocks the Upload method.
	UploadFunc func(in *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)

	// UploadWithContextFunc mocks the UploadWithContext method.
	UploadWithContextFunc func(ctx context.Context, in *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)

	// calls tracks calls to the methods.
	calls struct {
		// Upload holds details about calls to the Upload method.
		Upload []struct {
			// In is the in argument value.
			In *s3manager.UploadInput
			// Options is the options argument value.
			Options []func(*s3manager.Uploader)
		}
		// UploadWithContext holds details about calls to the UploadWithContext method.
		UploadWithContext []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *s3manager.UploadInput
			// Options is the options argument value.
			Options []func(*s3manager.Uploader)
		}
	}
	lockUpload            sync.RWMutex
	lockUploadWithContext sync.RWMutex
}

// Upload calls UploadFunc.
func (mock *S3SDKUploaderMock) Upload(in *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	if mock.UploadFunc == nil {
		panic("S3SDKUploaderMock.UploadFunc: method is nil but S3SDKUploader.Upload was just called")
	}
	callInfo := struct {
		In      *s3manager.UploadInput
		Options []func(*s3manager.Uploader)
	}{
		In:      in,
		Options: options,
	}
	mock.lockUpload.Lock()
	mock.calls.Upload = append(mock.calls.Upload, callInfo)
	mock.lockUpload.Unlock()
	return mock.UploadFunc(in, options...)
}

// UploadCalls gets all the calls that were made to Upload.
// Check the length with:
//     len(mockedS3SDKUploader.UploadCalls())
func (mock *S3SDKUploaderMock) UploadCalls() []struct {
	In      *s3manager.UploadInput
	Options []func(*s3manager.Uploader)
} {
	var calls []struct {
		In      *s3manager.UploadInput
		Options []func(*s3manager.Uploader)
	}
	mock.lockUpload.RLock()
	calls = mock.calls.Upload
	mock.lockUpload.RUnlock()
	return calls
}

// UploadWithContext calls UploadWithContextFunc.
func (mock *S3SDKUploaderMock) UploadWithContext(ctx context.Context, in *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	if mock.UploadWithContextFunc == nil {
		panic("S3SDKUploaderMock.UploadWithContextFunc: method is nil but S3SDKUploader.UploadWithContext was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		In      *s3manager.UploadInput
		Options []func(*s3manager.Uploader)
	}{
		Ctx:     ctx,
		In:      in,
		Options: options,
	}
	mock.lockUploadWithContext.Lock()
	mock.calls.UploadWithContext = append(mock.calls.UploadWithContext, callInfo)
	mock.lockUploadWithContext.Unlock()
	return mock.UploadWithContextFunc(ctx, in, options...)
}

// UploadWithContextCalls gets all the calls that were made to UploadWithContext.
// Check the length with:
//     len(mockedS3SDKUploader.UploadWithContextCalls())
func (mock *S3SDKUploaderMock) UploadWithContextCalls() []struct {
	Ctx     context.Context
	In      *s3manager.UploadInput
	Options []func(*s3manager.Uploader)
} {
	var calls []struct {
		Ctx     context.Context
		In      *s3manager.UploadInput
		Options []func(*s3manager.Uploader)
	}
	mock.lockUploadWithContext.RLock()
	calls = mock.calls.UploadWithContext
	mock.lockUploadWithContext.RUnlock()
	return calls
}
