package dao

import (
	"adorable-star/model"
	"time"
)

var User = &UserDAO{}

type UserDAO struct{}

func (*UserDAO) GetUserByID(id int) (*model.User, error) {
	var user model.User
	err := DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (*UserDAO) GetUserByUsernameOrEmail(name string) (*model.User, error) {
	var user model.User
	err := DB.Table("user").Where("email = ? OR username = ?", name, name).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (*UserDAO) InsertUser(email string, username string, password string) error {
	user := &model.User{
		Email:      email,
		Username:   username,
		Password:   password,
		ActiveTime: time.Now(),
	}
	return DB.Create(user).Error
}

func (*UserDAO) UpdateActiveTime(id int) error {
	return DB.Model(&model.User{ID: id}).Update("active_time", time.Now()).Error
}
