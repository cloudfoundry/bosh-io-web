package tarball

import (
	"path/filepath"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
)

type Builder struct {
	fs     boshsys.FileSystem
	runner boshsys.CmdRunner

	logTag string
	logger boshlog.Logger
}

func NewBuilder(
	fs boshsys.FileSystem,
	runner boshsys.CmdRunner,
	logger boshlog.Logger,
) Builder {
	return Builder{
		fs:     fs,
		runner: runner,

		logTag: "Builder",
		logger: logger,
	}
}

func (tr Builder) Build(manifestPath string) (string, error) {
	tr.logger.Debug(tr.logTag, "Building tarball from '%s'", manifestPath)

	tgzPath := manifestPath + ".tgz"

	cmd := boshsys.Command{
		Name: "bosh2",
		Args: []string{"create-release", manifestPath, "--tarball", tgzPath},

		WorkingDir: tr.buildReleaseWorkingDir(manifestPath),
	}

	_, _, _, err := tr.runner.RunComplexCommand(cmd)
	if err != nil {
		return "", bosherr.WrapError(err, "Running bosh create release")
	}

	return tgzPath, nil
}

func (tr Builder) CleanUp(tgzPath string) error {
	return tr.fs.RemoveAll(tgzPath)
}

func (tr Builder) buildReleaseWorkingDir(manifestPath string) string {
	workingDir := filepath.Dir(filepath.Dir(manifestPath))

	// If newer release format is used get rid of subdirectory in releases/:
	// e.g. $dir/releases/relname-1.yml
	//      $dir/releases/relname/relname-1.yml
	if !tr.fs.FileExists(filepath.Join(workingDir, "releases")) {
		workingDir = filepath.Dir(workingDir)
	}

	return workingDir
}
