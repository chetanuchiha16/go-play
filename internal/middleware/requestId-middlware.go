package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func RequestIdMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		request_id := fmt.Sprintf("%d", time.Now().UnixNano())
		ctx := context.WithValue(r.Context(), "request_id", request_id) // take request context and store in our context
		w.Header().Set("X-Request-ID", request_id)
		next.ServeHTTP(w, r.WithContext(ctx)) // now take our updated context and put it back to requst context
	})
}