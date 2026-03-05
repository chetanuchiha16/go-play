package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chetanuchiha16/go-play/internal/config"
	"github.com/chetanuchiha16/go-play/internal/database"
	"github.com/chetanuchiha16/go-play/internal/domain/user"
	"github.com/chetanuchiha16/go-play/internal/middleware"
	"github.com/getkin/kin-openapi/openapi3" // The missing import to fix the compiler error
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option" // Import for security options
)

func main() {
	cfg := config.Load()
	pool, err := database.NewPool(context.Background(), cfg.DATABASE_URL)
	if err != nil {
		log.Fatal("error connecting to the db")
	}
	defer pool.Close()

	userStore := database.NewStore(pool)
	userService := user.NewService(userStore)
	userHandler := user.NewHandler(userService)

	s := fuego.NewServer(
		fuego.WithAddr("localhost:8080"),
		// This EXACT block fixes the compiler error by using the expected openapi3 types
		fuego.WithSecurity(openapi3.SecuritySchemes{
			"bearerAuth": &openapi3.SecuritySchemeRef{
				Value: openapi3.NewJWTSecurityScheme(),
			},
		}),
	)

	fuego.Use(s, middleware.RequestIdMiddleWare)
	fuego.Use(s, middleware.LoggerMiddleware)

	fuego.Post(s, "/users", userHandler.CreateUser)
	fuego.Get(s, "/users/{id}", userHandler.GetUser)
	fuego.Get(s, "/users", userHandler.ListUser)
	fuego.Post(s, "/login", userHandler.Login)

	authGroup := fuego.Group(s, "")
	fuego.Use(authGroup, middleware.AuthMiddleware)

	// Use option.Security to tell Swagger this specific route needs the token
	fuego.Delete(authGroup, "/users/{id}", userHandler.DeleteUser, option.Security(openapi3.SecurityRequirement{"bearerAuth": []string{}}))

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
