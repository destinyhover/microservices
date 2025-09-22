package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/destinyhover/microservices/order/config"
	"github.com/destinyhover/microservices/order/internal/adapters/db"
	gserver "github.com/destinyhover/microservices/order/internal/adapters/grpc"
	app "github.com/destinyhover/microservices/order/internal/application/core/api"
	"github.com/destinyhover/microservices/order/internal/application/fakes"
)

func main() {
	// (необязательно) аккуратный текстовый хендлер для slog
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		fmt.Println(err)
		return
	}
	pp := fakes.NoopPayment{}
	// 1) Собираем application-сервис, реализующий ports.APIPort
	api := app.NewApplication(dbAdapter, pp)

	// 2) Берём порт из окружения (APPLICATION_PORT)
	port := config.GetApplicationPort()

	// 3) Стартуем gRPC-сервер (твой адаптер)
	server := gserver.NewAdapter(api, port)
	slog.Info("booting order service", "env", config.GetEnv(), "port", port)
	server.Run()
}
