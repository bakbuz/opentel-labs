package telem

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
)

// Initializes an OTLP exporter, and configures the corresponding trace and metric providers.
func InitProvider(ctx context.Context, otelAgentAddr string, serviceNameKey string) func() {
	// step 1
	log.Println("opentelemetry initializing 1/8:", otelAgentAddr, serviceNameKey)

	// step 2
	log.Println("opentelemetry initializing 2/8: resource initializing")
	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String(serviceNameKey),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %s", err)
	}

	// step 3
	log.Println("opentelemetry initializing 3/8: metric exporter initializing")

	metricExp, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(otelAgentAddr),
	)
	if err != nil {
		log.Fatalf("failed to create the collector metric exporter: %s", err)
	}

	// step 4
	log.Println("opentelemetry initializing 4/8: meter provider initializing")

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				metricExp,
				sdkmetric.WithInterval(2*time.Second),
			),
		),
	)
	otel.SetMeterProvider(meterProvider)

	// step 5
	log.Println("opentelemetry initializing 5/8: grpc trace client initializing")

	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otelAgentAddr),
		otlptracegrpc.WithDialOption(grpc.WithBlock()))
	sctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	// step 6
	log.Println("opentelemetry initializing 6/8: trace exporter initializing")

	traceExp, err := otlptrace.New(sctx, traceClient)
	if err != nil {
		log.Fatalf("failed to create the collector trace exporter: %s", err)
	}

	// step 7
	log.Println("opentelemetry initializing 7/8: tracer provider initializing")

	bsp := sdktrace.NewBatchSpanProcessor(traceExp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(tracerProvider)

	// step 8
	log.Println("opentelemetry initializing 8/8: opentelemetry initialized")

	return func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := traceExp.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
		// pushes any last exports to the receiver
		if err := meterProvider.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}
}
