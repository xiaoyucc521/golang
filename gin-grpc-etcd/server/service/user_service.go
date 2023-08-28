package service

import (
	"context"
	"fmt"
	"log"

	"gin-grpc-etcd/idl/pb"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u UserService) Ping(ctx context.Context, empty *pb.Empty) (*pb.PingResponse, error) {
	return &pb.PingResponse{
		Message: "请求成功",
	}, nil
}

func (u UserService) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
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
