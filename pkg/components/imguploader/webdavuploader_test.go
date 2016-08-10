package imguploader

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUploadToWebdav(t *testing.T) {
	webdavUploader, _ := NewWebdavImageUploader("http://localhost:9999/", "username", "password")

	Convey("Can upload image to webdav server", t, func() {
		path, err := webdavUploader.Upload("~/dev/cyber-crime-mr-robot.png")

		So(err, ShouldBeNil)
		So(path, ShouldNotEqual, "")
	})
}
