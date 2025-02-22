package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	analyticsgrpc "github.com/nikita-reshetnyak/analytics/internal/grpc/analytics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log        *slog.Logger
	grpcServer *grpc.Server
	port       int
}

func New(log *slog.Logger, analyticsService analyticsgrpc.Analytics, port int) *App {
	recoverOpt := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))
			return status.Errorf(codes.Internal, "internal error")
		}),
	}
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			//logging.StartCall, logging.FinishCall,
			logging.PayloadReceived, logging.PayloadSent,
		),
		// Add any other option (check functions starting with logging.With).
	}
	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoverOpt...),
		logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
	))
	analyticsgrpc.Register(gRPCServer, analyticsService)

	return &App{
		log:        log,
		port:       port,
		grpcServer: gRPCServer,
	}
}
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}
func (a *App) Run() error {
	const op = "grpcapp.Run"
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	a.log.Info("grpc server is started", slog.String("addr", listener.Addr().String()))
	if err := a.grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
func (a *App) Stop() {
	const op = "grpcapp.Stop"
	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))
	a.grpcServer.GracefulStop()
}
