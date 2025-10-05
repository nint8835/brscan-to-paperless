package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/nint8835/brscan-to-paperless/pkg/proto"
	"github.com/nint8835/brscan-to-paperless/pkg/worker"
)

type Server struct {
	pb.UnimplementedBrscanToPaperlessServer

	socketPath string
	logger     *slog.Logger
	worker     *worker.Worker
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

func (s *Server) Trigger(ctx context.Context, req *pb.TriggerRequest) (*pb.TriggerResponse, error) {
	s.logger.Info("Trigger called", "option", req.Option.String())

	scannedPages, err := s.worker.Scan()
	if errors.Is(err, worker.ErrTaskOngoing) {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	} else if err != nil {
		s.logger.Error("Failed to scan", "err", err)
		return nil, status.Error(codes.Internal, "Failed to scan")
	}

	return &pb.TriggerResponse{
		PagesScanned: uint32(scannedPages),
	}, nil
}

func New(socketPath string) (*Server, error) {
	workerInst, err := worker.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create worker: %w", err)
	}

	return &Server{
		socketPath: socketPath,
		logger: slog.Default().With(
			slog.String("component", "server"),
		),
		worker: workerInst,
	}, nil
}
