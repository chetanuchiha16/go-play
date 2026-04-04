package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/chetanuchiha16/go-play/pkg"
)

const id pkg.ContextKey = "request_id"

func (MiddlewareManager) RequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		request_id := fmt.Sprintf("%d", time.Now().UnixNano())
		ctx := context.WithValue(r.Context(), id, request_id) // take request context and store in our context
		w.Header().Set("X-Request-ID", request_id)
		next.ServeHTTP(w, r.WithContext(ctx)) // now take our updated context and put it back to requst context
	})
}
