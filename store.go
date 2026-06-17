package logstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/dracory/neat"
	contractsorm "github.com/dracory/neat/contracts/database/orm"
	contractsschema "github.com/dracory/neat/contracts/database/schema"
	neatuid "github.com/dracory/neat/support/uid"
	"github.com/dromara/carbon/v2"
)

// == INTERFACE ===============================================================

// StoreInterface defines the interface for a log store
type StoreInterface interface {
	// GetDB returns the underlying *sql.DB
	GetDB() *sql.DB

	// GetLogTableName returns the log table name
	GetLogTableName() string
	// SetLogTableName sets the log table name
	SetLogTableName(logTableName string)

	// MigrateDown drops the log table
	MigrateDown(ctx context.Context, tx ...*sql.Tx) error

	// MigrateUp creates the log table
	MigrateUp(ctx context.Context, tx ...*sql.Tx) error

	// EnableDebug enables or disables debug mode
	EnableDebug(debug bool)

	// Log adds a log entry
	Log(logEntry LogInterface) error

	// Debug adds a debug log
	Debug(message string) error

	// DebugWithContext adds a debug log with context data
	DebugWithContext(message string, context interface{}) error

	// Error adds an error log
	Error(message string) error

	// ErrorWithContext adds an error log with context data
	ErrorWithContext(message string, context interface{}) error

	// Fatal adds a fatal log
	Fatal(message string) error

	// FatalWithContext adds a fatal log with context data
	FatalWithContext(message string, context interface{}) error

	// Info adds an info log
	Info(message string) error

	// InfoWithContext adds an info log with context data
	InfoWithContext(message string, context interface{}) error

	// Panic adds a panic log and calls panic(message) after logging
	Panic(message string)

	// PanicWithContext adds a panic log with context data and calls panic(message) after logging
	PanicWithContext(message string, context interface{})

	// Trace adds a trace log
	Trace(message string) error

	// TraceWithContext adds a trace log with context data
	TraceWithContext(message string, context interface{}) error

	// Warn adds a warn log
	Warn(message string) error

	// WarnWithContext adds a warn log with context data
	WarnWithContext(message string, context interface{}) error

	LogCount(ctx context.Context, query LogQueryInterface) (int64, error)
	LogCreate(ctx context.Context, logEntry LogInterface) error
	LogList(ctx context.Context, query LogQueryInterface) ([]LogInterface, error)
	LogDelete(ctx context.Context, logEntry LogInterface) error
	LogDeleteByID(ctx context.Context, id string) error
	LogDeleteByIDs(ctx context.Context, ids []string) error
	LogFindByID(ctx context.Context, id string) (LogInterface, error)
}

// == TYPE ====================================================================

var _ StoreInterface = (*storeImplementation)(nil)

// storeImplementation implements StoreInterface for log operations.
type storeImplementation struct {
	logTableName       string
	db                 *neat.Database
	automigrateEnabled bool
	debugEnabled       bool
	logger             *slog.Logger
}

// NewStoreOptions define the options for creating a new log store
type NewStoreOptions struct {
	LogTableName       string
	DB                 *sql.DB
	AutomigrateEnabled bool
	DebugEnabled       bool
}

// NewStore creates a new log store
func NewStore(opts NewStoreOptions) (StoreInterface, error) {
	if opts.DB == nil {
		return nil, errors.New("log store: DB is required")
	}

	if opts.LogTableName == "" {
		return nil, errors.New("log store: logTableName is required")
	}

	neatDB, err := neat.NewFromSQLDB(opts.DB)
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	store := &storeImplementation{
		logTableName:       opts.LogTableName,
		db:                 neatDB,
		automigrateEnabled: opts.AutomigrateEnabled,
		debugEnabled:       opts.DebugEnabled,
		logger:             logger,
	}

	if store.automigrateEnabled {
		if err := store.MigrateUp(context.Background()); err != nil {
			return nil, err
		}
	}

	return store, nil
}

// == DB ======================================================================

// GetDB returns the underlying *sql.DB.
func (st *storeImplementation) GetDB() *sql.DB {
	db, _ := st.db.DB()
	return db
}

// == TABLE NAME ==============================================================

// GetLogTableName returns the log table name
func (st *storeImplementation) GetLogTableName() string {
	return st.logTableName
}

// SetLogTableName sets the log table name
func (st *storeImplementation) SetLogTableName(logTableName string) {
	st.logTableName = logTableName
}

// == MIGRATE =================================================================

