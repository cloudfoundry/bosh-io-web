package stemcell

var (
	supportedDistros = []Distro{
		Distro{
			Name:      "Ubuntu Trusty",
			OSMatches: []StemcellOSMatch{{OSName: "ubuntu", OSVersion: "trusty"}},
		},

		Distro{
			Name:      "CentOS 6.5",
			OSMatches: []StemcellOSMatch{{OSName: "centos", OSVersion: ""}},
		},

		Distro{
			Name: "Ubuntu Lucid",
			OSMatches: []StemcellOSMatch{
				{OSName: "ubuntu", OSVersion: "lucid"},
				{OSName: "ubuntu", OSVersion: ""},
			},
		},
	}

	unknownDistro = Distro{Name: "Unknown", MatchAll: true}
)
