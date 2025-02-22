package app

import (
	"log/slog"

	grpcapp "github.com/nikita-reshetnyak/analytics/internal/app/grpc"
	"github.com/nikita-reshetnyak/analytics/internal/services/analytics"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int) *App {
	analyticsService := analytics.New(log)
	grpcApp := grpcapp.New(log, analyticsService, grpcPort)
	return &App{GRPCServer: grpcApp}
}
