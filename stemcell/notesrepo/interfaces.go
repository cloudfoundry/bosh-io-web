package notesrepo

type NotesRepository interface {
	Find(version string) (NoteRec, bool, error)
	Save(version string, noteRec NoteRec) error
}
