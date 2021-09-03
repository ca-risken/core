// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: iam/access_token.proto

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

type ListAccessTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId     uint32 `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	Name          string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	AccessTokenId uint32 `protobuf:"varint,3,opt,name=access_token_id,json=accessTokenId,proto3" json:"access_token_id,omitempty"`
}

func (x *ListAccessTokenRequest) Reset() {
	*x = ListAccessTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_access_token_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAccessTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAccessTokenRequest) ProtoMessage() {}

func (x *ListAccessTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_access_token_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAccessTokenRequest.ProtoReflect.Descriptor instead.
func (*ListAccessTokenRequest) Descriptor() ([]byte, []int) {
	return file_iam_access_token_proto_rawDescGZIP(), []int{0}
}

func (x *ListAccessTokenRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *ListAccessTokenRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ListAccessTokenRequest) GetAccessTokenId() uint32 {
	if x != nil {
		return x.AccessTokenId
	}
	return 0
}

type ListAccessTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccessToken []*AccessToken `protobuf:"bytes,1,rep,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
}

func (x *ListAccessTokenResponse) Reset() {
	*x = ListAccessTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_access_token_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAccessTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAccessTokenResponse) ProtoMessage() {}

func (x *ListAccessTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_iam_access_token_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAccessTokenResponse.ProtoReflect.Descriptor instead.
func (*ListAccessTokenResponse) Descriptor() ([]byte, []int) {
	return file_iam_access_token_proto_rawDescGZIP(), []int{1}
}

func (x *ListAccessTokenResponse) GetAccessToken() []*AccessToken {
	if x != nil {
		return x.AccessToken
	}
	return nil
}

type AuthenticateAccessTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId      uint32 `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	AccessTokenId  uint32 `protobuf:"varint,2,opt,name=access_token_id,json=accessTokenId,proto3" json:"access_token_id,omitempty"`
	PlainTextToken string `protobuf:"bytes,3,opt,name=plain_text_token,json=plainTextToken,proto3" json:"plain_text_token,omitempty"`
}

func (x *AuthenticateAccessTokenRequest) Reset() {
	*x = AuthenticateAccessTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_access_token_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthenticateAccessTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticateAccessTokenRequest) ProtoMessage() {}

func (x *AuthenticateAccessTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_access_token_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticateAccessTokenRequest.ProtoReflect.Descriptor instead.
func (*AuthenticateAccessTokenRequest) Descriptor() ([]byte, []int) {
	return file_iam_access_token_proto_rawDescGZIP(), []int{2}
}

func (x *AuthenticateAccessTokenRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *AuthenticateAccessTokenRequest) GetAccessTokenId() uint32 {
	if x != nil {
		return x.AccessTokenId
	}
	return 0
}

func (x *AuthenticateAccessTokenRequest) GetPlainTextToken() string {
	if x != nil {
		return x.PlainTextToken
	}
	return ""
}

type AuthenticateAccessTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccessToken *AccessToken `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
}

func (x *AuthenticateAccessTokenResponse) Reset() {
	*x = AuthenticateAccessTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_access_token_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthenticateAccessTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticateAccessTokenResponse) ProtoMessage() {}

func (x *AuthenticateAccessTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_iam_access_token_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticateAccessTokenResponse.ProtoReflect.Descriptor instead.
func (*AuthenticateAccessTokenResponse) Descriptor() ([]byte, []int) {
	return file_iam_access_token_proto_rawDescGZIP(), []int{3}
}

func (x *AuthenticateAccessTokenResponse) GetAccessToken() *AccessToken {
	if x != nil {
		return x.AccessToken
	}
	return nil
}

type PutAccessTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId   uint32                `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	AccessToken *AccessTokenForUpsert `protobuf:"bytes,2,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
}

func (x *PutAccessTokenRequest) Reset() {
	*x = PutAccessTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_access_token_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutAccessTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutAccessTokenRequest) ProtoMessage() {}

