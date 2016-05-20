package tarball

import (
	boshblob "github.com/cloudfoundry/bosh-utils/blobstore"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
	bhreltarsrepo "github.com/cppforlife/bosh-hub/release/releasetarsrepo"
)

type Uploader struct {
	blobstore boshblob.Blobstore

	logTag string
	logger boshlog.Logger
}

func NewUploader(
	blobstore boshblob.Blobstore,
	logger boshlog.Logger,
) Uploader {
	return Uploader{
		blobstore: blobstore,

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

	relTarRec := bhreltarsrepo.ReleaseTarballRec{
		BlobID: blobID,
		SHA1:   sha,
	}

	err = relVerRec.SetTarball(relTarRec)
	if err != nil {
		return bosherr.WrapError(err, "Saving release into release tarballs repository")
	}

	return nil
}
