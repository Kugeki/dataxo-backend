package restapi

import (
	"dataxo-backend-game-ms/pkg/httplog"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"time"
)

// LogFormatter to implement middlewares.LogFormatter
type LogFormatter struct {
	log          *slog.Logger
	defaultLevel slog.Level
}

func NewLogFormatter(log *slog.Logger, defaultLevel slog.Level) *LogFormatter {
	return &LogFormatter{
		log:          log.With(slog.String("context", "rest logging middleware")),
		defaultLevel: defaultLevel,
	}
}

func (l *LogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	return NewLogEntry(l.log, httplog.New(l.log, l.defaultLevel), r)
}

// LogEntry to implement middlewares.LogEntry
type LogEntry struct {
	log    *slog.Logger
	reqLog *httplog.RequestLogger

	r *http.Request
}

func NewLogEntry(log *slog.Logger, reqLog *httplog.RequestLogger, r *http.Request) *LogEntry {
	return &LogEntry{log: log, reqLog: reqLog, r: r}
}

func (e *LogEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, extra interface{}) {
	ctx := e.r.Context()
	e.reqLog.LogEnd(ctx, middleware.GetReqID(ctx), status, elapsed)
}

func (e *LogEntry) Panic(v interface{}, stack []byte) {
	e.log.Error("something panic",
		slog.Any("panic", v),
		slog.String("stack", string(stack)),
	)
}
