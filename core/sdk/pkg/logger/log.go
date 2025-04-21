package logger

import (
	"go-project/core/debug/writer"
	"go-project/core/logger"
	"go-project/core/sdk/pkg"
	"io"
	"os"
)

func Setup(opts ...Option) logger.Logger {
	op := setDefault()

	for _, o := range opts {
		o(&op)
	}

	if !pkg.PathExist(op.path) {
		err := pkg.PathCreate(op.path)
		if err != nil {
			logger.Fatal("create log path error: %s", err.Error())
		}
	}

	var (
		err    error
		output io.Writer
	)

	switch op.stdout {
	case "file":
		output, err = writer.NewFileWriter(
			writer.WithPath(op.path),
			writer.WithCap(op.cap<<10),
		)
		if err != nil {
			logger.Fatal("create log file error: %s", err.Error())
		}
	default:
		output = os.Stdout
	}

	var level logger.Level
	level, err = logger.GetLevel(op.level)

	if err != nil {
		logger.Fatalf("get log level error: %s", err.Error())
	}

	switch op.driver {
	case "zap":
		// logger.DefaultLogger, err =
	}

	return nil
}
