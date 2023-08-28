package main

import (
	"context"
	"fmt"
	"grpc/idl/pb"
	"log"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

// NewUserService 工厂函数创建服务实例
func NewUserService() *UserService {
	return &UserService{}
}

// Ping 空函数，用与检测客户端和服务端是否通畅
func (u *UserService) Ping(ctx context.Context, empty *pb.Empty) (*pb.PingResponse, error) {
	return &pb.PingResponse{
		Message: "请求成功",
	}, nil
}

// Login 登录
func (u *UserService) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Println(fmt.Sprintf("User Login: %v", request.Username))
	if request.GetUsername() == "admin" && request.GetPassword() == "123456" {
		return &pb.LoginResponse{
			Code:    1,
			Message: "登录成功",
		}, nil
	}

	return &pb.LoginResponse{
		Code:    0,
		Message: "登录失败",
	}, nil
}
