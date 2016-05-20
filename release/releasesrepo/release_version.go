package releasesrepo

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	semiver "github.com/cppforlife/go-semi-semantic/version"

	bhnotesrepo "github.com/cppforlife/bosh-hub/release/notesrepo"
	bhreltarsrepo "github.com/cppforlife/bosh-hub/release/releasetarsrepo"
)

type ReleaseVersionRec struct {
	notesRepo       bhnotesrepo.NotesRepository
	releaseTarsRepo bhreltarsrepo.ReleaseTarballsRepository
	avatarsResolver avatarsResolver

	Source     string
	VersionRaw string

	avatarURL string
}

type ReleaseVersionRecSorting []ReleaseVersionRec

func (r ReleaseVersionRec) String() string {
	return fmt.Sprintf("%s %s", r.Source, r.VersionRaw)
}

// AsSource returns Source object based on the Source string
// todo refactor to remove Source string
func (r ReleaseVersionRec) AsSource() Source {
	return Source{
		Full:            r.Source,
		avatarsResolver: r.avatarsResolver,
	}
}

// Version returns parsed version
// todo non-memoized lazy loading is expensive
func (r ReleaseVersionRec) Version() semiver.Version {
	// Validate should not allow invalid version to be saved
	ver, err := semiver.NewVersionFromString(r.VersionRaw)
	if err != nil {
		panic(fmt.Sprintf("Version '%s' is not valid: %s", r.VersionRaw, err))
	}

	return ver
}

func (r ReleaseVersionRec) AvatarURL() string {
	return r.AsSource().AvatarURL()
}

func (r ReleaseVersionRec) Tarball() (bhreltarsrepo.ReleaseTarballRec, error) {
	return r.releaseTarsRepo.Find(r.Source, r.VersionRaw)
}

func (r ReleaseVersionRec) SetTarball(relTarRec bhreltarsrepo.ReleaseTarballRec) error {
	return r.releaseTarsRepo.Save(r.Source, r.VersionRaw, relTarRec)
}

func (r ReleaseVersionRec) Notes() (bhnotesrepo.NoteRec, bool, error) {
	return r.notesRepo.Find(r.Source, r.VersionRaw)
}

func (r ReleaseVersionRec) SetNotes(noteRec bhnotesrepo.NoteRec) error {
	return r.notesRepo.Save(r.Source, r.VersionRaw, noteRec)
}

func (r ReleaseVersionRec) Equals(other ReleaseVersionRec) bool {
	return r.Source == other.Source && r.VersionRaw == other.VersionRaw
}

func (r ReleaseVersionRec) Validate() error {
	if len(r.Source) == 0 {
		return bosherr.Error("Expected source to be non-empty")
	}

	if len(r.VersionRaw) == 0 {
		return bosherr.Error("Expected version to be non-empty")
	}

	_, err := semiver.NewVersionFromString(r.VersionRaw)
	if err != nil {
		return bosherr.WrapError(err, "Expected version to be a valid version")
	}

	return nil
}

func (s ReleaseVersionRecSorting) Len() int           { return len(s) }
func (s ReleaseVersionRecSorting) Less(i, j int) bool { return s[i].Version().IsLt(s[j].Version()) }
func (s ReleaseVersionRecSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
