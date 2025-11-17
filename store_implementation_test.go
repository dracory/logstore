package logstore

import (
	"database/sql"
	"testing"
	"time"

	"github.com/dracory/uid"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	dsn := ":memory:?parseTime=true"
	db, err := sql.Open("sqlite3", dsn)

	if err != nil {
		panic(err)
	}

	return db
}

func TestStoreCreate(t *testing.T) {
	db := InitDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log_create",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	if store == nil {
		t.Fatal("Store could not be created")
	}
}

func TestNewStore_Error_LogTableNameRequired(t *testing.T) {
	db := InitDB()

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "",
		AutomigrateEnabled: true,
	})

	if err == nil {
		t.Fatalf("expected error for empty LogTableName, got nil (store=%v)", store)
	}
}

func TestNewStore_Error_DBRequired(t *testing.T) {
	store, err := NewStore(NewStoreOptions{
		DB:                 nil,
		LogTableName:       "logs",
		AutomigrateEnabled: true,
	})

	if err == nil {
		t.Fatalf("expected error for nil DB, got nil (store=%v)", store)
	}
}

// func TestWithAutoMigrate(t *testing.T) {
// 	db := InitDB()

// 	// Initializes automigrateEnabled to False
// 	s := Store{
// 		logTableName:       "log_with_automigrate_false",
// 		db:                 db,
// 		automigrateEnabled: false,
// 	}

// 	// Modified to True
// 	f := WithAutoMigrate(true)
// 	f(&s)

// 	// Test Results
// 	if s.automigrateEnabled != true {
// 		t.Fatal("automigrateEnabled: Expected [true] received [%v]", s.automigrateEnabled)
// 	}

// 	// Initializes automigrateEnabled to True
// 	s = Store{
// 		logTableName:       "log_with_automigrate_true",
// 		db:                 db,
// 		automigrateEnabled: true,
// 	}

// 	// Modified to True
// 	f = WithAutoMigrate(false)
// 	f(&s)

// 	// Test Results
// 	if s.automigrateEnabled == true {
// 		t.Fatal("automigrateEnabled: Expected [true] received [%v]", s.automigrateEnabled)
// 	}
// }

// func TestWithDb(t *testing.T) {
// 	db := InitDB()

// 	s := Store{
// 		logTableName:       "LogTable",
// 		db:                 nil,
// 		automigrateEnabled: false,
// 	}

// 	f := WithDb(db)

// 	// DB has to be initialized now
// 	f(&s)

// 	// db non Nil expected
// 	if s.db == nil {
// 		t.Fatal("db initialization failed")
// 	}
// }

// func TestWithTableName(t *testing.T) {
// 	s := Store{
// 		logTableName:       "",
// 		db:                 nil,
// 		automigrateEnabled: false,
// 	}
// 	// TC: 1
// 	table_name := "Table1"
// 	f := WithTableName(table_name)
// 	f(&s)
// 	if s.logTableName != table_name {
// 		t.Fatal("Expected logTableName [%v], received [%v]", table_name, s.logTableName)
// 	}
// 	// TC: 2
// 	table_name = "Table2"
// 	f = WithTableName(table_name)
// 	f(&s)
// 	if s.logTableName != table_name {
// 		t.Fatal("Expected logTableName [%v], received [%v]", table_name, s.logTableName)
// 	}
// }

func Test_Store_AutoMigrate(t *testing.T) {
	db := InitDB()

	// Initializes automigrateEnabled to False
	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log_with_automigrate",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	errAutomigrate := s.AutoMigrate()

	if errAutomigrate != nil {
		t.Fatal("Store could not be automigrated: " + errAutomigrate.Error())
	}
}

func Test_Store_Log(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	now := time.Now().UTC()
	logEntry := NewLogWithData(
		uid.HumanUid(),
		LEVEL_DEBUG,
		"Test Message",
		"Test Context",
		now,
	)

	err = s.Log(logEntry)
	if err != nil {
		t.Fatal("Unexpected error: ", err.Error())
	}
}

func Test_Store_Debug(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	err = s.Debug("debug")
	if err != nil {
		t.Fatal("Unexpected error: ", err.Error())
	}
}

func Test_Store_DebugWithContext(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	err = s.DebugWithContext("debug", "Debug Message")
	if err != nil {
		t.Fatal("Unexpected error: ", err.Error())
	}
}

func Test_Store_Error(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	err = s.Error("error")
	if err != nil {
		t.Fatal("Unexpected error: ", err.Error())
	}
}

