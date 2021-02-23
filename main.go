package main

import (
	"github.com/audit/handler/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
)

func main() {

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Group(func(r chi.Router) {
		r.Mount("/users", user.DefaultHandler())
	})

	log.Fatal(http.ListenAndServe(":3000", mux))
}
