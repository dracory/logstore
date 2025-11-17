package logstore

const COLUMN_CONTEXT = "context"
const COLUMN_ID = "id"
const COLUMN_LEVEL = "level"
const COLUMN_MESSAGE = "message"
const COLUMN_TIME = "time"

// Log levels
const (
	// LevelTrace trace level
	LEVEL_TRACE = "trace"
	// LevelDebug debug level
	LEVEL_DEBUG = "debug"
	// LevelError error level
	LEVEL_ERROR = "error"
	// LevelFatal fatal level
	LEVEL_FATAL = "fatal"
	// LevelInfo info level
	LEVEL_INFO = "info"
	// LevelPanic panic level
	LEVEL_PANIC = "panic"
	// LevelWarning warning level
	LEVEL_WARNING = "warning"
)
