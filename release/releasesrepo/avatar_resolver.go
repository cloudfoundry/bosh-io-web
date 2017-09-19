package releasesrepo

import (
	"gopkg.in/yaml.v2"
	"path/filepath"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type predefinedAvatarsResolver struct {
	releasesDir string
	fs          boshsys.FileSystem
}

func (r predefinedAvatarsResolver) Resolve(location string) string {
	defs, err := r.defs()
	if err != nil {
		return ""
	}

	for _, def := range defs {
		if def.Location() == location {
			return def.URL
		}
	}

	return ""
}

type avatarDefYAML struct {
	RepoURL string `yaml:"repo_url"`
	URL     string
}

func (d avatarDefYAML) Location() string {
	return strings.TrimPrefix(d.RepoURL, "https://")
}

func (r predefinedAvatarsResolver) defs() ([]avatarDefYAML, error) {
	contents, err := r.fs.ReadFileString(filepath.Join(r.releasesDir, "avatars.yml"))
	if err != nil {
		return nil, bosherr.WrapError(err, "Reading releases")
	}

	var defs []avatarDefYAML

	err = yaml.Unmarshal([]byte(contents), &defs)
	if err != nil {
		return nil, bosherr.WrapError(err, "Unmarshaling releases")
	}

	return defs, nil
}
