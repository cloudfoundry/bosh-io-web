package stemcell

import (
	"fmt"

	bhstemsrepo "github.com/bosh-io/web/stemcell/stemsrepo"
)

type Distro struct {
	NameName string // e.g. 'ubuntu-trusty' TODO rename this to Name, Name to Title/Label
	Name     string // e.g. 'Ubuntu Trusty'
	Sort     uint8  // smaller == more important

	Deprecated bool

	OSMatches                []StemcellOSMatch
	SupportedInfrastructures Infrastructures
}

type StemcellOSMatch struct {
	OSName    string // e.g. ubuntu
	OSVersion string // e.g. trusty, ''
}

func (m StemcellOSMatch) Name() string {
	format := "%s-%s"
	if m.OSName == "windows" {
		format = "%s%s"
	}

	return fmt.Sprintf(format, m.OSName, m.OSVersion)
}

func (d Distro) IsVisible(includeDeprecated bool) bool {
	return !d.Deprecated || (d.Deprecated && includeDeprecated)
}

func (d Distro) Matches(s bhstemsrepo.Stemcell) bool {
	for _, m := range d.OSMatches {
		if m.Matches(s) {
			return true
		}
	}

	return false
}

func (m StemcellOSMatch) Matches(s bhstemsrepo.Stemcell) bool {
	return s.OSName() == m.OSName && s.OSVersion() == m.OSVersion
}
