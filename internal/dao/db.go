package dao

import (
	"adorable-star/internal/global"
	"adorable-star/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Init() {
	DB, _ = gorm.Open(sqlite.Open(global.GetCwd()+"/dev.db"), &gorm.Config{
		NamingStrategy:         schema.NamingStrategy{SingularTable: true},
		SkipDefaultTransaction: true,
	})
	DB.AutoMigrate(&model.User{}, &model.JupiterData{}, &model.Course{}, &model.Assignment{}, &model.Message{})
}
