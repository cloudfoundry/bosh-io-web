package stemcell

type StemcellRef struct {
	ManifestName string
	DocPage      string
}

type StemcellRefs []StemcellRef

var KnownStemcells = StemcellRefs{
	NewStemcellRef("bosh-aws-xen-hvm-ubuntu-trusty-go_agent", "init-aws-v1"),
	NewStemcellRef("bosh-openstack-kvm-ubuntu-trusty-go_agent", "init-openstack-v1"),
	NewStemcellRef("bosh-vsphere-esxi-ubuntu-trusty-go_agent", "init-vsphere-v1"),
	NewStemcellRef("bosh-vcloud-esxi-ubuntu-trusty-go_agent", "init-vcloud-v1"),
	NewStemcellRef("bosh-azure-hyperv-ubuntu-trusty-go_agent", "init-azure-v1"),
	NewStemcellRef("bosh-softlayer-xen-ubuntu-trusty-go_agent", "init-softlayer-v1"),
}

func NewStemcellRef(manifestName, docPage string) StemcellRef {
	return StemcellRef{
		ManifestName: manifestName,
		DocPage:      docPage,
	}
}

func (refs StemcellRefs) FindByDocPage(docPage string) (StemcellRef, bool) {
	for _, r := range refs {
		if r.DocPage == docPage {
			return r, true
		}
	}

	return StemcellRef{}, false
}
