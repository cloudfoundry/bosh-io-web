package release

import (
	"fmt"

	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
)

type ReleaseRef struct {
	Source     Source
	PrettyName string
	DocPage    string
}

type ReleaseRefs []ReleaseRef

var BOSH = NewReleaseRef("github.com/cloudfoundry/bosh", "BOSH", "")

var KnownCPIs = ReleaseRefs{
	NewReleaseRef("github.com/cloudfoundry-incubator/bosh-aws-cpi-release", "AWS", "init-aws-v1"),
	NewReleaseRef("github.com/cloudfoundry-incubator/bosh-openstack-cpi-release", "OpenStack", "init-openstack-v1"),
	NewReleaseRef("github.com/cloudfoundry-incubator/bosh-vsphere-cpi-release", "vSphere", "init-vsphere-v1"),
	NewReleaseRef("github.com/cloudfoundry-incubator/bosh-vcloud-cpi-release", "vCloud", "init-vcloud-v1"),
	NewReleaseRef("github.com/cloudfoundry-incubator/bosh-azure-cpi-release", "Azure", "init-azure-v1"),
	NewReleaseRef("github.com/cloudfoundry-incubator/bosh-softlayer-cpi-release", "SoftLayer", "init-softlayer-v1"),
}

func NewReleaseRef(fullSource, prettyName, docPage string) ReleaseRef {
	return ReleaseRef{
		Source:     NewSource(bhrelsrepo.Source{Full: fullSource}),
		PrettyName: prettyName,
		DocPage:    docPage,
	}
}

func (c ReleaseRef) DocPagePath() string { return fmt.Sprintf("/docs/%s.html", c.DocPage) }

func (c ReleaseRef) DocPageLink() string {
	return fmt.Sprintf("<a href='%s'>Initializing a BOSH environment on %s</a>", c.DocPagePath(), c.PrettyName)
}

func (c ReleaseRefs) FindByShortName(name string) (ReleaseRef, bool) {
	for _, cpi := range c {
		if cpi.Source.ShortName() == name {
			return cpi, true
		}
	}

	return ReleaseRef{}, false
}

func (c ReleaseRefs) FindByDocPage(docPage string) (ReleaseRef, bool) {
	for _, cpi := range c {
		if cpi.DocPage == docPage {
			return cpi, true
		}
	}

	return ReleaseRef{}, false
}
