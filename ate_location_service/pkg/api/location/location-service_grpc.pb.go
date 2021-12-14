// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package location

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

// LocationServiceClient is the client API for LocationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LocationServiceClient interface {
	AddLocation(ctx context.Context, in *AddLocationRequest, opts ...grpc.CallOption) (*AddLocationResponse, error)
	UpdateLocation(ctx context.Context, in *UpdateLocationRequest, opts ...grpc.CallOption) (*UpdateLocationResponse, error)
	ViewLocations(ctx context.Context, in *VoidNoParams, opts ...grpc.CallOption) (*AllLocationResponse, error)
	UpdateCurrentLocation(ctx context.Context, in *CurrentLocationRequest, opts ...grpc.CallOption) (*CurrentLocationResponse, error)
	GetCurrentLocation(ctx context.Context, in *VoidNoParams, opts ...grpc.CallOption) (*CurrentLocationResponse, error)
}

type locationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLocationServiceClient(cc grpc.ClientConnInterface) LocationServiceClient {
	return &locationServiceClient{cc}
}

func (c *locationServiceClient) AddLocation(ctx context.Context, in *AddLocationRequest, opts ...grpc.CallOption) (*AddLocationResponse, error) {
	out := new(AddLocationResponse)
	err := c.cc.Invoke(ctx, "/location.LocationService/AddLocation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *locationServiceClient) UpdateLocation(ctx context.Context, in *UpdateLocationRequest, opts ...grpc.CallOption) (*UpdateLocationResponse, error) {
	out := new(UpdateLocationResponse)
	err := c.cc.Invoke(ctx, "/location.LocationService/UpdateLocation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *locationServiceClient) ViewLocations(ctx context.Context, in *VoidNoParams, opts ...grpc.CallOption) (*AllLocationResponse, error) {
	out := new(AllLocationResponse)
	err := c.cc.Invoke(ctx, "/location.LocationService/ViewLocations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *locationServiceClient) UpdateCurrentLocation(ctx context.Context, in *CurrentLocationRequest, opts ...grpc.CallOption) (*CurrentLocationResponse, error) {
	out := new(CurrentLocationResponse)
	err := c.cc.Invoke(ctx, "/location.LocationService/UpdateCurrentLocation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *locationServiceClient) GetCurrentLocation(ctx context.Context, in *VoidNoParams, opts ...grpc.CallOption) (*CurrentLocationResponse, error) {
	out := new(CurrentLocationResponse)
	err := c.cc.Invoke(ctx, "/location.LocationService/GetCurrentLocation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LocationServiceServer is the server API for LocationService service.
// All implementations must embed UnimplementedLocationServiceServer
// for forward compatibility
type LocationServiceServer interface {
	AddLocation(context.Context, *AddLocationRequest) (*AddLocationResponse, error)
	UpdateLocation(context.Context, *UpdateLocationRequest) (*UpdateLocationResponse, error)
	ViewLocations(context.Context, *VoidNoParams) (*AllLocationResponse, error)
	UpdateCurrentLocation(context.Context, *CurrentLocationRequest) (*CurrentLocationResponse, error)
	GetCurrentLocation(context.Context, *VoidNoParams) (*CurrentLocationResponse, error)
	mustEmbedUnimplementedLocationServiceServer()
}

// UnimplementedLocationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLocationServiceServer struct {
}

func (UnimplementedLocationServiceServer) AddLocation(context.Context, *AddLocationRequest) (*AddLocationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddLocation not implemented")
}
func (UnimplementedLocationServiceServer) UpdateLocation(context.Context, *UpdateLocationRequest) (*UpdateLocationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateLocation not implemented")
}
func (UnimplementedLocationServiceServer) ViewLocations(context.Context, *VoidNoParams) (*AllLocationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ViewLocations not implemented")
}
func (UnimplementedLocationServiceServer) UpdateCurrentLocation(context.Context, *CurrentLocationRequest) (*CurrentLocationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCurrentLocation not implemented")
}
func (UnimplementedLocationServiceServer) GetCurrentLocation(context.Context, *VoidNoParams) (*CurrentLocationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCurrentLocation not implemented")
}
func (UnimplementedLocationServiceServer) mustEmbedUnimplementedLocationServiceServer() {}

// UnsafeLocationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LocationServiceServer will
// result in compilation errors.
type UnsafeLocationServiceServer interface {
	mustEmbedUnimplementedLocationServiceServer()
}

func RegisterLocationServiceServer(s grpc.ServiceRegistrar, srv LocationServiceServer) {
	s.RegisterService(&LocationService_ServiceDesc, srv)
}

func _LocationService_AddLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddLocationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocationServiceServer).AddLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/location.LocationService/AddLocation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocationServiceServer).AddLocation(ctx, req.(*AddLocationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocationService_UpdateLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateLocationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocationServiceServer).UpdateLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/location.LocationService/UpdateLocation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocationServiceServer).UpdateLocation(ctx, req.(*UpdateLocationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocationService_ViewLocations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VoidNoParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocationServiceServer).ViewLocations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/location.LocationService/ViewLocations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocationServiceServer).ViewLocations(ctx, req.(*VoidNoParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocationService_UpdateCurrentLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CurrentLocationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocationServiceServer).UpdateCurrentLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/location.LocationService/UpdateCurrentLocation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocationServiceServer).UpdateCurrentLocation(ctx, req.(*CurrentLocationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LocationService_GetCurrentLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VoidNoParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocationServiceServer).GetCurrentLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/location.LocationService/GetCurrentLocation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocationServiceServer).GetCurrentLocation(ctx, req.(*VoidNoParams))
	}
	return interceptor(ctx, in, info, handler)
}

// LocationService_ServiceDesc is the grpc.ServiceDesc for LocationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LocationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "location.LocationService",
	HandlerType: (*LocationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddLocation",
			Handler:    _LocationService_AddLocation_Handler,
		},
		{
			MethodName: "UpdateLocation",
			Handler:    _LocationService_UpdateLocation_Handler,
		},
		{
			MethodName: "ViewLocations",
			Handler:    _LocationService_ViewLocations_Handler,
		},
		{
			MethodName: "UpdateCurrentLocation",
			Handler:    _LocationService_UpdateCurrentLocation_Handler,
		},
		{
			MethodName: "GetCurrentLocation",
			Handler:    _LocationService_GetCurrentLocation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "location-service.proto",
}