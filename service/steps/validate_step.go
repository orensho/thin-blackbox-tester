package steps

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	_ "image/png" // png decode support
	"time"

	"github.com/chromedp/chromedp"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const validateStepType = "validate-step"
const pollFunction = `() => {
		const unloaded_images = [...document.querySelectorAll("img")].filter((image) => { return image.height == 0 })
		return unloaded_images.length == 0
	}
`

type validateStepConf struct {
	Hash     string `validate:"required"`
	Selector string `validate:"required"`
}

type validateStep struct {
	name string
	conf validateStepConf
}

func (s *validateStep) GetType() string {
	return validateStepType
}

func (s *validateStep) GetName() string {
	return s.name
}

func (s *validateStep) Init(name string, input map[string]interface{}) error {
	var conf validateStepConf
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

func (s *validateStep) Run(logger *log.Entry) chromedp.Tasks {
	logger.Infof("Running validate step with conf %+v", s.conf)

	validationTasks := chromedp.Tasks{}

	validationTasks = append(validationTasks, chromedp.ActionFunc(func(ctx context.Context) error {
		var buf []byte

		const timeoutMS = 100
		// wait for the image to load, this is required because isolated pages are ""loaded"" immediately
		// and then JS pops in elements in the background
		chromedp.PollFunction(pollFunction, nil, chromedp.WithPollingTimeout(timeoutMS*time.Millisecond))

		err := chromedp.Screenshot(s.conf.Selector, &buf, chromedp.NodeVisible).Do(ctx)

		if err != nil {
			logger.Error("failed to take screenshot")
		}

		return s.validateScreenshot(logger, buf)
	}))

	return validationTasks
}

func (s *validateStep) validateScreenshot(logger *log.Entry, buf []byte) error {
	logger.Infof("validating test screenshot")

	hash := md5.Sum(buf)
	hashString := hex.EncodeToString(hash[:])

	if hashString != s.conf.Hash {
		return errors.Errorf("different hash result. Expected <= %s got %s.", s.conf.Hash, hashString)
	}

	logger.Info("Done validating image hash")
	return nil
}
