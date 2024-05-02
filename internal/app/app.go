package app

import (
	"time"
	grpcapp "weather_back_api_getway/internal/app/grpc"
	"weather_back_api_getway/internal/repositories"
	"weather_back_api_getway/internal/services"
)

type App struct {
	GRPCServer *grpcapp.App
	db         *repositories.Database
	tokenTTL   time.Duration
}

func New(
	gRPCPort int,
	db *repositories.Database,
	tokenTTL time.Duration,
) *App {

	service := services.New(db, db, db, tokenTTL)

	gRPCServer := grpcapp.New(gRPCPort, service)

	return &App{
		GRPCServer: gRPCServer,
		db:         db,
		tokenTTL:   tokenTTL,
	}
}
