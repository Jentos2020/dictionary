package pg

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm/logger"
)

type SlogAdapter struct {
	log   *slog.Logger
	level logger.LogLevel
}

func NewSlogAdapter(log *slog.Logger) logger.Interface {
	return &SlogAdapter{log: log, level: logger.Error} // по умолчанию только ошибки
}

func (l *SlogAdapter) LogMode(level logger.LogLevel) logger.Interface {
	return &SlogAdapter{log: l.log, level: level}
}

func (l *SlogAdapter) Info(ctx context.Context, msg string, data ...interface{}) {
	l.log.InfoContext(ctx, fmt.Sprintf(msg, data...))
}

func (l *SlogAdapter) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.log.WarnContext(ctx, fmt.Sprintf(msg, data...))
}

func (l *SlogAdapter) Error(ctx context.Context, msg string, data ...interface{}) {
	l.log.ErrorContext(ctx, fmt.Sprintf(msg, data...))
}

func (l *SlogAdapter) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level == logger.Silent {
		return
	}

	sql, rows := fc()
	duration := time.Since(begin)

	if err != nil && l.level >= logger.Error {
		l.log.ErrorContext(ctx, "SQL error",
			slog.String("query", sql),
			slog.Int64("rows", rows),
			slog.Duration("elapsed", duration),
			slog.String("error", err.Error()),
		)
		return
	}

	if l.level >= logger.Info {
		l.log.DebugContext(ctx, "SQL executed",
			slog.String("query", sql),
			slog.Int64("rows", rows),
			slog.Duration("elapsed", duration),
		)
	}
}
