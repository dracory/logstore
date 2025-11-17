package logstore

import (
	"errors"
)

// Errors
var (
	ErrLogTableNameRequired = errors.New("log store: logTableName is required")
	ErrDBRequired           = errors.New("log store: DB is required")
)

// StoreInterface defines the interface for a log store
type StoreInterface interface {
	// AutoMigrate creates the necessary database tables
	AutoMigrate() error

	// EnableDebug enables or disables debug mode
	EnableDebug(debug bool)

	// Log adds a log entry
	Log(logEntry *Log) error

	// Debug adds a debug log
	Debug(message string) error

	// DebugWithContext adds a debug log with context data
	DebugWithContext(message string, context interface{}) error

	// Error adds an error log
	Error(message string) error

	// ErrorWithContext adds an error log with context data
	ErrorWithContext(message string, context interface{}) error

	// Fatal adds a fatal log
	Fatal(message string) error

	// FatalWithContext adds a fatal log with context data
	FatalWithContext(message string, context interface{}) error

	// Info adds an info log
	Info(message string) error

	// InfoWithContext adds an info log with context data
	InfoWithContext(message string, context interface{}) error

	// Panic adds a panic log and calls panic(message) after logging
	Panic(message string)

	// PanicWithContext adds a panic log with context data and calls panic(message) after logging
	PanicWithContext(message string, context interface{})

	// Trace adds a trace log
	Trace(message string) error

	// TraceWithContext adds a trace log with context data
	TraceWithContext(message string, context interface{}) error

	// Warn adds a warn log
	Warn(message string) error

	// WarnWithContext adds a warn log with context data
	WarnWithContext(message string, context interface{}) error
}
