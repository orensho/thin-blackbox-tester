package config

type TesterConfig struct {
	Definitions map[string]Definition `yaml:"definitions"`
	Flows       map[string]Flow       `yaml:"flows"`
}

type Definition struct {
	Type   string                 `yaml:"type"`
	Config map[string]interface{} `yaml:"config"`
}

type Flow struct {
	Config FlowConfig `yaml:"config"`
	Steps  []string   `yaml:"steps"`
}

type FlowConfig struct {
	Frequency string  `yaml:"frequency"`
	Timeout   *string `yaml:"timeout,omitempty"`
}
