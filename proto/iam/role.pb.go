// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.11
// source: iam/role.proto

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

type ListRoleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId     uint32 `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	Name          string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	UserId        uint32 `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	AccessTokenId uint32 `protobuf:"varint,4,opt,name=access_token_id,json=accessTokenId,proto3" json:"access_token_id,omitempty"`
}

func (x *ListRoleRequest) Reset() {
	*x = ListRoleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_role_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRoleRequest) ProtoMessage() {}

func (x *ListRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_role_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRoleRequest.ProtoReflect.Descriptor instead.
func (*ListRoleRequest) Descriptor() ([]byte, []int) {
	return file_iam_role_proto_rawDescGZIP(), []int{0}
}

func (x *ListRoleRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *ListRoleRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ListRoleRequest) GetUserId() uint32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ListRoleRequest) GetAccessTokenId() uint32 {
	if x != nil {
		return x.AccessTokenId
	}
	return 0
}

type ListRoleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoleId []uint32 `protobuf:"varint,1,rep,packed,name=role_id,json=roleId,proto3" json:"role_id,omitempty"`
}

func (x *ListRoleResponse) Reset() {
	*x = ListRoleResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_role_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRoleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRoleResponse) ProtoMessage() {}

func (x *ListRoleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_iam_role_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRoleResponse.ProtoReflect.Descriptor instead.
func (*ListRoleResponse) Descriptor() ([]byte, []int) {
	return file_iam_role_proto_rawDescGZIP(), []int{1}
}

func (x *ListRoleResponse) GetRoleId() []uint32 {
	if x != nil {
		return x.RoleId
	}
	return nil
}

type GetRoleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId uint32 `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	RoleId    uint32 `protobuf:"varint,2,opt,name=role_id,json=roleId,proto3" json:"role_id,omitempty"`
}

func (x *GetRoleRequest) Reset() {
	*x = GetRoleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_role_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRoleRequest) ProtoMessage() {}

func (x *GetRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_role_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRoleRequest.ProtoReflect.Descriptor instead.
func (*GetRoleRequest) Descriptor() ([]byte, []int) {
	return file_iam_role_proto_rawDescGZIP(), []int{2}
}

func (x *GetRoleRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *GetRoleRequest) GetRoleId() uint32 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

type GetRoleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Role *Role `protobuf:"bytes,1,opt,name=role,proto3" json:"role,omitempty"`
}

func (x *GetRoleResponse) Reset() {
	*x = GetRoleResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_role_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRoleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRoleResponse) ProtoMessage() {}

func (x *GetRoleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_iam_role_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRoleResponse.ProtoReflect.Descriptor instead.
func (*GetRoleResponse) Descriptor() ([]byte, []int) {
	return file_iam_role_proto_rawDescGZIP(), []int{3}
}

func (x *GetRoleResponse) GetRole() *Role {
	if x != nil {
		return x.Role
	}
	return nil
}

type PutRoleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId uint32         `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	Role      *RoleForUpsert `protobuf:"bytes,2,opt,name=role,proto3" json:"role,omitempty"`
}

func (x *PutRoleRequest) Reset() {
	*x = PutRoleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_role_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutRoleRequest) ProtoMessage() {}

func (x *PutRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_role_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutRoleRequest.ProtoReflect.Descriptor instead.
func (*PutRoleRequest) Descriptor() ([]byte, []int) {
	return file_iam_role_proto_rawDescGZIP(), []int{4}
}

func (x *PutRoleRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *PutRoleRequest) GetRole() *RoleForUpsert {
	if x != nil {
		return x.Role
	}
	return nil
}

type PutRoleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Role *Role `protobuf:"bytes,1,opt,name=role,proto3" json:"role,omitempty"`
}

func (x *PutRoleResponse) Reset() {
	*x = PutRoleResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_role_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutRoleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutRoleResponse) ProtoMessage() {}

