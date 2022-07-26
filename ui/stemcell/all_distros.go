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
	ubuntuXenialDistro = Distro{
		NameName: "ubuntu-xenial",
		Name:     "Ubuntu Xenial",

		OSMatches: []StemcellOSMatch{
			{OSName: "ubuntu", OSVersion: "xenial"},
		},

		SupportedInfrastructures: allInfrastructures,

		Sort: 3,
	}
	ubuntuTrustyDistro = Distro{
		NameName: "ubuntu-trusty",
		Name:     "Ubuntu Trusty",

		OSMatches: []StemcellOSMatch{
			{OSName: "ubuntu", OSVersion: "trusty"},
		},

		SupportedInfrastructures: trustyInfrastructures,

		Sort: 4,
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

		Sort: 5,
	}
	windows1803Distro = Distro{
		NameName: "windows1803",
		Name:     "Windows 1803",

		OSMatches: []StemcellOSMatch{
			{OSName: "windows", OSVersion: "1803"},
		},

		SupportedInfrastructures: Infrastructures{
			awsInfrastructure,
			googleInfrastructure,
			azureInfrastructure,
		},

		Sort: 6,
	}
	windows2016Distro = Distro{
		NameName: "windows2016",
		Name:     "Windows 2016",

		OSMatches: []StemcellOSMatch{
			{OSName: "windows", OSVersion: "2016"},
		},

		SupportedInfrastructures: Infrastructures{
			awsInfrastructure,
			googleInfrastructure,
			azureInfrastructure,
		},

		Sort: 7,
	}
	windows2012R2Distro = Distro{
		NameName: "windows2012R2",
		Name:     "Windows 2012R2",

		OSMatches: []StemcellOSMatch{
			{OSName: "windows", OSVersion: "2012R2"},
		},

		SupportedInfrastructures: Infrastructures{
			awsInfrastructure,
			googleInfrastructure,
			azureInfrastructure,
		},

		Sort: 8,
	}
	centos7Distro = Distro{
		NameName: "centos-7",
		Name:     "CentOS 7",

		OSMatches: []StemcellOSMatch{
			{OSName: "centos", OSVersion: "7"},
		},

		SupportedInfrastructures: Infrastructures{
			awsInfrastructure,
			googleInfrastructure,
			azureInfrastructure,
			openstackInfrastructure,
			vsphereInfrastructure,
			wardenInfrastructure,
		},

		Sort: 9,
	}
)

var (
	allDistros = []Distro{
		ubuntuJammyDistro,
		ubuntuBionicDistro,
		ubuntuXenialDistro,
		ubuntuTrustyDistro,
		windows2019Distro,
		windows1803Distro,
		windows2016Distro,
		windows2012R2Distro,
		centos7Distro,
	}
)
