package logger

type options struct {
	driver string
	path   string
	stdout string
	level  string
}

type Option func(*options)

func setDefault() options {
	return options{
		driver: "default",
		path:   "temp/logs",
		stdout: "default",
		level:  "warn",
	}
}

func WithPath(s string) Option {
	return func(o *options) {
		o.path = s
	}
}

func WithStdout(s string) Option {
	return func(o *options) {
		o.stdout = s
	}
}

func WithLevel(s string) Option {
	return func(o *options) {
		o.level = s
	}
}
