/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/Dima-salang/PomoLite/timer"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
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
			fmt.Println(color.RedString("❌ Error opening database: %v", err))
			return
		}
		defer storage.Close()

		pomoStats, err := storage.ComputePomoStats(timeframe)
		if err != nil {
			fmt.Println(color.RedString("❌ Error computing stats: %v", err))
			return
		}

		// Headline
		fmt.Println(color.CyanString("\n📊 Pomodoro Statistics (%s)\n", timeframe))

		// Tabwriter for aligned columns
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		// Summary stats
		fmt.Fprintf(w, "%s\t%s\n", color.YellowString("Total Sessions:"), fmt.Sprintf("%d", pomoStats.TotalSessions))
		fmt.Fprintf(w, "%s\t%s\n", color.YellowString("Total Work Duration:"), pomoStats.TotalWorkDuration.String())
		fmt.Fprintf(w, "%s\t%s\n", color.YellowString("Average Session:"), pomoStats.AverageSessionDuration.String())
		fmt.Fprintf(w, "%s\t%s\n", color.YellowString("Longest Session:"), pomoStats.LongestSession.String())
		fmt.Fprintf(w, "%s\t%s\n", color.YellowString("Shortest Session:"), pomoStats.ShortestSession.String())
		w.Flush()

		fmt.Println()

		// Sessions per label
		if len(pomoStats.PomosPerLabel) > 0 {
			fmt.Println(color.GreenString("📌 Sessions per Label:"))
			w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			for label, count := range pomoStats.PomosPerLabel {
				fmt.Fprintf(w, "  %s\t%d\n", color.MagentaString(label), count)
			}
			w.Flush()
			fmt.Println()
		}

		// Time spent per label
		if len(pomoStats.TimeSpentPerLabel) > 0 {
			fmt.Println(color.GreenString("⏱ Time Spent per Label:"))
			w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			for label, dur := range pomoStats.TimeSpentPerLabel {
				fmt.Fprintf(w, "  %s\t%s\n", color.MagentaString(label), formatDuration(dur))
			}
			w.Flush()
			fmt.Println()
		}
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

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%02dh %02dm %02ds", h, m, s)
	}
	return fmt.Sprintf("%02dm %02ds", m, s)
}