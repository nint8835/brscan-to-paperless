package worker

import (
	"errors"
	"fmt"

	"github.com/fewebahr/sane"
)

func (w *Worker) Scan() (int, error) {
	locked := w.mutex.TryLock()
	if !locked {
		return 0, ErrTaskOngoing
	}
	defer w.mutex.Unlock()

	scannedImages := []*sane.Image{}

	for {
		w.logger.Info("Requesting image from scanner...")
		img, err := w.conn.ReadImage()
		if errors.Is(err, sane.ErrEmpty) {
			w.logger.Info("No more images to scan.")
			break
		} else if err != nil {
			return 0, fmt.Errorf("failed to read image from scanner: %w", err)
		}

		w.logger.Info("Scanned image", "bounds", img.Bounds())
		scannedImages = append(scannedImages, img)
	}

	w.logger.Info("Completed scanning", "numImages", len(scannedImages))
	return len(scannedImages), nil
}
