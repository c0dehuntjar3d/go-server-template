package logger

import (
	"app/config"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type ZereLogger struct {
	logger *zerolog.Logger
}

var _ Interface = (*ZereLogger)(nil)

func NewZere(cfg *config.Config) (*ZereLogger, error) {
	var l zerolog.Level

	switch strings.ToLower(cfg.Log.Level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	skipFrameCount := 3
	logger := zerolog.New(os.Stdout).With().Timestamp().CallerWithSkipFrameCount(
		zerolog.CallerSkipFrameCount + skipFrameCount).Logger()

	return &ZereLogger{
		logger: &logger,
	}, nil
}

func (l *ZereLogger) Debug(message string, args ...Field) {
	l.logger.Debug().Msgf(message, args)
}

func (l *ZereLogger) Info(message string, args ...Field) {
	l.logger.Info().Msgf(message, args)
}

func (l *ZereLogger) Warn(message string, args ...Field) {
	l.logger.Error().Msgf(message, args)
}

func (l *ZereLogger) Error(message string, args ...Field) {
	if l.logger.GetLevel() == zerolog.DebugLevel {
		l.Debug(message, args...)
	}

	l.logger.Error().Msgf(message, args)
}

func (l *ZereLogger) Fatal(message string, args ...Field) {
	l.logger.Fatal().Msgf(message, args)
}
