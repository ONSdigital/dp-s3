// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	v3 "github.com/ONSdigital/dp-s3/v3"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"sync"
)

// Ensure, that S3SDKClientMock does implement v3.S3SDKClient.
// If this is not the case, regenerate this file with moq.
var _ v3.S3SDKClient = &S3SDKClientMock{}

// S3SDKClientMock is a mock implementation of v3.S3SDKClient.
//
//	func TestSomethingThatUsesS3SDKClient(t *testing.T) {
//
//		// make and configure a mocked v3.S3SDKClient
//		mockedS3SDKClient := &S3SDKClientMock{
//			CompleteMultipartUploadFunc: func(ctx context.Context, in *s3.CompleteMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CompleteMultipartUploadOutput, error) {
//				panic("mock out the CompleteMultipartUpload method")
//			},
//			CreateMultipartUploadFunc: func(ctx context.Context, in *s3.CreateMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CreateMultipartUploadOutput, error) {
//				panic("mock out the CreateMultipartUpload method")
//			},
//			GetBucketPolicyFunc: func(ctx context.Context, in *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
//				panic("mock out the GetBucketPolicy method")
//			},
//			GetObjectFunc: func(ctx context.Context, in *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
//				panic("mock out the GetObject method")
//			},
//			HeadBucketFunc: func(ctx context.Context, in *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
//				panic("mock out the HeadBucket method")
//			},
//			HeadObjectFunc: func(ctx context.Context, in *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
//				panic("mock out the HeadObject method")
//			},
//			ListMultipartUploadsFunc: func(ctx context.Context, in *s3.ListMultipartUploadsInput, optFns ...func(*s3.Options)) (*s3.ListMultipartUploadsOutput, error) {
//				panic("mock out the ListMultipartUploads method")
//			},
//			ListObjectsFunc: func(ctx context.Context, in *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
//				panic("mock out the ListObjects method")
//			},
//			ListPartsFunc: func(ctx context.Context, in *s3.ListPartsInput, optFns ...func(*s3.Options)) (*s3.ListPartsOutput, error) {
//				panic("mock out the ListParts method")
//			},
//			PutBucketPolicyFunc: func(ctx context.Context, in *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error) {
//				panic("mock out the PutBucketPolicy method")
//			},
//			UploadPartFunc: func(ctx context.Context, in *s3.UploadPartInput, optFns ...func(*s3.Options)) (*s3.UploadPartOutput, error) {
//				panic("mock out the UploadPart method")
//			},
//		}
//
//		// use mockedS3SDKClient in code that requires v3.S3SDKClient
//		// and then make assertions.
//
//	}
type S3SDKClientMock struct {
	// CompleteMultipartUploadFunc mocks the CompleteMultipartUpload method.
	CompleteMultipartUploadFunc func(ctx context.Context, in *s3.CompleteMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CompleteMultipartUploadOutput, error)

	// CreateMultipartUploadFunc mocks the CreateMultipartUpload method.
	CreateMultipartUploadFunc func(ctx context.Context, in *s3.CreateMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CreateMultipartUploadOutput, error)

	// GetBucketPolicyFunc mocks the GetBucketPolicy method.
	GetBucketPolicyFunc func(ctx context.Context, in *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error)

	// GetObjectFunc mocks the GetObject method.
	GetObjectFunc func(ctx context.Context, in *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)

	// HeadBucketFunc mocks the HeadBucket method.
	HeadBucketFunc func(ctx context.Context, in *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error)

	// HeadObjectFunc mocks the HeadObject method.
	HeadObjectFunc func(ctx context.Context, in *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)

	// ListMultipartUploadsFunc mocks the ListMultipartUploads method.
	ListMultipartUploadsFunc func(ctx context.Context, in *s3.ListMultipartUploadsInput, optFns ...func(*s3.Options)) (*s3.ListMultipartUploadsOutput, error)

	// ListObjectsFunc mocks the ListObjects method.
	ListObjectsFunc func(ctx context.Context, in *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error)

	// ListPartsFunc mocks the ListParts method.
	ListPartsFunc func(ctx context.Context, in *s3.ListPartsInput, optFns ...func(*s3.Options)) (*s3.ListPartsOutput, error)

	// PutBucketPolicyFunc mocks the PutBucketPolicy method.
	PutBucketPolicyFunc func(ctx context.Context, in *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error)

	// UploadPartFunc mocks the UploadPart method.
	UploadPartFunc func(ctx context.Context, in *s3.UploadPartInput, optFns ...func(*s3.Options)) (*s3.UploadPartOutput, error)

	// calls tracks calls to the methods.
	calls struct {
		// CompleteMultipartUpload holds details about calls to the CompleteMultipartUpload method.
		CompleteMultipartUpload []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *s3.CompleteMultipartUploadInput
			// OptFns is the optFns argument value.
			OptFns []func(*s3.Options)
		}
		// CreateMultipartUpload holds details about calls to the CreateMultipartUpload method.
		CreateMultipartUpload []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *s3.CreateMultipartUploadInput
			// OptFns is the optFns argument value.
			OptFns []func(*s3.Options)
		}
		// GetBucketPolicy holds details about calls to the GetBucketPolicy method.
		GetBucketPolicy []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *s3.GetBucketPolicyInput
			// OptFns is the optFns argument value.
			OptFns []func(*s3.Options)
		}
		// GetObject holds details about calls to the GetObject method.
		GetObject []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *s3.GetObjectInput
			// OptFns is the optFns argument value.
			OptFns []func(*s3.Options)
		}
		// HeadBucket holds details about calls to the HeadBucket method.
		HeadBucket []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *s3.HeadBucketInput
			// OptFns is the optFns argument value.
			OptFns []func(*s3.Options)
		}
		// HeadObject holds details about calls to the HeadObject method.
		HeadObject []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *s3.HeadObjectInput
			// OptFns is the optFns argument value.
			OptFns []func(*s3.Options)
		}
		// ListMultipartUploads holds details about calls to the ListMultipartUploads method.
		ListMultipartUploads []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *s3.ListMultipartUploadsInput
			// OptFns is the optFns argument value.
			OptFns []func(*s3.Options)
		}
		// ListObjects holds details about calls to the ListObjects method.
		ListObjects []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *s3.ListObjectsInput
			// OptFns is the optFns argument value.
			OptFns []func(*s3.Options)
		}
		// ListParts holds details about calls to the ListParts method.
		ListParts []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *s3.ListPartsInput
			// OptFns is the optFns argument value.
			OptFns []func(*s3.Options)
		}
		// PutBucketPolicy holds details about calls to the PutBucketPolicy method.
		PutBucketPolicy []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *s3.PutBucketPolicyInput
			// OptFns is the optFns argument value.
			OptFns []func(*s3.Options)
		}
		// UploadPart holds details about calls to the UploadPart method.
		UploadPart []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// In is the in argument value.
			In *s3.UploadPartInput
			// OptFns is the optFns argument value.
			OptFns []func(*s3.Options)
		}
	}
	lockCompleteMultipartUpload sync.RWMutex
	lockCreateMultipartUpload   sync.RWMutex
	lockGetBucketPolicy         sync.RWMutex
	lockGetObject               sync.RWMutex
	lockHeadBucket              sync.RWMutex
	lockHeadObject              sync.RWMutex
	lockListMultipartUploads    sync.RWMutex
	lockListObjects             sync.RWMutex
	lockListParts               sync.RWMutex
	lockPutBucketPolicy         sync.RWMutex
	lockUploadPart              sync.RWMutex
}

