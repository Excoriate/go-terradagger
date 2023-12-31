package daggerio

import (
	"context"
	"os"
	"testing"

	"github.com/Excoriate/go-terradagger/pkg/o11y"

	"github.com/stretchr/testify/assert"
)

func TestResolveDaggerLogConfig(t *testing.T) {
	tests := []struct {
		name              string
		enableErrorsOnly  bool
		expectedLogOutput *os.File
	}{
		{
			name:              "Should return the default log output",
			enableErrorsOnly:  false,
			expectedLogOutput: os.Stdout,
		},
		// Another test case
		{
			name:              "Should return the stderr log output",
			enableErrorsOnly:  true,
			expectedLogOutput: os.Stderr,
		},
	}

	bc := NewBackend(context.Background(), o11y.NewLogger(o11y.LoggerOptions{
		EnableJSONHandler: false,
		EnableStdError:    true,
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logOutput := bc.ResolveDaggerLogConfig(tt.enableErrorsOnly)
			assert.Equal(t, tt.expectedLogOutput, logOutput)
		})
	}
}
