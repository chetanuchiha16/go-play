package middleware

import (
	"net/http"
	"time"

	"github.com/chetanuchiha16/go-play/internal/logger"
)



type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) { // our handler calls this version
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
func (MiddlewareManager) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //HadleFunc is of type func, this is a type conversion so it satisfy the interface and can call ServeHttp method
		start := time.Now()
		lrw := &loggingResponseWriter{w, http.StatusOK}

		next.ServeHTTP(lrw, r) // call next by providing lrw and r

		reqId, ok := r.Context().Value("request_id").(string)
		if !ok {
			reqId = "unknown"
		}

		logger.Log.Info().
			Str("id", reqId).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", lrw.statusCode).
			Dur("latency", time.Since(start)).
			Msgf("%s %s %d", r.Method, r.URL.Path, lrw.statusCode)
	})
}
