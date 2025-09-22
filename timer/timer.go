package timer

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
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
	bar := progressbar.NewOptions64(
		int64(duration.Seconds()),
		progressbar.OptionSetWidth(30),
		progressbar.OptionShowCount(),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetDescription(color.CyanString(
			"▶ %s [%02d:%02d]", label,
			int(duration.Minutes()), int(duration.Seconds())%60,
		)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerPadding: "░",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for remaining := duration; remaining > 0; {
		select {
		case <-ticker.C:
			if !pt.PauseFlag.Load() {
				remaining -= time.Second
				bar.Add(1)

				// Update label with time left (mm:ss), cyan text
				bar.Describe(color.CyanString("▶ %s [%02d:%02d]", label,
					int(remaining.Minutes()), int(remaining.Seconds())%60))
			}
		case cmd := <-pt.ControlChan:
			switch cmd {
			case "pause":
				pt.PauseFlag.Store(true)
				bar.Describe(color.YellowString("⏸ Paused - press 'r' to resume"))
			case "resume":
				pt.PauseFlag.Store(false)
				bar.Describe(color.GreenString("▶ Resumed: %s", label))
			case "quit":
				fmt.Println(color.RedString("\n⏹ Timer stopped early."))
				return false
			}
		}
	}

	// Completion feedback
	fmt.Println(color.GreenString("\n✅ %s completed!", label))

	// Desktop notification
	err := beeep.Notify(pt.WorkLabel, fmt.Sprintf("%s completed!", label), "")

	// Terminal beep (may or may not work depending on system)
	fmt.Print("\a")
	beeep.Beep(500, 200)

	if err != nil {
		fmt.Println(color.RedString("Error sending notification: %v", err))
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

func CheckInput(minutes int, breakMinutes int) bool {
	// if minutes, breakMinutes, or workSeconds is less than or equal to 0, return false
	if minutes <= 0 || breakMinutes <= 0 {
		fmt.Println("Error: Invalid input. Minutes, break, and seconds must be greater than 0.")
		return false
	}

	return true
}
