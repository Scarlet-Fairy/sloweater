package main

import (
	"context"
	"flag"
	"github.com/Scarlet-Fairy/sloweater/pb"
	"github.com/Scarlet-Fairy/sloweater/pkg/config/localStatic"
	nomadOrchestrator "github.com/Scarlet-Fairy/sloweater/pkg/orchestrator/nomad"
	redisPubSub "github.com/Scarlet-Fairy/sloweater/pkg/pubsub/redis"
	redisRepo "github.com/Scarlet-Fairy/sloweater/pkg/repository/redis"
	"github.com/Scarlet-Fairy/sloweater/pkg/service"
	"github.com/Scarlet-Fairy/sloweater/pkg/transport"
	"github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/go-redis/redis/v8"
	consulApi "github.com/hashicorp/consul/api"
	nomadApi "github.com/hashicorp/nomad/api"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net"
	"os"
)

var ctx = context.Background()

func main() {
	fs := flag.NewFlagSet("sloweater", flag.ExitOnError)
	var (
		grpcAddr    = fs.String("grpc-addr", ":8082", "gRPC listen address")
		redisUrl    = flag.String("redis-url", "redis://localhost:6379", "redis url where publish complete events")
		registryUrl = flag.String("registry-url", "localhost:5000", "docker image registry url where cobold push builded images")
		nomadUrl    = flag.String("nomad-url", "localhost:4646", "nomad api url")
		consulUrl   = flag.String("consul-url", "localhost:8500", "consul api url")
		// tracingHost = flag.String("tracing-host", "localhost", "host where send traces")
		// tracingPort = flag.String("tracing-port", "6831", "port of the host where send traces")
	)
	_ = fs.Parse(os.Args[1:])

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	redisClient, err := newRedisClient(*redisUrl)
	if err != nil {
		logger.Log("err", err.Error(), "msg", "Unable to connect to redis", "redis-url", *redisUrl)
		os.Exit(1)
	}

	nomadClient, err := newNomadClient(*nomadUrl)
	if err != nil {
		logger.Log("err", err.Error(), "msg", "Unable to connect to nomad", "nomad-url", *nomadUrl)
		os.Exit(1)
	}

	_, err = newConsulClient(*consulUrl)
	if err != nil {
		logger.Log("err", err.Error(), "msg", "Unable to connect to consul", "consul-url", *consulUrl)
		os.Exit(1)
	}

	configs := localStatic.NewConfig()

	redisRepository := redisRepo.New(redisClient)
	redisPubSub := redisPubSub.New(redisClient, redisRepository)
	nomadOrchestrator := nomadOrchestrator.New(nomadClient, *configs, *registryUrl)

	svc := service.NewService(nomadOrchestrator, redisRepository, redisPubSub, logger)
	endpoints := service.NewEndpoints(svc)
	grpcServer := transport.NewGRPCServer(endpoints, logger)

	var g run.Group
	{
		grpcListener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "gRPC", "addr", *grpcAddr)

			baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
			pb.RegisterSchedulerServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(err error) {
			_ = grpcListener.Close()
		})
	}

	logger.Log("exit", g.Run())
}

func newRedisClient(url string) (*redis.Client, error) {
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, errors.Wrap(err, "redis client")
	}

	client := redis.NewClient(options)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.New("unable to connect to redis")
	}

	return client, nil
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
