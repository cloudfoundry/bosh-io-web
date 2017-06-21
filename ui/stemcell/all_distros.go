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
			Name: "Windows 2012R2",

			OSMatches: []StemcellOSMatch{
				{OSName: "windows", OSVersion: "2012R2"},
			},

			Sort: 2,
		},

		Distro{
			Name: "CentOS 7.x",

			OSMatches: []StemcellOSMatch{
				{OSName: "centos", OSVersion: "7"},
			},

			Sort: 3,
		},
	}
)
