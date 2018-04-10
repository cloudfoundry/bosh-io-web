package main

import (
	"encoding/json"

	yaml "gopkg.in/yaml.v2"

	"github.com/bosh-io/web/controllers"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type Config struct {
	Repos ReposOptions

	APIKey string

	Analytics AnalyticsConfig
}

type AnalyticsConfig struct {
	GoogleAnalyticsID string
}

func NewConfigFromPath(path string, fs boshsys.FileSystem) (Config, error) {
	var config Config

	bytes, err := fs.ReadFile(path)
	if err != nil {
		return config, bosherr.WrapErrorf(err, "Reading config %s", path)
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, bosherr.WrapError(err, "Unmarshalling config")
	}

	return config, nil
}

func LoadRedirects(fs boshsys.FileSystem) (controllers.RedirectsConfig, error) {
	var redirects controllers.RedirectsConfig

	bytes, err := fs.ReadFile("conf/redirects.yml")
	if err != nil {
		return redirects, bosherr.WrapErrorf(err, "Reading redirects")
	}

	err = yaml.Unmarshal(bytes, &redirects)
	if err != nil {
		return redirects, bosherr.WrapErrorf(err, "Unmarshalling redirects")
	}

	return redirects, nil
}
