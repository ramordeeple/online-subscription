package logger

import (
	"go.uber.org/zap"
)

var log *zap.Logger

func Init() error {
	var err error
	log, err = zap.NewDevelopment()
	return err
}

func Info(msg string, args ...zap.Field) {
	if log != nil {
		log.Info(msg, args...)
	}
}

func Error(msg string, args ...zap.Field) {
	if log != nil {
		log.Error(msg, args...)
	}
}

func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}
