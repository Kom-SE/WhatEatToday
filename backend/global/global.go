package global

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	DB      *gorm.DB      //数据库连接实例
	RedisDB *redis.Client // Redis连接实例
)
