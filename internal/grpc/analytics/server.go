package grpcanalytics

import (
	"context"
	"time"

	v1 "github.com/nikita-reshetnyak/analytics-protos/gen/go/analytics/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Analytics interface {
	SendEvent(ctx context.Context, name string, date time.Time)
}
type serverApi struct {
	v1.UnimplementedAnalyticsServer
	analytics Analytics
}

func Register(gRPCServer *grpc.Server, analytics Analytics) {
	v1.RegisterAnalyticsServer(gRPCServer, &serverApi{analytics: analytics})
}
func (s *serverApi) SendEvent(ctx context.Context, eq *v1.Event) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
