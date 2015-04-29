package stemcell

import (
	"encoding/json"
	"fmt"

	semiver "github.com/cppforlife/go-semi-semantic/version"
	humanize "github.com/dustin/go-humanize"

	bhstemsrepo "github.com/cppforlife/bosh-hub/stemcell/stemsrepo"
)

var (
	prettyInfNames = map[string]string{
		"aws":       "AWS",
		"openstack": "OpenStack",
		"vsphere":   "vSphere",
		"vcloud":    "vCloud",
		"warden":    "BOSH Lite", // todo warden and boshlite are flipped
	}

	prettyHvNames = map[string]string{
		"xen":      "Xen",
		"xen-hvm":  "Xen-HVM",
		"esxi":     "ESXi",
		"kvm":      "KVM",
		"boshlite": "Warden",
	}
)

type Stemcell struct {
	ManifestName string
	FriendlyName string
	Version      semiver.Version

	OSName    string
	OSVersion string

	RegularSource *StemcellSource
	LightSource   *StemcellSource
}

type StemcellSource struct {
	URL  string `json:"url"`
	Size uint64 `json:"size"`
	MD5  string `json:"md5"`

	UpdatedAt string `json:"-"`
}

type stemcellAPIRecord struct {
	Name    string `json:"name"`
	Version string `json:"version"`

	// Use StemcellSource for now for convenience
	Regular *StemcellSource `json:"regular,omitempty"`
	Light   *StemcellSource `json:"light,omitempty"`
}

type StemcellFriendlyNameSorting []Stemcell

type StemcellVersionSorting []Stemcell

func NewStemcell(s bhstemsrepo.Stemcell) Stemcell {
	infName, ok := prettyInfNames[s.InfName()]
	if !ok {
		infName = s.InfName()
	}

	hvName, ok := prettyHvNames[s.HvName()]
	if !ok {
		hvName = s.HvName()
	}

	optionalDiskFormat := ""
	if len(s.DiskFormat()) > 0 {
		optionalDiskFormat = fmt.Sprintf(" (%s)", s.DiskFormat())
	}

	stemcell := &Stemcell{
		ManifestName: s.Name(),
		FriendlyName: fmt.Sprintf("%s %s%s", infName, hvName, optionalDiskFormat),
		Version:      s.Version(),

		OSName:    s.OSName(),
		OSVersion: s.OSVersion(),
	}

	stemcell.AddAsSource(s)

	return *stemcell
}

func (s *Stemcell) AddAsSource(s_ bhstemsrepo.Stemcell) {
	source := &StemcellSource{
		URL:  s_.URL(),
		Size: s_.Size(),
		MD5:  s_.MD5(),

		UpdatedAt: s_.UpdatedAt(),
	}

	if s_.IsLight() {
		s.LightSource = source
	} else {
		s.RegularSource = source
	}
}

func (s Stemcell) UserVisibleDownloadURL() string {
	// todo make domain configurable
	return fmt.Sprintf("https://bosh.io/d/stemcells/%s?v=%s", s.ManifestName, s.Version)
}

func (s Stemcell) UserVisibleLatestDownloadURL() string {
	// todo make domain configurable
	return fmt.Sprintf("https://bosh.io/d/stemcells/%s", s.ManifestName)
}

func (s Stemcell) ActualDownloadURL() string {
	// Prefer light stemcells
	if s.LightSource != nil {
		return s.LightSource.URL
	}

	return s.RegularSource.URL
}

func (s Stemcell) MarshalJSON() ([]byte, error) {
	record := stemcellAPIRecord{
		Name:    s.ManifestName,
		Version: s.Version.AsString(),
		Regular: s.RegularSource,
		Light:   s.LightSource,
	}

	return json.Marshal(record)
}

func (s Stemcell) AllVersionsURL() string { return fmt.Sprintf("/stemcells/%s", s.ManifestName) }

func (s StemcellSource) FormattedSize() string { return humanize.Bytes(s.Size) }

func (s StemcellFriendlyNameSorting) Len() int           { return len(s) }
func (s StemcellFriendlyNameSorting) Less(i, j int) bool { return s[i].ManifestName < s[j].ManifestName }
func (s StemcellFriendlyNameSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s StemcellVersionSorting) Len() int           { return len(s) }
func (s StemcellVersionSorting) Less(i, j int) bool { return s[i].Version.IsLt(s[j].Version) }
func (s StemcellVersionSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
