package releasetarsrepo

type ReleaseTarballsRepository interface {
	Find(source, version string) (ReleaseTarballRec, bool, error)
	Save(source, version string, relTarRec ReleaseTarballRec) error
	// todo figure out source/version vs relVerRec
}
