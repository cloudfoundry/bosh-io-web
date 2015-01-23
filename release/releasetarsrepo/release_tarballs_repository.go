package releasetarsrepo

import (
	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
)

type ReleaseTarballsRepository interface {
	Find(bhrelsrepo.ReleaseVersionRec) (ReleaseTarballRec, bool, error)
	Save(bhrelsrepo.ReleaseVersionRec, ReleaseTarballRec) error
}
