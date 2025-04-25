package grpcapp

import (
	authgrpc "A/internal/grpc/auth"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	log  *slog.Logger
	gRPC *grpc.Server
	port int
}

// New создает и возвращает новый экземпляр gRPC приложения.
func New(log *slog.Logger, authService authgrpc.Auth, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)

	return &App{
		log:  log,
		gRPC: gRPCServer,
		port: port,
	}
}

// MustRun запускает gRPC сервер, но если возникнет ошибка, то
// вызывает panic.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run запускает gRPC сервер.
func (a *App) Run() error {
	const op = "app/grpc/app.go/run"

	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("starting gRPC server", slog.String("addr", l.Addr().String()))

	if err := a.gRPC.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop останавливает gRPC сервер.
func (a *App) Stop() {
	const op = "app/grpc/app.go/stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC server")

	a.gRPC.GracefulStop()

}
