package checksumsrepo

type ChecksumsRepository interface {
	Find(key string) (ChecksumRec, error)
	Save(key string, rec ChecksumRec) error
}
