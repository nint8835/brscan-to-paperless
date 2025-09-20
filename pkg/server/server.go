package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/nint8835/brscan-to-paperless/pkg/proto"
)

type Server struct {
	pb.UnimplementedBrscanToPaperlessServer

	socketPath string
	logger     *slog.Logger
}

func (s *Server) Serve() error {
	socketDir := filepath.Dir(s.socketPath)

	s.logger.Info("Starting server", "socketPath", s.socketPath)

	if _, err := os.Stat(socketDir); os.IsNotExist(err) {
		// TODO: Is this the right permission?
		if err := os.MkdirAll(socketDir, 0o755); err != nil {
			return err
		}
	}

	if _, err := os.Stat(s.socketPath); !os.IsNotExist(err) {
		if err := os.Remove(s.socketPath); err != nil {
			return fmt.Errorf("failed to remove existing socket file: %w", err)
		}
	}

	listener, err := net.Listen("unix", s.socketPath)
	if err != nil {
		return fmt.Errorf("failed to listen on socket: %w", err)
	}

	server := grpc.NewServer()
	pb.RegisterBrscanToPaperlessServer(server, s)

	return server.Serve(listener)
}

func (s *Server) TestRequest(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	s.logger.Info("Received TestRequest")
	return &emptypb.Empty{}, nil
}

func New(socketPath string) *Server {
	return &Server{
		socketPath: socketPath,
		logger: slog.Default().With(
			slog.String("component", "server"),
		),
	}
}
