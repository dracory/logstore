package logstore

import (
	"time"
)

const (
	// LevelTrace trace level
	LevelTrace = "trace"
	// LevelDebug debug level
	LevelDebug = "debug"
	// LevelError error level
	LevelError = "error"
	// LevelFatal fatal level
	LevelFatal = "fatal"
	// LevelInfo info level
	LevelInfo = "info"
	// LevelPanic panic level
	LevelPanic = "panic"
	// LevelWarning warning level
	LevelWarning = "warning"
)

// Log type
type Log struct {
	ID      string
	Level   string
	Message string
	Context string
	Time    *time.Time
}
