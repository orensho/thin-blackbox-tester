package steps

import (
	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
)

type StepInterface interface {
	GetType() string
	GetName() string
	Init(name string, conf map[string]interface{}) error
	Run(logger *log.Entry) chromedp.Tasks
}
