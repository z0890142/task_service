package app

import (
	"context"
	"fmt"
	"task_service/config"

	"github.com/redis/go-redis/v9"
)

func InitCacheHook(app *Application) error {
	cacheConfig := config.GetConfig().Cache

	addr := fmt.Sprintf("%s:%v", cacheConfig.Host, cacheConfig.Port)

	redisOpt := &redis.Options{
		Addr: addr, // Redis 伺服器位址
	}

	if cacheConfig.Password != "" {
		redisOpt.Password = cacheConfig.Password
	}

	rdb := redis.NewClient(redisOpt)
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("InitCacheHook: %v", err)
	}

	app.cacheClient = rdb
	return nil
}
