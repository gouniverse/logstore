package logstore

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlserver"
	"github.com/gouniverse/uid"
)

// Store defines a session store
type Store struct {
	logTableName       string
	db                 *sql.DB
	dbDriverName       string
	automigrateEnabled bool
	debug              bool
}

// StoreOption options for the cache store
type StoreOption func(*Store)

// WithAutoMigrate sets the table name for the cache store
func WithAutoMigrate(automigrateEnabled bool) StoreOption {
	return func(s *Store) {
		s.automigrateEnabled = automigrateEnabled
	}
}

// WithDb sets the database for the setting store
func WithDb(db *sql.DB) StoreOption {
	return func(s *Store) {
		s.db = db
		s.dbDriverName = s.DriverName(s.db)
	}
}

// WithDebug prints the SQL queries
func WithDebug(debug bool) StoreOption {
	return func(s *Store) {
		s.debug = debug
	}
}

// WithTableName sets the table name for the cache store
func WithTableName(logTableName string) StoreOption {
	return func(s *Store) {
		s.logTableName = logTableName
	}
}

// NewStore creates a new entity store
func NewStore(opts ...StoreOption) (*Store, error) {
	store := &Store{}

	for _, opt := range opts {
		opt(store)
	}

	if store.db == nil {
		return nil, errors.New("log store: db is required")
	}

	if store.dbDriverName == "" {
		return nil, errors.New("log store: dbDriverName is required")
	}

	if store.logTableName == "" {
		return nil, errors.New("log store: logTableName is required")
	}

	if store.automigrateEnabled {
		store.AutoMigrate()
	}

	return store, nil
}

// AutoMigrate auto migrate
func (st *Store) AutoMigrate() error {
	sql := st.SqlCreateTable()

	if st.debug {
		log.Println(sql)
	}

	_, err := st.db.Exec(sql)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// DriverName finds the driver name from database
func (st *Store) DriverName(db *sql.DB) string {
	dv := reflect.ValueOf(db.Driver())
	driverFullName := dv.Type().String()
	if strings.Contains(driverFullName, "mysql") {
		return "mysql"
	}
	if strings.Contains(driverFullName, "postgres") || strings.Contains(driverFullName, "pq") {
		return "postgres"
	}
	if strings.Contains(driverFullName, "sqlite") {
		return "sqlite"
	}
	if strings.Contains(driverFullName, "mssql") {
		return "mssql"
	}
	return driverFullName
}

// EnableDebug - enables the debug option
func (st *Store) EnableDebug(debug bool) {
	st.debug = debug
}

// Log adds a log
func (st *Store) Log(logEntry *Log) (bool, error) {
	if logEntry.ID == "" {
		logEntry.ID = uid.MicroUid()
	}
	if logEntry.Time == nil {
		now := time.Now()
		logEntry.Time = &now
	}

	var sqlStr string
	sqlStr, _, _ = goqu.Dialect(st.dbDriverName).Insert(st.logTableName).Rows(logEntry).ToSQL()

	if st.debug {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		if st.debug {
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

// SqlCreateTable returns a SQL string for creating the setting table
func (st *Store) SqlCreateTable() string {
	sqlMysql := `
	CREATE TABLE IF NOT EXISTS ` + st.logTableName + ` (
	  id varchar(40) NOT NULL PRIMARY KEY,
	  level varchar(40) NOT NULL,
	  message varchar(510) NOT NULL,
	  context longtext,
	  time datetime NOT NULL
	);
	`

	sqlPostgres := `
	CREATE TABLE IF NOT EXISTS "` + st.logTableName + `" (
	  "id" varchar(40) NOT NULL PRIMARY KEY,
	  "level" varchar(40) NOT NULL,
	  "message" varchar(510) NOT NULL,
	  "context" longtext,
	  "time" timestamptz(6) NOT NULL
	)
	`

	sqlSqlite := `
	CREATE TABLE IF NOT EXISTS "` + st.logTableName + `" (
	  "id" varchar(40) NOT NULL PRIMARY KEY,
	  "level" varchar(40) NOT NULL,
	  "message" varchar(510) NOT NULL,
	  "context" longtext,
	  "time" datetime NOT NULL
	)
	`

	sql := "unsupported driver '" + st.dbDriverName + "'"

	if st.dbDriverName == "mysql" {
		sql = sqlMysql
	}
	if st.dbDriverName == "postgres" {
		sql = sqlPostgres
	}
	if st.dbDriverName == "sqlite" {
		sql = sqlSqlite
	}

	return sql
}
