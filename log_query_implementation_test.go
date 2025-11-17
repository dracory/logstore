package logstore

import (
	"strings"
	"testing"
)

func Test_LogQueryImplementation_Validate_MessageAndContextTerms(t *testing.T) {
	// message_contains empty should error
	q := LogQuery().SetMessageContains("")
	if err := q.(*logQueryImplementation).Validate(); err == nil {
		t.Fatal("expected error for empty message_contains, got nil")
	}

	// message_not_contains empty should error
	q = LogQuery().SetMessageNotContains("")
	if err := q.(*logQueryImplementation).Validate(); err == nil {
		t.Fatal("expected error for empty message_not_contains, got nil")
	}

	// context_contains empty should error
	q = LogQuery().SetContextContains("")
	if err := q.(*logQueryImplementation).Validate(); err == nil {
		t.Fatal("expected error for empty context_contains, got nil")
	}

	// context_not_contains empty should error
	q = LogQuery().SetContextNotContains("")
	if err := q.(*logQueryImplementation).Validate(); err == nil {
		t.Fatal("expected error for empty context_not_contains, got nil")
	}

	// Non-empty terms should pass validation
	q = LogQuery().
		SetMessageContains("foo").
		SetMessageNotContains("bar").
		SetContextContains("baz").
		SetContextNotContains("qux")

	if err := q.(*logQueryImplementation).Validate(); err != nil {
		t.Fatalf("expected no error for valid terms, got: %v", err)
	}
}

// helper store stub implementing just enough of StoreInterface for ToSelectDataset

type logQueryTestStore struct{}

func (s *logQueryTestStore) AutoMigrate() error                    { return nil }
func (s *logQueryTestStore) EnableDebug(debug bool)                {}
func (s *logQueryTestStore) GetDriverName() string                 { return "sqlite3" }
func (s *logQueryTestStore) GetLogTableName() string               { return "logs" }
func (s *logQueryTestStore) Log(logEntry LogInterface) error       { return nil }
func (s *logQueryTestStore) LogCreate(logEntry LogInterface) error { return nil }
func (s *logQueryTestStore) Debug(message string) error            { return nil }
func (s *logQueryTestStore) DebugWithContext(message string, context interface{}) error {
	return nil
}
func (s *logQueryTestStore) Error(message string) error { return nil }
func (s *logQueryTestStore) ErrorWithContext(message string, context interface{}) error {
	return nil
}
func (s *logQueryTestStore) Fatal(message string) error { return nil }
func (s *logQueryTestStore) FatalWithContext(message string, context interface{}) error {
	return nil
}
func (s *logQueryTestStore) Info(message string) error { return nil }
func (s *logQueryTestStore) InfoWithContext(message string, context interface{}) error {
	return nil
}
func (s *logQueryTestStore) Panic(message string)                                 {}
func (s *logQueryTestStore) PanicWithContext(message string, context interface{}) {}
func (s *logQueryTestStore) Trace(message string) error                           { return nil }
func (s *logQueryTestStore) TraceWithContext(message string, context interface{}) error {
	return nil
}
func (s *logQueryTestStore) Warn(message string) error { return nil }
func (s *logQueryTestStore) WarnWithContext(message string, context interface{}) error {
	return nil
}
func (s *logQueryTestStore) LogList(query LogQueryInterface) ([]LogInterface, error) {
	return nil, nil
}
func (s *logQueryTestStore) LogDelete(logEntry LogInterface) error       { return nil }
func (s *logQueryTestStore) LogDeleteByID(id string) error               { return nil }
func (s *logQueryTestStore) LogFindByID(id string) (LogInterface, error) { return nil, nil }

func (s *logQueryTestStore) LogCount(query LogQueryInterface) (int, error) {
	return 0, nil
}

func Test_LogQueryImplementation_ToSelectDataset_MessageAndContextFilters(t *testing.T) {
	st := &logQueryTestStore{}

	q := LogQuery().
		SetMessageContains("error").
		SetMessageNotContains("debug").
		SetContextContains("user").
		SetContextNotContains("trace")

	selectDataset, _, err := q.ToSelectDataset(st)
	if err != nil {
		t.Fatalf("unexpected error from ToSelectDataset: %v", err)
	}

	sql, args, err := selectDataset.Prepared(true).ToSQL()
	if err != nil {
		t.Fatalf("unexpected error generating SQL: %v", err)
	}

	// Basic checks that our LIKE/NOT LIKE clauses are present in the SQL
	expectedSqlTerms := []string{
		"message", "LIKE",
		"message", "NOT LIKE",
		"context", "LIKE",
		"context", "NOT LIKE",
	}

	for _, term := range expectedSqlTerms {
		if !strings.Contains(sql, term) {
			t.Fatalf("expected SQL to contain %q, got: %s", term, sql)
		}
	}

	// And the concrete LIKE values should be present in the prepared args
	if len(args) != 4 {
		t.Fatalf("expected 4 args, got %d (%v)", len(args), args)
	}

	expectedArgs := []string{"%error%", "%debug%", "%user%", "%trace%"}
	for i, expected := range expectedArgs {
		if args[i] != expected {
			t.Fatalf("expected arg[%d] = %q, got %v", i, expected, args[i])
		}
	}
}
