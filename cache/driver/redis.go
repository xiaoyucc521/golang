package driver

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/maps"
)

type RedisCache struct {
	client *redis.Client
	Driver
}

var RedisOptions = map[string]interface{}{
	"addr":       ":6379",
	"password":   "",
	"db":         0,
	"timeout":    0,
	"expire":     0,
	"prefix":     "",
	"tag_prefix": "",
}

// NewRedisClient 创建 redis 链接
func NewRedisClient(options map[string]interface{}) (*RedisCache, error) {
	// 俩map merge,结果为第一个参数
	maps.Copy(RedisOptions, options)

	client := redis.NewClient(&redis.Options{
		Addr:     RedisOptions["addr"].(string),
		Password: RedisOptions["password"].(string),
		DB:       RedisOptions["db"].(int),
	})

	// 测试Redis连接
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
		Driver: Driver{Options: RedisOptions},
	}, nil
}

// Get 读取缓存
func (r *RedisCache) Get(key string) ([]byte, error) {
	ctx := context.Background()
	key = r.GetCacheKey(key)
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	return val, nil
}

// Set 写入缓存
func (r *RedisCache) Set(key string, value []byte, expiration time.Duration) error {
	ctx := context.Background()
	key = r.GetCacheKey(key)
	err := r.client.Set(ctx, key, value, expiration).Err()
	return err
}

// Delete 删除缓存
func (r *RedisCache) Delete(key string) error {
	ctx := context.Background()
	key = r.GetCacheKey(key)
	err := r.client.Del(ctx, key).Err()
	return err
}
