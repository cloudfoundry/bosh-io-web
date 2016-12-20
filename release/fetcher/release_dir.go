package fetcher

import (
	"path/filepath"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
	bprelman "github.com/cppforlife/bosh-provisioner/release/manifest"
)

type ReleaseDir struct {
	releaseDir  string
	cleanUpFunc func() error

	fs     boshsys.FileSystem
	logTag string
	logger boshlog.Logger
}

func NewReleaseDir(
	releaseDir string,
	cleanUpFunc func() error,
	fs boshsys.FileSystem,
	logger boshlog.Logger,
) ReleaseDir {
	return ReleaseDir{
		releaseDir:  releaseDir,
		cleanUpFunc: cleanUpFunc,

		fs:     fs,
		logTag: "ReleaseDir",
		logger: logger,
	}
}

func (r ReleaseDir) ReleaseManifests() (map[string]bprelman.Manifest, error) {
	return r.findReleaseManifestPaths(r.releaseDir)
}

func (r ReleaseDir) Close() error { return r.cleanUpFunc() }

func (r ReleaseDir) findReleaseManifestPaths(releaseDir string) (map[string]bprelman.Manifest, error) {
	manifests := map[string]bprelman.Manifest{}

	baseDir := releaseDir

	releaseMatches, err := r.fs.Glob(filepath.Join(baseDir, "releases/**/*"))
	if err != nil {
		return nil, bosherr.WrapError(err, "Globbing releases/**/* directory")
	}

	nonDirReleaseMatches, err := r.fs.Glob(filepath.Join(baseDir, "releases/*"))
	if err != nil {
		return nil, bosherr.WrapError(err, "Globbing releases/* directory")
	}

	for _, releaseMatch := range nonDirReleaseMatches {
		releaseMatches = append(releaseMatches, releaseMatch)
	}

	r.logger.Debug(r.logTag, "Found '%d' release metadata files", len(releaseMatches))

	// Release manifest file names will be unique even if they are in release sub-directories
	// e.g. relname-1.yml, relname-0.567.yml
	lookedAtFileNames := map[string]struct{}{}

	for _, releaseMatch := range releaseMatches {
		baseFileName := filepath.Base(releaseMatch)

		// index.yml is not a release manifest or it's not a yml file
		if baseFileName == "index.yml" || filepath.Ext(baseFileName) != ".yml" {
			continue
		}

		// releaseMatches and nonDirReleaseMatches will contain dup items
		if _, ok := lookedAtFileNames[baseFileName]; ok {
			r.logger.Debug(r.logTag, "Skipping over dup release match '%s'", releaseMatch)
			continue
		}

		manifest, err := bprelman.NewManifestFromPath(releaseMatch, r.fs)
		if err != nil {
			// todo should this failure be shown somewhere?
			// e.g. https://github.com/cloudfoundry/bosh/blob/master/release/releases/bosh-1.yml
			r.logger.Debug(r.logTag, "Failed building manifest '%s'", releaseMatch)
			continue
		}

		lookedAtFileNames[baseFileName] = struct{}{}

		r.logger.Debug(r.logTag, "Reading manifest '%s'", releaseMatch)

		manifests[releaseMatch] = manifest
	}

	return manifests, nil
}
