// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.0--rc1
// source: common.proto

package common

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ServerServer_CallFuncOnServer_FullMethodName = "/common.ServerServer/CallFuncOnServer"
)

// ServerServerClient is the client API for ServerServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServerServerClient interface {
	CallFuncOnServer(ctx context.Context, in *Text, opts ...grpc.CallOption) (*Text, error)
}

type serverServerClient struct {
	cc grpc.ClientConnInterface
}

func NewServerServerClient(cc grpc.ClientConnInterface) ServerServerClient {
	return &serverServerClient{cc}
}

func (c *serverServerClient) CallFuncOnServer(ctx context.Context, in *Text, opts ...grpc.CallOption) (*Text, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Text)
	err := c.cc.Invoke(ctx, ServerServer_CallFuncOnServer_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServerServerServer is the server API for ServerServer service.
// All implementations must embed UnimplementedServerServerServer
// for forward compatibility.
type ServerServerServer interface {
	CallFuncOnServer(context.Context, *Text) (*Text, error)
	mustEmbedUnimplementedServerServerServer()
}

// UnimplementedServerServerServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedServerServerServer struct{}

func (UnimplementedServerServerServer) CallFuncOnServer(context.Context, *Text) (*Text, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CallFuncOnServer not implemented")
}
func (UnimplementedServerServerServer) mustEmbedUnimplementedServerServerServer() {}
func (UnimplementedServerServerServer) testEmbeddedByValue()                      {}

// UnsafeServerServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServerServerServer will
// result in compilation errors.
type UnsafeServerServerServer interface {
	mustEmbedUnimplementedServerServerServer()
}

func RegisterServerServerServer(s grpc.ServiceRegistrar, srv ServerServerServer) {
	// If the following call pancis, it indicates UnimplementedServerServerServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ServerServer_ServiceDesc, srv)
}

func _ServerServer_CallFuncOnServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Text)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServerServer).CallFuncOnServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ServerServer_CallFuncOnServer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServerServer).CallFuncOnServer(ctx, req.(*Text))
	}
	return interceptor(ctx, in, info, handler)
}

// ServerServer_ServiceDesc is the grpc.ServiceDesc for ServerServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ServerServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "common.ServerServer",
	HandlerType: (*ServerServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CallFuncOnServer",
			Handler:    _ServerServer_CallFuncOnServer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "common.proto",
}

const (
	ClientServer_CallFuncOnClient_FullMethodName = "/common.ClientServer/CallFuncOnClient"
)

// ClientServerClient is the client API for ClientServer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClientServerClient interface {
	CallFuncOnClient(ctx context.Context, in *Text, opts ...grpc.CallOption) (*Text, error)
}

type clientServerClient struct {
	cc grpc.ClientConnInterface
}

func NewClientServerClient(cc grpc.ClientConnInterface) ClientServerClient {
	return &clientServerClient{cc}
}

func (c *clientServerClient) CallFuncOnClient(ctx context.Context, in *Text, opts ...grpc.CallOption) (*Text, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Text)
	err := c.cc.Invoke(ctx, ClientServer_CallFuncOnClient_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClientServerServer is the server API for ClientServer service.
// All implementations must embed UnimplementedClientServerServer
// for forward compatibility.
type ClientServerServer interface {
	CallFuncOnClient(context.Context, *Text) (*Text, error)
	mustEmbedUnimplementedClientServerServer()
}

// UnimplementedClientServerServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedClientServerServer struct{}

func (UnimplementedClientServerServer) CallFuncOnClient(context.Context, *Text) (*Text, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CallFuncOnClient not implemented")
}
func (UnimplementedClientServerServer) mustEmbedUnimplementedClientServerServer() {}
func (UnimplementedClientServerServer) testEmbeddedByValue()                      {}

// UnsafeClientServerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClientServerServer will
// result in compilation errors.
type UnsafeClientServerServer interface {
	mustEmbedUnimplementedClientServerServer()
}

func RegisterClientServerServer(s grpc.ServiceRegistrar, srv ClientServerServer) {
	// If the following call pancis, it indicates UnimplementedClientServerServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ClientServer_ServiceDesc, srv)
}

func _ClientServer_CallFuncOnClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Text)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServerServer).CallFuncOnClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClientServer_CallFuncOnClient_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServerServer).CallFuncOnClient(ctx, req.(*Text))
	}
	return interceptor(ctx, in, info, handler)
}

// ClientServer_ServiceDesc is the grpc.ServiceDesc for ClientServer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClientServer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "common.ClientServer",
	HandlerType: (*ClientServerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CallFuncOnClient",
			Handler:    _ClientServer_CallFuncOnClient_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "common.proto",
}