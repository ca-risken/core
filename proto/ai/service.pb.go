// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.29.3
// source: ai/service.proto

package ai

import (
	context "context"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
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

type ChatAIRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Question    string         `protobuf:"bytes,1,opt,name=question,proto3" json:"question,omitempty"` // Required
	ChatHistory []*ChatHistory `protobuf:"bytes,2,rep,name=chat_history,json=chatHistory,proto3" json:"chat_history,omitempty"`
}

func (x *ChatAIRequest) Reset() {
	*x = ChatAIRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ai_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatAIRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatAIRequest) ProtoMessage() {}

func (x *ChatAIRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ai_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatAIRequest.ProtoReflect.Descriptor instead.
func (*ChatAIRequest) Descriptor() ([]byte, []int) {
	return file_ai_service_proto_rawDescGZIP(), []int{0}
}

func (x *ChatAIRequest) GetQuestion() string {
	if x != nil {
		return x.Question
	}
	return ""
}

func (x *ChatAIRequest) GetChatHistory() []*ChatHistory {
	if x != nil {
		return x.ChatHistory
	}
	return nil
}

type ChatAIResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Answer string `protobuf:"bytes,1,opt,name=answer,proto3" json:"answer,omitempty"`
}

func (x *ChatAIResponse) Reset() {
	*x = ChatAIResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ai_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChatAIResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChatAIResponse) ProtoMessage() {}

func (x *ChatAIResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ai_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChatAIResponse.ProtoReflect.Descriptor instead.
func (*ChatAIResponse) Descriptor() ([]byte, []int) {
	return file_ai_service_proto_rawDescGZIP(), []int{1}
}

func (x *ChatAIResponse) GetAnswer() string {
	if x != nil {
		return x.Answer
	}
	return ""
}

type GenerateReportRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Prompt    string `protobuf:"bytes,1,opt,name=prompt,proto3" json:"prompt,omitempty"`                         // Required
	ProjectId uint32 `protobuf:"varint,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"` // Required
}

func (x *GenerateReportRequest) Reset() {
	*x = GenerateReportRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ai_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenerateReportRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerateReportRequest) ProtoMessage() {}

