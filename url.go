// file: url.go
//
// Contains string manipulation methods to obtain an S3 URL in the different styles supported by AWS
// and translate from one to another.
package s3

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// URLStyle is the type to define the URL style iota enumeration corresponding an S3 url (path, virtualHosted, etc)
type URLStyle int

// Possible S3 URL format styles, as defined in https://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html
const (
	// PathStyle example: 'https://s3-eu-west-1.amazonaws.com/myBucket/my/s3/object/key'
	PathStyle = iota
	// GlobalPathStyle example: 'https://s3.amazonaws.com/myBucket/my/s3/object/key'
	GlobalPathStyle
	// VirtualHostedStyle example: 'https://myBucket.s3-eu-west-1.amazonaws.com/my/s3/object/key'
	VirtualHostedStyle
	// GlobalVirtualHostedStyle example: 'https://myBucket.s3.amazonaws.com/my/s3/object/key'
	GlobalVirtualHostedStyle
	// AliasVirtualHostedStyle example: 'https://myBucket/my/s3/object/key'
	AliasVirtualHostedStyle
)

var urlStyles = []string{"Path", "GlobalPath", "VirtualHosted", "GlobalVirtualHosted", "AliasVirtualHosted"}

// Values of the format styles
func (style URLStyle) String() string {
	return urlStyles[style]
}

// S3Url represents an S3 URL with bucketName, key and region (optional). This struct is
// intended to be used for S3 URL string manipulation/translation in its possible format styles.
type S3Url struct {
	Scheme     string
	Region     string
	BucketName string
	Key        string
}

// String returns the S3 URL string in the requested format style.
func (s3Url *S3Url) String(style URLStyle) (string, error) {
	switch style {
	case PathStyle:
		if len(s3Url.Region) == 0 {
			return "", errors.New("path style format requires a region")
		}
		urlFormat := "%s://s3-%s.amazonaws.com/%s/%s"
		return fmt.Sprintf(urlFormat, s3Url.Scheme, s3Url.Region, s3Url.BucketName, s3Url.Key), nil
	case GlobalPathStyle:
		urlFormat := "%s://s3.amazonaws.com/%s/%s"
		return fmt.Sprintf(urlFormat, s3Url.Scheme, s3Url.BucketName, s3Url.Key), nil
	case VirtualHostedStyle:
		if len(s3Url.Region) == 0 {
			return "", errors.New("virtual-hosted style format requires a region")
		}
		urlFormat := "%s://%s.s3-%s.amazonaws.com/%s"
		return fmt.Sprintf(urlFormat, s3Url.Scheme, s3Url.BucketName, s3Url.Region, s3Url.Key), nil
	case GlobalVirtualHostedStyle:
		urlFormat := "%s://%s.s3.amazonaws.com/%s"
		return fmt.Sprintf(urlFormat, s3Url.Scheme, s3Url.BucketName, s3Url.Key), nil
	case AliasVirtualHostedStyle:
		urlFormat := "%s://%s/%s"
		return fmt.Sprintf(urlFormat, s3Url.Scheme, s3Url.BucketName, s3Url.Key), nil
	}
	return "", errors.New("undefined style")
}

// ParseURL creates an S3Url struct from the provided rawULR and format style
func ParseURL(rawURL string, style URLStyle) (*S3Url, error) {
	switch style {
	case PathStyle:
		return ParsePathStyleURL(rawURL)
	case GlobalPathStyle:
		return ParseGlobalPathStyleURL(rawURL)
	case VirtualHostedStyle:
		return ParseVirtualHostedURL(rawURL)
	case GlobalVirtualHostedStyle:
		return ParseGlobalVirtualHostedURL(rawURL)
	case AliasVirtualHostedStyle:
		return ParseAliasVirtualHostedURL(rawURL)
	}
	return nil, errors.New("undefined style")
}

// ParsePathStyleURL creates an S3Url struct from the provided path-style url string
// Example: 'https://s3-eu-west-1.amazonaws.com/myBucket/my/s3/object/key'.
func ParsePathStyleURL(pathStyleURL string) (*S3Url, error) {

	parsedUrl, err := url.Parse(pathStyleURL)
	if err != nil {
		return nil, err
	}

	region := strings.TrimSuffix(strings.TrimPrefix(parsedUrl.Host, "s3-"), ".amazonaws.com")
	if len(region) == 0 {
		return nil, fmt.Errorf("wrong region in path-style url: %s", pathStyleURL)
	}

	bucketName, key, err := parsePath(parsedUrl)
	if err != nil {
		return nil, err
	}

	return NewURLWithScheme(parsedUrl.Scheme, region, bucketName, key)
}

