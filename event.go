package event

import "context"

type Handler[T any] interface {
	Handle(ctx context.Context, e T)
}

type HandlerFunc[T any] func(ctx context.Context, e T)

func (f HandlerFunc[T]) Handle(ctx context.Context, e T) {
	f(ctx, e)
}

type Emitter[T any] struct {
	handlers []Handler[T]
}

func (m *Emitter[T]) Handle(h Handler[T]) {
	m.handlers = append(m.handlers, h)
}

func (m *Emitter[T]) HandleFunc(f func(ctx context.Context, e T)) {
	m.handlers = append(m.handlers, HandlerFunc[T](f))
}

func (m *Emitter[T]) Trigger(ctx context.Context, e T) {
	for _, h := range m.handlers {
		h.Handle(ctx, e)
	}
}
