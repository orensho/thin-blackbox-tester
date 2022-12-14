package config

type MetricsSettings struct {
	Enabled      bool   `env:"METRICS_ENABLE" envDefault:"true"`
	Environment  string `env:"METRICS_ENVIRONMENT" envDefault:"local"`
	MetricPrefix string `env:"METRICS_PREFIX" envDefault:"fg_blackbox"`
	MetricPort   string `env:"METRICS_PORT" envDefault:"8888"`
}

func (s *MetricsSettings) Evaluate() error {
	return nil
}

func (s *MetricsSettings) Validate() error {
	return nil
}
