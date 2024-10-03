package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/leonardocartaxo/expenses-analyzer/internal/user"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/users", user.Router)
	addr := ":3000"
	fmt.Printf("Starting GO API service at %s ...\n", addr)
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}
