// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.2
// source: internal/infrastructure/grpc/userServ/user_service.proto

package usergrpc

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

// UserServiceGRPCClient is the client API for UserServiceGRPC service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceGRPCClient interface {
	GetAddressDetail(ctx context.Context, in *GetDetailAddressRequest, opts ...grpc.CallOption) (*GetDetailAddressResponse, error)
}

type userServiceGRPCClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceGRPCClient(cc grpc.ClientConnInterface) UserServiceGRPCClient {
	return &userServiceGRPCClient{cc}
}

func (c *userServiceGRPCClient) GetAddressDetail(ctx context.Context, in *GetDetailAddressRequest, opts ...grpc.CallOption) (*GetDetailAddressResponse, error) {
	out := new(GetDetailAddressResponse)
	err := c.cc.Invoke(ctx, "/UserServiceGRPC/GetAddressDetail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceGRPCServer is the server API for UserServiceGRPC service.
// All implementations must embed UnimplementedUserServiceGRPCServer
// for forward compatibility
type UserServiceGRPCServer interface {
	GetAddressDetail(context.Context, *GetDetailAddressRequest) (*GetDetailAddressResponse, error)
	mustEmbedUnimplementedUserServiceGRPCServer()
}

// UnimplementedUserServiceGRPCServer must be embedded to have forward compatible implementations.
type UnimplementedUserServiceGRPCServer struct {
}

func (UnimplementedUserServiceGRPCServer) GetAddressDetail(context.Context, *GetDetailAddressRequest) (*GetDetailAddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAddressDetail not implemented")
}
func (UnimplementedUserServiceGRPCServer) mustEmbedUnimplementedUserServiceGRPCServer() {}

// UnsafeUserServiceGRPCServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceGRPCServer will
// result in compilation errors.
type UnsafeUserServiceGRPCServer interface {
	mustEmbedUnimplementedUserServiceGRPCServer()
}

func RegisterUserServiceGRPCServer(s grpc.ServiceRegistrar, srv UserServiceGRPCServer) {
	s.RegisterService(&UserServiceGRPC_ServiceDesc, srv)
}

func _UserServiceGRPC_GetAddressDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDetailAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceGRPCServer).GetAddressDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/UserServiceGRPC/GetAddressDetail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceGRPCServer).GetAddressDetail(ctx, req.(*GetDetailAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserServiceGRPC_ServiceDesc is the grpc.ServiceDesc for UserServiceGRPC service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserServiceGRPC_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "UserServiceGRPC",
	HandlerType: (*UserServiceGRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAddressDetail",
			Handler:    _UserServiceGRPC_GetAddressDetail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/infrastructure/grpc/userServ/user_service.proto",
}
