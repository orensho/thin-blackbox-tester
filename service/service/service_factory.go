package service

import "github.com/pkg/errors"

type Factory struct {
	ConfigurationService *ConfigurationService
	MetricsService       MetricsServiceInterface
}

func NewServiceFactory(configurationService *ConfigurationService) (*Factory, error) {
	metricsService, err := NewMetricsService(configurationService.MetricsSettings)
	if err != nil {
		return nil, errors.Wrap(err, "Failed creating metrics service")
	}

	sf := Factory{
		ConfigurationService: configurationService,
		MetricsService:       metricsService,
	}

	return &sf, nil
}
