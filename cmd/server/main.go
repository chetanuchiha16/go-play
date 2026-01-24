package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"github.com/chetanuchiha16/go-play/internal/config"
	"github.com/chetanuchiha16/go-play/internal/domain/user"
	"github.com/chetanuchiha16/go-play/pkg/database"
)

func main() {

	cfg := config.Load()
	pool, err := database.NewPool(context.Background(), cfg.DATABASE_URL)
	if err != nil {
		log.Fatal("error connecting to the db")
	}

	userStore := user.NewStore(pool)
	userService := user.NewService(userStore)
	userHandler := user.NewHandler(userService)

	mux := http.NewServeMux()

	// mux.HandleFunc("/", user.Hello_Hina)

	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUser)
	mux.HandleFunc("DELETE /users/{id}", userHandler.DeleteUser)

	fmt.Println("listening")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}

}
