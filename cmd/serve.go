package cmd

import (
	"github.com/spf13/cobra"

	"github.com/nint8835/brscan-to-paperless/pkg/server"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the brscan-to-paperless daemon.",
	Run: func(cmd *cobra.Command, args []string) {
		server := server.New(socketPath)
		err := server.Serve()
		checkErr(err, "Failed to serve server")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
