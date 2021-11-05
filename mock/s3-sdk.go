// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"github.com/ONSdigital/dp-s3"
	"github.com/aws/aws-sdk-go/service/s3"
	"sync"
)

var (
	lockS3SDKClientMockCompleteMultipartUpload sync.RWMutex
	lockS3SDKClientMockCreateMultipartUpload   sync.RWMutex
	lockS3SDKClientMockGetObject               sync.RWMutex
	lockS3SDKClientMockHeadBucket              sync.RWMutex
	lockS3SDKClientMockHeadObject              sync.RWMutex
	lockS3SDKClientMockListMultipartUploads    sync.RWMutex
	lockS3SDKClientMockListParts               sync.RWMutex
	lockS3SDKClientMockUploadPart              sync.RWMutex
)

// Ensure, that S3SDKClientMock does implement s3client.S3SDKClient.
// If this is not the case, regenerate this file with moq.
var _ s3client.S3SDKClient = &S3SDKClientMock{}

// S3SDKClientMock is a mock implementation of s3client.S3SDKClient.
//
//     func TestSomethingThatUsesS3SDKClient(t *testing.T) {
//
//         // make and configure a mocked s3client.S3SDKClient
//         mockedS3SDKClient := &S3SDKClientMock{
//             CompleteMultipartUploadFunc: func(in *s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error) {
// 	               panic("mock out the CompleteMultipartUpload method")
//             },
//             CreateMultipartUploadFunc: func(in *s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error) {
// 	               panic("mock out the CreateMultipartUpload method")
//             },
//             GetObjectFunc: func(in *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
// 	               panic("mock out the GetObject method")
//             },
//             HeadBucketFunc: func(in *s3.HeadBucketInput) (*s3.HeadBucketOutput, error) {
// 	               panic("mock out the HeadBucket method")
//             },
//             HeadObjectFunc: func(in *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
// 	               panic("mock out the HeadObject method")
//             },
//             ListMultipartUploadsFunc: func(in *s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error) {
// 	               panic("mock out the ListMultipartUploads method")
//             },
//             ListPartsFunc: func(in *s3.ListPartsInput) (*s3.ListPartsOutput, error) {
// 	               panic("mock out the ListParts method")
//             },
//             UploadPartFunc: func(in *s3.UploadPartInput) (*s3.UploadPartOutput, error) {
// 	               panic("mock out the UploadPart method")
//             },
//         }
//
//         // use mockedS3SDKClient in code that requires s3client.S3SDKClient
//         // and then make assertions.
//
//     }
type S3SDKClientMock struct {
	// CompleteMultipartUploadFunc mocks the CompleteMultipartUpload method.
	CompleteMultipartUploadFunc func(in *s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error)

	// CreateMultipartUploadFunc mocks the CreateMultipartUpload method.
	CreateMultipartUploadFunc func(in *s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error)

	// GetObjectFunc mocks the GetObject method.
	GetObjectFunc func(in *s3.GetObjectInput) (*s3.GetObjectOutput, error)

	// HeadBucketFunc mocks the HeadBucket method.
	HeadBucketFunc func(in *s3.HeadBucketInput) (*s3.HeadBucketOutput, error)

	// HeadObjectFunc mocks the HeadObject method.
	HeadObjectFunc func(in *s3.HeadObjectInput) (*s3.HeadObjectOutput, error)

	// ListMultipartUploadsFunc mocks the ListMultipartUploads method.
	ListMultipartUploadsFunc func(in *s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error)

	// ListPartsFunc mocks the ListParts method.
	ListPartsFunc func(in *s3.ListPartsInput) (*s3.ListPartsOutput, error)

	// UploadPartFunc mocks the UploadPart method.
	UploadPartFunc func(in *s3.UploadPartInput) (*s3.UploadPartOutput, error)

	// calls tracks calls to the methods.
	calls struct {
		// CompleteMultipartUpload holds details about calls to the CompleteMultipartUpload method.
		CompleteMultipartUpload []struct {
			// In is the in argument value.
			In *s3.CompleteMultipartUploadInput
		}
		// CreateMultipartUpload holds details about calls to the CreateMultipartUpload method.
		CreateMultipartUpload []struct {
			// In is the in argument value.
			In *s3.CreateMultipartUploadInput
		}
		// GetObject holds details about calls to the GetObject method.
		GetObject []struct {
			// In is the in argument value.
			In *s3.GetObjectInput
		}
		// HeadBucket holds details about calls to the HeadBucket method.
		HeadBucket []struct {
			// In is the in argument value.
			In *s3.HeadBucketInput
		}
		// HeadObject holds details about calls to the HeadObject method.
		HeadObject []struct {
			// In is the in argument value.
			In *s3.HeadObjectInput
		}
		// ListMultipartUploads holds details about calls to the ListMultipartUploads method.
		ListMultipartUploads []struct {
			// In is the in argument value.
			In *s3.ListMultipartUploadsInput
		}
		// ListParts holds details about calls to the ListParts method.
		ListParts []struct {
			// In is the in argument value.
			In *s3.ListPartsInput
		}
		// UploadPart holds details about calls to the UploadPart method.
		UploadPart []struct {
			// In is the in argument value.
			In *s3.UploadPartInput
		}
	}
}

