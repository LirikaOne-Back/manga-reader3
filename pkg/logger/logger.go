package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Level представляет уровень логирования
type Level string

const (
	// LevelDebug - самый детальный уровень логирования
	LevelDebug Level = "debug"
	// LevelInfo - стандартный уровень логирования
	LevelInfo Level = "info"
	// LevelWarn - уровень для предупреждений
	LevelWarn Level = "warn"
	// LevelError - уровень для ошибок
	LevelError Level = "error"
)

// Config содержит настройки логгера
type Config struct {
	Level      Level
	JSONFormat bool
	Output     io.Writer
	WithSource bool
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() Config {
	return Config{
		Level:      LevelInfo,
		JSONFormat: false,
		Output:     os.Stdout,
		WithSource: true,
	}
}

// New создает и настраивает новый логгер slog
func New(cfg Config) *slog.Logger {
	var level slog.Level
	switch cfg.Level {
	case LevelDebug:
		level = slog.LevelDebug
	case LevelInfo:
		level = slog.LevelInfo
	case LevelWarn:
		level = slog.LevelWarn
	case LevelError:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var handler slog.Handler

	handlerOpts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(time.Now().Format(time.RFC3339))
			}
			return a
		},
	}

	if cfg.JSONFormat {
		handler = slog.NewJSONHandler(cfg.Output, handlerOpts)
	} else {
		handler = slog.NewTextHandler(cfg.Output, handlerOpts)
	}

	if cfg.WithSource {
		handler = &sourceHandler{
			Handler: handler,
		}
	}

	return slog.New(handler)
}

// sourceHandler добавляет информацию о файле/строке в лог
type sourceHandler struct {
	slog.Handler
}

func (h *sourceHandler) Handle(ctx context.Context, r slog.Record) error {
	// Определяем место вызова лога
	_, file, line, ok := runtime.Caller(4) // Нужно настроить уровень runtime.Caller
	if ok {
		r.AddAttrs(slog.String("source", filepath.Base(file)+":"+itoa(line)))
	}
	return h.Handler.Handle(ctx, r)
}

// itoa - быстрая конвертация int в string без выделения памяти
func itoa(i int) string {
	// Для большинства строк номеров строк 4 символов будет достаточно
	var buf [4]byte
	pos := len(buf)
	for i >= 10 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	pos--
	buf[pos] = byte('0' + i)
	return string(buf[pos:])
}

// RequestLogger создает middleware для логирования HTTP запросов
func RequestLogger(logger *slog.Logger) func(string, ...any) {
	return func(msg string, args ...any) {
		logger.Info(msg, args...)
	}
}
