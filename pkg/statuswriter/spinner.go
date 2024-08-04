package statuswriter

import (
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

type callback []func()

func (cb *callback) Add(fn ...func()) {
	*cb = append(([]func())(*cb), fn...)
}

func (cb *callback) Invoke() {
	for _, fn := range *cb {
		fn()
	}
	*cb = callback{}
}

type spinnerStatusWriter struct {
	barSupplier func() *progressbar.ProgressBar
	bar         *progressbar.ProgressBar
	onStop      callback
}

// SetStatus implements StatusWriter.
func (obj *spinnerStatusWriter) SetStatus(desc string) {
	obj.bar.Describe(desc)
}

// Start implements StatusWriter.
func (obj *spinnerStatusWriter) Start() {
	obj.bar = obj.barSupplier()
	ticker := time.NewTicker(100 * time.Millisecond)
	obj.onStop.Add(
		func() {
			ticker.Stop()
			_ = obj.bar.Clear()
			_ = obj.bar.Finish()
			_ = obj.bar.Exit()
		})
	go func() {
		for range ticker.C {
			_ = obj.bar.Add(1)
		}
	}()
}

// Stop implements StatusWriter.
func (obj *spinnerStatusWriter) Stop() {
	obj.onStop.Invoke()
	obj.bar = nil
}

func (obj *spinnerStatusWriter) WithStatusbar(supplier func() *progressbar.ProgressBar) *spinnerStatusWriter {
	obj.barSupplier = supplier
	return obj
}

func createDefaultProgressbar() *progressbar.ProgressBar {
	return progressbar.NewOptions64(
		-1,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(10),
		progressbar.OptionSpinnerType(11),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionClearOnFinish(),
	)
}

func Spinner() *spinnerStatusWriter {
	return &spinnerStatusWriter{barSupplier: createDefaultProgressbar}
}
