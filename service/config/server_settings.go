package config

import (
	"time"

	"github.com/pkg/errors"
)

type ServerSettings struct {
	LocalListenIP       string `env:"SERVER_LOCAL_LISTEN_IP" envDefault:"127.0.0.1"`
	LocalListenPort     string `env:"SERVER_LOCAL_LISTEN_PORT" envDefault:"8080"`
	ShutDownGracePeriod string `env:"SERVER_SHUTDOWN_GRACE_PERIOD" envDefault:"10s"`

	ParsedShutDownGracePeriod time.Duration
}

func (s *ServerSettings) Evaluate() error {
	period, err := time.ParseDuration(s.ShutDownGracePeriod)
	if err != nil {
		return errors.Wrap(err, "Unable to parse SHUTDOWN_GRACE_PERIOD")
	}
	s.ParsedShutDownGracePeriod = period

	return nil
}

func (s *ServerSettings) Validate() error {
	return nil
}
