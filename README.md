dp-s3
================
Client to interact with AWS S3

### Getting started

#### Setting up AWS credentials

In order to access AWS S3, this library will require your access key id and access secret key. You can either setup a default profile in ~/.aws/credentials file:
```
[default]
aws_access_key_id=<id>
aws_secret_access_key=<secret>
region=eu-west-1
```

Or export the values as environmental variables:
```
export AWS_ACCESS_KEY_ID=<id>
export AWS_SECRET_ACCESS_KEY=<secret>
```

More information in [Amazon documentation](https://docs.aws.amazon.com/cli/latest/userguide//cli-chap-configure.html)


#### Setting up IAM policy

The functionality implemented by this library requires that the user has some permissions defined by an IAM policy.

- Health-check functionality performs a HEAD bucket operation, requiring allowed `s3:ListBucket` for all resources.

- Get functionality requires allowed `s3:GetObject` for the objects under the hierarchy you want to allow (e.g. `my-bucket/prefix/*`).

- Upload (PUT) functionality requires allowed `s3:PutObject` for the objects under the hierarchy you want to allow (e.g. `my-bucket/prefix/*`).

- Multipart upload functionality requires allowed `s3:PutObject`, `s3:GetObject`, `s3:AbortMultipartUpload`, `s3:ListMultipartUploadParts` for objects under the hierarchy you want to allow (e.g. `my-bucket/prefix/*`); and `s3:ListBucketMultipartUploads` for the bucket (e.g. `my-bucket`).

Please, see our [terraform repository](https://github.com/ONSdigital/dp-setup/tree/develop/terraform) for more information.

#### S3 Client Usage

The S3 client wraps the AWS SDK s3 client and offers functionality to read objects from S3 and upload objects using `multipart upload`, which is an AWS SDK functionality to perform uploads in chunks. More information [here](https://docs.aws.amazon.com/AmazonS3/latest/userguide/mpuoverview.html)

The client contains a bucket and region, note that the bucket needs to be created in the region that you provide in order to access it.

There are 2 available constructors:
- Constructor without AWS session (will create a new session):
```
s3cli := s3client.NewClient(<region>, <bucket>)
```
- Constructor with AWS session (will reuse the provided session):
```
s3cli := s3client.NewClientWithSession(<bucket>, <awsSession>)
```
It is recommended to create a single AWS session in your service and reuse it if you need other clients. The client offers a session getter: `s3cli.Session()`

The S3 client exposes functions to get or upload files using the vanilla aws sdk, or the s3crypto wrapper, which allows you to provide a psk (pre-shared key) for encryption.

Functions that have the suffix `WithPSK` allow you to provide a psk for encryption. For example:
- Get an un-encrypted object from S3
```
file, err := s3cli.Get("my/s3/file")
```
- Get an encrypted object from S3, using a psk:
```
file, err := s3cli.GetWithPSK("my/s3/file", psk)
```


#### Uploader Usage

The Uploader is a higher level S3 client that wraps the SDK uploader, from s3manager package, as well as the lower level S3 client.
This offers functionality to put objects in S3 in a single func call, hiding the low level details of chunking. More information [here](https://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#Uploader)

Similarly to the s3 client, you can create an uploader and establish a new session, or reuse an existing one:

- Constructor without AWS session (will create a new session):
```
s3Uploader := s3client.NewUploader(<region>, <bucket>)
```
- Constructor with AWS session (will reuse the provided session):
```
s3Uploader := s3client.NewUploaderWithSession(<bucket>, <awsSession>)
```

Similarly to the s3 client, it is recommended to reuse AWS sessions between clients/uploaderes.

Functions that have the suffix `WithPSK` allow you to provide a psk for encryption. For example:
- Upload an un-encrypted object to S3
```
result, err := s3Uploader.Upload(&s3manager.UploadInput{
		Body:   file.Reader,
		Key:    &filename,
	})
```
- Upload an encrypted object to S3, using a psk:
```
result, err := s3Uploader.UploadWithPSK(&s3manager.UploadInput{
		Body:   file.Reader,
		Key:    &filename,
	}, psk)
```

#### URL Usage

S3Url is a structure intended to be used for S3 URL string manipulation in its different formats. To create a new structure you need to provide region, bucketName and object key,
and optionally the scheme:

```
s3Url, err := func NewURL(<region>, <bucket>, <s3ObjcetKey>)
s3Url, err := func NewURLWithScheme(<scheme>, <region>, <bucket>, <s3ObjcetKey>)
```

If you want to parse a URL into an s3Url object, you can use `ParseURL()` method, providing the format style:

```
s3Url, err := ParseURL(<rawURL>, <URLStyle>)
```

Once you have a valid s3Url object, you can obtain the URL string representation in the required format style by calling `String()` method:

```
str, err := s3Url.String(<URLStyle>)
```

##### Valid URL format Styles

The following URL styles are supported:

- PathStyle: `https://s3-eu-west-1.amazonaws.com/myBucket/my/s3/object/key`
- GlobalPathStyle: `https://s3.amazonaws.com/myBucket/my/s3/object/key`
- VirtualHostedStyle: `https://myBucket.s3-eu-west-1.amazonaws.com/my/s3/object/key`
- GlobalVirtualHostedStyle: `https://myBucket.s3.amazonaws.com/my/s3/object/key`
- AliasVirtualHostedStyle: '`https://myBucket/my/s3/object/key`

More information in [S3 official documentation](https://docs.aws.amazon.com/AmazonS3/latest/dev/VirtualHosting.html)

### Health package

The S3 checker function performs a [HEAD bucket](https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#S3.HeadBucket) operation . The health check will succeed only if the bucket can be accessed using the client (i.e. client must be authenticated correctly, bucket must exist and have been created in the same region as the client).

Read the [Health Check Specification](https://github.com/ONSdigital/dp/blob/master/standards/HEALTH_CHECK_SPECIFICATION.md) for details.

After creating an S3 client as described above, call s3 health checker with `s3cli.Checker(context.Background())` and this will return a check object:

```
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

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2020, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