func (x *GenerateReportRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ai_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerateReportRequest.ProtoReflect.Descriptor instead.
func (*GenerateReportRequest) Descriptor() ([]byte, []int) {
	return file_ai_service_proto_rawDescGZIP(), []int{2}
}

func (x *GenerateReportRequest) GetPrompt() string {
	if x != nil {
		return x.Prompt
	}
	return ""
}

func (x *GenerateReportRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

type GenerateReportResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Report string `protobuf:"bytes,1,opt,name=report,proto3" json:"report,omitempty"`
}

func (x *GenerateReportResponse) Reset() {
	*x = GenerateReportResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ai_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GenerateReportResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GenerateReportResponse) ProtoMessage() {}

func (x *GenerateReportResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ai_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GenerateReportResponse.ProtoReflect.Descriptor instead.
func (*GenerateReportResponse) Descriptor() ([]byte, []int) {
	return file_ai_service_proto_rawDescGZIP(), []int{3}
}

func (x *GenerateReportResponse) GetReport() string {
	if x != nil {
		return x.Report
	}
	return ""
}

var File_ai_service_proto protoreflect.FileDescriptor

var file_ai_service_proto_rawDesc = []byte{
	0x0a, 0x10, 0x61, 0x69, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x07, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x69, 0x1a, 0x0f, 0x61, 0x69, 0x2f,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6d, 0x0a, 0x0d, 0x43, 0x68, 0x61, 0x74, 0x41, 0x49, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x08, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10,
	0x01, 0x52, 0x08, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x37, 0x0a, 0x0c, 0x63,
	0x68, 0x61, 0x74, 0x5f, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x69, 0x2e, 0x43, 0x68, 0x61, 0x74,
	0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x52, 0x0b, 0x63, 0x68, 0x61, 0x74, 0x48, 0x69, 0x73,
	0x74, 0x6f, 0x72, 0x79, 0x22, 0x28, 0x0a, 0x0e, 0x43, 0x68, 0x61, 0x74, 0x41, 0x49, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x22, 0x60,
	0x0a, 0x15, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x06, 0x70, 0x72, 0x6f, 0x6d, 0x70,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01,
	0x52, 0x06, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x12, 0x26, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x07, 0xfa, 0x42,
	0x04, 0x2a, 0x02, 0x20, 0x00, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64,
	0x22, 0x30, 0x0a, 0x16, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x32, 0x99, 0x01, 0x0a, 0x09, 0x41, 0x49, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x39, 0x0a, 0x06, 0x43, 0x68, 0x61, 0x74, 0x41, 0x49, 0x12, 0x16, 0x2e, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x61, 0x69, 0x2e, 0x43, 0x68, 0x61, 0x74, 0x41, 0x49, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x17, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x69, 0x2e, 0x43, 0x68, 0x61,
	0x74, 0x41, 0x49, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x51, 0x0a, 0x0e, 0x47,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x1e, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x69, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x61, 0x69, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x24,
	0x5a, 0x22, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x61, 0x2d,
	0x72, 0x69, 0x73, 0x6b, 0x65, 0x6e, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x61, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ai_service_proto_rawDescOnce sync.Once
	file_ai_service_proto_rawDescData = file_ai_service_proto_rawDesc
)

func file_ai_service_proto_rawDescGZIP() []byte {
	file_ai_service_proto_rawDescOnce.Do(func() {
		file_ai_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_ai_service_proto_rawDescData)
	})
	return file_ai_service_proto_rawDescData
}

var file_ai_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_ai_service_proto_goTypes = []interface{}{
	(*ChatAIRequest)(nil),          // 0: core.ai.ChatAIRequest
	(*ChatAIResponse)(nil),         // 1: core.ai.ChatAIResponse
	(*GenerateReportRequest)(nil),  // 2: core.ai.GenerateReportRequest
	(*GenerateReportResponse)(nil), // 3: core.ai.GenerateReportResponse
	(*ChatHistory)(nil),            // 4: core.ai.ChatHistory
}
var file_ai_service_proto_depIdxs = []int32{
	4, // 0: core.ai.ChatAIRequest.chat_history:type_name -> core.ai.ChatHistory
	0, // 1: core.ai.AIService.ChatAI:input_type -> core.ai.ChatAIRequest
	2, // 2: core.ai.AIService.GenerateReport:input_type -> core.ai.GenerateReportRequest
	1, // 3: core.ai.AIService.ChatAI:output_type -> core.ai.ChatAIResponse
	3, // 4: core.ai.AIService.GenerateReport:output_type -> core.ai.GenerateReportResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_ai_service_proto_init() }
func file_ai_service_proto_init() {
	if File_ai_service_proto != nil {
		return
	}
	file_ai_entity_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_ai_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatAIRequest); i {
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
		file_ai_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChatAIResponse); i {
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
		file_ai_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenerateReportRequest); i {
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
		file_ai_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GenerateReportResponse); i {
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
			RawDescriptor: file_ai_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ai_service_proto_goTypes,
		DependencyIndexes: file_ai_service_proto_depIdxs,
		MessageInfos:      file_ai_service_proto_msgTypes,
	}.Build()
	File_ai_service_proto = out.File
	file_ai_service_proto_rawDesc = nil
	file_ai_service_proto_goTypes = nil
	file_ai_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AIServiceClient is the client API for AIService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AIServiceClient interface {
	ChatAI(ctx context.Context, in *ChatAIRequest, opts ...grpc.CallOption) (*ChatAIResponse, error)
	GenerateReport(ctx context.Context, in *GenerateReportRequest, opts ...grpc.CallOption) (*GenerateReportResponse, error)
}

type aIServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAIServiceClient(cc grpc.ClientConnInterface) AIServiceClient {
	return &aIServiceClient{cc}
}

func (c *aIServiceClient) ChatAI(ctx context.Context, in *ChatAIRequest, opts ...grpc.CallOption) (*ChatAIResponse, error) {
	out := new(ChatAIResponse)
	err := c.cc.Invoke(ctx, "/core.ai.AIService/ChatAI", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aIServiceClient) GenerateReport(ctx context.Context, in *GenerateReportRequest, opts ...grpc.CallOption) (*GenerateReportResponse, error) {
	out := new(GenerateReportResponse)
	err := c.cc.Invoke(ctx, "/core.ai.AIService/GenerateReport", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AIServiceServer is the server API for AIService service.
type AIServiceServer interface {
	ChatAI(context.Context, *ChatAIRequest) (*ChatAIResponse, error)
	GenerateReport(context.Context, *GenerateReportRequest) (*GenerateReportResponse, error)
}

// UnimplementedAIServiceServer can be embedded to have forward compatible implementations.
type UnimplementedAIServiceServer struct {
}

func (*UnimplementedAIServiceServer) ChatAI(context.Context, *ChatAIRequest) (*ChatAIResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChatAI not implemented")
}
func (*UnimplementedAIServiceServer) GenerateReport(context.Context, *GenerateReportRequest) (*GenerateReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GenerateReport not implemented")
}

func RegisterAIServiceServer(s *grpc.Server, srv AIServiceServer) {
	s.RegisterService(&_AIService_serviceDesc, srv)
}

func _AIService_ChatAI_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChatAIRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AIServiceServer).ChatAI(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.ai.AIService/ChatAI",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AIServiceServer).ChatAI(ctx, req.(*ChatAIRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AIService_GenerateReport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AIServiceServer).GenerateReport(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.ai.AIService/GenerateReport",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AIServiceServer).GenerateReport(ctx, req.(*GenerateReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AIService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "core.ai.AIService",
	HandlerType: (*AIServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ChatAI",
			Handler:    _AIService_ChatAI_Handler,
		},
		{
			MethodName: "GenerateReport",
			Handler:    _AIService_GenerateReport_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ai/service.proto",
}