func (x *PutAccessTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_access_token_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutAccessTokenRequest.ProtoReflect.Descriptor instead.
func (*PutAccessTokenRequest) Descriptor() ([]byte, []int) {
	return file_iam_access_token_proto_rawDescGZIP(), []int{4}
}

func (x *PutAccessTokenRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *PutAccessTokenRequest) GetAccessToken() *AccessTokenForUpsert {
	if x != nil {
		return x.AccessToken
	}
	return nil
}

type PutAccessTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccessToken *AccessToken `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
}

func (x *PutAccessTokenResponse) Reset() {
	*x = PutAccessTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_access_token_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutAccessTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutAccessTokenResponse) ProtoMessage() {}

func (x *PutAccessTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_iam_access_token_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutAccessTokenResponse.ProtoReflect.Descriptor instead.
func (*PutAccessTokenResponse) Descriptor() ([]byte, []int) {
	return file_iam_access_token_proto_rawDescGZIP(), []int{5}
}

func (x *PutAccessTokenResponse) GetAccessToken() *AccessToken {
	if x != nil {
		return x.AccessToken
	}
	return nil
}

type DeleteAccessTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId     uint32 `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	AccessTokenId uint32 `protobuf:"varint,2,opt,name=access_token_id,json=accessTokenId,proto3" json:"access_token_id,omitempty"`
}

func (x *DeleteAccessTokenRequest) Reset() {
	*x = DeleteAccessTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_access_token_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteAccessTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteAccessTokenRequest) ProtoMessage() {}

func (x *DeleteAccessTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_access_token_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteAccessTokenRequest.ProtoReflect.Descriptor instead.
func (*DeleteAccessTokenRequest) Descriptor() ([]byte, []int) {
	return file_iam_access_token_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteAccessTokenRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *DeleteAccessTokenRequest) GetAccessTokenId() uint32 {
	if x != nil {
		return x.AccessTokenId
	}
	return 0
}

type AttachAccessTokenRoleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId     uint32 `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	AccessTokenId uint32 `protobuf:"varint,2,opt,name=access_token_id,json=accessTokenId,proto3" json:"access_token_id,omitempty"`
	RoleId        uint32 `protobuf:"varint,3,opt,name=role_id,json=roleId,proto3" json:"role_id,omitempty"`
}

func (x *AttachAccessTokenRoleRequest) Reset() {
	*x = AttachAccessTokenRoleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_access_token_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AttachAccessTokenRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AttachAccessTokenRoleRequest) ProtoMessage() {}

func (x *AttachAccessTokenRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_access_token_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AttachAccessTokenRoleRequest.ProtoReflect.Descriptor instead.
func (*AttachAccessTokenRoleRequest) Descriptor() ([]byte, []int) {
	return file_iam_access_token_proto_rawDescGZIP(), []int{7}
}

func (x *AttachAccessTokenRoleRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *AttachAccessTokenRoleRequest) GetAccessTokenId() uint32 {
	if x != nil {
		return x.AccessTokenId
	}
	return 0
}

func (x *AttachAccessTokenRoleRequest) GetRoleId() uint32 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

type AttachAccessTokenRoleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccessTokenRole *AccessTokenRole `protobuf:"bytes,1,opt,name=access_token_role,json=accessTokenRole,proto3" json:"access_token_role,omitempty"`
}

func (x *AttachAccessTokenRoleResponse) Reset() {
	*x = AttachAccessTokenRoleResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_access_token_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AttachAccessTokenRoleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AttachAccessTokenRoleResponse) ProtoMessage() {}

func (x *AttachAccessTokenRoleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_iam_access_token_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AttachAccessTokenRoleResponse.ProtoReflect.Descriptor instead.
func (*AttachAccessTokenRoleResponse) Descriptor() ([]byte, []int) {
	return file_iam_access_token_proto_rawDescGZIP(), []int{8}
}

func (x *AttachAccessTokenRoleResponse) GetAccessTokenRole() *AccessTokenRole {
	if x != nil {
		return x.AccessTokenRole
	}
	return nil
}

type DetachAccessTokenRoleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId     uint32 `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	AccessTokenId uint32 `protobuf:"varint,2,opt,name=access_token_id,json=accessTokenId,proto3" json:"access_token_id,omitempty"`
	RoleId        uint32 `protobuf:"varint,3,opt,name=role_id,json=roleId,proto3" json:"role_id,omitempty"`
}

