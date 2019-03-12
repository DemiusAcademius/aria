// Code generated by protoc-gen-go. DO NOT EDIT.
// source: image-builder.proto

package api

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

type ResponseCode int32

const (
	ResponseCode_OK ResponseCode = 0
)

var ResponseCode_name = map[int32]string{
	0: "OK",
}
var ResponseCode_value = map[string]int32{
	"OK": 0,
}

func (x ResponseCode) String() string {
	return proto.EnumName(ResponseCode_name, int32(x))
}
func (ResponseCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_image_builder_11535617c483f665, []int{0}
}

type GitRepo_ProjectType int32

const (
	GitRepo_DOTNET GitRepo_ProjectType = 0
	GitRepo_JAVA   GitRepo_ProjectType = 1
	GitRepo_NODEJS GitRepo_ProjectType = 2
)

var GitRepo_ProjectType_name = map[int32]string{
	0: "DOTNET",
	1: "JAVA",
	2: "NODEJS",
}
var GitRepo_ProjectType_value = map[string]int32{
	"DOTNET": 0,
	"JAVA":   1,
	"NODEJS": 2,
}

func (x GitRepo_ProjectType) String() string {
	return proto.EnumName(GitRepo_ProjectType_name, int32(x))
}
func (GitRepo_ProjectType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_image_builder_11535617c483f665, []int{1, 0}
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_image_builder_11535617c483f665, []int{0}
}
func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (dst *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(dst, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type GitRepo struct {
	ProjectType          GitRepo_ProjectType `protobuf:"varint,1,opt,name=projectType,proto3,enum=api.GitRepo_ProjectType" json:"projectType,omitempty"`
	Provider             string              `protobuf:"bytes,2,opt,name=provider,proto3" json:"provider,omitempty"`
	Repo                 string              `protobuf:"bytes,3,opt,name=repo,proto3" json:"repo,omitempty"`
	Username             string              `protobuf:"bytes,4,opt,name=username,proto3" json:"username,omitempty"`
	Password             string              `protobuf:"bytes,5,opt,name=password,proto3" json:"password,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *GitRepo) Reset()         { *m = GitRepo{} }
func (m *GitRepo) String() string { return proto.CompactTextString(m) }
func (*GitRepo) ProtoMessage()    {}
func (*GitRepo) Descriptor() ([]byte, []int) {
	return fileDescriptor_image_builder_11535617c483f665, []int{1}
}
func (m *GitRepo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GitRepo.Unmarshal(m, b)
}
func (m *GitRepo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GitRepo.Marshal(b, m, deterministic)
}
func (dst *GitRepo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GitRepo.Merge(dst, src)
}
func (m *GitRepo) XXX_Size() int {
	return xxx_messageInfo_GitRepo.Size(m)
}
func (m *GitRepo) XXX_DiscardUnknown() {
	xxx_messageInfo_GitRepo.DiscardUnknown(m)
}

var xxx_messageInfo_GitRepo proto.InternalMessageInfo

func (m *GitRepo) GetProjectType() GitRepo_ProjectType {
	if m != nil {
		return m.ProjectType
	}
	return GitRepo_DOTNET
}

func (m *GitRepo) GetProvider() string {
	if m != nil {
		return m.Provider
	}
	return ""
}

func (m *GitRepo) GetRepo() string {
	if m != nil {
		return m.Repo
	}
	return ""
}

func (m *GitRepo) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *GitRepo) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type Namespace struct {
	Namespace            string   `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Namespace) Reset()         { *m = Namespace{} }
func (m *Namespace) String() string { return proto.CompactTextString(m) }
func (*Namespace) ProtoMessage()    {}
func (*Namespace) Descriptor() ([]byte, []int) {
	return fileDescriptor_image_builder_11535617c483f665, []int{2}
}
func (m *Namespace) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Namespace.Unmarshal(m, b)
}
func (m *Namespace) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Namespace.Marshal(b, m, deterministic)
}
func (dst *Namespace) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Namespace.Merge(dst, src)
}
func (m *Namespace) XXX_Size() int {
	return xxx_messageInfo_Namespace.Size(m)
}
func (m *Namespace) XXX_DiscardUnknown() {
	xxx_messageInfo_Namespace.DiscardUnknown(m)
}

