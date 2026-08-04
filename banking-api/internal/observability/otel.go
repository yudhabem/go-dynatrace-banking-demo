package observability

import (
	"context"
	"errors"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

const defaultServiceName = "banking-api"

type Providers struct {
	loggerProvider *log.LoggerProvider
	meterProvider  *metric.MeterProvider
	tracerProvider *trace.TracerProvider
}

func Init(ctx context.Context) (*Providers, error) {
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == "" {
		return &Providers{}, nil
	}

	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = defaultServiceName
	}

	res := resource.NewWithAttributes("", attribute.String("service.name", serviceName))

	traceExporter, err := otlptracehttp.New(ctx)
	if err != nil {
		return nil, err
	}

	metricExporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return nil, err
	}

	logExporter, err := otlploghttp.New(ctx)
	if err != nil {
		return nil, err
	}

	providers := &Providers{
		tracerProvider: trace.NewTracerProvider(
			trace.WithBatcher(traceExporter),
			trace.WithResource(res),
		),
		meterProvider: metric.NewMeterProvider(
			metric.WithReader(metric.NewPeriodicReader(metricExporter)),
			metric.WithResource(res),
		),
		loggerProvider: log.NewLoggerProvider(
			log.WithProcessor(log.NewBatchProcessor(logExporter)),
			log.WithResource(res),
		),
	}

	otel.SetTracerProvider(providers.tracerProvider)
	otel.SetMeterProvider(providers.meterProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return providers, nil
}

func (p *Providers) LoggerProvider() *log.LoggerProvider {
	return p.loggerProvider
}

func (p *Providers) Shutdown(ctx context.Context) error {
	if p == nil {
		return nil
	}

	var errs []error
	if p.loggerProvider != nil {
		errs = append(errs, p.loggerProvider.Shutdown(ctx))
	}
	if p.meterProvider != nil {
		errs = append(errs, p.meterProvider.Shutdown(ctx))
	}
	if p.tracerProvider != nil {
		errs = append(errs, p.tracerProvider.Shutdown(ctx))
	}

	return errors.Join(errs...)
}
