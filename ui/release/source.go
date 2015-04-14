package release

import (
	"fmt"

	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
)

type Source struct {
	src bhrelsrepo.Source
}

type SourceSorting []Source

func NewSource(src bhrelsrepo.Source) Source {
	return Source{src: src}
}

func (s Source) Full() string { return s.src.Full }

func (s Source) Short() string { return s.src.Short() }

func (s Source) AvatarURL() string { return s.src.AvatarURL() }

func (s Source) FromGithub() bool { return s.src.FromGithub() }

func (s Source) GithubURL() string { return s.src.GithubURL() }

func (s Source) String() string { return s.src.Full }

func (s Source) URL() string {
	return fmt.Sprintf("/releases/%s", s.src.Full)
}

func (s SourceSorting) Len() int           { return len(s) }
func (s SourceSorting) Less(i, j int) bool { return s[i].Full() < s[j].Full() }
func (s SourceSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
