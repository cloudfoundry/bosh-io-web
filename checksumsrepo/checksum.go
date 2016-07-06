package checksumsrepo

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
)

type ChecksumRec struct {
	SHA1 string
}

func (r ChecksumRec) Validate() error {
	if len(r.SHA1) != 40 {
		return bosherr.New("Expected SHA1 to be 40 chars long")
	}

	return nil
}
