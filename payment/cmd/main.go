package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/destinyhover/microservices/payment/config"
	"github.com/destinyhover/microservices/payment/internal/adapters/db"
	gserver "github.com/destinyhover/microservices/payment/internal/adapters/grpc"
	app "github.com/destinyhover/microservices/payment/internal/application/core/api"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

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
