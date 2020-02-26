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
			So(s3Url.Region(), ShouldEqual, expectedRegion)
			So(s3Url.BucketName(), ShouldEqual, expectedBucketName)
			So(s3Url.Key(), ShouldEqual, expectedKey)
		})

		Convey("Path style URL string is formatted as expected", func() {
			expectedStr := fmt.Sprintf("https://s3-%s.amazonaws.com/%s/%s", expectedRegion, expectedBucketName, expectedKey)
			urlStr, err := s3Url.String(s3client.StylePath)
			So(err, ShouldBeNil)
			So(urlStr, ShouldEqual, expectedStr)
		})

		Convey("Global Path style URL string is formatted as expected", func() {
			expectedStr := fmt.Sprintf("https://s3.amazonaws.com/%s/%s", expectedBucketName, expectedKey)
			urlStr, err := s3Url.String(s3client.StyleGlobalPath)
			So(err, ShouldBeNil)
			So(urlStr, ShouldEqual, expectedStr)
		})

		Convey("Virtual hosted style URL string is formatted as expected", func() {
			expectedStr := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", expectedBucketName, expectedRegion, expectedKey)
			urlStr, err := s3Url.String(s3client.StyleVirtualHosted)
			So(err, ShouldBeNil)
			So(urlStr, ShouldEqual, expectedStr)
		})

		Convey("Global virtual hosted style URL string is formatted as expected", func() {
			expectedStr := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", expectedBucketName, expectedKey)
			urlStr, err := s3Url.String(s3client.StyleGlobalVirtualHosted)
			So(err, ShouldBeNil)
			So(urlStr, ShouldEqual, expectedStr)
		})

		Convey("DNS Alias virtual hosted style URL string is formatted as expected", func() {
			expectedStr := fmt.Sprintf("https://%s/%s", expectedBucketName, expectedKey)
			urlStr, err := s3Url.String(s3client.StyleAliasVirtualHosted)
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
			So(s3Url.Region(), ShouldEqual, "")
			So(s3Url.BucketName(), ShouldEqual, expectedBucketName)
			So(s3Url.Key(), ShouldEqual, expectedKey)
		})

		Convey("Path style URL string is formatted as expected", func() {
			_, err := s3Url.String(s3client.StylePath)
			So(err, ShouldNotBeNil)
		})

		Convey("Global Path style URL string is formatted as expected", func() {
			expectedStr := fmt.Sprintf("https://s3.amazonaws.com/%s/%s", expectedBucketName, expectedKey)
			urlStr, err := s3Url.String(s3client.StyleGlobalPath)
			So(err, ShouldBeNil)
			So(urlStr, ShouldEqual, expectedStr)
		})

		Convey("Virtual hosted style URL string is formatted as expected", func() {
			_, err := s3Url.String(s3client.StyleVirtualHosted)
			So(err, ShouldNotBeNil)
		})

		Convey("Global virtual hosted style URL string is formatted as expected", func() {
			expectedStr := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", expectedBucketName, expectedKey)
			urlStr, err := s3Url.String(s3client.StyleGlobalVirtualHosted)
			So(err, ShouldBeNil)
			So(urlStr, ShouldEqual, expectedStr)
		})

		Convey("DNS Alias virtual hosted style URL string is formatted as expected", func() {
			expectedStr := fmt.Sprintf("https://%s/%s", expectedBucketName, expectedKey)
			urlStr, err := s3Url.String(s3client.StyleAliasVirtualHosted)
			So(err, ShouldBeNil)
			So(urlStr, ShouldEqual, expectedStr)
		})
	})
}

