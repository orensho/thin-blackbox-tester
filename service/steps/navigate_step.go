package steps

import (
	"github.com/chromedp/chromedp"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const navigateStepType = "navigate-step"

type navigateStepConf struct {
	URL string `validate:"required"`
}

type navigateStep struct {
	name string
	conf navigateStepConf
}

func (s *navigateStep) GetType() string {
	return navigateStepType
}

func (s *navigateStep) GetName() string {
	return s.name
}

func (s *navigateStep) Init(name string, input map[string]interface{}) error {
	var conf navigateStepConf
	err := mapstructure.Decode(input, &conf)
	if err != nil {
		return errors.Wrapf(err, "failed parsing step '%s' configuration", s.GetType())
	}

	// validate conf using validate tags
	err = validator.New().Struct(conf)
	if err != nil {
		return errors.Wrapf(err, "failed validating step '%s' configuration", s.GetType())
	}

	s.name = name
	s.conf = conf

	return nil
}

func (s *navigateStep) Run(logger *log.Entry, proxy string) chromedp.Tasks {
	logger.Infof("navigating to %s", s.conf.URL)

	return chromedp.Tasks{
		chromedp.Navigate(s.conf.URL),
		// wait for page elements
		chromedp.WaitReady(`/html/body/p`, chromedp.NodeVisible),
		chromedp.WaitReady(`img`, chromedp.NodeVisible, chromedp.ByQuery),
	}
}