func (x *DetachAccessTokenRoleRequest) Reset() {
	*x = DetachAccessTokenRoleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_access_token_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DetachAccessTokenRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DetachAccessTokenRoleRequest) ProtoMessage() {}

func (x *DetachAccessTokenRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_access_token_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DetachAccessTokenRoleRequest.ProtoReflect.Descriptor instead.
func (*DetachAccessTokenRoleRequest) Descriptor() ([]byte, []int) {
	return file_iam_access_token_proto_rawDescGZIP(), []int{9}
}

func (x *DetachAccessTokenRoleRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *DetachAccessTokenRoleRequest) GetAccessTokenId() uint32 {
	if x != nil {
		return x.AccessTokenId
	}
	return 0
}

func (x *DetachAccessTokenRoleRequest) GetRoleId() uint32 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

var File_iam_access_token_proto protoreflect.FileDescriptor

var file_iam_access_token_proto_rawDesc = []byte{
	0x0a, 0x16, 0x69, 0x61, 0x6d, 0x2f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x69,
	0x61, 0x6d, 0x1a, 0x10, 0x69, 0x61, 0x6d, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x73, 0x0a, 0x16, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d,
	0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x26, 0x0a, 0x0f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0d, 0x61, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x64, 0x22, 0x53, 0x0a, 0x17, 0x4c, 0x69, 0x73,
	0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x91,
	0x01, 0x0a, 0x1e, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x41,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64,
	0x12, 0x26, 0x0a, 0x0f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0d, 0x61, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x10, 0x70, 0x6c, 0x61, 0x69,
	0x6e, 0x5f, 0x74, 0x65, 0x78, 0x74, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x70, 0x6c, 0x61, 0x69, 0x6e, 0x54, 0x65, 0x78, 0x74, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x22, 0x5b, 0x0a, 0x1f, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61,
	0x74, 0x65, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x63, 0x6f,
	0x72, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22,
	0x79, 0x0a, 0x15, 0x50, 0x75, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x41, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x46, 0x6f, 0x72, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x52, 0x0b, 0x61,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x52, 0x0a, 0x16, 0x50, 0x75,
	0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x63, 0x6f, 0x72,
	0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x61,
	0x0a, 0x18, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0f, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x0d, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x49,
	0x64, 0x22, 0x7e, 0x0a, 0x1c, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x41, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64,
	0x12, 0x26, 0x0a, 0x0f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0d, 0x61, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x72, 0x6f, 0x6c, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x72, 0x6f, 0x6c, 0x65, 0x49,
	0x64, 0x22, 0x66, 0x0a, 0x1d, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x41, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x45, 0x0a, 0x11, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x5f, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e,
	0x63, 0x6f, 0x72, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x6f, 0x6c, 0x65, 0x52, 0x0f, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x6f, 0x6c, 0x65, 0x22, 0x7e, 0x0a, 0x1c, 0x44, 0x65, 0x74,
	0x61, 0x63, 0x68, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x6f,
	0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0f, 0x61, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x0d, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x64,
	0x12, 0x17, 0x0a, 0x07, 0x72, 0x6f, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x06, 0x72, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x43, 0x79, 0x62, 0x65, 0x72, 0x41, 0x67, 0x65,
	0x6e, 0x74, 0x2f, 0x6d, 0x69, 0x6d, 0x6f, 0x73, 0x61, 0x2d, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x69, 0x61, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_iam_access_token_proto_rawDescOnce sync.Once
	file_iam_access_token_proto_rawDescData = file_iam_access_token_proto_rawDesc
)

func file_iam_access_token_proto_rawDescGZIP() []byte {
	file_iam_access_token_proto_rawDescOnce.Do(func() {
		file_iam_access_token_proto_rawDescData = protoimpl.X.CompressGZIP(file_iam_access_token_proto_rawDescData)
	})
	return file_iam_access_token_proto_rawDescData
}

var file_iam_access_token_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_iam_access_token_proto_goTypes = []interface{}{
	(*ListAccessTokenRequest)(nil),          // 0: core.iam.ListAccessTokenRequest
	(*ListAccessTokenResponse)(nil),         // 1: core.iam.ListAccessTokenResponse
	(*AuthenticateAccessTokenRequest)(nil),  // 2: core.iam.AuthenticateAccessTokenRequest
	(*AuthenticateAccessTokenResponse)(nil), // 3: core.iam.AuthenticateAccessTokenResponse
	(*PutAccessTokenRequest)(nil),           // 4: core.iam.PutAccessTokenRequest
	(*PutAccessTokenResponse)(nil),          // 5: core.iam.PutAccessTokenResponse
	(*DeleteAccessTokenRequest)(nil),        // 6: core.iam.DeleteAccessTokenRequest
	(*AttachAccessTokenRoleRequest)(nil),    // 7: core.iam.AttachAccessTokenRoleRequest
	(*AttachAccessTokenRoleResponse)(nil),   // 8: core.iam.AttachAccessTokenRoleResponse
	(*DetachAccessTokenRoleRequest)(nil),    // 9: core.iam.DetachAccessTokenRoleRequest
	(*AccessToken)(nil),                     // 10: core.iam.AccessToken
	(*AccessTokenForUpsert)(nil),            // 11: core.iam.AccessTokenForUpsert
	(*AccessTokenRole)(nil),                 // 12: core.iam.AccessTokenRole
}
var file_iam_access_token_proto_depIdxs = []int32{
	10, // 0: core.iam.ListAccessTokenResponse.access_token:type_name -> core.iam.AccessToken
	10, // 1: core.iam.AuthenticateAccessTokenResponse.access_token:type_name -> core.iam.AccessToken
	11, // 2: core.iam.PutAccessTokenRequest.access_token:type_name -> core.iam.AccessTokenForUpsert
	10, // 3: core.iam.PutAccessTokenResponse.access_token:type_name -> core.iam.AccessToken
	12, // 4: core.iam.AttachAccessTokenRoleResponse.access_token_role:type_name -> core.iam.AccessTokenRole
	5,  // [5:5] is the sub-list for method output_type
	5,  // [5:5] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_iam_access_token_proto_init() }
func file_iam_access_token_proto_init() {
	if File_iam_access_token_proto != nil {
		return
	}
	file_iam_entity_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_iam_access_token_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAccessTokenRequest); i {
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
		file_iam_access_token_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAccessTokenResponse); i {
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
		file_iam_access_token_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthenticateAccessTokenRequest); i {
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
		file_iam_access_token_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthenticateAccessTokenResponse); i {
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
		file_iam_access_token_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutAccessTokenRequest); i {
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
		file_iam_access_token_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutAccessTokenResponse); i {
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
		file_iam_access_token_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteAccessTokenRequest); i {
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
		file_iam_access_token_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AttachAccessTokenRoleRequest); i {
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
		file_iam_access_token_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AttachAccessTokenRoleResponse); i {
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
		file_iam_access_token_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DetachAccessTokenRoleRequest); i {
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
			RawDescriptor: file_iam_access_token_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_iam_access_token_proto_goTypes,
		DependencyIndexes: file_iam_access_token_proto_depIdxs,
		MessageInfos:      file_iam_access_token_proto_msgTypes,
	}.Build()
	File_iam_access_token_proto = out.File
	file_iam_access_token_proto_rawDesc = nil
	file_iam_access_token_proto_goTypes = nil
	file_iam_access_token_proto_depIdxs = nil
}