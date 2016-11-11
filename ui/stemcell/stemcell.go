package stemcell

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"

	semiver "github.com/cppforlife/go-semi-semantic/version"
	humanize "github.com/dustin/go-humanize"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	bhstemsrepo "github.com/cppforlife/bosh-hub/stemcell/stemsrepo"
)

var (
	prettyInfNames = map[string]string{
		"aws":       "AWS",
		"openstack": "OpenStack",
		"vsphere":   "vSphere",
		"vcloud":    "vCloud",
		"azure":     "Azure",
		"google":    "Google",
		"softlayer": "SoftLayer",
		"warden":    "BOSH Lite", // todo warden and boshlite are flipped
	}

	prettyHvNames = map[string]string{
		"xen":      "Xen",
		"xen-hvm":  "Xen-HVM",
		"esxi":     "ESXi",
		"kvm":      "KVM",
		"hyperv":   "Hyper-V",
		"boshlite": "Warden",
	}
)

type Stemcell struct {
	stemRec bhstemsrepo.Stemcell

	ManifestName string
	Version      semiver.Version

	OSName    string
	OSVersion string

	RegularSource *StemcellSource
	LightSource   *StemcellSource

	// todo china stemcell will be consolidated into light stemcell at some point
	LightChinaSource *StemcellSource

	// memoized notes
	notesInMarkdown *[]byte
}

type StemcellSource struct {
	friendlyName string

	isLight    bool
	isForChina bool

	URL  string `json:"url"`
	Size uint64 `json:"size"`
	MD5  string `json:"md5"`
	SHA1 string `json:"sha1,omitempty"`

	UpdatedAt string `json:"-"`
}

type stemcellAPIRecord struct {
	Name    string `json:"name"`
	Version string `json:"version"`

	// Use StemcellSource for now for convenience
	Regular *StemcellSource `json:"regular,omitempty"`
	Light   *StemcellSource `json:"light,omitempty"`

	// todo china stemcell will be consolidated into light stemcell at some point
	LightChina *StemcellSource `json:"light_china,omitempty"`
}

type StemcellManifestNameSorting []Stemcell

type StemcellVersionSorting []Stemcell

func NewStemcell(s bhstemsrepo.Stemcell) Stemcell {
	stemcell := &Stemcell{
		stemRec: s,

		ManifestName: s.Name(),
		Version:      s.Version(),

		OSName:    s.OSName(),
		OSVersion: s.OSVersion(),
	}

	stemcell.AddAsSource(s)

	return *stemcell
}

func (s *Stemcell) AddAsSource(s_ bhstemsrepo.Stemcell) {
	infName, ok := prettyInfNames[s_.InfName()]
	if !ok {
		infName = s_.InfName()
	}

	hvName, ok := prettyHvNames[s_.HvName()]
	if !ok {
		hvName = s_.HvName()
	}

	optionalDiskFormat := ""
	if len(s_.DiskFormat()) > 0 {
		optionalDiskFormat = fmt.Sprintf(" (%s)", s_.DiskFormat())
	}

	optionalLight := ""
	if s_.IsLight() {
		if s_.IsForChina() {
			optionalLight = " Light China"
		} else {
			optionalLight = " Light"
		}
	}

	source := &StemcellSource{
		friendlyName: fmt.Sprintf("%s %s%s%s", infName, hvName, optionalDiskFormat, optionalLight),

		URL:  s_.URL(),
		Size: s_.Size(),
		MD5:  s_.MD5(),
		SHA1: s_.SHA1(),

		UpdatedAt: s_.UpdatedAt(),
	}

	if s_.IsLight() {
		source.isLight = true

		if s_.IsForChina() {
			source.isForChina = true
			s.LightChinaSource = source
		} else {
			s.LightSource = source
		}
	} else {
		s.RegularSource = source
	}
}

func (s *Stemcell) Sources() []*StemcellSource {
	sources := []*StemcellSource{}

	if s.RegularSource != nil {
		sources = append(sources, s.RegularSource)
	}

	if s.LightSource != nil {
		sources = append(sources, s.LightSource)
	}

	if s.LightChinaSource != nil {
		sources = append(sources, s.LightChinaSource)
	}

	return sources
}

func (s Stemcell) UserVisibleDownloadURL() string {
	return fmt.Sprintf("https://bosh.io/d/stemcells/%s?v=%s", s.ManifestName, s.Version)
}

func (s Stemcell) UserVisibleLatestDownloadURL() string {
	return fmt.Sprintf("https://bosh.io/d/stemcells/%s", s.ManifestName)
}

func (s Stemcell) ActualDownloadURL(preferLight bool, mustBeForChina bool) (string, error) {
	// todo remove china variation
	if mustBeForChina {
		if s.LightChinaSource != nil {
			return s.LightChinaSource.URL, nil
		}

		return "", errors.New("No light stemcell for China found")
	}

	if preferLight {
		if s.LightSource != nil {
			return s.LightSource.URL, nil
		}
	}

	return s.RegularSource.URL, nil
}

func (s Stemcell) SHA1() string {
	if s.LightSource != nil {
		return s.LightSource.SHA1
	}
	return s.RegularSource.SHA1
}

func (s *Stemcell) NotesInMarkdown() (template.HTML, error) {
	if s.notesInMarkdown == nil {
		// Do not care about found -> no UI indicator
		noteRec, _, err := s.stemRec.Notes()
		if err != nil {
			return template.HTML(""), err
		}

		unsafeMarkdown := blackfriday.MarkdownCommon([]byte(noteRec.Content))
		safeMarkdown := bluemonday.UGCPolicy().SanitizeBytes(unsafeMarkdown)

		s.notesInMarkdown = &safeMarkdown
	}

	// todo sanitized markdown
	return template.HTML(*s.notesInMarkdown), nil
}

func (s Stemcell) MarshalJSON() ([]byte, error) {
	record := stemcellAPIRecord{
		Name:       s.ManifestName,
		Version:    s.Version.AsString(),
		Regular:    s.RegularSource,
		Light:      s.LightSource,
		LightChina: s.LightChinaSource,
	}

	return json.Marshal(record)
}

func (s Stemcell) AllVersionsURL() string { return fmt.Sprintf("/stemcells/%s", s.ManifestName) }

func (s StemcellSource) UserVisibleDownloadURL() string { return s.URL }
func (s StemcellSource) FriendlyName() string           { return s.friendlyName }
func (s StemcellSource) FormattedSize() string          { return humanize.Bytes(s.Size) }

func (s StemcellManifestNameSorting) Len() int           { return len(s) }
func (s StemcellManifestNameSorting) Less(i, j int) bool { return s[i].ManifestName < s[j].ManifestName }
func (s StemcellManifestNameSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s StemcellVersionSorting) Len() int           { return len(s) }
func (s StemcellVersionSorting) Less(i, j int) bool { return s[i].Version.IsLt(s[j].Version) }
func (s StemcellVersionSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
