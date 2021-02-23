package auth

import (
	"github.com/afrizuko/kilango/model"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/lestrrat-go/jwx/jwt"
	"log"
	"net/http"
	"os"
	"time"
)

type Handler struct {
	http.Handler
	service   model.AuthService
	tokenAuth *jwtauth.JWTAuth
}

func NewHandler(service model.AuthService) *Handler {

	handler := new(Handler)
	handler.service = service
	secret := os.Getenv("HASH_SECRET")
	if secret == "" {
		secret = "BbcWorldServices"
	}
	handler.tokenAuth = jwtauth.New("HS256", []byte(secret), secret)

	mux := chi.NewRouter()
	if os.Getenv("STATE") != "prod" {
		mux.Use(render.SetContentType(render.ContentTypeJSON))
		mux.Group(func(r chi.Router) {
			r.Use(handler.Verifier)
			r.Use(Authenticator)
			r.Get("/user", handler.GetProfile)
		})
	} else {
		mux.Get("/user", handler.GetProfile)
	}

	mux.Post("/", handler.Authenticate)
	handler.Handler = mux
	return handler
}

func DefaultHandler() *Handler {

	handler := new(Handler)
	handler.service = model.NewAuthServiceImpl()
	secret := os.Getenv("HASH_SECRET")
	if secret == "" {
		secret = "BbcWorldServices"
	}
	handler.tokenAuth = jwtauth.New("HS256", []byte(secret), nil)

	mux := chi.NewRouter()
	if os.Getenv("STATE") != "prod" {
		mux.Use(render.SetContentType(render.ContentTypeJSON))
	}

	mux.Post("/", handler.Authenticate)
	mux.Get("/user", handler.GetProfile)
	handler.Handler = mux
	return handler
}

func (h *Handler) Authenticate(writer http.ResponseWriter, request *http.Request) {

	var authRequest model.AuthRequest
	_ = render.Bind(request, &authRequest)

	user, err := h.service.Authenticate(authRequest)
	if err != nil {
		render.Status(request, http.StatusUnauthorized)
		render.JSON(writer, request, map[string]string{"error": err.Error()})
		return
	}

	claims := map[string]interface{}{"user_id": user.ID}
	jwtauth.SetIssuedNow(claims)
	jwtauth.SetExpiryIn(claims, 15*time.Minute)
	_, tokenString, _ := h.tokenAuth.Encode(claims)

	_, err = jwtauth.VerifyToken(h.tokenAuth, tokenString)
	if err != nil {
		log.Println(err)
	}

	response := model.AuthResponse{
		Token:     tokenString,
		ExpiresAt: claims["exp"].(int64),
		TokenType: "Bearer",
	}
	render.JSON(writer, request, response)
}

func (h *Handler) GetProfile(writer http.ResponseWriter, request *http.Request) {
	var _, claims, _ = jwtauth.FromContext(request.Context())
	userID, ok := claims["user_id"]

	if !ok {
		render.Status(request, http.StatusNotFound)
		render.JSON(writer, request, map[string]string{"error": "invalid token"})
		return
	}

	if val, ok := userID.(uint); ok {
		userProfile, _ := h.service.GetUserProfile(val)
		render.JSON(writer, request, userProfile)
		return
	}

	render.Status(request, http.StatusUnauthorized)
	render.JSON(writer, request, map[string]string{"error": "un-parsable token"})
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": err.Error()})
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, map[string]string{"error": http.StatusText(http.StatusUnauthorized)})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h *Handler) Verifier(next http.Handler) http.Handler {
	return jwtauth.Verify(h.tokenAuth, jwtauth.TokenFromHeader)(next)
}
