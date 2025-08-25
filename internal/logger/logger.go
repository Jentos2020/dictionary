package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

type CustomLogger struct {
	w io.Writer
}

var (
	Red    = "\033[91m"
	Green  = "\033[92m"
	Blue   = "\033[94m"
	Yellow = "\033[93m"
	Reset  = "\033[0m"
)

var levelColor = map[slog.Level]string{
	slog.LevelDebug: Green,
	slog.LevelInfo:  Blue,
	slog.LevelWarn:  Yellow,
	slog.LevelError: Red,
}

func CustomFormat(w io.Writer) *CustomLogger {
	return &CustomLogger{w: w}
}

func (h *CustomLogger) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *CustomLogger) Handle(ctx context.Context, r slog.Record) error {
	ts := r.Time.Format("2006-01-02 15:04:05.0000")
	color := levelColor[r.Level]
	level := r.Level.String()
	msg := fmt.Sprintf("%s%s%s %s", color, level, Reset, r.Message)

	r.Attrs(func(attr slog.Attr) bool {
		msg += fmt.Sprintf(" %s=%v", attr.Key, attr.Value)
		return true
	})

	_, err := fmt.Fprintf(h.w, "%s %s\n", ts, msg)
	return err
}

func (h *CustomLogger) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *CustomLogger) WithGroup(name string) slog.Handler {
	return h
}
