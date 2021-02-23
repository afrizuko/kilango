package main

import (
	"github.com/afrizuko/kilango/handler/auth"
	"github.com/afrizuko/kilango/handler/user"
	"github.com/afrizuko/kilango/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	secret := os.Getenv("HASH_SECRET")
	if secret == "" {
		secret = "BbcWorldServices"
	}

	tokenAuth = jwtauth.New("HS256", []byte(secret), secret)
}

func main() {

	mux := chi.NewRouter()
	mux.Use(util.Logger("router"))
	mux.Use(middleware.Recoverer)

	mux.Use(middleware.RequestID)
	mux.Use(middleware.URLFormat)

	// unauthenticated APIs
	authH := auth.DefaultHandler()
	mux.Group(func(r chi.Router) {
		r.Mount("/api/auth", authH)
	})

	//authenticated APIs
	mux.Group(func(r chi.Router) {
		r.Use(authH.Verifier)
		r.Use(auth.Authenticator)

		r.Mount("/api/users", user.DefaultHandler())
	})

	log.Fatal(http.ListenAndServe(":3000", mux))
}
