// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package sso

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UsersServiceClient is the client API for UsersService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UsersServiceClient interface {
	GetUser(ctx context.Context, in *GetUserIn, opts ...grpc.CallOption) (*GetUserOut, error)
	GetUserByEmail(ctx context.Context, in *GetUserByEmailIn, opts ...grpc.CallOption) (*GetUserOut, error)
	GetUsers(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetUsersOut, error)
	GetMe(ctx context.Context, in *GetMeIn, opts ...grpc.CallOption) (*GetUserOut, error)
}

type usersServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUsersServiceClient(cc grpc.ClientConnInterface) UsersServiceClient {
	return &usersServiceClient{cc}
}

func (c *usersServiceClient) GetUser(ctx context.Context, in *GetUserIn, opts ...grpc.CallOption) (*GetUserOut, error) {
	out := new(GetUserOut)
	err := c.cc.Invoke(ctx, "/users.UsersService/GetUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersServiceClient) GetUserByEmail(ctx context.Context, in *GetUserByEmailIn, opts ...grpc.CallOption) (*GetUserOut, error) {
	out := new(GetUserOut)
	err := c.cc.Invoke(ctx, "/users.UsersService/GetUserByEmail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersServiceClient) GetUsers(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetUsersOut, error) {
	out := new(GetUsersOut)
	err := c.cc.Invoke(ctx, "/users.UsersService/GetUsers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersServiceClient) GetMe(ctx context.Context, in *GetMeIn, opts ...grpc.CallOption) (*GetUserOut, error) {
	out := new(GetUserOut)
	err := c.cc.Invoke(ctx, "/users.UsersService/GetMe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UsersServiceServer is the server API for UsersService service.
// All implementations must embed UnimplementedUsersServiceServer
// for forward compatibility
type UsersServiceServer interface {
	GetUser(context.Context, *GetUserIn) (*GetUserOut, error)
	GetUserByEmail(context.Context, *GetUserByEmailIn) (*GetUserOut, error)
	GetUsers(context.Context, *emptypb.Empty) (*GetUsersOut, error)
	GetMe(context.Context, *GetMeIn) (*GetUserOut, error)
	mustEmbedUnimplementedUsersServiceServer()
}

// UnimplementedUsersServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUsersServiceServer struct {
}

func (UnimplementedUsersServiceServer) GetUser(context.Context, *GetUserIn) (*GetUserOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (UnimplementedUsersServiceServer) GetUserByEmail(context.Context, *GetUserByEmailIn) (*GetUserOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByEmail not implemented")
}
func (UnimplementedUsersServiceServer) GetUsers(context.Context, *emptypb.Empty) (*GetUsersOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUsers not implemented")
}
func (UnimplementedUsersServiceServer) GetMe(context.Context, *GetMeIn) (*GetUserOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMe not implemented")
}
func (UnimplementedUsersServiceServer) mustEmbedUnimplementedUsersServiceServer() {}

// UnsafeUsersServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UsersServiceServer will
// result in compilation errors.
type UnsafeUsersServiceServer interface {
	mustEmbedUnimplementedUsersServiceServer()
}

func RegisterUsersServiceServer(s grpc.ServiceRegistrar, srv UsersServiceServer) {
	s.RegisterService(&UsersService_ServiceDesc, srv)
}

func _UsersService_GetUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).GetUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.UsersService/GetUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).GetUser(ctx, req.(*GetUserIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _UsersService_GetUserByEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserByEmailIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).GetUserByEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.UsersService/GetUserByEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).GetUserByEmail(ctx, req.(*GetUserByEmailIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _UsersService_GetUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).GetUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.UsersService/GetUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).GetUsers(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _UsersService_GetMe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMeIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServiceServer).GetMe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/users.UsersService/GetMe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServiceServer).GetMe(ctx, req.(*GetMeIn))
	}
	return interceptor(ctx, in, info, handler)
}

// UsersService_ServiceDesc is the grpc.ServiceDesc for UsersService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UsersService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "users.UsersService",
	HandlerType: (*UsersServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUser",
			Handler:    _UsersService_GetUser_Handler,
		},
		{
			MethodName: "GetUserByEmail",
			Handler:    _UsersService_GetUserByEmail_Handler,
		},
		{
			MethodName: "GetUsers",
			Handler:    _UsersService_GetUsers_Handler,
		},
		{
			MethodName: "GetMe",
			Handler:    _UsersService_GetMe_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sso/users.proto",
}
