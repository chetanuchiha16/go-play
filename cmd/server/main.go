package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	swaggerui "github.com/alexliesenfeld/go-swagger-ui"
	"github.com/chetanuchiha16/go-play/internal/api"
	"github.com/chetanuchiha16/go-play/internal/config"
	"github.com/chetanuchiha16/go-play/internal/database"
	"github.com/chetanuchiha16/go-play/internal/domain/auth"
	"github.com/chetanuchiha16/go-play/internal/domain/user"
	"github.com/chetanuchiha16/go-play/internal/middleware"
	"github.com/chetanuchiha16/go-play/internal/server"
)

//go:generate oapi-codegen -generate types,std-http,spec -package api -o ../../internal/api/generated.go ../../api/openapi.yaml

func main() {
	// ── Config & DB ───────────────────────────────────────────────
	cfg := config.Load()
	pool, err := database.NewPool(context.Background(), cfg.DATABASE_URL)
	if err != nil {
		log.Fatal("error connecting to the db")
	}
	defer pool.Close()
	store := database.NewStore(pool)

	// ── Services ──────────────────────────────────────────────────
	authService := auth.NewAuthService(store, []byte(cfg.JWT_SECRET))
	userService := user.NewUserService(store)

	// ── Domain Handlers ───────────────────────────────────────────
	userHandler := user.NewHandler(userService)
	authHandler := auth.NewAuthHandler(authService)

	// ── Composed Server (satisfies api.ServerInterface) ───────────
	srv := server.NewServer(userHandler, authHandler, pool)

	// ── Router ────────────────────────────────────────────────────
	mux := http.NewServeMux()

	// Mount oapi-codegen generated routes with per-operation auth middleware.
	api.HandlerWithOptions(srv, api.StdHTTPServerOptions{
		BaseRouter: mux,
		Middlewares: []api.MiddlewareFunc{
			authBearerMiddleware(authService),
		},
	})

	// Serve OpenAPI spec as YAML
	mux.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/openapi.yaml")
	})

	// Swagger UI at /docs
	mux.Handle("GET /docs", swaggerui.NewHandler(
		swaggerui.WithSpecURL("/openapi.yaml"),
		swaggerui.WithPersistAuthorization(true),
		swaggerui.WithTryItOutEnabled(true),
	))

	// ── Global middleware stack ────────────────────────────────────
	mw := middleware.NewMiddlewareManager()
	var handler http.Handler = mux
	handler = mw.LoggerMiddleware(handler)
	handler = mw.RequestIdMiddleware(handler)
	handler = mw.CorsMiddleware(handler)

	// ── HTTP Server ───────────────────────────────────────────────
	httpServer := &http.Server{
		Addr:              ":9999",
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("Server listening on %s", httpServer.Addr)
		log.Printf("Swagger UI: http://localhost%s/docs", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// ── Graceful shutdown ─────────────────────────────────────────
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("Server exiting")
}

// authBearerMiddleware returns an oapi-codegen MiddlewareFunc that enforces
// JWT Bearer auth ONLY on operations that have security requirements.
func authBearerMiddleware(authSvc auth.AuthService) api.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			scopes := r.Context().Value(api.BearerAuthScopes)
			if scopes == nil {
				next.ServeHTTP(w, r)
				return
			}
			authSvc.AuthMiddleware(next).ServeHTTP(w, r)
		})
	}
}
