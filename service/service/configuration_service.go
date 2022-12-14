package service

import (
	"github.com/orensho/thin-slack-blackbox-tester/service/config"
	"github.com/pkg/errors"
)

func CreateConfigurationService(configReader config.FgBlackboxConfigReader) (*ConfigurationService, error) {
	// Load Configuration
	serverSettings, err := configReader.LoadServerSettings()
	if err != nil {
		return nil, errors.Wrap(err, "Failed Loading Server-Settings")
	}

	testerSettings, err := configReader.LoadTesterSettings()
	if err != nil {
		return nil, errors.Wrap(err, "Failed Loading Tester-Settings")
	}

	metricsSettings, err := configReader.LoadMetricsSettings()
	if err != nil {
		return nil, errors.Wrap(err, "Failed Loading Metrics-Settings")
	}

	testerConfig, err := configReader.LoadTesterConfig(testerSettings)
	if err != nil {
		return nil, errors.Wrap(err, "Failed Loading tester config")
	}

	proxiesSettings, err := configReader.LoadProxiesSettings()
	if err != nil {
		return nil, errors.Wrap(err, "Failed Loading Proxies-Settings")
	}

	return &ConfigurationService{
		ServerSettings:  serverSettings,
		TesterSettings:  testerSettings,
		TesterConfig:    testerConfig,
		MetricsSettings: metricsSettings,
		ProxiesSettings: proxiesSettings,
	}, nil
}

type ConfigurationService struct {
	ServerSettings  *config.ServerSettings
	TesterSettings  *config.TesterSettings
	TesterConfig    *config.TesterConfig
	MetricsSettings *config.MetricsSettings
	ProxiesSettings *config.ProxiesSettings
}
