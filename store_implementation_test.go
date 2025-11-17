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

	query := LogQuery().SetLevel(LevelDebug)
	logs, err := s.LogList(query)
	if err != nil {
		t.Fatal("Unexpected error from LogList: ", err.Error())
	}

	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	if logs[0].Level != LevelDebug {
		t.Fatal("LogList did not return the expected log level")
	}
}
