package cache

import (
	"cache/driver"
	"errors"
	"time"
)

type Cache interface {
	// Get 读取缓存
	Get(key string) ([]byte, error)
	// Set 写入缓存
	Set(key string, value []byte, expiration time.Duration) error
	// Delete 删除缓存
	Delete(key string) error
}

func NewCache(cacheType string, options map[string]interface{}) (Cache, error) {
	switch cacheType {
	case "file":
		return nil, nil
	case "redis":
		return driver.NewRedisClient(options)
	default:
		return nil, errors.New("不支持的缓存类型")
	}
}
