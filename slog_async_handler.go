package logstore

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
)

// SlogAsyncBufferSizeDefault bounds the pending-record queue used by
// SlogAsyncHandler. Sized so a burst of log calls cannot exhaust memory;
// when full, new records are dropped rather than blocking the caller on a
// slow database write.
const SlogAsyncBufferSizeDefault = 4096

// SlogAsyncHandler wraps SlogHandler and performs the real Handle work on
// a single background goroutine. SlogHandler issues a synchronous DB
// INSERT per record; calling it inline adds a database round-trip to every
// logged operation, which is expensive on request hot paths (e.g. a
// per-request access log).
//
// Records are cloned before enqueueing (slog.Record shares its attrs
// backing array with the caller). A full queue drops the record and counts
// it — logging must never block or fail the caller.
//
// WithAttrs/WithGroup children share the same queue and drain worker —
// Drain on any handler drains them all, and records logged after Drain
// are dropped silently. Drain must be called only once logging has
// stopped (typically during shutdown, before closing the database).
type SlogAsyncHandler struct {
	inner   *SlogHandler
	state   *slogAsyncState
	dropped *atomic.Uint64
}

type slogAsyncState struct {
	queue  chan slog.Record
	wg     sync.WaitGroup
	once   sync.Once
	closed atomic.Bool
}

// NewSlogAsyncHandler starts the drain worker and returns an async handler
// writing to the given log store. Call Drain on shutdown, before closing
// the underlying database, so buffered records are not written into a
// closed connection.
func NewSlogAsyncHandler(logStore StoreInterface) *SlogAsyncHandler {
	return NewSlogAsyncHandlerWithBuffer(logStore, SlogAsyncBufferSizeDefault)
}

// NewSlogAsyncHandlerWithBuffer is NewSlogAsyncHandler with a custom queue
// capacity. Values <= 0 fall back to SlogAsyncBufferSizeDefault.
func NewSlogAsyncHandlerWithBuffer(logStore StoreInterface, capacity int) *SlogAsyncHandler {
	if capacity <= 0 {
		capacity = SlogAsyncBufferSizeDefault
	}
	inner := NewSlogHandler(logStore)
	state := &slogAsyncState{queue: make(chan slog.Record, capacity)}
	state.wg.Add(1)
	go func() {
		defer state.wg.Done()
		for record := range state.queue {
			_ = inner.Handle(context.Background(), record)
		}
	}()
	return &SlogAsyncHandler{inner: inner, state: state, dropped: &atomic.Uint64{}}
}

func (h *SlogAsyncHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *SlogAsyncHandler) Handle(_ context.Context, record slog.Record) error {
	// Closed queue: drop silently instead of panicking on send to a closed
	// channel — in-flight requests may still log while Drain runs.
	if h.state.closed.Load() {
		return nil
	}
	select {
	case h.state.queue <- record.Clone():
	default:
		h.drop()
	}
	return nil
}

func (h *SlogAsyncHandler) drop() {
	if n := h.dropped.Add(1); n == 1 || n%1000 == 0 {
		fmt.Fprintf(os.Stderr,
			"logstore async handler: dropped %d log records (buffer full, capacity %d)\n",
			n, cap(h.state.queue))
	}
}

func (h *SlogAsyncHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SlogAsyncHandler{inner: h.inner.WithAttrs(attrs).(*SlogHandler), state: h.state, dropped: h.dropped}
}

func (h *SlogAsyncHandler) WithGroup(name string) slog.Handler {
	return &SlogAsyncHandler{inner: h.inner.WithGroup(name).(*SlogHandler), state: h.state, dropped: h.dropped}
}

// Drain closes the queue and waits for the worker to flush all pending
// records to the wrapped handler. Records logged after Drain are dropped
// silently. Safe to call more than once.
func (h *SlogAsyncHandler) Drain() {
	h.state.once.Do(func() {
		h.state.closed.Store(true)
		close(h.state.queue)
		h.state.wg.Wait()
	})
}

// DroppedCount returns how many records were dropped because the queue
// was full.
func (h *SlogAsyncHandler) DroppedCount() uint64 {
	return h.dropped.Load()
}
