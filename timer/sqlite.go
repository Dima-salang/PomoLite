package timer

import (
	"database/sql"
	"fmt"
	"time"
	_ "github.com/mattn/go-sqlite3"
)


type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(path string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := initTable(db); err != nil {
		return nil, err
	}
	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage)Close() error {
	return s.db.Close()
}

// save the timer data to the database
func (s *SQLiteStorage)SaveTimerData(label string, startTime time.Time, endTime time.Time) error {
	_, err := s.db.Exec(`
		INSERT INTO pomodoro (label, start_time, end_time)
		VALUES (?, ?, ?)
	`, label, startTime, endTime)


	fmt.Printf("Timer saved successfully with label: %s, start time: %s, end time: %s\n", label, startTime, endTime)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStorage)ListSessions(count int) ([]Session, error) {
	// if count is 0, return all sessions

	query := `
		SELECT id, label, start_time, end_time
		FROM pomodoro
		ORDER BY start_time DESC
	`
	var rows *sql.Rows
	var err error
	if count > 0 {
		query += " LIMIT ?"
		rows, err = s.db.Query(query, count)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
	} else {
		rows, err = s.db.Query(query)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
	}

	var sessions []Session
	for rows.Next() {
		var session Session
		if err := rows.Scan(&session.ID, &session.Label, &session.StartTime, &session.EndTime); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

// create table
func initTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS pomodoro (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			label TEXT NOT NULL,
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL
		)
	`)
	if err != nil {
		return err
	}
	return nil
}