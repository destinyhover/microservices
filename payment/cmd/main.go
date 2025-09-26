package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/destinyhover/microservices/payment/config"
	"github.com/destinyhover/microservices/payment/internal/adapters/db"
	gserver "github.com/destinyhover/microservices/payment/internal/adapters/grpc"
	app "github.com/destinyhover/microservices/payment/internal/application/core/api"
	"github.com/destinyhover/microservices/payment/internal/telemetry"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
	// Building up  tracing
	ctx := context.Background()
	shutdown, err := telemetry.SetupProvider(ctx, telemetry.Config{
		ServiceName:    "payment",
		ServiceVersion: "1.0.0",
		Endpoint:       config.GetOLTPEndpoint(),
		Insecure:       config.IsOTLPInsecure(),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			slog.Error("failed to set up telemetry", "err", err)
			os.Exit(1)
		}
	}()

	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		fmt.Println(err)
		return
	}

	api := app.NewApplication(dbAdapter)

	port, err := strconv.Atoi(config.GetPaymentPort())
	if err != nil {
		fmt.Println(err)
		return
	}

	server := gserver.NewAdapter(api, port)
	slog.Info("booting payment service", "env", config.GetEnv(), "port", port)
	server.Run()
}
