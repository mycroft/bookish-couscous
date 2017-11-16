// Code generated by protoc-gen-go. DO NOT EDIT.
// source: client.proto

/*
Package main is a generated protocol buffer package.

It is generated from these files:
	client.proto

It has these top-level messages:
	HelloRequest
	HelloReply
*/
package main

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The request message containing the user's name.
type HelloRequest struct {
	Uid uint32 `protobuf:"varint,1,opt,name=uid" json:"uid,omitempty"`
}

func (m *HelloRequest) Reset()                    { *m = HelloRequest{} }
func (m *HelloRequest) String() string            { return proto.CompactTextString(m) }
func (*HelloRequest) ProtoMessage()               {}
func (*HelloRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *HelloRequest) GetUid() uint32 {
	if m != nil {
		return m.Uid
	}
	return 0
}

// The response message containing the greetings
type HelloReply struct {
	BestFriend        uint32 `protobuf:"varint,1,opt,name=best_friend,json=bestFriend" json:"best_friend,omitempty"`
	Crush             uint32 `protobuf:"varint,2,opt,name=crush" json:"crush,omitempty"`
	MostSeen          uint32 `protobuf:"varint,3,opt,name=most_seen,json=mostSeen" json:"most_seen,omitempty"`
	MutualLove        uint32 `protobuf:"varint,4,opt,name=mutual_love,json=mutualLove" json:"mutual_love,omitempty"`
	MutualLoveAllTime uint32 `protobuf:"varint,5,opt,name=mutual_love_all_time,json=mutualLoveAllTime" json:"mutual_love_all_time,omitempty"`
}

func (m *HelloReply) Reset()                    { *m = HelloReply{} }
func (m *HelloReply) String() string            { return proto.CompactTextString(m) }
func (*HelloReply) ProtoMessage()               {}
func (*HelloReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *HelloReply) GetBestFriend() uint32 {
	if m != nil {
		return m.BestFriend
	}
	return 0
}

func (m *HelloReply) GetCrush() uint32 {
	if m != nil {
		return m.Crush
	}
	return 0
}

func (m *HelloReply) GetMostSeen() uint32 {
	if m != nil {
		return m.MostSeen
	}
	return 0
}

func (m *HelloReply) GetMutualLove() uint32 {
	if m != nil {
		return m.MutualLove
	}
	return 0
}

func (m *HelloReply) GetMutualLoveAllTime() uint32 {
	if m != nil {
		return m.MutualLoveAllTime
	}
	return 0
}

func init() {
	proto.RegisterType((*HelloRequest)(nil), "main.HelloRequest")
	proto.RegisterType((*HelloReply)(nil), "main.HelloReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Greeter service

type GreeterClient interface {
	SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
}

type greeterClient struct {
	cc *grpc.ClientConn
}

func NewGreeterClient(cc *grpc.ClientConn) GreeterClient {
	return &greeterClient{cc}
}

func (c *greeterClient) SayHello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := grpc.Invoke(ctx, "/main.Greeter/SayHello", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Greeter service

type GreeterServer interface {
	SayHello(context.Context, *HelloRequest) (*HelloReply, error)
}

func RegisterGreeterServer(s *grpc.Server, srv GreeterServer) {
	s.RegisterService(&_Greeter_serviceDesc, srv)
}

func _Greeter_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.Greeter/SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServer).SayHello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Greeter_serviceDesc = grpc.ServiceDesc{
	ServiceName: "main.Greeter",
	HandlerType: (*GreeterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _Greeter_SayHello_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "client.proto",
}

func init() { proto.RegisterFile("client.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 253 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0x4f, 0x4b, 0x33, 0x31,
	0x10, 0xc6, 0xdf, 0x7d, 0xdb, 0x6a, 0x9d, 0xad, 0x50, 0x43, 0x0f, 0x41, 0x05, 0xcb, 0x9e, 0x3c,
	0xad, 0x50, 0xcf, 0x1e, 0xac, 0xa0, 0x1e, 0x3c, 0x94, 0xd6, 0x7b, 0xc8, 0xd6, 0x91, 0x06, 0x26,
	0xc9, 0x9a, 0x3f, 0x85, 0xfd, 0x5a, 0x7e, 0x42, 0x49, 0x56, 0xb1, 0x9e, 0x92, 0xf9, 0xcd, 0x8f,
	0x81, 0xe7, 0x81, 0xc9, 0x96, 0x14, 0x9a, 0x50, 0xb7, 0xce, 0x06, 0xcb, 0x86, 0x5a, 0x2a, 0x53,
	0xcd, 0x61, 0xf2, 0x8c, 0x44, 0x76, 0x8d, 0x1f, 0x11, 0x7d, 0x60, 0x53, 0x18, 0x44, 0xf5, 0xc6,
	0x8b, 0x79, 0x71, 0x7d, 0xba, 0x4e, 0xdf, 0xea, 0xb3, 0x00, 0xf8, 0x56, 0x5a, 0xea, 0xd8, 0x15,
	0x94, 0x0d, 0xfa, 0x20, 0xde, 0x9d, 0x42, 0xf3, 0x23, 0x42, 0x42, 0x8f, 0x99, 0xb0, 0x19, 0x8c,
	0xb6, 0x2e, 0xfa, 0x1d, 0xff, 0x9f, 0x57, 0xfd, 0xc0, 0x2e, 0xe0, 0x44, 0x5b, 0x1f, 0x84, 0x47,
	0x34, 0x7c, 0x90, 0x37, 0xe3, 0x04, 0x36, 0x88, 0x26, 0xdd, 0xd4, 0x31, 0x44, 0x49, 0x82, 0xec,
	0x1e, 0xf9, 0xb0, 0xbf, 0xd9, 0xa3, 0x17, 0xbb, 0x47, 0x76, 0x03, 0xb3, 0x03, 0x41, 0x48, 0x22,
	0x11, 0x94, 0x46, 0x3e, 0xca, 0xe6, 0xd9, 0xaf, 0x79, 0x4f, 0xf4, 0xaa, 0x34, 0x2e, 0xee, 0xe0,
	0xf8, 0xc9, 0x21, 0x06, 0x74, 0x6c, 0x01, 0xe3, 0x8d, 0xec, 0x72, 0x02, 0xc6, 0xea, 0x14, 0xba,
	0x3e, 0x4c, 0x7c, 0x3e, 0xfd, 0xc3, 0x5a, 0xea, 0xaa, 0x7f, 0xcb, 0x4b, 0x28, 0x1b, 0x92, 0xbb,
	0xba, 0x2f, 0x6c, 0x59, 0x3e, 0xe4, 0x77, 0x95, 0x7a, 0x5b, 0x15, 0xcd, 0x51, 0x2e, 0xf0, 0xf6,
	0x2b, 0x00, 0x00, 0xff, 0xff, 0x68, 0xdf, 0x3a, 0xa1, 0x50, 0x01, 0x00, 0x00,
}
