package utils

import (
	"main/global"
	"time"
)

const (
	RefreshTokenTime = 30 * 24 * time.Hour // 刷新令牌的过期时间为30天
	RefreshPrefix    = "rtoken:"           // 刷新令牌的前缀
)

// Redis存储刷新令牌rtoken
func SetRtoken(userID uint8, rtoken string) error {
	// 设置rtoken的过期时间为30天
	err := global.RedisDB.Set(RefreshPrefix+string(userID), rtoken, RefreshTokenTime).Err()
	if err != nil {
		return err
	}
	return nil
}
