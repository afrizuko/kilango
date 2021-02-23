package model

import (
	"gorm.io/gorm"
	"net/http"
	"time"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Name      string         `json:"name"`
	Username  string         `json:"username" gorm:"uniqueIndex,length:10"`
	Password  string         `json:"password" gorm:"length:10"`
	Status    string         `json:"status"`
}

func (u User) Bind(*http.Request) error {
	return nil
}

type UserService interface {
	GetUsers(page, limit int) ([]User, error)
	GetUser(id uint) (User, error)
	CreateUser(user *User) error
	ModifyUser(id uint, user *User) error
	DeleteUser(id uint) error
	PurgeUser(id uint) error
}
