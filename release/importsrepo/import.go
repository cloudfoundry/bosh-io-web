package importsrepo

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	semiver "github.com/cppforlife/go-semi-semantic/version"
)

type ImportRec struct {
	RelSource string
	Version   string
}

func (r ImportRec) Validate() error {
	if len(r.RelSource) == 0 {
		return bosherr.Error("Expected release source to be non-empty")
	}

	if len(r.Version) == 0 {
		return bosherr.Error("Expected version to be non-empty")
	}

	_, err := semiver.NewVersionFromString(r.Version)
	if err != nil {
		return bosherr.WrapError(err, "Expected version to be a valid version")
	}

	return nil
}
