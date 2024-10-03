package stemcell

var (
	ubuntuJammyDistro = Distro{
		NameName: "ubuntu-jammy",
		Name:     "Ubuntu Jammy",

		OSMatches: []StemcellOSMatch{
			{OSName: "ubuntu", OSVersion: "jammy"},
		},

		SupportedInfrastructures: allInfrastructures,

		Sort: 1,
	}
	ubuntuBionicDistro = Distro{
		NameName: "ubuntu-bionic",
		Name:     "Ubuntu Bionic",

		OSMatches: []StemcellOSMatch{
			{OSName: "ubuntu", OSVersion: "bionic"},
		},

		SupportedInfrastructures: allInfrastructures,

		Sort: 2,
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

		Sort: 3,
	}
)

var (
	allDistros = []Distro{
		ubuntuJammyDistro,
		ubuntuBionicDistro,
		windows2019Distro,
	}
)
