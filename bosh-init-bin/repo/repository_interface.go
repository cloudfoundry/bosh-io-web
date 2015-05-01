package repo

import (
	semiver "github.com/cppforlife/go-semi-semantic/version"
)

type BinaryGroup struct {
	Version  semiver.Version
	Binaries []Binary
}

type Binary interface {
	Description() string

	Version() semiver.Version
	UpdatedAt() string

	Size() uint64
	MD5() string

	Platform() string // e.g. linux, darwin
	Arch() string     // e.g. amd64

	URL() string
}

type Repository interface {
	FindLatest() ([]BinaryGroup, error)
}
