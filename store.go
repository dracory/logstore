package logstore

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlserver"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/sb"
)

// Store defines a session store
type storeImplementation struct {
	logTableName       string
	db                 *sql.DB
	dbDriverName       string
	automigrateEnabled bool
	debugEnabled       bool
}

// NewStoreOptions define the options for creating a new session store
type NewStoreOptions struct {
	LogTableName       string
	DB                 *sql.DB
	DbDriverName       string
	AutomigrateEnabled bool
	DebugEnabled       bool
}

// NewStore creates a new session store
func NewStore(opts NewStoreOptions) (*storeImplementation, error) {
	store := &storeImplementation{
		logTableName:       opts.LogTableName,
		automigrateEnabled: opts.AutomigrateEnabled,
		db:                 opts.DB,
		dbDriverName:       opts.DbDriverName,
		debugEnabled:       opts.DebugEnabled,
	}

	if store.logTableName == "" {
		return nil, errors.New("log store: logTableName is required")
	}

	if store.db == nil {
		return nil, errors.New("log store: DB is required")
	}

	if store.dbDriverName == "" {
		store.dbDriverName = sb.DatabaseDriverName(store.db)
	}

	if store.automigrateEnabled {
		store.AutoMigrate()
	}

	return store, nil
}

// AutoMigrate auto migrate
func (st *storeImplementation) AutoMigrate() error {
	sql := st.SqlCreateTable()

	if st.debugEnabled {
		log.Println(sql)
	}

	_, err := st.db.Exec(sql)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// EnableDebug - enables the debug option
func (st *storeImplementation) EnableDebug(debug bool) {
	st.debugEnabled = debug
}

// Log adds a log
func (st *storeImplementation) Log(logEntry *Log) error {
	if logEntry.ID == "" {
		logEntry.ID = uid.MicroUid()
	}
	if logEntry.Time == nil {
		t := carbon.Now(carbon.UTC).StdTime()
		logEntry.Time = &t
	}

	sqlStr, sqlParams, err := goqu.Dialect(st.dbDriverName).
		Insert(st.logTableName).
		Rows(logEntry).
		Prepared(true).
		ToSQL()

	if err != nil {
		return err
	}

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	_, err = st.db.Exec(sqlStr, sqlParams...)

	if err != nil {
		if st.debugEnabled {
			log.Println(err.Error())
		}
		return err
	}

	return nil
}

// Debug adds a debug log
func (st *storeImplementation) Debug(message string) error {
	log := Log{
		Level:   LevelDebug,
		Message: message,
	}
	return st.Log(&log)
}

// DebugWithContext adds a debug log with context data
func (st *storeImplementation) DebugWithContext(message string, context interface{}) error {
	contextBytes, err := json.Marshal(context)

	if err != nil {
		log.Println(err)
		contextBytes = []byte("JSON encode error")
	}

	log := Log{
		Level:   LevelDebug,
		Message: message,
		Context: string(contextBytes),
	}
	return st.Log(&log)
}

// Error adds an error log
func (st *storeImplementation) Error(message string) error {
	log := Log{
		Level:   LevelError,
		Message: message,
	}
	return st.Log(&log)
}

// ErrorWithContext adds an error log with context data
func (st *storeImplementation) ErrorWithContext(message string, context interface{}) error {
	contextBytes, err := json.Marshal(context)

	if err != nil {
		log.Println(err)
		contextBytes = []byte("JSON encode error")
	}

	log := Log{
		Level:   LevelError,
		Message: message,
		Context: string(contextBytes),
	}
	return st.Log(&log)
}

// Fatal adds an fatal log and calls os.Exit(1) after logging
func (st *storeImplementation) Fatal(message string) error {
	log := Log{
		Level:   LevelFatal,
		Message: message,
	}

	err := st.Log(&log)
	// os.Exit(1)
	return err
}

// FatalWithContext adds a fatal log with context data and calls os.Exit(1) after logging
func (st *storeImplementation) FatalWithContext(message string, context interface{}) error {
	contextBytes, err := json.Marshal(context)

	if err != nil {
		log.Println(err)
		contextBytes = []byte("JSON encode error")
	}

	log := Log{
		Level:   LevelFatal,
		Message: message,
		Context: string(contextBytes),
	}

	err = st.Log(&log)
	// os.Exit(1)
	return err
}

// Info adds an info log
func (st *storeImplementation) Info(message string) error {
	log := Log{
		Level:   LevelInfo,
		Message: message,
	}
	return st.Log(&log)
}

// InfoWithContext adds an info log with context data
func (st *storeImplementation) InfoWithContext(message string, context interface{}) error {
	contextBytes, err := json.Marshal(context)

	if err != nil {
		log.Println(err)
		contextBytes = []byte("JSON encode error")
	}

	log := Log{
		Level:   LevelInfo,
		Message: message,
		Context: string(contextBytes),
	}
	return st.Log(&log)
}

// Panic adds an panic log and calls panic(message) after logging
func (st *storeImplementation) Panic(message string) {
	log := Log{
		Level:   LevelPanic,
		Message: message,
	}

	st.Log(&log)
	panic(message)
}

// PanicWithContext adds a panic log with context data and calls panic(message) after logging
func (st *storeImplementation) PanicWithContext(message string, context interface{}) {
	contextBytes, err := json.Marshal(context)

	if err != nil {
		log.Println(err)
		contextBytes = []byte("JSON encode error")
	}

	log := Log{
		Level:   LevelFatal,
		Message: message,
		Context: string(contextBytes),
	}

	st.Log(&log)
	panic(message)
}

// Trace adds a trace log
func (st *storeImplementation) Trace(message string) error {
	log := Log{
		Level:   LevelTrace,
		Message: message,
	}

	return st.Log(&log)
}

// TraceWithContext adds a trace log with context data
func (st *storeImplementation) TraceWithContext(message string, context interface{}) error {
	contextBytes, err := json.Marshal(context)

	if err != nil {
		log.Println(err)
		contextBytes = []byte("JSON encode error")
	}

	log := Log{
		Level:   LevelTrace,
		Message: message,
		Context: string(contextBytes),
	}

	return st.Log(&log)
}

// Warn adds a warn log
func (st *storeImplementation) Warn(message string) error {
	log := Log{
		Level:   LevelWarning,
		Message: message,
	}

	return st.Log(&log)
}

// WarnWithContext adds a warn log with context data
func (st *storeImplementation) WarnWithContext(message string, context interface{}) error {
	contextBytes, err := json.Marshal(context)

	if err != nil {
		log.Println(err)
		contextBytes = []byte("JSON encode error")
	}

	log := Log{
		Level:   LevelWarning,
		Message: message,
		Context: string(contextBytes),
	}

	return st.Log(&log)
}
