package notesrepo

type NotesRepository interface {
	Find(source, version string) (NoteRec, bool, error)
	Save(source, version string, noteRec NoteRec) error
	// todo figure out source/version vs relVerRec
}
