package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"app/pkg/config"
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
		l.logger.Debug(withColor(Cyan, message))
	} else {
		fields := mapFields(args...)
		l.logger.Debug(withColor(Cyan, message), fields...)
	}
}

func (l *ZapLogger) Info(message string, args ...Field) {
	if len(args) == 0 {
		l.logger.Info(withColor(Green, message))
	} else {
		fields := mapFields(args...)
		l.logger.Info(withColor(Green, message), fields...)
	}
}

func (l *ZapLogger) Warn(message string, args ...Field) {
	if len(args) == 0 {
		l.logger.Warn(withColor(Yellow, message))
	} else {
		fields := mapFields(args...)
		l.logger.Warn(withColor(Yellow, message), fields...)
	}
}

func (l *ZapLogger) Error(message string, args ...Field) {
	if len(args) == 0 {
		l.logger.Error(withColor(Red, message))
	} else {
		fields := mapFields(args...)
		l.logger.Error(withColor(Red, message), fields...)
	}
}

func (l *ZapLogger) Fatal(message string, args ...Field) {
	if len(args) == 0 {
		l.logger.Fatal(withColor(Red, message))
	} else {
		fields := mapFields(args...)
		l.logger.Fatal(withColor(Red, message), fields...)
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
