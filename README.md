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

You can access AWS S3 to get objects and do multipart uploads by creating a new client using the `NewClient()` function in client.go with the right region and bucketName,
or `NewClientWithSession()` if you already have an established AWS session.
Please, note that you will only be able to see S3 buckets created in a particular region using a client accessing that region.

```
s3cli := s3client.NewClient(<region>, <bucket>)
s3cli.Get(<S3ObjectKey>)
...
```

```
s3cli := NewClientWithSession(<bucket>, <awsSession>)
s3cli.Get(<S3ObjectKey>)
...
```

#### Uploader Usage

You can access AWS S3 to upload (PUT) objects by creating a new uploader using the `NewUploader()` function in uploader.go with the right region and bucketName,
or `NewUploaderWithSession()` if you already have an established AWS session.
Please, note that you will only be able to see S3 buckets created in a particular region using a client accessing that region.

```
s3Uploader := s3client.NewUploader(<region>, <bucket>)
s3Uploader.Upload(<input>)
...
```

```
s3Uploader := NewUploaderWithSession(<bucket>, <awsSession>)
s3Uploader.Upload(<input>)
...
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
