package worker

import "errors"

var ErrNoScannerFound = errors.New("no Brother scanner found")

var ErrTaskOngoing = errors.New("a task is already ongoing")
