package log

import (
	"cosmos-server/pkg/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//go:generate mockgen -destination=./mock/log_mock.go -package=log cosmos-server/pkg/log Logger

type Logger interface {
	Infof(format string, args ...any)
	Infow(msg string, keysAndValues ...any)
	Errorf(format string, args ...any)
	Debugf(format string, args ...any)
	Warnf(format string, args ...any)
}

type logger struct {
	zapLogger *zap.SugaredLogger
}

func NewLogger(logConfig config.LogConfig) (Logger, error) {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = getLoggingLevel(logConfig.Level)

	zapConfig.OutputPaths = []string{"stdout"}
	zapConfig.ErrorOutputPaths = []string{"stdout"}

	zapConfig.DisableCaller = true

	zapConfig.EncoderConfig.TimeKey = "time"
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zapLogger, err := zapConfig.Build()

	if err != nil {
		return nil, err
	}
	sugarZapLogger := zapLogger.Sugar()

	return &logger{
		zapLogger: sugarZapLogger,
	}, nil
}

func getLoggingLevel(logLevel string) zap.AtomicLevel {
	switch logLevel {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}

func (l *logger) Infof(format string, args ...any) {
	l.zapLogger.Infof(format, args...)
}

func (l *logger) Errorf(format string, args ...any) {
	l.zapLogger.Errorf(format, args...)
}

func (l *logger) Debugf(format string, args ...any) {
	l.zapLogger.Debugf(format, args...)
}

func (l *logger) Warnf(format string, args ...any) {
	l.zapLogger.Warnf(format, args...)
}

func (l *logger) Infow(msg string, keysAndValues ...any) {
	l.zapLogger.Infow(msg, keysAndValues...)
}
