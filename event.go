package event

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
)

var defaultLogger atomic.Pointer[slog.Logger]

// SetDefaultLogger sets the logger used by emitters that were not given one
// with the Logger option. A nil l turns this logging off.
func SetDefaultLogger(l *slog.Logger) {
	defaultLogger.Store(l)
}

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

// New returns an Emitter configured with opts.
func New[T any](opts ...Option) *Emitter[T] {
	o := emitterOptions{}
	for _, opt := range opts {
		opt.apply(&o)
	}
	return &Emitter[T]{logger: o.logger}
}

func (m *Emitter[T]) Handle(h Handler[T]) {
	m.handlers = append(m.handlers, h)
}

func (m *Emitter[T]) HandleFunc(f func(ctx context.Context, e T)) {
	m.handlers = append(m.handlers, HandlerFunc[T](f))
}

func (m *Emitter[T]) Trigger(ctx context.Context, e T) {
	logger := m.logger
	if logger == nil {
		logger = defaultLogger.Load()
	}
	if logger != nil && logger.Enabled(ctx, slog.LevelDebug) {
		logger.LogAttrs(ctx, slog.LevelDebug, "event triggered",
			slog.String("module", "gocrumb/event"),
			slog.String("type", fmt.Sprintf("%T", e)),
			slog.Any("event", e),
		)
	}
	for _, h := range m.handlers {
		h.Handle(ctx, e)
	}
}
