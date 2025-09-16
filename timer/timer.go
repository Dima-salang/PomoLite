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
	WorkDuration time.Duration
	BreakDuration time.Duration
	WorkLabel string
	PauseFlag atomic.Bool
	ControlChan chan string
}

func NewPomodoroTimer(workDuration time.Duration, breakDuration time.Duration, workLabel string) *PomodoroTimer {
	return &PomodoroTimer{
		WorkDuration: workDuration,
		BreakDuration: breakDuration,
		WorkLabel: workLabel,
		PauseFlag: atomic.Bool{},
		ControlChan: make(chan string),
	}
}


func (pt *PomodoroTimer) Start() bool {
	fmt.Printf("Starting Pomodoro Timer for %s with a break of %s\n", pt.WorkDuration.String(), pt.BreakDuration.String())

	if (!pt.CountDownStart("Work", pt.WorkDuration)) {
		return false
	}

	fmt.Println("\nWork completed, good job! Take a break.")

	if (!pt.CountDownStart("Break", pt.BreakDuration)) {
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

	err := beeep.Notify("Pomodoro Timer", "Timer completed.", "")
	if err != nil {
		fmt.Println("Error: ", err)
	}



	return true
}


// listening for commands
func ListenForCommands(controlChan chan <- string) {
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
			}
	}
}
	
