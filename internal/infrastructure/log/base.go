package log

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/noxhalley/funken/config"
)

type loggerCtxKey struct{}

const messageKey = "message"

var (
	keys         []string
	logMapCtxKey = loggerCtxKey{}
)

// Initialize: initializes the logger with default handler
func Initialize(w io.Writer, cfg *config.Config, keyInput []string) {
	keys = append(keys, keyInput...)
	level := slog.LevelInfo
	if cfg.App.LogLevel == "DEBUG" {
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if v, ok := a.Value.Any().(time.Duration); ok {
				a.Value = slog.StringValue(v.String())
			}
			if a.Key != slog.MessageKey {
				return a
			}
			a.Key = messageKey
			return a
		},
	}

	slog.SetDefault(slog.New(&handler{
		Handler: slog.NewJSONHandler(w, opts),
	}))
}

// AddLogValToCtx adds a key-val pair to the context in sync.Map for thread safely
// this value automatically added to the log record with defaultHandler
func AddLogValToCtx(ctx context.Context, key string, val interface{}) context.Context {
	m, ok := ctx.Value(logMapCtxKey).(*sync.Map)
	if !ok {
		m = &sync.Map{}
	}
	m.Store(key, val)
	return context.WithValue(ctx, logMapCtxKey, m)
}

func Group(key string, args ...any) slog.Attr {
	return slog.Group(key, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, msg, args...)
}

func Debug(ctx context.Context, msg string, args ...any) {
	slog.DebugContext(ctx, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
}

func Fatal(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
	os.Exit(1)
}
