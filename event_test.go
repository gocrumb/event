package event

import (
	"context"
	"log/slog"
	"testing"
)

type captureHandler struct {
	enabled bool
	records []slog.Record
}

func (h *captureHandler) Enabled(context.Context, slog.Level) bool { return h.enabled }
func (h *captureHandler) Handle(_ context.Context, r slog.Record) error {
	h.records = append(h.records, r)
	return nil
}
func (h *captureHandler) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h *captureHandler) WithGroup(string) slog.Handler      { return h }

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

func TestTriggerLogsDebugRecord(t *testing.T) {
	h := &captureHandler{enabled: true}
	em := New[orderPlaced](Logger(slog.New(h)))
	em.Trigger(context.Background(), orderPlaced{OrderID: 7})

	if len(h.records) != 1 {
		t.Fatalf("records = %d, want 1", len(h.records))
	}
	r := h.records[0]
	if r.Level != slog.LevelDebug {
		t.Errorf("level = %v, want Debug", r.Level)
	}
	if r.Message != "event triggered" {
		t.Errorf("message = %q, want %q", r.Message, "event triggered")
	}

	attrs := map[string]slog.Value{}
	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value
		return true
	})
	if got := attrs["module"].String(); got != "gocrumb/event" {
		t.Errorf("module attr = %q, want %q", got, "gocrumb/event")
	}
	if got := attrs["type"].String(); got != "event.orderPlaced" {
		t.Errorf("type attr = %q, want %q", got, "event.orderPlaced")
	}
	if _, ok := attrs["event"]; !ok {
		t.Errorf("missing event attr")
	}
}

func TestTriggerSkipsLoggingWhenDebugDisabled(t *testing.T) {
	h := &captureHandler{enabled: false}
	em := Emitter[orderPlaced]{logger: slog.New(h)}
	em.Trigger(context.Background(), orderPlaced{})
	if len(h.records) != 0 {
		t.Errorf("records = %d, want 0", len(h.records))
	}
}

func TestTriggerStillInvokesHandlersWhenLogging(t *testing.T) {
	h := &captureHandler{enabled: true}
	em := New[orderPlaced](Logger(slog.New(h)))
	called := false
	em.HandleFunc(func(ctx context.Context, e orderPlaced) {
		called = true
	})
	em.Trigger(context.Background(), orderPlaced{})
	if !called {
		t.Error("handler was not invoked")
	}
}

func TestTriggerDoesNotTouchDefaultLoggerWhenNil(t *testing.T) {
	h := &captureHandler{enabled: true}
	prev := slog.Default()
	slog.SetDefault(slog.New(h))
	t.Cleanup(func() { slog.SetDefault(prev) })

	var em Emitter[orderPlaced]
	em.Trigger(context.Background(), orderPlaced{})
	if len(h.records) != 0 {
		t.Errorf("records = %d, want 0 (nil Logger must not fall back to slog.Default)", len(h.records))
	}
}

func TestTriggerUsesDefaultLogger(t *testing.T) {
	h := &captureHandler{enabled: true}
	SetDefaultLogger(slog.New(h))
	t.Cleanup(func() { SetDefaultLogger(nil) })

	var em Emitter[orderPlaced]
	em.Trigger(context.Background(), orderPlaced{})
	if len(h.records) != 1 {
		t.Fatalf("records = %d, want 1", len(h.records))
	}

	own := &captureHandler{enabled: true}
	New[orderPlaced](Logger(slog.New(own))).Trigger(context.Background(), orderPlaced{})
	if len(own.records) != 1 || len(h.records) != 1 {
		t.Errorf("own records = %d, default records = %d; want 1, 1", len(own.records), len(h.records))
	}
}
