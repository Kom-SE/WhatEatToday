package controllers

import (
	"context"
	"encoding/json"
	"log"
	"main/global"
	"main/models"
	"time"

	"github.com/robfig/cron/v3"
)

const (
	TopRecipesRedisKey = "top_recipes:likes"
	TopRecipesExpire   = 1 * time.Hour // 缓存1小时
	TopRecipesLimit    = 10            // 获取前10条
)

// TopRecipeCache 缓存的菜谱结构
type TopRecipeCache struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Images      string    `json:"images"`
	Likes       uint      `json:"likes"`
	FoodID      string    `json:"food_id"`
	AuthorID    uint      `json:"author_id"`
	CookTime    string    `json:"cook_time"`
	Process     string    `json:"process"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Index       int       `json:"index,omitempty"` // 可选字段，表示排名
}

// 每30分钟更新前十点赞数的菜谱
func UpdateTopRecipes() {
	log.Println("开始更新前十点赞数菜谱...")

	var topRecipes []models.Recipe
	// 查询点赞数最多的前10个菜谱，并预加载用户信息
	err := global.DB.Order("likes DESC").
		Limit(TopRecipesLimit).
		Find(&topRecipes).Error

	if err != nil {
		log.Printf("查询热门菜谱失败: %v", err)
		return
	}

	// 转换为缓存结构
	var cacheRecipes []TopRecipeCache
	for _, recipe := range topRecipes {
		cacheRecipe := TopRecipeCache{
			ID:          recipe.ID,
			Name:        recipe.Name,
			Description: recipe.Description,
			Images:      recipe.Images,
			Likes:       recipe.Likes,
			FoodID:      recipe.FoodID,
			AuthorID:    recipe.AuthorID,
			CookTime:    recipe.CookTime,
			Process:     recipe.Process,
			CreatedAt:   recipe.CreatedAt,
			UpdatedAt:   recipe.UpdatedAt,
		}
		cacheRecipes = append(cacheRecipes, cacheRecipe)
	}

	// 序列化为JSON
	jsonData, err := json.Marshal(cacheRecipes)
	if err != nil {
		log.Printf("序列化热门菜谱失败: %v", err)
		return
	}

	// 存储到Redis
	ctx := context.Background()
	err = global.RedisDB.Set(ctx, TopRecipesRedisKey, string(jsonData), TopRecipesExpire).Err()
	if err != nil {
		log.Printf("存储热门菜谱到Redis失败: %v", err)
		return
	}

	log.Printf("成功更新前十热门菜谱到Redis，共 %d 个", len(cacheRecipes))

	// // 打印详细信息
	// for i, recipe := range cacheRecipes {
	// 	log.Printf("第%d名: ID=%d, 标题=%s, 点赞数=%d",
	// 		i+1, recipe.ID, recipe.Name, recipe.Likes)
	// }
}

// 从Redis获取热门菜谱
func GetTopRecipesFromCache() ([]TopRecipeCache, error) {
	ctx := context.Background()

	// 从Redis获取
	jsonData, err := global.RedisDB.Get(ctx, TopRecipesRedisKey).Result()
	if err != nil {
		log.Printf("从Redis获取热门菜谱失败: %v", err)
		// 如果Redis中没有，直接从数据库查询并缓存
		UpdateTopRecipes()

		// 重新尝试从Redis获取
		jsonData, err = global.RedisDB.Get(ctx, TopRecipesRedisKey).Result()
		if err != nil {
			return nil, err
		}
	}

	// 反序列化
	var recipes []TopRecipeCache
	err = json.Unmarshal([]byte(jsonData), &recipes)
	if err != nil {
		log.Printf("反序列化热门菜谱失败: %v", err)
		return nil, err
	}

	return recipes, nil
}

// 定时任务设置
func StartCronJobs() *cron.Cron {
	c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(log.Writer(), "cron: ", log.LstdFlags))))

	// 每30分钟执行一次 (0 */30 * * * *)
	_, err := c.AddFunc("*/30 * * * *", UpdateTopRecipes)
	if err != nil {
		log.Fatalf("添加热门菜谱更新任务失败: %v", err)
	}

	// 启动定时任务
	c.Start()
	log.Println("定时任务已启动 - 每30分钟更新热门菜谱")

	// 程序启动时立即执行一次
	go UpdateTopRecipes()

	return c
}

// 停止定时任务
func StopCronJobs(c *cron.Cron) {
	if c != nil {
		c.Stop()
		log.Println("定时任务已停止")
	}
}

// 清理Redis中的热门菜谱缓存
func ClearTopRecipesCache() error {
	ctx := context.Background()
	err := global.RedisDB.Del(ctx, TopRecipesRedisKey).Err()
	if err != nil {
		log.Printf("清理热门菜谱缓存失败: %v", err)
		return err
	}
	log.Println("热门菜谱缓存已清理")
	return nil
}

// 获取缓存统计信息
func GetCacheStats() map[string]interface{} {
	ctx := context.Background()

	// 检查缓存是否存在
	exists, _ := global.RedisDB.Exists(ctx, TopRecipesRedisKey).Result()

	stats := map[string]interface{}{
		"cache_exists": exists > 0,
		"cache_key":    TopRecipesRedisKey,
		"expire_time":  TopRecipesExpire.String(),
		"limit":        TopRecipesLimit,
	}

	if exists > 0 {
		// 获取TTL
		ttl, _ := global.RedisDB.TTL(ctx, TopRecipesRedisKey).Result()
		stats["ttl"] = ttl.String()

		// 获取缓存大小
		size, _ := global.RedisDB.StrLen(ctx, TopRecipesRedisKey).Result()
		stats["cache_size_bytes"] = size
	}

	return stats
}
