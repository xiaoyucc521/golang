package main

import (
	"context"
	"fmt"
	clientV3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	cli, err := clientV3.New(clientV3.Config{
		Endpoints:   []string{"xiaoyucc521.com:2379"},
		DialTimeout: time.Second * 5, // 建立连接失败的超时时
	})

	if err != nil {
		fmt.Println(fmt.Sprintf("---------connect to etcd failed, err:%v", err))
		return
	}

	defer func(cli *clientV3.Client) {
		err := cli.Close()
		if err != nil {
			fmt.Println(fmt.Sprintf("---------close to etcd failed, err:%v", err))
		}
	}(cli)

	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	key := "test_etcd_key"
	value := "test_etcd_value"
	put, err := cli.Put(ctx, key, value)
	if err != nil {
		fmt.Println(fmt.Sprintf("---------put to etcd failed, err:%v", err))
		return
	}
	fmt.Println(put)
	cancel()

	// get
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	get, err := cli.Get(ctx, key)
	if err != nil {
		fmt.Println(fmt.Sprintf("---------get to etcd failed, err:%v", err))
		return
	}
	fmt.Println(get.Kvs)
	for _, v := range get.Kvs {
		fmt.Println(fmt.Sprintf("------key:%s value:%s", v.Key, v.Value))
	}
	cancel()

	// del
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	response, err := cli.Delete(ctx, key)
	if err != nil {
		fmt.Println(fmt.Sprintf("---------delete to etcd failed, err:%v", err))
		return
	}
	fmt.Println(response)

	//value1 := `[{"path":"D:/learn/go/log-collector-lmh/log_agent_etcd/log_file/log1","topic":"log1"},{"path":"D:/learn/go/log-collector-lmh/log_agent_etcd/log_file/log2","topic":"log2"}]`
	//fmt.Println(value1)
}
