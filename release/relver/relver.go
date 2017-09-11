package relver

import (
	"path/filepath"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
)

type RelVer struct {
	source     string
	versionRaw string

	releasesIndexDir string
	fs               boshsys.FileSystem
}

func (r RelVer) Read(fileName string) ([]byte, error) {
	foundPaths, err := r.fs.Glob(filepath.Join(r.releasesIndexDir, r.source, "*-"+r.versionRaw, fileName))
	if err != nil {
		return nil, bosherr.WrapError(err, "Globbing release versions")
	}

	if len(foundPaths) != 1 {
		return nil, bosherr.WrapError(err, "Finding release version")
	}

	contents, err := r.fs.ReadFile(foundPaths[0])
	if err != nil {
		return nil, bosherr.WrapError(err, "Reading release file")
	}

	return contents, err
}

func (r RelVer) ReadOptinal(fileName string) ([]byte, bool, error) {
	foundPaths, err := r.fs.Glob(filepath.Join(r.releasesIndexDir, r.source, "*-"+r.versionRaw, fileName))
	if err != nil {
		return nil, false, bosherr.WrapError(err, "Globbing release versions")
	}

	if len(foundPaths) > 1 {
		return nil, false, bosherr.WrapError(err, "Finding release version")
	}

	if len(foundPaths) == 1 {
		contents, err := r.fs.ReadFile(foundPaths[0])
		if err != nil {
			return nil, false, bosherr.WrapError(err, "Reading release file")
		}

		return contents, true, err
	}

	return nil, false, nil
}
