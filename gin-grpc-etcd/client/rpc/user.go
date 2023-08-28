package rpc

import (
	"context"
	"gin-grpc-etcd/idl/pb"
)

func UserPing(ctx context.Context, req *pb.Empty) (*pb.PingResponse, error) {
	ping, err := userClient.Ping(ctx, req)
	if err != nil {
		return nil, err
	}

	return ping, nil
}

func UserLogin(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	login, err := userClient.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return login, nil
}