// CompleteMultipartUpload calls CompleteMultipartUploadFunc.
func (mock *S3SDKClientMock) CompleteMultipartUpload(in *s3.CompleteMultipartUploadInput) (*s3.CompleteMultipartUploadOutput, error) {
	if mock.CompleteMultipartUploadFunc == nil {
		panic("S3SDKClientMock.CompleteMultipartUploadFunc: method is nil but S3SDKClient.CompleteMultipartUpload was just called")
	}
	callInfo := struct {
		In *s3.CompleteMultipartUploadInput
	}{
		In: in,
	}
	lockS3SDKClientMockCompleteMultipartUpload.Lock()
	mock.calls.CompleteMultipartUpload = append(mock.calls.CompleteMultipartUpload, callInfo)
	lockS3SDKClientMockCompleteMultipartUpload.Unlock()
	return mock.CompleteMultipartUploadFunc(in)
}

// CompleteMultipartUploadCalls gets all the calls that were made to CompleteMultipartUpload.
// Check the length with:
//     len(mockedS3SDKClient.CompleteMultipartUploadCalls())
func (mock *S3SDKClientMock) CompleteMultipartUploadCalls() []struct {
	In *s3.CompleteMultipartUploadInput
} {
	var calls []struct {
		In *s3.CompleteMultipartUploadInput
	}
	lockS3SDKClientMockCompleteMultipartUpload.RLock()
	calls = mock.calls.CompleteMultipartUpload
	lockS3SDKClientMockCompleteMultipartUpload.RUnlock()
	return calls
}

// CreateMultipartUpload calls CreateMultipartUploadFunc.
func (mock *S3SDKClientMock) CreateMultipartUpload(in *s3.CreateMultipartUploadInput) (*s3.CreateMultipartUploadOutput, error) {
	if mock.CreateMultipartUploadFunc == nil {
		panic("S3SDKClientMock.CreateMultipartUploadFunc: method is nil but S3SDKClient.CreateMultipartUpload was just called")
	}
	callInfo := struct {
		In *s3.CreateMultipartUploadInput
	}{
		In: in,
	}
	lockS3SDKClientMockCreateMultipartUpload.Lock()
	mock.calls.CreateMultipartUpload = append(mock.calls.CreateMultipartUpload, callInfo)
	lockS3SDKClientMockCreateMultipartUpload.Unlock()
	return mock.CreateMultipartUploadFunc(in)
}

