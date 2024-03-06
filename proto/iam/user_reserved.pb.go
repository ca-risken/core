// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.25.3
// source: iam/user_reserved.proto

package iam

import (
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

type ListUserReservedRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId  uint32 `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	UserIdpKey string `protobuf:"bytes,2,opt,name=user_idp_key,json=userIdpKey,proto3" json:"user_idp_key,omitempty"`
}

func (x *ListUserReservedRequest) Reset() {
	*x = ListUserReservedRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_user_reserved_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListUserReservedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListUserReservedRequest) ProtoMessage() {}

func (x *ListUserReservedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_user_reserved_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListUserReservedRequest.ProtoReflect.Descriptor instead.
func (*ListUserReservedRequest) Descriptor() ([]byte, []int) {
	return file_iam_user_reserved_proto_rawDescGZIP(), []int{0}
}

func (x *ListUserReservedRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *ListUserReservedRequest) GetUserIdpKey() string {
	if x != nil {
		return x.UserIdpKey
	}
	return ""
}

type ListUserReservedResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserReserved []*UserReserved `protobuf:"bytes,1,rep,name=user_reserved,json=userReserved,proto3" json:"user_reserved,omitempty"`
}

func (x *ListUserReservedResponse) Reset() {
	*x = ListUserReservedResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_user_reserved_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListUserReservedResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListUserReservedResponse) ProtoMessage() {}

func (x *ListUserReservedResponse) ProtoReflect() protoreflect.Message {
	mi := &file_iam_user_reserved_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListUserReservedResponse.ProtoReflect.Descriptor instead.
func (*ListUserReservedResponse) Descriptor() ([]byte, []int) {
	return file_iam_user_reserved_proto_rawDescGZIP(), []int{1}
}

func (x *ListUserReservedResponse) GetUserReserved() []*UserReserved {
	if x != nil {
		return x.UserReserved
	}
	return nil
}

type PutUserReservedRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId    uint32                 `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	UserReserved *UserReservedForUpsert `protobuf:"bytes,2,opt,name=user_reserved,json=userReserved,proto3" json:"user_reserved,omitempty"`
}

func (x *PutUserReservedRequest) Reset() {
	*x = PutUserReservedRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_user_reserved_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutUserReservedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutUserReservedRequest) ProtoMessage() {}

func (x *PutUserReservedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_user_reserved_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutUserReservedRequest.ProtoReflect.Descriptor instead.
func (*PutUserReservedRequest) Descriptor() ([]byte, []int) {
	return file_iam_user_reserved_proto_rawDescGZIP(), []int{2}
}

func (x *PutUserReservedRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *PutUserReservedRequest) GetUserReserved() *UserReservedForUpsert {
	if x != nil {
		return x.UserReserved
	}
	return nil
}

type PutUserReservedResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserReserved *UserReserved `protobuf:"bytes,1,opt,name=user_reserved,json=userReserved,proto3" json:"user_reserved,omitempty"`
}

func (x *PutUserReservedResponse) Reset() {
	*x = PutUserReservedResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_user_reserved_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutUserReservedResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutUserReservedResponse) ProtoMessage() {}

