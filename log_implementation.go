package logstore

import (
	"time"

	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

// logImplementation is the concrete implementation of LogInterface
type logImplementation struct {
	id      string
	level   string
	message string
	context string
	time    time.Time
}

var _ LogInterface = (*logImplementation)(nil)

// NewLog creates a new log with the current UTC time
func NewLog() LogInterface {
	return &logImplementation{
		id:      uid.MicroUid(),
		level:   LevelInfo,
		message: "",
		context: "",
		time:    time.Now().UTC(),
	}
}

// NewLogWithData creates a new log from existing values
func NewLogWithData(id, level, message, context string, t time.Time) LogInterface {
	return &logImplementation{
		id:      id,
		level:   level,
		message: message,
		context: context,
		time:    t,
	}
}

func (l *logImplementation) GetID() string {
	return l.id
}

func (l *logImplementation) SetID(id string) LogInterface {
	l.id = id
	return l
}

func (l *logImplementation) GetLevel() string {
	return l.level
}

func (l *logImplementation) SetLevel(level string) LogInterface {
	l.level = level
	return l
}

func (l *logImplementation) GetMessage() string {
	return l.message
}

func (l *logImplementation) SetMessage(message string) LogInterface {
	l.message = message
	return l
}

func (l *logImplementation) GetContext() string {
	return l.context
}

func (l *logImplementation) SetContext(context string) LogInterface {
	l.context = context
	return l
}

func (l *logImplementation) GetTime() time.Time {
	return l.time
}

func (l *logImplementation) SetTime(t time.Time) LogInterface {
	l.time = t
	return l
}

func (l *logImplementation) GetTimeCarbon() *carbon.Carbon {
	return carbon.CreateFromStdTime(l.time, carbon.UTC)
}

func (l *logImplementation) SetTimeCarbon(t *carbon.Carbon) LogInterface {
	l.time = t.StdTime()
	return l
}
