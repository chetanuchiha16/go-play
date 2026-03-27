package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/chetanuchiha16/go-play/internal/config"
	"github.com/chetanuchiha16/go-play/internal/database"
	"github.com/chetanuchiha16/go-play/internal/domain/auth"
	"github.com/chetanuchiha16/go-play/internal/domain/user"
	"github.com/chetanuchiha16/go-play/internal/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
)

// openApiHandler returns a Swagger UI 5.x handler that supports OpenAPI 3.1.0
func openApiHandler(specURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		html := strings.ReplaceAll(swaggerUIHTML, "{{.SpecURL}}", specURL)
		w.Write([]byte(html))
	})
}

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({
      url: "{{.SpecURL}}",
      dom_id: '#swagger-ui',
      presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
      layout: "BaseLayout",
      persistAuthorization: true,
    });
  </script>
</body>
</html>`

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
		),
	)

	mw := middleware.NewMiddlewareManager()
	fuego.Use(s, mw.CorsMiddleware)
	fuego.Use(s, mw.RequestIdMiddleWare)
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

	stop := make(chan os.Signal, 1)
	go func() {
		log.Println("listening at localhost:8080")
		s.Run()
	}()
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	fmt.Println("kill recieved")
	fmt.Println("Shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("resourse releasing")
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("forcfully shutting down %v", err)
	}
	pool.Close()

}
