package bootstrap

import (
	"context"
	"time"

	otelpyroscope "github.com/grafana/otel-profiling-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func InitTracer() (trace.TracerProvider, *sdktrace.TracerProvider, error) {
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithURLPath("/v1/traces"),
		otlptracehttp.WithTimeout(5*time.Second),
	)
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		return nil, nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	enhancedTracerProvider := otelpyroscope.NewTracerProvider(tp)
	otel.SetTracerProvider(enhancedTracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))
	return enhancedTracerProvider, tp, err
}
