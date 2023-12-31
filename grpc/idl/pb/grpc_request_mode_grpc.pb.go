// 定义proto语法版本，这里指定使用proto3版本

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.22.5
// source: grpc_request_mode.proto

// 这里随便定义个包名

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	GrpcRequestMode_Ping_FullMethodName         = "/grpc.user.grpcRequestMode/Ping"
	GrpcRequestMode_ServerStream_FullMethodName = "/grpc.user.grpcRequestMode/ServerStream"
	GrpcRequestMode_Sum_FullMethodName          = "/grpc.user.grpcRequestMode/Sum"
	GrpcRequestMode_ClientStream_FullMethodName = "/grpc.user.grpcRequestMode/ClientStream"
	GrpcRequestMode_Streaming_FullMethodName    = "/grpc.user.grpcRequestMode/Streaming"
)

// GrpcRequestModeClient is the client API for GrpcRequestMode service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GrpcRequestModeClient interface {
	// 一、简单模式（Simple RPC）: 客户端发起请求并等待服务端响应，就是普通的 Ping-Pong 模式。
	// 空函数，用与检测客户端和服务端是否通畅
	Ping(ctx context.Context, in *GrpcEmpty, opts ...grpc.CallOption) (*GrpcPingResponse, error)
	// 二、服务端流式（Server-side streaming RPC）: 服务端发送数据，客户端接收数据。客户端发送请求到服务器，拿到一个流去读取返回的消息序列。客户端读取返回的流，直到里面没有任何消息。
	ServerStream(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (GrpcRequestMode_ServerStreamClient, error)
	// 三、客户端流式（Client-side streaming RPC）: 与服务端数据流模式相反，这次是客户端源源不断的向服务端发送数据流，而在发送结束后，由服务端返回一个响应。
	// 计算求和的方式来测试服务端流
	Sum(ctx context.Context, opts ...grpc.CallOption) (GrpcRequestMode_SumClient, error)
	ClientStream(ctx context.Context, opts ...grpc.CallOption) (GrpcRequestMode_ClientStreamClient, error)
	// 四、双向流式（Bidirectional streaming RPC）: 双方使用读写流去发送一个消息序列，两个流独立操作，双方可以同时发送和同时接收。
	Streaming(ctx context.Context, opts ...grpc.CallOption) (GrpcRequestMode_StreamingClient, error)
}

type grpcRequestModeClient struct {
	cc grpc.ClientConnInterface
}

func NewGrpcRequestModeClient(cc grpc.ClientConnInterface) GrpcRequestModeClient {
	return &grpcRequestModeClient{cc}
}

