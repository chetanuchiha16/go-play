package main

import (
	"context"
	"log"

	"github.com/chetanuchiha16/go-play/internal/config"
	"github.com/chetanuchiha16/go-play/internal/database"
	"github.com/chetanuchiha16/go-play/internal/domain/user"
	"github.com/chetanuchiha16/go-play/internal/middleware"
	"github.com/go-fuego/fuego"
)

func main() {
	cfg := config.Load()
	pool, err := database.NewPool(context.Background(), cfg.DATABASE_URL)
	if err != nil {
		log.Fatal("error connecting to the db")
	}
	defer pool.Close()

	// 1. Setup our Store and Service
	userStore := database.NewStore(pool)
	userService := user.NewService(userStore)
	userHandler := user.NewHandler(userService)

	// 2. Initialize Fuego Server
	s := fuego.NewServer(
		fuego.WithAddr("localhost:8080"),
	)

	// 3. Add Middlewares (Fuego accepts standard middleware)
	fuego.Use(s, middleware.RequestIdMiddleWare)
	fuego.Use(s, middleware.LoggerMiddleware)

	// 4. Routes (Fuego handles the methods and paths)
	fuego.Post(s, "/users", userHandler.CreateUser)
	fuego.Get(s, "/users/{id}", userHandler.GetUser)
	fuego.Get(s, "/users", userHandler.ListUser)
	fuego.Post(s, "/login", userHandler.Login)

	// Secure routes
	authGroup := fuego.Group(s, "")
	fuego.Use(authGroup, middleware.AuthMiddleware)
	fuego.Delete(authGroup, "/users/{id}", userHandler.DeleteUser)

	// 5. Run (Handles graceful shutdown automatically)
	s.Run()
}