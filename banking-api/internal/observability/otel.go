package observability

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const ServiceName = "banking-api"

func Init() *sdktrace.TracerProvider {

	res, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(ServiceName),
		),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return tp
}
