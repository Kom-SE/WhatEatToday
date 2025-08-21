package utils

import (
	"fmt"
	"main/global"
	"time"
)

const (
	RefreshTokenTime = 30 * 24 * time.Hour // 刷新令牌的过期时间为30天
	RefreshPrefix    = "rtoken:"           // 刷新令牌的前缀
)

// Redis存储刷新令牌rtoken
func SetRtoken(userID uint8, rtoken string) error {
	rname := fmt.Sprintf("%s%d", RefreshPrefix, userID)
	ctx := global.RedisDB.Context() // 获取Redis上下文
	// 设置rtoken的过期时间为30天
	err := global.RedisDB.Set(ctx, rname, rtoken, RefreshTokenTime).Err()
	if err != nil {
		return err
	}
	return nil
}
