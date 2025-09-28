package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lmittmann/tint"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/nint8835/brscan-to-paperless/pkg/proto"
)

func checkErr(err error, msg string) {
	if err != nil {
		slog.Error(msg, "err", err)
		os.Exit(1)
	}
}

func initLogging() {
	level, validLevel := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}[strings.ToLower(logLevel)]
	if !validLevel {
		slog.Error("Invalid log level", "level", logLevel)
		os.Exit(1)
	}

	slog.SetDefault(slog.New(
		tint.NewHandler(
			os.Stderr,
			&tint.Options{
				TimeFormat: time.Kitchen,
				Level:      level,
			},
		),
	))
}

func createClient() (pb.BrscanToPaperlessClient, *grpc.ClientConn, error) {
	socketAbsPath, err := filepath.Abs(socketPath)
	checkErr(err, "Failed to get absolute path of socket")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get absolute path of socket: %w", err)
	}
	connStr := fmt.Sprintf("unix://%s", socketAbsPath)

	conn, err := grpc.NewClient(
		connStr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	return pb.NewBrscanToPaperlessClient(conn), conn, nil
}
