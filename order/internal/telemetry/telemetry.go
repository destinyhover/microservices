package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	Endpoint       string
	Insecure       bool
}

func SetupProvider(ctx context.Context, cfg Config) (func(context.Context) error, error) {
	exp, err := newExporter(ctx, cfg)
	if err != nil {
		return nil, err
	}
	res, err := newResource(ctx, cfg)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exp)),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp.Shutdown, nil
}

func newExporter(ctx context.Context, cfg Config) (sdktrace.SpanExporter, error) {
	if cfg.Endpoint == "" {
		return stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)
	}
	// OTLP/gRPC → Jaeger/Collector
	var opts []otlptracegrpc.Option
	opts = append(opts, otlptracegrpc.WithEndpoint(cfg.Endpoint))
	if cfg.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}
	dctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	exp, err := otlptracegrpc.New(dctx, opts...)
	if err != nil {
		return nil, err
	}
	return exp, nil
}
func newResource(ctx context.Context, cfg Config) (*resource.Resource, error) {
	attrs := []attribute.KeyValue{
		semconv.ServiceName(cfg.ServiceName),
	}
	if cfg.ServiceVersion != "" {
		attrs = append(attrs, semconv.ServiceVersion(cfg.ServiceVersion))
	}
	if cfg.Environment != "" {
		attrs = append(attrs, semconv.DeploymentEnvironmentName(cfg.Environment))
	}
	return resource.New(ctx,
		resource.WithHost(),
		resource.WithSchemaURL(semconv.SchemaURL),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(attrs...))
}
