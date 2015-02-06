package stemcell

import (
	bhstemsrepo "github.com/cppforlife/bosh-hub/stemcell/stemsrepo"
)

type Distro struct {
	Name string // e.g. 'Ubuntu Trusty'
	Sort uint8  // smaller == more important

	OSMatches []StemcellOSMatch
}

type StemcellOSMatch struct {
	OSName    string // e.g. ubuntu
	OSVersion string // e.g. trusty, ''
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
