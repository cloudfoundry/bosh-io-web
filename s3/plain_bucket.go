package s3

import (
	"strings"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
)

type PlainBucket struct {
	url string

	logTag string
	logger boshlog.Logger
}

func NewPlainBucket(url string, logger boshlog.Logger) Bucket {
	return PlainBucket{url, "PlainBucket", logger}
}

func (b PlainBucket) Files() ([]File, error) {
	page := NewPlainBucketPage(b.url, &b, b.logger)

	b.logger.Debug(b.logTag, "Fetching all pages for bucket at url '%s'", b.url)

	var allFiles []File

	for page != nil {
		files, err := page.Files()
		if err != nil {
			return allFiles, bosherr.WrapError(err, "Fetching bucket page")
		}

		allFiles = append(allFiles, files...)

		page, err = page.Next()
		if err != nil {
			return allFiles, bosherr.WrapError(err, "Finding out next bucket page")
		}
	}

	return allFiles, nil
}

func (b PlainBucket) URL() string { return b.url }

func (b PlainBucket) ObjectURL(key string) (string, error) {
	// do not use path since it will collapse http://
	return strings.TrimSuffix(b.url, "/") + "/" + key, nil
}
