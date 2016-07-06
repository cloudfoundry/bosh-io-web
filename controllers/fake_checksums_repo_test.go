package controllers_test

import (
	bhchecks "github.com/cppforlife/bosh-hub/checksumsrepo"
)

type FakeChecksumsRepository struct {
	SavedKey    string
	SavedRecord bhchecks.ChecksumRec
}

func (r *FakeChecksumsRepository) Find(key string) (bhchecks.ChecksumRec, error) {
	return bhchecks.ChecksumRec{}, nil
}

func (r *FakeChecksumsRepository) Save(key string, rec bhchecks.ChecksumRec) error {
	r.SavedKey = key
	r.SavedRecord = rec
	return nil
}
