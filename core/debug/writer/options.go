package writer

type Options struct {
	path   string
	suffix string
	cap    uint
}

func setup() Options {
	return Options{
		path:   "/tmp/go-project",
		suffix: "log",
	}
}

type Option func(*Options)

func WithPath(s string) Option {
	return func(o *Options) {
		o.path = s
	}
}

func WithSuffix(s string) Option {
	return func(o *Options) {
		o.suffix = s
	}
}

func WithCap(u uint) Option {
	return func(o *Options) {
		o.cap = u
	}
}
