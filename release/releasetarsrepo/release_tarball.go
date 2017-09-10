package releasetarsrepo

import (
	"fmt"
	"strings"

	bhs3 "github.com/cppforlife/bosh-hub/s3"
)

type ReleaseTarballRec struct {
	urlFactory bhs3.URLFactory

	source     string
	versionRaw string

	BlobID string
	SHA1   string
}

func (r ReleaseTarballRec) ActualDownloadURL() (string, error) {
	path := "/" + r.BlobID

	fileName := fmt.Sprintf("%s-%s.tgz", r.sourceShortName(), r.versionRaw)

	return r.urlFactory.New(path, fileName).String()
}

func (r ReleaseTarballRec) sourceShortName() string {
	parts := strings.Split(r.source, "/")

	return parts[len(parts)-1]
}
