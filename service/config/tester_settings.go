package config

type TesterSettings struct {
	ConfigFilename   string `env:"TESTER_CONFIG_FILENAME" envDefault:"config.yaml"`
	ConfigFolder     string `env:"TESTER_CONFIG_FOLDER" envDefault:"configuration"`
	ShowDebugBrowser bool   `env:"TESTER_SHOW_DEBUG_BROWSER" envDefault:"false"`
	Environment      string `env:"TESTER_ENVIRONMENT" envDefault:"dev"`
}

func (s *TesterSettings) Evaluate() error {
	return nil
}

func (s *TesterSettings) Validate() error {
	return nil
}
