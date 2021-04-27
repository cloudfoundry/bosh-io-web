package stemsrepo

import (
	semiver "github.com/cppforlife/go-semi-semantic/version"

	bhnotesrepo "github.com/bosh-io/web/stemcell/notesrepo"
)

type Stemcell interface {
	Name() string

	Version() semiver.Version
	UpdatedAt() string

	Size() uint64
	MD5() string
	SHA1() string // could be empty
	SHA256() string

	InfName() string    // e.g. aws
	HvName() string     // e.g. kvm
	DiskFormat() string // e.g. raw

	OSName() string    // e.g. Ubuntu
	OSVersion() string // e.g. Trusty

	IsLight() bool
	IsForChina() bool

	URL() string

	Notes() (bhnotesrepo.NoteRec, bool, error)
}

type StemcellsRepository interface {
	FindAll(string) ([]Stemcell, error)
}
