package dao

import (
	"adorable-star/internal/model"
	"adorable-star/internal/pkg/util"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Init() {
	// Init writers
	writer, _ := os.OpenFile(util.GetCwd()+"/storage/log/sql.log", os.O_CREATE, os.ModeAppend)
	writers := []io.Writer{writer}
	if gin.Mode() != gin.ReleaseMode {
		writers = append(writers, os.Stdout)
	}

	// Init logger
	logger := logger.New(
		log.New(io.MultiWriter(writers...), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
		},
	)

	// Open database connection
	DB, _ = gorm.Open(sqlite.Open(util.GetCwd()+"/storage/db/dev.db"), &gorm.Config{
		NamingStrategy:         schema.NamingStrategy{SingularTable: true},
		SkipDefaultTransaction: true,
		DisableAutomaticPing:   true,
		PrepareStmt:            true,
		Logger:                 logger,
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
