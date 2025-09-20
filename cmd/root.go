package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

var logLevel string

var rootCmd = &cobra.Command{
	Use:   "brscan-to-paperless",
	Short: "Integration between brscan / brscan-skey and Paperless-ngx.",
}

func Execute() {
	err := rootCmd.Execute()
	checkErr(err, "Failed to execute")
}

func init() {
	cobra.OnInitialize(initLogging)

	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Set the logging level (debug, info, warn, error)")
	_ = rootCmd.RegisterFlagCompletionFunc(
		"log-level",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			var matches []string
			for _, level := range []string{"debug", "info", "warn", "error"} {
				if strings.HasPrefix(level, strings.ToLower(toComplete)) {
					matches = append(matches, level)
				}
			}
			return matches, cobra.ShellCompDirectiveNoFileComp
		},
	)
}
