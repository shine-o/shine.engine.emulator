// Code generated by protoc-gen-go. DO NOT EDIT.
// source: login.proto

package login

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = proto.Marshal
	_ = fmt.Errorf
	_ = math.Inf
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type User struct {
	UserName             string   `protobuf:"bytes,1,opt,name=UserName,proto3" json:"UserName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_67c21677aa7f4e4f, []int{0}
}

func (m *User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_User.Unmarshal(m, b)
}

func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_User.Marshal(b, m, deterministic)
}

func (m *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(m, src)
}

func (m *User) XXX_Size() int {
	return xxx_messageInfo_User.Size(m)
}

func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetUserName() string {
	if m != nil {
		return m.UserName
	}
	return ""
}

type UserInfo struct {
	UserID               uint64   `protobuf:"varint,1,opt,name=UserID,proto3" json:"UserID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserInfo) Reset()         { *m = UserInfo{} }
func (m *UserInfo) String() string { return proto.CompactTextString(m) }
func (*UserInfo) ProtoMessage()    {}
func (*UserInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_67c21677aa7f4e4f, []int{1}
}

func (m *UserInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserInfo.Unmarshal(m, b)
}

func (m *UserInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserInfo.Marshal(b, m, deterministic)
}

func (m *UserInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserInfo.Merge(m, src)
}

func (m *UserInfo) XXX_Size() int {
	return xxx_messageInfo_UserInfo.Size(m)
}

func (m *UserInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_UserInfo.DiscardUnknown(m)
}

var xxx_messageInfo_UserInfo proto.InternalMessageInfo

func (m *UserInfo) GetUserID() uint64 {
	if m != nil {
		return m.UserID
	}
	return 0
}

func init() {
	proto.RegisterType((*User)(nil), "login.User")
	proto.RegisterType((*UserInfo)(nil), "login.UserInfo")
}

func init() {
	proto.RegisterFile("login.proto", fileDescriptor_67c21677aa7f4e4f)
}

var fileDescriptor_67c21677aa7f4e4f = []byte{
	// 127 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0xc9, 0x4f, 0xcf,
	0xcc, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x73, 0x94, 0x94, 0xb8, 0x58, 0x42,
	0x8b, 0x53, 0x8b, 0x84, 0xa4, 0xb8, 0x38, 0x40, 0xb4, 0x5f, 0x62, 0x6e, 0xaa, 0x04, 0xa3, 0x02,
	0xa3, 0x06, 0x67, 0x10, 0x9c, 0xaf, 0xa4, 0x04, 0x91, 0xf3, 0xcc, 0x4b, 0xcb, 0x17, 0x12, 0xe3,
	0x62, 0x03, 0xb3, 0x5d, 0xc0, 0xaa, 0x58, 0x82, 0xa0, 0x3c, 0x23, 0x33, 0x2e, 0x56, 0x1f, 0x90,
	0x81, 0x42, 0xba, 0x5c, 0xdc, 0x8e, 0xc9, 0xc9, 0xf9, 0xa5, 0x79, 0x25, 0x60, 0xf5, 0xdc, 0x7a,
	0x10, 0x4b, 0x41, 0xca, 0xa4, 0xf8, 0x91, 0x38, 0x20, 0x59, 0x25, 0x86, 0x24, 0x36, 0xb0, 0x6b,
	0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x74, 0x13, 0x8e, 0x92, 0x9c, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ context.Context
	_ grpc.ClientConnInterface
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// LoginClient is the client API for Login service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LoginClient interface {
	AccountInfo(ctx context.Context, in *User, opts ...grpc.CallOption) (*UserInfo, error)
}

type loginClient struct {
	cc grpc.ClientConnInterface
}

func NewLoginClient(cc grpc.ClientConnInterface) LoginClient {
	return &loginClient{cc}
}

func (c *loginClient) AccountInfo(ctx context.Context, in *User, opts ...grpc.CallOption) (*UserInfo, error) {
	out := new(UserInfo)
	err := c.cc.Invoke(ctx, "/login.Login/AccountInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LoginServer is the server API for Login service.
type LoginServer interface {
	AccountInfo(context.Context, *User) (*UserInfo, error)
}

// UnimplementedLoginServer can be embedded to have forward compatible implementations.
type UnimplementedLoginServer struct{}

func (*UnimplementedLoginServer) AccountInfo(ctx context.Context, req *User) (*UserInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AccountInfo not implemented")
}

func RegisterLoginServer(s *grpc.Server, srv LoginServer) {
	s.RegisterService(&_Login_serviceDesc, srv)
}

func _Login_AccountInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(User)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoginServer).AccountInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/login.Login/AccountInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoginServer).AccountInfo(ctx, req.(*User))
	}
	return interceptor(ctx, in, info, handler)
}

var _Login_serviceDesc = grpc.ServiceDesc{
	ServiceName: "login.Login",
	HandlerType: (*LoginServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AccountInfo",
			Handler:    _Login_AccountInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "login.proto",
}