// MigrateUp creates the log table
func (st *storeImplementation) MigrateUp(ctx context.Context, tx ...*sql.Tx) error {
	if st.db.Schema().HasTable(st.logTableName) {
		if st.debugEnabled {
			st.logger.Info("MigrateUp: table already exists", "table", st.logTableName)
		}
		return nil
	}

	err := st.db.Schema().Create(st.logTableName, func(table contractsschema.Blueprint) {
		table.String(COLUMN_ID, 40)
		table.Primary(COLUMN_ID)
		table.String(COLUMN_LEVEL, 20)
		table.Text(COLUMN_MESSAGE)
		table.Text(COLUMN_CONTEXT)
		table.DateTime(COLUMN_TIME)
	})

	if err != nil {
		if st.debugEnabled {
			st.logger.Error("MigrateUp failed", "error", err)
		}
		return err
	}

	return nil
}

// MigrateDown drops the log table
func (st *storeImplementation) MigrateDown(ctx context.Context, tx ...*sql.Tx) error {
	if !st.db.Schema().HasTable(st.logTableName) {
		if st.debugEnabled {
			st.logger.Info("MigrateDown: table does not exist", "table", st.logTableName)
		}
		return nil
	}

	err := st.db.Schema().Drop(st.logTableName)
	if err != nil {
		if st.debugEnabled {
			st.logger.Error("MigrateDown failed", "error", err)
		}
		return err
	}
	return nil
}

// AutoMigrate auto migrate (deprecated - use MigrateUp)
func (st *storeImplementation) AutoMigrate() error {
	return st.MigrateUp(context.Background())
}

// == DEBUG ===================================================================

// EnableDebug enables or disables debug mode.
func (st *storeImplementation) EnableDebug(debug bool) {
	st.debugEnabled = debug
	if debug {
		st.db.EnableDebug()
		st.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		st.db.DisableDebug()
		st.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
}

// == CONVENIENCE LOGGERS =====================================================

// Log adds a log (shortcut for LogCreate)
func (st *storeImplementation) Log(logEntry LogInterface) error {
	if logEntry == nil {
		return errors.New("log entry is nil")
	}

	return st.LogCreate(context.Background(), logEntry)
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
		st.logger.Error("JSON encode error", "error", err)
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
		st.logger.Error("JSON encode error", "error", err)
		contextBytes = []byte("JSON encode error")
	}

	logEntry := NewLog().
		SetLevel(LEVEL_ERROR).
		SetMessage(message).
		SetContext(string(contextBytes))
	return st.Log(logEntry)
}

// Fatal adds an fatal log
func (st *storeImplementation) Fatal(message string) error {
	logEntry := NewLog().
		SetLevel(LEVEL_FATAL).
		SetMessage(message)

	err := st.Log(logEntry)
	// os.Exit(1)
	return err
}

// FatalWithContext adds a fatal log with context data
func (st *storeImplementation) FatalWithContext(message string, context interface{}) error {
	contextBytes, err := json.Marshal(context)
	if err != nil {
		st.logger.Error("JSON encode error", "error", err)
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
		st.logger.Error("JSON encode error", "error", err)
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
		st.logger.Error("JSON encode error", "error", err)
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
		st.logger.Error("JSON encode error", "error", err)
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
		st.logger.Error("JSON encode error", "error", err)
		contextBytes = []byte("JSON encode error")
	}

	logEntry := NewLog().
		SetLevel(LEVEL_WARNING).
		SetMessage(message).
		SetContext(string(contextBytes))

	return st.Log(logEntry)
}

// == CRUD ====================================================================

// LogCreate adds a log
func (st *storeImplementation) LogCreate(ctx context.Context, logEntry LogInterface) error {
	if logEntry == nil {
		return errors.New("log entry is nil")
	}

	id := logEntry.GetID()
	if id == "" {
		id = neatuid.GenerateShortID()
		logEntry.SetID(id)
	}

	t := logEntry.GetTime()
	if t.IsZero() {
		t = time.Now().UTC()
		logEntry.SetTime(t)
	}

	row := map[string]any{
		COLUMN_ID:      logEntry.GetID(),
		COLUMN_LEVEL:   logEntry.GetLevel(),
		COLUMN_MESSAGE: logEntry.GetMessage(),
		COLUMN_CONTEXT: logEntry.GetContext(),
		COLUMN_TIME:    logEntry.GetTime(),
	}

	return st.db.Query().Table(st.logTableName).Create(row)
}

// LogDelete deletes a log
func (st *storeImplementation) LogDelete(ctx context.Context, logEntry LogInterface) error {
	if logEntry == nil {
		return errors.New("log entry is nil")
	}

	return st.LogDeleteByID(ctx, logEntry.GetID())
}

// LogDeleteByID deletes a log by ID
func (st *storeImplementation) LogDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("log id is empty")
	}

	_, err := st.db.Query().
		Table(st.logTableName).
		Where(COLUMN_ID+" = ?", id).
		Delete()

	return err
}

// LogDeleteByIDs deletes multiple logs by their IDs in a single query
func (st *storeImplementation) LogDeleteByIDs(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	_, err := st.db.Query().
		Table(st.logTableName).
		WhereIn(COLUMN_ID, args).
		Delete()

	return err
}

