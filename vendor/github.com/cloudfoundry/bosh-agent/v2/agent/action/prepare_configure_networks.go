package action

import (
	"errors"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	boshplatform "github.com/cloudfoundry/bosh-agent/v2/platform"
	boshsettings "github.com/cloudfoundry/bosh-agent/v2/settings"
)

type PrepareConfigureNetworksAction struct {
	platform        boshplatform.Platform
	settingsService boshsettings.Service
}

func NewPrepareConfigureNetworks(
	platform boshplatform.Platform,
	settingsService boshsettings.Service,
) PrepareConfigureNetworksAction {
	return PrepareConfigureNetworksAction{
		platform:        platform,
		settingsService: settingsService,
	}
}

func (a PrepareConfigureNetworksAction) IsAsynchronous(_ ProtocolVersion) bool {
	return false
}

func (a PrepareConfigureNetworksAction) IsPersistent() bool {
	return false
}

func (a PrepareConfigureNetworksAction) IsLoggable() bool {
	return true
}

func (a PrepareConfigureNetworksAction) Run() (string, error) {
	err := a.settingsService.InvalidateSettings()
	if err != nil {
		return "", bosherr.WrapError(err, "Invalidating settings")
	}

	err = a.platform.PrepareForNetworkingChange()
	if err != nil {
		return "", bosherr.WrapError(err, "Preparing for networking change")
	}

	return "ok", nil
}

func (a PrepareConfigureNetworksAction) Resume() (interface{}, error) {
	return nil, errors.New("not supported")
}

func (a PrepareConfigureNetworksAction) Cancel() error {
	return errors.New("not supported")
}
