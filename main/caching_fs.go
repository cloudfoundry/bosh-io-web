package main

import (
	"os"
	"sync"
	"path/filepath"

	boshsys "github.com/cloudfoundry/bosh-agent/system"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
)

type CachingFileSystem struct {
	fs boshsys.FileSystem

	readCache map[string][]byte
	readCacheLock sync.RWMutex

	globCache map[string][]string
	globCacheLock sync.RWMutex

	logTag string
	logger boshlog.Logger
}

var _ boshsys.FileSystem = &CachingFileSystem{}

func NewCachingFileSystem(fs boshsys.FileSystem, logger boshlog.Logger) *CachingFileSystem {
	return &CachingFileSystem{
		fs: fs,

		readCache: map[string][]byte{},
		globCache: map[string][]string{},

		logTag: "CachingFileSystem",
		logger: logger,
	}
}

func (f CachingFileSystem) HomeDir(username string) (path string, err error) {
  return f.fs.HomeDir(username)
}

func (f CachingFileSystem) MkdirAll(path string, perm os.FileMode) (err error) {
  return f.fs.MkdirAll(path, perm)
}

func (f CachingFileSystem) RemoveAll(fileOrDir string) (err error) {
  return f.fs.RemoveAll(fileOrDir)
}

func (f CachingFileSystem) Chown(path, username string) (err error) {
  return f.fs.Chown(path, username)
}

func (f CachingFileSystem) Chmod(path string, perm os.FileMode) (err error) {
  return f.fs.Chmod(path, perm)
}

func (f CachingFileSystem) OpenFile(path string, flag int, perm os.FileMode) (boshsys.ReadWriteCloseStater, error) {
  return f.fs.OpenFile(path, flag, perm)
}

func (f CachingFileSystem) WriteFileString(path, content string) (err error) {
  return f.fs.WriteFileString(path, content)
}

func (f CachingFileSystem) WriteFile(path string, content []byte) (err error) {
  return f.fs.WriteFile(path, content)
}

func (f CachingFileSystem) ConvergeFileContents(path string, content []byte) (written bool, err error) {
  return f.fs.ConvergeFileContents(path, content)
}

func (f *CachingFileSystem) ReadFileString(path string) (string, error) {
	bytes, err := f.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (f CachingFileSystem) ReadFile(path string) ([]byte, error) {
	f.readCacheLock.Lock()
	defer f.readCacheLock.Unlock()

	if content, found := f.readCache[path]; found {
		f.logger.Debug(f.logTag, "hit: read[%s]", path)
		return content, nil
	} else {
		f.logger.Debug(f.logTag, "miss: read[%s]", path)
	}

  content, err := f.fs.ReadFile(path)
  if err == nil {
  	f.readCache[path] = content
  }

  return content, err
}

func (f CachingFileSystem) FileExists(path string) bool {
  return f.fs.FileExists(path)
}

func (f CachingFileSystem) Rename(oldPath, newPath string) (err error) {
  return f.fs.Rename(oldPath, newPath)
}

func (f CachingFileSystem) Symlink(oldPath, newPath string) (err error) {
  return f.fs.Symlink(oldPath, newPath)
}

func (f CachingFileSystem) ReadLink(symlinkPath string) (targetPath string, err error) {
  return f.fs.ReadLink(symlinkPath)
}

func (f CachingFileSystem) CopyFile(srcPath, dstPath string) (err error) {
  return f.fs.CopyFile(srcPath, dstPath)
}

func (f CachingFileSystem) TempFile(prefix string) (file *os.File, err error) {
  return f.fs.TempFile(prefix)
}

func (f CachingFileSystem) TempDir(prefix string) (path string, err error) {
  return f.fs.TempDir(prefix)
}

func (f CachingFileSystem) Glob(pattern string) ([]string, error) {
	f.globCacheLock.Lock()
	defer f.globCacheLock.Unlock()

	if matches, found := f.globCache[pattern]; found {
		f.logger.Debug(f.logTag, "hit: glob[%s]", pattern)
		return matches, nil
	} else {
		f.logger.Debug(f.logTag, "miss: glob[%s]", pattern)
	}

  matches, err := f.fs.Glob(pattern)
  if err == nil {
  	f.globCache[pattern] = matches
  }

  return matches, err
}

func (f CachingFileSystem) Walk(root string, walkFunc filepath.WalkFunc) error {
  return f.fs.Walk(root, walkFunc)
}
