package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func Init(level string) error {
	var zapLevel zapcore.Level

	switch strings.ToUpper(level) {
	case "DEBUG":
		zapLevel = zapcore.DebugLevel
	case "INFO":
		zapLevel = zapcore.InfoLevel
	case "WARN":
		zapLevel = zapcore.WarnLevel
	case "ERROR":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapLevel)

	var err error
	log, err = cfg.Build()
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

func Get() *zap.Logger {
	return log
}

func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}

func Fatal(msg string, args ...zap.Field) {
	if log != nil {
		log.Fatal(msg, args...)
	}
}
