package cmd

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Dima-salang/PomoLite/timer"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// sessionsCmd represents the sessions command
var sessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "list the sessions",
	Long: `List the sessions.

	This command lists all the sessions saved in the database.
	Each session includes the label, start time, and end time.
	The sessions are ordered by start time in descending order.`,
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		storage, err := timer.NewSQLiteStorage("./pomodoro.db")
		if err != nil {
			fmt.Println(color.RedString("Error: %v", err))
			return
		}
		defer storage.Close()

		sessions, err := storage.ListSessions(limit)
		if err != nil {
			fmt.Println(color.RedString("Error: %v", err))
			return
		}
		if len(sessions) == 0 {
			fmt.Println(color.YellowString("No sessions found."))
			return
		}

		// Regex to strip ANSI color codes for length calculations
		ansi := regexp.MustCompile("\x1b\\[[0-9;]*m")
		visibleLen := func(s string) int {
			return len([]rune(ansi.ReplaceAllString(s, "")))
		}
		padRightANSI := func(s string, width int) string {
			v := visibleLen(s)
			if v >= width {
				return s
			}
			return s + strings.Repeat(" ", width-v)
		}

		// Compute column widths (start/end use fixed format length 19)
		idW := visibleLen("ID")
		labelW := visibleLen("Label")
		startW := visibleLen("Start Time") // header, but we'll ensure at least 19
		endW := visibleLen("End Time")
		durW := visibleLen("Duration")

		// Iterate to figure out max widths
		for _, s := range sessions {
			idStr := fmt.Sprintf("%d", s.ID)
			if len(idStr) > idW {
				idW = len(idStr)
			}
			if visibleLen(s.Label) > labelW {
				labelW = visibleLen(s.Label)
			}
			startStr := s.StartTime.Format("2006-01-02 15:04:05")
			if visibleLen(startStr) > startW {
				startW = visibleLen(startStr)
			}
			endStr := s.EndTime.Format("2006-01-02 15:04:05")
			if visibleLen(endStr) > endW {
				endW = visibleLen(endStr)
			}
			durStr := s.EndTime.Sub(s.StartTime).Round(time.Second).String()
			if visibleLen(durStr) > durW {
				durW = visibleLen(durStr)
			}
		}

		// Header (colored)
		hID := color.CyanString("ID")
		hLabel := color.CyanString("Label")
		hStart := color.CyanString("Start Time")
		hEnd := color.CyanString("End Time")
		hDur := color.CyanString("Duration")

		// Print header and separator
		sepLen := idW + labelW + startW + endW + durW + 4*2 // 4 gaps of "  "
		fmt.Printf("%s  %s  %s  %s  %s\n",
			padRightANSI(hID, idW),
			padRightANSI(hLabel, labelW),
			padRightANSI(hStart, startW),
			padRightANSI(hEnd, endW),
			padRightANSI(hDur, durW),
		)
		fmt.Println(strings.Repeat("-", sepLen))

		// Rows (alternating label color)
		for i, s := range sessions {
			idStr := fmt.Sprintf("%d", s.ID)
			labelColored := color.GreenString(s.Label)
			if i%2 == 1 {
				labelColored = color.YellowString(s.Label)
			}
			startStr := s.StartTime.Format("2006-01-02 15:04:05")
			endStr := s.EndTime.Format("2006-01-02 15:04:05")
			durStr := s.EndTime.Sub(s.StartTime).Round(time.Second).String()
			durColored := color.MagentaString("%s", durStr)

			fmt.Printf("%s  %s  %s  %s  %s\n",
				padRightANSI(idStr, idW),
				padRightANSI(labelColored, labelW),
				padRightANSI(startStr, startW),
				padRightANSI(endStr, endW),
				padRightANSI(durColored, durW),
			)
		}

		// Footer note if --count was used
		if limit > 0 {
			fmt.Println()
			fmt.Println(color.HiBlackString("Showing last %d session(s).", limit))
		}
	},
}

func init() {
	rootCmd.AddCommand(sessionsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sessionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sessionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	sessionsCmd.Flags().IntP("limit", "l", 0, "number of sessions to list")
}
