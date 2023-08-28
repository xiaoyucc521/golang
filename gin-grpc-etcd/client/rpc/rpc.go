package rpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"

	"gin-grpc-etcd/config"
	"gin-grpc-etcd/discovery"
	"gin-grpc-etcd/discovery/balancer/weight"
	"gin-grpc-etcd/idl/pb"
)

var (
	Register   *discovery.Resolver
	ctx        context.Context
	CancelFunc context.CancelFunc

	userClient pb.UserServiceClient
)

func Init() {
	// etcd 地址
	endpoints := config.Conf.Etcd.Addr

	Register = discovery.NewResolver(endpoints, logrus.New())
	resolver.Register(Register)
	ctx, CancelFunc = context.WithTimeout(context.Background(), 3*time.Second)

	defer CancelFunc()
	defer Register.Close()

	initClient(discovery.Server{
		Project: config.Conf.Project.Name,
		Name:    config.Conf.Services["user"].Name,
		Version: config.Conf.Server.Version,
	}, &userClient)

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
		opts = append(opts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, weight.Name)))
		//opts = append(opts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
	}

	conn, err = grpc.DialContext(ctx, addr, opts...)
	return conn, err
}
