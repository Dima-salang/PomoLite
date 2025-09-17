/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/Dima-salang/PomoLite/timer"
	"github.com/fatih/color"
)

/*
TIMEFRAME possible values:
- all
- today
- week
- month
- year
*/


// statCmd represents the stat command
var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "Displays the statistics for the timeframe specified",
	Long: `Displays the statistics for the timeframe specified. Possible values for timeframe are:
- all
- today
- week
- month
- year`,
	Run: func(cmd *cobra.Command, args []string) {
		timeframe, _ := cmd.Flags().GetString("timeframe")
		storage, err := timer.NewSQLiteStorage("./pomodoro.db")
		if err != nil {
			fmt.Println(color.RedString("Error: %v", err))
			return
		}
		defer storage.Close()
		pomoStats, err := storage.ComputePomoStats(timeframe)
		if err != nil {
			fmt.Println(color.RedString("Error: %v", err))
			return
		}
		fmt.Println(pomoStats)
	},
}

func init() {
	rootCmd.AddCommand(statCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	statCmd.Flags().StringP("timeframe", "t", "all", "timeframe for stats")
}
