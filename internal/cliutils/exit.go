package cliutils

import (
	"os"

	"github.com/glasskube/glasskube/internal/telemetry"
)

// ExitSuccess shuts down with exit code 0
func ExitSuccess() {
	telemetry.Exit()
	os.Exit(0)
}

// ExitWithError ends the process with exit code 1
func ExitWithError() {
	telemetry.ExitWithError()
	os.Exit(1)
}

// ExitFromSignal ends the process with exit code 1
func ExitFromSignal(sig *os.Signal) {
	if sig == nil {
		sig = &os.Interrupt
	}
	telemetry.ExitFromSignal(*sig)
	os.Exit(1)
}
