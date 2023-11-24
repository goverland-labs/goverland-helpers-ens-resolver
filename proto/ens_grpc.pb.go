// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.0
// source: ens.proto

package proto

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

// EnsClient is the client API for Ens service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EnsClient interface {
	ResolveAddresses(ctx context.Context, in *ResolveAddressesRequest, opts ...grpc.CallOption) (*ResolveResponse, error)
	ResolveDomains(ctx context.Context, in *ResolveDomainsRequest, opts ...grpc.CallOption) (*ResolveResponse, error)
}

type ensClient struct {
	cc grpc.ClientConnInterface
}

func NewEnsClient(cc grpc.ClientConnInterface) EnsClient {
	return &ensClient{cc}
}

func (c *ensClient) ResolveAddresses(ctx context.Context, in *ResolveAddressesRequest, opts ...grpc.CallOption) (*ResolveResponse, error) {
	out := new(ResolveResponse)
	err := c.cc.Invoke(ctx, "/proto.Ens/ResolveAddresses", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ensClient) ResolveDomains(ctx context.Context, in *ResolveDomainsRequest, opts ...grpc.CallOption) (*ResolveResponse, error) {
	out := new(ResolveResponse)
	err := c.cc.Invoke(ctx, "/proto.Ens/ResolveDomains", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EnsServer is the server API for Ens service.
// All implementations must embed UnimplementedEnsServer
// for forward compatibility
type EnsServer interface {
	ResolveAddresses(context.Context, *ResolveAddressesRequest) (*ResolveResponse, error)
	ResolveDomains(context.Context, *ResolveDomainsRequest) (*ResolveResponse, error)
	mustEmbedUnimplementedEnsServer()
}

// UnimplementedEnsServer must be embedded to have forward compatible implementations.
type UnimplementedEnsServer struct {
}

func (UnimplementedEnsServer) ResolveAddresses(context.Context, *ResolveAddressesRequest) (*ResolveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResolveAddresses not implemented")
}
func (UnimplementedEnsServer) ResolveDomains(context.Context, *ResolveDomainsRequest) (*ResolveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResolveDomains not implemented")
}
func (UnimplementedEnsServer) mustEmbedUnimplementedEnsServer() {}

// UnsafeEnsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EnsServer will
// result in compilation errors.
type UnsafeEnsServer interface {
	mustEmbedUnimplementedEnsServer()
}

func RegisterEnsServer(s grpc.ServiceRegistrar, srv EnsServer) {
	s.RegisterService(&Ens_ServiceDesc, srv)
}

func _Ens_ResolveAddresses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResolveAddressesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnsServer).ResolveAddresses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Ens/ResolveAddresses",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnsServer).ResolveAddresses(ctx, req.(*ResolveAddressesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Ens_ResolveDomains_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResolveDomainsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EnsServer).ResolveDomains(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Ens/ResolveDomains",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EnsServer).ResolveDomains(ctx, req.(*ResolveDomainsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Ens_ServiceDesc is the grpc.ServiceDesc for Ens service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Ens_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Ens",
	HandlerType: (*EnsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ResolveAddresses",
			Handler:    _Ens_ResolveAddresses_Handler,
		},
		{
			MethodName: "ResolveDomains",
			Handler:    _Ens_ResolveDomains_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ens.proto",
}
