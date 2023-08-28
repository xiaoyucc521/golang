package main

import (
	"fmt"
	"google.golang.org/grpc"
	"grpc/idl/pb"
	"net"
)

func main() {
	// 监听端口
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("listen on :8888")

	// 创建 user gRPC 服务端
	grpcServer := grpc.NewServer()
	// 注册实现的服务实例
	userService := NewUserService()
	pb.RegisterUserServiceServer(grpcServer, userService)

	// 创建 grpc_request_model gRPC 服务端
	grpcRequestModelService := newGrpcRequestModel()
	pb.RegisterGrpcRequestModeServer(grpcServer, grpcRequestModelService)

	// 启动 grpc 服务
	fmt.Println("gRPC is running...")
	err = grpcServer.Serve(listen)
	if err != nil {
		fmt.Println(err)
		return
	}
}
