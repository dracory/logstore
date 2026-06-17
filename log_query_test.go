package logstore

import (
	"context"
	"database/sql"
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

// helper store stub implementing just enough of StoreInterface for testing
type logQueryTestStore struct{}

func (s *logQueryTestStore) SetLogTableName(logTableName string)                  {}
func (s *logQueryTestStore) MigrateDown(ctx context.Context, tx ...*sql.Tx) error { return nil }
func (s *logQueryTestStore) MigrateUp(ctx context.Context, tx ...*sql.Tx) error   { return nil }
func (s *logQueryTestStore) EnableDebug(debug bool)                               {}
func (s *logQueryTestStore) GetDB() *sql.DB                                       { return nil }
func (s *logQueryTestStore) GetLogTableName() string                              { return "logs" }
func (s *logQueryTestStore) Log(logEntry LogInterface) error                      { return nil }
func (s *logQueryTestStore) LogCreate(ctx context.Context, logEntry LogInterface) error {
	return nil
}
func (s *logQueryTestStore) Debug(message string) error { return nil }
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
func (s *logQueryTestStore) LogList(ctx context.Context, query LogQueryInterface) ([]LogInterface, error) {
	return nil, nil
}
func (s *logQueryTestStore) LogDelete(ctx context.Context, logEntry LogInterface) error { return nil }
func (s *logQueryTestStore) LogDeleteByID(ctx context.Context, id string) error         { return nil }
func (s *logQueryTestStore) LogDeleteByIDs(ctx context.Context, ids []string) error     { return nil }
func (s *logQueryTestStore) LogFindByID(ctx context.Context, id string) (LogInterface, error) {
	return nil, nil
}

func (s *logQueryTestStore) LogCount(ctx context.Context, query LogQueryInterface) (int64, error) {
	return 0, nil
}