func Test_Store_ErrorWithContext(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	err = s.ErrorWithContext("error", "Error Message")
	if err != nil {
		t.Fatal("Unexpected error: ", err.Error())
	}
}

// Fatal methods uses system level API to terminate program (os.Exit)
func Test_Store_Fatal(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	err = s.Fatal("fatal")
	if err != nil {
		t.Fatal("Unexpected error: ", err.Error())
	}
}

func Test_Store_FatalWithContext(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	err = s.FatalWithContext("fatal", "Fatal Message")
	if err != nil {
		t.Fatal("Unexpected error: ", err.Error())
	}
}

func Test_Store_Info(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	err = s.Info("Info")
	if err != nil {
		t.Fatal("Unexpected error: ", err.Error())
	}
}

func Test_Store_InfoWithContext(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	err = s.InfoWithContext("Info", "Info Message")
	if err != nil {
		t.Fatal("Unexpected error: ", err.Error())
	}
}

func Test_Store_Trace(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	err = s.Trace("trace")
	if err != nil {
		t.Fatal("Unexpected error: ", err.Error())
	}
}

func Test_Store_TraceWithContext(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	err = s.TraceWithContext("trace", "Trace Message")
	if err != nil {
		t.Fatal("Unexpected error: ", err.Error())
	}
}

func Test_Store_Warn(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	err = s.Warn("warn")
	if err != nil {
		t.Fatal("Unexpected error: ", err.Error())
	}
}

func Test_Store_WarnWithContext(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	err = s.WarnWithContext("warn", "Warning Message")
	if err != nil {
		t.Fatal("Unexpected error: ", err.Error())
	}
}

func Test_Store_GetDriverNameAndLogTableName(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:           db,
		LogTableName: "log_getters",
		DbDriverName: "sqlite3",
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	if s.GetDriverName() != "sqlite3" {
		t.Fatal("GetDriverName returned unexpected value")
	}

	if s.GetLogTableName() != "log_getters" {
		t.Fatal("GetLogTableName returned unexpected value")
	}
}

func Test_Store_LogList(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log_list",
		DbDriverName:       "sqlite3",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	err = s.Debug("debug message")
	if err != nil {
		t.Fatal("Unexpected error: ", err.Error())
	}

	err = s.Error("error message")
	if err != nil {
		t.Fatal("Unexpected error: ", err.Error())
	}

	query := LogQuery().SetLevel(LEVEL_DEBUG)
	logs, err := s.LogList(query)
	if err != nil {
		t.Fatal("Unexpected error from LogList: ", err.Error())
	}

	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	if logs[0].GetLevel() != LEVEL_DEBUG {
		t.Fatal("LogList did not return the expected log level")
	}
}

func Test_Store_LogCreateAndFindByID(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log_create_find",
		DbDriverName:       "sqlite3",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	entry := NewLog().
		SetLevel(LEVEL_INFO).
		SetMessage("create and find test")

	if err := s.LogCreate(entry); err != nil {
		t.Fatal("Unexpected error from LogCreate: ", err.Error())
	}

	id := entry.GetID()
	if id == "" {
		t.Fatal("LogCreate did not assign an ID")
	}

	found, err := s.LogFindByID(id)
	if err != nil {
		t.Fatal("Unexpected error from LogFindByID: ", err.Error())
	}

	if found == nil {
		t.Fatal("LogFindByID returned nil entry")
	}

	if found.GetID() != id {
		t.Fatal("LogFindByID returned entry with unexpected ID")
	}
}

func Test_Store_LogDelete(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log_delete",
		DbDriverName:       "sqlite3",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	entry := NewLog().
		SetLevel(LEVEL_ERROR).
		SetMessage("delete test")

	if err := s.LogCreate(entry); err != nil {
		t.Fatal("Unexpected error from LogCreate: ", err.Error())
	}

	id := entry.GetID()
	if id == "" {
		t.Fatal("LogCreate did not assign an ID")
	}

	if err := s.LogDelete(entry); err != nil {
		t.Fatal("Unexpected error from LogDelete: ", err.Error())
	}

	found, err := s.LogFindByID(id)
	if err != nil {
		t.Fatal("Unexpected error from LogFindByID: ", err.Error())
	}

	if found != nil {
		t.Fatal("LogDelete did not remove the entry")
	}
}

