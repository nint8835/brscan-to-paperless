package cmd

import (
	"context"
	"log/slog"

	"github.com/spf13/cobra"

	pb "github.com/nint8835/brscan-to-paperless/pkg/proto"
	"github.com/nint8835/brscan-to-paperless/pkg/utils"
)

var triggerCmd = &cobra.Command{
	Use:   "trigger <option>",
	Short: "Trigger an action in the brscan-to-paperless service.",
	Long: "Trigger an action in the brscan-to-paperless service. This is intended to be used by the brscan-skey " +
		"daemon, but can be used to manually trigger actions without having to touch your printer.",
	Args: cobra.ExactArgs(1),
	ValidArgs: []string{
		"file",
		"ocr",
		"image",
		"email",
	},
	Run: func(cmd *cobra.Command, args []string) {
		triggerOpt := map[string]pb.TriggerOption{
			"file":  pb.TriggerOption_TRIGGER_OPTION_FILE,
			"ocr":   pb.TriggerOption_TRIGGER_OPTION_OCR,
			"image": pb.TriggerOption_TRIGGER_OPTION_IMAGE,
			"email": pb.TriggerOption_TRIGGER_OPTION_EMAIL,
		}[args[0]]

		client, conn, err := createClient()
		checkErr(err, "Failed to create gRPC client")
		defer utils.DeferredClose(conn)

		resp, err := client.Trigger(context.Background(), &pb.TriggerRequest{
			Option: triggerOpt,
		})
		checkErr(err, "Failed to trigger action")

		slog.Info("Successfully triggered action", "option", args[0], "response", resp)
	},
}

func init() {
	rootCmd.AddCommand(triggerCmd)
}