var xxx_messageInfo_Namespace proto.InternalMessageInfo

func (m *Namespace) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

type Project struct {
	Namespace            string   `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Project              string   `protobuf:"bytes,2,opt,name=project,proto3" json:"project,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Project) Reset()         { *m = Project{} }
func (m *Project) String() string { return proto.CompactTextString(m) }
func (*Project) ProtoMessage()    {}
func (*Project) Descriptor() ([]byte, []int) {
	return fileDescriptor_image_builder_11535617c483f665, []int{3}
}
func (m *Project) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Project.Unmarshal(m, b)
}
func (m *Project) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Project.Marshal(b, m, deterministic)
}
func (dst *Project) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Project.Merge(dst, src)
}
func (m *Project) XXX_Size() int {
	return xxx_messageInfo_Project.Size(m)
}
func (m *Project) XXX_DiscardUnknown() {
	xxx_messageInfo_Project.DiscardUnknown(m)
}

var xxx_messageInfo_Project proto.InternalMessageInfo

func (m *Project) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *Project) GetProject() string {
	if m != nil {
		return m.Project
	}
	return ""
}

type Response struct {
	Code                 ResponseCode `protobuf:"varint,1,opt,name=code,proto3,enum=api.ResponseCode" json:"code,omitempty"`
	ErrorDesctiption     string       `protobuf:"bytes,2,opt,name=errorDesctiption,proto3" json:"errorDesctiption,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_image_builder_11535617c483f665, []int{4}
}
func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (dst *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(dst, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetCode() ResponseCode {
	if m != nil {
		return m.Code
	}
	return ResponseCode_OK
}

func (m *Response) GetErrorDesctiption() string {
	if m != nil {
		return m.ErrorDesctiption
	}
	return ""
}

type BuildResponse struct {
	Project              *Project     `protobuf:"bytes,1,opt,name=project,proto3" json:"project,omitempty"`
	Code                 ResponseCode `protobuf:"varint,2,opt,name=code,proto3,enum=api.ResponseCode" json:"code,omitempty"`
	ErrorDesctiption     string       `protobuf:"bytes,3,opt,name=errorDesctiption,proto3" json:"errorDesctiption,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *BuildResponse) Reset()         { *m = BuildResponse{} }
func (m *BuildResponse) String() string { return proto.CompactTextString(m) }
func (*BuildResponse) ProtoMessage()    {}
func (*BuildResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_image_builder_11535617c483f665, []int{5}
}
func (m *BuildResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BuildResponse.Unmarshal(m, b)
}
func (m *BuildResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BuildResponse.Marshal(b, m, deterministic)
}
func (dst *BuildResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BuildResponse.Merge(dst, src)
}
func (m *BuildResponse) XXX_Size() int {
	return xxx_messageInfo_BuildResponse.Size(m)
}
func (m *BuildResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_BuildResponse.DiscardUnknown(m)
}

var xxx_messageInfo_BuildResponse proto.InternalMessageInfo

func (m *BuildResponse) GetProject() *Project {
	if m != nil {
		return m.Project
	}
	return nil
}

func (m *BuildResponse) GetCode() ResponseCode {
	if m != nil {
		return m.Code
	}
	return ResponseCode_OK
}

func (m *BuildResponse) GetErrorDesctiption() string {
	if m != nil {
		return m.ErrorDesctiption
	}
	return ""
}

func init() {
	proto.RegisterType((*Empty)(nil), "api.Empty")
	proto.RegisterType((*GitRepo)(nil), "api.GitRepo")
	proto.RegisterType((*Namespace)(nil), "api.Namespace")
	proto.RegisterType((*Project)(nil), "api.Project")
	proto.RegisterType((*Response)(nil), "api.Response")
	proto.RegisterType((*BuildResponse)(nil), "api.BuildResponse")
	proto.RegisterEnum("api.ResponseCode", ResponseCode_name, ResponseCode_value)
	proto.RegisterEnum("api.GitRepo_ProjectType", GitRepo_ProjectType_name, GitRepo_ProjectType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ImageBuilderClient is the client API for ImageBuilder service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ImageBuilderClient interface {
	Init(ctx context.Context, in *GitRepo, opts ...grpc.CallOption) (*Response, error)
	Pull(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error)
	BuildProject(ctx context.Context, in *Project, opts ...grpc.CallOption) (*BuildResponse, error)
	BuildNamespace(ctx context.Context, in *Namespace, opts ...grpc.CallOption) (ImageBuilder_BuildNamespaceClient, error)
	BuildAll(ctx context.Context, in *Empty, opts ...grpc.CallOption) (ImageBuilder_BuildAllClient, error)
}

type imageBuilderClient struct {
	cc *grpc.ClientConn
}

func NewImageBuilderClient(cc *grpc.ClientConn) ImageBuilderClient {
	return &imageBuilderClient{cc}
}

func (c *imageBuilderClient) Init(ctx context.Context, in *GitRepo, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.ImageBuilder/Init", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imageBuilderClient) Pull(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/api.ImageBuilder/Pull", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imageBuilderClient) BuildProject(ctx context.Context, in *Project, opts ...grpc.CallOption) (*BuildResponse, error) {
	out := new(BuildResponse)
	err := c.cc.Invoke(ctx, "/api.ImageBuilder/BuildProject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *imageBuilderClient) BuildNamespace(ctx context.Context, in *Namespace, opts ...grpc.CallOption) (ImageBuilder_BuildNamespaceClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ImageBuilder_serviceDesc.Streams[0], "/api.ImageBuilder/BuildNamespace", opts...)
	if err != nil {
		return nil, err
	}
	x := &imageBuilderBuildNamespaceClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ImageBuilder_BuildNamespaceClient interface {
	Recv() (*BuildResponse, error)
	grpc.ClientStream
}

type imageBuilderBuildNamespaceClient struct {
	grpc.ClientStream
}

func (x *imageBuilderBuildNamespaceClient) Recv() (*BuildResponse, error) {
	m := new(BuildResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *imageBuilderClient) BuildAll(ctx context.Context, in *Empty, opts ...grpc.CallOption) (ImageBuilder_BuildAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ImageBuilder_serviceDesc.Streams[1], "/api.ImageBuilder/BuildAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &imageBuilderBuildAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ImageBuilder_BuildAllClient interface {
	Recv() (*BuildResponse, error)
	grpc.ClientStream
}

type imageBuilderBuildAllClient struct {
	grpc.ClientStream
}

func (x *imageBuilderBuildAllClient) Recv() (*BuildResponse, error) {
	m := new(BuildResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ImageBuilderServer is the server API for ImageBuilder service.
type ImageBuilderServer interface {
	Init(context.Context, *GitRepo) (*Response, error)
	Pull(context.Context, *Empty) (*Response, error)
	BuildProject(context.Context, *Project) (*BuildResponse, error)
	BuildNamespace(*Namespace, ImageBuilder_BuildNamespaceServer) error
	BuildAll(*Empty, ImageBuilder_BuildAllServer) error
}

func RegisterImageBuilderServer(s *grpc.Server, srv ImageBuilderServer) {
	s.RegisterService(&_ImageBuilder_serviceDesc, srv)
}

func _ImageBuilder_Init_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GitRepo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImageBuilderServer).Init(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ImageBuilder/Init",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImageBuilderServer).Init(ctx, req.(*GitRepo))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImageBuilder_Pull_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImageBuilderServer).Pull(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ImageBuilder/Pull",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImageBuilderServer).Pull(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImageBuilder_BuildProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Project)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ImageBuilderServer).BuildProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ImageBuilder/BuildProject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ImageBuilderServer).BuildProject(ctx, req.(*Project))
	}
	return interceptor(ctx, in, info, handler)
}

func _ImageBuilder_BuildNamespace_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Namespace)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ImageBuilderServer).BuildNamespace(m, &imageBuilderBuildNamespaceServer{stream})
}

