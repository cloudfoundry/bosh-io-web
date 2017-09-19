package eventlog

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

const (
	ConfigDeviceTypeJSON = "json"
	ConfigDeviceTypeText = "text"
)

type Config struct {
	DeviceType string `json:"device_type"`
}

func (c Config) Validate() error {
	switch c.DeviceType {
	case ConfigDeviceTypeJSON:
		return nil
	case ConfigDeviceTypeText:
		return nil
	default:
		return bosherr.Errorf("Unknown device type '%s'", c.DeviceType)
	}
}
