package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc/idl/pb"
	"io"
	"log"
	"strconv"
)

// grpcConn 打开 grpc 服务端链接 当前是禁用链接和加密的
func grpcConn(addr string) (*grpc.ClientConn, error) {
	// 打开 grpc 服务端链接 当前是禁用链接和加密的
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// grpcClose 关闭 grpc 连接
func grpcClose(conn *grpc.ClientConn) {
	err := conn.Close()
	fmt.Println("进入关闭")
	if err != nil {
		panic(err)
	}
}

func main() {
	// 创建链接
	conn, _ := grpcConn(":8888")
	// 关闭链接
	defer grpcClose(conn)

	// 创建 user gRPC 客户端
	userClient := pb.NewUserServiceClient(conn)
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
	userResp1, err := userClient.Login(context.Background(), userReq1)
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("Login Result: %v", userResp1))

	fmt.Println(fmt.Sprintf("------------ %v ------------", "grpc_request_mode"))

	// 创建 grpc_request_model gRPC 客户端
	grpcRequestModelClient := pb.NewGrpcRequestModeClient(conn)

	grpcRequestModelReq := &pb.GrpcEmpty{}

	grpcRequestModelResp, grpcRequestModelErr := grpcRequestModelClient.Ping(context.Background(), grpcRequestModelReq)
	if grpcRequestModelErr != nil {
		panic(grpcRequestModelErr)
	}
	fmt.Println(fmt.Sprintf("Ping Result: %v", grpcRequestModelResp))

	sumCli, err := grpcRequestModelClient.Sum(context.Background())
	if err != nil {
		panic(err)
	}
	_ = sumCli.Send(&pb.SumRequest{Num: int64(1)})
	_ = sumCli.Send(&pb.SumRequest{Num: int64(2)})
	_ = sumCli.Send(&pb.SumRequest{Num: int64(3)})
	_ = sumCli.Send(&pb.SumRequest{Num: int64(4)})
	recv, err := sumCli.CloseAndRecv()
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("Sum Result: %v", recv))

	clientStream(grpcRequestModelClient, "这是客户端发送的数据")

	ServerStream(grpcRequestModelClient, &pb.StreamRequest{Input: "我是一只小老虎"})
	streaming(grpcRequestModelClient)

}

func clientStream(client pb.GrpcRequestModeClient, input string) {
	stream, err := client.ClientStream(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, s := range input {
		fmt.Println(fmt.Sprintf("Client Stream Send: %v", string(s)))
		// 发送数据
		err := stream.Send(&pb.StreamRequest{Input: string(s)})
		if err != nil {
			log.Fatal(err)
		}
	}

	// 接受并关闭流
	recv, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("Client Stream Recv: %v", recv.Output))
}

func ServerStream(client pb.GrpcRequestModeClient, request *pb.StreamRequest) {
	fmt.Println(fmt.Sprintf("Server Stream Send: %v", request.Input))
	// 发送数据
	stream, err := client.ServerStream(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}

	var str string
	for {
		// 接受数据
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		str += res.Output
		fmt.Println(fmt.Sprintf("Server Stream Recv: %v", res.Output))
	}
	fmt.Println(fmt.Sprintf("Server Stream Recv: %v", str))
}

func streaming(client pb.GrpcRequestModeClient) {
	streamingClient, err := client.Streaming(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		// 发送数据
		err := streamingClient.Send(&pb.StreamRequest{Input: strconv.Itoa(i)})
		if err != nil {
			log.Fatal(err)
		}

		// 接受数据
		recv, err := streamingClient.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(fmt.Sprintf("Streaming Recv: %v", recv.Output))
	}
	err = streamingClient.CloseSend()
	if err != nil {
		log.Fatal(err)
	}
}
