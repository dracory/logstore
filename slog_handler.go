package logstore

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
)

var _ slog.Handler = (*SlogHandler)(nil) // verify it extends the slog interface

type SlogHandler struct {
	slogHandler slog.Handler
	buffer      *bytes.Buffer
	mutex       *sync.Mutex
	logStore    StoreInterface
}

func NewSlogHandler(logStore StoreInterface) *SlogHandler {
	return &SlogHandler{
		slogHandler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
		buffer:   &bytes.Buffer{},
		mutex:    &sync.Mutex{},
		logStore: logStore,
	}
}

func (handler *SlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return handler.slogHandler.Enabled(ctx, level)
}

func (handler *SlogHandler) Handle(ctx context.Context, record slog.Record) error {
	level := record.Level.String()
	message := record.Message
	attrs, err := handler.computeAttrs(ctx, record)

	if err != nil {
		return fmt.Errorf("error when calling computeAttrs: %w", err)
	}

	if level == slog.LevelDebug.String() {
		return handler.logStore.DebugWithContext(message, attrs)
	}

	if level == slog.LevelInfo.String() {
		return handler.logStore.InfoWithContext(message, attrs)
	}

	if level == slog.LevelWarn.String() {
		return handler.logStore.WarnWithContext(message, attrs)
	}

	if level == slog.LevelError.String() {
		return handler.logStore.ErrorWithContext(message, attrs)
	}

	return handler.logStore.FatalWithContext(message, attrs)
}

func (handler *SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SlogHandler{
		slogHandler: handler.slogHandler.WithAttrs(attrs),
		buffer:      handler.buffer,
		mutex:       handler.mutex,
		logStore:    handler.logStore,
	}
}

func (handler *SlogHandler) WithGroup(name string) slog.Handler {
	return &SlogHandler{
		slogHandler: handler.slogHandler.WithGroup(name),
		buffer:      handler.buffer,
		mutex:       handler.mutex,
		logStore:    handler.logStore,
	}
}

func (handler *SlogHandler) computeAttrs(
	ctx context.Context,
	r slog.Record,
) (map[string]any, error) {
	handler.mutex.Lock()

	defer func() {
		handler.buffer.Reset()
		handler.mutex.Unlock()
	}()

	if err := handler.slogHandler.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	attrs := map[string]any{}

	r.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value.Any()
		return true
	})

	return attrs, nil
}
