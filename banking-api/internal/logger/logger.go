package logger

import (
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func New(loggerProviders ...log.LoggerProvider) *zap.Logger {

	if Log != nil {
		return Log
	}

	var loggerProvider log.LoggerProvider
	if len(loggerProviders) > 0 {
		loggerProvider = loggerProviders[0]
	}

	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	if loggerProvider != nil {
		l = l.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(core, otelzap.NewCore("banking-api", otelzap.WithLoggerProvider(loggerProvider)))
		}))
	}

	Log = l

	return Log
}