func (x *PutUserReservedResponse) ProtoReflect() protoreflect.Message {
	mi := &file_iam_user_reserved_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutUserReservedResponse.ProtoReflect.Descriptor instead.
func (*PutUserReservedResponse) Descriptor() ([]byte, []int) {
	return file_iam_user_reserved_proto_rawDescGZIP(), []int{3}
}

func (x *PutUserReservedResponse) GetUserReserved() *UserReserved {
	if x != nil {
		return x.UserReserved
	}
	return nil
}

type DeleteUserReservedRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId  uint32 `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	ReservedId uint32 `protobuf:"varint,2,opt,name=reserved_id,json=reservedId,proto3" json:"reserved_id,omitempty"`
}

func (x *DeleteUserReservedRequest) Reset() {
	*x = DeleteUserReservedRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_user_reserved_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteUserReservedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUserReservedRequest) ProtoMessage() {}

func (x *DeleteUserReservedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_user_reserved_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUserReservedRequest.ProtoReflect.Descriptor instead.
func (*DeleteUserReservedRequest) Descriptor() ([]byte, []int) {
	return file_iam_user_reserved_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteUserReservedRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *DeleteUserReservedRequest) GetReservedId() uint32 {
	if x != nil {
		return x.ReservedId
	}
	return 0
}

var File_iam_user_reserved_proto protoreflect.FileDescriptor

var file_iam_user_reserved_proto_rawDesc = []byte{
	0x0a, 0x17, 0x69, 0x61, 0x6d, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x63, 0x6f, 0x72, 0x65, 0x2e,
	0x69, 0x61, 0x6d, 0x1a, 0x10, 0x69, 0x61, 0x6d, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5a, 0x0a, 0x17, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65,
	0x72, 0x52, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12,
	0x20, 0x0a, 0x0c, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x70, 0x5f, 0x6b, 0x65, 0x79, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x70, 0x4b, 0x65,
	0x79, 0x22, 0x57, 0x0a, 0x18, 0x4c, 0x69, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a,
	0x0d, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x52, 0x0c, 0x75, 0x73,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x22, 0x7d, 0x0a, 0x16, 0x50, 0x75,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x49, 0x64, 0x12, 0x44, 0x0a, 0x0d, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x64, 0x46, 0x6f, 0x72, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x52, 0x0c, 0x75, 0x73, 0x65,
	0x72, 0x52, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x22, 0x56, 0x0a, 0x17, 0x50, 0x75, 0x74,
	0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a, 0x0d, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x64, 0x52, 0x0c, 0x75, 0x73, 0x65, 0x72, 0x52, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x64, 0x22, 0x5b, 0x0a, 0x19, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52,
	0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d,
	0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x1f, 0x0a,
	0x0b, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x0a, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x49, 0x64, 0x42, 0x25,
	0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x61, 0x2d,
	0x72, 0x69, 0x73, 0x6b, 0x65, 0x6e, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x69, 0x61, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_iam_user_reserved_proto_rawDescOnce sync.Once
	file_iam_user_reserved_proto_rawDescData = file_iam_user_reserved_proto_rawDesc
)

func file_iam_user_reserved_proto_rawDescGZIP() []byte {
	file_iam_user_reserved_proto_rawDescOnce.Do(func() {
		file_iam_user_reserved_proto_rawDescData = protoimpl.X.CompressGZIP(file_iam_user_reserved_proto_rawDescData)
	})
	return file_iam_user_reserved_proto_rawDescData
}

var file_iam_user_reserved_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_iam_user_reserved_proto_goTypes = []interface{}{
	(*ListUserReservedRequest)(nil),   // 0: core.iam.ListUserReservedRequest
	(*ListUserReservedResponse)(nil),  // 1: core.iam.ListUserReservedResponse
	(*PutUserReservedRequest)(nil),    // 2: core.iam.PutUserReservedRequest
	(*PutUserReservedResponse)(nil),   // 3: core.iam.PutUserReservedResponse
	(*DeleteUserReservedRequest)(nil), // 4: core.iam.DeleteUserReservedRequest
	(*UserReserved)(nil),              // 5: core.iam.UserReserved
	(*UserReservedForUpsert)(nil),     // 6: core.iam.UserReservedForUpsert
}
var file_iam_user_reserved_proto_depIdxs = []int32{
	5, // 0: core.iam.ListUserReservedResponse.user_reserved:type_name -> core.iam.UserReserved
	6, // 1: core.iam.PutUserReservedRequest.user_reserved:type_name -> core.iam.UserReservedForUpsert
	5, // 2: core.iam.PutUserReservedResponse.user_reserved:type_name -> core.iam.UserReserved
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_iam_user_reserved_proto_init() }
func file_iam_user_reserved_proto_init() {
	if File_iam_user_reserved_proto != nil {
		return
	}
	file_iam_entity_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_iam_user_reserved_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListUserReservedRequest); i {
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
		file_iam_user_reserved_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListUserReservedResponse); i {
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
		file_iam_user_reserved_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutUserReservedRequest); i {
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
		file_iam_user_reserved_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutUserReservedResponse); i {
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
		file_iam_user_reserved_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteUserReservedRequest); i {
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
			RawDescriptor: file_iam_user_reserved_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_iam_user_reserved_proto_goTypes,
		DependencyIndexes: file_iam_user_reserved_proto_depIdxs,
		MessageInfos:      file_iam_user_reserved_proto_msgTypes,
	}.Build()
	File_iam_user_reserved_proto = out.File
	file_iam_user_reserved_proto_rawDesc = nil
	file_iam_user_reserved_proto_goTypes = nil
	file_iam_user_reserved_proto_depIdxs = nil
}
