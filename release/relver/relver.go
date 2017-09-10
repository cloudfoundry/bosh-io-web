package relver

import (
	"path/filepath"
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
)

type RelVer struct {
	source string
	versionRaw string

	releasesIndexDir string
	fs boshsys.FileSystem
}

func (r RelVer) ReadV1(filename string, out interface{}) error {
	foundPaths, err := r.fs.Glob(filepath.Join(r.releasesIndexDir, r.source, "*-"+r.versionRaw, filename+".v1.yml"))
	if err != nil {
		return bosherr.WrapError(err, "Globbing release versions")
	}

	if len(foundPaths) != 1 {
		return bosherr.WrapError(err, "Finding release version")
	}

	contents, err := r.fs.ReadFile(foundPaths[0])
	if err != nil {
		return bosherr.WrapError(err, "Reading release job file")
	}

	err = json.Unmarshal(contents, &out)
	if err != nil {
		return bosherr.WrapError(err, "Unmarshaling release jobs")
	}

	return nil
}