// CompleteMultipartUpload calls CompleteMultipartUploadFunc.
func (mock *S3SDKClientMock) CompleteMultipartUpload(ctx context.Context, in *s3.CompleteMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CompleteMultipartUploadOutput, error) {
	if mock.CompleteMultipartUploadFunc == nil {
		panic("S3SDKClientMock.CompleteMultipartUploadFunc: method is nil but S3SDKClient.CompleteMultipartUpload was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		In     *s3.CompleteMultipartUploadInput
		OptFns []func(*s3.Options)
	}{
		Ctx:    ctx,
		In:     in,
		OptFns: optFns,
	}
	mock.lockCompleteMultipartUpload.Lock()
	mock.calls.CompleteMultipartUpload = append(mock.calls.CompleteMultipartUpload, callInfo)
	mock.lockCompleteMultipartUpload.Unlock()
	return mock.CompleteMultipartUploadFunc(ctx, in, optFns...)
}

// CompleteMultipartUploadCalls gets all the calls that were made to CompleteMultipartUpload.
// Check the length with:
//
//	len(mockedS3SDKClient.CompleteMultipartUploadCalls())
func (mock *S3SDKClientMock) CompleteMultipartUploadCalls() []struct {
	Ctx    context.Context
	In     *s3.CompleteMultipartUploadInput
	OptFns []func(*s3.Options)
} {
	var calls []struct {
		Ctx    context.Context
		In     *s3.CompleteMultipartUploadInput
		OptFns []func(*s3.Options)
	}
	mock.lockCompleteMultipartUpload.RLock()
	calls = mock.calls.CompleteMultipartUpload
	mock.lockCompleteMultipartUpload.RUnlock()
	return calls
}

// CreateMultipartUpload calls CreateMultipartUploadFunc.
func (mock *S3SDKClientMock) CreateMultipartUpload(ctx context.Context, in *s3.CreateMultipartUploadInput, optFns ...func(*s3.Options)) (*s3.CreateMultipartUploadOutput, error) {
	if mock.CreateMultipartUploadFunc == nil {
		panic("S3SDKClientMock.CreateMultipartUploadFunc: method is nil but S3SDKClient.CreateMultipartUpload was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		In     *s3.CreateMultipartUploadInput
		OptFns []func(*s3.Options)
	}{
		Ctx:    ctx,
		In:     in,
		OptFns: optFns,
	}
	mock.lockCreateMultipartUpload.Lock()
	mock.calls.CreateMultipartUpload = append(mock.calls.CreateMultipartUpload, callInfo)
	mock.lockCreateMultipartUpload.Unlock()
	return mock.CreateMultipartUploadFunc(ctx, in, optFns...)
}

// CreateMultipartUploadCalls gets all the calls that were made to CreateMultipartUpload.
// Check the length with:
//
//	len(mockedS3SDKClient.CreateMultipartUploadCalls())
func (mock *S3SDKClientMock) CreateMultipartUploadCalls() []struct {
	Ctx    context.Context
	In     *s3.CreateMultipartUploadInput
	OptFns []func(*s3.Options)
} {
	var calls []struct {
		Ctx    context.Context
		In     *s3.CreateMultipartUploadInput
		OptFns []func(*s3.Options)
	}
	mock.lockCreateMultipartUpload.RLock()
	calls = mock.calls.CreateMultipartUpload
	mock.lockCreateMultipartUpload.RUnlock()
	return calls
}

// GetBucketPolicy calls GetBucketPolicyFunc.
func (mock *S3SDKClientMock) GetBucketPolicy(ctx context.Context, in *s3.GetBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.GetBucketPolicyOutput, error) {
	if mock.GetBucketPolicyFunc == nil {
		panic("S3SDKClientMock.GetBucketPolicyFunc: method is nil but S3SDKClient.GetBucketPolicy was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		In     *s3.GetBucketPolicyInput
		OptFns []func(*s3.Options)
	}{
		Ctx:    ctx,
		In:     in,
		OptFns: optFns,
	}
	mock.lockGetBucketPolicy.Lock()
	mock.calls.GetBucketPolicy = append(mock.calls.GetBucketPolicy, callInfo)
	mock.lockGetBucketPolicy.Unlock()
	return mock.GetBucketPolicyFunc(ctx, in, optFns...)
}

// GetBucketPolicyCalls gets all the calls that were made to GetBucketPolicy.
// Check the length with:
//
//	len(mockedS3SDKClient.GetBucketPolicyCalls())
func (mock *S3SDKClientMock) GetBucketPolicyCalls() []struct {
	Ctx    context.Context
	In     *s3.GetBucketPolicyInput
	OptFns []func(*s3.Options)
} {
	var calls []struct {
		Ctx    context.Context
		In     *s3.GetBucketPolicyInput
		OptFns []func(*s3.Options)
	}
	mock.lockGetBucketPolicy.RLock()
	calls = mock.calls.GetBucketPolicy
	mock.lockGetBucketPolicy.RUnlock()
	return calls
}

// GetObject calls GetObjectFunc.
func (mock *S3SDKClientMock) GetObject(ctx context.Context, in *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	if mock.GetObjectFunc == nil {
		panic("S3SDKClientMock.GetObjectFunc: method is nil but S3SDKClient.GetObject was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		In     *s3.GetObjectInput
		OptFns []func(*s3.Options)
	}{
		Ctx:    ctx,
		In:     in,
		OptFns: optFns,
	}
	mock.lockGetObject.Lock()
	mock.calls.GetObject = append(mock.calls.GetObject, callInfo)
	mock.lockGetObject.Unlock()
	return mock.GetObjectFunc(ctx, in, optFns...)
}

// GetObjectCalls gets all the calls that were made to GetObject.
// Check the length with:
//
//	len(mockedS3SDKClient.GetObjectCalls())
func (mock *S3SDKClientMock) GetObjectCalls() []struct {
	Ctx    context.Context
	In     *s3.GetObjectInput
	OptFns []func(*s3.Options)
} {
	var calls []struct {
		Ctx    context.Context
		In     *s3.GetObjectInput
		OptFns []func(*s3.Options)
	}
	mock.lockGetObject.RLock()
	calls = mock.calls.GetObject
	mock.lockGetObject.RUnlock()
	return calls
}

// HeadBucket calls HeadBucketFunc.
func (mock *S3SDKClientMock) HeadBucket(ctx context.Context, in *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
	if mock.HeadBucketFunc == nil {
		panic("S3SDKClientMock.HeadBucketFunc: method is nil but S3SDKClient.HeadBucket was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		In     *s3.HeadBucketInput
		OptFns []func(*s3.Options)
	}{
		Ctx:    ctx,
		In:     in,
		OptFns: optFns,
	}
	mock.lockHeadBucket.Lock()
	mock.calls.HeadBucket = append(mock.calls.HeadBucket, callInfo)
	mock.lockHeadBucket.Unlock()
	return mock.HeadBucketFunc(ctx, in, optFns...)
}

// HeadBucketCalls gets all the calls that were made to HeadBucket.
// Check the length with:
//
//	len(mockedS3SDKClient.HeadBucketCalls())
func (mock *S3SDKClientMock) HeadBucketCalls() []struct {
	Ctx    context.Context
	In     *s3.HeadBucketInput
	OptFns []func(*s3.Options)
} {
	var calls []struct {
		Ctx    context.Context
		In     *s3.HeadBucketInput
		OptFns []func(*s3.Options)
	}
	mock.lockHeadBucket.RLock()
	calls = mock.calls.HeadBucket
	mock.lockHeadBucket.RUnlock()
	return calls
}

// HeadObject calls HeadObjectFunc.
func (mock *S3SDKClientMock) HeadObject(ctx context.Context, in *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	if mock.HeadObjectFunc == nil {
		panic("S3SDKClientMock.HeadObjectFunc: method is nil but S3SDKClient.HeadObject was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		In     *s3.HeadObjectInput
		OptFns []func(*s3.Options)
	}{
		Ctx:    ctx,
		In:     in,
		OptFns: optFns,
	}
	mock.lockHeadObject.Lock()
	mock.calls.HeadObject = append(mock.calls.HeadObject, callInfo)
	mock.lockHeadObject.Unlock()
	return mock.HeadObjectFunc(ctx, in, optFns...)
}

// HeadObjectCalls gets all the calls that were made to HeadObject.
// Check the length with:
//
//	len(mockedS3SDKClient.HeadObjectCalls())
func (mock *S3SDKClientMock) HeadObjectCalls() []struct {
	Ctx    context.Context
	In     *s3.HeadObjectInput
	OptFns []func(*s3.Options)
} {
	var calls []struct {
		Ctx    context.Context
		In     *s3.HeadObjectInput
		OptFns []func(*s3.Options)
	}
	mock.lockHeadObject.RLock()
	calls = mock.calls.HeadObject
	mock.lockHeadObject.RUnlock()
	return calls
}

// ListMultipartUploads calls ListMultipartUploadsFunc.
func (mock *S3SDKClientMock) ListMultipartUploads(ctx context.Context, in *s3.ListMultipartUploadsInput, optFns ...func(*s3.Options)) (*s3.ListMultipartUploadsOutput, error) {
	if mock.ListMultipartUploadsFunc == nil {
		panic("S3SDKClientMock.ListMultipartUploadsFunc: method is nil but S3SDKClient.ListMultipartUploads was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		In     *s3.ListMultipartUploadsInput
		OptFns []func(*s3.Options)
	}{
		Ctx:    ctx,
		In:     in,
		OptFns: optFns,
	}
	mock.lockListMultipartUploads.Lock()
	mock.calls.ListMultipartUploads = append(mock.calls.ListMultipartUploads, callInfo)
	mock.lockListMultipartUploads.Unlock()
	return mock.ListMultipartUploadsFunc(ctx, in, optFns...)
}

// ListMultipartUploadsCalls gets all the calls that were made to ListMultipartUploads.
// Check the length with:
//
//	len(mockedS3SDKClient.ListMultipartUploadsCalls())
func (mock *S3SDKClientMock) ListMultipartUploadsCalls() []struct {
	Ctx    context.Context
	In     *s3.ListMultipartUploadsInput
	OptFns []func(*s3.Options)
} {
	var calls []struct {
		Ctx    context.Context
		In     *s3.ListMultipartUploadsInput
		OptFns []func(*s3.Options)
	}
	mock.lockListMultipartUploads.RLock()
	calls = mock.calls.ListMultipartUploads
	mock.lockListMultipartUploads.RUnlock()
	return calls
}

// ListObjects calls ListObjectsFunc.
func (mock *S3SDKClientMock) ListObjects(ctx context.Context, in *s3.ListObjectsInput, optFns ...func(*s3.Options)) (*s3.ListObjectsOutput, error) {
	if mock.ListObjectsFunc == nil {
		panic("S3SDKClientMock.ListObjectsFunc: method is nil but S3SDKClient.ListObjects was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		In     *s3.ListObjectsInput
		OptFns []func(*s3.Options)
	}{
		Ctx:    ctx,
		In:     in,
		OptFns: optFns,
	}
	mock.lockListObjects.Lock()
	mock.calls.ListObjects = append(mock.calls.ListObjects, callInfo)
	mock.lockListObjects.Unlock()
	return mock.ListObjectsFunc(ctx, in, optFns...)
}

// ListObjectsCalls gets all the calls that were made to ListObjects.
// Check the length with:
//
//	len(mockedS3SDKClient.ListObjectsCalls())
func (mock *S3SDKClientMock) ListObjectsCalls() []struct {
	Ctx    context.Context
	In     *s3.ListObjectsInput
	OptFns []func(*s3.Options)
} {
	var calls []struct {
		Ctx    context.Context
		In     *s3.ListObjectsInput
		OptFns []func(*s3.Options)
	}
	mock.lockListObjects.RLock()
	calls = mock.calls.ListObjects
	mock.lockListObjects.RUnlock()
	return calls
}

// ListParts calls ListPartsFunc.
func (mock *S3SDKClientMock) ListParts(ctx context.Context, in *s3.ListPartsInput, optFns ...func(*s3.Options)) (*s3.ListPartsOutput, error) {
	if mock.ListPartsFunc == nil {
		panic("S3SDKClientMock.ListPartsFunc: method is nil but S3SDKClient.ListParts was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		In     *s3.ListPartsInput
		OptFns []func(*s3.Options)
	}{
		Ctx:    ctx,
		In:     in,
		OptFns: optFns,
	}
	mock.lockListParts.Lock()
	mock.calls.ListParts = append(mock.calls.ListParts, callInfo)
	mock.lockListParts.Unlock()
	return mock.ListPartsFunc(ctx, in, optFns...)
}

// ListPartsCalls gets all the calls that were made to ListParts.
// Check the length with:
//
//	len(mockedS3SDKClient.ListPartsCalls())
func (mock *S3SDKClientMock) ListPartsCalls() []struct {
	Ctx    context.Context
	In     *s3.ListPartsInput
	OptFns []func(*s3.Options)
} {
	var calls []struct {
		Ctx    context.Context
		In     *s3.ListPartsInput
		OptFns []func(*s3.Options)
	}
	mock.lockListParts.RLock()
	calls = mock.calls.ListParts
	mock.lockListParts.RUnlock()
	return calls
}

// PutBucketPolicy calls PutBucketPolicyFunc.
func (mock *S3SDKClientMock) PutBucketPolicy(ctx context.Context, in *s3.PutBucketPolicyInput, optFns ...func(*s3.Options)) (*s3.PutBucketPolicyOutput, error) {
	if mock.PutBucketPolicyFunc == nil {
		panic("S3SDKClientMock.PutBucketPolicyFunc: method is nil but S3SDKClient.PutBucketPolicy was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		In     *s3.PutBucketPolicyInput
		OptFns []func(*s3.Options)
	}{
		Ctx:    ctx,
		In:     in,
		OptFns: optFns,
	}
	mock.lockPutBucketPolicy.Lock()
	mock.calls.PutBucketPolicy = append(mock.calls.PutBucketPolicy, callInfo)
	mock.lockPutBucketPolicy.Unlock()
	return mock.PutBucketPolicyFunc(ctx, in, optFns...)
}

// PutBucketPolicyCalls gets all the calls that were made to PutBucketPolicy.
// Check the length with:
//
//	len(mockedS3SDKClient.PutBucketPolicyCalls())
func (mock *S3SDKClientMock) PutBucketPolicyCalls() []struct {
	Ctx    context.Context
	In     *s3.PutBucketPolicyInput
	OptFns []func(*s3.Options)
} {
	var calls []struct {
		Ctx    context.Context
		In     *s3.PutBucketPolicyInput
		OptFns []func(*s3.Options)
	}
	mock.lockPutBucketPolicy.RLock()
	calls = mock.calls.PutBucketPolicy
	mock.lockPutBucketPolicy.RUnlock()
	return calls
}

// UploadPart calls UploadPartFunc.
func (mock *S3SDKClientMock) UploadPart(ctx context.Context, in *s3.UploadPartInput, optFns ...func(*s3.Options)) (*s3.UploadPartOutput, error) {
	if mock.UploadPartFunc == nil {
		panic("S3SDKClientMock.UploadPartFunc: method is nil but S3SDKClient.UploadPart was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		In     *s3.UploadPartInput
		OptFns []func(*s3.Options)
	}{
		Ctx:    ctx,
		In:     in,
		OptFns: optFns,
	}
	mock.lockUploadPart.Lock()
	mock.calls.UploadPart = append(mock.calls.UploadPart, callInfo)
	mock.lockUploadPart.Unlock()
	return mock.UploadPartFunc(ctx, in, optFns...)
}

// UploadPartCalls gets all the calls that were made to UploadPart.
// Check the length with:
//
//	len(mockedS3SDKClient.UploadPartCalls())
func (mock *S3SDKClientMock) UploadPartCalls() []struct {
	Ctx    context.Context
	In     *s3.UploadPartInput
	OptFns []func(*s3.Options)
} {
	var calls []struct {
		Ctx    context.Context
		In     *s3.UploadPartInput
		OptFns []func(*s3.Options)
	}
	mock.lockUploadPart.RLock()
	calls = mock.calls.UploadPart
	mock.lockUploadPart.RUnlock()
	return calls
}
