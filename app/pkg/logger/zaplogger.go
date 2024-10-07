package logger

import (
	"app/config"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger          *zap.Logger
	errorLogEnabled bool
}

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
)

var supportedLoggingLevels = map[string]zapcore.Level{
	DebugLevel: zapcore.DebugLevel,
	InfoLevel:  zapcore.InfoLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
}

var _ Interface = (*ZapLogger)(nil)

func NewZap(cfg *config.Log) (*ZapLogger, error) {
	level := supportedLoggingLevels[strings.ToLower(cfg.Level)]

	loggerCfg := zap.Config{
		Encoding:         cfg.Encoding,
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{cfg.OutputPath},
		ErrorOutputPaths: []string{cfg.OutputPath},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "Time",
			LevelKey:       "Level",
			NameKey:        "Name",
			CallerKey:      "Caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "Message",
			StacktraceKey:  "Stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	logger, err := loggerCfg.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{
		logger:          logger,
		errorLogEnabled: cfg.ErrorEnabled,
	}, nil
}

func (l *ZapLogger) Debug(message string, args ...Field) {
	if len(args) == 0 {
		l.logger.Debug(message)
	} else {
		fields := mapFields(args...)
		l.logger.Debug(message, fields...)
	}
}

func (l *ZapLogger) Info(message string, args ...Field) {
	if len(args) == 0 {
		l.logger.Info(message)
	} else {
		fields := mapFields(args...)
		l.logger.Info(message, fields...)
	}
}

func (l *ZapLogger) Warn(message string, args ...Field) {
	if len(args) == 0 {
		l.logger.Warn(message)
	} else {
		fields := mapFields(args...)
		l.logger.Warn(message, fields...)
	}
}

func (l *ZapLogger) Error(message string, args ...Field) {
	if l.logger.Level() == zapcore.DebugLevel {
		l.Debug(message, args...)
		return
	}

	if len(args) == 0 {
		l.logger.Error(message)
	} else {
		fields := mapFields(args...)
		l.logger.Error(message, fields...)
	}
}

func (l *ZapLogger) Fatal(message string, args ...Field) {
	if len(args) == 0 {
		l.logger.Fatal(message)
	} else {
		fields := mapFields(args...)
		l.logger.Fatal(message, fields...)
	}
}

func mapFields(args ...Field) []zapcore.Field {
	fields := make([]zapcore.Field, 0, len(args))

	for _, arg := range args {
		switch v := arg.Value.(type) {
		case string:
			fields = append(fields, zap.String(arg.Key, v))
		case int:
			fields = append(fields, zap.Int(arg.Key, v))
		case bool:
			fields = append(fields, zap.Bool(arg.Key, v))
		case error:
			fields = append(fields, zap.Error(v))
		default:
			fields = append(fields, zap.Any(fmt.Sprintf(arg.Key), v))
		}
	}

	return fields
}