package main

import (
	"fmt"
	"sync"

	"snowflake/pkg/snowflake"
)

var work = snowflake.NewWorker(30, 30)

func main() {
	var wg sync.WaitGroup

	count := 10000
	ch := make(chan int64, count)

	wg.Add(count)
	defer close(ch)
	//并发 count个goroutine 进行 snowFlake ID 生成
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			id := work.GenerateID()
			ch <- id
		}()
	}
	wg.Wait()
	m := make(map[int64]int)
	for i := 0; i < count; i++ {
		id := <-ch
		// 如果 map 中存在为 id 的 key, 说明生成的 snowflake ID 有重复
		_, ok := m[id]
		if ok {
			fmt.Printf("repeat id %d\n", id)
			return
		}
		// 将 id 作为 key 存入 map
		m[id] = i
	}
	// 成功生成 snowflake ID
	fmt.Println("All", len(m), "snowflake ID Get successed!")
}
