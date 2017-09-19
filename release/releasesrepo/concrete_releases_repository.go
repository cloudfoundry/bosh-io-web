package releasesrepo

import (
	"encoding/json"
	"path/filepath"
	"sort"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"gopkg.in/yaml.v2"

	bhnotesrepo "github.com/bosh-io/web/release/notesrepo"
	bhreltarsrepo "github.com/bosh-io/web/release/releasetarsrepo"
)

type CRRepository struct {
	avatarsResolver  avatarsResolver
	releasesDir      string
	releasesIndexDir string

	releaseTarsRepo bhreltarsrepo.ReleaseTarballsRepository
	notesRepo       bhnotesrepo.NotesRepository

	fs boshsys.FileSystem

	logTag string
	logger boshlog.Logger
}

func NewConcreteReleasesRepository(
	releasesDir string,
	releasesIndexDir string,
	notesRepo bhnotesrepo.NotesRepository,
	releaseTarsRepo bhreltarsrepo.ReleaseTarballsRepository,
	fs boshsys.FileSystem,
	logger boshlog.Logger,
) CRRepository {
	return CRRepository{
		avatarsResolver:  predefinedAvatarsResolver{releasesDir, fs},
		releasesDir:      releasesDir,
		releasesIndexDir: releasesIndexDir,

		releaseTarsRepo: releaseTarsRepo,
		notesRepo:       notesRepo,

		fs: fs,

		logTag: "CRRepository",
		logger: logger,
	}
}

func (r CRRepository) ListCurated() ([]ReleaseVersionRec, error) {
	var relVerRecs []ReleaseVersionRec

	defs, err := r.defs()
	if err != nil {
		return nil, err
	}

	for _, def := range defs {
		if !def.OnHomePage() {
			continue
		}

		recs, err := r.FindAll(def.TrimmedURL())
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
	var sources []Source

	defs, err := r.defs()
	if err != nil {
		return nil, err
	}

	for _, def := range defs {
		sources = append(sources, Source{Full: def.TrimmedURL()})
	}

	for i, _ := range sources {
		sources[i].avatarsResolver = r.avatarsResolver
	}

	return sources, nil
}

type releaseV1YAML struct {
	Version string `yaml:"Version"`
}

func (r CRRepository) FindAll(source string) ([]ReleaseVersionRec, error) {
	var relVerRecs []ReleaseVersionRec

	if len(source) == 0 {
		return relVerRecs, bosherr.Error("Expected source to be non-empty")
	}

	foundPaths, err := r.fs.Glob(filepath.Join(r.releasesIndexDir, source, "*", "release.v1.yml"))
	if err != nil {
		return nil, bosherr.WrapError(err, "Globbing release versions")
	}

	for _, path := range foundPaths {
		contents, err := r.fs.ReadFileString(path)
		if err != nil {
			return nil, bosherr.WrapError(err, "Reading release version")
		}

		var v1 releaseV1YAML

		err = json.Unmarshal([]byte(contents), &v1)
		if err != nil {
			return nil, bosherr.WrapError(err, "Unmarshaling release version")
		}

		relVerRec, err := r.Find(source, v1.Version)
		if err != nil {
			return nil, bosherr.WrapError(err, "Building release version")
		}

		relVerRecs = append(relVerRecs, relVerRec)
	}

	return relVerRecs, nil
}

func (r CRRepository) FindLatest(source string) (ReleaseVersionRec, error) {
	var relVerRec ReleaseVersionRec

	if len(source) == 0 {
		return relVerRec, bosherr.Error("Expected source to be non-empty")
	}

	relVerRecs, err := r.FindAll(source)
	if err != nil {
		return relVerRec, bosherr.WrapError(err, "Finding release version records")
	}

	if len(relVerRecs) == 0 {
		return relVerRec, bosherr.Error("Expected to find at least one release version record")
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

type releaseDefYAML struct {
	URL        string
	Categories []string
}

func (d releaseDefYAML) OnHomePage() bool {
	for _, cat := range d.Categories {
		if cat == "homepage" {
			return true
		}
	}
	return false
}

func (d releaseDefYAML) TrimmedURL() string {
	return strings.TrimPrefix(d.URL, "https://")
}

func (r CRRepository) defs() ([]releaseDefYAML, error) {
	contents, err := r.fs.ReadFileString(filepath.Join(r.releasesDir, "index.yml"))
	if err != nil {
		return nil, bosherr.WrapError(err, "Reading releases")
	}

	var defs []releaseDefYAML

	err = yaml.Unmarshal([]byte(contents), &defs)
	if err != nil {
		return nil, bosherr.WrapError(err, "Unmarshaling releases")
	}

	return defs, nil
}
