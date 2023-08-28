package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"grpc-etcd/config"
	"grpc-etcd/discovery"
	"grpc-etcd/discovery/balancer/weight"
	"grpc-etcd/idl/pb"
	"log"
	"time"
)

var (
	Register   *discovery.Resolver
	ctx        context.Context
	CancelFunc context.CancelFunc

	userClient pb.UserServiceClient
)

func main() {
	// 初始化配置
	config.InitConfig()

	endpoints := config.Conf.Etcd.Address
	Register = discovery.NewResolver(endpoints, logrus.New())
	resolver.Register(Register)

	serverInfo := discovery.Server{
		Project: config.Conf.Project.Name,          // 项目名
		Name:    config.Conf.Services["user"].Name, // 服务名
		Version: config.Conf.Server.Version,
	}

	ctx, CancelFunc = context.WithTimeout(context.Background(), 3*time.Second)
	defer CancelFunc()

	initClient(serverInfo, &userClient)

	userReq := &pb.Empty{}
	userResp, err := userClient.Ping(context.Background(), userReq)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("Ping Result: %v", userResp))
	userReq1 := &pb.LoginRequest{
		Username: "admin",
		Password: "123456",
	}

	for i := 0; i < 20; i++ {
		userResp1, err := userClient.Login(context.Background(), userReq1)
		if err != nil {
			panic(err)
		}
		fmt.Println(fmt.Sprintf("Login Result: %v", userResp1))
		time.Sleep(time.Second)
	}
}

func initClient(server discovery.Server, client interface{}) {
	conn, err := connectServer(server)
	if err != nil {
		panic(err)
	}

	switch c := client.(type) {
	case *pb.UserServiceClient:
		*c = pb.NewUserServiceClient(conn)
	}
}

func connectServer(server discovery.Server) (conn *grpc.ClientConn, err error) {
	opts := []grpc.DialOption{
		// 当前是没有使用证书（权限）认证的，不安全
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	keyPrefix := discovery.BuildPrefix(server)
	// Scheme://Authority/Endpoint
	addr := fmt.Sprintf("%s://%s/%s", Register.Scheme(), "", keyPrefix)

	// 实现负载均衡
	if config.Conf.Services[server.Name].LoadBalance {
		log.Printf("load balance enabled for %s\n", server.Name)
		opts = append(opts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, weight.Name))) // roundrobin.Name
	}

	conn, err = grpc.DialContext(ctx, addr, opts...)
	return conn, err
}
