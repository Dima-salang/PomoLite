package cmd

import (
	"fmt"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var minutes int
var breakMinutes int
var workDuration time.Duration
var breakDuration time.Duration

var controlChan chan bool
var pauseFlag bool

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start a Pomodoro Timer.",
	Long: `Start a Pomodoro Timer.

	FLAGS:
	-m : minutes of work
	-b : minutes of break`,
	Run: func(cmd *cobra.Command, args []string) {
		controlChan = make(chan bool)
		defer close(controlChan)

		for {
			workDuration = time.Duration(minutes) * time.Minute
			breakDuration = time.Duration(breakMinutes) * time.Minute

			fmt.Printf("Starting Pomodoro Timer for %s with a break of %s\n", workDuration.String(), breakDuration.String())

			workComplete := make(chan bool)
			breakComplete := make(chan bool)

			go func() {
				countDownStart(workDuration, workComplete)
				fmt.Println("\nWork completed, good job! Take a break.")
				close(workComplete)
			}()

			select {
			case <-workComplete: // Work completed, move to break
			}

			go func() {
				countDownStart(breakDuration, breakComplete)
				fmt.Println("\nBreak completed. Back to work.")
				close(breakComplete)
			}()

			select {
			case <-breakComplete: // Break completed, repeat the cycle
			}
		}
	},
}

func countDownStart(duration time.Duration, completionChan chan bool) {
	bar := progressbar.Default(int64(duration.Seconds()))
	ticker := time.NewTicker(time.Second)
	pauseFlag = false
	defer ticker.Stop()

	for remaining := duration; remaining > 0; {
		select {
		case <-ticker.C:
			if !pauseFlag {
				remaining -= time.Second
				bar.Add(1)
			}
		case <-controlChan:
			pauseTimer()
			fmt.Println("\nTimer paused. Press Enter to resume.")
			waitForEnter()
			resumeTimer()
		}
	}

	err := beeep.Notify("Pomodoro Timer", "Timer completed.", "")
	if err != nil {
		panic(err)
	}



	completionChan <- true // Notify that the countdown is complete
}

func waitForEnter() {
	// Wait for Enter without blocking
	go func() {
		var input string
		fmt.Scanln(&input)
		controlChan <- true // Notify that Enter is pressed
	}()
}

func pauseTimer() {
	pauseFlag = true
}

func resumeTimer() {
	pauseFlag = false
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().IntVarP(&minutes, "minutes", "m", 25, "minutes to work")
	startCmd.Flags().IntVarP(&breakMinutes, "break", "b", 5, "minutes to take a break")
}
