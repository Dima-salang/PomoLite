package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/Dima-salang/PomoLite/timer"

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
		count, _ := cmd.Flags().GetInt("count")
		storage, err := timer.NewSQLiteStorage("./pomodoro.db")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		defer storage.Close()
		sessions, err := storage.ListSessions(count)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tLabel\tStart Time\tEnd Time\tDuration")
		for _, s := range sessions {
			duration := s.EndTime.Sub(s.StartTime).Round(time.Second)
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
				s.ID,
				s.Label,
				s.StartTime.Format("2006-01-02 15:04:05"),
				s.EndTime.Format("2006-01-02 15:04:05"),
				duration,
			)
		}
		w.Flush()
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
