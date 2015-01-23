package stemcell

type Distro struct {
	Name string // e.g. 'Ubuntu Trusty'

	OSMatches []StemcellOSMatch

	MatchAll bool
}

type StemcellOSMatch struct {
	OSName    string // e.g. ubuntu
	OSVersion string // e.g. trusty, ''
}

func (d Distro) Matches(s Stemcell) bool {
	if d.MatchAll {
		return true
	}

	for _, m := range d.OSMatches {
		if m.Matches(s) {
			return true
		}
	}

	return false
}

func (m StemcellOSMatch) Matches(s Stemcell) bool {
	return s.OSName == m.OSName && s.OSVersion == m.OSVersion
}
