package driver

import (
	"fmt"
	"time"
)

type FileCache struct {
	FilePath string
	Driver
}

var FileOptions = map[string]interface{}{
	"expire":     0,
	"prefix":     "",
	"path":       "",
	"tag_prefix": "",
}

func NewFileCache(options map[string]interface{}) error {

	return nil
}

// Get 读取缓存
func (f *FileCache) Get(key string) ([]byte, error) {
	return nil, nil
}

// Set 写入缓存
func (f *FileCache) Set(key string, value []byte, expiration time.Duration) error {
	return nil
}

// Delete 删除缓存
func (f *FileCache) Delete(key string) error {
	return nil
}

// GetCacheKey 获取实际的缓存标识
func (f *FileCache) GetCacheKey(key string) string {
	return fmt.Sprintf("%s%s", f.Driver.Options["prefix"], key)
}
