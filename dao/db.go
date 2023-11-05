package dao

import (
	"adorable-star/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Init() {
	DB, _ = gorm.Open(sqlite.Open("./dev.db"), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	DB.AutoMigrate(&model.User{}, &model.JupiterData{}, &model.Course{}, &model.Assignment{}, &model.Message{})
}