// CreateMultipartUploadCalls gets all the calls that were made to CreateMultipartUpload.
// Check the length with:
//     len(mockedS3SDKClient.CreateMultipartUploadCalls())
func (mock *S3SDKClientMock) CreateMultipartUploadCalls() []struct {
	In *s3.CreateMultipartUploadInput
} {
	var calls []struct {
		In *s3.CreateMultipartUploadInput
	}
	lockS3SDKClientMockCreateMultipartUpload.RLock()
	calls = mock.calls.CreateMultipartUpload
	lockS3SDKClientMockCreateMultipartUpload.RUnlock()
	return calls
}

// GetObject calls GetObjectFunc.
func (mock *S3SDKClientMock) GetObject(in *s3.GetObjectInput) (*s3.GetObjectOutput, error) {
	if mock.GetObjectFunc == nil {
		panic("S3SDKClientMock.GetObjectFunc: method is nil but S3SDKClient.GetObject was just called")
	}
	callInfo := struct {
		In *s3.GetObjectInput
	}{
		In: in,
	}
	lockS3SDKClientMockGetObject.Lock()
	mock.calls.GetObject = append(mock.calls.GetObject, callInfo)
	lockS3SDKClientMockGetObject.Unlock()
	return mock.GetObjectFunc(in)
}

// GetObjectCalls gets all the calls that were made to GetObject.
// Check the length with:
//     len(mockedS3SDKClient.GetObjectCalls())
func (mock *S3SDKClientMock) GetObjectCalls() []struct {
	In *s3.GetObjectInput
} {
	var calls []struct {
		In *s3.GetObjectInput
	}
	lockS3SDKClientMockGetObject.RLock()
	calls = mock.calls.GetObject
	lockS3SDKClientMockGetObject.RUnlock()
	return calls
}

// HeadBucket calls HeadBucketFunc.
func (mock *S3SDKClientMock) HeadBucket(in *s3.HeadBucketInput) (*s3.HeadBucketOutput, error) {
	if mock.HeadBucketFunc == nil {
		panic("S3SDKClientMock.HeadBucketFunc: method is nil but S3SDKClient.HeadBucket was just called")
	}
	callInfo := struct {
		In *s3.HeadBucketInput
	}{
		In: in,
	}
	lockS3SDKClientMockHeadBucket.Lock()
	mock.calls.HeadBucket = append(mock.calls.HeadBucket, callInfo)
	lockS3SDKClientMockHeadBucket.Unlock()
	return mock.HeadBucketFunc(in)
}

// HeadBucketCalls gets all the calls that were made to HeadBucket.
// Check the length with:
//     len(mockedS3SDKClient.HeadBucketCalls())
func (mock *S3SDKClientMock) HeadBucketCalls() []struct {
	In *s3.HeadBucketInput
} {
	var calls []struct {
		In *s3.HeadBucketInput
	}
	lockS3SDKClientMockHeadBucket.RLock()
	calls = mock.calls.HeadBucket
	lockS3SDKClientMockHeadBucket.RUnlock()
	return calls
}

// HeadObject calls HeadObjectFunc.
func (mock *S3SDKClientMock) HeadObject(in *s3.HeadObjectInput) (*s3.HeadObjectOutput, error) {
	if mock.HeadObjectFunc == nil {
		panic("S3SDKClientMock.HeadObjectFunc: method is nil but S3SDKClient.HeadObject was just called")
	}
	callInfo := struct {
		In *s3.HeadObjectInput
	}{
		In: in,
	}
	lockS3SDKClientMockHeadObject.Lock()
	mock.calls.HeadObject = append(mock.calls.HeadObject, callInfo)
	lockS3SDKClientMockHeadObject.Unlock()
	return mock.HeadObjectFunc(in)
}

// HeadObjectCalls gets all the calls that were made to HeadObject.
// Check the length with:
//     len(mockedS3SDKClient.HeadObjectCalls())
func (mock *S3SDKClientMock) HeadObjectCalls() []struct {
	In *s3.HeadObjectInput
} {
	var calls []struct {
		In *s3.HeadObjectInput
	}
	lockS3SDKClientMockHeadObject.RLock()
	calls = mock.calls.HeadObject
	lockS3SDKClientMockHeadObject.RUnlock()
	return calls
}

