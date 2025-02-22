package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nikita-reshetnyak/analytics/internal/app"
	"github.com/nikita-reshetnyak/analytics/internal/config"
)

var (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	config := config.MustLoad()
	logger := setupLogger(config.Env)
	application := app.New(logger, config.GRPC.Port)
	go func() {
		application.GRPCServer.MustRun()
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	application.GRPCServer.Stop()
	logger.Info("Gracefully stopped")

}
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}
	return log
}
