package releasesrepo

import (
	"encoding/json"
	"path/filepath"
	"regexp"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"

	bprel "github.com/cppforlife/bosh-provisioner/release"
)

type CRVRepository struct {
	releasesIndexDir string
	fs               boshsys.FileSystem
	logger           boshlog.Logger
}

func NewConcreteReleaseVersionsRepository(
	releasesIndexDir string,
	fs boshsys.FileSystem,
	logger boshlog.Logger,
) CRVRepository {
	return CRVRepository{
		releasesIndexDir: releasesIndexDir,
		fs:               fs,
		logger:           logger,
	}
}

var (
	sourceChars  = regexp.MustCompile(`\Agithub.com/[a-zA-Z\-0-9\/_]+\z`)
	versionChars = regexp.MustCompile(`\A[a-zA-Z-0-9\._+-]+\z`)
)

func (r CRVRepository) Find(relVerRec ReleaseVersionRec) (bprel.Release, error) {
	var rel bprel.Release

	if !sourceChars.MatchString(relVerRec.Source) {
		return rel, bosherr.New("Release version: Invalid source")
	}

	if !versionChars.MatchString(relVerRec.VersionRaw) {
		return rel, bosherr.New("Invalid version")
	}

	foundPaths, err := r.fs.Glob(filepath.Join(r.releasesIndexDir, relVerRec.Source, "*-"+relVerRec.VersionRaw, "release.v1.yml"))
	if err != nil {
		return rel, bosherr.WrapError(err, "Globbing release versions")
	}

	if len(foundPaths) != 1 {
		return rel, bosherr.WrapError(err, "Finding release version")
	}

	contents, err := r.fs.ReadFile(foundPaths[0])
	if err != nil {
		return rel, bosherr.WrapError(err, "Reading release file")
	}

	err = json.Unmarshal(contents, &rel)
	if err != nil {
		return rel, bosherr.WrapError(err, "Unmarshaling release")
	}

	return rel, nil
}
