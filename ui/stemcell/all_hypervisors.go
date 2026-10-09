package stemcell

var (
	xenHypervisor = Hypervisor{
		Name:  "xen",
		Title: "Xen",
	}
	xenhvmHypervisor = Hypervisor{
		Name:  "xen-hvm",
		Title: "Xen-HVM",
	}
	esxiHypervisor = Hypervisor{
		Name:  "esxi",
		Title: "ESXi",
	}
	kvmHypervisor = Hypervisor{
		Name:  "kvm",
		Title: "KVM",
	}
	hypervHypervisor = Hypervisor{
		Name:  "hyperv",
		Title: "Hyper-V",
	}
	boshliteHypervisor = Hypervisor{
		Name:  "boshlite",
		Title: "BOSH Lite",
	}
	containerHypervisor = Hypervisor{
		Name:  "container",
		Title: "Container",
	}
)
var (
	allHypervisors = Hypervisors{
		xenHypervisor,
		xenhvmHypervisor,
		esxiHypervisor,
		kvmHypervisor,
		hypervHypervisor,
		boshliteHypervisor,
		containerHypervisor,
	}
)
