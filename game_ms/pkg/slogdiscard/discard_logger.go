package slogdiscard

import (
	"context"
	"log/slog"
)

var discardLogger *slog.Logger
var discardHandler slog.Handler

func Logger() *slog.Logger {
	if discardLogger == nil {
		discardLogger = New()
	}
	return discardLogger
}

func Handler() slog.Handler {
	if discardHandler == nil {
		discardHandler = NewHandler()
	}
	return discardHandler
}

func New() *slog.Logger {
	return slog.New(NewHandler())
}

func NewHandler() slog.Handler {
	return &DiscardHandler{}
}

func LoggerIfNil(log *slog.Logger) *slog.Logger {
	if log == nil {
		return Logger()
	}

	return log
}

type DiscardHandler struct{}

func (h DiscardHandler) Enabled(context.Context, slog.Level) bool {
	return false
}

func (h DiscardHandler) Handle(context.Context, slog.Record) error {
	return nil
}

func (h DiscardHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h DiscardHandler) WithGroup(name string) slog.Handler {
	return h
}
