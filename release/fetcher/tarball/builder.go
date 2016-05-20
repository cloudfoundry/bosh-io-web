package tarball

import (
	"path/filepath"
	"regexp"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

var (
	/*
		Running command 'bosh create release releases/bosh-90.yml':
		...
		Generated /Users/pivotal/workspace/bosh/release/releases/bosh-90.tgz
		Release size: 123.9M
	*/
	builderReleasePathRegex = regexp.MustCompile(`Generated (.+\.tgz)`)
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

	cmd := boshsys.Command{
		Name: "bosh",
		Args: []string{"create", "release", manifestPath},

		WorkingDir: tr.buildReleaseWorkingDir(manifestPath),
	}

	stdout, _, _, err := tr.runner.RunComplexCommand(cmd)
	if err != nil {
		return "", bosherr.WrapError(err, "Running bosh create release")
	}

	pathMatches := builderReleasePathRegex.FindStringSubmatch(stdout)
	if len(pathMatches) != 2 {
		return "", bosherr.WrapErrorf(err, "tgz path was not found in '%s'", stdout)
	}

	return pathMatches[1], nil
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
