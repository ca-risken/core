// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0-devel
// 	protoc        v3.12.1
// source: iam/service.proto

package iam

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// AuthnRequest tokenからユーザを識別します
type AuthnRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *AuthnRequest) Reset() {
	*x = AuthnRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthnRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthnRequest) ProtoMessage() {}

func (x *AuthnRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthnRequest.ProtoReflect.Descriptor instead.
func (*AuthnRequest) Descriptor() ([]byte, []int) {
	return file_iam_service_proto_rawDescGZIP(), []int{0}
}

func (x *AuthnRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type AuthnResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *AuthnResponse) Reset() {
	*x = AuthnResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthnResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthnResponse) ProtoMessage() {}

func (x *AuthnResponse) ProtoReflect() protoreflect.Message {
	mi := &file_iam_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthnResponse.ProtoReflect.Descriptor instead.
func (*AuthnResponse) Descriptor() ([]byte, []int) {
	return file_iam_service_proto_rawDescGZIP(), []int{1}
}

func (x *AuthnResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

// AuthzRequest
// ユーザからのリクエストに対して、アクションやリソースへの認可を行います
type AuthzRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId       uint32 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`                  // UserID,(e.g.)111
	ProjectId    uint32 `protobuf:"varint,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`         // ProjectID,(e.g.)1001
	ActionName   string `protobuf:"bytes,3,opt,name=action_name,json=actionName,proto3" json:"action_name,omitempty"`       // Service&API_name(<service_name>/<API>format),(e.g.)`finding/GetFinding`
	ResourceName string `protobuf:"bytes,4,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"` // System_resource_name(<prefix>/<resouorce_name>format),(e.g.)`aws:accessAnalyzer/samplebucket`
}

func (x *AuthzRequest) Reset() {
	*x = AuthzRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthzRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthzRequest) ProtoMessage() {}

func (x *AuthzRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthzRequest.ProtoReflect.Descriptor instead.
func (*AuthzRequest) Descriptor() ([]byte, []int) {
	return file_iam_service_proto_rawDescGZIP(), []int{2}
}

func (x *AuthzRequest) GetUserId() uint32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *AuthzRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *AuthzRequest) GetActionName() string {
	if x != nil {
		return x.ActionName
	}
	return ""
}

func (x *AuthzRequest) GetResourceName() string {
	if x != nil {
		return x.ResourceName
	}
	return ""
}

type AuthzResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *AuthzResponse) Reset() {
	*x = AuthzResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthzResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthzResponse) ProtoMessage() {}

func (x *AuthzResponse) ProtoReflect() protoreflect.Message {
	mi := &file_iam_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthzResponse.ProtoReflect.Descriptor instead.
func (*AuthzResponse) Descriptor() ([]byte, []int) {
	return file_iam_service_proto_rawDescGZIP(), []int{3}
}

func (x *AuthzResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

var File_iam_service_proto protoreflect.FileDescriptor

var file_iam_service_proto_rawDesc = []byte{
	0x0a, 0x11, 0x69, 0x61, 0x6d, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x08, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x1a, 0x10, 0x69,
	0x61, 0x6d, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x24, 0x0a, 0x0c, 0x41, 0x75, 0x74, 0x68, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x33, 0x0a, 0x0d, 0x41, 0x75, 0x74, 0x68, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x22, 0x8c, 0x01, 0x0a, 0x0c, 0x41,
	0x75, 0x74, 0x68, 0x7a, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x1f, 0x0a, 0x0d, 0x41, 0x75, 0x74,
	0x68, 0x7a, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x32, 0x8d, 0x01, 0x0a, 0x0a, 0x49,
	0x41, 0x4d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x40, 0x0a, 0x0d, 0x41, 0x75, 0x74,
	0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x64, 0x12, 0x16, 0x2e, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x17, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x41, 0x75,
	0x74, 0x68, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3d, 0x0a, 0x0a, 0x41,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x64, 0x12, 0x16, 0x2e, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x7a, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x17, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x41, 0x75, 0x74,
	0x68, 0x7a, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x43, 0x79, 0x62, 0x65, 0x72, 0x41, 0x67,
	0x65, 0x6e, 0x74, 0x2f, 0x6d, 0x69, 0x6d, 0x6f, 0x73, 0x61, 0x2d, 0x63, 0x6f, 0x72, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x69, 0x61, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_iam_service_proto_rawDescOnce sync.Once
	file_iam_service_proto_rawDescData = file_iam_service_proto_rawDesc
)

func file_iam_service_proto_rawDescGZIP() []byte {
	file_iam_service_proto_rawDescOnce.Do(func() {
		file_iam_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_iam_service_proto_rawDescData)
	})
	return file_iam_service_proto_rawDescData
}

var file_iam_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_iam_service_proto_goTypes = []interface{}{
	(*AuthnRequest)(nil),  // 0: core.iam.AuthnRequest
	(*AuthnResponse)(nil), // 1: core.iam.AuthnResponse
	(*AuthzRequest)(nil),  // 2: core.iam.AuthzRequest
	(*AuthzResponse)(nil), // 3: core.iam.AuthzResponse
	(*User)(nil),          // 4: core.iam.User
}
var file_iam_service_proto_depIdxs = []int32{
	4, // 0: core.iam.AuthnResponse.user:type_name -> core.iam.User
	0, // 1: core.iam.IAMService.Authenticated:input_type -> core.iam.AuthnRequest
	2, // 2: core.iam.IAMService.Authorized:input_type -> core.iam.AuthzRequest
	1, // 3: core.iam.IAMService.Authenticated:output_type -> core.iam.AuthnResponse
	3, // 4: core.iam.IAMService.Authorized:output_type -> core.iam.AuthzResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_iam_service_proto_init() }
func file_iam_service_proto_init() {
	if File_iam_service_proto != nil {
		return
	}
	file_iam_entity_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_iam_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthnRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_iam_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthnResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_iam_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthzRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_iam_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthzResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_iam_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_iam_service_proto_goTypes,
		DependencyIndexes: file_iam_service_proto_depIdxs,
		MessageInfos:      file_iam_service_proto_msgTypes,
	}.Build()
	File_iam_service_proto = out.File
	file_iam_service_proto_rawDesc = nil
	file_iam_service_proto_goTypes = nil
	file_iam_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// IAMServiceClient is the client API for IAMService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type IAMServiceClient interface {
	// 認証（リクエストユーザを識別します）
	Authenticated(ctx context.Context, in *AuthnRequest, opts ...grpc.CallOption) (*AuthnResponse, error)
	// 認可（ユーザがリクエストしたアクションや、リソースに対しての認可を行います）
	Authorized(ctx context.Context, in *AuthzRequest, opts ...grpc.CallOption) (*AuthzResponse, error)
}

type iAMServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewIAMServiceClient(cc grpc.ClientConnInterface) IAMServiceClient {
	return &iAMServiceClient{cc}
}

func (c *iAMServiceClient) Authenticated(ctx context.Context, in *AuthnRequest, opts ...grpc.CallOption) (*AuthnResponse, error) {
	out := new(AuthnResponse)
	err := c.cc.Invoke(ctx, "/core.iam.IAMService/Authenticated", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iAMServiceClient) Authorized(ctx context.Context, in *AuthzRequest, opts ...grpc.CallOption) (*AuthzResponse, error) {
	out := new(AuthzResponse)
	err := c.cc.Invoke(ctx, "/core.iam.IAMService/Authorized", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IAMServiceServer is the server API for IAMService service.
type IAMServiceServer interface {
	// 認証（リクエストユーザを識別します）
	Authenticated(context.Context, *AuthnRequest) (*AuthnResponse, error)
	// 認可（ユーザがリクエストしたアクションや、リソースに対しての認可を行います）
	Authorized(context.Context, *AuthzRequest) (*AuthzResponse, error)
}

// UnimplementedIAMServiceServer can be embedded to have forward compatible implementations.
type UnimplementedIAMServiceServer struct {
}

func (*UnimplementedIAMServiceServer) Authenticated(context.Context, *AuthnRequest) (*AuthnResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authenticated not implemented")
}
func (*UnimplementedIAMServiceServer) Authorized(context.Context, *AuthzRequest) (*AuthzResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authorized not implemented")
}

func RegisterIAMServiceServer(s *grpc.Server, srv IAMServiceServer) {
	s.RegisterService(&_IAMService_serviceDesc, srv)
}

func _IAMService_Authenticated_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthnRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IAMServiceServer).Authenticated(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.iam.IAMService/Authenticated",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IAMServiceServer).Authenticated(ctx, req.(*AuthnRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IAMService_Authorized_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthzRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IAMServiceServer).Authorized(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.iam.IAMService/Authorized",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IAMServiceServer).Authorized(ctx, req.(*AuthzRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _IAMService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "core.iam.IAMService",
	HandlerType: (*IAMServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Authenticated",
			Handler:    _IAMService_Authenticated_Handler,
		},
		{
			MethodName: "Authorized",
			Handler:    _IAMService_Authorized_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "iam/service.proto",
}
