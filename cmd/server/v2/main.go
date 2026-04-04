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
	"github.com/chetanuchiha16/go-play/internal/config"
	"github.com/chetanuchiha16/go-play/internal/database"
	"github.com/chetanuchiha16/go-play/internal/domain/auth"
	"github.com/chetanuchiha16/go-play/internal/domain/user"
	"github.com/chetanuchiha16/go-play/internal/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
)

// openApiHandler returns a professional Swagger UI 5.x handler
func openApiHandler(specURL string) http.Handler {
	return swaggerui.NewHandler(
		swaggerui.WithSpecURL(specURL),
		swaggerui.WithPersistAuthorization(true),
		swaggerui.WithTryItOutEnabled(true),
	)
}

func main() {
	cfg := config.Load()
	pool, err := database.NewPool(context.Background(), cfg.DATABASE_URL)
	if err != nil {
		log.Fatal("error connecting to the db")
	}
	defer pool.Close()
	store := database.NewStore(pool)

	authService := auth.NewAuthService(store, []byte(cfg.JWT_SECRET))
	authHandler := auth.NewAuthHandler(authService)

	userService := user.NewUserService(store)
	userHandler := user.NewUserHandler(userService)

	s := fuego.NewServer(

		fuego.WithSecurity(openapi3.SecuritySchemes{
			"bearerAuth": &openapi3.SecuritySchemeRef{
				Value: openapi3.NewJWTSecurityScheme(),
			},
		}),

		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				UIHandler: openApiHandler,
			}),
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				// This overrides the default "Cheatsheet" information
				Info: &openapi3.Info{
					Title:       "Go-Play",
					Version:     "1.0.0",
					Description: " ", // Use a space or your own custom text
				},
			}),
		),

		
	)

	// s.OpenAPI.Description() = 

	mw := middleware.NewMiddlewareManager()
	fuego.Use(s, mw.CorsMiddleware)
	fuego.Use(s, mw.RequestIdMiddleware)
	fuego.Use(s, mw.LoggerMiddleware)

	fuego.Get(s, "/health", func(c fuego.ContextNoBody) (string, error) {

		err := pool.Ping(c.Context())
		if err != nil {
			return "Service unavailable", fuego.InternalServerError{
				Status: http.StatusServiceUnavailable,
				Detail: err.Error(),
				Title:  "Database connection failed",
			}
		}
		return "OK", nil
	})
	authHandler.RegisterAuthRoutes(s, authService.AuthMiddleware)
	userHandler.RegisterUserRoutes(s, authService.AuthMiddleware)

	go func() {
    if err := s.Run(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("listen: %s\n", err)
    }
}()

// Wait for interrupt signal
stop := make(chan os.Signal, 1)
signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
<-stop

log.Println("Shutting down gracefully...")

// Create context with timeout for shutdown
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := s.Shutdown(ctx); err != nil {
    log.Fatal("Server forced to shutdown: ", err)
}

log.Println("Server exiting")

}
