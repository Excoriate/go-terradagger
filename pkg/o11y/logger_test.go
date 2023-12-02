package o11y

import "testing"

func TestNewLogger(t *testing.T) {
  t.Run("stdout with text format", func(t *testing.T) {
    log := NewLogger(LoggerOptions{EnableJSONHandler: false, EnableStdError: false}).(*LogImpl)
    log.Info("tests")
  })
}
