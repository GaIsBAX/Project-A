package app

import (
	"log/slog"
	"time"

	grpcApp "A/internal/app/grpc"
)

type App struct {
	GRPCSrv *grpcApp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {

	// todo: инициализировать хранилище

	// todo: инициализировать auth service

	grpcApp := grpcApp.New(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
