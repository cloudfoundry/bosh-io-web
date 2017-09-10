package main

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshsys "github.com/cloudfoundry/bosh-agent/system"

	bhctrls "github.com/cppforlife/bosh-hub/controllers"
	bhnoteimporter "github.com/cppforlife/bosh-hub/release/noteimporter"
	bhstemsimp "github.com/cppforlife/bosh-hub/stemcell/importer"
	bhstemnoteimporter "github.com/cppforlife/bosh-hub/stemcell/noteimporter"
)

type Config struct {
	Repos ReposOptions

	APIKey string

	Analytics AnalyticsConfig

	// Does not start web server; just does background work
	ActAsWorker bool

	ChecksumPrivs []bhctrls.ChecksumReqMatch

	ReleaseNoteImporter  bhnoteimporter.FactoryOptions
	StemcellNoteImporter bhstemnoteimporter.FactoryOptions

	StemcellImporter    bhstemsimp.FactoryOptions
}

type AnalyticsConfig struct {
	GoogleAnalyticsID string
}

func NewConfigFromPath(path string, fs boshsys.FileSystem) (Config, error) {
	var config Config

	bytes, err := fs.ReadFile(path)
	if err != nil {
		return config, bosherr.WrapError(err, "Reading config %s", path)
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, bosherr.WrapError(err, "Unmarshalling config")
	}

	return config, nil
}

func (c Config) Validate() error {
	for i, match := range c.ChecksumPrivs {
		err := match.Validate()
		if err != nil {
			return bosherr.WrapError(err, "Validating ChecksumPrivs[%d]", i)
		}
	}

	return nil
}
