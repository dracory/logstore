package logstore

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlserver"
	"github.com/dracory/sb"
	"github.com/dracory/uid"
	"github.com/dromara/carbon/v2"
)

// Store defines a session store
type storeImplementation struct {
	logTableName       string
	db                 *sql.DB
	dbDriverName       string
	automigrateEnabled bool
	debugEnabled       bool
}

func (st *storeImplementation) GetDriverName() string {
	return st.dbDriverName
}

func (st *storeImplementation) GetLogTableName() string {
	return st.logTableName
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

// Log adds a log (shortcut for LogCreate)
func (st *storeImplementation) Log(logEntry LogInterface) error {
	if logEntry == nil {
		return errors.New("log entry is nil")
	}

	return st.LogCreate(logEntry)
}

// Debug adds a debug log
func (st *storeImplementation) Debug(message string) error {
	logEntry := NewLog().
		SetLevel(LEVEL_DEBUG).
		SetMessage(message)
	return st.Log(logEntry)
}

// DebugWithContext adds a debug log with context data
func (st *storeImplementation) DebugWithContext(message string, context interface{}) error {
	contextBytes, err := json.Marshal(context)

	if err != nil {
		log.Println(err)
		contextBytes = []byte("JSON encode error")
	}

	logEntry := NewLog().
		SetLevel(LEVEL_DEBUG).
		SetMessage(message).
		SetContext(string(contextBytes))
	return st.Log(logEntry)
}

// Error adds an error log
func (st *storeImplementation) Error(message string) error {
	logEntry := NewLog().
		SetLevel(LEVEL_ERROR).
		SetMessage(message)
	return st.Log(logEntry)
}

// ErrorWithContext adds an error log with context data
func (st *storeImplementation) ErrorWithContext(message string, context interface{}) error {
	contextBytes, err := json.Marshal(context)

	if err != nil {
		log.Println(err)
		contextBytes = []byte("JSON encode error")
	}

	logEntry := NewLog().
		SetLevel(LEVEL_ERROR).
		SetMessage(message).
		SetContext(string(contextBytes))
	return st.Log(logEntry)
}

// Fatal adds an fatal log and calls os.Exit(1) after logging
func (st *storeImplementation) Fatal(message string) error {
	logEntry := NewLog().
		SetLevel(LEVEL_FATAL).
		SetMessage(message)

	err := st.Log(logEntry)
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

	logEntry := NewLog().
		SetLevel(LEVEL_FATAL).
		SetMessage(message).
		SetContext(string(contextBytes))

	err = st.Log(logEntry)
	// os.Exit(1)
	return err
}

// Info adds an info log
func (st *storeImplementation) Info(message string) error {
	logEntry := NewLog().
		SetLevel(LEVEL_INFO).
		SetMessage(message)
	return st.Log(logEntry)
}

// InfoWithContext adds an info log with context data
func (st *storeImplementation) InfoWithContext(message string, context interface{}) error {
	contextBytes, err := json.Marshal(context)

	if err != nil {
		log.Println(err)
		contextBytes = []byte("JSON encode error")
	}

	logEntry := NewLog().
		SetLevel(LEVEL_INFO).
		SetMessage(message).
		SetContext(string(contextBytes))
	return st.Log(logEntry)
}

// Panic adds an panic log and calls panic(message) after logging
func (st *storeImplementation) Panic(message string) {
	logEntry := NewLog().
		SetLevel(LEVEL_PANIC).
		SetMessage(message)

	st.Log(logEntry)
	panic(message)
}

// PanicWithContext adds a panic log with context data and calls panic(message) after logging
func (st *storeImplementation) PanicWithContext(message string, context interface{}) {
	contextBytes, err := json.Marshal(context)

	if err != nil {
		log.Println(err)
		contextBytes = []byte("JSON encode error")
	}

	logEntry := NewLog().
		SetLevel(LEVEL_FATAL).
		SetMessage(message).
		SetContext(string(contextBytes))

	st.Log(logEntry)
	panic(message)
}

// Trace adds a trace log
func (st *storeImplementation) Trace(message string) error {
	logEntry := NewLog().
		SetLevel(LEVEL_TRACE).
		SetMessage(message)

	return st.Log(logEntry)
}

// TraceWithContext adds a trace log with context data
func (st *storeImplementation) TraceWithContext(message string, context interface{}) error {
	contextBytes, err := json.Marshal(context)

	if err != nil {
		log.Println(err)
		contextBytes = []byte("JSON encode error")
	}

	logEntry := NewLog().
		SetLevel(LEVEL_TRACE).
		SetMessage(message).
		SetContext(string(contextBytes))

	return st.Log(logEntry)
}

// Warn adds a warn log
func (st *storeImplementation) Warn(message string) error {
	logEntry := NewLog().
		SetLevel(LEVEL_WARNING).
		SetMessage(message)

	return st.Log(logEntry)
}

// WarnWithContext adds a warn log with context data
func (st *storeImplementation) WarnWithContext(message string, context interface{}) error {
	contextBytes, err := json.Marshal(context)

	if err != nil {
		log.Println(err)
		contextBytes = []byte("JSON encode error")
	}

	logEntry := NewLog().
		SetLevel(LEVEL_WARNING).
		SetMessage(message).
		SetContext(string(contextBytes))

	return st.Log(logEntry)
}

// LogCreate adds a log
func (st *storeImplementation) LogCreate(logEntry LogInterface) error {
	if logEntry == nil {
		return errors.New("log entry is nil")
	}

	id := logEntry.GetID()
	if id == "" {
		id = uid.MicroUid()
		logEntry.SetID(id)
	}

	t := logEntry.GetTime()
	if t.IsZero() {
		t = time.Now().UTC()
		logEntry.SetTime(t)
	}

	sqlStr, sqlParams, err := goqu.Dialect(st.dbDriverName).
		Insert(st.logTableName).
		Rows(struct {
			ID      string    `db:"id"`
			Level   string    `db:"level"`
			Message string    `db:"message"`
			Context string    `db:"context"`
			Time    time.Time `db:"time"`
		}{
			ID:      logEntry.GetID(),
			Level:   logEntry.GetLevel(),
			Message: logEntry.GetMessage(),
			Context: logEntry.GetContext(),
			Time:    logEntry.GetTime(),
		}).
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

// LogDelete deletes a log
func (st *storeImplementation) LogDelete(logEntry LogInterface) error {
	if logEntry == nil {
		return errors.New("log entry is nil")
	}

	return st.LogDeleteByID(logEntry.GetID())
}

// LogDeleteByID deletes a log by ID
func (st *storeImplementation) LogDeleteByID(id string) error {
	if id == "" {
		return errors.New("log id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(st.dbDriverName).
		Delete(st.logTableName).
		Prepared(true).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr, params...)

	return err
}

// LogFindByID finds a log by ID
func (st *storeImplementation) LogFindByID(id string) (LogInterface, error) {
	if id == "" {
		return nil, errors.New("log id is empty")
	}

	list, err := st.LogList(LogQuery().
		SetID(id).
		SetLimit(1))

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (st *storeImplementation) LogList(query LogQueryInterface) ([]LogInterface, error) {
	if query == nil {
		query = LogQuery()
	}

	q, columns, err := query.ToSelectDataset(st)
	if err != nil {
		return []LogInterface{}, err
	}

	sqlStr, sqlParams, errSql := q.Prepared(true).Select(columns...).ToSQL()
	if errSql != nil {
		return []LogInterface{}, errSql
	}

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	db := sb.NewDatabase(st.db, st.dbDriverName)
	modelMaps, err := db.SelectToMapString(sqlStr, sqlParams...)
	if err != nil {
		return []LogInterface{}, err
	}

	list := []LogInterface{}

	for _, modelMap := range modelMaps {
		id := modelMap[COLUMN_ID]
		level := modelMap[COLUMN_LEVEL]
		message := modelMap[COLUMN_MESSAGE]
		context := modelMap[COLUMN_CONTEXT]
		var t time.Time
		if v, ok := modelMap[COLUMN_TIME]; ok && v != "" {
			parsed := carbon.Parse(v, carbon.UTC).StdTime()
			t = parsed
		}

		logEntry := NewLogWithData(id, level, message, context, t)
		list = append(list, logEntry)
	}

	return list, nil
}
