package opentel

import (
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type OpenTel struct {
	ServiceName      string
	ServiceVersion   string
	ExporterEndpoint string
}

func NewOpenTel(serviceName, serviceVersion, exporterEndpoint string) *OpenTel {
	return &OpenTel{
		ServiceName:      serviceName,
		ServiceVersion:   serviceVersion,
		ExporterEndpoint: exporterEndpoint,
	}
}

func (o *OpenTel) GetTracer() trace.Tracer {
	var logger = log.New(os.Stderr, "zipkin-example ", log.Ldate|log.Ltime|log.Llongfile)
	exporter, err := zipkin.New(
		o.ExporterEndpoint,
	)
	if err != nil {
		logger.Fatalf("failed to create Zipkin exporter: %v", err)
	}

	batcher := sdktrace.NewBatchSpanProcessor(exporter)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(o.ServiceName),
			semconv.ServiceVersionKey.String(o.ServiceVersion),
		)),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	tracer := otel.Tracer("io.opentelemetry.traces.goapp")
	return tracer
}
