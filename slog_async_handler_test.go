package logstore

import (
	"log/slog"
	"sync"
	"testing"
	"time"
)

// countingStore captures StoreInterface writes for assertions.
type countingStore struct {
	StoreInterface
	mu    sync.Mutex
	count int
}

func (s *countingStore) InfoWithContext(message string, context interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.count++
	return nil
}

func (s *countingStore) total() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.count
}

// TestSlogAsyncHandler verifies records are delivered to the store
// asynchronously and Drain flushes everything before shutdown.
func TestSlogAsyncHandler(t *testing.T) {
	store := &countingStore{}
	handler := NewSlogAsyncHandler(store)
	logger := slog.New(handler)

	logger.Info("first")
	logger.Info("second")

	handler.Drain()

	if got := store.total(); got != 2 {
		t.Fatalf("expected 2 records after drain, got %d", got)
	}
	if handler.DroppedCount() != 0 {
		t.Fatalf("expected 0 dropped records, got %d", handler.DroppedCount())
	}
}

// TestSlogAsyncHandlerWithAttrs verifies a child logger built via With
// still reaches the store — guards against a dropped logStore reference
// in the wrapped handler chain.
func TestSlogAsyncHandlerWithAttrs(t *testing.T) {
	store := &countingStore{}
	handler := NewSlogAsyncHandler(store)
	logger := slog.New(handler).With("key", "value").WithGroup("group")

	logger.Info("child record")

	handler.Drain()

	if got := store.total(); got != 1 {
		t.Fatalf("expected 1 record after drain, got %d", got)
	}
}

// TestSlogAsyncHandlerDoesNotBlock verifies Handle returns promptly even
// when the store would be slow — the point of the async wrapper.
func TestSlogAsyncHandlerDoesNotBlock(t *testing.T) {
	store := &countingStore{}
	handler := NewSlogAsyncHandler(store)
	logger := slog.New(handler)

	start := time.Now()
	for i := 0; i < 10; i++ {
		logger.Info("burst")
	}
	if elapsed := time.Since(start); elapsed > 50*time.Millisecond {
		t.Fatalf("Handle blocked for %v, expected near-instant enqueue", elapsed)
	}

	handler.Drain()
}

// TestSlogAsyncHandlerDropOnFull verifies a full queue drops records
// instead of blocking, and counts the drops.
func TestSlogAsyncHandlerDropOnFull(t *testing.T) {
	store := &blockingStore{release: make(chan struct{})}
	handler := NewSlogAsyncHandlerWithBuffer(store, 2)
	logger := slog.New(handler)

	// Flood far beyond capacity; whether the worker has claimed a record
	// yet is irrelevant — the queue can hold at most cap+1 records.
	for i := 0; i < 50; i++ {
		logger.Info("flood")
	}

	deadline := time.Now().Add(2 * time.Second)
	for handler.DroppedCount() == 0 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	close(store.release)
	handler.Drain()

	if handler.DroppedCount() == 0 {
		t.Fatal("expected some records to be dropped on a full buffer")
	}
}

// TestSlogAsyncHandlerHandleAfterDrain verifies logging after Drain drops
// silently instead of panicking on a closed channel.
func TestSlogAsyncHandlerHandleAfterDrain(t *testing.T) {
	store := &countingStore{}
	handler := NewSlogAsyncHandler(store)
	logger := slog.New(handler)

	logger.Info("before drain")
	handler.Drain()

	// Must not panic — in-flight requests may still log during shutdown.
	logger.Info("after drain")

	if got := store.total(); got != 1 {
		t.Fatalf("expected only the pre-drain record, got %d", got)
	}
}

type blockingStore struct {
	StoreInterface
	release chan struct{}
}

func (s *blockingStore) InfoWithContext(message string, context interface{}) error {
	<-s.release
	return nil
}
