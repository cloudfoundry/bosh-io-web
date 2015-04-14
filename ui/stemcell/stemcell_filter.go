package stemcell

type StemcellFilter struct {
	Name string

	IncludeDeprecatedDistros bool
}

func (f StemcellFilter) ShowingAllVersions() bool {
	return len(f.Name) > 0
}

func (f StemcellFilter) Limit() int {
	if len(f.Name) > 0 {
		return 30
	}

	return 1
}
