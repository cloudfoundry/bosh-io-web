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
			Name: "CentOS 6.x",

			OSMatches: []StemcellOSMatch{
				{OSName: "centos", OSVersion: ""},
			},

			Sort: 2,
		},
	}
)
