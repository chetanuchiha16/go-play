package main

import (
	"fmt"
	"net/http"

	"github.com/chetanuchiha16/go-play/internal/domain/user"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", user.Hello_Hina)

	mux.HandleFunc("POST /users", user.CreateUser)
	mux.HandleFunc("GET /users/{id}", user.GetUser)
	mux.HandleFunc("DELETE /users/{id}", user.DeleteUser)

	fmt.Println("listening")
	http.ListenAndServe(":8080", mux)
}
