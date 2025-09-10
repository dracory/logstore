package logstore

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/dracory/uid"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filepath string) *sql.DB {
	os.Remove(filepath) // remove database
	dsn := filepath + "?parseTime=true"
	db, err := sql.Open("sqlite3", dsn)

	if err != nil {
		panic(err)
	}

	return db
}

func TestStoreCreate(t *testing.T) {
	db := InitDB("test_log_store_create.db")

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
// 	db := InitDB("test_log_store_automigrate.db")

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
// 	db := InitDB("test_log_store_with_automigrate.db")

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
	db := InitDB("test_log_store_automigrate.db")

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
	db := InitDB("test_log_store_log.db")

	s, err := NewStore(NewStoreOptions{
		DB:                 db,
		LogTableName:       "log",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("Store could not be created: " + err.Error())
	}

	time := time.Now()
	log := Log{
		ID:      uid.HumanUid(),
		Level:   LevelDebug,
		Message: "Test Message",
		Context: "Test Context",
		Time:    &time,
	}

	err = s.Log(&log)
	if err != nil {
		t.Fatal("Unexpected error: ", err.Error())
	}
}

func Test_Store_Debug(t *testing.T) {
	db := InitDB("test_log_store_debug.db")

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
	db := InitDB("test_log_store_debug_with_cotext.db")

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
	db := InitDB("test_log_store_log.db")

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
	db := InitDB("test_log_store_error_with_context.db")

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
	db := InitDB("test_log_store_log.db")

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
	db := InitDB("test_log_store_log.db")

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
	db := InitDB("test_log_store_info.db")

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
	db := InitDB("test_log_store_info_with_context.db")

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
	db := InitDB("test_log_store_trace.db")

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
	db := InitDB("test_log_store_trace_with_context.db")

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
	db := InitDB("test_log_store_warn.db")

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
	db := InitDB("test_log_store_warn_with_context.db")

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
