package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/twistingmercury/monitoring-example/handlers"

	"github.com/rs/zerolog"

	"github.com/twistingmercury/monitoring/health"
	"github.com/twistingmercury/monitoring/logs"
	"github.com/twistingmercury/monitoring/metrics"
	"github.com/twistingmercury/monitoring/traces"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // please don't use in PRODUTION environments!!!

	"github.com/gin-gonic/gin"
)

var buildDate = time.Now().String()

const (
	HTTP               = "http"
	GRPC               = "grpc"
	OtelProtocolEnvVar = "OTEL_PROTOCOL"
	OtelCollectorEP    = "OTEL_COLLECTOR_EP"
	apiVersion         = "1.0.0"
	serviceName        = "monex"
	commit             = "123456"
	env                = "local"
	metricsPort        = "9090"
)

func main() {
	// 1: initialize the logger first, so we can log any errors that occur during initialization
	logs.Initialize(zerolog.DebugLevel, apiVersion, serviceName, buildDate, commit, env, os.Stdout)
	otelProtocol := os.Getenv(OtelProtocolEnvVar)
	if len(otelProtocol) == 0 {
		//logs.Fatal("OTEL_PROTOCOL is not set - shutting down")
		log.Logger.Fatal().Msg("OTEL_PROTOCOL is not set - shutting down")
	}

	// 2: create trace.SpanExporter so that we can initialize the tracer
	exporter, err := newSpanExporter(otelProtocol)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("failed to initialize the exporter - shutting down")
	}

	// 3: initialize the tracer using the exporter from step 2
	shutdown, err := traces.Initialize(exporter, apiVersion, serviceName, buildDate, commit, env)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("failed to initialize the tracer - shutting down")
	}

	// 4: defer the shutdown of the exporter, tracer, and Prometheus app metrics
	defer func() {
		ctx := context.Background()
		exporter.Shutdown(ctx)
		shutdown(ctx)
	}()

	// 5: initialize the Prometheus app metrics
	metrics.Initialize(metricsPort, "example")
	metrics.Publish()

	// 6: initialize the healthcheck dependencies
	deps := []health.DependencyDescriptor{
		{Name: "Golang Site", Type: "Website", Connection: "https://golang.org/"},
		{Name: "sql dB check", Type: "database", HandlerFunc: handlers.CheckMSSQL},
	}

	// note: in production, you should use gin.ReleaseMode
	gin.SetMode(gin.DebugMode)

	// create a new gin router with no middleware; we'll add our own
	router := gin.New()

	// finally, configure gin to use the middleware
	router.Use(
		gin.Recovery(),
		metrics.GinMetricsMiddleWare(),
		traces.GinTracingMiddleware(),
		logs.GinLoggingMiddleware())

	router.GET("/ping", handlers.PingHandler)
	router.GET("/pong", handlers.PongHandler)
	router.GET("/healthcheck", health.Handler("examples", deps...))

	log.Logger.
		Info().
		Int("port", 8080).Msg("starting server")

	if err := router.Run(":8080"); err != nil {
		log.Logger.Fatal().Err(err).Msg("failed to start the server")
	}
}

// newSpanExporter creates a new span exporter based on the protocol
func newSpanExporter(protocol string) (exporter trace.SpanExporter, err error) {
	otelCollectorEP := os.Getenv(OtelCollectorEP)
	if len(otelCollectorEP) == 0 {
		err = errors.New("OTEL_COLLECTOR_EP is not set - shutting down")
		return
	}

	exCtx := context.Background()
	switch strings.ToLower(protocol) {
	case GRPC:
		log.Logger.Debug().Msg("using grpc exporter")
		grpcConn, _ := grpc.Dial(otelCollectorEP, grpc.WithTransportCredentials(insecure.NewCredentials()))
		exporter, err = otlptracegrpc.New(exCtx, otlptracegrpc.WithGRPCConn(grpcConn))
	case HTTP:
		log.Logger.Debug().Msg("using http exporter")
		exporter, err = traces.NewHTTPExporter(exCtx, otelCollectorEP, otlptracehttp.WithInsecure())
	default:
		err = fmt.Errorf("invalid protocol: %s", protocol)
	}

	return
}
