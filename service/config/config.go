package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/hashicorp/go-multierror"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Supplies Loaded Evaluated and Validated configuration
type FgBlackboxConfigReader struct {
	*conf.ConfigReader
}

func NewFigBlackboxConfigReader() FgBlackboxConfigReader {
	r := FgBlackboxConfigReader{}
	r.Init()

	return r
}

func (c *FgBlackboxConfigReader) LoadProxiesSettings() (*ProxiesSettings, error) {
	proxies := ProxiesSettings{}

	return &proxies, nil
}

func (c *FgBlackboxConfigReader) LoadLoggerSettings() (*logging.LoggerSettings, error) {
	loggerConfig := logging.LoggerSettings{}
	err := c.LoadEnvironmentVariables(&loggerConfig)

	return &loggerConfig, err
}

func (c *FgBlackboxConfigReader) LoadAdministrativeSettings() (*AdministrativeSettings, error) {
	administrativeSettings := AdministrativeSettings{}
	err := c.LoadEnvironmentVariables(&administrativeSettings)

	return &administrativeSettings, err
}

func (c *FgBlackboxConfigReader) LoadServerSettings() (*ServerSettings, error) {
	serverSettings := ServerSettings{}
	err := c.LoadEnvironmentVariables(&serverSettings)

	return &serverSettings, err
}

func (c *FgBlackboxConfigReader) LoadTesterSettings() (*TesterSettings, error) {
	testerSettings := TesterSettings{}
	err := c.LoadEnvironmentVariables(&testerSettings)

	return &testerSettings, err
}

func (c *FgBlackboxConfigReader) LoadMetricsSettings() (*MetricsSettings, error) {
	metricsSettings := MetricsSettings{}
	err := c.LoadEnvironmentVariables(&metricsSettings)

	return &metricsSettings, err
}

func (c *FgBlackboxConfigReader) LoadTesterConfig(testerSettings *TesterSettings) (*TesterConfig, error) {
	testerConfig := TesterConfig{}
	configFilePath := path.Join(testerSettings.ConfigFolder, testerSettings.Environment, testerSettings.ConfigFilename)
	yamlFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed reading file: %s", configFilePath)
	}

	yamlFileWithEnv, err := expendEnvVars(string(yamlFile))
	if err != nil {
		return nil, errors.Wrapf(err, "failed while expanding environment variables in file: %s", configFilePath)
	}

	err = yaml.Unmarshal([]byte(yamlFileWithEnv), &testerConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed unmarshalling file: %s", configFilePath)
	}

	return &testerConfig, err
}

func expendEnvVars(text string) (string, error) {
	var missingVars *multierror.Error
	yamlFileWithEnv := os.Expand(text, func(varName string) string {
		value := os.Getenv(varName)
		if value == "" {
			missingVars = multierror.Append(missingVars, fmt.Errorf("missing environment variable %s", varName))
		}

		return value
	})

	return yamlFileWithEnv, missingVars.ErrorOrNil()
}

func (c *FgBlackboxConfigReader) Init() {
	// Load environment variables from configuration file
	err := c.LoadEnvironmentVariablesFromFile()
	if err != nil {
		log.WithError(err).Warn("Unable to load configuration file. Continue with environment variables.")
	}
}
