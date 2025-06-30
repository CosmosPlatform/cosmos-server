package log

import (
	"cosmos-server/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
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

func (l *logger) Infof(format string, args ...interface{}) {
	l.zapLogger.Infof(format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.zapLogger.Errorf(format, args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.zapLogger.Debugf(format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.zapLogger.Warnf(format, args...)
}

func (l *logger) Infow(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Infow(msg, keysAndValues...)
}
