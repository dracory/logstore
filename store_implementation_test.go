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
