/*
Copyright © 2023 LUIS GABRIELLE PUTAN <luisgabrielle1026@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const asciiArt = `

██████╗  ██████╗ ███╗   ███╗ ██████╗ ██╗     ██╗████████╗███████╗
██╔══██╗██╔═══██╗████╗ ████║██╔═══██╗██║     ██║╚══██╔══╝██╔════╝
██████╔╝██║   ██║██╔████╔██║██║   ██║██║     ██║   ██║   █████╗  
██╔═══╝ ██║   ██║██║╚██╔╝██║██║   ██║██║     ██║   ██║   ██╔══╝  
██║     ╚██████╔╝██║ ╚═╝ ██║╚██████╔╝███████╗██║   ██║   ███████╗
╚═╝      ╚═════╝ ╚═╝     ╚═╝ ╚═════╝ ╚══════╝╚═╝   ╚═╝   ╚══════╝
                                                                 


`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pomo",
	Short: "Lightweight CLI Pomodoro Timer",
	Long: fmt.Sprintf(`%s PomoLite is a lightweight CLI Pomodoro application desgned for students to efficiently accomplish tasks and maximize their learning potential
########################################################################################

Example usage:
	
pomo start -m 30 -b 5 :: Starts a 30-minute Pomodoro timer with a 5-minute break

To pause the Pomodoro Timer, you can press the 'p' button on your keyboard.

To unpause or resume the Pomodoro Timer, press 'r'.

To quit the Pomodoro Timer and save it for stats, press 'q'.



Developed by PUTAN LUIS GABRIELLE <luisgabrielle1026@gmail.com>

#######################################################################################`, asciiArt),
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.PomoLite.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


