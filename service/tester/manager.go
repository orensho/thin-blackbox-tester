package tester

import (
	"context"

	"github.com/orensho/thin-slack-blackbox-tester/service/config"
	"github.com/orensho/thin-slack-blackbox-tester/service/service"
	"github.com/orensho/thin-slack-blackbox-tester/service/steps"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

type ManagerInterface interface {
	Init(conf *config.TesterConfig) error
	Start()
	Stop()
}

type managerImpl struct {
	testerSettings *config.TesterSettings
	stepsFactory   steps.StepFactoryInterface
	cron           *cron.Cron
	testContext    context.Context
	testCancel     context.CancelFunc
	metricsService service.MetricsServiceInterface
}

func NewManager(
	rootCtx context.Context,
	stepsFactory steps.StepFactoryInterface,
	testerSettings *config.TesterSettings,
	metricsService service.MetricsServiceInterface,
) ManagerInterface {
	testContext, testCancel := context.WithCancel(rootCtx)
	cronLogger := cron.PrintfLogger(log.StandardLogger())

	return &managerImpl{
		testerSettings: testerSettings,
		metricsService: metricsService,
		stepsFactory:   stepsFactory,
		cron: cron.New(cron.WithChain(
			cron.Recover(cronLogger),
			cron.SkipIfStillRunning(cronLogger),
		)),
		testContext: testContext,
		testCancel:  testCancel,
	}
}

func (m *managerImpl) Start() {
	m.cron.Start()
}

func (m *managerImpl) Stop() {
	// cancel all running flows
	m.testCancel()

	// stop the cron
	doneCtx := m.cron.Stop()

	// wait for cron jobs to stop
	<-doneCtx.Done()
}

func (m *managerImpl) Init(conf *config.TesterConfig) error {
	// create all test flows
	for flowName, flowDefinition := range conf.Flows {
		// create flow steps
		flowSteps, err := m.createFlowSteps(conf.Definitions, flowDefinition.Steps)
		if err != nil {
			return errors.Wrapf(err, "Failed creating flow '%s' steps", flowName)
		}

		// create the flow
		flow, err := newFlow(
			m.testContext,
			flowName,
			flowDefinition.Config,
			flowSteps,
			m.metricsService,
		)
		if err != nil {
			return errors.Wrapf(err, "Failed creating flow '%s'", flowName)
		}

		// add flow to the scheduler
		_, err = m.cron.AddJob(flowDefinition.Config.Frequency, flow)
		if err != nil {
			return errors.Wrapf(err, "Failed scheduling flow '%s'", flowName)
		}
	}

	return nil
}

func (m *managerImpl) createFlowSteps(stepsDefinition map[string]config.Definition, flowStepNames []string) ([]steps.StepInterface, error) {
	var flowSteps []steps.StepInterface

	for _, stepName := range flowStepNames {
		stepDefinition, ok := stepsDefinition[stepName]
		if !ok {
			return nil, errors.Errorf("undefined step '%s'", stepName)
		}

		// create the step
		step, err := m.stepsFactory.NewStep(stepDefinition.Type)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed creating step '%s'", stepName)
		}

		// init the step configuration
		err = step.Init(stepName, stepDefinition.Config)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed initializing step '%s'", stepName)
		}

		flowSteps = append(flowSteps, step)
	}

	return flowSteps, nil
}
