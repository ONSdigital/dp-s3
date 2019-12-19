package s3client_test

import (
	"testing"

	s3 "github.com/ONSdigital/dp-s3"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {

	const expectedBucketName = "csv-bucket"
	const expectedFilePath = "dir1/test-file.csv"
	const rawURL = "s3://" + expectedBucketName + "/" + expectedFilePath

	Convey("Given an instance of s3.URL with a valid fully qualified S3 URL", t, func() {

		s3URL, err := s3.NewURL(rawURL)
		So(err, ShouldBeNil)

		Convey("When the bucket name is requested", func() {

			bucketName := s3URL.BucketName()

			Convey("It should provide the expected value", func() {
				So(bucketName, ShouldEqual, expectedBucketName)
			})
		})

		Convey("When the file path is requested", func() {

			filePath := s3URL.Path()

			Convey("It should provide the expected value", func() {
				So(filePath, ShouldEqual, expectedFilePath)
			})
		})
	})
}
