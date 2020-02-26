package s3client

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// URLStyle is the type to define the URL style iota enumeration corresponding an S3 url (path, virtualHosted, etc)
type URLStyle int

// Possible S3 URL styles
const (
	StylePath = iota
	StyleGlobalPath
	StyleVirtualHosted
	StyleGlobalVirtualHosted
	StyleAliasVirtualHosted
)

var urlStyles = []string{"Path", "GlobalPath", "VirtualHosted", "GlobalVirtualHosted", "AliasVirtualHosted"}

// Values of the format styles
func (style URLStyle) String() string {
	return urlStyles[style]
}

// S3Url represents an S3 URL with bucketName, key and region (optional). This struct is
// intended to be used for S3 URL string manipulation/translation in its possible format styles.
type S3Url struct {
	region     string
	bucketName string
	key        string
}

// Region returns the region defined by the url
func (s3Url *S3Url) Region() string {
	return s3Url.region
}

// BucketName returns the bucket name defined by the url
func (s3Url *S3Url) BucketName() string {
	return s3Url.bucketName
}

// Key returns the object key defined by the url
func (s3Url *S3Url) Key() string {
	return s3Url.key
}

// String returns the S3 URL string in the requested format style.
// Possible formats are defined in https://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html
// PathStyle example: 'https://s3-eu-west-1.amazonaws.com/myBucket/my/s3/object/key'
// GlobalPathStyle example: 'https://s3.amazonaws.com/myBucket/my/s3/object/key'
// VirtualHostedStyle example: 'https://myBucket.s3-eu-west-1.amazonaws.com/my/s3/object/key'
// GlobalVirtualHostedStyle example: 'https://myBucket.s3.amazonaws.com/my/s3/object/key'
// AliasVirtualHostedStyle example: 'https://myBucket/my/s3/object/key'
func (s3Url *S3Url) String(style URLStyle) (string, error) {
	switch style {
	case StylePath:
		if len(s3Url.region) == 0 {
			return "", errors.New("Path style format requires a region")
		}
		url := "https://s3-%s.amazonaws.com/%s/%s"
		return fmt.Sprintf(url, s3Url.region, s3Url.bucketName, s3Url.key), nil
	case StyleGlobalPath:
		url := "https://s3.amazonaws.com/%s/%s"
		return fmt.Sprintf(url, s3Url.bucketName, s3Url.key), nil
	case StyleVirtualHosted:
		if len(s3Url.region) == 0 {
			return "", errors.New("Virtual-hosted style format requires a region")
		}
		url := "https://%s.s3-%s.amazonaws.com/%s"
		return fmt.Sprintf(url, s3Url.bucketName, s3Url.region, s3Url.key), nil
	case StyleGlobalVirtualHosted:
		url := "https://%s.s3.amazonaws.com/%s"
		return fmt.Sprintf(url, s3Url.bucketName, s3Url.key), nil
	case StyleAliasVirtualHosted:
		url := "https://%s/%s"
		return fmt.Sprintf(url, s3Url.bucketName, s3Url.key), nil
	}
	return "", errors.New("Undefined style")
}

// ParseURL creates an S3Url struct from the provided rawULR and format style
func ParseURL(rawURL string, style URLStyle) (*S3Url, error) {
	switch style {
	case StylePath:
		return ParsePathStyleURL(rawURL)
	case StyleGlobalPath:
		return ParseGlobalPathStyleURL(rawURL)
	case StyleVirtualHosted:
		return ParseVirtualHostedURL(rawURL)
	case StyleGlobalVirtualHosted:
		return ParseGlobalVirtualHostedURL(rawURL)
	case StyleAliasVirtualHosted:
		return ParseAliasVirtualHostedURL(rawURL)
	}
	return nil, errors.New("Undefined style")
}

// ParsePathStyleURL creates an S3Url struct from the provided path-style url string
// Example: 'https://s3-eu-west-1.amazonaws.com/myBucket/my/s3/object/key'
func ParsePathStyleURL(pathStyleURL string) (*S3Url, error) {

	url, err := url.Parse(pathStyleURL)
	if err != nil {
		return nil, err
	}

	region := strings.TrimSuffix(strings.TrimPrefix(url.Host, "s3-"), ".amazonaws.com")
	if len(region) == 0 {
		return nil, fmt.Errorf("wrong region in path-style url: %s", pathStyleURL)
	}

	bucketName, key, err := parsePath(url)
	if err != nil {
		return nil, err
	}

	return NewURL(region, bucketName, key)
}

