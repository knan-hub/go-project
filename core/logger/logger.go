package logger

type Logger interface {
	Init(opts ...Option) error
	Log(l Level, v ...interface{})
	Logf(l Level, s string, v ...interface{})
	Options() Options
	Fields(fields map[string]interface{}) Logger
	String() string
}

var DefaultLogger Logger

func Init(opts ...Option) error {
	return DefaultLogger.Init(opts...)
}

func Log(l Level, v ...interface{}) {
	DefaultLogger.Log(l, v...)
}

func Logf(l Level, s string, v ...interface{}) {
	DefaultLogger.Logf(l, s, v...)
}

func Fields(fields map[string]interface{}) Logger {
	return DefaultLogger.Fields(fields)
}

func String() string {
	return DefaultLogger.String()
}
