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

- Health-check functionality performs a HEAD bucket operation, requiring allowed `s3:HeadBucket` for all resources.

- Get functionality requires allowed `s3:GetObject` for the objects under the hierarchy you want to allow (e.g. `my-bucket/prefix/*`).

- Upload (PUT) functionality requires allowed `s3:PutObject` for the objects under the hierarchy you want to allow (e.g. `my-bucket/prefix/*`).

- Multipart upload functionality requires allowed `s3:PutObject`, `s3:GetObject`, `s3:AbortMultipartUpload`, `s3:ListMultipartUploadParts` for objects under the hierarcy you want to allow (e.g. `my-bucket/prefix/*`); and `s3:ListBucketMultipartUploads` for the bucket (e.g. `my-bucket`).

Please, see our [terraform repository](https://github.com/ONSdigital/dp-setup/tree/develop/terraform) for more information.

#### Usage

You can access AWS S3 creating a new client using the New() function in client.go providing the right region. Please, note that you will only be able to see S3 buckets created in a particular region using a client accessing that region.

```
s3cli := s3client.New(<region>)
s3cli.Get(<url>)
```

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
