package imguploader

import (
	"reflect"
	"testing"

	"github.com/grafana/grafana/pkg/setting"

	. "github.com/smartystreets/goconvey/convey"
)

func TestImageUploaderFactory(t *testing.T) {
	Convey("Can create image uploader for ", t, func() {
		Convey("S3ImageUploader", func() {
			var err error
			err = setting.NewConfigContext(&setting.CommandLineArgs{
				HomePath: "../../",
			})

			So(err, ShouldBeNil)

			sec, err := setting.Cfg.NewSection("image.uploader")
			sec.NewKey("remote", "s3")

			s3sec, err := setting.Cfg.NewSection("image.uploader.s3")
			s3sec.NewKey("bucket_url", "bucket_url")
			s3sec.NewKey("access_key", "access_key")
			s3sec.NewKey("secret_key", "secret_key")

			So(err, ShouldBeNil)

			uploader, err := NewImageUploader()

			So(err, ShouldBeNil)
			So(reflect.TypeOf(uploader), ShouldEqual, "S3Uploader")
		})

		SkipConvey("Webdav uploader", func() {
			var err error
			err = setting.NewConfigContext(&setting.CommandLineArgs{
				HomePath: "../../",
			})

			So(err, ShouldBeNil)

			sec, err := setting.Cfg.NewSection("image.uploader")
			sec.NewKey("remote", "webdav")

			s3sec, err := setting.Cfg.NewSection("image.uploader.webdav")
			s3sec.NewKey("url", "webdavUrl")
			s3sec.NewKey("username", "username")
			s3sec.NewKey("password", "password")

			So(err, ShouldBeNil)

			uploader, err := NewImageUploader()

			So(err, ShouldBeNil)
			So(reflect.TypeOf(uploader), ShouldEqual, "WebDavUploader")
		})
	})
}
