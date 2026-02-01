package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Define simple ANSI codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	// colorYellow = "\033[33m"
)

func init() {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
		NoColor:    false, // Ensure colors are enabled
	}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //HadleFunc is of type func
		start := time.Now()
		lrw := &loggingResponseWriter{w, http.StatusOK}

		next.ServeHTTP(lrw, r)

		// 1. Determine the color for the status code ONLY
		var statusColor string
		switch {
		case lrw.statusCode >= 500:
			statusColor = colorRed
		case lrw.statusCode >= 400:
			statusColor = colorRed
		default:
			statusColor = colorGreen
		}

		// 2. Format the status code with the color strings
		coloredStatus := fmt.Sprintf("%s%d%s", statusColor, lrw.statusCode, colorReset)

		// 3. Print the log
		// Note: We use log.Info() for everything now so the whole line 
		// isn't red—only the status code we just formatted.
		log.Info().Msgf("%-3s %s %s %s",
			r.Method,
			r.URL.Path,
			coloredStatus,
			time.Since(start),
		)
	})
}