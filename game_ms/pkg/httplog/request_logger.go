package httplog

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type RequestLogger struct {
	l            *slog.Logger
	DefaultLevel slog.Level
}

func New(l *slog.Logger, defaultLevel slog.Level) *RequestLogger {
	return &RequestLogger{l: l, DefaultLevel: defaultLevel}
}

func (l *RequestLogger) LogBegin(ctx context.Context, r *http.Request, requestID string) {
	l.l.Log(ctx, l.DefaultLevel, "request incoming",
		slog.String("req_id", requestID),
		slog.String("method", r.Method),
		slog.String("url", r.URL.String()),
		slog.String("remote_addr", r.RemoteAddr),
		slog.String("user_agent", r.UserAgent()),
		slog.String("proto", r.Proto),
	)
}

func (l *RequestLogger) LogEnd(ctx context.Context, requestID string, status int, elapsed time.Duration) {
	level := l.DefaultLevel
	switch {
	case status >= 500:
		level = slog.LevelError
	case status >= 400:
		level = slog.LevelWarn
	}

	l.l.Log(ctx, level, "request completed",
		slog.String("req_id", requestID),
		slog.Int("status", status),
		slog.Duration("elapsed", elapsed),
	)
}
