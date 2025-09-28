package utils

import (
	"io"
	"log/slog"
)

// DeferCloser is a utility function to handle errors from deferred close functions.
func DeferCloser(obj io.Closer) {
	err := obj.Close()
	if err != nil {
		slog.Error("Failed to close resource", "err", err)
	}
}
