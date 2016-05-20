package watchersrepo

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	semiver "github.com/cppforlife/go-semi-semantic/version"
)

type WatcherRec struct {
	RelSource string

	// Versions below minimum version will not be imported
	MinVersionRaw string
}

func (r WatcherRec) MinVersion() semiver.Version {
	// Validate should not allow invalid version to be saved
	ver, err := semiver.NewVersionFromString(r.MinVersionRaw)
	if err != nil {
		panic(fmt.Sprintf("Version '%s' is not valid: %s", r.MinVersionRaw, err))
	}

	return ver
}

func (r WatcherRec) Validate() error {
	if len(r.RelSource) == 0 {
		return bosherr.New("Expected release source to be non-empty")
	}

	if len(r.MinVersionRaw) == 0 {
		return bosherr.New("Expected min version to be non-empty")
	}

	_, err := semiver.NewVersionFromString(r.MinVersionRaw)
	if err != nil {
		return bosherr.WrapError(err, "Expected min version to be a valid version")
	}

	return nil
}
