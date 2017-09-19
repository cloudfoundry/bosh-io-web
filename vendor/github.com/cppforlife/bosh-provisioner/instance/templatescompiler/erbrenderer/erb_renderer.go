package erbrenderer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

const erbRendererLogTag = "ERBRenderer"

type ERBRenderer struct {
	fs      boshsys.FileSystem
	runner  boshsys.CmdRunner
	context TemplateEvaluationContext
	logger  boshlog.Logger

	rendererScript string
}

func NewERBRenderer(
	fs boshsys.FileSystem,
	runner boshsys.CmdRunner,
	context TemplateEvaluationContext,
	logger boshlog.Logger,
) ERBRenderer {
	return ERBRenderer{
		fs:      fs,
		runner:  runner,
		context: context,
		logger:  logger,

		rendererScript: templateEvaluationContextRb,
	}
}

func (r ERBRenderer) Render(srcPath, dstPath string) error {
	r.logger.Debug(erbRendererLogTag, "Rendering template %s", dstPath)

	dirPath := filepath.Dir(dstPath)

	err := r.fs.MkdirAll(dirPath, os.FileMode(0755))
	if err != nil {
		return bosherr.WrapErrorf(err, "Creating directory %s", dirPath)
	}

	rendererScriptPath, err := r.writeRendererScript()
	if err != nil {
		return err
	}

	contextPath, err := r.writeContext()
	if err != nil {
		return err
	}

	// Use ruby to compile job templates
	command := boshsys.Command{
		Name: r.determineRubyExePath(),
		Args: []string{rendererScriptPath, contextPath, srcPath, dstPath},
	}

	_, _, _, err = r.runner.RunComplexCommand(command)
	if err != nil {
		return bosherr.WrapError(err, "Running ruby")
	}

	return nil
}

func (r ERBRenderer) writeRendererScript() (string, error) {
	// todo use temp path; it's same everytime?
	path := "/tmp/erb-render.rb"

	err := r.fs.WriteFileString(path, r.rendererScript)
	if err != nil {
		return "", bosherr.WrapError(err, "Writing renderer script")
	}

	return path, nil
}

func (r ERBRenderer) writeContext() (string, error) {
	contextBytes, err := json.Marshal(r.context)
	if err != nil {
		return "", bosherr.WrapError(err, "Marshalling context")
	}

	// todo use temp path?
	path := "/tmp/erb-context.json"

	err = r.fs.WriteFileString(path, string(contextBytes))
	if err != nil {
		return "", bosherr.WrapError(err, "Writing context")
	}

	return path, nil
}

func (r ERBRenderer) determineRubyExePath() string {
	rubies := []string{
		"ruby",
		"/usr/bin/ruby1.9.3",
		"/opt/vagrant_ruby/bin/ruby",
	}

	// ruby 1.8.7 fails with "no such file to load -- rubygems"
	rubyCheck := "require 'rubygems'; puts 'works'"

	for _, path := range rubies {
		stdout, _, _, err := r.runner.RunCommand(path, "-e", rubyCheck)
		if err != nil {
			continue
		}

		if strings.Contains(stdout, "works") {
			return path
		}
	}

	return rubies[0] // give up
}
