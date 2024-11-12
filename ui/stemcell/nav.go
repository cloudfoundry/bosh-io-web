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
			distronav.Add(nav.Link{
				Title: infra.Title,
				URL:   fmt.Sprintf("/stemcells/bosh-%s-%s-%s-go_agent", infra.Name, infra.SupportedHypervisors[0].Hypervisor.Name, distro.OSMatches[0].Name()),
			})
		}

		allnav.Add(distronav)
	}

	root.Add(allnav)
	root.Activate("#stemcells")

	return root
}
