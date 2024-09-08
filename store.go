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
	"github.com/golang-module/carbon/v2"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/uid"
)

// Store defines a session store
type Store struct {
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
func NewStore(opts NewStoreOptions) (*Store, error) {
	store := &Store{
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
func (st *Store) AutoMigrate() error {
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
func (st *Store) EnableDebug(debug bool) {
	st.debugEnabled = debug
}

// Log adds a log
func (st *Store) Log(logEntry *Log) (bool, error) {
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

	if err == nil {
		log.Println(sqlStr)
	}

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	_, err = st.db.Exec(sqlStr, sqlParams...)

	if err != nil {
		if st.debugEnabled {
			log.Println(err.Error())
		}
		return false, err
	}

	return true, nil
}

// Debug adds a debug log
func (st *Store) Debug(message string) (bool, error) {
	log := Log{
		Level:   LevelDebug,
		Message: message,
	}
	return st.Log(&log)
}

// DebugWithContext adds a debug log with context data
func (st *Store) DebugWithContext(message string, context interface{}) (bool, error) {
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
func (st *Store) Error(message string) (bool, error) {
	log := Log{
		Level:   LevelError,
		Message: message,
	}
	return st.Log(&log)
}

// ErrorWithContext adds an error log with context data
func (st *Store) ErrorWithContext(message string, context interface{}) (bool, error) {
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
func (st *Store) Fatal(message string) (bool, error) {
	log := Log{
		Level:   LevelFatal,
		Message: message,
	}

	result, err := st.Log(&log)
	// os.Exit(1)
	return result, err
}

// FatalWithContext adds a fatal log with context data and calls os.Exit(1) after logging
func (st *Store) FatalWithContext(message string, context interface{}) (bool, error) {
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

	result, err := st.Log(&log)
	// os.Exit(1)
	return result, err
}

// Info adds an info log
func (st *Store) Info(message string) (bool, error) {
	log := Log{
		Level:   LevelInfo,
		Message: message,
	}
	return st.Log(&log)
}

// InfoWithContext adds an info log with context data
func (st *Store) InfoWithContext(message string, context interface{}) (bool, error) {
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
func (st *Store) Panic(message string) {
	log := Log{
		Level:   LevelPanic,
		Message: message,
	}

	st.Log(&log)
	panic(message)
}

// PanicWithContext adds a panic log with context data and calls panic(message) after logging
func (st *Store) PanicWithContext(message string, context interface{}) {
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
func (st *Store) Trace(message string) (bool, error) {
	log := Log{
		Level:   LevelTrace,
		Message: message,
	}
	return st.Log(&log)
}

// TraceWithContext adds a trace log with context data
func (st *Store) TraceWithContext(message string, context interface{}) (bool, error) {
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
func (st *Store) Warn(message string) (bool, error) {
	log := Log{
		Level:   LevelWarning,
		Message: message,
	}
	return st.Log(&log)
}

// WarnWithContext adds a warn log with context data
func (st *Store) WarnWithContext(message string, context interface{}) (bool, error) {
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
