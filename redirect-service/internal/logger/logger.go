package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/hoggir/re-path/redirect-service/internal/config"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)

	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}

type appLogger struct {
	logger *slog.Logger
}

func NewLogger(cfg *config.Config) Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: getLogLevel(cfg.App.Env),
	}

	if cfg.App.Env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return &appLogger{
		logger: slog.New(handler),
	}
}

func getLogLevel(env string) slog.Level {
	switch env {
	case "production":
		return slog.LevelInfo
	case "development":
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}

func (l *appLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *appLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *appLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *appLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *appLogger) Fatal(msg string, args ...any) {
	l.logger.Error(msg, args...)
	os.Exit(1)
}

func (l *appLogger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

func (l *appLogger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

func (l *appLogger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

func (l *appLogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}
