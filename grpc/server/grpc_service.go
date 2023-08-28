package main

import (
	"context"
	"fmt"
	"grpc/idl/pb"
	"io"
	"log"
	"strconv"
)

type grpcRequestMode struct {
	pb.UnimplementedGrpcRequestModeServer
}

// newGrpcRequestModel 工厂函数创建服务实例
func newGrpcRequestModel() *grpcRequestMode {
	return &grpcRequestMode{}
}

func (g *grpcRequestMode) Ping(ctx context.Context, empty *pb.GrpcEmpty) (*pb.GrpcPingResponse, error) {
	return &pb.GrpcPingResponse{
		Message: "请求成功",
	}, nil
}

func (g *grpcRequestMode) ServerStream(in *pb.StreamRequest, stream pb.GrpcRequestMode_ServerStreamServer) error {
	input := in.Input
	fmt.Println(fmt.Sprintf("Server Stream Recv: %v", input))

	str := "我知道了，你是一只小老虎"

	for _, s := range str {
		err := stream.Send(&pb.StreamResponse{Output: string(s)})
		if err != nil {
			return err
		}
	}
	return nil
}

// Sum 案例--客户端流式处理
// 计算求和的方式来测试服务端流
func (g *grpcRequestMode) Sum(request pb.GrpcRequestMode_SumServer) error {
	var sum int64 = 0
	for {
		reqObj, err := request.Recv()
		if err == io.EOF {
			log.Println(fmt.Sprintf("Recv Sum err: %v", err))
			return request.SendAndClose(&pb.SumResponse{Result: sum})
		} else if err == nil {
			log.Println(fmt.Sprintf("get client request param = %v", reqObj.Num))
			sum += reqObj.Num
		} else {
			return err
		}
	}
}

func (g *grpcRequestMode) ClientStream(stream pb.GrpcRequestMode_ClientStreamServer) error {
	output := ""
	for {
		// 接受数据
		recv, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		output += recv.Input
	}
	fmt.Println(output)

	// 回应并关闭流
	return stream.SendAndClose(&pb.StreamResponse{Output: "这是服务端回应的数据"})
}

func (g *grpcRequestMode) Streaming(stream pb.GrpcRequestMode_StreamingServer) error {
	for n := 0; ; {
		// 接受数据
		Recv, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		v, _ := strconv.Atoi(Recv.Input)
		fmt.Println(fmt.Sprintf("Server Recv: %v", v))
		n += v

		fmt.Println(fmt.Sprintf("Server Stream Send: %v", n))
		// 发送数据
		err = stream.Send(&pb.StreamResponse{Output: strconv.Itoa(n)})
		if err != nil {
			return err
		}
	}
}
