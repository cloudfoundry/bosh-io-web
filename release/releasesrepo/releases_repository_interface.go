package releasesrepo

import (
	bprel "github.com/bosh-dep-forks/bosh-provisioner/release"
)

type ReleasesRepository interface {
	// todo return sources with ReleaseVersionRec
	ListCurated() ([]ReleaseVersionRec, error)
	ListAll() ([]Source, error)

	FindAll(source string) ([]ReleaseVersionRec, error)
	FindLatest(source string) (ReleaseVersionRec, error)
	Find(source, version string) (ReleaseVersionRec, error)
}

type ReleaseVersionsRepository interface {
	Find(ReleaseVersionRec) (bprel.Release, error)
}

type avatarsResolver interface {
	Resolve(string) string
}
