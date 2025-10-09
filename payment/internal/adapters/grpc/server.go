package grpc

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/destinyhover/microservices-proto/golang/payment"
	"github.com/destinyhover/microservices/payment/config"
	"github.com/destinyhover/microservices/payment/internal/ports"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Adapter struct {
	api    ports.APIPort
	port   int
	server *grpc.Server
	payment.UnimplementedPaymentServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{api: api, port: port}
}

func (a Adapter) Run() {

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		slog.Error("failed to listen", "port", a.port, "err", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))

	a.server = grpcServer
	payment.RegisterPaymentServer(grpcServer, a)
	if config.GetEnv() == "development" {
		reflection.Register(grpcServer)
	}

	slog.Info("starting payment service on", "port", a.port)
	if err := grpcServer.Serve(listen); err != nil {
		slog.Error("failed to serve grpc", "port", a.port, "err", err)
		os.Exit(1)
	}
}

func (a Adapter) Stop() {
	a.server.Stop()
}
