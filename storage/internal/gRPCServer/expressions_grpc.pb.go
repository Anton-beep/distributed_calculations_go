// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: expressions.proto

package gRPCServer

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

// ExpressionsServiceClient is the client API for ExpressionsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExpressionsServiceClient interface {
	GetUpdates(ctx context.Context, in *Empty, opts ...grpc.CallOption) (ExpressionsService_GetUpdatesClient, error)
	ConfirmStartCalculating(ctx context.Context, in *Expression, opts ...grpc.CallOption) (*Confirm, error)
	PostResult(ctx context.Context, in *Expression, opts ...grpc.CallOption) (*Message, error)
	KeepAlive(ctx context.Context, in *KeepAliveMsg, opts ...grpc.CallOption) (*Empty, error)
	GetOperationsAndTimes(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*OperationsAndTimes, error)
}

type expressionsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewExpressionsServiceClient(cc grpc.ClientConnInterface) ExpressionsServiceClient {
	return &expressionsServiceClient{cc}
}

func (c *expressionsServiceClient) GetUpdates(ctx context.Context, in *Empty, opts ...grpc.CallOption) (ExpressionsService_GetUpdatesClient, error) {
	stream, err := c.cc.NewStream(ctx, &ExpressionsService_ServiceDesc.Streams[0], "/storage.ExpressionsService/GetUpdates", opts...)
	if err != nil {
		return nil, err
	}
	x := &expressionsServiceGetUpdatesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ExpressionsService_GetUpdatesClient interface {
	Recv() (*Expression, error)
	grpc.ClientStream
}

type expressionsServiceGetUpdatesClient struct {
	grpc.ClientStream
}

func (x *expressionsServiceGetUpdatesClient) Recv() (*Expression, error) {
	m := new(Expression)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *expressionsServiceClient) ConfirmStartCalculating(ctx context.Context, in *Expression, opts ...grpc.CallOption) (*Confirm, error) {
	out := new(Confirm)
	err := c.cc.Invoke(ctx, "/storage.ExpressionsService/ConfirmStartCalculating", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *expressionsServiceClient) PostResult(ctx context.Context, in *Expression, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/storage.ExpressionsService/PostResult", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *expressionsServiceClient) KeepAlive(ctx context.Context, in *KeepAliveMsg, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/storage.ExpressionsService/KeepAlive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *expressionsServiceClient) GetOperationsAndTimes(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*OperationsAndTimes, error) {
	out := new(OperationsAndTimes)
	err := c.cc.Invoke(ctx, "/storage.ExpressionsService/GetOperationsAndTimes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExpressionsServiceServer is the server API for ExpressionsService service.
// All implementations must embed UnimplementedExpressionsServiceServer
// for forward compatibility
type ExpressionsServiceServer interface {
	GetUpdates(*Empty, ExpressionsService_GetUpdatesServer) error
	ConfirmStartCalculating(context.Context, *Expression) (*Confirm, error)
	PostResult(context.Context, *Expression) (*Message, error)
	KeepAlive(context.Context, *KeepAliveMsg) (*Empty, error)
	GetOperationsAndTimes(context.Context, *Empty) (*OperationsAndTimes, error)
	mustEmbedUnimplementedExpressionsServiceServer()
}

// UnimplementedExpressionsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedExpressionsServiceServer struct {
}

func (UnimplementedExpressionsServiceServer) GetUpdates(*Empty, ExpressionsService_GetUpdatesServer) error {
	return status.Errorf(codes.Unimplemented, "method GetUpdates not implemented")
}
func (UnimplementedExpressionsServiceServer) ConfirmStartCalculating(context.Context, *Expression) (*Confirm, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfirmStartCalculating not implemented")
}
func (UnimplementedExpressionsServiceServer) PostResult(context.Context, *Expression) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostResult not implemented")
}
func (UnimplementedExpressionsServiceServer) KeepAlive(context.Context, *KeepAliveMsg) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KeepAlive not implemented")
}
func (UnimplementedExpressionsServiceServer) GetOperationsAndTimes(context.Context, *Empty) (*OperationsAndTimes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOperationsAndTimes not implemented")
}
func (UnimplementedExpressionsServiceServer) mustEmbedUnimplementedExpressionsServiceServer() {}

// UnsafeExpressionsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExpressionsServiceServer will
// result in compilation errors.
type UnsafeExpressionsServiceServer interface {
	mustEmbedUnimplementedExpressionsServiceServer()
}

func RegisterExpressionsServiceServer(s grpc.ServiceRegistrar, srv ExpressionsServiceServer) {
	s.RegisterService(&ExpressionsService_ServiceDesc, srv)
}

func _ExpressionsService_GetUpdates_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ExpressionsServiceServer).GetUpdates(m, &expressionsServiceGetUpdatesServer{stream})
}

type ExpressionsService_GetUpdatesServer interface {
	Send(*Expression) error
	grpc.ServerStream
}

type expressionsServiceGetUpdatesServer struct {
	grpc.ServerStream
}

func (x *expressionsServiceGetUpdatesServer) Send(m *Expression) error {
	return x.ServerStream.SendMsg(m)
}

func _ExpressionsService_ConfirmStartCalculating_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Expression)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExpressionsServiceServer).ConfirmStartCalculating(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/storage.ExpressionsService/ConfirmStartCalculating",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExpressionsServiceServer).ConfirmStartCalculating(ctx, req.(*Expression))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExpressionsService_PostResult_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Expression)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExpressionsServiceServer).PostResult(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/storage.ExpressionsService/PostResult",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExpressionsServiceServer).PostResult(ctx, req.(*Expression))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExpressionsService_KeepAlive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeepAliveMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExpressionsServiceServer).KeepAlive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/storage.ExpressionsService/KeepAlive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExpressionsServiceServer).KeepAlive(ctx, req.(*KeepAliveMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExpressionsService_GetOperationsAndTimes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExpressionsServiceServer).GetOperationsAndTimes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/storage.ExpressionsService/GetOperationsAndTimes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExpressionsServiceServer).GetOperationsAndTimes(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// ExpressionsService_ServiceDesc is the grpc.ServiceDesc for ExpressionsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ExpressionsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "storage.ExpressionsService",
	HandlerType: (*ExpressionsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ConfirmStartCalculating",
			Handler:    _ExpressionsService_ConfirmStartCalculating_Handler,
		},
		{
			MethodName: "PostResult",
			Handler:    _ExpressionsService_PostResult_Handler,
		},
		{
			MethodName: "KeepAlive",
			Handler:    _ExpressionsService_KeepAlive_Handler,
		},
		{
			MethodName: "GetOperationsAndTimes",
			Handler:    _ExpressionsService_GetOperationsAndTimes_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetUpdates",
			Handler:       _ExpressionsService_GetUpdates_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "expressions.proto",
}
