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
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	TokenType string `json:"token_type"`
}

type AuthService interface {
	Authenticate(req AuthRequest) (User, error)
	GetUserProfile(id uint) (User, error)
}
