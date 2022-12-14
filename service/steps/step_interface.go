package steps

import (
	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

//go:generate mockery -name StepInterface -inpkg -case=underscore -output MockStepInterface
type StepInterface interface {
	GetType() string
	GetName() string
	Init(name string, conf map[string]interface{}) error
	Run(logger *log.Entry, proxy string) chromedp.Tasks
}