// ListMultipartUploads calls ListMultipartUploadsFunc.
func (mock *S3SDKClientMock) ListMultipartUploads(in *s3.ListMultipartUploadsInput) (*s3.ListMultipartUploadsOutput, error) {
	if mock.ListMultipartUploadsFunc == nil {
		panic("S3SDKClientMock.ListMultipartUploadsFunc: method is nil but S3SDKClient.ListMultipartUploads was just called")
	}
	callInfo := struct {
		In *s3.ListMultipartUploadsInput
	}{
		In: in,
	}
	lockS3SDKClientMockListMultipartUploads.Lock()
	mock.calls.ListMultipartUploads = append(mock.calls.ListMultipartUploads, callInfo)
	lockS3SDKClientMockListMultipartUploads.Unlock()
	return mock.ListMultipartUploadsFunc(in)
}

// ListMultipartUploadsCalls gets all the calls that were made to ListMultipartUploads.
// Check the length with:
//     len(mockedS3SDKClient.ListMultipartUploadsCalls())
func (mock *S3SDKClientMock) ListMultipartUploadsCalls() []struct {
	In *s3.ListMultipartUploadsInput
} {
	var calls []struct {
		In *s3.ListMultipartUploadsInput
	}
	lockS3SDKClientMockListMultipartUploads.RLock()
	calls = mock.calls.ListMultipartUploads
	lockS3SDKClientMockListMultipartUploads.RUnlock()
	return calls
}

// ListParts calls ListPartsFunc.
func (mock *S3SDKClientMock) ListParts(in *s3.ListPartsInput) (*s3.ListPartsOutput, error) {
	if mock.ListPartsFunc == nil {
		panic("S3SDKClientMock.ListPartsFunc: method is nil but S3SDKClient.ListParts was just called")
	}
	callInfo := struct {
		In *s3.ListPartsInput
	}{
		In: in,
	}
	lockS3SDKClientMockListParts.Lock()
	mock.calls.ListParts = append(mock.calls.ListParts, callInfo)
	lockS3SDKClientMockListParts.Unlock()
	return mock.ListPartsFunc(in)
}

// ListPartsCalls gets all the calls that were made to ListParts.
// Check the length with:
//     len(mockedS3SDKClient.ListPartsCalls())
func (mock *S3SDKClientMock) ListPartsCalls() []struct {
	In *s3.ListPartsInput
} {
	var calls []struct {
		In *s3.ListPartsInput
	}
	lockS3SDKClientMockListParts.RLock()
	calls = mock.calls.ListParts
	lockS3SDKClientMockListParts.RUnlock()
	return calls
}

// UploadPart calls UploadPartFunc.
func (mock *S3SDKClientMock) UploadPart(in *s3.UploadPartInput) (*s3.UploadPartOutput, error) {
	if mock.UploadPartFunc == nil {
		panic("S3SDKClientMock.UploadPartFunc: method is nil but S3SDKClient.UploadPart was just called")
	}
	callInfo := struct {
		In *s3.UploadPartInput
	}{
		In: in,
	}
	lockS3SDKClientMockUploadPart.Lock()
	mock.calls.UploadPart = append(mock.calls.UploadPart, callInfo)
	lockS3SDKClientMockUploadPart.Unlock()
	return mock.UploadPartFunc(in)
}

// UploadPartCalls gets all the calls that were made to UploadPart.
// Check the length with:
//     len(mockedS3SDKClient.UploadPartCalls())
func (mock *S3SDKClientMock) UploadPartCalls() []struct {
	In *s3.UploadPartInput
} {
	var calls []struct {
		In *s3.UploadPartInput
	}
	lockS3SDKClientMockUploadPart.RLock()
	calls = mock.calls.UploadPart
	lockS3SDKClientMockUploadPart.RUnlock()
	return calls
}
