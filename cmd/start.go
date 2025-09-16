package cmd

import (
	"fmt"
	"time"

	"github.com/Dima-salang/PomoLite/timer"
	"github.com/eiannone/keyboard"
	"github.com/spf13/cobra"
)

var minutes int
var breakMinutes int
var workSeconds int
var label string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start a Pomodoro Timer.",
	Long: `Start a Pomodoro Timer.

	FLAGS:
	-l : label for the work session
	-m : minutes of work
	-s : seconds of work
	-b : minutes of break`,
	Run: func(cmd *cobra.Command, args []string) {
		// check for the validity of the input
		if !timer.CheckInput(minutes, breakMinutes, workSeconds) {
			return
		}
		var storage timer.Storage
		storage, err := timer.NewSQLiteStorage("./pomodoro.db")

		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		defer storage.Close()

		totalWorkDuration := time.Duration(minutes)*time.Minute + time.Duration(workSeconds)*time.Second
		totalBreakDuration := time.Duration(breakMinutes) * time.Minute

		pt := timer.NewPomodoroTimer(totalWorkDuration, totalBreakDuration, label)
		go timer.ListenForCommands(pt.ControlChan)

		defer keyboard.Close()
		for {
			ok := pt.Start()
			if !ok {
				pt.EndTime = time.Now()
				storage.SaveTimerData(pt.WorkLabel, pt.StartTime, pt.EndTime)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&label, "label", "l", "Work", "label for the work session")
	startCmd.Flags().IntVarP(&minutes, "minutes", "m", 0, "minutes to work")
	startCmd.Flags().IntVarP(&workSeconds, "seconds", "s", 0, "seconds to work")
	startCmd.Flags().IntVarP(&breakMinutes, "break", "b", 5, "minutes to take a break")
}
