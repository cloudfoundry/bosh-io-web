// Package manifest represents internal structure of a release job.
package manifest

import (
	"github.com/cloudfoundry-incubator/candiedyaml"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Manifest struct {
	Job Job
}

type Job struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`

	TemplateNames TemplateNames `yaml:"templates"`

	PackageNames []string `yaml:"packages"`

	PropertyMappings PropertyMappings `yaml:"properties"`
}

type TemplateNames map[string]string

type PropertyMappings map[string]PropertyDefinition

type PropertyDefinition struct {
	Description string `yaml:"description"`

	// Non-raw field is populated by the validator.
	DefaultRaw interface{} `yaml:"default"`
	Default    interface{}

	// Non-raw field is populated by the validator.
	ExampleRaw interface{} `yaml:"example"`
	Example    interface{}

	Examples []PropertyExampleDefinition `yaml:"examples"`
}

type PropertyExampleDefinition struct {
	Description string

	// Non-raw field is populated by the validator.
	ValueRaw interface{} `yaml:"value"`
	Value    interface{}
}

func NewManifestFromPath(path string, fs boshsys.FileSystem) (Manifest, error) {
	bytes, err := fs.ReadFile(path)
	if err != nil {
		return Manifest{}, bosherr.WrapErrorf(err, "Reading manifest %s", path)
	}

	return NewManifestFromBytes(bytes)
}

func NewManifestFromBytes(bytes []byte) (Manifest, error) {
	var manifest Manifest
	var job Job

	err := candiedyaml.Unmarshal(bytes, &job)
	if err != nil {
		return manifest, bosherr.WrapError(err, "Parsing job")
	}

	manifest.Job = job

	err = NewSyntaxValidator(&manifest).Validate()
	if err != nil {
		return Manifest{}, bosherr.WrapError(err, "Validating manifest syntactically")
	}

	return manifest, nil
}

/*
# Example for job.MF
name: dummy

description: ...

templates:
  dummy_ctl: bin/dummy_ctl

packages:
- dummy_package
- dummy_package2

properties:
  dummy_value:
    description: Some value for the dummy job
    default: 300
    example: ...
    examples:
    - description: Some description
      value: ...
*/
