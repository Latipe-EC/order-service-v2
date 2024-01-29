// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.2
// source: internal/infrastructure/grpc/productServ/product_service.proto

package productgrpc

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

// ProductServiceGRPCClient is the client API for ProductServiceGRPC service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProductServiceGRPCClient interface {
	CheckInStock(ctx context.Context, in *GetPurchaseProductRequest, opts ...grpc.CallOption) (*GetPurchaseItemResponse, error)
	UpdateQuantity(ctx context.Context, in *UpdateProductQuantityRequest, opts ...grpc.CallOption) (*UpdateProductQuantityResponse, error)
}

type productServiceGRPCClient struct {
	cc grpc.ClientConnInterface
}

func NewProductServiceGRPCClient(cc grpc.ClientConnInterface) ProductServiceGRPCClient {
	return &productServiceGRPCClient{cc}
}

func (c *productServiceGRPCClient) CheckInStock(ctx context.Context, in *GetPurchaseProductRequest, opts ...grpc.CallOption) (*GetPurchaseItemResponse, error) {
	out := new(GetPurchaseItemResponse)
	err := c.cc.Invoke(ctx, "/protobuf.ProductServiceGRPC/CheckInStock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productServiceGRPCClient) UpdateQuantity(ctx context.Context, in *UpdateProductQuantityRequest, opts ...grpc.CallOption) (*UpdateProductQuantityResponse, error) {
	out := new(UpdateProductQuantityResponse)
	err := c.cc.Invoke(ctx, "/protobuf.ProductServiceGRPC/UpdateQuantity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProductServiceGRPCServer is the server API for ProductServiceGRPC service.
// All implementations must embed UnimplementedProductServiceGRPCServer
// for forward compatibility
type ProductServiceGRPCServer interface {
	CheckInStock(context.Context, *GetPurchaseProductRequest) (*GetPurchaseItemResponse, error)
	UpdateQuantity(context.Context, *UpdateProductQuantityRequest) (*UpdateProductQuantityResponse, error)
	mustEmbedUnimplementedProductServiceGRPCServer()
}

// UnimplementedProductServiceGRPCServer must be embedded to have forward compatible implementations.
type UnimplementedProductServiceGRPCServer struct {
}

func (UnimplementedProductServiceGRPCServer) CheckInStock(context.Context, *GetPurchaseProductRequest) (*GetPurchaseItemResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckInStock not implemented")
}
func (UnimplementedProductServiceGRPCServer) UpdateQuantity(context.Context, *UpdateProductQuantityRequest) (*UpdateProductQuantityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateQuantity not implemented")
}
func (UnimplementedProductServiceGRPCServer) mustEmbedUnimplementedProductServiceGRPCServer() {}

// UnsafeProductServiceGRPCServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProductServiceGRPCServer will
// result in compilation errors.
type UnsafeProductServiceGRPCServer interface {
	mustEmbedUnimplementedProductServiceGRPCServer()
}

func RegisterProductServiceGRPCServer(s grpc.ServiceRegistrar, srv ProductServiceGRPCServer) {
	s.RegisterService(&ProductServiceGRPC_ServiceDesc, srv)
}

func _ProductServiceGRPC_CheckInStock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPurchaseProductRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceGRPCServer).CheckInStock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.ProductServiceGRPC/CheckInStock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceGRPCServer).CheckInStock(ctx, req.(*GetPurchaseProductRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductServiceGRPC_UpdateQuantity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProductQuantityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceGRPCServer).UpdateQuantity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.ProductServiceGRPC/UpdateQuantity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceGRPCServer).UpdateQuantity(ctx, req.(*UpdateProductQuantityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProductServiceGRPC_ServiceDesc is the grpc.ServiceDesc for ProductServiceGRPC service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProductServiceGRPC_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.ProductServiceGRPC",
	HandlerType: (*ProductServiceGRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckInStock",
			Handler:    _ProductServiceGRPC_CheckInStock_Handler,
		},
		{
			MethodName: "UpdateQuantity",
			Handler:    _ProductServiceGRPC_UpdateQuantity_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/infrastructure/grpc/productServ/product_service.proto",
}