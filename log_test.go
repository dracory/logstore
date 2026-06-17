package logstore

import (
	"testing"
	"time"

	"github.com/dromara/carbon/v2"
)

func Test_LogImplementation_GettersSetters(t *testing.T) {
	log := NewLog().
		SetID("id-123").
		SetLevel(LEVEL_INFO).
		SetMessage("message").
		SetContext("ctx")

	if log.GetID() != "id-123" {
		t.Fatalf("expected ID 'id-123', got '%s'", log.GetID())
	}

	if log.GetLevel() != LEVEL_INFO {
		t.Fatalf("expected level '%s', got '%s'", LEVEL_INFO, log.GetLevel())
	}

	if log.GetMessage() != "message" {
		t.Fatalf("expected message 'message', got '%s'", log.GetMessage())
	}

	if log.GetContext() != "ctx" {
		t.Fatalf("expected context 'ctx', got '%s'", log.GetContext())
	}
}

func Test_LogImplementation_DefaultConstructorValues(t *testing.T) {
	log := NewLog()

	if log.GetID() == "" {
		t.Fatal("expected NewLog to set a non-empty ID")
	}

	if log.GetLevel() != LEVEL_INFO {
		t.Fatalf("expected default level '%s', got '%s'", LEVEL_INFO, log.GetLevel())
	}

	if log.GetMessage() != "" {
		t.Fatalf("expected default message to be empty, got '%s'", log.GetMessage())
	}

	if log.GetContext() != "" {
		t.Fatalf("expected default context to be empty, got '%s'", log.GetContext())
	}
}

func Test_LogImplementation_TimeAndCarbon(t *testing.T) {
	// NewLog() should initialize time to a non-zero UTC value
	log := NewLog()

	initial := log.GetTime()
	if initial.IsZero() {
		t.Fatal("expected NewLog time to be non-zero")
	}

	// Ensure GetTimeCarbon reflects the same instant in UTC
	c := log.GetTimeCarbon()
	if c == nil {
		t.Fatal("expected GetTimeCarbon to return non-nil carbon instance")
	}

	if !c.StdTime().Equal(initial) {
		t.Fatal("carbon time does not match underlying time")
	}

	// SetTimeCarbon should update the underlying time
	future := carbon.Now(carbon.UTC).AddMinutes(5)
	log.SetTimeCarbon(future)

	if !log.GetTime().Equal(future.StdTime()) {
		t.Fatal("SetTimeCarbon did not update underlying time correctly")
	}

	// SetTime should also be reflected by GetTimeCarbon
	manual := time.Now().UTC().Add(10 * time.Minute)
	log.SetTime(manual)

	c2 := log.GetTimeCarbon()
	if !c2.StdTime().Equal(manual) {
		t.Fatal("SetTime did not update carbon view correctly")
	}
}
