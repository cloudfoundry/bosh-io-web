package stemcell

var (
	awsInfrastructure = Infrastructure{
		Name:             "aws",
		Title:            "Amazon Web Services",
		DocumentationURL: "/docs/aws-cpi/",
		SupportedHypervisors: InfrastructureHypervisors{
			{
				Hypervisor: xenhvmHypervisor,
			},
			{
				Hypervisor: xenHypervisor,
				Deprecated: true,
			},
		},
	}
	googleInfrastructure = Infrastructure{
		Name:             "google",
		Title:            "Google Cloud Platform",
		DocumentationURL: "/docs/google-cpi/",
		SupportedHypervisors: InfrastructureHypervisors{
			{
				Hypervisor: kvmHypervisor,
			},
		},
	}
	azureInfrastructure = Infrastructure{
		Name:             "azure",
		Title:            "Microsoft Azure",
		DocumentationURL: "/docs/azure-cpi/",
		SupportedHypervisors: InfrastructureHypervisors{
			{
				Hypervisor: hypervHypervisor,
			},
		},
	}
	openstackInfrastructure = Infrastructure{
		Name:             "openstack",
		Title:            "OpenStack",
		DocumentationURL: "/docs/openstack-cpi/",
		SupportedHypervisors: InfrastructureHypervisors{
			{
				Hypervisor: kvmHypervisor,
			},
		},
	}
	softlayerInfrastructure = Infrastructure{
		Name:             "softlayer",
		Title:            "SoftLayer",
		DocumentationURL: "/docs/softlayer-cpi/",
		SupportedHypervisors: InfrastructureHypervisors{
			{
				Hypervisor: xenHypervisor,
			},
		},
	}
	vcloudInfrastructure = Infrastructure{
		Name:             "vcloud",
		Title:            "VMware vCloud",
		DocumentationURL: "/docs/vcloud-cpi/",
		SupportedHypervisors: InfrastructureHypervisors{
			{
				Hypervisor: esxiHypervisor,
			},
		},
	}
	vsphereInfrastructure = Infrastructure{
		Name:             "vsphere",
		Title:            "VMware vSphere",
		DocumentationURL: "/docs/vsphere-cpi/",
		SupportedHypervisors: InfrastructureHypervisors{
			{
				Hypervisor: esxiHypervisor,
			},
		},
	}
	wardenInfrastructure = Infrastructure{
		Name:             "warden",
		Title:            "Warden (BOSH Lite)",
		DocumentationURL: "/docs/warden-cpi/",
		SupportedHypervisors: InfrastructureHypervisors{
			{
				Hypervisor: boshliteHypervisor,
			},
		},
	}
	alicloudInfrastructure = Infrastructure{
		Name:             "alicloud",
		Title:            "Alibaba Cloud",
		DocumentationURL: "/docs/alicloud-cpi/",
		SupportedHypervisors: InfrastructureHypervisors{
			{
				Hypervisor: kvmHypervisor,
			},
		},
	}
)

var (
	trustyInfrastructures = Infrastructures{
		awsInfrastructure,
		googleInfrastructure,
		azureInfrastructure,
		openstackInfrastructure,
		softlayerInfrastructure,
		vcloudInfrastructure,
		vsphereInfrastructure,
		wardenInfrastructure,
	}
	allInfrastructures = Infrastructures{
		awsInfrastructure,
		googleInfrastructure,
		azureInfrastructure,
		alicloudInfrastructure,
		openstackInfrastructure,
		softlayerInfrastructure,
		vcloudInfrastructure,
		vsphereInfrastructure,
		wardenInfrastructure,
	}
	nobleInfrastructures = Infrastructures{
		awsInfrastructure,
		googleInfrastructure,
		azureInfrastructure,
		openstackInfrastructure,
		vsphereInfrastructure,
		wardenInfrastructure,
	}
)
