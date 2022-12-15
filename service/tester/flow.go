package tester

import (
	"context"
	"errors"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/orensho/thin-slack-blackbox-tester/service/service"

	"github.com/chromedp/chromedp"
	"github.com/orensho/thin-slack-blackbox-tester/service/config"
	"github.com/orensho/thin-slack-blackbox-tester/service/steps"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type flow struct {
	name           string
	config         config.FlowConfig
	steps          []steps.StepInterface
	rootCtx        context.Context
	metricsService service.MetricsServiceInterface

	timeout time.Duration // calculated from config
}

var DefaultTimeout = time.Minute * 10

func newFlow(
	rootCtx context.Context,
	name string,
	conf config.FlowConfig,
	flowSteps []steps.StepInterface,
	metricsService service.MetricsServiceInterface,
) (*flow, error) {
	flow := &flow{
		name:           name,
		config:         conf,
		steps:          flowSteps,
		timeout:        DefaultTimeout,
		rootCtx:        rootCtx,
		metricsService: metricsService,
	}

	if conf.Timeout != nil {
		timeout, err := time.ParseDuration(*conf.Timeout)
		if err != nil {
			return nil, err
		}
		flow.timeout = timeout
	}

	return flow, nil
}

func (f *flow) Run() {
	logger := log.WithFields(log.Fields{
		"flow":  f.name,
		"runId": uuid.NewV4(), // unique id per run for easy logs debugging
	})
	logger.Infof("Starting flow %s", f.name)

	browserCtx, cancelFunc := f.createTabContext(logger)
	defer cancelFunc() // releases resources
	f.stepsRun(browserCtx, logger)

	logger.Infof("Finished flow successfully %s", f.name)
}

func (f *flow) stepsRun(browserCtx context.Context, logger *log.Entry) {
	for _, step := range f.steps {
		step := step // avoid closure

		logger = logger.WithFields(log.Fields{
			"step": step.GetName(),
		})

		logger.Infof("executing step '%s'", step.GetName())
		stepStartTime := time.Now()

		// execute the tasks returned from the step
		err := chromedp.Run(browserCtx, step.Run(logger))

		ms := float64(time.Since(stepStartTime).Nanoseconds()) / 1e6
		logger.Infof("flow duration %fms", ms)
		errMetrics := f.metricsService.ReportStepTestDuration(browserCtx, f.name, ms, step.GetName())
		if errMetrics != nil {
			logger.WithError(errMetrics).Error("failed reporting step duration")
		}

		// failure
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				errReport := f.metricsService.ReportStepTestTimeout(f.rootCtx, f.name, step.GetName())
				if errReport != nil {
					logger.WithError(errReport).Error("failed reporting step timeout")
				}

				logger.WithError(err).
					Errorf("flow timeout after %s in step '%s'", f.timeout.String(), step.GetName())
			} else {
				errReport := f.metricsService.ReportStepTestError(f.rootCtx, f.name, step.GetName())
				if errReport != nil {
					logger.WithError(errReport).Error("failed reporting step error")
				}

				var cdpErr *runtime.ExceptionDetails
				if errors.As(err, &cdpErr) && cdpErr.Exception != nil {
					logger.
						WithError(err).
						WithField("ErrClassName", cdpErr.Exception.ClassName).
						WithField("ErrDescription", cdpErr.Exception.Description).
						Errorf("executing step '%s' returned an error, stopping flow", step.GetName())
				} else {
					logger.WithError(err).
						Errorf("executing step '%s' returned an error, stopping flow", step.GetName())
				}
			}
			return
		}

		// success
		err = f.metricsService.ReportStepTestSuccess(f.rootCtx, f.name, step.GetName())
		if err != nil {
			logger.WithError(err).Error("failed reporting step success")
		}
		logger.Infof("finished successfully executing step '%s'", step.GetName())
	}
}

func (f *flow) createTabContext(logger *log.Entry) (context.Context, context.CancelFunc) {
	// create flow context with timeout
	flowContext, flowCancel := context.WithTimeout(f.rootCtx, f.timeout)

	// create the flags for the headless browser
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.UserAgent("Firefox/80"),
	)

	// create a allocator with the new flags
	allocCtx, allocCancel := chromedp.NewExecAllocator(flowContext, opts...)

	// create browser context with the allocator and logging
	browserCtx, browserCancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithDebugf(logger.Debugf),
		chromedp.WithLogf(logger.Infof),
		chromedp.WithErrorf(logger.Errorf))

	return browserCtx, func() {
		// call all created cancel func
		browserCancel()
		allocCancel()
		flowCancel()
	}
}
