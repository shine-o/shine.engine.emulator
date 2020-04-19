// Code generated by protoc-gen-go. DO NOT EDIT.
// source: master.proto

package zm

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type MapQuery struct {
	Name                 string   `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MapQuery) Reset()         { *m = MapQuery{} }
func (m *MapQuery) String() string { return proto.CompactTextString(m) }
func (*MapQuery) ProtoMessage()    {}
func (*MapQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_f9c348dec43a6705, []int{0}
}

func (m *MapQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MapQuery.Unmarshal(m, b)
}
func (m *MapQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MapQuery.Marshal(b, m, deterministic)
}
func (m *MapQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MapQuery.Merge(m, src)
}
func (m *MapQuery) XXX_Size() int {
	return xxx_messageInfo_MapQuery.Size(m)
}
func (m *MapQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_MapQuery.DiscardUnknown(m)
}

var xxx_messageInfo_MapQuery proto.InternalMessageInfo

func (m *MapQuery) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type ConnectionInfo struct {
	IP                   string   `protobuf:"bytes,1,opt,name=IP,proto3" json:"IP,omitempty"`
	Port                 int32    `protobuf:"varint,2,opt,name=Port,proto3" json:"Port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConnectionInfo) Reset()         { *m = ConnectionInfo{} }
func (m *ConnectionInfo) String() string { return proto.CompactTextString(m) }
func (*ConnectionInfo) ProtoMessage()    {}
func (*ConnectionInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_f9c348dec43a6705, []int{1}
}

func (m *ConnectionInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConnectionInfo.Unmarshal(m, b)
}
func (m *ConnectionInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConnectionInfo.Marshal(b, m, deterministic)
}
func (m *ConnectionInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConnectionInfo.Merge(m, src)
}
func (m *ConnectionInfo) XXX_Size() int {
	return xxx_messageInfo_ConnectionInfo.Size(m)
}
func (m *ConnectionInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ConnectionInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ConnectionInfo proto.InternalMessageInfo

func (m *ConnectionInfo) GetIP() string {
	if m != nil {
		return m.IP
	}
	return ""
}

func (m *ConnectionInfo) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

type ZoneDetails struct {
	Maps                 []string        `protobuf:"bytes,1,rep,name=maps,proto3" json:"maps,omitempty"`
	Conn                 *ConnectionInfo `protobuf:"bytes,2,opt,name=conn,proto3" json:"conn,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *ZoneDetails) Reset()         { *m = ZoneDetails{} }
func (m *ZoneDetails) String() string { return proto.CompactTextString(m) }
func (*ZoneDetails) ProtoMessage()    {}
func (*ZoneDetails) Descriptor() ([]byte, []int) {
	return fileDescriptor_f9c348dec43a6705, []int{2}
}

func (m *ZoneDetails) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ZoneDetails.Unmarshal(m, b)
}
func (m *ZoneDetails) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ZoneDetails.Marshal(b, m, deterministic)
}
func (m *ZoneDetails) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ZoneDetails.Merge(m, src)
}
func (m *ZoneDetails) XXX_Size() int {
	return xxx_messageInfo_ZoneDetails.Size(m)
}
func (m *ZoneDetails) XXX_DiscardUnknown() {
	xxx_messageInfo_ZoneDetails.DiscardUnknown(m)
}

var xxx_messageInfo_ZoneDetails proto.InternalMessageInfo

func (m *ZoneDetails) GetMaps() []string {
	if m != nil {
		return m.Maps
	}
	return nil
}

func (m *ZoneDetails) GetConn() *ConnectionInfo {
	if m != nil {
		return m.Conn
	}
	return nil
}

type ZoneRegistered struct {
	Success              bool     `protobuf:"varint,1,opt,name=Success,proto3" json:"Success,omitempty"`
	ZoneID               string   `protobuf:"bytes,2,opt,name=ZoneID,proto3" json:"ZoneID,omitempty"`
	World                *World   `protobuf:"bytes,3,opt,name=World,proto3" json:"World,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ZoneRegistered) Reset()         { *m = ZoneRegistered{} }
func (m *ZoneRegistered) String() string { return proto.CompactTextString(m) }
func (*ZoneRegistered) ProtoMessage()    {}
func (*ZoneRegistered) Descriptor() ([]byte, []int) {
	return fileDescriptor_f9c348dec43a6705, []int{3}
}

func (m *ZoneRegistered) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ZoneRegistered.Unmarshal(m, b)
}
func (m *ZoneRegistered) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ZoneRegistered.Marshal(b, m, deterministic)
}
func (m *ZoneRegistered) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ZoneRegistered.Merge(m, src)
}
func (m *ZoneRegistered) XXX_Size() int {
	return xxx_messageInfo_ZoneRegistered.Size(m)
}
func (m *ZoneRegistered) XXX_DiscardUnknown() {
	xxx_messageInfo_ZoneRegistered.DiscardUnknown(m)
}

