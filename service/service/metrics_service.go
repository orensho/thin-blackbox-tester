package service

import (
	"context"

	"github.com/orensho/thin-slack-blackbox-tester/service/config"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

//go:generate mockery -name MetricsServiceInterface -inpkg -case=underscore -output MockMetricsServiceInterface
type MetricsServiceInterface interface {
	ReportStepTestSuccess(ctx context.Context, flowName string, stepName string, proxyName string, customerName string) error
	ReportStepTestError(ctx context.Context, flowName string, stepName string, proxyName string, customerName string) error
	ReportStepTestTimeout(ctx context.Context, flowName string, stepName string, proxyName string, customerName string) error
	ReportStepTestDuration(ctx context.Context, flowName string, ms float64, stepName string, proxyName string, customerName string) error
}

var (
	keyFlow        = tag.MustNewKey("flow")
	keyStep        = tag.MustNewKey("step")
	keyProxy       = tag.MustNewKey("proxy")
	keyCustomer    = tag.MustNewKey("customer")
	keyEnvironment = tag.MustNewKey("environment")
)

type metricsService struct {
	settings          *config.MetricsSettings
	testsStepErrors   *stats.Int64Measure
	testsStepSuccess  *stats.Int64Measure
	testsStepTimeout  *stats.Int64Measure
	testsStepDuration *stats.Float64Measure
}

func NewMetricsService(settings *config.MetricsSettings) (MetricsServiceInterface, error) {
	s := &metricsService{
		settings: settings,
	}

	err := s.initMetrics()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *metricsService) ReportStepTestSuccess(ctx context.Context, flowName string, stepName string, proxyName string, customerName string) error { //nolint // line length
	ctx, err := s.createStepMeasurementContext(ctx, flowName, stepName, proxyName, customerName)
	if err != nil {
		return errors.Wrap(err, "Failed setting tags on context")
	}

	stats.Record(ctx, s.testsStepSuccess.M(1))

	return nil
}

func (s *metricsService) ReportStepTestError(ctx context.Context, flowName string, stepName string, proxyName string, customerName string) error { //nolint // line length
	ctx, err := s.createStepMeasurementContext(ctx, flowName, stepName, proxyName, customerName)
	if err != nil {
		return errors.Wrap(err, "Failed setting tags on context")
	}

	stats.Record(ctx, s.testsStepErrors.M(1))

	return nil
}

func (s *metricsService) ReportStepTestTimeout(ctx context.Context, flowName string, stepName string, proxyName string, customerName string) error { //nolint // line length
	ctx, err := s.createStepMeasurementContext(ctx, flowName, stepName, proxyName, customerName)
	if err != nil {
		return errors.Wrap(err, "Failed setting tags on context")
	}

	stats.Record(ctx, s.testsStepTimeout.M(1))

	return nil
}

func (s *metricsService) ReportStepTestDuration(ctx context.Context, flowName string, ms float64, stepName string, proxyName string, customerName string) error { //nolint // line length
	ctx, err := s.createStepMeasurementContext(ctx, flowName, stepName, proxyName, customerName)
	if err != nil {
		return errors.Wrap(err, "Failed setting tags on context")
	}

	stats.Record(ctx, s.testsStepDuration.M(ms))

	return nil
}

func (s *metricsService) createStepMeasurementContext(ctx context.Context, flowName string, stepName string, proxyName string, customerName string) (context.Context, error) { //nolint // line length
	return tag.New(ctx,
		tag.Upsert(keyFlow, flowName),
		tag.Upsert(keyStep, stepName),
		tag.Upsert(keyProxy, proxyName),
		tag.Upsert(keyCustomer, customerName),
		tag.Upsert(keyEnvironment, s.settings.Environment))
}

func (s *metricsService) initMetrics() error {
	s.testsStepDuration = stats.Float64("tests/latency", "The latency in milliseconds per test step", stats.UnitMilliseconds)
	s.testsStepErrors = stats.Int64("tests/errors", "The number of step errors", stats.UnitDimensionless)
	s.testsStepTimeout = stats.Int64("tests/timeouts", "The number of step timeouts", stats.UnitDimensionless)
	s.testsStepSuccess = stats.Int64("tests/success", "The number of step successes", stats.UnitDimensionless)

	latencyStepView := &view.View{
		Name:        "step_latency_distribution",
		Measure:     s.testsStepDuration,
		Description: "The distribution of the flows latencies",

		// Latency in buckets:
		// [>=0ms, >=500ms, >=1s, >=10s, >=30s, >=60s, >=90s, >=120s]
		//nolint:gomnd //false positive
		Aggregation: view.Distribution(500, 1000, 10000, 30000, 60000, 90000, 120000),
		TagKeys:     []tag.Key{keyFlow, keyStep, keyProxy, keyCustomer, keyEnvironment},
	}

	errorStepCountView := &view.View{
		Name:        "step_errors_counter",
		Measure:     s.testsStepErrors,
		Description: "The number of failed steps",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{keyFlow, keyStep, keyProxy, keyCustomer, keyEnvironment},
	}

	successStepCountView := &view.View{
		Name:        "step_success_counter",
		Measure:     s.testsStepSuccess,
		Description: "The number of successful steps",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{keyFlow, keyStep, keyProxy, keyCustomer, keyEnvironment},
	}

	timeoutStepCountView := &view.View{
		Name:        "step_timeout_counter",
		Measure:     s.testsStepTimeout,
		Description: "The number of steps timeouts",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{keyFlow, keyStep, keyProxy, keyCustomer, keyEnvironment},
	}

	// Register the views
	if err := view.Register(latencyStepView, errorStepCountView, successStepCountView, timeoutStepCountView); err != nil {
		return errors.Wrap(err, "Failed to register views")
	}

	return nil
}
