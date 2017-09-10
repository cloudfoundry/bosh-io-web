package releasetarsrepo

type ReleaseTarballsRepository interface {
	Find(source, version string) (ReleaseTarballRec, error)
}
