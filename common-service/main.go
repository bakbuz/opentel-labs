package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"maydere.com/opentel-labs/common-service/handler"
	"maydere.com/opentel-labs/common-service/pb"
	"maydere.com/opentel-labs/common-service/telem"
)

var (
	// an example app flag with the default value
	tcpPort       int    = 9002
	otelAgentAddr string = "0.0.0.0:4317"
)

func init() {
	// print error stack to the log messages
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// flag tcp port
	tcpPortVal, ok := os.LookupEnv("TCP_PORT")
	if !ok {
		log.Fatal("TCP_PORT is not set")
	}
	tcpPortInt, err := strconv.Atoi(tcpPortVal)
	if err != nil {
		log.Fatal("TCP_PORT is not a valid number", err)
	}
	tcpPort = tcpPortInt

	// flag otel agent address
	otelAgentAddrVal, ok := os.LookupEnv("OTEL_AGENT_ADDR")
	if !ok {
		log.Fatal("OTEL_AGENT_ADDR is not set")
	}
	otelAgentAddr = otelAgentAddrVal
}

func main() {
	// this is the root context for entire app, we will cancel the root context when user press ctrl+c
	ctx, quit := signal.NotifyContext(context.Background(), os.Interrupt)
	defer quit()

	// send all log messages to stdout
	logWriter := zerolog.SyncWriter(os.Stdout)

	// if the program starts on the console use colored log writer
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		logWriter = zerolog.NewConsoleWriter()
	}

	// init logger
	logger := zerolog.New(logWriter).With().Timestamp().Logger()

	// add logger to context
	ctx = logger.WithContext(ctx)

	// parse the given app flags
	flag.Parse()

	// execute app
	if err := run(ctx); err != nil {
		logger.Fatal().Stack().Err(err).Msgf("program exited with an error: %+v", err)
	}
}

// run is the entry point for the app, it should live in this function
func run(ctx context.Context) error {
	// get logger from the context
	log := zerolog.Ctx(ctx)

	// an example log output
	log.Info().Int("tcp_port", tcpPort).Str("otel_agent_addr", otelAgentAddr).Msg("common-service application is starting. v0.0.1")

	// init tracer provider
	tracerProvider := telem.InitProvider(otelAgentAddr, "common-service")
	defer tracerProvider()

	// handler
	srv := handler.NewHandler("")

	// run
	errChConsumer := make(chan error)
	errChRPC := make(chan error)

	grpcCtx, grpcCancel := context.WithCancel(ctx)
	defer grpcCancel()

	go func(ctx context.Context, errCh chan error) {
		defer close(errCh)

		lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", tcpPort))
		if err != nil {
			errCh <- errors.WithStack(err)
			return
		}

		grpcServer := grpc.NewServer(
			grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
			//grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		)

		pb.RegisterCommonServiceServer(grpcServer, srv)

		log.Debug().Msgf("server listening at %v", lis.Addr())

		go func() {
			<-ctx.Done()
			grpcServer.Stop()
		}()

		if err := grpcServer.Serve(lis); err != nil {
			errCh <- errors.WithStack(err)
			return
		}
	}(grpcCtx, errChRPC)

	select {
	case <-ctx.Done():
		return errors.WithStack(ctx.Err())

	case err := <-errChConsumer:
		return errors.WithStack(err)

	case err := <-errChRPC:
		return errors.WithStack(err)
	}
}
