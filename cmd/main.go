package main

import (
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"weather_back_api_getway/config"
	"weather_back_api_getway/internal/app"
	"weather_back_api_getway/internal/repositories"
)

func main() {

	config := config.MustLoad()

	db := repositories.New(config.Database)
	defer repositories.Stop(db)

	application := app.New(config.GRPCConfig.Port, db, config.TokenTTL)
	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sign := <-stop

	slog.Info("stopping server:", sign)

	application.GRPCServer.Stop()

	slog.Info("Successfully stopped server")
}