func (x *PutRoleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_iam_role_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutRoleResponse.ProtoReflect.Descriptor instead.
func (*PutRoleResponse) Descriptor() ([]byte, []int) {
	return file_iam_role_proto_rawDescGZIP(), []int{5}
}

func (x *PutRoleResponse) GetRole() *Role {
	if x != nil {
		return x.Role
	}
	return nil
}

type DeleteRoleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId uint32 `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	RoleId    uint32 `protobuf:"varint,2,opt,name=role_id,json=roleId,proto3" json:"role_id,omitempty"`
}

func (x *DeleteRoleRequest) Reset() {
	*x = DeleteRoleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_role_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRoleRequest) ProtoMessage() {}

func (x *DeleteRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_role_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRoleRequest.ProtoReflect.Descriptor instead.
func (*DeleteRoleRequest) Descriptor() ([]byte, []int) {
	return file_iam_role_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteRoleRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *DeleteRoleRequest) GetRoleId() uint32 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

type AttachRoleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId uint32 `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	UserId    uint32 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	RoleId    uint32 `protobuf:"varint,3,opt,name=role_id,json=roleId,proto3" json:"role_id,omitempty"`
}

func (x *AttachRoleRequest) Reset() {
	*x = AttachRoleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_role_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AttachRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AttachRoleRequest) ProtoMessage() {}

func (x *AttachRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_role_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AttachRoleRequest.ProtoReflect.Descriptor instead.
func (*AttachRoleRequest) Descriptor() ([]byte, []int) {
	return file_iam_role_proto_rawDescGZIP(), []int{7}
}

func (x *AttachRoleRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *AttachRoleRequest) GetUserId() uint32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *AttachRoleRequest) GetRoleId() uint32 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

type AttachRoleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserRole *UserRole `protobuf:"bytes,1,opt,name=user_role,json=userRole,proto3" json:"user_role,omitempty"`
}

func (x *AttachRoleResponse) Reset() {
	*x = AttachRoleResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_role_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AttachRoleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AttachRoleResponse) ProtoMessage() {}

func (x *AttachRoleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_iam_role_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AttachRoleResponse.ProtoReflect.Descriptor instead.
func (*AttachRoleResponse) Descriptor() ([]byte, []int) {
	return file_iam_role_proto_rawDescGZIP(), []int{8}
}

func (x *AttachRoleResponse) GetUserRole() *UserRole {
	if x != nil {
		return x.UserRole
	}
	return nil
}

type DetachRoleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId uint32 `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	UserId    uint32 `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	RoleId    uint32 `protobuf:"varint,3,opt,name=role_id,json=roleId,proto3" json:"role_id,omitempty"`
}

func (x *DetachRoleRequest) Reset() {
	*x = DetachRoleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_role_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DetachRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DetachRoleRequest) ProtoMessage() {}

func (x *DetachRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_role_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DetachRoleRequest.ProtoReflect.Descriptor instead.
func (*DetachRoleRequest) Descriptor() ([]byte, []int) {
	return file_iam_role_proto_rawDescGZIP(), []int{9}
}

func (x *DetachRoleRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *DetachRoleRequest) GetUserId() uint32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *DetachRoleRequest) GetRoleId() uint32 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

var File_iam_role_proto protoreflect.FileDescriptor

var file_iam_role_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x69, 0x61, 0x6d, 0x2f, 0x72, 0x6f, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x1a, 0x10, 0x69, 0x61, 0x6d, 0x2f,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x85, 0x01, 0x0a,
	0x0f, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0f,
	0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0d, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x49, 0x64, 0x22, 0x2b, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x6f, 0x6c, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x72, 0x6f, 0x6c, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x06, 0x72, 0x6f, 0x6c, 0x65, 0x49,
	0x64, 0x22, 0x48, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x72, 0x6f, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x06, 0x72, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x22, 0x35, 0x0a, 0x0f, 0x47,
	0x65, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22,
	0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x04, 0x72, 0x6f,
	0x6c, 0x65, 0x22, 0x5c, 0x0a, 0x0e, 0x50, 0x75, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x49, 0x64, 0x12, 0x2b, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x17, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x52, 0x6f, 0x6c,
	0x65, 0x46, 0x6f, 0x72, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65,
	0x22, 0x35, 0x0a, 0x0f, 0x50, 0x75, 0x74, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0e, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x52, 0x6f, 0x6c,
	0x65, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x22, 0x4b, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x72,
	0x6f, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x72, 0x6f,
	0x6c, 0x65, 0x49, 0x64, 0x22, 0x64, 0x0a, 0x11, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x52, 0x6f,
	0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x17, 0x0a, 0x07, 0x72, 0x6f, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x06, 0x72, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x22, 0x45, 0x0a, 0x12, 0x41, 0x74,
	0x74, 0x61, 0x63, 0x68, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x2f, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x52, 0x6f, 0x6c,
	0x65, 0x22, 0x64, 0x0a, 0x11, 0x44, 0x65, 0x74, 0x61, 0x63, 0x68, 0x52, 0x6f, 0x6c, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x17,
	0x0a, 0x07, 0x72, 0x6f, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x06, 0x72, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x61, 0x2d, 0x72, 0x69, 0x73, 0x6b, 0x65, 0x6e, 0x2f,
	0x63, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x69, 0x61, 0x6d, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_iam_role_proto_rawDescOnce sync.Once
	file_iam_role_proto_rawDescData = file_iam_role_proto_rawDesc
)

