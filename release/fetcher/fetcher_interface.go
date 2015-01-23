package fetcher

type Fetcher interface {
	Fetch(string) (ReleaseDir, error)
}
