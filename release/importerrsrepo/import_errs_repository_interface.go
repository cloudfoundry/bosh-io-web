package importerrsrepo

import (
	bhimpsrepo "github.com/cppforlife/bosh-hub/release/importsrepo"
)

type ImportErrsRepository interface {
	ListAll() ([]ImportErrRec, error)
	Add(ImportErrRec) error
	Remove(bhimpsrepo.ImportRec) error
	Contains(bhimpsrepo.ImportRec) (bool, error)
}
