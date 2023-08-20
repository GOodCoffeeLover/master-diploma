// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: api/sandbox.proto

package sandbox

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

// SandboxClient is the client API for Sandbox service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SandboxClient interface {
	// Sends a greeting
	Execute(ctx context.Context, opts ...grpc.CallOption) (Sandbox_ExecuteClient, error)
}

type sandboxClient struct {
	cc grpc.ClientConnInterface
}

func NewSandboxClient(cc grpc.ClientConnInterface) SandboxClient {
	return &sandboxClient{cc}
}

func (c *sandboxClient) Execute(ctx context.Context, opts ...grpc.CallOption) (Sandbox_ExecuteClient, error) {
	stream, err := c.cc.NewStream(ctx, &Sandbox_ServiceDesc.Streams[0], "/sandbox.sandbox/Execute", opts...)
	if err != nil {
		return nil, err
	}
	x := &sandboxExecuteClient{stream}
	return x, nil
}

type Sandbox_ExecuteClient interface {
	Send(*ExecuteRequest) error
	Recv() (*ExecuteResponse, error)
	grpc.ClientStream
}

type sandboxExecuteClient struct {
	grpc.ClientStream
}

func (x *sandboxExecuteClient) Send(m *ExecuteRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *sandboxExecuteClient) Recv() (*ExecuteResponse, error) {
	m := new(ExecuteResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SandboxServer is the server API for Sandbox service.
// All implementations must embed UnimplementedSandboxServer
// for forward compatibility
type SandboxServer interface {
	// Sends a greeting
	Execute(Sandbox_ExecuteServer) error
	mustEmbedUnimplementedSandboxServer()
}

// UnimplementedSandboxServer must be embedded to have forward compatible implementations.
type UnimplementedSandboxServer struct {
}

func (UnimplementedSandboxServer) Execute(Sandbox_ExecuteServer) error {
	return status.Errorf(codes.Unimplemented, "method Execute not implemented")
}
func (UnimplementedSandboxServer) mustEmbedUnimplementedSandboxServer() {}

// UnsafeSandboxServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SandboxServer will
// result in compilation errors.
type UnsafeSandboxServer interface {
	mustEmbedUnimplementedSandboxServer()
}

func RegisterSandboxServer(s grpc.ServiceRegistrar, srv SandboxServer) {
	s.RegisterService(&Sandbox_ServiceDesc, srv)
}

func _Sandbox_Execute_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SandboxServer).Execute(&sandboxExecuteServer{stream})
}

type Sandbox_ExecuteServer interface {
	Send(*ExecuteResponse) error
	Recv() (*ExecuteRequest, error)
	grpc.ServerStream
}

type sandboxExecuteServer struct {
	grpc.ServerStream
}

func (x *sandboxExecuteServer) Send(m *ExecuteResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *sandboxExecuteServer) Recv() (*ExecuteRequest, error) {
	m := new(ExecuteRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Sandbox_ServiceDesc is the grpc.ServiceDesc for Sandbox service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Sandbox_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sandbox.sandbox",
	HandlerType: (*SandboxServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Execute",
			Handler:       _Sandbox_Execute_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "api/sandbox.proto",
}
