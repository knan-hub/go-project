package logger

import (
	"context"
	"io"
	"os"
)

type Options struct {
	Level           Level
	Fields          map[string]interface{}
	Out             io.Writer
	CallerSkipCount int
	Context         context.Context
	Name            string
}

type Option func(*Options)

func WithLevel(l Level) Option {
	return func(o *Options) {
		o.Level = l
	}
}

func WithFields(f map[string]interface{}) Option {
	return func(o *Options) {
		o.Fields = f
	}
}

func WithOut(w io.Writer) Option {
	return func(o *Options) {
		o.Out = w
	}
}

func WithCallerSkipCount(i int) Option {
	return func(o *Options) {
		o.CallerSkipCount = i
	}
}

func WithName(s string) Option {
	return func(o *Options) {
		o.Name = s
	}
}

func SetOption(k, v interface{}) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}

func DefaultOptions() Options {
	return Options{
		Level:           InfoLevel,
		Fields:          make(map[string]interface{}),
		Out:             os.Stderr,
		CallerSkipCount: 3,
		Context:         context.Background(),
		Name:            "",
	}
}