var xxx_messageInfo_ZoneRegistered proto.InternalMessageInfo

func (m *ZoneRegistered) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *ZoneRegistered) GetZoneID() string {
	if m != nil {
		return m.ZoneID
	}
	return ""
}

func (m *ZoneRegistered) GetWorld() *World {
	if m != nil {
		return m.World
	}
	return nil
}

type World struct {
	IP                   string   `protobuf:"bytes,1,opt,name=IP,proto3" json:"IP,omitempty"`
	Port                 int32    `protobuf:"varint,2,opt,name=Port,proto3" json:"Port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *World) Reset()         { *m = World{} }
func (m *World) String() string { return proto.CompactTextString(m) }
func (*World) ProtoMessage()    {}
func (*World) Descriptor() ([]byte, []int) {
	return fileDescriptor_f9c348dec43a6705, []int{4}
}

func (m *World) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_World.Unmarshal(m, b)
}
func (m *World) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_World.Marshal(b, m, deterministic)
}
func (m *World) XXX_Merge(src proto.Message) {
	xxx_messageInfo_World.Merge(m, src)
}
func (m *World) XXX_Size() int {
	return xxx_messageInfo_World.Size(m)
}
func (m *World) XXX_DiscardUnknown() {
	xxx_messageInfo_World.DiscardUnknown(m)
}

var xxx_messageInfo_World proto.InternalMessageInfo

func (m *World) GetIP() string {
	if m != nil {
		return m.IP
	}
	return ""
}

func (m *World) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func init() {
	proto.RegisterType((*MapQuery)(nil), "zm.MapQuery")
	proto.RegisterType((*ConnectionInfo)(nil), "zm.ConnectionInfo")
	proto.RegisterType((*ZoneDetails)(nil), "zm.ZoneDetails")
	proto.RegisterType((*ZoneRegistered)(nil), "zm.ZoneRegistered")
	proto.RegisterType((*World)(nil), "zm.World")
}

func init() {
	proto.RegisterFile("master.proto", fileDescriptor_f9c348dec43a6705)
}

var fileDescriptor_f9c348dec43a6705 = []byte{
	// 288 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x51, 0x4f, 0x4b, 0xfb, 0x40,
	0x10, 0x6d, 0xd2, 0x3f, 0xbf, 0x66, 0x1a, 0xf2, 0x83, 0x39, 0x48, 0xf0, 0xa0, 0x61, 0x0f, 0x52,
	0x10, 0x82, 0x54, 0xfd, 0x04, 0xf6, 0xb2, 0x87, 0x48, 0x5c, 0x0f, 0x05, 0x6f, 0x6b, 0xba, 0x6a,
	0xa0, 0xd9, 0x8d, 0xbb, 0xe9, 0xc1, 0x7e, 0x7a, 0x99, 0x6d, 0x02, 0x16, 0x3c, 0x78, 0x9b, 0x99,
	0x7d, 0xf3, 0xe6, 0xbd, 0xb7, 0x10, 0x37, 0xd2, 0x75, 0xca, 0xe6, 0xad, 0x35, 0x9d, 0xc1, 0xf0,
	0xd0, 0xb0, 0x0b, 0x98, 0x17, 0xb2, 0x7d, 0xda, 0x2b, 0xfb, 0x85, 0x08, 0x93, 0x47, 0xd9, 0xa8,
	0x34, 0xc8, 0x82, 0x65, 0x24, 0x7c, 0xcd, 0xee, 0x20, 0x79, 0x30, 0x5a, 0xab, 0xaa, 0xab, 0x8d,
	0xe6, 0xfa, 0xcd, 0x60, 0x02, 0x21, 0x2f, 0x7b, 0x4c, 0xc8, 0x4b, 0xda, 0x2a, 0x8d, 0xed, 0xd2,
	0x30, 0x0b, 0x96, 0x53, 0xe1, 0x6b, 0xc6, 0x61, 0xf1, 0x62, 0xb4, 0x5a, 0xab, 0x4e, 0xd6, 0x3b,
	0x47, 0x90, 0x46, 0xb6, 0x2e, 0x0d, 0xb2, 0x31, 0x11, 0x53, 0x8d, 0x57, 0x30, 0xa9, 0x8c, 0xd6,
	0x7e, 0x6d, 0xb1, 0xc2, 0xfc, 0xd0, 0xe4, 0xa7, 0x87, 0x84, 0x7f, 0x67, 0x15, 0x24, 0x44, 0x25,
	0xd4, 0x7b, 0x4d, 0xd2, 0xd5, 0x16, 0x53, 0xf8, 0xf7, 0xbc, 0xaf, 0x2a, 0xe5, 0x9c, 0x57, 0x31,
	0x17, 0x43, 0x8b, 0x67, 0x30, 0x23, 0x2c, 0x5f, 0x7b, 0xd6, 0x48, 0xf4, 0x1d, 0x5e, 0xc2, 0x74,
	0x63, 0xec, 0x6e, 0x9b, 0x8e, 0xfd, 0xb1, 0x88, 0x8e, 0xf9, 0x81, 0x38, 0xce, 0xd9, 0x75, 0x0f,
	0xf8, 0x8b, 0xb9, 0xd5, 0x27, 0xcc, 0x0a, 0x1f, 0x23, 0xde, 0x00, 0x6c, 0x3e, 0x94, 0x55, 0xdc,
	0x15, 0xb2, 0xc5, 0x98, 0x68, 0x87, 0x30, 0xcf, 0x7f, 0x71, 0xc4, 0x46, 0x78, 0x0f, 0xf1, 0xe0,
	0x84, 0xb4, 0xe1, 0x7f, 0x42, 0xfd, 0x88, 0xea, 0xb8, 0x76, 0x6a, 0x98, 0x8d, 0x5e, 0x67, 0xfe,
	0xc3, 0x6e, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0xc3, 0xcb, 0x3c, 0x3a, 0xc0, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// MasterClient is the client API for Master service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MasterClient interface {
	WhereIsMap(ctx context.Context, in *MapQuery, opts ...grpc.CallOption) (*ConnectionInfo, error)
	RegisterZone(ctx context.Context, in *ZoneDetails, opts ...grpc.CallOption) (*ZoneRegistered, error)
}

type masterClient struct {
	cc grpc.ClientConnInterface
}

func NewMasterClient(cc grpc.ClientConnInterface) MasterClient {
	return &masterClient{cc}
}

func (c *masterClient) WhereIsMap(ctx context.Context, in *MapQuery, opts ...grpc.CallOption) (*ConnectionInfo, error) {
	out := new(ConnectionInfo)
	err := c.cc.Invoke(ctx, "/zm.Master/WhereIsMap", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterClient) RegisterZone(ctx context.Context, in *ZoneDetails, opts ...grpc.CallOption) (*ZoneRegistered, error) {
	out := new(ZoneRegistered)
	err := c.cc.Invoke(ctx, "/zm.Master/RegisterZone", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MasterServer is the server API for Master service.
type MasterServer interface {
	WhereIsMap(context.Context, *MapQuery) (*ConnectionInfo, error)
	RegisterZone(context.Context, *ZoneDetails) (*ZoneRegistered, error)
}

// UnimplementedMasterServer can be embedded to have forward compatible implementations.
type UnimplementedMasterServer struct {
}

func (*UnimplementedMasterServer) WhereIsMap(ctx context.Context, req *MapQuery) (*ConnectionInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WhereIsMap not implemented")
}
func (*UnimplementedMasterServer) RegisterZone(ctx context.Context, req *ZoneDetails) (*ZoneRegistered, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterZone not implemented")
}

func RegisterMasterServer(s *grpc.Server, srv MasterServer) {
	s.RegisterService(&_Master_serviceDesc, srv)
}

func _Master_WhereIsMap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MapQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterServer).WhereIsMap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zm.Master/WhereIsMap",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterServer).WhereIsMap(ctx, req.(*MapQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _Master_RegisterZone_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ZoneDetails)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterServer).RegisterZone(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/zm.Master/RegisterZone",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterServer).RegisterZone(ctx, req.(*ZoneDetails))
	}
	return interceptor(ctx, in, info, handler)
}

var _Master_serviceDesc = grpc.ServiceDesc{
	ServiceName: "zm.Master",
	HandlerType: (*MasterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "WhereIsMap",
			Handler:    _Master_WhereIsMap_Handler,
		},
		{
			MethodName: "RegisterZone",
			Handler:    _Master_RegisterZone_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "master.proto",
}
