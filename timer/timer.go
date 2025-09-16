package timer

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/gen2brain/beeep"
	"github.com/schollz/progressbar/v3"
)

type PomodoroTimer struct {
	WorkDuration  time.Duration
	BreakDuration time.Duration
	WorkLabel     string
	PauseFlag     atomic.Bool
	ControlChan   chan string
	StartTime     time.Time
	EndTime       time.Time
}

func NewPomodoroTimer(workDuration time.Duration, breakDuration time.Duration, workLabel string) *PomodoroTimer {
	return &PomodoroTimer{
		WorkDuration:  workDuration,
		BreakDuration: breakDuration,
		WorkLabel:     workLabel,
		PauseFlag:     atomic.Bool{},
		ControlChan:   make(chan string),
		StartTime:     time.Now(),
		EndTime:       time.Now(),
	}
}

func (pt *PomodoroTimer) Start() bool {
	fmt.Printf("Starting Pomodoro Timer for %s for %s with a break of %s\n", pt.WorkLabel, pt.WorkDuration.String(), pt.BreakDuration.String())

	if !pt.CountDownStart(pt.WorkLabel, pt.WorkDuration) {
		return false
	}

	fmt.Println("\nWork completed, good job! Take a break.")

	if !pt.CountDownStart(pt.WorkLabel, pt.BreakDuration) {
		return false
	}

	fmt.Println("\nBreak completed. Back to work.")

	return true
}

func (pt *PomodoroTimer) CountDownStart(label string, duration time.Duration) bool {
	bar := progressbar.Default(int64(duration.Seconds()))
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for remaining := duration; remaining > 0; {
		select {
		case <-ticker.C:
			if !pt.PauseFlag.Load() {
				remaining -= time.Second
				bar.Add(1)
			}
		case cmd := <-pt.ControlChan:
			switch cmd {
			case "pause":
				pt.PauseFlag.Store(true)
				fmt.Println("\nTimer paused. Press r to resume.")
			case "resume":
				pt.PauseFlag.Store(false)
			case "quit":
				return false
			}

		}
	}

	err := beeep.Notify(pt.WorkLabel, "Timer completed.", "")
	if err != nil {
		fmt.Println("Error: ", err)
	}

	return true
}

// listening for commands
func ListenForCommands(controlChan chan<- string) {
	if err := keyboard.Open(); err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer keyboard.Close()
	for {
		cmd_input, _, err := keyboard.GetKey()
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		switch cmd_input {
		case 'p':
			controlChan <- "pause"
		case 'r':
			controlChan <- "resume"
		case 'q':
			controlChan <- "quit"
			return
		}
	}
}

func CheckInput(minutes int, breakMinutes int, workSeconds int) bool {
	if minutes < 0 || breakMinutes < 0 || workSeconds < 0 {
		fmt.Println("Error: Invalid input. Minutes, break, and seconds must be non-negative.")
		return false
	}

	if workSeconds >= 60 {
		fmt.Println("Error: Invalid input. Seconds must be less than 60.")
		return false
	}

	return true
}
