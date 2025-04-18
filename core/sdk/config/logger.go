package config

import "go-project/core/sdk/pkg/logger"

type Logger struct {
	Path      string
	Stdout    string
	Level     string
	EnabledDB bool
}

func (l Logger) Setup() {
	logger.Setup(
		logger.WithPath(l.Path),
		logger.WithStdout(l.Stdout),
		logger.WithLevel(l.Level),
	)
}

var LoggerConfig = new(Logger)
