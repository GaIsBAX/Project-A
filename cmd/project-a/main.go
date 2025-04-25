package main

import (
	"A/internal/app"
	"A/internal/config"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// Точка входа в приложение
func main() {
	// Подключение конфига
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// Подключение логгера
	log := setupLogger(cfg.Env)
	fmt.Println(log)

	// Подключение приложения
	app := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	// Запуск приложения
	go app.GRPCSrv.MustRun()

	// Остановка приложения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Ожидание сигнала
	sign := <-stop

	log.Info("stopping app", slog.String("signal", sign.String()))

	// Корректная остановка приложения
	app.GRPCSrv.Stop()

	log.Info("app stop")

}

// Настройка логгера
// TODO: вынести в отдельный пакет и костомизировать в будущем
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
