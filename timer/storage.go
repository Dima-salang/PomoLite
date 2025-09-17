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

type PomoStats struct {
	TotalWorkDuration time.Duration
	TotalSessions int
	AverageSessionDuration time.Duration
	LongestSession time.Duration
	ShortestSession time.Duration
	HighestSessionLabel map[string]time.Duration
	TimeSpentPerLabel map[string]time.Duration
	PomosPerLabel map[string]int
}
