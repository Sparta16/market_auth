package grpcapp

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net"
	"weather_back_api_getway/internal/controllers"
)

type App struct {
	gRPCServer *grpc.Server
	port       int
}

func New(port int, auth controllers.Auth) *App {

	gRPCServer := grpc.NewServer()

	controllers.Register(gRPCServer, auth)

	return &App{
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return err
	}

	slog.Info("gRPC server created", l)
	slog.Info("Listening on port %d", a.port)

	if err := a.gRPCServer.Serve(l); err != nil {
		slog.Error("Failed to serve gRPC server", l, err)
		return err
	}

	return nil
}

func (a *App) Stop() {

	log.Printf("Stopped gRPC server on port %d", a.port)

	a.gRPCServer.GracefulStop()
}
