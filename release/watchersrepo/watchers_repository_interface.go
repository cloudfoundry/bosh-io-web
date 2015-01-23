package watchersrepo

type WatchersRepository interface {
	ListAll() ([]WatcherRec, error)
	Add(string, string) error
	Remove(string) error
}