// ParseGlobalPathStyleURL creates an S3Url struct from the provided global-path-style url string
// Example: 'https://s3.amazonaws.com/myBucket/my/s3/object/key'
// This method is compatible with PathStyle format (if region is present in the URL, it will be ignored)
func ParseGlobalPathStyleURL(gpURL string) (*S3Url, error) {
	parsedUrl, err := url.Parse(gpURL)
	if err != nil {
		return nil, err
	}

	bucketName, key, err := parsePath(parsedUrl)
	if err != nil {
		return nil, err
	}

	return NewURLWithScheme(parsedUrl.Scheme, "", bucketName, key)
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
	parsedUrl, err := url.Parse(vhURL)
	if err != nil {
		return nil, err
	}

	splittedHost := strings.Split(parsedUrl.Host, ".")
	if len(splittedHost) < 4 {
		return nil, fmt.Errorf("could not find bucket name or region in virtual-hosted-style url %s", vhURL)
	}

	region := strings.TrimPrefix(splittedHost[len(splittedHost)-3], "s3-")
	if len(region) == 0 {
		return nil, fmt.Errorf("wrong region in virtual-hosted-style url: %s", vhURL)
	}

	bucketName := strings.TrimSuffix(parsedUrl.Host, fmt.Sprintf(".s3-%s.amazonaws.com", region))
	if len(bucketName) == 0 {
		return nil, fmt.Errorf("wrong bucket name in virtual-hosted-style url: %s", vhURL)
	}

	key := strings.TrimPrefix(parsedUrl.Path, "/")
	if len(key) == 0 {
		return nil, fmt.Errorf("wrong key in virtual-hosted-style url: %s", vhURL)
	}

	return NewURLWithScheme(parsedUrl.Scheme, region, bucketName, key)
}

// ParseGlobalVirtualHostedURL creates an S3Url struct from the provided global-virtual-hosted-style url string
// Example: 'https://myBucket.s3.amazonaws.com/my/s3/object/key'
func ParseGlobalVirtualHostedURL(gvhURL string) (*S3Url, error) {
	parsedUrl, err := url.Parse(gvhURL)
	if err != nil {
		return nil, err
	}

	bucketName := strings.TrimSuffix(parsedUrl.Host, ".s3.amazonaws.com")
	if len(bucketName) == 0 {
		return nil, fmt.Errorf("wrong bucketName in global virtual hosted style url: %s", gvhURL)
	}

	key := strings.TrimPrefix(parsedUrl.Path, "/")
	if len(key) == 0 {
		return nil, fmt.Errorf("wrong key in global virtual hosted style url: %s", gvhURL)
	}
	return NewURLWithScheme(parsedUrl.Scheme, "", bucketName, key)
}

// ParseAliasVirtualHostedURL creates an S3Url struct from the provided dns-alias-virtual-hosted-style url string
// Example: 'https://myBucket/my/s3/object/key'
func ParseAliasVirtualHostedURL(avhURL string) (*S3Url, error) {
	parsedUrl, err := url.Parse(avhURL)
	if err != nil {
		return nil, err
	}

	bucketName := parsedUrl.Host
	if len(bucketName) == 0 {
		return nil, fmt.Errorf("wrong bucketName in DNS-alias-virtual-hosted-style url: %s", avhURL)
	}

	key := strings.TrimPrefix(parsedUrl.Path, "/")
	if len(key) == 0 {
		return nil, fmt.Errorf("wrong key in global virtual hosted style url: %s", avhURL)
	}

	return NewURLWithScheme(parsedUrl.Scheme, "", bucketName, key)
}

// NewURL instantiates a new S3Url struct with the provided region, bucket name and object key
func NewURL(region, bucketName, key string) (*S3Url, error) {
	return NewURLWithScheme("https", region, bucketName, key)
}

// NewURLWithScheme instantiates a new S3Url struct with the provided scheme, region, bucket and object key
func NewURLWithScheme(scheme, region, bucketName, key string) (*S3Url, error) {
	if len(bucketName) == 0 {
		return nil, errors.New("bucketName required")
	}
	if len(key) == 0 {
		return nil, errors.New("key required")
	}
	if len(scheme) == 0 {
		scheme = "https"
	}
	return &S3Url{
		Scheme:     scheme,
		Region:     region,
		BucketName: bucketName,
		Key:        key,
	}, nil
}
