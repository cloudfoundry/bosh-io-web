package downloader

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

const gitDownloaderLogTag = "GitDownloader"

type GitDownloader struct {
	fs     boshsys.FileSystem
	runner boshsys.CmdRunner
	logger boshlog.Logger
}

func NewGitDownloader(
	fs boshsys.FileSystem,
	runner boshsys.CmdRunner,
	logger boshlog.Logger,
) GitDownloader {
	return GitDownloader{fs: fs, runner: runner, logger: logger}
}

func (d GitDownloader) Download(url string) (string, error) {
	file, err := d.fs.TempFile("downloader-GitDownloader")
	if err != nil {
		return "", bosherr.WrapError(err, "Creating download destination")
	}

	d.logger.Debug(gitDownloaderLogTag, "Planning to download '%s' to '%s'", url, file.Name())

	defer file.Close()

	err = d.fs.RemoveAll(file.Name())
	if err != nil {
		return "", bosherr.WrapError(err, "Clearing download destination")
	}

	dir := file.Name()

	_, _, _, err = d.runner.RunCommand("git", "clone", url, dir, "--depth", "1")
	if err != nil {
		return "", bosherr.WrapError(err, "Cloning git repository")
	}

	return dir, nil
}

func (d GitDownloader) CleanUp(dir string) error {
	return d.fs.RemoveAll(dir)
}
