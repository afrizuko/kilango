package model

import (
	"errors"
	"github.com/afrizuko/kilango/util"
	"gorm.io/gorm"
	"log"
)

var userServiceImpl *UserServiceImpl

type UserServiceImpl struct {
	*gorm.DB
}

func NewUserServiceImpl() *UserServiceImpl {

	// promote reusability
	if userServiceImpl != nil {
		return userServiceImpl
	}

	userServiceImpl := new(UserServiceImpl)
	userServiceImpl.DB = NewConnection()
	if err := userServiceImpl.DB.AutoMigrate(&User{}); err != nil {
		log.Fatal(err)
	}

	userServiceImpl.SetupDefaultUser()
	return userServiceImpl
}

func (s *UserServiceImpl) SetupDefaultUser() {

	var user User
	s.DB.Where("username=?", "SYSTEM").First(&user)

	pin, _ := util.HashPin("SYSTEM")
	if user.ID < 1 {
		s.DB.Save(&User{
			Name:     "SYSTEM ADMIN",
			Username: "SYSTEM",
			Password: pin,
			Status:   "A",
		})
	}
}

func (s *UserServiceImpl) GetUsers(page, limit int) ([]User, error) {
	var users []User
	if err := s.DB.Offset(page).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserServiceImpl) GetUser(id uint) (User, error) {
	var user User
	if err := s.DB.First(&user, id).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (s *UserServiceImpl) CreateUser(user *User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	return s.DB.Save(user).Error
}

func (s *UserServiceImpl) ModifyUser(id uint, user *User) error {
	var result User
	if err := s.DB.First(&result, id).Error; err != nil {
		return err
	}
	return s.DB.Save(user).Error
}

func (s *UserServiceImpl) DeleteUser(id uint) error {
	if err := s.DB.Delete(&User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) PurgeUser(id uint) error {

	if err := s.DB.Unscoped().Delete(&User{
		ID: id,
	}).Error; err != nil {
		return err
	}
	return nil
}
