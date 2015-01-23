package importsrepo

type ImportsRepository interface {
	ListAll() ([]ImportRec, error)
	Push(string, string) error
	Pull() (ImportRec, bool, error)
	Remove(ImportRec) error
}
