package stemsrepo

import (
	semiver "github.com/cppforlife/go-semi-semantic/version"
)

type Stemcell interface {
	Name() string

	Version() semiver.Version
	UpdatedAt() string

	Size() uint64
	MD5() string

	InfName() string // e.g. aws
	HvName() string  // e.g. kvm

	OSName() string    // e.g. Ubuntu
	OSVersion() string // e.g. Trusty

	IsLight() bool
	IsDeprecated() bool

	URL() string
}

type StemcellsRepository interface {
	FindAll() ([]Stemcell, error)
}
