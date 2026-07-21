package stemcell

var (
	ubuntuResoluteDistro = Distro{
		NameName:        "ubuntu-resolute",
		Name:            "Ubuntu Resolute",
		NoGoAgentSuffix: true,

		OSMatches: []StemcellOSMatch{
			{OSName: "ubuntu", OSVersion: "resolute"},
		},

		SupportedInfrastructures: resoluteInfrastructures,

		Sort: 1,
	}
	ubuntuNobleDistro = Distro{
		NameName:        "ubuntu-noble",
		Name:            "Ubuntu Noble",
		NoGoAgentSuffix: true,

		OSMatches: []StemcellOSMatch{
			{OSName: "ubuntu", OSVersion: "noble"},
		},

		SupportedInfrastructures: nobleInfrastructures,

		Sort: 2,
	}
	ubuntuJammyDistro = Distro{
		NameName: "ubuntu-jammy",
		Name:     "Ubuntu Jammy",

		OSMatches: []StemcellOSMatch{
			{OSName: "ubuntu", OSVersion: "jammy"},
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
		ubuntuResoluteDistro,
		ubuntuNobleDistro,
		ubuntuJammyDistro,
		windows2019Distro,
	}
)
