package importerrsrepo

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	bhimpsrepo "github.com/cppforlife/bosh-hub/release/importsrepo"
)

type ImportErrRec struct {
	ImportRec bhimpsrepo.ImportRec

	Err string
}

func (r ImportErrRec) Validate() error {
	err := r.ImportRec.Validate()
	if err != nil {
		return bosherr.WrapError(err, "Validating import")
	}

	if len(r.Err) == 0 {
		return bosherr.New("Expected error to be non-empty")
	}

	return nil
}
