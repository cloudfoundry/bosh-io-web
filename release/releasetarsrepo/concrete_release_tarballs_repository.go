package releasetarsrepo

import (
	"encoding/xml"
	"regexp"
	"path/filepath"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	bhs3 "github.com/cppforlife/bosh-hub/s3"
)

type CRTRepository struct {
	releasesIndexDir      string
	urlFactory bhs3.URLFactory
	fs boshsys.FileSystem
	logger     boshlog.Logger
}

func NewConcreteReleaseTarballsRepository(
	releasesIndexDir string,
	urlFactory bhs3.URLFactory,
	fs boshsys.FileSystem,
	logger boshlog.Logger,
) CRTRepository {
	return CRTRepository{
		releasesIndexDir:      releasesIndexDir,
		urlFactory: urlFactory,
		fs: fs,
		logger:     logger,
	}
}

var (
	sourceChars = regexp.MustCompile(`\Agithub.com/[a-zA-Z\-0-9\/_]+\z`)
	versionChars = regexp.MustCompile(`\A[a-zA-Z-0-9\._+-]+\z`)
)

func (r CRTRepository) Find(source, version string) (ReleaseTarballRec, error) {
	var relTarRec ReleaseTarballRec

	if !sourceChars.MatchString(source) {
		return relTarRec, bosherr.New("Release tarball: Invalid source")
	}

	if !versionChars.MatchString(version) {
		return relTarRec, bosherr.New("Invalid version")
	}

	foundPaths, err := r.fs.Glob(filepath.Join(r.releasesIndexDir, source, "*-"+version, "source.meta4"))
	if err != nil {
		return relTarRec, bosherr.WrapError(err, "Globbing release versions")
	}

	if len(foundPaths) != 1 {
		return relTarRec, bosherr.WrapError(err, "Finding release version")
	}

	relTarRec.urlFactory = r.urlFactory
	relTarRec.source = source
	relTarRec.versionRaw = version

	var meta4 Metalink

	contents, err := r.fs.ReadFile(foundPaths[0])
	if err != nil {
		return relTarRec, bosherr.WrapError(err, "Reading meta4 file")
	}

	err = xml.Unmarshal(contents, &meta4)
	if err != nil {
		return relTarRec, bosherr.WrapError(err, "Unmarshaling meta4")
	}

	relTarRec.BlobID = filepath.Base(meta4.Files[0].URLs[0].URL)
	relTarRec.SHA1   = meta4.Files[0].Hashes[0].Hash

	return relTarRec, nil
}
