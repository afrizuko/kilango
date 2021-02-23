package user

import (
	"errors"
	"github.com/afrizuko/kilango/model"
)

type Stub struct {
	counter uint
	users   map[uint]model.User
}

func NewStub() *Stub {
	stub := new(Stub)
	stub.counter = 1
	stub.users = map[uint]model.User{}
	for i := 1; i <= 5; i++ {
		stub.users[uint(i)] = model.User{
			ID: uint(i),
		}
	}
	return stub
}

func (s *Stub) GetUsers(_, limit int) ([]model.User, error) {
	var users []model.User
	count := 0
	for _, user := range s.users {
		if count >= limit {
			break
		}
		users = append(users, user)
		count++
	}
	return users, nil
}

func (s *Stub) GetUser(id uint) (model.User, error) {
	if user, exists := s.users[id]; exists {
		return user, nil
	}

	return model.User{}, errors.New("user not found")
}

func (s *Stub) CreateUser(user *model.User) error {
	s.counter++
	user.ID = s.counter
	s.users[s.counter] = *user
	return nil
}

func (s *Stub) ModifyUser(id uint, user *model.User) error {
	if _, exists := s.users[id]; exists {
		s.users[id] = *user
		return nil
	}
	return errors.New("user with specified id not found")
}

func (s *Stub) DeleteUser(id uint) error {
	delete(s.users, id)
	return nil
}

func (s *Stub) PurgeUser(id uint) error {
	return s.DeleteUser(id)
}
