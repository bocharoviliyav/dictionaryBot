package log

import (
	"go.uber.org/zap"
)

var logger *zap.Logger
var sugar *zap.SugaredLogger

func init() {
	logger, _ = zap.NewProduction(zap.AddCallerSkip(1))
	defer logger.Sync()
	sugar = logger.Sugar()
}

func Info(message string, args ...interface{}) {
	sugar.Infof(message, args)
}

func Error(message string, args ...interface{}) {
	sugar.Errorf(message, args)
}

func Debug(message string, args ...interface{}) {
	sugar.Debugf(message, args)
}