type ImageBuilder_BuildNamespaceServer interface {
	Send(*BuildResponse) error
	grpc.ServerStream
}

type imageBuilderBuildNamespaceServer struct {
	grpc.ServerStream
}

func (x *imageBuilderBuildNamespaceServer) Send(m *BuildResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _ImageBuilder_BuildAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ImageBuilderServer).BuildAll(m, &imageBuilderBuildAllServer{stream})
}

type ImageBuilder_BuildAllServer interface {
	Send(*BuildResponse) error
	grpc.ServerStream
}

type imageBuilderBuildAllServer struct {
	grpc.ServerStream
}

func (x *imageBuilderBuildAllServer) Send(m *BuildResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _ImageBuilder_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.ImageBuilder",
	HandlerType: (*ImageBuilderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Init",
			Handler:    _ImageBuilder_Init_Handler,
		},
		{
			MethodName: "Pull",
			Handler:    _ImageBuilder_Pull_Handler,
		},
		{
			MethodName: "BuildProject",
			Handler:    _ImageBuilder_BuildProject_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "BuildNamespace",
			Handler:       _ImageBuilder_BuildNamespace_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "BuildAll",
			Handler:       _ImageBuilder_BuildAll_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "image-builder.proto",
}

func init() { proto.RegisterFile("image-builder.proto", fileDescriptor_image_builder_11535617c483f665) }

var fileDescriptor_image_builder_11535617c483f665 = []byte{
	// 435 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x53, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0xcd, 0x26, 0x6e, 0x3e, 0x26, 0x6e, 0x64, 0x06, 0x09, 0x59, 0x11, 0x87, 0x6a, 0x11, 0xa8,
	0x54, 0xc2, 0xa0, 0x70, 0x41, 0xdc, 0x5c, 0x12, 0xa1, 0x16, 0x29, 0xa9, 0x4c, 0xc4, 0x8d, 0x83,
	0x6b, 0x8f, 0xd0, 0xa2, 0xd8, 0xbb, 0x5a, 0x3b, 0xa0, 0xfe, 0x05, 0x7e, 0x24, 0xbf, 0x83, 0x23,
	0xf2, 0x76, 0xfd, 0x11, 0x5a, 0x90, 0xb8, 0xe5, 0xbd, 0x79, 0x79, 0xf3, 0x66, 0x76, 0x0c, 0x0f,
	0x45, 0x16, 0x7f, 0xa1, 0x17, 0xd7, 0x7b, 0xb1, 0x4b, 0x49, 0x07, 0x4a, 0xcb, 0x52, 0xe2, 0x20,
	0x56, 0x82, 0x8f, 0xe0, 0x68, 0x95, 0xa9, 0xf2, 0x86, 0xff, 0x64, 0x30, 0x7a, 0x2f, 0xca, 0x88,
	0x94, 0xc4, 0xb7, 0x30, 0x55, 0x5a, 0x7e, 0xa5, 0xa4, 0xdc, 0xde, 0x28, 0xf2, 0xd9, 0x09, 0x3b,
	0x9d, 0x2d, 0xfc, 0x20, 0x56, 0x22, 0xb0, 0x92, 0xe0, 0xaa, 0xad, 0x47, 0x5d, 0x31, 0xce, 0x61,
	0xac, 0xb4, 0xfc, 0x26, 0x52, 0xd2, 0x7e, 0xff, 0x84, 0x9d, 0x4e, 0xa2, 0x06, 0x23, 0x82, 0xa3,
	0x49, 0x49, 0x7f, 0x60, 0x78, 0xf3, 0xbb, 0xd2, 0xef, 0x0b, 0xd2, 0x79, 0x9c, 0x91, 0xef, 0xdc,
	0xea, 0x6b, 0x6c, 0xbc, 0xe2, 0xa2, 0xf8, 0x2e, 0x75, 0xea, 0x1f, 0x59, 0x2f, 0x8b, 0xf9, 0x4b,
	0x98, 0x76, 0x32, 0x20, 0xc0, 0x70, 0xb9, 0xd9, 0xae, 0x57, 0x5b, 0xaf, 0x87, 0x63, 0x70, 0x2e,
	0xc3, 0x4f, 0xa1, 0xc7, 0x2a, 0x76, 0xbd, 0x59, 0xae, 0x2e, 0x3f, 0x7a, 0x7d, 0xfe, 0x1c, 0x26,
	0xeb, 0x38, 0xa3, 0x42, 0xc5, 0x09, 0xe1, 0x63, 0x98, 0xe4, 0x35, 0x30, 0xf3, 0x4d, 0xa2, 0x96,
	0xe0, 0x21, 0x8c, 0xac, 0xf7, 0xbf, 0x85, 0xe8, 0xc3, 0xc8, 0xce, 0x6e, 0x67, 0xad, 0x21, 0xff,
	0x0c, 0xe3, 0x88, 0x0a, 0x25, 0xf3, 0x82, 0xf0, 0x29, 0x38, 0x89, 0x4c, 0xeb, 0x3d, 0x3e, 0x30,
	0x7b, 0xac, 0x8b, 0xef, 0x64, 0x4a, 0x91, 0x29, 0xe3, 0x19, 0x78, 0xa4, 0xb5, 0xd4, 0x4b, 0x2a,
	0x92, 0x52, 0xa8, 0x52, 0xc8, 0xdc, 0xba, 0xde, 0xe1, 0xf9, 0x0f, 0x06, 0xc7, 0xe7, 0xd5, 0x6b,
	0x36, 0x4d, 0x9e, 0xb5, 0x51, 0xaa, 0x3e, 0xd3, 0x85, 0x6b, 0xfa, 0xd8, 0x39, 0x9a, 0x60, 0x4d,
	0x98, 0xfe, 0xff, 0x87, 0x19, 0xdc, 0x1f, 0xe6, 0xec, 0x11, 0xb8, 0x5d, 0x07, 0x1c, 0x42, 0x7f,
	0xf3, 0xc1, 0xeb, 0x2d, 0x7e, 0x31, 0x70, 0x2f, 0xaa, 0xc3, 0x3b, 0xbf, 0xbd, 0xbb, 0xaa, 0xf7,
	0x45, 0x2e, 0x4a, 0x74, 0xbb, 0xa7, 0x34, 0x3f, 0x3e, 0xc8, 0xc0, 0x7b, 0xf8, 0x04, 0x9c, 0xab,
	0xfd, 0x6e, 0x87, 0x60, 0x0a, 0xe6, 0x3c, 0xef, 0x8a, 0x16, 0xe0, 0x1a, 0xdb, 0xfa, 0xa1, 0x0e,
	0xc6, 0x9d, 0xa3, 0x41, 0x07, 0x1b, 0xe2, 0x3d, 0x7c, 0x03, 0x33, 0x43, 0xb5, 0x77, 0x30, 0x33,
	0xba, 0x06, 0xdf, 0xff, 0xbf, 0x57, 0x0c, 0x03, 0x18, 0x1b, 0x32, 0xfc, 0x23, 0xd6, 0x5f, 0xf4,
	0xd7, 0x43, 0xf3, 0x89, 0xbd, 0xfe, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x41, 0xa6, 0x8a, 0x18, 0x79,
	0x03, 0x00, 0x00,
}
