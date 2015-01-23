package releasetarsrepo

import (
	"fmt"

	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
	bhs3 "github.com/cppforlife/bosh-hub/s3"
)

type ReleaseTarballRec struct {
	urlFactory bhs3.URLFactory
	relVerRec  bhrelsrepo.ReleaseVersionRec

	BlobID string
	SHA1   string
}

func (r ReleaseTarballRec) ActualDownloadURL() (string, error) {
	path := "/" + r.BlobID

	fileName := fmt.Sprintf("%s-%s.tgz", r.relVerRec.SourceShortName(), r.relVerRec.VersionRaw)

	return r.urlFactory.New(path, fileName).String()
}
