package grpc

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/destinyhover/microservices-proto/golang/order"
	"github.com/destinyhover/microservices/order/config"
	"github.com/destinyhover/microservices/order/internal/ports"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Adapter struct {
	api    ports.APIPort
	port   int
	server *grpc.Server
	order.UnimplementedOrderServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{api: api, port: port}
}

func (a *Adapter) Run() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		slog.Error("failed to listen", "port", a.port, "err", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))

	a.server = grpcServer
	order.RegisterOrderServer(grpcServer, a)
	if config.GetEnv() == "development" {
		reflection.Register(grpcServer)
	}

	slog.Info("starting order service on", "port", a.port)
	if err := grpcServer.Serve(listen); err != nil {
		slog.Error("gRPC serve failed", "port", a.port, "err", err)
		os.Exit(1)
	}
}

func (a Adapter) Stop() {
	if a.server != nil {
		a.server.GracefulStop()
	}
}
