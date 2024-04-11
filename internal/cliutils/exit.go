package cliutils

import (
	"context"
	"os"

	"github.com/glasskube/glasskube/internal/telemetry"
)

// ExitSuccess shuts down with exit code 0
func ExitSuccess(ctx context.Context) {
	telemetry.Exit(ctx)
	os.Exit(0)
}

// ExitWithError ends the process with exit code 1
func ExitWithError(ctx context.Context) {
	telemetry.ExitWithError(ctx)
	os.Exit(1)
}

// ExitFromSignal ends the process with exit code 1
func ExitFromSignal(ctx context.Context, sig *os.Signal) {
	telemetry.ExitFromSignal(ctx, *sig)
	os.Exit(1)
}
