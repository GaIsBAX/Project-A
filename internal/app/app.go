package app

import (
	"log/slog"
	"time"

	grpcapp "A/internal/app/grpc"
	"A/internal/services/auth"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {

	// todo: инициализировать хранилище

	// todo: инициализировать auth service
	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
