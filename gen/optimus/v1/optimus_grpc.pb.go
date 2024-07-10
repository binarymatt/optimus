// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: optimus/v1/optimus.proto

package optimusv1

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

// OptimusLogServiceClient is the client API for OptimusLogService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OptimusLogServiceClient interface {
	// Store saves an event(s) onto the processing pipeline
	StoreLogEvent(ctx context.Context, in *StoreLogEventRequest, opts ...grpc.CallOption) (*StoreLogEventResponse, error)
}

type optimusLogServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOptimusLogServiceClient(cc grpc.ClientConnInterface) OptimusLogServiceClient {
	return &optimusLogServiceClient{cc}
}

func (c *optimusLogServiceClient) StoreLogEvent(ctx context.Context, in *StoreLogEventRequest, opts ...grpc.CallOption) (*StoreLogEventResponse, error) {
	out := new(StoreLogEventResponse)
	err := c.cc.Invoke(ctx, "/optimus.v1.OptimusLogService/StoreLogEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OptimusLogServiceServer is the server API for OptimusLogService service.
// All implementations should embed UnimplementedOptimusLogServiceServer
// for forward compatibility
type OptimusLogServiceServer interface {
	// Store saves an event(s) onto the processing pipeline
	StoreLogEvent(context.Context, *StoreLogEventRequest) (*StoreLogEventResponse, error)
}

// UnimplementedOptimusLogServiceServer should be embedded to have forward compatible implementations.
type UnimplementedOptimusLogServiceServer struct {
}

func (UnimplementedOptimusLogServiceServer) StoreLogEvent(context.Context, *StoreLogEventRequest) (*StoreLogEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreLogEvent not implemented")
}

// UnsafeOptimusLogServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OptimusLogServiceServer will
// result in compilation errors.
type UnsafeOptimusLogServiceServer interface {
	mustEmbedUnimplementedOptimusLogServiceServer()
}

func RegisterOptimusLogServiceServer(s grpc.ServiceRegistrar, srv OptimusLogServiceServer) {
	s.RegisterService(&OptimusLogService_ServiceDesc, srv)
}

func _OptimusLogService_StoreLogEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreLogEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OptimusLogServiceServer).StoreLogEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/optimus.v1.OptimusLogService/StoreLogEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OptimusLogServiceServer).StoreLogEvent(ctx, req.(*StoreLogEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OptimusLogService_ServiceDesc is the grpc.ServiceDesc for OptimusLogService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OptimusLogService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "optimus.v1.OptimusLogService",
	HandlerType: (*OptimusLogServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StoreLogEvent",
			Handler:    _OptimusLogService_StoreLogEvent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "optimus/v1/optimus.proto",
}
