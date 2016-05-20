package watchersrepo

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	bpindex "github.com/cppforlife/bosh-provisioner/index"
)

const cwRepositoryLogTag = "ConcreteWatchersRepository"

type CWRepository struct {
	index  bpindex.Index
	logger boshlog.Logger
}

type watcherRecKey struct {
	RelSource string
}

func NewConcreteWatchersRepository(index bpindex.Index, logger boshlog.Logger) CWRepository {
	return CWRepository{
		index:  index,
		logger: logger,
	}
}

func (r CWRepository) ListAll() ([]WatcherRec, error) {
	var watcherRecs []WatcherRec

	err := r.index.List(&watcherRecs)
	if err != nil {
		return nil, bosherr.WrapError(err, "Finding watchers")
	}

	return watcherRecs, nil
}

func (r CWRepository) Add(relSource, minVersion string) error {
	watcherRec := WatcherRec{
		RelSource:     relSource,
		MinVersionRaw: minVersion,
	}

	err := watcherRec.Validate()
	if err != nil {
		return bosherr.WrapError(err, "Validating watcher")
	}

	// Since key is unique it is ok to dup the item; value is overriden
	err = r.index.Save(watcherRecKey{relSource}, watcherRec)
	if err != nil {
		return bosherr.WrapError(err, "Saving watcher")
	}

	return nil
}

func (r CWRepository) Remove(relSource string) error {
	if len(relSource) == 0 {
		return bosherr.New("Expected release source to be non-empty")
	}

	err := r.index.Remove(watcherRecKey{relSource})
	if err != nil {
		return bosherr.WrapError(err, "Removing watcher")
	}

	return nil
}
