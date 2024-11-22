package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"smart-hub/config"
	"sync"
)

type zapLogger struct {
	log *zap.Logger
}

var (
	globalLogger Logger
	once         sync.Once
)

func InitLogger(cfg *config.LogConfig) error {
	var err error
	once.Do(func() {
		zapConfig := zap.NewProductionConfig()

		level, err := zapcore.ParseLevel(cfg.Level)
		if err != nil {
			level = zapcore.DebugLevel
		}
		zapConfig.Level.SetLevel(level)
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		logger, err := zapConfig.Build()
		if err != nil {
			panic(err)
		}

		globalLogger = &zapLogger{
			log: logger,
		}
	})
	return err
}

func GetLogger() Logger {
	if globalLogger == nil {
		err := InitLogger(&config.LogConfig{})
		if err != nil {
			return nil
		}
	}
	return globalLogger
}

func toZapField(field interface{}) zap.Field {
	switch f := field.(type) {
	case error:
		return zap.Error(f)
	case zap.Field:
		return f
	default:
		return zap.Any("field", f)
	}
}

func convertFields(fields ...interface{}) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = toZapField(field)
	}
	return zapFields
}

func (l *zapLogger) Debug(message string, fields ...interface{}) {
	l.log.Debug(message, convertFields(fields...)...)
}

func (l *zapLogger) Info(message string, fields ...interface{}) {
	l.log.Info(message, convertFields(fields...)...)
}

func (l *zapLogger) Warn(message string, fields ...interface{}) {
	l.log.Warn(message, convertFields(fields...)...)
}

func (l *zapLogger) Error(message string, fields ...interface{}) {
	l.log.Error(message, convertFields(fields...)...)
}

func Debug(message string, fields ...interface{}) {
	GetLogger().Debug(message, fields...)
}

func Info(message string, fields ...interface{}) {
	GetLogger().Info(message, fields...)
}

func Warn(message string, fields ...interface{}) {
	GetLogger().Warn(message, fields...)
}

func Error(message string, fields ...interface{}) {
	GetLogger().Error(message, fields...)
}
