package steps

import (
	"time"

	"github.com/chromedp/chromedp"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const waitStepType = "wait-step"

type waitStepConf struct {
	Duration string `validate:"required"`

	durationParsed time.Duration
}

type waitStep struct {
	name string
	conf waitStepConf
}

func (s *waitStep) GetType() string {
	return waitStepType
}

func (s *waitStep) GetName() string {
	return s.name
}

func (s *waitStep) Init(name string, input map[string]interface{}) error {
	var conf waitStepConf
	err := mapstructure.Decode(input, &conf)
	if err != nil {
		return errors.Wrapf(err, "failed parsing step '%s' configuration", s.GetType())
	}

	// validate conf using validate tags
	err = validator.New().Struct(conf)
	if err != nil {
		return errors.Wrapf(err, "failed validating step '%s' configuration", s.GetType())
	}

	parsedDuration, err := time.ParseDuration(conf.Duration)
	if err != nil {
		return errors.Wrapf(err, "failed parsing step '%s' duration", s.GetType())
	}
	conf.durationParsed = parsedDuration

	s.name = name
	s.conf = conf

	return nil
}

func (s *waitStep) Run(logger *log.Entry) chromedp.Tasks {
	logger.Infof("Waiting for %s", s.conf.Duration)

	return chromedp.Tasks{
		chromedp.Sleep(s.conf.durationParsed),
	}
}