func Test_Store_LogDeleteByID(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log_delete_by_id",
		DbDriverName:       "sqlite3",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	entry := NewLog().
		SetLevel(LEVEL_WARNING).
		SetMessage("delete by id test")

	if err := s.LogCreate(entry); err != nil {
		t.Fatal("Unexpected error from LogCreate: ", err.Error())
	}

	id := entry.GetID()
	if id == "" {
		t.Fatal("LogCreate did not assign an ID")
	}

	if err := s.LogDeleteByID(id); err != nil {
		t.Fatal("Unexpected error from LogDeleteByID: ", err.Error())
	}

	found, err := s.LogFindByID(id)
	if err != nil {
		t.Fatal("Unexpected error from LogFindByID: ", err.Error())
	}

	if found != nil {
		t.Fatal("LogDeleteByID did not remove the entry")
	}
}

func Test_Store_LogDeleteByID_Error_EmptyID(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log_delete_by_id_error",
		DbDriverName:       "sqlite3",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	if err := s.LogDeleteByID(""); err == nil {
		t.Fatal("expected error from LogDeleteByID(\"\"), got nil")
	}
}

func Test_Store_LogCount_Basic(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log_count_basic",
		DbDriverName:       "sqlite3",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	// create two debug logs and one error log
	if err := s.Debug("debug 1"); err != nil {
		t.Fatalf("unexpected error creating debug log: %v", err)
	}
	if err := s.Debug("debug 2"); err != nil {
		t.Fatalf("unexpected error creating debug log: %v", err)
	}
	if err := s.Error("error 1"); err != nil {
		t.Fatalf("unexpected error creating error log: %v", err)
	}

	// count all logs
	countAll, err := s.LogCount(LogQuery())
	if err != nil {
		t.Fatalf("unexpected error from LogCount (all): %v", err)
	}
	if countAll != 3 {
		t.Fatalf("expected 3 logs in total, got %d", countAll)
	}

	// count only debug logs
	countDebug, err := s.LogCount(LogQuery().SetLevel(LEVEL_DEBUG))
	if err != nil {
		t.Fatalf("unexpected error from LogCount (debug): %v", err)
	}
	if countDebug != 2 {
		t.Fatalf("expected 2 debug logs, got %d", countDebug)
	}
}

func Test_Store_LogCount_IgnoresLimitAndOffset(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log_count_paging",
		DbDriverName:       "sqlite3",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	// create several logs
	for i := 0; i < 5; i++ {
		if err := s.Info("info log"); err != nil {
			t.Fatalf("unexpected error creating info log: %v", err)
		}
	}

	// apply limit/offset on the query; LogCount should still return total
	query := LogQuery().
		SetLevel(LEVEL_INFO).
		SetLimit(2).
		SetOffset(1)

	count, err := s.LogCount(query)
	if err != nil {
		t.Fatalf("unexpected error from LogCount with paging: %v", err)
	}
	if count != 5 {
		t.Fatalf("expected LogCount to ignore limit/offset and return 5, got %d", count)
	}
}

func Test_Store_LogFindByID_Error_EmptyID(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log_find_by_id_error",
		DbDriverName:       "sqlite3",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	if _, err := s.LogFindByID(""); err == nil {
		t.Fatal("expected error from LogFindByID(\"\"), got nil")
	}
}

func Test_Store_LogCount_Error_FromInvalidQuery(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log_count_error",
		DbDriverName:       "sqlite3",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	query := LogQuery().SetLimit(-1)
	if _, err := s.LogCount(query); err == nil {
		t.Fatal("expected error from LogCount with invalid query (negative limit), got nil")
	}
}

func Test_Store_LogList_NoMatches(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log_list_no_matches",
		DbDriverName:       "sqlite3",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	if err := s.Info("info only"); err != nil {
		t.Fatalf("unexpected error creating info log: %v", err)
	}

	logs, err := s.LogList(LogQuery().SetLevel(LEVEL_ERROR))
	if err != nil {
		t.Fatalf("unexpected error from LogList with no matches: %v", err)
	}
	if len(logs) != 0 {
		t.Fatalf("expected 0 logs for non-matching query, got %d", len(logs))
	}
}

func Test_Store_LogCount_NoMatches(t *testing.T) {
	db := InitDB()

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log_count_no_matches",
		DbDriverName:       "sqlite3",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	if err := s.Info("info only"); err != nil {
		t.Fatalf("unexpected error creating info log: %v", err)
	}

	count, err := s.LogCount(LogQuery().SetLevel(LEVEL_ERROR))
	if err != nil {
		t.Fatalf("unexpected error from LogCount with no matches: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected LogCount to return 0 for non-matching query, got %d", count)
	}
}
