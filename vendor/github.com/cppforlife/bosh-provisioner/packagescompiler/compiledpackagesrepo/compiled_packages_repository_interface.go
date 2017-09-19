package compiledpackagesrepo

import (
	bprel "github.com/cppforlife/bosh-provisioner/release"
)

type CompiledPackageRecord struct {
	BlobID string
	SHA1   string
}

// CompiledPackagesRepository maintains list of compiled packages as blobs
// todo account for stemcell
type CompiledPackagesRepository interface {
	Find(bprel.Package) (CompiledPackageRecord, bool, error)
	Save(bprel.Package, CompiledPackageRecord) error
}
