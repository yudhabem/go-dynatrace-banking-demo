package logger

import "go.uber.org/zap"

var Log *zap.Logger

func New() *zap.Logger {

	if Log != nil {
		return Log
	}

	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	Log = l
	return Log
}
