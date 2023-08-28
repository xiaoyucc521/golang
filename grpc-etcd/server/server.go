package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"grpc-etcd/config"
	"grpc-etcd/discovery"
	"grpc-etcd/idl/pb"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 初始化配置
	config.InitConfig()

	endpoints := config.Conf.Etcd.Address
	etcd := discovery.NewRegister(endpoints, logrus.New())
	address := config.Conf.Server.Host

	// 当前服务配置
	serverInfo := discovery.Server{
		Project: config.Conf.Project.Name,   // 项目名
		Name:    config.Conf.Server.Name,    // 服务名
		Addr:    config.Conf.Server.Host,    // 服务地址
		Version: config.Conf.Server.Version, // 服务版本
		Weight:  config.Conf.Server.Weight,  // 服务权重
	}

	log.Println(serverInfo)
	// 服务注册
	if err := etcd.Register(serverInfo, 10); err != nil {
		panic(fmt.Sprintf("etcdRegister.Register ERROR: %v", err))
	}

	// 创建 user gRPC 服务端
	grpcServer := grpc.NewServer()
	// 注册实现的服务实例
	userService := NewUserService()
	pb.RegisterUserServiceServer(grpcServer, userService)

	// 监听端口
	listen, err := net.Listen("tcp", address)
	if err != nil {
		panic(fmt.Sprintf("net.Listen ERROR: %v", err))
	}

	err = grpcServer.Serve(listen)
	if err != nil {
		panic(fmt.Sprintf("grpcServer.Serve ERROR: %v", err))
	}

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
