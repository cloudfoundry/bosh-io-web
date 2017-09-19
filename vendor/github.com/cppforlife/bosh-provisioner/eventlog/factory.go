package eventlog

import (
	"os"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type Factory struct {
	config Config
	logger boshlog.Logger
}

func NewFactory(config Config, logger boshlog.Logger) Factory {
	return Factory{config: config, logger: logger}
}

func (f Factory) NewLog() Log {
	var device Device

	switch f.config.DeviceType {
	case ConfigDeviceTypeJSON:
		device = NewJSONDevice(os.Stdout)
	case ConfigDeviceTypeText:
		device = NewTextDevice(os.Stdout)
	default:
		// config should be validated before using it with a factory
		panic(bosherr.Errorf("Unknown device type '%s'", f.config.DeviceType))
	}

	return NewLog(device, f.logger)
}
