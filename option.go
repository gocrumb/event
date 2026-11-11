package event

import "log/slog"

// Option configures an Emitter.
type Option interface {
	apply(*emitterOptions)
}

// OptionFunc adapts a function to an Option.
type OptionFunc func(*emitterOptions)

func (f OptionFunc) apply(o *emitterOptions) {
	f(o)
}

type emitterOptions struct {
	logger *slog.Logger
}

// Logger makes the emitter log each triggered event to l at debug level.
func Logger(l *slog.Logger) Option {
	return OptionFunc(func(o *emitterOptions) {
		o.logger = l
	})
}
