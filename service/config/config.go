package config

import (
	"fmt"
	envParser "github.com/caarlos0/env"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
)

type BlackboxConfigReader struct {
}

func NewBlackboxConfigReader() BlackboxConfigReader {
	r := BlackboxConfigReader{}

	return r
}

func (c *BlackboxConfigReader) LoadServerSettings() (*ServerSettings, error) {
	serverSettings := ServerSettings{}
	err := loadEnvironmentVariables(&serverSettings)

	return &serverSettings, err
}

func (c *BlackboxConfigReader) LoadTesterSettings() (*TesterSettings, error) {
	testerSettings := TesterSettings{}
	err := loadEnvironmentVariables(&testerSettings)

	return &testerSettings, err
}

func (c *BlackboxConfigReader) LoadMetricsSettings() (*MetricsSettings, error) {
	metricsSettings := MetricsSettings{}
	err := loadEnvironmentVariables(&metricsSettings)

	return &metricsSettings, err
}

func (c *BlackboxConfigReader) LoadTesterConfig(testerSettings *TesterSettings) (*TesterConfig, error) {
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

func loadEnvironmentVariables(settings SettingsInterface) error {

	err := envParser.Parse(settings)
	if err != nil {

		return errors.Wrap(err, "Failed to load environment variables.")
	}

	if err := settings.Validate(); err != nil {

		return errors.Wrap(err, "setting validation failed, recheck your settings")
	}

	if err := settings.Evaluate(); err != nil {

		return errors.Wrap(err, "settings evaluation failed, recheck your settings")
	}

	log.Infof("Loaded Settings %#v", settings)
	return nil
}
