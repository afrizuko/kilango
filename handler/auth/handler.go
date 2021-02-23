package auth

import (
	"github.com/audit/model"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"net/http"
	"os"
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
	handler.tokenAuth = jwtauth.New("HS256", []byte(secret), nil)

	mux := chi.NewRouter()
	if os.Getenv("STATE") != "prod" {
		mux.Use(render.SetContentType(render.ContentTypeJSON))
	}

	mux.Post("/", handler.Authenticate)
	handler.Handler = mux
	return handler
}

func (h Handler) Authenticate(writer http.ResponseWriter, request *http.Request) {

	var authRequest model.AuthRequest
	_ = render.Bind(request, &authRequest)

	user, err := h.service.Authenticate(authRequest)
	if err != nil {
		render.Status(request, http.StatusUnauthorized)
		render.JSON(writer, request, map[string]string{"error": err.Error()})
		return
	}
	_, tokenString, _ := h.tokenAuth.Encode(map[string]interface{}{"user_id": user.ID})

	response := model.AuthResponse{
		Token:     tokenString,
		ExpiresAt: 0,
	}
	render.JSON(writer, request, response)
}
