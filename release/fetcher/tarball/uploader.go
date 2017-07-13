package tarball

import (
	"strings"

	boshblob "github.com/cloudfoundry/bosh-agent/blobstore"
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"

	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
	bhreltarsrepo "github.com/cppforlife/bosh-hub/release/releasetarsrepo"
)

type Uploader struct {
	blobstore boshblob.Blobstore
	runner    boshsys.CmdRunner

	logTag string
	logger boshlog.Logger
}

func NewUploader(
	blobstore boshblob.Blobstore,
	runner boshsys.CmdRunner,
	logger boshlog.Logger,
) Uploader {
	return Uploader{
		blobstore: blobstore,
		runner:    runner,

		logTag: "Uploader",
		logger: logger,
	}
}

func (u Uploader) Upload(relVerRec bhrelsrepo.ReleaseVersionRec, tgzPath string) error {
	u.logger.Info(u.logTag, "Uploading '%s' for '%v'", tgzPath, relVerRec)

	blobID, sha, err := u.blobstore.Create(tgzPath)
	if err != nil {
		return bosherr.WrapError(err, "Creating release tarball")
	}

	cmd := boshsys.Command{
		Name: "shasum",
		Args: []string{"-a", "256", tgzPath},
	}

	sha256, _, _, err := u.runner.RunComplexCommand(cmd)
	if err != nil {
		return bosherr.WrapError(err, "Running bosh create release")
	}

	relTarRec := bhreltarsrepo.ReleaseTarballRec{
		BlobID: blobID,
		SHA1:   sha,
		Digest: "sha256:" + strings.Trim(sha256, " "),
	}

	err = relVerRec.SetTarball(relTarRec)
	if err != nil {
		return bosherr.WrapError(err, "Saving release into release tarballs repository")
	}

	return nil
}
