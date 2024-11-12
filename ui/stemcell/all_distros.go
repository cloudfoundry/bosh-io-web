package stemcell

var (
	ubuntuNobleDistro = Distro{
		NameName:        "ubuntu-noble",
		Name:            "Ubuntu Noble",
		NoGoAgentSuffix: true,

		OSMatches: []StemcellOSMatch{
			{OSName: "ubuntu", OSVersion: "noble"},
		},

		SupportedInfrastructures: nobleInfrastructures,

		Sort: 1,
	}
	ubuntuJammyDistro = Distro{
		NameName: "ubuntu-jammy",
		Name:     "Ubuntu Jammy",

		OSMatches: []StemcellOSMatch{
			{OSName: "ubuntu", OSVersion: "jammy"},
		},

		SupportedInfrastructures: allInfrastructures,

		Sort: 2,
	}
	ubuntuBionicDistro = Distro{
		NameName: "ubuntu-bionic",
		Name:     "Ubuntu Bionic",

		OSMatches: []StemcellOSMatch{
			{OSName: "ubuntu", OSVersion: "bionic"},
		},

		SupportedInfrastructures: allInfrastructures,

		Sort: 3,
	}
	windows2019Distro = Distro{
		NameName: "windows2019",
		Name:     "Windows 2019",

		OSMatches: []StemcellOSMatch{
			{OSName: "windows", OSVersion: "2019"},
		},

		SupportedInfrastructures: Infrastructures{
			awsInfrastructure,
			googleInfrastructure,
			azureInfrastructure,
		},

		Sort: 4,
	}
)

var (
	allDistros = []Distro{
		ubuntuNobleDistro,
		ubuntuJammyDistro,
		ubuntuBionicDistro,
		windows2019Distro,
	}
)
