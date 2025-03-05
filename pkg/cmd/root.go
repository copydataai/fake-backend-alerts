package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "alertcli",
	Short: "A CLI for testing alert management providers",
	Long: `A CLI tool designed to test and stress test alert management providers 
like OpsGenie, PagerDuty, etc. It can send individual alerts or run 
predefined scenarios to generate high volumes of alerts.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(scenarioCmd)
	rootCmd.AddCommand(versionCmd)
}
