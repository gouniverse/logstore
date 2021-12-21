package logstore

import (
	"time"

	"github.com/gouniverse/uid"
	"gorm.io/gorm"
)

const (
	// LevelTrace trace level
	LevelTrace = "trace"
	// LevelDebug debug level
	LevelDebug = "debug"
	// LevelError error level
	LevelError = "error"
	// LevelFatal fatal level
	LevelFatal = "fatal"
	// LevelInfo info level
	LevelInfo = "info"
	// LevelPanic panic level
	LevelPanic = "panic"
	// LevelWarning warning level
	LevelWarning = "warning"
)

// Log type
type Log struct {
	ID      string     `gorm:"type:varchar(40);column:id;primary_key;"`
	Level   string     `gorm:"type:varchar(10);column:level;"`
	Message string     `gorm:"type:varchar(510);column:message;"`
	Context string     `gorm:"type:longtext;column:context;"`
	Time    *time.Time `gorm:"type:datetime;column:time;DEFAULT NULL;"`
}

// BeforeCreate adds UID to model
func (l *Log) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uid.HumanUid()
	l.ID = uuid
	return nil
}
