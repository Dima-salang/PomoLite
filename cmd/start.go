package cmd

import (
	"time"
	"github.com/spf13/cobra"
	"github.com/Dima-salang/PomoLite/timer"
)

var minutes int
var breakMinutes int

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start a Pomodoro Timer.",
	Long: `Start a Pomodoro Timer.

	FLAGS:
	-m : minutes of work
	-b : minutes of break`,
	Run: func(cmd *cobra.Command, args []string) {
		pt := timer.NewPomodoroTimer(time.Duration(minutes)*time.Minute, time.Duration(breakMinutes)*time.Minute, "Work")
		go timer.ListenForCommands(pt.ControlChan)

		for {
			ok := pt.Start()
			if !ok {
				return
			}
		}
	},
}


func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().IntVarP(&minutes, "minutes", "m", 25, "minutes to work")
	startCmd.Flags().IntVarP(&breakMinutes, "break", "b", 5, "minutes to take a break")
}
