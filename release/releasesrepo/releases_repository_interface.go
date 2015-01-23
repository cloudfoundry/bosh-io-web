package releasesrepo

import (
	bprel "github.com/cppforlife/bosh-provisioner/release"
)

type Source string

type ReleasesRepository interface {
	// todo return sources with ReleaseVersionRec
	ListCurated() ([]ReleaseVersionRec, error)
	ListAll() ([]Source, error)

	FindAll(string) ([]ReleaseVersionRec, bool, error)
	FindLatest(string) (ReleaseVersionRec, bool, error)

	Add(ReleaseVersionRec) error
	Contains(ReleaseVersionRec) (bool, error)
}

type ReleaseVersionsRepository interface {
	Find(ReleaseVersionRec) (bprel.Release, bool, error)
	Save(ReleaseVersionRec, bprel.Release) error
}