// ParseGlobalPathStyleURL creates an S3Url struct from the provided global-path-style url string
// Example: 'https://s3.amazonaws.com/myBucket/my/s3/object/key'
func ParseGlobalPathStyleURL(gpURL string) (*S3Url, error) {
	url, err := url.Parse(gpURL)
	if err != nil {
		return nil, err
	}

	bucketName, key, err := parsePath(url)
	if err != nil {
		return nil, err
	}

	return NewURL("", bucketName, key)
}

func parsePath(url *url.URL) (bucketName string, key string, err error) {
	splittedPath := strings.Split(url.Path, "/")
	if len(splittedPath) < 3 {
		return "", "", fmt.Errorf("could not find bucket or filename in file path-style url %s", url.String())
	}

	bucketName = splittedPath[1]
	if len(bucketName) == 0 {
		return "", "", fmt.Errorf("missing bucket name in path-style url %s", url.String())
	}

	key = strings.TrimPrefix(url.Path, fmt.Sprintf("/%s/", bucketName))
	if len(key) == 0 {
		return "", "", fmt.Errorf("missing s3 object key in path-style url %s", url.String())
	}
	return
}

// ParseVirtualHostedURL creates an S3Url struct from the provided virtual-hosted-style url string
// Example: 'https://myBucket.s3-eu-west-1.amazonaws.com/my/s3/object/key'
func ParseVirtualHostedURL(vhURL string) (*S3Url, error) {
	url, err := url.Parse(vhURL)
	if err != nil {
		return nil, err
	}

	splittedHost := strings.Split(url.Host, ".")
	if len(splittedHost) < 4 {
		return nil, fmt.Errorf("could not find bucket name or region in virtual-hosted-style url %s", vhURL)
	}

	region := strings.TrimPrefix(splittedHost[len(splittedHost)-3], "s3-")
	if len(region) == 0 {
		return nil, fmt.Errorf("wrong region in virtual-hosted-style url: %s", vhURL)
	}

	bucketName := strings.TrimSuffix(url.Host, fmt.Sprintf(".s3-%s.amazonaws.com", region))
	if len(bucketName) == 0 {
		return nil, fmt.Errorf("wrong bucket name in virtual-hosted-style url: %s", vhURL)
	}

	key := strings.TrimPrefix(url.Path, "/")
	if len(key) == 0 {
		return nil, fmt.Errorf("wrong key in virtual-hosted-style url: %s", vhURL)
	}

	return NewURL(region, bucketName, key)
}

// ParseGlobalVirtualHostedURL creates an S3Url struct from the provided global-virtual-hosted-style url string
// Example: 'https://myBucket.s3.amazonaws.com/my/s3/object/key'
func ParseGlobalVirtualHostedURL(gvhURL string) (*S3Url, error) {
	url, err := url.Parse(gvhURL)
	if err != nil {
		return nil, err
	}

	bucketName := strings.TrimSuffix(url.Host, ".s3.amazonaws.com")
	if len(bucketName) == 0 {
		return nil, fmt.Errorf("wrong bucketName in global virtual hosted style url: %s", gvhURL)
	}

	key := strings.TrimPrefix(url.Path, "/")
	if len(key) == 0 {
		return nil, fmt.Errorf("wrong key in global virtual hosted style url: %s", gvhURL)
	}

	return NewURL("", bucketName, key)
}

// ParseAliasVirtualHostedURL creates an S3Url struct from the provided dns-alias-virtual-hosted-style url string
// Example: 'https://myBucket/my/s3/object/key'
func ParseAliasVirtualHostedURL(avhURL string) (*S3Url, error) {
	url, err := url.Parse(avhURL)
	if err != nil {
		return nil, err
	}

	bucketName := url.Host
	if len(bucketName) == 0 {
		return nil, fmt.Errorf("wrong bucketName in DNS-alias-virtual-hosted-style url: %s", avhURL)
	}

	key := strings.TrimPrefix(url.Path, "/")
	if len(key) == 0 {
		return nil, fmt.Errorf("wrong key in global virtual hosted style url: %s", avhURL)
	}

	return NewURL("", bucketName, key)
}

// NewURL instantiates a new S3Url struct with the provided region, bucket name and object key
func NewURL(region, bucketName, key string) (*S3Url, error) {
	if len(bucketName) == 0 {
		return nil, errors.New("bucketName required")
	}
	if len(key) == 0 {
		return nil, errors.New("key required")
	}
	return &S3Url{
		region:     region,
		bucketName: bucketName,
		key:        key,
	}, nil
}
