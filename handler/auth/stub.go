package auth

import (
	"errors"
	"fmt"
	"github.com/afrizuko/kilango/model"
)

type Stub struct {
	counter uint
	users   map[string]model.User
}

func NewStub() *Stub {
	stub := new(Stub)
	stub.counter = 1
	stub.users = map[string]model.User{}
	for i := 1; i <= 5; i++ {
		username := fmt.Sprintf("%d", i)
		stub.users[username] = model.User{
			ID:       uint(i),
			Username: username,
			Password: "00" + username,
		}
	}
	return stub
}

func (s *Stub) Authenticate(req model.AuthRequest) (model.User, error) {
	if user, exists := s.users[req.Username]; exists && user.Password == req.Password {
		return user, nil
	}
	return model.User{}, errors.New("user with specified credentials not found")
}

func (s *Stub) GetUserProfile(id uint) (model.User, error) {
	for _, user := range s.users {
		if user.ID == id {
			return user, nil
		}
	}
	return model.User{}, errors.New("user with specified credentials not found")
}
