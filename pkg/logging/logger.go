package logging

import (
	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, keysAndValues ...any)
}

type DefaultLogger struct {
	logger *zap.SugaredLogger
}

func NewDefaultLogger() *DefaultLogger {
	logger, err := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()

	if err != nil {
		panic(err)
	}

	return &DefaultLogger{logger: logger.Sugar()}
}

func (z *DefaultLogger) Info(msg string, keysAndValues ...any) {
	z.logger.Infow(msg, keysAndValues...)
}
