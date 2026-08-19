package store

import (
	"context"
	"log/slog"
	"time"

	"ragbox/config"

	"github.com/redis/go-redis/v9"
)

var GlobalCache *Cache

func init() {
	GlobalCache = NewCache(context.Background())
}

type Cache struct {
	source *redis.Client
}

func (c *Cache) Get(key string) (string, error) {
	val, err := c.source.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (c *Cache) SetWithExpire(key string, value string, expiration time.Duration) error {
	err := c.source.Set(context.Background(), key, value, expiration).Err()
	return err
}

func NewCache(ctx context.Context) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Config.Cache.Addr,
		Password: config.Config.Cache.Password, // no password set
		DB:       config.Config.Cache.DBIndex,  // use default DB
	})

	// 测试连接
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		slog.Error("Failed to connect to Redis:", "err", err)
		panic(err)
	}
	slog.Info("Redis connected:", "pong", pong)
	return &Cache{source: rdb}
}
