package main

import (
	"context"
	"flag"
	"github.com/Scarlet-Fairy/sloweater/pb"
	"github.com/Scarlet-Fairy/sloweater/pkg/config/localStatic"
	"github.com/Scarlet-Fairy/sloweater/pkg/endpoint"
	nomadOrchestrator "github.com/Scarlet-Fairy/sloweater/pkg/orchestrator/nomad"
	"github.com/Scarlet-Fairy/sloweater/pkg/service"
	grpcTransport "github.com/Scarlet-Fairy/sloweater/pkg/transport/grpc"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	consulApi "github.com/hashicorp/consul/api"
	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/oklog/run"
	"google.golang.org/grpc"
	"net"
	"os"
	"strings"
)

var ctx = context.Background()

var (
	grpcAddr    = flag.String("grpc-addr", ":8082", "gRPC server listen address")
	registryUrl = flag.String("registry-url", "localhost:5000", "docker image registry url where cobold push builded images")
	nomadUrl    = flag.String("nomad-url", "http://localhost:4646", "nomad api url")
	consulUrl   = flag.String("consul-url", "http://localhost:8500", "consul api url")
	// tracingHost = flag.String("tracing-host", "localhost", "host where send traces")
	// tracingPort = flag.String("tracing-port", "6831", "port of the host where send traces")
)

func main() {
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
		logger = level.NewInjector(logger, level.InfoValue())
	}

	debugLogger := level.Debug(logger)
	errorLogger := level.Error(logger)

	debugLogger.Log("args", strings.Join(os.Args[1:], " "))

	nomadClient, err := newNomadClient(*nomadUrl)
	if err != nil {
		errorLogger.Log("err", err.Error(), "msg", "Unable to connect to nomad", "nomad-url", *nomadUrl)
		os.Exit(1)
	}

	consulClient, err := newConsulClient(*consulUrl)
	if err != nil {
		errorLogger.Log("err", err.Error(), "msg", "Unable to connect to consul", "consul-url", *consulUrl)
		os.Exit(1)
	}

	configs := localStatic.NewConfig()

	nomadOrchestrator := nomadOrchestrator.New(
		nomadClient,
		consulClient,
		*configs,
		log.With(
			logger,
			"component", "orchestrator",
			"layer", "service",
		),
		*registryUrl,
	)

	svc := service.NewService(nomadOrchestrator, log.With(logger, "layer", "service"))
	endpoints := endpoint.NewEndpoints(svc, log.With(logger, "layer", "endpoint"))
	grpcServer := grpcTransport.NewGRPCServer(endpoints, log.With(logger, "layer", "transport"))

	var g run.Group
	{
		grpcListener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errorLogger.Log(
				"layer", "transport",
				"transport", "gRPC",
				"during", "Listen",
				"err", err,
			)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log(
				"layer", "transport",
				"transport", "gRPC",
				"addr", *grpcAddr,
			)

			baseServer := grpc.NewServer(
				grpc.UnaryInterceptor(kitgrpc.Interceptor),
			)
			pb.RegisterSchedulerServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(err error) {
			_ = grpcListener.Close()
		})
	}

	logger.Log("exit", g.Run())
}

func newNomadClient(url string) (*nomadApi.Client, error) {
	config := nomadApi.DefaultConfig()
	config.Address = url

	return nomadApi.NewClient(config)
}

func newConsulClient(url string) (*consulApi.Client, error) {
	config := consulApi.DefaultConfig()
	config.Address = url

	return consulApi.NewClient(config)
}
