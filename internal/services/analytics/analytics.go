package analytics

import (
	"context"
	"log/slog"
	"time"
)

type Analytics struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Analytics {
	return &Analytics{log: log}
}
func (a *Analytics) SendEvent(ctx context.Context, name string, date time.Time) error {
	const op = "Analytics.SendEvent"
	a.log.With(
		slog.String("op", op),
		slog.String("name", name),
		slog.String("date", date.String()),
	)
	return nil
}
