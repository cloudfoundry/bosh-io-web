package release

import "github.com/bosh-io/web/ui/nav"

func Navigation() nav.Link {
	root := nav.Link{Title: "Releases"}

	allnav := nav.Link{Title: "Releases", URL: "#releases", IsSection: true}
	allnav.Add(nav.Link{
		Title: "Browse Releases",
		URL:   "/releases",
	})
	root.Add(allnav)

	return root
}
