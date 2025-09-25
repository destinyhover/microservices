package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/destinyhover/microservices/order/config"
	"github.com/destinyhover/microservices/order/internal/adapters/db"
	gserver "github.com/destinyhover/microservices/order/internal/adapters/grpc"
	"github.com/destinyhover/microservices/order/internal/adapters/payment"
	app "github.com/destinyhover/microservices/order/internal/application/core/api"
	"github.com/destinyhover/microservices/order/internal/telemetry"
)

func main() {
	// (необязательно) аккуратный текстовый хендлер для slog
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// Building up  tracing
	ctx := context.Background()
	shutdown, err := telemetry.SetupProvider(ctx, telemetry.Config{
		ServiceName:    "order",
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

	pPort := config.GetPaymentSourceURL()
	pAdapter, err := payment.NewAdapter(pPort)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 1) Собираем application-сервис, реализующий ports.APIPort
	api := app.NewApplication(dbAdapter, pAdapter)

	// 2) Берём порт из окружения (APPLICATION_PORT)
	port := config.GetApplicationPort()

	// 3) Стартуем gRPC-сервер (твой адаптер)
	server := gserver.NewAdapter(api, port)
	slog.Info("booting order service", "env", config.GetEnv(), "port", port)
	server.Run()
}
