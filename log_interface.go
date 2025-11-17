package logstore

import (
	"time"

	"github.com/dromara/carbon/v2"
)

// LogInterface defines the public API for a log entry
type LogInterface interface {
	GetID() string
	SetID(id string) LogInterface

	GetLevel() string
	SetLevel(level string) LogInterface

	GetMessage() string
	SetMessage(message string) LogInterface

	GetContext() string
	SetContext(context string) LogInterface

	GetTime() time.Time
	SetTime(t time.Time) LogInterface

	GetTimeCarbon() *carbon.Carbon
	SetTimeCarbon(t *carbon.Carbon) LogInterface
}
