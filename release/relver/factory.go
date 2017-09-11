package relver

import (
	"regexp"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
)

type Factory struct {
	releasesIndexDir string
	fs               boshsys.FileSystem
	logger           boshlog.Logger
}

func NewFactory(releasesIndexDir string, fs boshsys.FileSystem, logger boshlog.Logger) Factory {
	return Factory{releasesIndexDir: releasesIndexDir, fs: fs, logger: logger}
}

var (
	sourceChars  = regexp.MustCompile(`\Agithub.com/[a-zA-Z\-0-9\/_]+\z`)
	versionChars = regexp.MustCompile(`\A[a-zA-Z-0-9\._+-]+\z`)
)

func (f Factory) Find(source, versionRaw string) (RelVer, error) {
	if !sourceChars.MatchString(source) {
		return RelVer{}, bosherr.New("Release version: Invalid source")
	}

	if !versionChars.MatchString(versionRaw) {
		return RelVer{}, bosherr.New("Invalid version")
	}

	return RelVer{source: source, versionRaw: versionRaw, releasesIndexDir: f.releasesIndexDir, fs: f.fs}, nil
}
