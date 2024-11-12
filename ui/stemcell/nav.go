package stemcell

import (
	"fmt"

	"github.com/bosh-io/web/ui/nav"
)

func Navigation() nav.Link {
	root := nav.Link{Title: "Stemcells"}

	allnav := nav.Link{Title: "Stemcells", URL: "#stemcells"}
	allnav.Add(nav.Link{
		Title: "Latest Versions",
		URL:   "/stemcells",
	})

	for _, distro := range allDistros {
		distronav := nav.Link{Title: distro.Name}

		for _, infra := range distro.SupportedInfrastructures {
			var stemcellUrl string
			if distro.NoGoAgentSuffix {
				stemcellUrl = fmt.Sprintf("/stemcells/bosh-%s-%s-%s", infra.Name, infra.SupportedHypervisors[0].Hypervisor.Name, distro.OSMatches[0].Name())
			} else {
				stemcellUrl = fmt.Sprintf("/stemcells/bosh-%s-%s-%s-go_agent", infra.Name, infra.SupportedHypervisors[0].Hypervisor.Name, distro.OSMatches[0].Name())
			}
			distronav.Add(nav.Link{
				Title: infra.Title,
				URL:   stemcellUrl,
			})
		}

		allnav.Add(distronav)
	}

	root.Add(allnav)
	root.Activate("#stemcells")

	return root
}
