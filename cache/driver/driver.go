package driver

import "fmt"

type Driver struct {
	Options map[string]interface{}
}

// GetCacheKey 获取实际的缓存标识
func (d *Driver) GetCacheKey(key string) string {
	return fmt.Sprintf("%s%s", d.Options["prefix"], key)
}
