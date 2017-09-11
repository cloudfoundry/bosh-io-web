package notesrepo

type NotesRepository interface {
	Find(source, version string) (NoteRec, bool, error)
}
