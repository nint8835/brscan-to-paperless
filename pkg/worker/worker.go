package worker

import (
	"fmt"
	"log/slog"
	"runtime"
	"sync"

	"github.com/fewebahr/sane"
)

type Worker struct {
	mutex  sync.Mutex
	logger *slog.Logger
	conn   *sane.Conn
}

func New() (*Worker, error) {
	err := sane.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SANE: %w", err)
	}

	worker := &Worker{
		mutex: sync.Mutex{},
		logger: slog.Default().With(
			slog.String("component", "worker"),
		),
	}
	runtime.AddCleanup(worker, func(_ any) {
		sane.Exit()
	}, "")

	devices, err := sane.Devices()
	if err != nil {
		return nil, fmt.Errorf("failed to list SANE devices: %w", err)
	}

	var brotherDevice *sane.Device

	for _, dev := range devices {
		if dev.Vendor == "Brother" {
			worker.logger.Debug("Found Brother device", "name", dev.Name, "model", dev.Model)
			brotherDevice = &dev
			break
		}
	}

	if brotherDevice == nil {
		return nil, ErrNoScannerFound
	}

	conn, err := sane.Open(brotherDevice.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to open SANE connection: %w", err)
	}
	worker.conn = conn

	return worker, nil
}
