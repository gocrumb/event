package event

import (
	"context"
	"testing"
)

type orderPlaced struct {
	OrderID int
}

func TestTriggerInvokesHandlersInOrder(t *testing.T) {
	var em Emitter[orderPlaced]
	var calls []string
	em.HandleFunc(func(ctx context.Context, e orderPlaced) {
		calls = append(calls, "a")
	})
	em.HandleFunc(func(ctx context.Context, e orderPlaced) {
		calls = append(calls, "b")
	})
	em.Trigger(context.Background(), orderPlaced{OrderID: 1})
	if got, want := calls, []string{"a", "b"}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("calls = %v, want %v", got, want)
	}
}

func TestTriggerPassesEvent(t *testing.T) {
	var em Emitter[orderPlaced]
	var got orderPlaced
	em.HandleFunc(func(ctx context.Context, e orderPlaced) {
		got = e
	})
	em.Trigger(context.Background(), orderPlaced{OrderID: 42})
	if got.OrderID != 42 {
		t.Errorf("OrderID = %d, want 42", got.OrderID)
	}
}

func TestTriggerPassesContext(t *testing.T) {
	type ctxKey struct{}
	var em Emitter[orderPlaced]
	var got any
	em.HandleFunc(func(ctx context.Context, e orderPlaced) {
		got = ctx.Value(ctxKey{})
	})
	ctx := context.WithValue(context.Background(), ctxKey{}, "v")
	em.Trigger(ctx, orderPlaced{})
	if got != "v" {
		t.Errorf("ctx value = %v, want %q", got, "v")
	}
}

type recordingHandler struct {
	called bool
}

func (h *recordingHandler) Handle(ctx context.Context, e orderPlaced) {
	h.called = true
}

func TestHandleRegistersHandler(t *testing.T) {
	var em Emitter[orderPlaced]
	h := &recordingHandler{}
	em.Handle(h)
	em.Trigger(context.Background(), orderPlaced{})
	if !h.called {
		t.Error("handler was not invoked")
	}
}

func TestTriggerOnZeroValueEmitterIsNoop(t *testing.T) {
	var em Emitter[orderPlaced]
	em.Trigger(context.Background(), orderPlaced{})
}
