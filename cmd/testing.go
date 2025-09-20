package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/nint8835/brscan-to-paperless/pkg/proto"
)

var testingCmd = &cobra.Command{
	Use:   "testing",
	Short: "Command for testing purposes.",
	Run: func(cmd *cobra.Command, args []string) {
		socketAbsPath, err := filepath.Abs(socketPath)
		checkErr(err, "Failed to get absolute path of socket")

		connStr := fmt.Sprintf("unix://%s", socketAbsPath)

		conn, err := grpc.NewClient(
			connStr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		checkErr(err, "Failed to create gRPC client")
		defer func() {
			err := conn.Close()
			checkErr(err, "Failed to close gRPC connection")
		}()

		daemonClient := pb.NewBrscanToPaperlessClient(conn)
		_, err = daemonClient.TestRequest(cmd.Context(), &emptypb.Empty{})
		checkErr(err, "Failed to call TestRequest")
	},
}

func init() {
	rootCmd.AddCommand(testingCmd)
}
