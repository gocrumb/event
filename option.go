package event

import "log/slog"

type Option interface {
	apply(*emitterOptions)
}

type OptionFunc func(*emitterOptions)

func (f OptionFunc) apply(o *emitterOptions) {
	f(o)
}

type emitterOptions struct {
	logger *slog.Logger
}

func Logger(l *slog.Logger) Option {
	return OptionFunc(func(o *emitterOptions) {
		o.logger = l
	})
}
