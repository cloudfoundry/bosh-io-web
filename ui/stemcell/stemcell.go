package stemcell

import (
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
	Name    string
	Version semiver.Version

	Size uint64
	MD5  string

	OSName    string
	OSVersion string

	IsLight      bool
	IsDeprecated bool

	URL string
}

type StemcellSorting []Stemcell

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

	return Stemcell{
		Name:    fmt.Sprintf("%s %s%s", infName, hvName, optionalDiskFormat),
		Version: s.Version(),

		Size: s.Size(),
		MD5:  s.MD5(),

		OSName:    s.OSName(),
		OSVersion: s.OSVersion(),

		IsLight:      s.IsLight(),
		IsDeprecated: s.IsDeprecated(),

		URL: s.URL(),
	}
}

func (s Stemcell) FormattedSize() string { return humanize.Bytes(s.Size) }

func (s Stemcell) HasURL() bool { return len(s.URL) > 0 }

func (s StemcellSorting) Len() int           { return len(s) }
func (s StemcellSorting) Less(i, j int) bool { return s[i].Name < s[j].Name }
func (s StemcellSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