// LogFindByID finds a log by ID
func (st *storeImplementation) LogFindByID(ctx context.Context, id string) (LogInterface, error) {
	if id == "" {
		return nil, errors.New("log id is empty")
	}

	list, err := st.LogList(ctx, LogQuery().
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

// LogList lists logs based on a query
func (st *storeImplementation) LogList(ctx context.Context, query LogQueryInterface) ([]LogInterface, error) {
	if query == nil {
		query = LogQuery()
	}

	if err := query.Validate(); err != nil {
		return []LogInterface{}, err
	}

	q := st.buildQuery(query)

	var results []map[string]any
	err := q.Get(&results)
	if err != nil {
		return []LogInterface{}, err
	}

	list := []LogInterface{}
	for _, result := range results {
		id := ""
		if v, ok := result[COLUMN_ID].(string); ok {
			id = v
		}
		level := ""
		if v, ok := result[COLUMN_LEVEL].(string); ok {
			level = v
		}
		message := ""
		if v, ok := result[COLUMN_MESSAGE].(string); ok {
			message = v
		}
		contextStr := ""
		if v, ok := result[COLUMN_CONTEXT].(string); ok {
			contextStr = v
		}

		var t time.Time
		if v, ok := result[COLUMN_TIME]; ok {
			switch vt := v.(type) {
			case time.Time:
				t = vt
			case string:
				t = carbon.Parse(vt, carbon.UTC).StdTime()
			}
		}

		logEntry := NewLogWithData(id, level, message, contextStr, t)
		list = append(list, logEntry)
	}

	return list, nil
}

// LogCount returns the total number of logs that match the given query
func (st *storeImplementation) LogCount(ctx context.Context, query LogQueryInterface) (int64, error) {
	if query == nil {
		query = LogQuery()
	}

	if err := query.Validate(); err != nil {
		return 0, err
	}

	q := st.buildQuery(query)

	var count int64
	err := q.Count(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// == QUERY BUILDER ==========================================================

// buildQuery builds a neat query from the log query interface.
func (st *storeImplementation) buildQuery(query LogQueryInterface) contractsorm.Query {
	q := st.db.Query().Table(st.logTableName)

	if query == nil {
		return q
	}

	if query.IsIDSet() && query.GetID() != "" {
		q = q.Where(COLUMN_ID+" = ?", query.GetID())
	}

	if query.IsIDInSet() && len(query.GetIDIn()) > 0 {
		args := make([]any, len(query.GetIDIn()))
		for i, id := range query.GetIDIn() {
			args[i] = id
		}
		q = q.WhereIn(COLUMN_ID, args)
	}

	if query.IsLevelSet() && query.GetLevel() != "" {
		q = q.Where(COLUMN_LEVEL+" = ?", query.GetLevel())
	}

	if query.IsLevelInSet() && len(query.GetLevelIn()) > 0 {
		args := make([]any, len(query.GetLevelIn()))
		for i, level := range query.GetLevelIn() {
			args[i] = level
		}
		q = q.WhereIn(COLUMN_LEVEL, args)
	}

	if query.IsMessageContainsSet() && query.GetMessageContains() != "" {
		q = q.Where(COLUMN_MESSAGE+" LIKE ?", "%"+query.GetMessageContains()+"%")
	}

	if query.IsMessageNotContainsSet() && query.GetMessageNotContains() != "" {
		q = q.Where(COLUMN_MESSAGE+" NOT LIKE ?", "%"+query.GetMessageNotContains()+"%")
	}

	if query.IsContextContainsSet() && query.GetContextContains() != "" {
		q = q.Where(COLUMN_CONTEXT+" LIKE ?", "%"+query.GetContextContains()+"%")
	}

	if query.IsContextNotContainsSet() && query.GetContextNotContains() != "" {
		q = q.Where(COLUMN_CONTEXT+" NOT LIKE ?", "%"+query.GetContextNotContains()+"%")
	}

	if query.IsTimeGteSet() && query.GetTimeGte() != "" {
		q = q.Where(COLUMN_TIME+" >= ?", query.GetTimeGte())
	}

	if query.IsTimeLteSet() && query.GetTimeLte() != "" {
		q = q.Where(COLUMN_TIME+" <= ?", query.GetTimeLte())
	}

	if query.IsLimitSet() && query.GetLimit() > 0 {
		q = q.Limit(query.GetLimit())
	}

	if query.IsOffsetSet() && query.GetOffset() > 0 {
		q = q.Offset(query.GetOffset())
	}

	if query.IsOrderBySet() && query.GetOrderBy() != "" {
		direction := "desc"
		if query.IsOrderDirectionSet() && query.GetOrderDirection() != "" {
			direction = query.GetOrderDirection()
		}
		q = q.OrderBy(query.GetOrderBy(), direction)
	}

	return q
}
