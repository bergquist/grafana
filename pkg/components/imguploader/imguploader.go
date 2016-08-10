package imguploader

import (
	"fmt"

	"github.com/grafana/grafana/pkg/setting"
)

type ImageUploader interface {
	Upload(path string) (string, error)
}

func NewImageUploader() (ImageUploader, error) {

	switch setting.ImageUploadProvider {
	case "s3":
		return NewS3Uploader(setting.S3TempImageStoreBucketUrl, setting.S3TempImageStoreAccessKey, setting.S3TempImageStoreSecretKey), nil
	case "webdav":
		return NewWebdavImageUploader("", "", "")
	}

	return nil, fmt.Errorf("could not find specified provider")
}
