// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package onboard

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

// OnboardServiceClient is the client API for OnboardService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OnboardServiceClient interface {
	Login(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*LoginUserResponse, error)
	Register(ctx context.Context, in *RegisterUserRequest, opts ...grpc.CallOption) (*RegisterUserResponse, error)
	NotificationRegister(ctx context.Context, in *NotificationUserRequest, opts ...grpc.CallOption) (*NotificationUserResponse, error)
}

type onboardServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOnboardServiceClient(cc grpc.ClientConnInterface) OnboardServiceClient {
	return &onboardServiceClient{cc}
}

func (c *onboardServiceClient) Login(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*LoginUserResponse, error) {
	out := new(LoginUserResponse)
	err := c.cc.Invoke(ctx, "/onboard.OnboardService/login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *onboardServiceClient) Register(ctx context.Context, in *RegisterUserRequest, opts ...grpc.CallOption) (*RegisterUserResponse, error) {
	out := new(RegisterUserResponse)
	err := c.cc.Invoke(ctx, "/onboard.OnboardService/register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *onboardServiceClient) NotificationRegister(ctx context.Context, in *NotificationUserRequest, opts ...grpc.CallOption) (*NotificationUserResponse, error) {
	out := new(NotificationUserResponse)
	err := c.cc.Invoke(ctx, "/onboard.OnboardService/notificationRegister", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OnboardServiceServer is the server API for OnboardService service.
// All implementations must embed UnimplementedOnboardServiceServer
// for forward compatibility
type OnboardServiceServer interface {
	Login(context.Context, *LoginUserRequest) (*LoginUserResponse, error)
	Register(context.Context, *RegisterUserRequest) (*RegisterUserResponse, error)
	NotificationRegister(context.Context, *NotificationUserRequest) (*NotificationUserResponse, error)
}

// UnimplementedOnboardServiceServer must be embedded to have forward compatible implementations.
type UnimplementedOnboardServiceServer struct {
}

func (UnimplementedOnboardServiceServer) Login(context.Context, *LoginUserRequest) (*LoginUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedOnboardServiceServer) Register(context.Context, *RegisterUserRequest) (*RegisterUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedOnboardServiceServer) NotificationRegister(context.Context, *NotificationUserRequest) (*NotificationUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotificationRegister not implemented")
}


func RegisterOnboardServiceServer(s grpc.ServiceRegistrar, srv OnboardServiceServer) {
	s.RegisterService(&OnboardService_ServiceDesc, srv)
}

func _OnboardService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OnboardServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/onboard.OnboardService/login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OnboardServiceServer).Login(ctx, req.(*LoginUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OnboardService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OnboardServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/onboard.OnboardService/register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OnboardServiceServer).Register(ctx, req.(*RegisterUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OnboardService_NotificationRegister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotificationUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OnboardServiceServer).NotificationRegister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/onboard.OnboardService/notificationRegister",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OnboardServiceServer).NotificationRegister(ctx, req.(*NotificationUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OnboardService_ServiceDesc is the grpc.ServiceDesc for OnboardService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OnboardService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "onboard.OnboardService",
	HandlerType: (*OnboardServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "login",
			Handler:    _OnboardService_Login_Handler,
		},
		{
			MethodName: "register",
			Handler:    _OnboardService_Register_Handler,
		},
		{
			MethodName: "notificationRegister",
			Handler:    _OnboardService_NotificationRegister_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "onboard-service.proto",
}
