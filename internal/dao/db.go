package dao

import (
	"adorable-star/internal/model"
	"adorable-star/internal/pkg/config"
	"adorable-star/internal/pkg/util"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB
var Redis *redis.Client

func InitDB() {
	// Build DSN string for PostgreSQL connection
	dsn := "user=" + config.Config.Postgresql.User +
		" password=" + config.Config.Postgresql.Password +
		" port=" + strconv.Itoa(config.Config.Postgresql.Port) +
		" sslmode=disable" +
		" TimeZone=Asia/Shanghai"

	if gin.Mode() == gin.ReleaseMode {
		dsn = dsn + " dbname=" + config.Config.Postgresql.DB + " host=" + config.Config.Postgresql.Host
	} else {
		dsn = dsn + " dbname=" + config.Config.Postgresql.DB + "-DEV" + " host=" + config.Config.Postgresql.DevHost
	}

	// Open database connection with PostgreSQL
	writer, _ := util.GetFileWriter("log/db.log")
	DB, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy:         schema.NamingStrategy{SingularTable: true},
		SkipDefaultTransaction: true,
		DisableAutomaticPing:   true,
		Logger: logger.New(log.New(writer, "\n", log.LstdFlags), logger.Config{
			Colorful: false,
			LogLevel: logger.Warn,
		}),
	})

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
	if err != nil && gin.Mode() == gin.ReleaseMode {
		panic(err)
	}
}
