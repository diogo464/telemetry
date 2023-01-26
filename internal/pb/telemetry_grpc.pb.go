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
	GetPropertyDescriptors(ctx context.Context, in *GetPropertyDescriptorsRequest, opts ...grpc.CallOption) (Telemetry_GetPropertyDescriptorsClient, error)
	GetProperties(ctx context.Context, in *GetPropertiesRequest, opts ...grpc.CallOption) (Telemetry_GetPropertiesClient, error)
	GetMetricDescriptors(ctx context.Context, in *GetMetricDescriptorsRequest, opts ...grpc.CallOption) (Telemetry_GetMetricDescriptorsClient, error)
	GetCaptureDescriptors(ctx context.Context, in *GetCaptureDescriptorsRequest, opts ...grpc.CallOption) (Telemetry_GetCaptureDescriptorsClient, error)
	GetEventDescriptors(ctx context.Context, in *GetEventDescriptorsRequest, opts ...grpc.CallOption) (Telemetry_GetEventDescriptorsClient, error)
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

func (c *telemetryClient) GetPropertyDescriptors(ctx context.Context, in *GetPropertyDescriptorsRequest, opts ...grpc.CallOption) (Telemetry_GetPropertyDescriptorsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Telemetry_ServiceDesc.Streams[0], "/telemetry.Telemetry/GetPropertyDescriptors", opts...)
	if err != nil {
		return nil, err
	}
	x := &telemetryGetPropertyDescriptorsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Telemetry_GetPropertyDescriptorsClient interface {
	Recv() (*PropertyDescriptor, error)
	grpc.ClientStream
}

type telemetryGetPropertyDescriptorsClient struct {
	grpc.ClientStream
}

func (x *telemetryGetPropertyDescriptorsClient) Recv() (*PropertyDescriptor, error) {
	m := new(PropertyDescriptor)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *telemetryClient) GetProperties(ctx context.Context, in *GetPropertiesRequest, opts ...grpc.CallOption) (Telemetry_GetPropertiesClient, error) {
	stream, err := c.cc.NewStream(ctx, &Telemetry_ServiceDesc.Streams[1], "/telemetry.Telemetry/GetProperties", opts...)
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

func (c *telemetryClient) GetMetricDescriptors(ctx context.Context, in *GetMetricDescriptorsRequest, opts ...grpc.CallOption) (Telemetry_GetMetricDescriptorsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Telemetry_ServiceDesc.Streams[2], "/telemetry.Telemetry/GetMetricDescriptors", opts...)
	if err != nil {
		return nil, err
	}
	x := &telemetryGetMetricDescriptorsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Telemetry_GetMetricDescriptorsClient interface {
	Recv() (*MetricDescriptor, error)
	grpc.ClientStream
}

type telemetryGetMetricDescriptorsClient struct {
	grpc.ClientStream
}

func (x *telemetryGetMetricDescriptorsClient) Recv() (*MetricDescriptor, error) {
	m := new(MetricDescriptor)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *telemetryClient) GetCaptureDescriptors(ctx context.Context, in *GetCaptureDescriptorsRequest, opts ...grpc.CallOption) (Telemetry_GetCaptureDescriptorsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Telemetry_ServiceDesc.Streams[3], "/telemetry.Telemetry/GetCaptureDescriptors", opts...)
	if err != nil {
		return nil, err
	}
	x := &telemetryGetCaptureDescriptorsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Telemetry_GetCaptureDescriptorsClient interface {
	Recv() (*CaptureDescriptor, error)
	grpc.ClientStream
}

type telemetryGetCaptureDescriptorsClient struct {
	grpc.ClientStream
}

func (x *telemetryGetCaptureDescriptorsClient) Recv() (*CaptureDescriptor, error) {
	m := new(CaptureDescriptor)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *telemetryClient) GetEventDescriptors(ctx context.Context, in *GetEventDescriptorsRequest, opts ...grpc.CallOption) (Telemetry_GetEventDescriptorsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Telemetry_ServiceDesc.Streams[4], "/telemetry.Telemetry/GetEventDescriptors", opts...)
	if err != nil {
		return nil, err
	}
	x := &telemetryGetEventDescriptorsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Telemetry_GetEventDescriptorsClient interface {
	Recv() (*EventDescriptor, error)
	grpc.ClientStream
}

type telemetryGetEventDescriptorsClient struct {
	grpc.ClientStream
}

func (x *telemetryGetEventDescriptorsClient) Recv() (*EventDescriptor, error) {
	m := new(EventDescriptor)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *telemetryClient) GetStream(ctx context.Context, in *GetStreamRequest, opts ...grpc.CallOption) (Telemetry_GetStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Telemetry_ServiceDesc.Streams[5], "/telemetry.Telemetry/GetStream", opts...)
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
	GetPropertyDescriptors(*GetPropertyDescriptorsRequest, Telemetry_GetPropertyDescriptorsServer) error
	GetProperties(*GetPropertiesRequest, Telemetry_GetPropertiesServer) error
	GetMetricDescriptors(*GetMetricDescriptorsRequest, Telemetry_GetMetricDescriptorsServer) error
	GetCaptureDescriptors(*GetCaptureDescriptorsRequest, Telemetry_GetCaptureDescriptorsServer) error
	GetEventDescriptors(*GetEventDescriptorsRequest, Telemetry_GetEventDescriptorsServer) error
	GetStream(*GetStreamRequest, Telemetry_GetStreamServer) error
	mustEmbedUnimplementedTelemetryServer()
}

