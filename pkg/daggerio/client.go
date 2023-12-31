package daggerio

import (
	"context"
	"io"
	"os"

	"dagger.io/dagger"
)

var defaultLogOutput io.Writer = os.Stdout // Default log output is stdout

type BackendClient interface {
	CreateDaggerBackend(ctx context.Context, options ...dagger.ClientOpt) (*dagger.Client, error)
	ResolveDaggerLogConfig(enableErrorsOnly bool) io.Writer
}

type Backend struct{}

// CreateDaggerBackend creates a new dagger client.
// If no options are passed, the default options are used.
func (b *Backend) CreateDaggerBackend(ctx context.Context, options ...dagger.ClientOpt) (*dagger.
	Client, error) {
	var c context.Context
	if ctx == nil {
		c = context.Background()
	} else {
		c = ctx
	}

	var daggerOptions []dagger.ClientOpt

	if len(options) == 0 {
		return dagger.Connect(c, dagger.WithLogOutput(os.Stderr))
	}

	// If options are passed, append them to the daggerOptions.
	daggerOptions = append(daggerOptions, options...)
	return dagger.Connect(c, daggerOptions...)
}

// ResolveDaggerLogConfig returns the log output for the dagger client.
func (b *Backend) ResolveDaggerLogConfig(enableErrorsOnly bool) io.Writer {
	if enableErrorsOnly {
		return os.Stderr
	}
	return defaultLogOutput
}
