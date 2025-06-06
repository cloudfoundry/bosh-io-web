package action

import (
	"errors"

	boshplatform "github.com/cloudfoundry/bosh-agent/v2/platform"
)

type DeleteARPEntriesActionArgs struct {
	Ips []string `json:"ips"`
}

type DeleteARPEntriesAction struct {
	platform boshplatform.Platform
}

func NewDeleteARPEntries(platform boshplatform.Platform) DeleteARPEntriesAction {
	return DeleteARPEntriesAction{
		platform: platform,
	}
}

func (a DeleteARPEntriesAction) IsAsynchronous(_ ProtocolVersion) bool {
	return false
}

func (a DeleteARPEntriesAction) IsPersistent() bool {
	return false
}

func (a DeleteARPEntriesAction) IsLoggable() bool {
	return true
}

func (a DeleteARPEntriesAction) Run(args DeleteARPEntriesActionArgs) (interface{}, error) {
	addresses := args.Ips
	for _, address := range addresses {
		_ = a.platform.DeleteARPEntryWithIP(address) //nolint:errcheck
	}

	return map[string]interface{}{}, nil
}

func (a DeleteARPEntriesAction) Resume() (interface{}, error) {
	return nil, errors.New("not supported")
}

func (a DeleteARPEntriesAction) Cancel() error {
	return errors.New("not supported")
}
