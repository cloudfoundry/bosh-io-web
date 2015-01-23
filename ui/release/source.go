package release

import (
	"fmt"
	"strings"
)

const (
	sourcePrefixGh = "github.com/"
)

type Source struct {
	Full  string
	Short string
}

type SourceSorting []Source

func NewSource(full string) Source {
	return Source{
		Full:  full,
		Short: strings.TrimPrefix(full, sourcePrefixGh),
	}
}

func (s Source) String() string { return s.Full }

func (s Source) URL() string {
	return fmt.Sprintf("/releases/%s", s.Full)
}

func (s Source) FromGithub() bool {
	return strings.Index(s.Full, sourcePrefixGh) == 0
}

func (s Source) GithubURL() string {
	return fmt.Sprintf("https://%s", s.Full)
}

func (s SourceSorting) Len() int           { return len(s) }
func (s SourceSorting) Less(i, j int) bool { return s[i].Full < s[j].Full }
func (s SourceSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
