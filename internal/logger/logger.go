package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel string = "DEBUG"
	InfoLevel  string = "INFO"
	WarnLevel  string = "WARN"
	ErrorLevel string = "ERROR"
)

func New(level string) (*zap.Logger, error) {
	var zapLevel zapcore.Level

	switch level {
	case DebugLevel:
		zapLevel = zapcore.DebugLevel
	case WarnLevel:
		zapLevel = zapcore.WarnLevel
	case ErrorLevel:
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapLevel)
	return cfg.Build()
}
