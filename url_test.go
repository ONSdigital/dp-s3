package s3client_test

import (
	"fmt"
	"testing"

	s3client "github.com/ONSdigital/dp-s3"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFullyDefinedUrl(t *testing.T) {

	const expectedBucketName = "csv-bucket"
	const expectedKey = "dir1/test-file.csv"
	const expectedRegion = "eu-west-1"

	Convey("Given an instance of S3Url with valid bucketName, region and object key", t, func() {
		s3Url, err := s3client.NewURL(expectedRegion, expectedBucketName, expectedKey)
		So(err, ShouldBeNil)

		Convey("Getters return the expected values", func() {
			So(s3Url.Region, ShouldEqual, expectedRegion)
			So(s3Url.BucketName, ShouldEqual, expectedBucketName)
			So(s3Url.Key, ShouldEqual, expectedKey)
		})

		Convey("Path style URL string is formatted as expected", func() {
			expectedStr := fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s", expectedRegion, expectedBucketName, expectedKey)
			urlStr, err := s3Url.String(s3client.PathStyle)
			So(err, ShouldBeNil)
			So(urlStr, ShouldEqual, expectedStr)
		})

		Convey("Global Path style URL string is formatted as expected", func() {
			expectedStr := fmt.Sprintf("https://s3.amazonaws.com/%s/%s", expectedBucketName, expectedKey)
			urlStr, err := s3Url.String(s3client.GlobalPathStyle)
			So(err, ShouldBeNil)
			So(urlStr, ShouldEqual, expectedStr)
		})

		Convey("Virtual hosted style URL string is formatted as expected", func() {
			expectedStr := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", expectedBucketName, expectedRegion, expectedKey)
			urlStr, err := s3Url.String(s3client.VirtualHostedStyle)
			So(err, ShouldBeNil)
			So(urlStr, ShouldEqual, expectedStr)
		})

		Convey("Global virtual hosted style URL string is formatted as expected", func() {
			expectedStr := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", expectedBucketName, expectedKey)
			urlStr, err := s3Url.String(s3client.GlobalVirtualHostedStyle)
			So(err, ShouldBeNil)
			So(urlStr, ShouldEqual, expectedStr)
		})

		Convey("DNS Alias virtual hosted style URL string is formatted as expected", func() {
			expectedStr := fmt.Sprintf("https://%s/%s", expectedBucketName, expectedKey)
			urlStr, err := s3Url.String(s3client.AliasVirtualHostedStyle)
			So(err, ShouldBeNil)
			So(urlStr, ShouldEqual, expectedStr)
		})
	})
}

func TestNoRegionUrl(t *testing.T) {

	const expectedBucketName = "csv-bucket"
	const expectedKey = "dir1/test-file.csv"

	Convey("Given an instance of S3Url with valid bucketName, region and object key", t, func() {
		s3Url, err := s3client.NewURL("", expectedBucketName, expectedKey)
		So(err, ShouldBeNil)

		Convey("Getters return the expected values", func() {
			So(s3Url.Region, ShouldEqual, "")
			So(s3Url.BucketName, ShouldEqual, expectedBucketName)
			So(s3Url.Key, ShouldEqual, expectedKey)
		})

		Convey("Path style URL string is formatted as expected", func() {
			_, err := s3Url.String(s3client.PathStyle)
			So(err, ShouldNotBeNil)
		})

		Convey("Global Path style URL string is formatted as expected", func() {
			expectedStr := fmt.Sprintf("https://s3.amazonaws.com/%s/%s", expectedBucketName, expectedKey)
			urlStr, err := s3Url.String(s3client.GlobalPathStyle)
			So(err, ShouldBeNil)
			So(urlStr, ShouldEqual, expectedStr)
		})

		Convey("Virtual hosted style URL string is formatted as expected", func() {
			_, err := s3Url.String(s3client.VirtualHostedStyle)
			So(err, ShouldNotBeNil)
		})

		Convey("Global virtual hosted style URL string is formatted as expected", func() {
			expectedStr := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", expectedBucketName, expectedKey)
			urlStr, err := s3Url.String(s3client.GlobalVirtualHostedStyle)
			So(err, ShouldBeNil)
			So(urlStr, ShouldEqual, expectedStr)
		})

		Convey("DNS Alias virtual hosted style URL string is formatted as expected", func() {
			expectedStr := fmt.Sprintf("https://%s/%s", expectedBucketName, expectedKey)
			urlStr, err := s3Url.String(s3client.AliasVirtualHostedStyle)
			So(err, ShouldBeNil)
			So(urlStr, ShouldEqual, expectedStr)
		})
	})
}

func TestParsing(t *testing.T) {

	const expectedBucketName = "csv-bucket"
	const expectedKey = "dir1/test-file.csv"
	const expectedRegion = "eu-west-1"

	Convey("Given S3 raw url strings in different acceptable formats", t, func() {

		expectedRegionalHttpsUrl, err := s3client.NewURL(expectedRegion, expectedBucketName, expectedKey)
		So(err, ShouldBeNil)
		expectedRegionalS3Url, err := s3client.NewURLWithScheme("s3", expectedRegion, expectedBucketName, expectedKey)
		So(err, ShouldBeNil)
		expectedGlobalHttpsUrl, err := s3client.NewURL("", expectedBucketName, expectedKey)
		So(err, ShouldBeNil)
		expectedGlobalS3Url, err := s3client.NewURLWithScheme("s3", "", expectedBucketName, expectedKey)
		So(err, ShouldBeNil)

		// URLs by style and expected generated s3Url objects
		urls := map[s3client.URLStyle]map[string]*s3client.S3Url{
			s3client.PathStyle: map[string]*s3client.S3Url{
				"https://s3-eu-west-1.amazonaws.com/csv-bucket/dir1/test-file.csv": expectedRegionalHttpsUrl,
				"s3://s3-eu-west-1.amazonaws.com/csv-bucket/dir1/test-file.csv":    expectedRegionalS3Url,
			},
			s3client.VirtualHostedStyle: map[string]*s3client.S3Url{
				"https://csv-bucket.s3-eu-west-1.amazonaws.com/dir1/test-file.csv": expectedRegionalHttpsUrl,
				"s3://csv-bucket.s3-eu-west-1.amazonaws.com/dir1/test-file.csv":    expectedRegionalS3Url,
			},
			s3client.GlobalPathStyle: map[string]*s3client.S3Url{
				"https://s3.amazonaws.com/csv-bucket/dir1/test-file.csv": expectedGlobalHttpsUrl,
				"s3://s3.amazonaws.com/csv-bucket/dir1/test-file.csv":    expectedGlobalS3Url,
			},
			s3client.GlobalVirtualHostedStyle: map[string]*s3client.S3Url{
				"https://csv-bucket.s3.amazonaws.com/dir1/test-file.csv": expectedGlobalHttpsUrl,
				"s3://csv-bucket.s3.amazonaws.com/dir1/test-file.csv":    expectedGlobalS3Url,
			},
			s3client.AliasVirtualHostedStyle: map[string]*s3client.S3Url{
				"https://csv-bucket/dir1/test-file.csv": expectedGlobalHttpsUrl,
				"s3://csv-bucket/dir1/test-file.csv":    expectedGlobalS3Url,
			},
		}

		Convey("Each format is correctly parsed, successfully retreiving bucket, key and region (if available)", func() {
			for style, urlMap := range urls {
				for url, expectedObject := range urlMap {
					s3Url, err := s3client.ParseURL(url, style)
					So(err, ShouldBeNil)
					So(s3Url, ShouldResemble, expectedObject)
				}
			}
		})

		Convey("A path-style url can be parsed as a global-path-style with empty region", func() {
			s3Url, err := s3client.ParseURL(
				"https://s3-eu-west-1.amazonaws.com/csv-bucket/dir1/test-file.csv", s3client.GlobalPathStyle)
			So(err, ShouldBeNil)
			So(s3Url, ShouldResemble, expectedGlobalHttpsUrl)
			s3Url, err = s3client.ParseURL(
				"s3://s3-eu-west-1.amazonaws.com/csv-bucket/dir1/test-file.csv", s3client.GlobalPathStyle)
			So(err, ShouldBeNil)
			So(s3Url, ShouldResemble, expectedGlobalS3Url)
		})

		Convey("Trying to parse an empty S3 raw url results in error ", func() {
			for style, _ := range urls {
				_, err := s3client.ParseURL("", style)
				So(err, ShouldNotBeNil)
			}
		})

		Convey("Tying to parse an s3 url that is missing the object key, results in error", func() {
			missingBucketUrl := "s3://some-file"
			for style, _ := range urls {
				_, err := s3client.ParseURL(missingBucketUrl, style)
				So(err, ShouldNotBeNil)
			}
		})

		Convey("Trying to parse an s3 url with empty bucket or key results in error", func() {
			emptyValuesUrl1 := "s3://///////"
			emptyValuesUrl2 := fmt.Sprintf("s3:/%s/", expectedBucketName)
			for style, _ := range urls {
				_, err := s3client.ParseURL(emptyValuesUrl1, style)
				So(err, ShouldNotBeNil)
				_, err = s3client.ParseURL(emptyValuesUrl2, style)
				So(err, ShouldNotBeNil)
			}
		})

		Convey("Tying to parse an s3 url that is missing the bucket name and object key results in error", func() {
			missingBucketUrl := "s3://"
			for style, _ := range urls {
				_, err := s3client.ParseURL(missingBucketUrl, style)
				So(err, ShouldNotBeNil)
			}
		})

		Convey("Trying to parse a malformed s3 url results in error", func() {
			malformedURL := "This%Url%Is%Malformed"
			for style, _ := range urls {
				_, err := s3client.ParseURL(malformedURL, style)
				So(err, ShouldNotBeNil)
			}
		})
	})
}

func TestNewURL(t *testing.T) {

	const expectedBucketName = "csv-bucket"
	const expectedKey = "dir1/test-file.csv"
	const expectedRegion = "eu-west-1"

	Convey("Given valid region, bucket name and key results in New creating the expected S3Url struct", t, func() {
		s3Url, err := s3client.NewURL(expectedRegion, expectedBucketName, expectedKey)
		So(err, ShouldBeNil)
		So(s3Url.Scheme, ShouldEqual, "https")
		So(s3Url.Region, ShouldEqual, expectedRegion)
		So(s3Url.BucketName, ShouldEqual, expectedBucketName)
		So(s3Url.Key, ShouldEqual, expectedKey)
	})

	Convey("Given an empty region, valid bucket name and key results in New creating the expected S3Url struct", t, func() {
		s3Url, err := s3client.NewURL("", expectedBucketName, expectedKey)
		So(err, ShouldBeNil)
		So(s3Url.Scheme, ShouldEqual, "https")
		So(s3Url.Region, ShouldEqual, "")
		So(s3Url.BucketName, ShouldEqual, expectedBucketName)
		So(s3Url.Key, ShouldEqual, expectedKey)
	})

	Convey("Given an empty bucket results in error trying to create a new S3Url", t, func() {
		_, err := s3client.NewURL(expectedRegion, "", expectedKey)
		So(err, ShouldNotBeNil)
	})

	Convey("Given an empty key results in error trying to create a new S3Url", t, func() {
		_, err := s3client.NewURL(expectedRegion, expectedBucketName, "")
		So(err, ShouldNotBeNil)
	})

	Convey("Given a non-defult scheme, valid region, bucket name and key results in NewURLWithScheme creating the expected S3Url struct", t, func() {
		s3Url, err := s3client.NewURLWithScheme("s3", expectedRegion, expectedBucketName, expectedKey)
		So(err, ShouldBeNil)
		So(s3Url.Scheme, ShouldEqual, "s3")
		So(s3Url.Region, ShouldEqual, expectedRegion)
		So(s3Url.BucketName, ShouldEqual, expectedBucketName)
		So(s3Url.Key, ShouldEqual, expectedKey)
	})
}
