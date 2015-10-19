package releasesrepo

import (
	"sort"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	bpindex "github.com/cppforlife/bosh-provisioner/index"

	bhnotesrepo "github.com/cppforlife/bosh-hub/release/notesrepo"
	bhreltarsrepo "github.com/cppforlife/bosh-hub/release/releasetarsrepo"
)

type CRRepository struct {
	predefinedSources []string
	avatarsResolver   avatarsResolver

	index           bpindex.Index
	releaseTarsRepo bhreltarsrepo.ReleaseTarballsRepository
	notesRepo       bhnotesrepo.NotesRepository

	logTag string
	logger boshlog.Logger
}

type sourceToRelVerRecKey struct {
	Source string
}

type predefinedAvatarsResolver struct {
	locationToURL map[string]string
}

func (r predefinedAvatarsResolver) Resolve(location string) string {
	return r.locationToURL[location]
}

func NewConcreteReleasesRepository(
	predefinedSources []string,
	predefinedAvatars map[string]string,
	index bpindex.Index,
	notesRepo bhnotesrepo.NotesRepository,
	releaseTarsRepo bhreltarsrepo.ReleaseTarballsRepository,
	logger boshlog.Logger,
) CRRepository {
	return CRRepository{
		predefinedSources: predefinedSources,
		avatarsResolver:   predefinedAvatarsResolver{predefinedAvatars},

		index:           index,
		releaseTarsRepo: releaseTarsRepo,
		notesRepo:       notesRepo,

		logTag: "CRRepository",
		logger: logger,
	}
}

func (r CRRepository) ListCurated() ([]ReleaseVersionRec, error) {
	var relVerRecs []ReleaseVersionRec

	for _, source := range r.predefinedSources {
		recs, err := r.FindAll(source)
		if err == nil {
			relVerRecs = append(relVerRecs, recs...)
		}
	}

	for i, _ := range relVerRecs {
		relVerRecs[i].avatarsResolver = r.avatarsResolver
	}

	return relVerRecs, nil
}

func (r CRRepository) ListAll() ([]Source, error) {
	var sourceKeys []sourceToRelVerRecKey
	var sources []Source

	err := r.index.ListKeys(&sourceKeys)
	if err != nil {
		return sources, bosherr.WrapError(err, "Listing all releases")
	}

	for _, sourceKey := range sourceKeys {
		sources = append(sources, Source{Full: sourceKey.Source})
	}

	for i, _ := range sources {
		sources[i].avatarsResolver = r.avatarsResolver
	}

	return sources, nil
}

func (r CRRepository) FindAll(source string) ([]ReleaseVersionRec, error) {
	var relVerRecs []ReleaseVersionRec

	if len(source) == 0 {
		return relVerRecs, bosherr.New("Expected source to be non-empty")
	}

	err := r.index.Find(sourceToRelVerRecKey{source}, &relVerRecs)
	if err != nil {
		return relVerRecs, bosherr.WrapError(err, "Finding release version records")
	}

	for i, _ := range relVerRecs {
		// Make sure to change real relVerRec
		relVerRecs[i].notesRepo = r.notesRepo
		relVerRecs[i].releaseTarsRepo = r.releaseTarsRepo
	}

	return relVerRecs, nil
}

func (r CRRepository) FindLatest(source string) (ReleaseVersionRec, error) {
	var relVerRec ReleaseVersionRec

	if len(source) == 0 {
		return relVerRec, bosherr.New("Expected source to be non-empty")
	}

	relVerRecs, err := r.FindAll(source)
	if err != nil {
		return relVerRec, bosherr.WrapError(err, "Finding release version records")
	}

	if len(relVerRecs) == 0 {
		return relVerRec, bosherr.New("Expected to find at least one release version record")
	}

	sort.Sort(ReleaseVersionRecSorting(relVerRecs))

	relVerRec = relVerRecs[len(relVerRecs)-1]

	return relVerRec, nil
}

func (r CRRepository) Find(source, version string) (ReleaseVersionRec, error) {
	relVerRec := ReleaseVersionRec{
		notesRepo:       r.notesRepo,
		releaseTarsRepo: r.releaseTarsRepo,

		Source:     source,
		VersionRaw: version,
	}

	err := relVerRec.Validate()
	if err != nil {
		return relVerRec, bosherr.WrapError(err, "Validating release version record")
	}

	return relVerRec, nil
}

func (r CRRepository) Add(relVerRec ReleaseVersionRec) error {
	r.logger.Debug(r.logTag, "Adding release '%v'", relVerRec)

	err := relVerRec.Validate()
	if err != nil {
		return bosherr.WrapError(err, "Validating release version record")
	}

	// Try to add release version record up to 5 times
	for i := 0; i < 5; i++ {
		tryAgain, err := r.tryToAddRelVerRec(relVerRec)
		if err != nil {
			return bosherr.WrapError(err, "Trying to add release version record")
		}

		if !tryAgain {
			return nil
		}
	}

	return bosherr.New("Failed to add release verison record several times")
}

func (r CRRepository) Contains(relVerRec ReleaseVersionRec) (bool, error) {
	r.logger.Debug(r.logTag, "Checking release '%v' existence", relVerRec)

	err := relVerRec.Validate()
	if err != nil {
		return false, bosherr.WrapError(err, "Validating release version record")
	}

	var relVerRecs []ReleaseVersionRec

	err = r.index.Find(sourceToRelVerRecKey{relVerRec.Source}, &relVerRecs)
	if err != nil {
		if err == bpindex.ErrNotFound {
			return false, nil
		}

		return false, bosherr.WrapError(err, "Finding release version records")
	}

	r.logger.Debug(r.logTag,
		"Found release versions '%v' for source '%s'", relVerRecs, relVerRec.Source)

	// Make sure version is not already in the list
	for _, existingRelVerRec := range relVerRecs {
		if existingRelVerRec.Equals(relVerRec) {
			return true, nil
		}
	}

	return false, nil
}

// tryToAddRelVerRec
// todo refactor to be generic add-retry key
func (r CRRepository) tryToAddRelVerRec(relVerRec ReleaseVersionRec) (bool, error) {
	var relVerRecs []ReleaseVersionRec

	// Find release for that source and disallow any modifications until it gets saved
	lockedRec, err := r.index.FindLocked(sourceToRelVerRecKey{relVerRec.Source}, &relVerRecs)

	// Release locked record no matter what
	// todo weird error handling
	defer lockedRec.Release()

	if err != nil {
		if err != bpindex.ErrNotFound {
			return false, bosherr.WrapError(err, "Finding release version records")
		}
	}

	// Make sure version is not already in the list
	for _, existingRelVerRec := range relVerRecs {
		if existingRelVerRec.Equals(relVerRec) {
			return false, nil
		}
	}

	// Add new release version to the release
	relVerRecs = append(relVerRecs, relVerRec)

	err = lockedRec.Save(relVerRecs)
	if err != nil {
		if err == bpindex.ErrChanged {
			// Try adding release version for the release
			return true, nil
		}

		return false, bosherr.WrapError(err, "Adding release version records")
	}

	// Successfully added release version for the release
	return false, nil
}
