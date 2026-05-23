package event

import (
	"context"
	"fmt"
	"log/slog"
)

type Handler[T any] interface {
	Handle(ctx context.Context, e T)
}

type HandlerFunc[T any] func(ctx context.Context, e T)

func (f HandlerFunc[T]) Handle(ctx context.Context, e T) {
	f(ctx, e)
}

type Emitter[T any] struct {
	handlers []Handler[T]
	logger   *slog.Logger
}

func (m *Emitter[T]) Handle(h Handler[T]) {
	m.handlers = append(m.handlers, h)
}

func (m *Emitter[T]) HandleFunc(f func(ctx context.Context, e T)) {
	m.handlers = append(m.handlers, HandlerFunc[T](f))
}

func (m *Emitter[T]) Trigger(ctx context.Context, e T) {
	if m.logger != nil && m.logger.Enabled(ctx, slog.LevelDebug) {
		m.logger.LogAttrs(ctx, slog.LevelDebug, "event triggered",
			slog.String("module", "gocrumb/event"),
			slog.String("type", fmt.Sprintf("%T", e)),
			slog.Any("event", e),
		)
	}
	for _, h := range m.handlers {
		h.Handle(ctx, e)
	}
}
