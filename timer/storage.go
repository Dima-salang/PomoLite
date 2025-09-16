package timer

// Storage interface for the timer using sqlite

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	SaveTimerData(label string, startTime time.Time, endTime time.Time) error
	ListSessions(count int) ([]Session, error)
	Close() error
}

type Session struct {
	ID        int
	Label     string
	StartTime time.Time
	EndTime   time.Time
}
