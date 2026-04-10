package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// Log is the globally accessible logger instance.
var Log zerolog.Logger

func init() {
	// Provide a safe default before Init() is explicitly called
	Log = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

// Init configures the global Logger based on the environment.
// In production, it logs as structured JSON.
// In development, it uses pretty console output.
func Init(env string) {
	if env == "production" {
		Log = zerolog.New(os.Stdout).
			With().
			Timestamp().
			Caller().
			Logger()
	} else {
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
			NoColor:    false,
		}
		Log = zerolog.New(output).
			With().
			Timestamp().
			Caller().
			Logger()
	}
}
