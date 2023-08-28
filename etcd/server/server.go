package main

import (
	"etcd/config"
	"etcd/discovery"
	"fmt"
	clientV3 "go.etcd.io/etcd/client/v3"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Register struct {
	cli           *clientV3.Client
	leaseID       clientV3.LeaseID // 租约ID
	keepAliveChan <-chan *clientV3.LeaseKeepAliveResponse
	key           string // key
	val           string // value
}

func main() {
	config.InitConfig()

	endpoints := config.Conf.Etcd.Address
	etcd := discovery.NewRegister(endpoints, log.Default())
	serverInfo := discovery.Server{
		Project: config.Conf.Project.Name, // 项目名
		Name:    config.Conf.Server.Name,  // 服务名
		Addr:    config.Conf.Server.Host,
		Version: config.Conf.Server.Version,
	}

	go func() {
		err := etcd.Register(serverInfo, 10)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// 监控系统型号，等待 ctrl + c 系统信号通知关闭
	exitCh := make(chan int, 1)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
		<-c
		etcd.Stop()
		time.Sleep(time.Second)
		close(exitCh)
	}()
	log.Println(fmt.Sprintf("exit %v", <-exitCh))
}
