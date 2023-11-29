package logger

import (
	"github.com/rs/zerolog"
	"os"
)

func NewLogger(logLevel zerolog.Level) zerolog.Logger {
	return zerolog.New(os.Stdout).
		Level(logLevel).
		With().Timestamp().
		Logger()
}
