package dao

import (
	"adorable-star/internal/model"
	"adorable-star/internal/pkg/util"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Init() {
	// Open database connection
	DB, _ = gorm.Open(sqlite.Open(util.GetCwd()+"/storage/db/dev.db"), &gorm.Config{
		NamingStrategy:         schema.NamingStrategy{SingularTable: true},
		SkipDefaultTransaction: true,
		DisableAutomaticPing:   true,
		PrepareStmt:            true,
	})

	// Use Write-Ahead Logging (WAL) mode
	DB.Exec("PRAGMA journal_mode=WAL;")
	DB.Exec("PRAGMA SYNCHRONOUS=NORMAL")

	// Set connection pool size
	db, _ := DB.DB()
	db.SetMaxIdleConns(5)

	// Migrate struct model to database
	DB.AutoMigrate(&model.User{}, &model.JupiterData{}, &model.Course{}, &model.Assignment{}, &model.Message{})
}
