package timer

import (
	"database/sql"
	"fmt"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

type TimeFrame struct {
	start time.Time
	end time.Time
}


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
		INSERT INTO sessions (label, start_time, end_time)
		VALUES (?, ?, ?)
	`, label, startTime.Unix(), endTime.Unix())


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
		FROM sessions
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

// STATS
func (s *SQLiteStorage)ComputePomoStats(timeframe string) (*PomoStats, error) {
	statsTimeFrame, err := resolveTimeFrame(timeframe)
	if err != nil {
		return nil, err
	}

	stats := &PomoStats{}
	stats.TotalWorkDuration, _ = computeTotalWorkDurationStats(statsTimeFrame, s.db)
	stats.TotalSessions, _ = computeTotalSessions(statsTimeFrame, s.db)
	stats.AverageSessionDuration, _ = computeAverageSessionDuration(statsTimeFrame, s.db)
	stats.LongestSession, _ = computeLongestSession(statsTimeFrame, s.db)
	stats.ShortestSession, _ = computeShortestSession(statsTimeFrame, s.db)
	stats.HighestSessionLabel, _ = computeHighestSessionLabel(statsTimeFrame, s.db)
	stats.TimeSpentPerLabel, _ = computeTimeSpentPerLabel(statsTimeFrame, s.db)
	stats.PomosPerLabel, _ = computePomosPerLabel(statsTimeFrame, s.db)

	return stats, nil
}


// resolve time frame

func resolveTimeFrame(timeframe string) (TimeFrame, error) {
    now := time.Now()
    statsTimeFrame := TimeFrame{}

    switch timeframe {
    case "all":
        // Zero times (no filtering)
        statsTimeFrame.start = time.Time{}
        statsTimeFrame.end = time.Time{}

    case "today":
        statsTimeFrame.start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
        statsTimeFrame.end = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), now.Location())

    case "week":
        weekday := int(now.Weekday())
        if weekday == 0 { // Sunday → make it 7
            weekday = 7
        }
        startOfWeek := now.AddDate(0, 0, -(weekday-1)) // Monday
        statsTimeFrame.start = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, now.Location())
        statsTimeFrame.end = statsTimeFrame.start.AddDate(0, 0, 7).Add(-time.Nanosecond)

    case "month":
        statsTimeFrame.start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
        statsTimeFrame.end = statsTimeFrame.start.AddDate(0, 1, 0).Add(-time.Nanosecond)

    case "year":
        statsTimeFrame.start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
        statsTimeFrame.end = statsTimeFrame.start.AddDate(1, 0, 0).Add(-time.Nanosecond)

    default:
        return TimeFrame{}, fmt.Errorf("invalid timeframe: %s", timeframe)
    }

    return statsTimeFrame, nil
}



// compute total work duration stats
func computeTotalWorkDurationStats(timeframe TimeFrame, db *sql.DB) (time.Duration, error) {
	query := `
		SELECT SUM(end_time - start_time)
		FROM sessions
		WHERE start_time BETWEEN ? AND ?
	`

	var totalSeconds sql.NullInt64
	err := db.QueryRow(query, timeframe.start.Unix(), timeframe.end.Unix()).Scan(&totalSeconds)
	if err != nil {
		return 0, err
	}

	if !totalSeconds.Valid {
		return 0, nil
	}

	fmt.Println(totalSeconds.Int64)

	return time.Duration(totalSeconds.Int64) * time.Second, nil
}

func computeTotalSessions(timeframe TimeFrame, db *sql.DB) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM sessions
		WHERE start_time BETWEEN ? AND ?
	`
	var count int
	err := db.QueryRow(query, timeframe.start.Unix(), timeframe.end.Unix()).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func computeAverageSessionDuration(timeframe TimeFrame, db *sql.DB) (time.Duration, error) {
	query := `
		SELECT AVG(end_time - start_time)
		FROM sessions
		WHERE start_time BETWEEN ? AND ?
	`
	var avgSeconds sql.NullInt64
	err := db.QueryRow(query, timeframe.start.Unix(), timeframe.end.Unix()).Scan(&avgSeconds)
	if err != nil {
		return 0, err
	}
	if !avgSeconds.Valid {
		return 0, nil
	}
	fmt.Println("Average Seconds: ", avgSeconds.Int64)
	return time.Duration(avgSeconds.Int64) * time.Second, nil
}

func computeLongestSession(timeframe TimeFrame, db *sql.DB) (time.Duration, error) {
	query := `SELECT MAX(end_time - start_time)
		FROM sessions
		WHERE start_time BETWEEN ? AND ?`
	var longestSeconds sql.NullInt64
	err := db.QueryRow(query, timeframe.start.Unix(), timeframe.end.Unix()).Scan(&longestSeconds)
	if err != nil {
		return 0, err
	}
	if !longestSeconds.Valid {
		return 0, nil
	}
	fmt.Println("Longest Seconds: ", longestSeconds.Int64)
	return time.Duration(longestSeconds.Int64) * time.Second, nil
}

func computeShortestSession(timeframe TimeFrame, db *sql.DB) (time.Duration, error) {
	query := `SELECT MIN(end_time - start_time)
		FROM sessions
		WHERE start_time BETWEEN ? AND ?`
	var shortestSeconds sql.NullInt64
	err := db.QueryRow(query, timeframe.start.Unix(), timeframe.end.Unix()).Scan(&shortestSeconds)
	if err != nil {
		return 0, err
	}
	if !shortestSeconds.Valid {
		return 0, nil
	}
	fmt.Println("Shortest Seconds: ", shortestSeconds.Int64)
	return time.Duration(shortestSeconds.Int64) * time.Second, nil
}

func computeHighestSessionLabel(timeframe TimeFrame, db *sql.DB) (map[string]time.Duration, error) {
	query := `SELECT label, MAX(end_time - start_time)
		FROM sessions
		WHERE start_time BETWEEN ? AND ?`
	highestSessionLabel := make(map[string]time.Duration)
	err := db.QueryRow(query, timeframe.start.Unix(), timeframe.end.Unix()).Scan(&highestSessionLabel)
	if err != nil {
		return nil, err
	}
	fmt.Println("Highest Session Label: ", highestSessionLabel)
	return highestSessionLabel, nil
}
func computeTimeSpentPerLabel(timeframe TimeFrame, db *sql.DB) (map[string]time.Duration, error) {
	query := `SELECT label, SUM(end_time - start_time)
		FROM sessions
		WHERE start_time BETWEEN ? AND ?`
	timeSpentPerLabel := make(map[string]time.Duration)
	rows, err := db.Query(query, timeframe.start.Unix(), timeframe.end.Unix())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var label string
		var duration time.Duration
		if err := rows.Scan(&label, &duration); err != nil {
			return nil, err
		}
		timeSpentPerLabel[label] = duration
	}
	fmt.Println("Time Spent Per Label: ", timeSpentPerLabel)
	return timeSpentPerLabel, nil
}

func computePomosPerLabel(timeframe TimeFrame, db *sql.DB) (map[string]int, error) {
	query := `SELECT label, COUNT(*)
		FROM sessions
		WHERE start_time BETWEEN ? AND ?
		GROUP BY label`
	pomosPerLabel := make(map[string]int)
	rows, err := db.Query(query, timeframe.start.Unix(), timeframe.end.Unix())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var label string
		var count int
		if err := rows.Scan(&label, &count); err != nil {
			return nil, err
		}
		pomosPerLabel[label] = count
	}
	fmt.Println("Pomos Per Label: ", pomosPerLabel)
	return pomosPerLabel, nil
}


// INITIAL SETUP

// create table
func initTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			label TEXT NOT NULL,
			start_time INTEGER NOT NULL,
			end_time INTEGER NOT NULL
		)
	`)
	if err != nil {
		return err
	}
	return nil
}