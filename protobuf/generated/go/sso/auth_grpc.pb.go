// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package sso

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

// SSOServiceClient is the client API for SSOService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SSOServiceClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
}

type sSOServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSSOServiceClient(cc grpc.ClientConnInterface) SSOServiceClient {
	return &sSOServiceClient{cc}
}

func (c *sSOServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/auth.SSOService/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sSOServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/auth.SSOService/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SSOServiceServer is the server API for SSOService service.
// All implementations must embed UnimplementedSSOServiceServer
// for forward compatibility
type SSOServiceServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	mustEmbedUnimplementedSSOServiceServer()
}

// UnimplementedSSOServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSSOServiceServer struct {
}

func (UnimplementedSSOServiceServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedSSOServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedSSOServiceServer) mustEmbedUnimplementedSSOServiceServer() {}

// UnsafeSSOServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SSOServiceServer will
// result in compilation errors.
type UnsafeSSOServiceServer interface {
	mustEmbedUnimplementedSSOServiceServer()
}

func RegisterSSOServiceServer(s grpc.ServiceRegistrar, srv SSOServiceServer) {
	s.RegisterService(&SSOService_ServiceDesc, srv)
}

func _SSOService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SSOServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth.SSOService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SSOServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SSOService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SSOServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth.SSOService/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SSOServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SSOService_ServiceDesc is the grpc.ServiceDesc for SSOService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SSOService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth.SSOService",
	HandlerType: (*SSOServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _SSOService_Register_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _SSOService_Login_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sso/auth.proto",
}
