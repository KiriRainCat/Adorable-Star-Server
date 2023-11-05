package service

import (
	"adorable-star/model"

	"gorm.io/gorm"
)

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db}
}

type UserService struct {
	db *gorm.DB
}

func (s *UserService) GetUserByUsernameOrEmail(name string) (*model.User, error) {
	var user model.User
	err := s.db.Table("user").Where("email = ? OR username = ?", name, name).First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) InsertUser(email string, username string, password string) error {
	user := &model.User{
		Email:    email,
		Username: username,
		Password: password,
	}
	err := s.db.Create(user).Error
	return err
}
