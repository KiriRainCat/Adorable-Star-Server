package dao

import (
	"adorable-star/internal/model"
	"adorable-star/internal/pkg/config"
	"adorable-star/internal/pkg/util"
	"strconv"

	"github.com/glebarez/sqlite"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB
var Redis *redis.Client

func InitDB() {
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

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Host + ":" + strconv.Itoa(config.Config.Redis.Port),
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.DB,
	})

	_, err := Redis.Ping().Result()
	if err != nil {
		panic(err)
	}
}
