package tests

import (
	"os"
	"testing"
	"time"

	"github.com/Dima-salang/PomoLite/timer"
)

var testStorage *timer.SQLiteStorage

func newTestSQLiteStorage(t *testing.T) *timer.SQLiteStorage {
	t.Helper()
	storage, err := timer.NewSQLiteStorage(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	return storage
}

func TestMain(m *testing.M) {
	storage, err := timer.NewSQLiteStorage(":memory:")
	if err != nil {
		panic(err)
	}
	testStorage = storage
	defer storage.Close()
	
	code := m.Run()

	storage.Close()

	os.Exit(code)
}


func TestSaveTimerData(t *testing.T) {
	storage := testStorage
	startTime := time.Date(2025, 9, 17, 12, 0, 0, 0, time.Local)
	endTime := time.Date(2025, 9, 17, 12, 15, 0, 0, time.Local)

	err := storage.SaveTimerData("Test", startTime, endTime)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListSessions(t *testing.T) {
	storage := testStorage

	startTime := time.Date(2025, 9, 17, 12, 0, 0, 0, time.Local)
	endTime := time.Date(2025, 9, 17, 12, 15, 0, 0, time.Local)

	err := storage.SaveTimerData("Test", startTime, endTime)
	if err != nil {
		t.Fatal(err)
	}

	sessions, err := storage.ListSessions(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(sessions) != 1 {
		t.Fatalf("Expected 1 session, got %d", len(sessions))
	}
}