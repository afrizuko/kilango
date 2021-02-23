package model

import (
	"errors"
	"github.com/afrizuko/kilango/util"
	"gorm.io/gorm"
)

var authServiceImpl *AuthServiceImpl

type AuthServiceImpl struct {
	*gorm.DB
}

func NewAuthServiceImpl() *AuthServiceImpl {
	// promote reusability
	if authServiceImpl != nil {
		return authServiceImpl
	}

	authServiceImpl := new(AuthServiceImpl)
	authServiceImpl.DB = NewConnection()
	return authServiceImpl
}

func (a *AuthServiceImpl) Authenticate(req AuthRequest) (User, error) {
	var user User
	err := a.DB.Where("username=?", req.Username).First(&user).Error

	if err != nil {
		return user, err
	}

	if !util.VerifyHash(user.Password, req.Password) {
		return user, errors.New("invalid credentials ")
	}
	return user, nil
}

func (a *AuthServiceImpl) GetUserProfile(id uint) (User, error) {
	var user User
	if err := a.DB.First(&user, id).Error; err != nil {
		return user, err
	}
	return user, nil
}
