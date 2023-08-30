package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	option := map[string]interface{}{
		"addr":     "127.0.0.1:6379",
		"password": "123456",
		"db":       0,
		"prefix":   "test:",
	}
	cache, _ := NewCache("redis", option)

	err := cache.Set("a", []byte("1"), time.Second*3600)
	if err != nil {
		fmt.Println(err)
	}
	res, _ := cache.Get("a")
	fmt.Println(string(res))
	_ = cache.Delete("a")
}
