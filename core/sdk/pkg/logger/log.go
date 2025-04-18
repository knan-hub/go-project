package logger

import (
	"go-project/core/logger"
	"go-project/core/sdk/pkg"
	"io"
	"log"
)

func Setup(opts ...Option) logger.Logger {
	op := setDefault()

	for _, o := range opts {
		o(&op)
	}

	if !pkg.PathExist(op.path) {
		err := pkg.PathCreate(op.path)
		if err != nil {
			log.Fatalf("create log path error: %s", err.Error())
		}
	}

	var (
		err    error
		output io.Writer
	)

	switch op.stdout {
	case "file":

	}

	return nil
}
