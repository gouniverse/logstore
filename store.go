package logstore

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Store defines a session store
type Store struct {
	logTableName       string
	db                 *gorm.DB
	automigrateEnabled bool
}

// StoreOption options for the cache store
type StoreOption func(*Store)

// WithAutoMigrate sets the table name for the cache store
func WithAutoMigrate(automigrateEnabled bool) StoreOption {
	return func(s *Store) {
		s.automigrateEnabled = automigrateEnabled
	}
}

// WithDriverAndDNS sets the driver and the DNS for the database for the cache store
func WithDriverAndDNS(driverName string, dsn string) StoreOption {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	return func(s *Store) {
		s.db = db
	}
}

// WithGormDb sets the GORM database for the cache store
func WithGormDb(db *gorm.DB) StoreOption {
	return func(s *Store) {
		s.db = db
	}
}

// WithTableName sets the table name for the cache store
func WithTableName(logTableName string) StoreOption {
	return func(s *Store) {
		s.logTableName = logTableName
	}
}

// NewStore creates a new entity store
func NewStore(opts ...StoreOption) *Store {
	store := &Store{}
	for _, opt := range opts {
		opt(store)
	}

	if store.logTableName == "" {
		log.Panic("User store: cacheTableName is required")
	}

	if store.automigrateEnabled == true {
		store.AutoMigrate()
	}

	return store
}

// AutoMigrate auto migrate
func (st *Store) AutoMigrate() {
	st.db.Table(st.logTableName).AutoMigrate(&Log{})
}

// Log adds a log
func (st *Store) Log(log *Log) bool {
	if log.Time == nil {
		now := time.Now()
		log.Time = &now
	}

	result := st.db.Table(st.logTableName).Create(&log)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	}

	return true
}

// Debug adds a debug log
func (st *Store) Debug(message string) bool {
	log := Log{
		Level:   LevelDebug,
		Message: message,
	}
	return st.Log(&log)
}

// DebugWithContext adds a debug log with context data
func (st *Store) DebugWithContext(message string, context interface{}) bool {
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
func (st *Store) Error(message string) bool {
	log := Log{
		Level:   LevelError,
		Message: message,
	}
	return st.Log(&log)
}

// ErrorWithContext adds an error log with context data
func (st *Store) ErrorWithContext(message string, context interface{}) bool {
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
func (st *Store) Fatal(message string) bool {
	log := Log{
		Level:   LevelFatal,
		Message: message,
	}

	result := st.Log(&log)
	os.Exit(1)
	return result
}

// FatalWithContext adds a fatal log with context data and calls os.Exit(1) after logging
func (st *Store) FatalWithContext(message string, context interface{}) bool {
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

	result := st.Log(&log)
	os.Exit(1)
	return result
}

// Info adds an info log
func (st *Store) Info(message string) bool {
	log := Log{
		Level:   LevelInfo,
		Message: message,
	}
	return st.Log(&log)
}

// InfoWithContext adds an info log with context data
func (st *Store) InfoWithContext(message string, context interface{}) bool {
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
func (st *Store) Panic(message string) bool {
	log := Log{
		Level:   LevelPanic,
		Message: message,
	}

	result := st.Log(&log)
	panic(message)
	return result
}

// PanicWithContext adds a panic log with context data and calls panic(message) after logging
func (st *Store) PanicWithContext(message string, context interface{}) bool {
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

	result := st.Log(&log)
	panic(message)
	return result
}

// Trace adds a trace log
func (st *Store) Trace(message string) bool {
	log := Log{
		Level:   LevelTrace,
		Message: message,
	}
	return st.Log(&log)
}

// TraceWithContext adds a trace log with context data
func (st *Store) TraceWithContext(message string, context interface{}) bool {
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
func (st *Store) Warn(message string) bool {
	log := Log{
		Level:   LevelWarning,
		Message: message,
	}
	return st.Log(&log)
}

// WarnWithContext adds a warn log with context data
func (st *Store) WarnWithContext(message string, context interface{}) bool {
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