func file_iam_role_proto_rawDescGZIP() []byte {
	file_iam_role_proto_rawDescOnce.Do(func() {
		file_iam_role_proto_rawDescData = protoimpl.X.CompressGZIP(file_iam_role_proto_rawDescData)
	})
	return file_iam_role_proto_rawDescData
}

var file_iam_role_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_iam_role_proto_goTypes = []interface{}{
	(*ListRoleRequest)(nil),    // 0: core.iam.ListRoleRequest
	(*ListRoleResponse)(nil),   // 1: core.iam.ListRoleResponse
	(*GetRoleRequest)(nil),     // 2: core.iam.GetRoleRequest
	(*GetRoleResponse)(nil),    // 3: core.iam.GetRoleResponse
	(*PutRoleRequest)(nil),     // 4: core.iam.PutRoleRequest
	(*PutRoleResponse)(nil),    // 5: core.iam.PutRoleResponse
	(*DeleteRoleRequest)(nil),  // 6: core.iam.DeleteRoleRequest
	(*AttachRoleRequest)(nil),  // 7: core.iam.AttachRoleRequest
	(*AttachRoleResponse)(nil), // 8: core.iam.AttachRoleResponse
	(*DetachRoleRequest)(nil),  // 9: core.iam.DetachRoleRequest
	(*Role)(nil),               // 10: core.iam.Role
	(*RoleForUpsert)(nil),      // 11: core.iam.RoleForUpsert
	(*UserRole)(nil),           // 12: core.iam.UserRole
}
var file_iam_role_proto_depIdxs = []int32{
	10, // 0: core.iam.GetRoleResponse.role:type_name -> core.iam.Role
	11, // 1: core.iam.PutRoleRequest.role:type_name -> core.iam.RoleForUpsert
	10, // 2: core.iam.PutRoleResponse.role:type_name -> core.iam.Role
	12, // 3: core.iam.AttachRoleResponse.user_role:type_name -> core.iam.UserRole
	4,  // [4:4] is the sub-list for method output_type
	4,  // [4:4] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_iam_role_proto_init() }
func file_iam_role_proto_init() {
	if File_iam_role_proto != nil {
		return
	}
	file_iam_entity_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_iam_role_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRoleRequest); i {
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
		file_iam_role_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRoleResponse); i {
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
		file_iam_role_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRoleRequest); i {
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
		file_iam_role_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRoleResponse); i {
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
		file_iam_role_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutRoleRequest); i {
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
		file_iam_role_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutRoleResponse); i {
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
		file_iam_role_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRoleRequest); i {
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
		file_iam_role_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AttachRoleRequest); i {
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
		file_iam_role_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AttachRoleResponse); i {
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
		file_iam_role_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DetachRoleRequest); i {
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
			RawDescriptor: file_iam_role_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_iam_role_proto_goTypes,
		DependencyIndexes: file_iam_role_proto_depIdxs,
		MessageInfos:      file_iam_role_proto_msgTypes,
	}.Build()
	File_iam_role_proto = out.File
	file_iam_role_proto_rawDesc = nil
	file_iam_role_proto_goTypes = nil
	file_iam_role_proto_depIdxs = nil
}
