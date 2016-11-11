package noteimporter

type NoteImporter interface {
	Import() error
}