func (c *grpcRequestModeClient) Ping(ctx context.Context, in *GrpcEmpty, opts ...grpc.CallOption) (*GrpcPingResponse, error) {
	out := new(GrpcPingResponse)
	err := c.cc.Invoke(ctx, GrpcRequestMode_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcRequestModeClient) ServerStream(ctx context.Context, in *StreamRequest, opts ...grpc.CallOption) (GrpcRequestMode_ServerStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &GrpcRequestMode_ServiceDesc.Streams[0], GrpcRequestMode_ServerStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &grpcRequestModeServerStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GrpcRequestMode_ServerStreamClient interface {
	Recv() (*StreamResponse, error)
	grpc.ClientStream
}

type grpcRequestModeServerStreamClient struct {
	grpc.ClientStream
}

func (x *grpcRequestModeServerStreamClient) Recv() (*StreamResponse, error) {
	m := new(StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *grpcRequestModeClient) Sum(ctx context.Context, opts ...grpc.CallOption) (GrpcRequestMode_SumClient, error) {
	stream, err := c.cc.NewStream(ctx, &GrpcRequestMode_ServiceDesc.Streams[1], GrpcRequestMode_Sum_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &grpcRequestModeSumClient{stream}
	return x, nil
}

type GrpcRequestMode_SumClient interface {
	Send(*SumRequest) error
	CloseAndRecv() (*SumResponse, error)
	grpc.ClientStream
}

type grpcRequestModeSumClient struct {
	grpc.ClientStream
}

func (x *grpcRequestModeSumClient) Send(m *SumRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *grpcRequestModeSumClient) CloseAndRecv() (*SumResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(SumResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *grpcRequestModeClient) ClientStream(ctx context.Context, opts ...grpc.CallOption) (GrpcRequestMode_ClientStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &GrpcRequestMode_ServiceDesc.Streams[2], GrpcRequestMode_ClientStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &grpcRequestModeClientStreamClient{stream}
	return x, nil
}

type GrpcRequestMode_ClientStreamClient interface {
	Send(*StreamRequest) error
	CloseAndRecv() (*StreamResponse, error)
	grpc.ClientStream
}

type grpcRequestModeClientStreamClient struct {
	grpc.ClientStream
}

func (x *grpcRequestModeClientStreamClient) Send(m *StreamRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *grpcRequestModeClientStreamClient) CloseAndRecv() (*StreamResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *grpcRequestModeClient) Streaming(ctx context.Context, opts ...grpc.CallOption) (GrpcRequestMode_StreamingClient, error) {
	stream, err := c.cc.NewStream(ctx, &GrpcRequestMode_ServiceDesc.Streams[3], GrpcRequestMode_Streaming_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &grpcRequestModeStreamingClient{stream}
	return x, nil
}

type GrpcRequestMode_StreamingClient interface {
	Send(*StreamRequest) error
	Recv() (*StreamResponse, error)
	grpc.ClientStream
}

type grpcRequestModeStreamingClient struct {
	grpc.ClientStream
}

func (x *grpcRequestModeStreamingClient) Send(m *StreamRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *grpcRequestModeStreamingClient) Recv() (*StreamResponse, error) {
	m := new(StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GrpcRequestModeServer is the server API for GrpcRequestMode service.
// All implementations must embed UnimplementedGrpcRequestModeServer
// for forward compatibility
type GrpcRequestModeServer interface {
	// 一、简单模式（Simple RPC）: 客户端发起请求并等待服务端响应，就是普通的 Ping-Pong 模式。
	// 空函数，用与检测客户端和服务端是否通畅
	Ping(context.Context, *GrpcEmpty) (*GrpcPingResponse, error)
	// 二、服务端流式（Server-side streaming RPC）: 服务端发送数据，客户端接收数据。客户端发送请求到服务器，拿到一个流去读取返回的消息序列。客户端读取返回的流，直到里面没有任何消息。
	ServerStream(*StreamRequest, GrpcRequestMode_ServerStreamServer) error
	// 三、客户端流式（Client-side streaming RPC）: 与服务端数据流模式相反，这次是客户端源源不断的向服务端发送数据流，而在发送结束后，由服务端返回一个响应。
	// 计算求和的方式来测试服务端流
	Sum(GrpcRequestMode_SumServer) error
	ClientStream(GrpcRequestMode_ClientStreamServer) error
	// 四、双向流式（Bidirectional streaming RPC）: 双方使用读写流去发送一个消息序列，两个流独立操作，双方可以同时发送和同时接收。
	Streaming(GrpcRequestMode_StreamingServer) error
	mustEmbedUnimplementedGrpcRequestModeServer()
}

// UnimplementedGrpcRequestModeServer must be embedded to have forward compatible implementations.
type UnimplementedGrpcRequestModeServer struct {
}

func (UnimplementedGrpcRequestModeServer) Ping(context.Context, *GrpcEmpty) (*GrpcPingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedGrpcRequestModeServer) ServerStream(*StreamRequest, GrpcRequestMode_ServerStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ServerStream not implemented")
}
func (UnimplementedGrpcRequestModeServer) Sum(GrpcRequestMode_SumServer) error {
	return status.Errorf(codes.Unimplemented, "method Sum not implemented")
}
func (UnimplementedGrpcRequestModeServer) ClientStream(GrpcRequestMode_ClientStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ClientStream not implemented")
}
func (UnimplementedGrpcRequestModeServer) Streaming(GrpcRequestMode_StreamingServer) error {
	return status.Errorf(codes.Unimplemented, "method Streaming not implemented")
}
func (UnimplementedGrpcRequestModeServer) mustEmbedUnimplementedGrpcRequestModeServer() {}

// UnsafeGrpcRequestModeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GrpcRequestModeServer will
// result in compilation errors.
type UnsafeGrpcRequestModeServer interface {
	mustEmbedUnimplementedGrpcRequestModeServer()
}

func RegisterGrpcRequestModeServer(s grpc.ServiceRegistrar, srv GrpcRequestModeServer) {
	s.RegisterService(&GrpcRequestMode_ServiceDesc, srv)
}

func _GrpcRequestMode_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GrpcEmpty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcRequestModeServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GrpcRequestMode_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcRequestModeServer).Ping(ctx, req.(*GrpcEmpty))
	}
	return interceptor(ctx, in, info, handler)
}

func _GrpcRequestMode_ServerStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(StreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GrpcRequestModeServer).ServerStream(m, &grpcRequestModeServerStreamServer{stream})
}

type GrpcRequestMode_ServerStreamServer interface {
	Send(*StreamResponse) error
	grpc.ServerStream
}

type grpcRequestModeServerStreamServer struct {
	grpc.ServerStream
}

func (x *grpcRequestModeServerStreamServer) Send(m *StreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _GrpcRequestMode_Sum_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GrpcRequestModeServer).Sum(&grpcRequestModeSumServer{stream})
}

type GrpcRequestMode_SumServer interface {
	SendAndClose(*SumResponse) error
	Recv() (*SumRequest, error)
	grpc.ServerStream
}

type grpcRequestModeSumServer struct {
	grpc.ServerStream
}

func (x *grpcRequestModeSumServer) SendAndClose(m *SumResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *grpcRequestModeSumServer) Recv() (*SumRequest, error) {
	m := new(SumRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GrpcRequestMode_ClientStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GrpcRequestModeServer).ClientStream(&grpcRequestModeClientStreamServer{stream})
}

type GrpcRequestMode_ClientStreamServer interface {
	SendAndClose(*StreamResponse) error
	Recv() (*StreamRequest, error)
	grpc.ServerStream
}

type grpcRequestModeClientStreamServer struct {
	grpc.ServerStream
}

func (x *grpcRequestModeClientStreamServer) SendAndClose(m *StreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *grpcRequestModeClientStreamServer) Recv() (*StreamRequest, error) {
	m := new(StreamRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GrpcRequestMode_Streaming_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GrpcRequestModeServer).Streaming(&grpcRequestModeStreamingServer{stream})
}

type GrpcRequestMode_StreamingServer interface {
	Send(*StreamResponse) error
	Recv() (*StreamRequest, error)
	grpc.ServerStream
}

type grpcRequestModeStreamingServer struct {
	grpc.ServerStream
}

func (x *grpcRequestModeStreamingServer) Send(m *StreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *grpcRequestModeStreamingServer) Recv() (*StreamRequest, error) {
	m := new(StreamRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GrpcRequestMode_ServiceDesc is the grpc.ServiceDesc for GrpcRequestMode service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GrpcRequestMode_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.user.grpcRequestMode",
	HandlerType: (*GrpcRequestModeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _GrpcRequestMode_Ping_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ServerStream",
			Handler:       _GrpcRequestMode_ServerStream_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Sum",
			Handler:       _GrpcRequestMode_Sum_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "ClientStream",
			Handler:       _GrpcRequestMode_ClientStream_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Streaming",
			Handler:       _GrpcRequestMode_Streaming_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "grpc_request_mode.proto",
}
