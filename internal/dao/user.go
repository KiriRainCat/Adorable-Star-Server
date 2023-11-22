package dao

import (
	"adorable-star/internal/model"
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

func (*UserDAO) GetUsers() (users []*model.User, err error) {
	err = DB.Find(&users).Error
	return
}

func (*UserDAO) GetActiveUsers() (users []*model.User, err error) {
	err = DB.Where("active_at > ?", time.Now().Add(-time.Hour*24)).Find(&users).Error
	return
}

func (*UserDAO) InsertUser(email string, username string, password string) error {
	user := &model.User{
		Email:    email,
		Username: username,
		Password: password,
		ActiveAt: time.Now(),
	}
	return DB.Create(user).Error
}

func (*UserDAO) UpdatePassword(id int, pwd string) error {
	return DB.Model(&model.User{ID: id}).Update("password", pwd).Error
}

func (*UserDAO) UpdateStatus(id int, status int) error {
	return DB.Model((&model.User{ID: id})).Update("status", status).Error
}

func (*UserDAO) UpdateActiveTime(id int) error {
	// Using update column to avoid `updated_at`` change
	return DB.Model(&model.User{ID: id}).UpdateColumn("active_at", time.Now()).Error
}