// UnimplementedTelemetryServer must be embedded to have forward compatible implementations.
type UnimplementedTelemetryServer struct {
}

func (UnimplementedTelemetryServer) GetSession(context.Context, *GetSessionRequest) (*GetSessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSession not implemented")
}
func (UnimplementedTelemetryServer) GetPropertyDescriptors(*GetPropertyDescriptorsRequest, Telemetry_GetPropertyDescriptorsServer) error {
	return status.Errorf(codes.Unimplemented, "method GetPropertyDescriptors not implemented")
}
func (UnimplementedTelemetryServer) GetProperties(*GetPropertiesRequest, Telemetry_GetPropertiesServer) error {
	return status.Errorf(codes.Unimplemented, "method GetProperties not implemented")
}
func (UnimplementedTelemetryServer) GetMetricDescriptors(*GetMetricDescriptorsRequest, Telemetry_GetMetricDescriptorsServer) error {
	return status.Errorf(codes.Unimplemented, "method GetMetricDescriptors not implemented")
}
func (UnimplementedTelemetryServer) GetCaptureDescriptors(*GetCaptureDescriptorsRequest, Telemetry_GetCaptureDescriptorsServer) error {
	return status.Errorf(codes.Unimplemented, "method GetCaptureDescriptors not implemented")
}
func (UnimplementedTelemetryServer) GetEventDescriptors(*GetEventDescriptorsRequest, Telemetry_GetEventDescriptorsServer) error {
	return status.Errorf(codes.Unimplemented, "method GetEventDescriptors not implemented")
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

func _Telemetry_GetPropertyDescriptors_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetPropertyDescriptorsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TelemetryServer).GetPropertyDescriptors(m, &telemetryGetPropertyDescriptorsServer{stream})
}

type Telemetry_GetPropertyDescriptorsServer interface {
	Send(*PropertyDescriptor) error
	grpc.ServerStream
}

type telemetryGetPropertyDescriptorsServer struct {
	grpc.ServerStream
}

func (x *telemetryGetPropertyDescriptorsServer) Send(m *PropertyDescriptor) error {
	return x.ServerStream.SendMsg(m)
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

func _Telemetry_GetMetricDescriptors_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetMetricDescriptorsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TelemetryServer).GetMetricDescriptors(m, &telemetryGetMetricDescriptorsServer{stream})
}

type Telemetry_GetMetricDescriptorsServer interface {
	Send(*MetricDescriptor) error
	grpc.ServerStream
}

type telemetryGetMetricDescriptorsServer struct {
	grpc.ServerStream
}

func (x *telemetryGetMetricDescriptorsServer) Send(m *MetricDescriptor) error {
	return x.ServerStream.SendMsg(m)
}

func _Telemetry_GetCaptureDescriptors_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetCaptureDescriptorsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TelemetryServer).GetCaptureDescriptors(m, &telemetryGetCaptureDescriptorsServer{stream})
}

type Telemetry_GetCaptureDescriptorsServer interface {
	Send(*CaptureDescriptor) error
	grpc.ServerStream
}

type telemetryGetCaptureDescriptorsServer struct {
	grpc.ServerStream
}

func (x *telemetryGetCaptureDescriptorsServer) Send(m *CaptureDescriptor) error {
	return x.ServerStream.SendMsg(m)
}

func _Telemetry_GetEventDescriptors_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetEventDescriptorsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TelemetryServer).GetEventDescriptors(m, &telemetryGetEventDescriptorsServer{stream})
}

type Telemetry_GetEventDescriptorsServer interface {
	Send(*EventDescriptor) error
	grpc.ServerStream
}

type telemetryGetEventDescriptorsServer struct {
	grpc.ServerStream
}

func (x *telemetryGetEventDescriptorsServer) Send(m *EventDescriptor) error {
	return x.ServerStream.SendMsg(m)
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
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetPropertyDescriptors",
			Handler:       _Telemetry_GetPropertyDescriptors_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetProperties",
			Handler:       _Telemetry_GetProperties_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetMetricDescriptors",
			Handler:       _Telemetry_GetMetricDescriptors_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetCaptureDescriptors",
			Handler:       _Telemetry_GetCaptureDescriptors_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetEventDescriptors",
			Handler:       _Telemetry_GetEventDescriptors_Handler,
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
