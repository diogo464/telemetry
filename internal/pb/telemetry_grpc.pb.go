// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.6
// source: internal/pb/telemetry.proto

package pb

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

// TelemetryClient is the client API for Telemetry service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TelemetryClient interface {
	GetSession(ctx context.Context, in *GetSessionRequest, opts ...grpc.CallOption) (*GetSessionResponse, error)
	GetProperties(ctx context.Context, in *GetPropertiesRequest, opts ...grpc.CallOption) (Telemetry_GetPropertiesClient, error)
	GetStreamDescriptors(ctx context.Context, in *GetStreamDescriptorsRequest, opts ...grpc.CallOption) (*GetStreamDescriptorsResponse, error)
	GetStream(ctx context.Context, in *GetStreamRequest, opts ...grpc.CallOption) (Telemetry_GetStreamClient, error)
}

type telemetryClient struct {
	cc grpc.ClientConnInterface
}

func NewTelemetryClient(cc grpc.ClientConnInterface) TelemetryClient {
	return &telemetryClient{cc}
}

func (c *telemetryClient) GetSession(ctx context.Context, in *GetSessionRequest, opts ...grpc.CallOption) (*GetSessionResponse, error) {
	out := new(GetSessionResponse)
	err := c.cc.Invoke(ctx, "/telemetry.Telemetry/GetSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *telemetryClient) GetProperties(ctx context.Context, in *GetPropertiesRequest, opts ...grpc.CallOption) (Telemetry_GetPropertiesClient, error) {
	stream, err := c.cc.NewStream(ctx, &Telemetry_ServiceDesc.Streams[0], "/telemetry.Telemetry/GetProperties", opts...)
	if err != nil {
		return nil, err
	}
	x := &telemetryGetPropertiesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Telemetry_GetPropertiesClient interface {
	Recv() (*Property, error)
	grpc.ClientStream
}

type telemetryGetPropertiesClient struct {
	grpc.ClientStream
}

func (x *telemetryGetPropertiesClient) Recv() (*Property, error) {
	m := new(Property)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *telemetryClient) GetStreamDescriptors(ctx context.Context, in *GetStreamDescriptorsRequest, opts ...grpc.CallOption) (*GetStreamDescriptorsResponse, error) {
	out := new(GetStreamDescriptorsResponse)
	err := c.cc.Invoke(ctx, "/telemetry.Telemetry/GetStreamDescriptors", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *telemetryClient) GetStream(ctx context.Context, in *GetStreamRequest, opts ...grpc.CallOption) (Telemetry_GetStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Telemetry_ServiceDesc.Streams[1], "/telemetry.Telemetry/GetStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &telemetryGetStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Telemetry_GetStreamClient interface {
	Recv() (*StreamSegment, error)
	grpc.ClientStream
}

type telemetryGetStreamClient struct {
	grpc.ClientStream
}

func (x *telemetryGetStreamClient) Recv() (*StreamSegment, error) {
	m := new(StreamSegment)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TelemetryServer is the server API for Telemetry service.
// All implementations must embed UnimplementedTelemetryServer
// for forward compatibility
type TelemetryServer interface {
	GetSession(context.Context, *GetSessionRequest) (*GetSessionResponse, error)
	GetProperties(*GetPropertiesRequest, Telemetry_GetPropertiesServer) error
	GetStreamDescriptors(context.Context, *GetStreamDescriptorsRequest) (*GetStreamDescriptorsResponse, error)
	GetStream(*GetStreamRequest, Telemetry_GetStreamServer) error
	mustEmbedUnimplementedTelemetryServer()
}

// UnimplementedTelemetryServer must be embedded to have forward compatible implementations.
type UnimplementedTelemetryServer struct {
}

func (UnimplementedTelemetryServer) GetSession(context.Context, *GetSessionRequest) (*GetSessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSession not implemented")
}
func (UnimplementedTelemetryServer) GetProperties(*GetPropertiesRequest, Telemetry_GetPropertiesServer) error {
	return status.Errorf(codes.Unimplemented, "method GetProperties not implemented")
}
func (UnimplementedTelemetryServer) GetStreamDescriptors(context.Context, *GetStreamDescriptorsRequest) (*GetStreamDescriptorsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStreamDescriptors not implemented")
}
func (UnimplementedTelemetryServer) GetStream(*GetStreamRequest, Telemetry_GetStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method GetStream not implemented")
}
func (UnimplementedTelemetryServer) mustEmbedUnimplementedTelemetryServer() {}

// UnsafeTelemetryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TelemetryServer will
// result in compilation errors.
type UnsafeTelemetryServer interface {
	mustEmbedUnimplementedTelemetryServer()
}

func RegisterTelemetryServer(s grpc.ServiceRegistrar, srv TelemetryServer) {
	s.RegisterService(&Telemetry_ServiceDesc, srv)
}

func _Telemetry_GetSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSessionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TelemetryServer).GetSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/telemetry.Telemetry/GetSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TelemetryServer).GetSession(ctx, req.(*GetSessionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Telemetry_GetProperties_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetPropertiesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TelemetryServer).GetProperties(m, &telemetryGetPropertiesServer{stream})
}

type Telemetry_GetPropertiesServer interface {
	Send(*Property) error
	grpc.ServerStream
}

type telemetryGetPropertiesServer struct {
	grpc.ServerStream
}

func (x *telemetryGetPropertiesServer) Send(m *Property) error {
	return x.ServerStream.SendMsg(m)
}

func _Telemetry_GetStreamDescriptors_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStreamDescriptorsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TelemetryServer).GetStreamDescriptors(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/telemetry.Telemetry/GetStreamDescriptors",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TelemetryServer).GetStreamDescriptors(ctx, req.(*GetStreamDescriptorsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Telemetry_GetStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TelemetryServer).GetStream(m, &telemetryGetStreamServer{stream})
}

type Telemetry_GetStreamServer interface {
	Send(*StreamSegment) error
	grpc.ServerStream
}

type telemetryGetStreamServer struct {
	grpc.ServerStream
}

func (x *telemetryGetStreamServer) Send(m *StreamSegment) error {
	return x.ServerStream.SendMsg(m)
}

// Telemetry_ServiceDesc is the grpc.ServiceDesc for Telemetry service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Telemetry_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "telemetry.Telemetry",
	HandlerType: (*TelemetryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetSession",
			Handler:    _Telemetry_GetSession_Handler,
		},
		{
			MethodName: "GetStreamDescriptors",
			Handler:    _Telemetry_GetStreamDescriptors_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetProperties",
			Handler:       _Telemetry_GetProperties_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetStream",
			Handler:       _Telemetry_GetStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "internal/pb/telemetry.proto",
}
