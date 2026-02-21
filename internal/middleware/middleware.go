package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chetanuchiha16/go-play/internal/domain/user"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Define simple ANSI codes
const (
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	// colorYellow = "\033[33m"
)

func init() {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
		NoColor:    false, // Ensure colors are enabled
	}
	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) { // our handler calls this version
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func RequestIdMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r * http.Request) {

		request_id := fmt.Sprintf("%d", time.Now().UnixNano())
		ctx := context.WithValue(r.Context(), "request_id", request_id)
		w.Header().Set("X-Request-ID", request_id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //HadleFunc is of type func, this is a type conversion so it satisfy the interface and can call ServeHttp method
		start := time.Now()
		lrw := &loggingResponseWriter{w, http.StatusOK}

		next.ServeHTTP(lrw, r) // call next by providing lrw and r

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
		reqId, ok := r.Context().Value("request_id").(string)
		if !ok {
			reqId = "unknown"
		}
		log.Info().
		Str("id ", reqId).
		Msgf("%-3s %s %s %s",
			r.Method,
			r.URL.Path,
			coloredStatus,
			time.Since(start),
		)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorisation header required", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorisation format", http.StatusUnauthorized)
			return
		}

		claims, err := user.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return // do not forget to return after hitting errors
		}
		ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
