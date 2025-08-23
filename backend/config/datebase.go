package config

import (
	"log"
	"main/global"
	"main/models"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() {
	// 启动mysql数据库服务
	dsn := AppConfig.Database.DSN
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to initialize database, got error: %v", err)
	}

	sqlDB, err := db.DB()

	sqlDB.SetMaxIdleConns(AppConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(AppConfig.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err != nil {
		log.Fatalf("Failed to configure database, got error: %v", err)
	}

	// // 自动迁移数据库表
	err = db.AutoMigrate(
		&models.User{},
		&models.Food{},
		&models.Label{},
		&models.Recipe{},
		&models.Comment{},
		&models.CommentLike{},
		&models.CollectedRecipes{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database tables, got error: %v", err)
	}

	global.DB = db
}

// 初始化Redis
func InitRedis() {
	// 启动Redis服务
	redisConfig := AppConfig.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password, // no password set
		DB:       redisConfig.DB,       // use default DB
		PoolSize: redisConfig.PoolSize,
	})

	global.RedisDB = rdb
}
