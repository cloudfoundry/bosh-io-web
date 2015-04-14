package releasesrepo

import (
	"fmt"
	"path"
	"strings"
)

const (
	sourcePrefixGh = "github.com/"
)

type Source struct {
	avatarsResolver avatarsResolver

	Full string
}

// Short returns 'abc/concourse' when source is 'github.com/abc/concourse'
func (s Source) Short() string {
	return strings.TrimPrefix(s.Full, sourcePrefixGh)
}

// ShortName returns 'concourse' when source is 'github.com/abc/concourse'
// todo better naming
func (s Source) ShortName() string {
	parts := strings.Split(s.Full, "/")
	return parts[len(parts)-1]
}

func (s Source) AvatarURL() string {
	return s.avatarsResolver.Resolve(s.locationName())
}

// LocationName returns 'github.com/abc' when source is 'github.com/abc/concourse'
func (s Source) locationName() string { return path.Dir(s.Full) }

func (s Source) FromGithub() bool {
	return strings.Index(s.Full, sourcePrefixGh) == 0
}

func (s Source) GithubURL() string {
	return fmt.Sprintf("https://%s", s.Full)
}
