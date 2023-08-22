package telem

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"google.golang.org/grpc"
)

// Initializes an OTLP exporter, and configures the corresponding trace and metric providers.
func InitProvider(ctx context.Context, otelAgentAddr string, serviceNameKey string) func() {
	log.Println("opentelemetry initializing 1/10:", otelAgentAddr, serviceNameKey)

	metricClient := otlpmetricgrpc.NewClient(
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(otelAgentAddr))

	log.Println("opentelemetry initializing 2/10: exporter metric client initializing")

	metricExp, err := otlpmetric.New(ctx, metricClient)
	if err != nil {
		log.Fatalf("failed to create the collector metric exporter: %s", err)
	}

	log.Println("opentelemetry initializing 3/10: exporter metric controller initializing")

	pusher := controller.New(
		processor.NewFactory(simple.NewWithHistogramDistribution(), metricExp),
		controller.WithExporter(metricExp),
		controller.WithCollectPeriod(2*time.Second),
	)
	global.SetMeterProvider(pusher)

	log.Println("opentelemetry initializing 4/10: exporter metric controller starting")

	if err = pusher.Start(ctx); err != nil {
		log.Fatalf("failed to start metric pusher: %s", err)
	}

	log.Println("opentelemetry initializing 5/10: grpc trace client initializing")

	traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(otelAgentAddr),
		otlptracegrpc.WithDialOption(grpc.WithBlock()))

	log.Println("opentelemetry initializing 6/10: grpc trace exporter initializing -> ", otelAgentAddr)

	traceExp, err := otlptrace.New(ctx, traceClient)
	if err != nil {
		log.Fatalf("failed to create the collector trace exporter: %s", err)
	}

	log.Println("opentelemetry initializing 7/10: trace resource configuring")

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceNameKey),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %s", err)
	}

	log.Println("opentelemetry initializing 8/10: tracer provider initializing")

	bsp := sdktrace.NewBatchSpanProcessor(traceExp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	log.Println("opentelemetry initializing 9/10: tracer provider initialized")

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tracerProvider)

	log.Println("opentelemetry initializing 10/10: opentelemetry initialized")

	return func() {
		log.Println("opentelemetry starting")

		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		if err := traceExp.Shutdown(cxt); err != nil {
			log.Fatalf("failed to otlp shutdown: %s", err)
			otel.Handle(err)
		}
		// pushes any last exports to the receiver
		if err := pusher.Stop(cxt); err != nil {
			log.Fatalf("failed to pusher stop: %s", err)
			otel.Handle(err)
		}

		log.Println("opentelemetry started")
	}
}