func TestParsing(t *testing.T) {

	const expectedBucketName = "csv-bucket"
	const expectedKey = "dir1/test-file.csv"
	const expectedRegion = "eu-west-1"

	styles := []s3client.URLStyle{
		s3client.StylePath,
		s3client.StyleGlobalPath,
		s3client.StyleVirtualHosted,
		s3client.StyleGlobalVirtualHosted,
		s3client.StyleAliasVirtualHosted,
	}

	Convey("Given S3 raw url strings in different acceptable formats", t, func() {
		// urls that define region
		regionalUrls := map[s3client.URLStyle][]string{
			s3client.StylePath: []string{
				"https://s3-eu-west-1.amazonaws.com/csv-bucket/dir1/test-file.csv",
				"s3://s3-eu-west-1.amazonaws.com/csv-bucket/dir1/test-file.csv"},
			s3client.StyleVirtualHosted: []string{
				"https://csv-bucket.s3-eu-west-1.amazonaws.com/dir1/test-file.csv",
				"s3://csv-bucket.s3-eu-west-1.amazonaws.com/dir1/test-file.csv"},
		}
		expectedRegionalS3Url, err := s3client.NewURL(expectedRegion, expectedBucketName, expectedKey)
		So(err, ShouldBeNil)

		// urls that don't define region
		globalUrls := map[s3client.URLStyle][]string{
			s3client.StyleGlobalPath: []string{
				"https://s3.amazonaws.com/csv-bucket/dir1/test-file.csv",
				"s3://s3.amazonaws.com/csv-bucket/dir1/test-file.csv"},
			s3client.StyleGlobalVirtualHosted: []string{
				"https://csv-bucket.s3.amazonaws.com/dir1/test-file.csv",
				"s3://csv-bucket.s3.amazonaws.com/dir1/test-file.csv"},
			s3client.StyleAliasVirtualHosted: []string{
				"https://csv-bucket/dir1/test-file.csv",
				"s3://csv-bucket/dir1/test-file.csv"},
		}
		expectedGlobalS3Url, err := s3client.NewURL("", expectedBucketName, expectedKey)
		So(err, ShouldBeNil)

		Convey("Each format is correctly parsed, successfully retreiving bucket, key and region (if available)", func() {
			for style, urls := range regionalUrls {
				for _, url := range urls {
					s3Url, err := s3client.ParseURL(url, style)
					So(err, ShouldBeNil)
					So(s3Url, ShouldResemble, expectedRegionalS3Url)
				}
			}
			for style, urls := range globalUrls {
				for _, url := range urls {
					s3Url, err := s3client.ParseURL(url, style)
					So(err, ShouldBeNil)
					So(s3Url, ShouldResemble, expectedGlobalS3Url)
				}
			}
		})

		Convey("A path-style url can be parsed as a global-path-style with empty region", func() {
			for _, url := range regionalUrls[s3client.StylePath] {
				s3Url, err := s3client.ParseURL(url, s3client.StyleGlobalPath)
				So(err, ShouldBeNil)
				So(s3Url, ShouldResemble, expectedGlobalS3Url)
			}
		})

		Convey("A global-path-style url can be parsed as a path-style with empty region", func() {
			for _, url := range regionalUrls[s3client.StyleGlobalPath] {
				s3Url, err := s3client.ParseURL(url, s3client.StylePath)
				So(err, ShouldBeNil)
				So(s3Url, ShouldResemble, expectedGlobalS3Url)
			}
		})
	})

	Convey("Trying to parse an empty S3 raw url results in error ", t, func() {
		for _, style := range styles {
			_, err := s3client.ParseURL("", style)
			So(err, ShouldNotBeNil)
		}
	})

	Convey("Tying to parse an s3 url that is missing the object key, results in error", t, func() {
		missingBucketUrl := "s3://some-file"
		for _, style := range styles {
			_, err := s3client.ParseURL(missingBucketUrl, style)
			So(err, ShouldNotBeNil)
		}
	})

	Convey("Trying to parse an s3 url with empty bucket or key results in error", t, func() {
		emptyValuesUrl1 := "s3://///////"
		emptyValuesUrl2 := fmt.Sprintf("s3:/%s/", expectedBucketName)
		for _, style := range styles {
			_, err := s3client.ParseURL(emptyValuesUrl1, style)
			So(err, ShouldNotBeNil)
			_, err = s3client.ParseURL(emptyValuesUrl2, style)
			So(err, ShouldNotBeNil)
		}
	})

	Convey("Tying to parse an s3 url that is missing the bucket name and object key results in error", t, func() {
		missingBucketUrl := "s3://"
		for _, style := range styles {
			_, err := s3client.ParseURL(missingBucketUrl, style)
			So(err, ShouldNotBeNil)
		}
	})

	Convey("Trying to parse a malformed s3 url results in error", t, func() {
		malformedURL := "This%Url%Is%Malformed"
		for _, style := range styles {
			_, err := s3client.ParseURL(malformedURL, style)
			So(err, ShouldNotBeNil)
		}
	})
}

func TestNewURL(t *testing.T) {

	const expectedBucketName = "csv-bucket"
	const expectedKey = "dir1/test-file.csv"
	const expectedRegion = "eu-west-1"

	Convey("Given valid region, bucket name and key results in New creating the expected S3Url struct", t, func() {
		s3Url, err := s3client.NewURL(expectedRegion, expectedBucketName, expectedKey)
		So(err, ShouldBeNil)
		So(s3Url.Region(), ShouldEqual, expectedRegion)
		So(s3Url.BucketName(), ShouldEqual, expectedBucketName)
		So(s3Url.Key(), ShouldEqual, expectedKey)
	})

	Convey("Given an empty region, valid bucket name and key results in New creating the expected S3Url struct", t, func() {
		s3Url, err := s3client.NewURL("", expectedBucketName, expectedKey)
		So(err, ShouldBeNil)
		So(s3Url.Region(), ShouldEqual, "")
		So(s3Url.BucketName(), ShouldEqual, expectedBucketName)
		So(s3Url.Key(), ShouldEqual, expectedKey)
	})

	Convey("Given an empty bucket results in error trying to create a new S3Url", t, func() {
		_, err := s3client.NewURL(expectedRegion, "", expectedKey)
		So(err, ShouldNotBeNil)
	})

	Convey("Given an empty key results in error trying to create a new S3Url", t, func() {
		_, err := s3client.NewURL(expectedRegion, expectedBucketName, "")
		So(err, ShouldNotBeNil)
	})
}
