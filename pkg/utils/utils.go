package utils

import (
	"io"
	"log/slog"
)

// DeferredClose is a utility function to handle errors from deferred close functions.
func DeferredClose(obj io.Closer) {
	err := obj.Close()
	if err != nil {
		slog.Error("Failed to close resource", "err", err)
	}
}
