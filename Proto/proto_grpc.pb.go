// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.3
// source: Proto/proto.proto

package Homework04

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
	ClientConnection_RequestAccess_FullMethodName = "/proto.ClientConnection/requestAccess"
	ClientConnection_Receive_FullMethodName       = "/proto.ClientConnection/receive"
	ClientConnection_Connection_FullMethodName    = "/proto.ClientConnection/Connection"
)

// ClientConnectionClient is the client API for ClientConnection service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClientConnectionClient interface {
	RequestAccess(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Empty, error)
	Receive(ctx context.Context, in *Response, opts ...grpc.CallOption) (*Empty, error)
	Connection(ctx context.Context, in *Greeting, opts ...grpc.CallOption) (*Empty, error)
}

type clientConnectionClient struct {
	cc grpc.ClientConnInterface
}

func NewClientConnectionClient(cc grpc.ClientConnInterface) ClientConnectionClient {
	return &clientConnectionClient{cc}
}

func (c *clientConnectionClient) RequestAccess(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, ClientConnection_RequestAccess_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientConnectionClient) Receive(ctx context.Context, in *Response, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, ClientConnection_Receive_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientConnectionClient) Connection(ctx context.Context, in *Greeting, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, ClientConnection_Connection_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClientConnectionServer is the server API for ClientConnection service.
// All implementations must embed UnimplementedClientConnectionServer
// for forward compatibility
type ClientConnectionServer interface {
	RequestAccess(context.Context, *Request) (*Empty, error)
	Receive(context.Context, *Response) (*Empty, error)
	Connection(context.Context, *Greeting) (*Empty, error)
	mustEmbedUnimplementedClientConnectionServer()
}

// UnimplementedClientConnectionServer must be embedded to have forward compatible implementations.
type UnimplementedClientConnectionServer struct {
}

func (UnimplementedClientConnectionServer) RequestAccess(context.Context, *Request) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestAccess not implemented")
}
func (UnimplementedClientConnectionServer) Receive(context.Context, *Response) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Receive not implemented")
}
func (UnimplementedClientConnectionServer) Connection(context.Context, *Greeting) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Connection not implemented")
}
func (UnimplementedClientConnectionServer) mustEmbedUnimplementedClientConnectionServer() {}

// UnsafeClientConnectionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClientConnectionServer will
// result in compilation errors.
type UnsafeClientConnectionServer interface {
	mustEmbedUnimplementedClientConnectionServer()
}

func RegisterClientConnectionServer(s grpc.ServiceRegistrar, srv ClientConnectionServer) {
	s.RegisterService(&ClientConnection_ServiceDesc, srv)
}

func _ClientConnection_RequestAccess_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientConnectionServer).RequestAccess(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClientConnection_RequestAccess_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientConnectionServer).RequestAccess(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClientConnection_Receive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Response)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientConnectionServer).Receive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClientConnection_Receive_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientConnectionServer).Receive(ctx, req.(*Response))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClientConnection_Connection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Greeting)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientConnectionServer).Connection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClientConnection_Connection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientConnectionServer).Connection(ctx, req.(*Greeting))
	}
	return interceptor(ctx, in, info, handler)
}

// ClientConnection_ServiceDesc is the grpc.ServiceDesc for ClientConnection service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClientConnection_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ClientConnection",
	HandlerType: (*ClientConnectionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "requestAccess",
			Handler:    _ClientConnection_RequestAccess_Handler,
		},
		{
			MethodName: "receive",
			Handler:    _ClientConnection_Receive_Handler,
		},
		{
			MethodName: "Connection",
			Handler:    _ClientConnection_Connection_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Proto/proto.proto",
}
