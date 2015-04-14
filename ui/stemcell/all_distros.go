package stemcell

var (
	allDistros = []Distro{
		Distro{
			Name: "Ubuntu Trusty",

			OSMatches: []StemcellOSMatch{
				{OSName: "ubuntu", OSVersion: "trusty"},
			},

			Sort: 1,
		},

		Distro{
			Name: "CentOS 7.x",

			OSMatches: []StemcellOSMatch{
				{OSName: "centos", OSVersion: "7"},
			},

			Sort: 2,
		},

		Distro{
			Name: "CentOS 6.x",

			OSMatches: []StemcellOSMatch{
				{OSName: "centos", OSVersion: ""},
			},

			Deprecated: true,

			Sort: 3,
		},
	}
)
