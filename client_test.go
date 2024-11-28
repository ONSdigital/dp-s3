package s3_test

import (
	"testing"

	dps3 "github.com/ONSdigital/dp-s3/v2"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewClient(t *testing.T) {
	Convey("Given an S3 bucket, region and env var AWS credentials", t, func() {

		bucket := "myBucket"
		region := "eu-north-1"
		awsAccessKey := "test"
		awsSecretKey := "test"

		t.Setenv("AWS_ACCESS_KEY_ID", awsAccessKey)
		t.Setenv("AWS_SECRET_ACCESS_KEY", awsSecretKey)

		Convey("When NewClient is called", func() {

			s3cli, err := dps3.NewClient(region, bucket)

			Convey("Then no error should be returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the expected client should be instantiated with the correct bucket, region, default endpoint and credentials", func() {
				So(s3cli, ShouldNotBeNil)
				So(s3cli.BucketName(), ShouldEqual, bucket)

				session := s3cli.Session()
				So(session.Config.Region, ShouldNotBeNil)
				So(*session.Config.Region, ShouldEqual, region)
				So(session.Config.Endpoint, ShouldBeNil) // a nil Endpoint means that the default AWS endpoints will be used

				creds, err := session.Config.Credentials.Get()
				So(err, ShouldBeNil)
				So(creds.AccessKeyID, ShouldEqual, awsAccessKey)
				So(creds.SecretAccessKey, ShouldEqual, awsSecretKey)
			})
		})
	})

	Convey("Given a valid S3 bucket, region and env var AWS credentials with an invalid AWS env var set", t, func() {

		bucket := "myBucket"
		region := "eu-north-1"
		awsAccessKey := "test"
		awsSecretKey := "test"

		t.Setenv("AWS_ACCESS_KEY_ID", awsAccessKey)
		t.Setenv("AWS_SECRET_ACCESS_KEY", awsSecretKey)
		t.Setenv("AWS_S3_USE_ARN_REGION", "invalid")

		Convey("When NewClient is called", func() {

			s3cli, err := dps3.NewClient(region, bucket)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("And no client should be returned", func() {
				So(s3cli, ShouldBeNil)
			})
		})
	})
}

func TestNewClientWithCredentials(t *testing.T) {
	Convey("Given an S3 bucket, region and credentials", t, func() {

		bucket := "myBucket"
		region := "eu-north-1"
		awsAccessKey := "test"
		awsSecretKey := "test"

		Convey("When NewClientWithCredentials is called", func() {

			s3cli, err := dps3.NewClientWithCredentials(region, bucket, awsAccessKey, awsSecretKey)

			Convey("Then no error should be returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the expected client should be instantiated with the correct bucket, region, default endpoint and credentials", func() {
				So(s3cli, ShouldNotBeNil)
				So(s3cli.BucketName(), ShouldEqual, bucket)

				session := s3cli.Session()
				So(session.Config.Region, ShouldNotBeNil)
				So(*session.Config.Region, ShouldEqual, region)
				So(session.Config.Endpoint, ShouldBeNil) // a nil Endpoint means that the default AWS endpoints will be used

				creds, err := session.Config.Credentials.Get()
				So(err, ShouldBeNil)
				So(creds.AccessKeyID, ShouldEqual, awsAccessKey)
				So(creds.SecretAccessKey, ShouldEqual, awsSecretKey)
			})
		})
	})

	Convey("Given an S3 bucket, region and empty credentials", t, func() {

		bucket := "myBucket"
		region := "eu-north-1"
		awsAccessKey := ""
		awsSecretKey := ""

		Convey("When NewClientWithCredentials is called", func() {

			s3cli, err := dps3.NewClientWithCredentials(region, bucket, awsAccessKey, awsSecretKey)

			Convey("Then no error should be returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the expected client should be instantiated with the correct bucket, region, default endpoint and credentials", func() {
				So(s3cli, ShouldNotBeNil)
				So(s3cli.BucketName(), ShouldEqual, bucket)

				session := s3cli.Session()
				So(session.Config.Region, ShouldNotBeNil)
				So(*session.Config.Region, ShouldEqual, region)
				So(session.Config.Endpoint, ShouldBeNil) // a nil endpoint means that the default AWS endpoints will be used

				creds, err := session.Config.Credentials.Get()
				So(err, ShouldNotBeNil)
				So(creds.AccessKeyID, ShouldEqual, awsAccessKey)
				So(creds.SecretAccessKey, ShouldEqual, awsSecretKey)
			})
		})
	})

	Convey("Given a valid S3 bucket, region and AWS credentials with an invalid AWS env var set", t, func() {

		bucket := "myBucket"
		region := "eu-north-1"
		awsAccessKey := "test"
		awsSecretKey := "test"

		t.Setenv("AWS_S3_USE_ARN_REGION", "invalid")

		Convey("When NewClientWithCredentials is called", func() {

			s3cli, err := dps3.NewClientWithCredentials(region, bucket, awsAccessKey, awsSecretKey)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("And no client should be returned", func() {
				So(s3cli, ShouldBeNil)
			})
		})
	})
}

func TestNewClientWithEndpoint(t *testing.T) {
	Convey("Given an S3 bucket, region, endpoint and env var AWS credentials", t, func() {

		bucket := "myBucket"
		region := "eu-north-1"
		endpoint := "http://some.endpoint.local"
		awsAccessKey := "test"
		awsSecretKey := "test"

		t.Setenv("AWS_ACCESS_KEY_ID", awsAccessKey)
		t.Setenv("AWS_SECRET_ACCESS_KEY", awsSecretKey)

		Convey("When NewClientWithEndpoint is called", func() {

			s3cli, err := dps3.NewClientWithEndpoint(region, bucket, endpoint)

			Convey("Then no error should be returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the expected client should be instantiated with the correct bucket, region, endpoint and credentials", func() {
				So(s3cli, ShouldNotBeNil)
				So(s3cli.BucketName(), ShouldEqual, bucket)

				session := s3cli.Session()
				So(session.Config.Region, ShouldNotBeNil)
				So(*session.Config.Region, ShouldEqual, region)
				So(session.Config.Endpoint, ShouldNotBeNil)
				So(*session.Config.Endpoint, ShouldEqual, endpoint)

				creds, err := session.Config.Credentials.Get()
				So(err, ShouldBeNil)
				So(creds.AccessKeyID, ShouldEqual, awsAccessKey)
				So(creds.SecretAccessKey, ShouldEqual, awsSecretKey)
			})
		})
	})

	Convey("Given a valid S3 bucket, region, endpoint and env var AWS credentials with an invalid AWS env var set", t, func() {

		bucket := "myBucket"
		region := "eu-north-1"
		endpoint := "http://some.endpoint.local"
		awsAccessKey := "test"
		awsSecretKey := "test"

		t.Setenv("AWS_ACCESS_KEY_ID", awsAccessKey)
		t.Setenv("AWS_SECRET_ACCESS_KEY", awsSecretKey)
		t.Setenv("AWS_S3_USE_ARN_REGION", "invalid")

		Convey("When NewClientWithEndpoint is called", func() {

			s3cli, err := dps3.NewClientWithEndpoint(region, bucket, endpoint)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("And no client should be returned", func() {
				So(s3cli, ShouldBeNil)
			})
		})
	})
}

func TestNewClientWithEndpointAndCredentials(t *testing.T) {
	Convey("Given an S3 bucket, region, endpoint and credentials", t, func() {

		bucket := "myBucket"
		region := "eu-north-1"
		endpoint := "http://some.endpoint.local"
		awsAccessKey := "test"
		awsSecretKey := "test"

		Convey("When NewClientWithEndpointAndCredentials is called", func() {

			s3cli, err := dps3.NewClientWithEndpointAndCredentials(region, bucket, endpoint, awsAccessKey, awsSecretKey)

			Convey("Then no error should be returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the expected client should be instantiated with the correct bucket, region, endpoint and credentials", func() {
				So(s3cli, ShouldNotBeNil)
				So(s3cli.BucketName(), ShouldEqual, bucket)

				session := s3cli.Session()
				So(session.Config.Region, ShouldNotBeNil)
				So(*session.Config.Region, ShouldEqual, region)
				So(session.Config.Endpoint, ShouldNotBeNil)
				So(*session.Config.Endpoint, ShouldEqual, endpoint)

				creds, err := session.Config.Credentials.Get()
				So(err, ShouldBeNil)
				So(creds.AccessKeyID, ShouldEqual, awsAccessKey)
				So(creds.SecretAccessKey, ShouldEqual, awsSecretKey)
			})
		})
	})

	Convey("Given an S3 bucket, region, endpoint and empty credentials", t, func() {

		bucket := "myBucket"
		region := "eu-north-1"
		endpoint := "http://some.endpoint.local"
		awsAccessKey := ""
		awsSecretKey := ""

		Convey("When NewClientWithEndpointAndCredentials is called", func() {

			s3cli, err := dps3.NewClientWithEndpointAndCredentials(region, bucket, endpoint, awsAccessKey, awsSecretKey)

			Convey("Then no error should be returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("And the expected client should be instantiated with the correct bucket, region, endpoint and credentials", func() {
				So(s3cli, ShouldNotBeNil)
				So(s3cli.BucketName(), ShouldEqual, bucket)

				session := s3cli.Session()
				So(session.Config.Region, ShouldNotBeNil)
				So(*session.Config.Region, ShouldEqual, region)
				So(session.Config.Endpoint, ShouldNotBeNil)
				So(*session.Config.Endpoint, ShouldEqual, endpoint)

				creds, err := session.Config.Credentials.Get()
				So(err, ShouldNotBeNil)
				So(creds.AccessKeyID, ShouldEqual, awsAccessKey)
				So(creds.SecretAccessKey, ShouldEqual, awsSecretKey)
			})
		})
	})

	Convey("Given a valid S3 bucket, region, endpoint and AWS credentials with an invalid AWS env var set", t, func() {

		bucket := "myBucket"
		region := "eu-north-1"
		endpoint := "http://some.endpoint.local"
		awsAccessKey := "test"
		awsSecretKey := "test"

		t.Setenv("AWS_S3_USE_ARN_REGION", "invalid")

		Convey("When NewClientWithEndpointAndCredentials is called", func() {

			s3cli, err := dps3.NewClientWithEndpointAndCredentials(region, bucket, endpoint, awsAccessKey, awsSecretKey)

			Convey("Then an error should be returned", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("And no client should be returned", func() {
				So(s3cli, ShouldBeNil)
			})
		})
	})
}
