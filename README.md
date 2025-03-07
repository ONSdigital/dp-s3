dp-s3
================
Client to interact with AWS S3

## Getting started

### Setting up AWS credentials

In order to access AWS S3, this library will require your access key id and access secret key. You can either setup a default profile in ~/.aws/credentials file:

```sh
[default]
aws_access_key_id=<id>
aws_secret_access_key=<secret>
region=eu-west-1
```

Or export the values as environmental variables:

```sh
export AWS_ACCESS_KEY_ID=<id>
export AWS_SECRET_ACCESS_KEY=<secret>
```

More information in [Amazon documentation](https://docs.aws.amazon.com/cli/latest/userguide//cli-chap-configure.html)


### Setting up IAM policy

The functionality implemented by this library requires that the user has some permissions defined by an IAM policy.

- Health-check functionality performs a HEAD bucket operation, requiring allowed `s3:ListBucket` for all resources.

- Get functionality requires allowed `s3:GetObject` for the objects under the hierarchy you want to allow (e.g. `my-bucket/prefix/*`).

- Upload (PUT) functionality requires allowed `s3:PutObject` for the objects under the hierarchy you want to allow (e.g. `my-bucket/prefix/*`).

- Multipart upload functionality requires allowed `s3:PutObject`, `s3:GetObject`, `s3:AbortMultipartUpload`, `s3:ListMultipartUploadParts` for objects under the hierarchy you want to allow (e.g. `my-bucket/prefix/*`); and `s3:ListBucketMultipartUploads` for the bucket (e.g. `my-bucket`).

Please, see our [terraform repository](https://github.com/ONSdigital/dp-setup/tree/awsb/terraform) for more information.

### S3 Client Usage

The S3 client wraps the necessary AWS SDK structs and offers functionality to check buckets, and read and write objects from/to S3.

The client is configured with a specific bucket and region, note that the bucket needs to be created in the region that you provide in order to access it.

There are 3 available constructors:

- Constructor without AWS config (will create a new config):

```golang
import dps3 "github.com/ONSdigital/dp-s3/v3"

s3cli := dps3.NewClient(ctx, region, bucketName)
```

- Constructor with AWS config (will reuse the provided config):

```golang
import dps3 "github.com/ONSdigital/dp-s3/v3"

s3cli := dps3.NewClientWithConfig(bucketName, cfg, optFns ...func(*s3.Options))
```

- Constructor without AWS config but with credentials (will create a new config)

```golang
import dps3 "github.com/ONSdigital/dp-s3/v3"

s3cli := dps3.NewClientWithCredentials(ctx, region, bucketName, awsAccessKey, awsSecretKey)
```

It is recommended to create a single AWS config in your service and reuse it if you need other clients. The client offers a config getter: `s3cli.Config()`

A bucket name getter is also offered for convenience: `s3cli.BucketName()`

#### Get

The S3 client exposes functions to get S3 objects by using the vanilla SDK or the crypto client, for user-defined encryption keys.

Functions that have the suffix `WithPSK` allow you to provide a psk for encryption. For example:

- Get an un-encrypted object from S3

```golang
file, err := s3cli.Get("my/s3/file")
```

- Get an encrypted object from S3, using a psk:

```golang
file, err := s3cli.GetWithPSK("my/s3/file", psk)
```

You can get a file's metadata via a Head call:

```golang
out, err := s3cli.Head("my/s3/file")
```

#### Upload

The client also wraps the AWS SDK manager uploader, which is a high level client to upload files which automatically splits large files into chunks and uploads them concurrently.

This offers functionality to put objects in S3 in a single func call, hiding the low level details of chunking. More information [here](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/s3/manager)

Functions that have the suffix `WithPSK` allow you to provide a psk for encryption. For example:

- Upload an un-encrypted object to S3

```golang
result, err := s3cli.Upload(
    ctx,
    &s3.PutObjectInput{
        Body:   file.Reader,
        Key:    &filename,
    },
)
```

- Upload an encrypted object to S3, using a psk:

```golang
result, err := s3cli.UploadWithPSK(
    ctx,
    &s3.PutObjectInput{
        Body:   file.Reader,
        Key:    &filename,
    },
    psk,
)
```

#### Multipart Upload

You may use the low-level AWS SDK s3 client [multipart upload](./upload_multipart.go) methods

 and upload objects using `multipart upload`, which is an AWS SDK functionality to perform uploads in chunks. More information [here](https://docs.aws.amazon.com/AmazonS3/latest/userguide/mpuoverview.html)

##### Chunk Size

The minimum chunk size allowed in [AWS S3 is 5 MegaBytes (MB)](https://docs.aws.amazon.com/AmazonS3/latest/API/API_CompleteMultipartUpload.html)
if any chunks (excluding the final chunk) are under this size a ErrChunkTooSmall error will be returned from UploadPart
and UploadPartWithPsk functions when all chunks have been uploaded.

#### URL

S3Url is a structure intended to be used for S3 URL string manipulation in its different formats. To create a new structure you need to provide region, bucketName and object key,
and optionally the scheme:

```golang
s3Url, err := func NewURL(region, bucket, s3ObjectKey)
s3Url, err := func NewURLWithScheme(scheme, region, bucket, s3ObjectKey)
```

If you want to parse a URL into an s3Url object, you can use `ParseURL()` method, providing the format style:

```golang
s3Url, err := ParseURL(rawURL, URLStyle)
```

Once you have a valid s3Url object, you can obtain the URL string representation in the required format style by calling `String()` method:

```golang
str, err := s3Url.String(URLStyle)
```

##### Valid URL format Styles

The following URL styles are supported:

- PathStyle: `https://s3-eu-west-1.amazonaws.com/myBucket/my/s3/object/key`
- GlobalPathStyle: `https://s3.amazonaws.com/myBucket/my/s3/object/key`
- VirtualHostedStyle: `https://myBucket.s3-eu-west-1.amazonaws.com/my/s3/object/key`
- GlobalVirtualHostedStyle: `https://myBucket.s3.amazonaws.com/my/s3/object/key`
- AliasVirtualHostedStyle: '`https://myBucket/my/s3/object/key`

More information in [S3 official documentation](https://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html)

#### Health check

The S3 checker function performs a [HEAD bucket](https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#S3.HeadBucket) operation . The health check will succeed only if the bucket can be accessed using the client (i.e. client must be authenticated correctly, bucket must exist and have been created in the same region as the client).

Read the [Health Check Specification](https://github.com/ONSdigital/dp/blob/master/standards/HEALTH_CHECK_SPECIFICATION.md) for details.

After creating an S3 client as described above, call s3 health checker with `s3cli.Checker(context.Background())` and this will return a check object:

```golang
{
    "name": "string",
    "status": "string",
    "message": "string",
    "status_code": "int",
    "last_checked": "ISO8601 - UTC date time",
    "last_success": "ISO8601 - UTC date time",
    "last_failure": "ISO8601 - UTC date time"
}
```

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

## License

Copyright © 2020, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
