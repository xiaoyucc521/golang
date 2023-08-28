package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"gin-grpc-etcd/config"
	"gin-grpc-etcd/discovery"
	"gin-grpc-etcd/idl/pb"
	"gin-grpc-etcd/server/service"
)

func main() {
	config.Init()

	// etcd 地址
	endpoints := config.Conf.Etcd.Addr

	// 创建 etcd 连接
	etcd := discovery.NewRegister(endpoints, logrus.New())

	// 当前服务配置信息
	serverInfo := discovery.Server{
		Project: config.Conf.Project.Name,
		Name:    config.Conf.Server.Name,
		Addr:    config.Conf.Server.Host,
		Version: config.Conf.Server.Version,
		Weight:  config.Conf.Server.Weight,
	}
	// 服务注册
	if err := etcd.Register(serverInfo, 10); err != nil {
		panic(fmt.Sprintf("etcdRegister.Register ERROR: %v", err))
	}

	// 创建 grpc
	grpcServer := grpc.NewServer()
	// 注册实现的服务实例
	userService := service.NewUserService()
	pb.RegisterUserServiceServer(grpcServer, userService)

	// 监听端口
	listen, err := net.Listen("tcp", config.Conf.Server.Host)
	if err != nil {
		panic(fmt.Sprintf("net.Listen ERROR: %v", err))
	}

	go func() {
		err = grpcServer.Serve(listen)
		if err != nil {
			panic(fmt.Sprintf("grpcServer.Serve ERROR: %v", err))
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
