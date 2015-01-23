package tarball

import (
	boshblob "github.com/cloudfoundry/bosh-agent/blobstore"
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
	bhreltarsrepo "github.com/cppforlife/bosh-hub/release/releasetarsrepo"
)

type Uploader struct {
	blobstore       boshblob.Blobstore
	releaseTarsRepo bhreltarsrepo.ReleaseTarballsRepository

	logTag string
	logger boshlog.Logger
}

func NewUploader(
	blobstore boshblob.Blobstore,
	releaseTarsRepo bhreltarsrepo.ReleaseTarballsRepository,
	logger boshlog.Logger,
) Uploader {
	return Uploader{
		blobstore:       blobstore,
		releaseTarsRepo: releaseTarsRepo,

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

	err = u.releaseTarsRepo.Save(relVerRec, relTarRec)
	if err != nil {
		return bosherr.WrapError(err, "Saving release into release tarballs repository")
	}

	return nil
}
