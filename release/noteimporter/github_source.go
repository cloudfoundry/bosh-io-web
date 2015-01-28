package noteimporter

import (
	"strings"

	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
)

type githubSource struct {
	Owner string
	Repo  string
}

func newGithubSource(source bhrelsrepo.Source) (githubSource, bool) {
	parts := strings.Split(string(source), "/")
	if len(parts) == 3 && parts[0] == "github.com" {
		return githubSource{Owner: parts[1], Repo: parts[2]}, true
	}

	// Not a github source
	return githubSource{}, false
}
