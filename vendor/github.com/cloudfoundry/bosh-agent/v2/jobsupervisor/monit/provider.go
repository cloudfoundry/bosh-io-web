package monit

import (
	"net/http"
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cloudfoundry/bosh-utils/httpclient"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	boshplatform "github.com/cloudfoundry/bosh-agent/v2/platform"
)

const (
	shortRetryStrategyAttempts = uint(20)
	longRetryStrategyAttempts  = uint(300)
	retryDelay                 = 1 * time.Second
	monitHost                  = "127.0.0.1:2822"
)

type ClientProvider interface {
	Get() (Client, error)
}

type clientProvider struct {
	platform        boshplatform.Platform
	logger          boshlog.Logger
	shortHTTPClient HTTPClient
	longHTTPClient  HTTPClient
}

func NewProvider(platform boshplatform.Platform, logger boshlog.Logger) ClientProvider {
	httpClient := http.DefaultClient

	shortHTTPClient := httpclient.NewRetryClient(
		httpClient,
		shortRetryStrategyAttempts,
		retryDelay,
		logger,
	)

	longHTTPClient := NewMonitRetryClient(
		httpClient,
		longRetryStrategyAttempts,
		shortRetryStrategyAttempts,
		retryDelay,
		logger,
	)

	return clientProvider{
		platform:        platform,
		logger:          logger,
		shortHTTPClient: shortHTTPClient,
		longHTTPClient:  longHTTPClient,
	}
}

func (p clientProvider) Get() (client Client, err error) {
	monitUser, monitPassword, err := p.platform.GetMonitCredentials()
	if err != nil {
		return nil, bosherr.WrapError(err, "Getting monit credentials")
	}

	return NewHTTPClient(
		monitHost,
		monitUser,
		monitPassword,
		p.shortHTTPClient,
		p.longHTTPClient,
		p.logger,
	), nil
}
