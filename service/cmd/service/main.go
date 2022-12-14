package main

import (
	"fmt"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/orensho/thin-slack-blackbox-tester/service/config"
	"github.com/orensho/thin-slack-blackbox-tester/service/server"
	"github.com/orensho/thin-slack-blackbox-tester/service/service"
	"github.com/orensho/thin-slack-blackbox-tester/service/steps"
	"github.com/orensho/thin-slack-blackbox-tester/service/tester"
	log "github.com/sirupsen/logrus"

	"net/http"
	"os"
	"time"
)

const blockingTimeout = 5 * time.Second

func main() {
	// disable the logger until reading the logger configuration
	log.SetLevel(log.PanicLevel)

	configReader := config.NewFigBlackboxConfigReader()
	logSettings, err := configReader.LoadLoggerSettings()
	if err != nil {
		logSettings = &logging.LoggerSettings{}
		log.Warn("Couldn't load logger settings")
	}
	logSettings.ServiceName = "fg-blackbox"

	logger := logging.InitLogger(logSettings)
	defer logger.Close()

	c, _ := os.Getwd()
	log.Infof("CWD: %v", c)

	configService, err := service.CreateConfigurationService(configReader)
	if err != nil {
		log.Panic(err)
	}

	rootCtx := app_context.NewRootAppContext()

	serviceFactory, err := service.NewServiceFactory(configService)
	if err != nil {
		log.Panic(err)
	}

	testManager := tester.NewManager(
		rootCtx.GetContext(),
		steps.NewStepFactory(),
		serviceFactory.ConfigurationService.TesterSettings,
		serviceFactory.MetricsService,
	)

	err = testManager.Init(serviceFactory.ConfigurationService.TesterConfig)
	if err != nil {
		log.WithError(err).Panic("failed initiating test manager")
	}

	FgBlackbox := server.NewFgBlackboxServer(configService.ServerSettings, serviceFactory)

	if configService.MetricsSettings.Enabled {
		pe, err := prometheus.NewExporter(prometheus.Options{
			Namespace: configService.MetricsSettings.MetricPrefix,
		})
		if err != nil {
			log.Panicf("Failed to create the Prometheus stats exporter: %v", err)
		}

		go func() {
			mux := http.NewServeMux()
			mux.Handle(config.MetricsEndpoint, pe)
			err := http.ListenAndServe(fmt.Sprintf(":%s", configService.MetricsSettings.MetricPort), mux)
			if err != nil {
				log.Panicf("Failed to start the Prometheus stats exporter: %v", err)
			}
		}()
	}

	forkContext := rootCtx.CreateFork()

	go func() {
		err := FgBlackbox.Start()
		if err != http.ErrServerClosed {
			log.WithField("error", err).Error("Failed to start fg-blackbox.")
		}
		forkContext.FinishBlocking(blockingTimeout)
	}()

	// start test manager
	testManager.Start()
	defer testManager.Stop()

	FgBlackbox.WaitForShutdown()
	rootCtx.FinishBlocking(blockingTimeout)
}
