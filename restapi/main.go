package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"maydere.com/opentel-labs/restapi/handler"
	"maydere.com/opentel-labs/restapi/pb"
	"maydere.com/opentel-labs/restapi/telem"

	echoSwagger "github.com/swaggo/echo-swagger"
	_ "maydere.com/opentel-labs/restapi/docs"
)

var (
	httpPort                 = 1919
	commonServiceAddr string = "localhost:9002"
	otelAgentAddr     string = "0.0.0.0:4317"
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// flag common-service address
	commonServiceAddrVal, ok := os.LookupEnv("COMMON_SERVICE_ADDR")
	if !ok {
		log.Fatal("COMMON_SERVICE_ADDR is not set")
	}
	commonServiceAddr = commonServiceAddrVal

	// flag otel agent address
	otelAgentAddrVal, ok := os.LookupEnv("OTEL_AGENT_ADDR")
	if !ok {
		log.Fatal("OTEL_AGENT_ADDR is not set")
	}
	otelAgentAddr = otelAgentAddrVal
}

// @title Restapi
// @version 1.0
// @description A set of APIs to allow applications interact to with the Restapi.

// @contact.name Maydere
// @contact.url http://www.maydere.com
// @contact.email destek@maydere.com

// @schemes http https
// @host localhost:1919
// @BasePath /v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	ctx, quit := signal.NotifyContext(context.Background(), os.Interrupt)
	defer quit()

	logWriter := zerolog.SyncWriter(os.Stdout)
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		logWriter = zerolog.NewConsoleWriter()
	}

	logger := zerolog.New(logWriter).With().Timestamp().Logger()

	ctx = logger.WithContext(ctx)

	flag.Parse()

	if err := run(ctx); err != nil {
		logger.Fatal().Stack().Err(err).Msgf("program exited with an error: %+v", err)
	}
}

func run(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)

	// gRPC Client
	// Set up a connection to the server.
	var commonConn grpc.ClientConn
	go func() {
		backoff.Retry(func() error {
			conn, err := grpc.DialContext(ctx,
				commonServiceAddr,
				grpc.WithInsecure(),
				grpc.WithBlock(),
				grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
			)
			if err != nil {
				return err
			} else {
				commonConn = *conn
				return nil
			}
		}, backoff.WithContext(backoff.NewConstantBackOff(time.Second*1), ctx))
	}()

	defer func() {
		emptyConnection := grpc.ClientConn{}
		if !reflect.DeepEqual(commonConn, emptyConnection) {
			commonConn.Close()
		}
	}()

	// echo
	e := echo.New()
	e.Use(middleware.BodyLimit("2M"))
	e.Pre(middleware.MethodOverride())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodPatch},
		AllowHeaders:     []string{},
		AllowCredentials: true,
	}))

	// echo router
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	api := e.Group("/v1")

	h := &handler.Handler{
		CommonServiceClient: pb.NewCommonServiceClient(&commonConn),
	}
	h.RegisterHandlers(api)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: e,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	srvErrCh := make(chan error)
	go func() {
		defer close(srvErrCh)

		if err := srv.ListenAndServe(); err != nil {
			srvErrCh <- errors.WithStack(err)
		}
	}()

	logger.Info().Str("listen_addr", srv.Addr).Str("common_service_addr", commonServiceAddr).Msg("restapi application is starting. v0.0.1")

	// init tracer proviler
	tracerProvider := telem.InitProvider(ctx, otelAgentAddr, "restapi")
	defer tracerProvider()

	select {
	case err := <-srvErrCh:
		return errors.WithStack(err)

	case <-ctx.Done():
		logger.Debug().Msg("graceful shutdown has been started.")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return errors.WithStack(err)
		}

		logger.Debug().Msg("graceful shutdown has been completed.")

		return nil
	}
}
