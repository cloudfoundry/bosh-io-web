package releasesrepo

import (
	bprel "github.com/cppforlife/bosh-provisioner/release"
)

type ReleasesRepository interface {
	// todo return sources with ReleaseVersionRec
	ListCurated() ([]ReleaseVersionRec, error)
	ListAll() ([]Source, error)

	FindAll(source string) ([]ReleaseVersionRec, bool, error)
	FindLatest(source string) (ReleaseVersionRec, bool, error)
	Find(source, version string) (ReleaseVersionRec, error)

	Add(ReleaseVersionRec) error
	Contains(ReleaseVersionRec) (bool, error)
}

type ReleaseVersionsRepository interface {
	Find(ReleaseVersionRec) (bprel.Release, bool, error)
	Save(ReleaseVersionRec, bprel.Release) error
}

type avatarsResolver interface {
	Resolve(string) string
}
