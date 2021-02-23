package model

import (
	"net/http"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a AuthRequest) Bind(_ *http.Request) error {
	return nil
}

type AuthResponse struct {
	Token     string  `json:"token"`
	ExpiresAt float32 `json:"expires_at"`
}

type AuthService interface {
	Authenticate(req AuthRequest) (User, error)
}
