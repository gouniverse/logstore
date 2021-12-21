package logstore

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestWithAutoMigrate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// TC: 1

	// Initializes automigrateEnabled to False
	s := Store{
		logTableName:       "LogTable",
		db:                 db,
		automigrateEnabled: false,
	}

	// Modified to True
	f := WithAutoMigrate(true)
	f(&s)

	// Test Results
	if s.automigrateEnabled != true {
		t.Fatalf("automigrateEnabled: Expected [true] received [%v]", s.automigrateEnabled)
	}

	// TC: 2

	// Initializes automigrateEnabled to True
	s = Store{
		logTableName:       "LogTable",
		db:                 db,
		automigrateEnabled: true,
	}

	// Modified to True
	f = WithAutoMigrate(false)
	f(&s)

	// Test Results
	if s.automigrateEnabled == true {
		t.Fatalf("automigrateEnabled: Expected [true] received [%v]", s.automigrateEnabled)
	}
}

func TestWithDriverAndDNS(t *testing.T) {
	s := Store{
		logTableName:       "LogTable",
		db:                 nil,
		automigrateEnabled: false,
	}

	// db is initialized to nil
	f := WithDriverAndDNS("test.sqlite", "test.sqlite")
	// DB has to be initialized now
	f(&s)

	// db non Nil expected
	if s.db == nil {
		t.Fatalf("db initialization failed")
	}
}

func TestWithGormDb(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// db is initialized to nil
	s := Store{
		logTableName:       "LogTable",
		db:                 nil,
		automigrateEnabled: false,
	}

	f := WithGormDb(db)
	// DB has to be initialized now
	f(&s)

	// db non Nil expected
	if s.db == nil {
		t.Fatalf("db initialization failed")
	}
}

func TestWithTableName(t *testing.T) {
	s := Store{
		logTableName:       "",
		db:                 nil,
		automigrateEnabled: false,
	}
	// TC: 1
	table_name := "Table1"
	f := WithTableName(table_name)
	f(&s)
	if s.logTableName != table_name {
		t.Fatalf("Expected logTableName [%v], received [%v]", table_name, s.logTableName)
	}
	// TC: 2
	table_name = "Table2"
	f = WithTableName(table_name)
	f(&s)
	if s.logTableName != table_name {
		t.Fatalf("Expected logTableName [%v], received [%v]", table_name, s.logTableName)
	}
}

func Test_Store_AutoMigrate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// TC: 1

	// Initializes automigrateEnabled to False
	s := Store{
		logTableName:       "LogTable",
		db:                 db,
		automigrateEnabled: true,
	}

	s.AutoMigrate()
}

func Test_Store_Log(t *testing.T) {
	log := Log{
		ID:      "test.sqlite",
		Level:   "test.sqlite",
		Message: "test.sqlite",
		Context: "test.sqlite",
		Time:    nil,
	}

	db, err := gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// TC: 1

	// Initializes automigrateEnabled to False
	s := Store{
		logTableName:       "LogTable",
		db:                 db,
		automigrateEnabled: true,
	}

	b := s.Log(&log)
	if b != true {
		t.Fatalf("Expected [true] received [%v]", b)
	}
}

func Test_Store_Debug(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// TC: 1

	// Initializes automigrateEnabled to False
	s := Store{
		logTableName:       "LogTable",
		db:                 db,
		automigrateEnabled: true,
	}

	b := s.Debug("debug")
	if b != true {
		t.Fatalf("Expected [true] received [%v]", b)
	}
}

func Test_Store_DebugWithContext(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// TC: 1

	// Initializes automigrateEnabled to False
	s := Store{
		logTableName:       "LogTable",
		db:                 db,
		automigrateEnabled: true,
	}

	b := s.DebugWithContext("debug", "Debug Message")
	if b != true {
		t.Fatalf("Expected [true] received [%v]", b)
	}
}

func Test_Store_Error(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// TC: 1

	// Initializes automigrateEnabled to False
	s := Store{
		logTableName:       "LogTable",
		db:                 db,
		automigrateEnabled: true,
	}

	b := s.Error("error")
	if b != true {
		t.Fatalf("Expected [true] received [%v]", b)
	}
}

func Test_Store_ErrorWithContext(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// TC: 1

	// Initializes automigrateEnabled to False
	s := Store{
		logTableName:       "LogTable",
		db:                 db,
		automigrateEnabled: true,
	}

	b := s.ErrorWithContext("error", "Error Message")
	if b != true {
		t.Fatalf("Expected [true] received [%v]", b)
	}
}

// Fatal methods uses system level API to terminate program (os.Exit)
/*
func Test_Store_Fatal(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// TC: 1

	// Initializes automigrateEnabled to False
	s := Store{
		logTableName:       "LogTable",
		db:                 db,
		automigrateEnabled: true,
	}

	b := s.Fatal("fatal")
	if b != true {
		t.Fatalf("Expected [true] received [%v]", b)
	}
}

func Test_Store_FatalWithContext(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// TC: 1

	// Initializes automigrateEnabled to False
	s := Store{
		logTableName:       "LogTable",
		db:                 db,
		automigrateEnabled: true,
	}

	b := s.FatalWithContext("fatal", "Fatal Message")
	if b != true {
		t.Fatalf("Expected [true] received [%v]", b)
	}
}
*/

func Test_Store_Info(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// TC: 1

	// Initializes automigrateEnabled to False
	s := Store{
		logTableName:       "LogTable",
		db:                 db,
		automigrateEnabled: true,
	}

	b := s.Info("Info")
	if b != true {
		t.Fatalf("Expected [true] received [%v]", b)
	}
}

func Test_Store_InfoWithContext(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// TC: 1

	// Initializes automigrateEnabled to False
	s := Store{
		logTableName:       "LogTable",
		db:                 db,
		automigrateEnabled: true,
	}

	b := s.InfoWithContext("Info", "Info Message")
	if b != true {
		t.Fatalf("Expected [true] received [%v]", b)
	}
}

func Test_Store_Trace(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// TC: 1

	// Initializes automigrateEnabled to False
	s := Store{
		logTableName:       "LogTable",
		db:                 db,
		automigrateEnabled: true,
	}

	b := s.Trace("trace")
	if b != true {
		t.Fatalf("Expected [true] received [%v]", b)
	}
}

func Test_Store_TraceWithContext(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// TC: 1

	// Initializes automigrateEnabled to False
	s := Store{
		logTableName:       "LogTable",
		db:                 db,
		automigrateEnabled: true,
	}

	b := s.TraceWithContext("trace", "Trace Message")
	if b != true {
		t.Fatalf("Expected [true] received [%v]", b)
	}
}

func Test_Store_Warn(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// TC: 1

	// Initializes automigrateEnabled to False
	s := Store{
		logTableName:       "LogTable",
		db:                 db,
		automigrateEnabled: true,
	}

	b := s.Warn("warn")
	if b != true {
		t.Fatalf("Expected [true] received [%v]", b)
	}
}

func Test_Store_WarnWithContext(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// TC: 1

	// Initializes automigrateEnabled to False
	s := Store{
		logTableName:       "LogTable",
		db:                 db,
		automigrateEnabled: true,
	}

	b := s.WarnWithContext("warn", "Warning Message")
	if b != true {
		t.Fatalf("Expected [true] received [%v]", b)
	}
}
