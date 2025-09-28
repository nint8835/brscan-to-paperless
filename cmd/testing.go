package cmd

import (
	"context"

	"github.com/spf13/cobra"

	pb "github.com/nint8835/brscan-to-paperless/pkg/proto"
	"github.com/nint8835/brscan-to-paperless/pkg/utils"
)

var testingCmd = &cobra.Command{
	Use:    "testing",
	Short:  "Command for testing purposes.",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		client, conn, err := createClient()
		checkErr(err, "Failed to create gRPC client")
		defer utils.DeferCloser(conn)

		_, err = client.Trigger(context.Background(), &pb.TriggerRequest{
			Option: pb.TriggerOption_TRIGGER_OPTION_FILE,
		})
		checkErr(err, "Failed to call TestRequest")
	},
}

func init() {
	rootCmd.AddCommand(testingCmd)
}
