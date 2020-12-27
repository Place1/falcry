// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protos

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// EventsClient is the client API for Events service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventsClient interface {
	Listen(ctx context.Context, in *ListenEventsReq, opts ...grpc.CallOption) (Events_ListenClient, error)
}

type eventsClient struct {
	cc grpc.ClientConnInterface
}

func NewEventsClient(cc grpc.ClientConnInterface) EventsClient {
	return &eventsClient{cc}
}

func (c *eventsClient) Listen(ctx context.Context, in *ListenEventsReq, opts ...grpc.CallOption) (Events_ListenClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Events_serviceDesc.Streams[0], "/protos.Events/Listen", opts...)
	if err != nil {
		return nil, err
	}
	x := &eventsListenClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Events_ListenClient interface {
	Recv() (*ListenEventsRes, error)
	grpc.ClientStream
}

type eventsListenClient struct {
	grpc.ClientStream
}

func (x *eventsListenClient) Recv() (*ListenEventsRes, error) {
	m := new(ListenEventsRes)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// EventsServer is the server API for Events service.
// All implementations must embed UnimplementedEventsServer
// for forward compatibility
type EventsServer interface {
	Listen(*ListenEventsReq, Events_ListenServer) error
	mustEmbedUnimplementedEventsServer()
}

// UnimplementedEventsServer must be embedded to have forward compatible implementations.
type UnimplementedEventsServer struct {
}

func (UnimplementedEventsServer) Listen(*ListenEventsReq, Events_ListenServer) error {
	return status.Errorf(codes.Unimplemented, "method Listen not implemented")
}
func (UnimplementedEventsServer) mustEmbedUnimplementedEventsServer() {}

// UnsafeEventsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventsServer will
// result in compilation errors.
type UnsafeEventsServer interface {
	mustEmbedUnimplementedEventsServer()
}

func RegisterEventsServer(s grpc.ServiceRegistrar, srv EventsServer) {
	s.RegisterService(&_Events_serviceDesc, srv)
}

func _Events_Listen_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListenEventsReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EventsServer).Listen(m, &eventsListenServer{stream})
}

type Events_ListenServer interface {
	Send(*ListenEventsRes) error
	grpc.ServerStream
}

type eventsListenServer struct {
	grpc.ServerStream
}

func (x *eventsListenServer) Send(m *ListenEventsRes) error {
	return x.ServerStream.SendMsg(m)
}

var _Events_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.Events",
	HandlerType: (*EventsServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Listen",
			Handler:       _Events_Listen_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "events.proto",
}
