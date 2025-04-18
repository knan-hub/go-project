package logger

import (
	"fmt"
	"os"
)

type Level int8

const (
	TraceLevel Level = iota - 2
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return "unknown"
	}
}

func (l Level) LevelForGorm() int {
	switch l {
	case FatalLevel, ErrorLevel:
		return 2
	case WarnLevel:
		return 3
	case InfoLevel, DebugLevel, TraceLevel:
		return 4
	default:
		return 1
	}
}

func (l Level) Enabled(lv Level) bool {
	return lv >= l
}

func GetLevel(s string) (Level, error) {
	switch s {
	case TraceLevel.String():
		return TraceLevel, nil
	case DebugLevel.String():
		return DebugLevel, nil
	case InfoLevel.String():
		return InfoLevel, nil
	case WarnLevel.String():
		return WarnLevel, nil
	case ErrorLevel.String():
		return ErrorLevel, nil
	case FatalLevel.String():
		return FatalLevel, nil
	default:
		return InfoLevel, fmt.Errorf("Unknown Level String: '%s', default to InfoLevel", s)
	}
}

func Info(args ...interface{}) {
	DefaultLogger.Log(InfoLevel, args)
}

func Infof(s string, args ...interface{}) {
	DefaultLogger.Logf(InfoLevel, s, args)
}

func Trace(args ...interface{}) {
	DefaultLogger.Log(TraceLevel, args)
}

func Tracef(s string, args ...interface{}) {
	DefaultLogger.Logf(TraceLevel, s, args)
}

func Debug(args ...interface{}) {
	DefaultLogger.Log(DebugLevel, args)
}

func Debugf(s string, args ...interface{}) {
	DefaultLogger.Logf(DebugLevel, s, args)
}

func Warn(args ...interface{}) {
	DefaultLogger.Log(WarnLevel, args)
}

func Warnf(s string, args ...interface{}) {
	DefaultLogger.Logf(WarnLevel, s, args)
}

func Error(args ...interface{}) {
	DefaultLogger.Log(ErrorLevel, args)
}

func Errorf(s string, args ...interface{}) {
	DefaultLogger.Logf(ErrorLevel, s, args)
}

func Fatal(args ...interface{}) {
	DefaultLogger.Log(FatalLevel, args)
	os.Exit(1)
}

func Fatalf(s string, args ...interface{}) {
	DefaultLogger.Logf(FatalLevel, s, args)
	os.Exit(1)
}

func V(lv Level, l Logger) bool {
	log := DefaultLogger
	if l != nil {
		log = l
	}
	return log.Options().Level <= lv
}
