// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: iam/policy.proto

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

type ListPolicyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId uint32 `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	Name      string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	RoleId    uint32 `protobuf:"varint,3,opt,name=role_id,json=roleId,proto3" json:"role_id,omitempty"`
}

func (x *ListPolicyRequest) Reset() {
	*x = ListPolicyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_policy_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListPolicyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPolicyRequest) ProtoMessage() {}

func (x *ListPolicyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_policy_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPolicyRequest.ProtoReflect.Descriptor instead.
func (*ListPolicyRequest) Descriptor() ([]byte, []int) {
	return file_iam_policy_proto_rawDescGZIP(), []int{0}
}

func (x *ListPolicyRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *ListPolicyRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ListPolicyRequest) GetRoleId() uint32 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

type ListPolicyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PolicyId []uint32 `protobuf:"varint,1,rep,packed,name=policy_id,json=policyId,proto3" json:"policy_id,omitempty"`
}

func (x *ListPolicyResponse) Reset() {
	*x = ListPolicyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_policy_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListPolicyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPolicyResponse) ProtoMessage() {}

func (x *ListPolicyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_iam_policy_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPolicyResponse.ProtoReflect.Descriptor instead.
func (*ListPolicyResponse) Descriptor() ([]byte, []int) {
	return file_iam_policy_proto_rawDescGZIP(), []int{1}
}

func (x *ListPolicyResponse) GetPolicyId() []uint32 {
	if x != nil {
		return x.PolicyId
	}
	return nil
}

type GetPolicyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId uint32 `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	PolicyId  uint32 `protobuf:"varint,2,opt,name=policy_id,json=policyId,proto3" json:"policy_id,omitempty"`
}

func (x *GetPolicyRequest) Reset() {
	*x = GetPolicyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_policy_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPolicyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPolicyRequest) ProtoMessage() {}

func (x *GetPolicyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_policy_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPolicyRequest.ProtoReflect.Descriptor instead.
func (*GetPolicyRequest) Descriptor() ([]byte, []int) {
	return file_iam_policy_proto_rawDescGZIP(), []int{2}
}

func (x *GetPolicyRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *GetPolicyRequest) GetPolicyId() uint32 {
	if x != nil {
		return x.PolicyId
	}
	return 0
}

type GetPolicyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Policy *Policy `protobuf:"bytes,1,opt,name=policy,proto3" json:"policy,omitempty"`
}

func (x *GetPolicyResponse) Reset() {
	*x = GetPolicyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_policy_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPolicyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPolicyResponse) ProtoMessage() {}

func (x *GetPolicyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_iam_policy_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPolicyResponse.ProtoReflect.Descriptor instead.
func (*GetPolicyResponse) Descriptor() ([]byte, []int) {
	return file_iam_policy_proto_rawDescGZIP(), []int{3}
}

func (x *GetPolicyResponse) GetPolicy() *Policy {
	if x != nil {
		return x.Policy
	}
	return nil
}

type PutPolicyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId uint32           `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	Policy    *PolicyForUpsert `protobuf:"bytes,2,opt,name=policy,proto3" json:"policy,omitempty"`
}

func (x *PutPolicyRequest) Reset() {
	*x = PutPolicyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_policy_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutPolicyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutPolicyRequest) ProtoMessage() {}

func (x *PutPolicyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_policy_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutPolicyRequest.ProtoReflect.Descriptor instead.
func (*PutPolicyRequest) Descriptor() ([]byte, []int) {
	return file_iam_policy_proto_rawDescGZIP(), []int{4}
}

func (x *PutPolicyRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *PutPolicyRequest) GetPolicy() *PolicyForUpsert {
	if x != nil {
		return x.Policy
	}
	return nil
}

type PutPolicyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Policy *Policy `protobuf:"bytes,1,opt,name=policy,proto3" json:"policy,omitempty"`
}

func (x *PutPolicyResponse) Reset() {
	*x = PutPolicyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_policy_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutPolicyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutPolicyResponse) ProtoMessage() {}

func (x *PutPolicyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_iam_policy_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutPolicyResponse.ProtoReflect.Descriptor instead.
func (*PutPolicyResponse) Descriptor() ([]byte, []int) {
	return file_iam_policy_proto_rawDescGZIP(), []int{5}
}

func (x *PutPolicyResponse) GetPolicy() *Policy {
	if x != nil {
		return x.Policy
	}
	return nil
}

type DeletePolicyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId uint32 `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	PolicyId  uint32 `protobuf:"varint,2,opt,name=policy_id,json=policyId,proto3" json:"policy_id,omitempty"`
}

func (x *DeletePolicyRequest) Reset() {
	*x = DeletePolicyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_policy_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeletePolicyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeletePolicyRequest) ProtoMessage() {}

func (x *DeletePolicyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_policy_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeletePolicyRequest.ProtoReflect.Descriptor instead.
func (*DeletePolicyRequest) Descriptor() ([]byte, []int) {
	return file_iam_policy_proto_rawDescGZIP(), []int{6}
}

func (x *DeletePolicyRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *DeletePolicyRequest) GetPolicyId() uint32 {
	if x != nil {
		return x.PolicyId
	}
	return 0
}

type AttachPolicyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId uint32 `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	RoleId    uint32 `protobuf:"varint,2,opt,name=role_id,json=roleId,proto3" json:"role_id,omitempty"`
	PolicyId  uint32 `protobuf:"varint,3,opt,name=policy_id,json=policyId,proto3" json:"policy_id,omitempty"`
}

func (x *AttachPolicyRequest) Reset() {
	*x = AttachPolicyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_policy_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AttachPolicyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AttachPolicyRequest) ProtoMessage() {}

func (x *AttachPolicyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_policy_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AttachPolicyRequest.ProtoReflect.Descriptor instead.
func (*AttachPolicyRequest) Descriptor() ([]byte, []int) {
	return file_iam_policy_proto_rawDescGZIP(), []int{7}
}

func (x *AttachPolicyRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *AttachPolicyRequest) GetRoleId() uint32 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

func (x *AttachPolicyRequest) GetPolicyId() uint32 {
	if x != nil {
		return x.PolicyId
	}
	return 0
}

type AttachPolicyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RolePolicy *RolePolicy `protobuf:"bytes,1,opt,name=role_policy,json=rolePolicy,proto3" json:"role_policy,omitempty"`
}

func (x *AttachPolicyResponse) Reset() {
	*x = AttachPolicyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_policy_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AttachPolicyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AttachPolicyResponse) ProtoMessage() {}

func (x *AttachPolicyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_iam_policy_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AttachPolicyResponse.ProtoReflect.Descriptor instead.
func (*AttachPolicyResponse) Descriptor() ([]byte, []int) {
	return file_iam_policy_proto_rawDescGZIP(), []int{8}
}

func (x *AttachPolicyResponse) GetRolePolicy() *RolePolicy {
	if x != nil {
		return x.RolePolicy
	}
	return nil
}

type DetachPolicyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId uint32 `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	RoleId    uint32 `protobuf:"varint,2,opt,name=role_id,json=roleId,proto3" json:"role_id,omitempty"`
	PolicyId  uint32 `protobuf:"varint,3,opt,name=policy_id,json=policyId,proto3" json:"policy_id,omitempty"`
}

func (x *DetachPolicyRequest) Reset() {
	*x = DetachPolicyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_iam_policy_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DetachPolicyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DetachPolicyRequest) ProtoMessage() {}

func (x *DetachPolicyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_iam_policy_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DetachPolicyRequest.ProtoReflect.Descriptor instead.
func (*DetachPolicyRequest) Descriptor() ([]byte, []int) {
	return file_iam_policy_proto_rawDescGZIP(), []int{9}
}

func (x *DetachPolicyRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *DetachPolicyRequest) GetRoleId() uint32 {
	if x != nil {
		return x.RoleId
	}
	return 0
}

func (x *DetachPolicyRequest) GetPolicyId() uint32 {
	if x != nil {
		return x.PolicyId
	}
	return 0
}

var File_iam_policy_proto protoreflect.FileDescriptor

var file_iam_policy_proto_rawDesc = []byte{
	0x0a, 0x10, 0x69, 0x61, 0x6d, 0x2f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x08, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x1a, 0x10, 0x69, 0x61,
	0x6d, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5f,
	0x0a, 0x11, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x72, 0x6f, 0x6c, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x72, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x22,
	0x31, 0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x08, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x49, 0x64, 0x22, 0x4e, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79,
	0x49, 0x64, 0x22, 0x3d, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x06, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x69,
	0x61, 0x6d, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x06, 0x70, 0x6f, 0x6c, 0x69, 0x63,
	0x79, 0x22, 0x64, 0x0a, 0x10, 0x50, 0x75, 0x74, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x49, 0x64, 0x12, 0x31, 0x0a, 0x06, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e,
	0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x46, 0x6f, 0x72, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74, 0x52,
	0x06, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x22, 0x3d, 0x0a, 0x11, 0x50, 0x75, 0x74, 0x50, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x06,
	0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x63,
	0x6f, 0x72, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x06,
	0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x22, 0x51, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a,
	0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09,
	0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x08, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x49, 0x64, 0x22, 0x6a, 0x0a, 0x13, 0x41, 0x74, 0x74,
	0x61, 0x63, 0x68, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12,
	0x17, 0x0a, 0x07, 0x72, 0x6f, 0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x06, 0x72, 0x6f, 0x6c, 0x65, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x6f, 0x6c, 0x69,
	0x63, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x70, 0x6f, 0x6c,
	0x69, 0x63, 0x79, 0x49, 0x64, 0x22, 0x4d, 0x0a, 0x14, 0x41, 0x74, 0x74, 0x61, 0x63, 0x68, 0x50,
	0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a,
	0x0b, 0x72, 0x6f, 0x6c, 0x65, 0x5f, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x69, 0x61, 0x6d, 0x2e, 0x52, 0x6f,
	0x6c, 0x65, 0x50, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x52, 0x0a, 0x72, 0x6f, 0x6c, 0x65, 0x50, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x22, 0x6a, 0x0a, 0x13, 0x44, 0x65, 0x74, 0x61, 0x63, 0x68, 0x50, 0x6f,
	0x6c, 0x69, 0x63, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x72, 0x6f,
	0x6c, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x72, 0x6f, 0x6c,
	0x65, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5f, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x49, 0x64,
	0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63,
	0x61, 0x2d, 0x72, 0x69, 0x73, 0x6b, 0x65, 0x6e, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x69, 0x61, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_iam_policy_proto_rawDescOnce sync.Once
	file_iam_policy_proto_rawDescData = file_iam_policy_proto_rawDesc
)

func file_iam_policy_proto_rawDescGZIP() []byte {
	file_iam_policy_proto_rawDescOnce.Do(func() {
		file_iam_policy_proto_rawDescData = protoimpl.X.CompressGZIP(file_iam_policy_proto_rawDescData)
	})
	return file_iam_policy_proto_rawDescData
}

var file_iam_policy_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_iam_policy_proto_goTypes = []interface{}{
	(*ListPolicyRequest)(nil),    // 0: core.iam.ListPolicyRequest
	(*ListPolicyResponse)(nil),   // 1: core.iam.ListPolicyResponse
	(*GetPolicyRequest)(nil),     // 2: core.iam.GetPolicyRequest
	(*GetPolicyResponse)(nil),    // 3: core.iam.GetPolicyResponse
	(*PutPolicyRequest)(nil),     // 4: core.iam.PutPolicyRequest
	(*PutPolicyResponse)(nil),    // 5: core.iam.PutPolicyResponse
	(*DeletePolicyRequest)(nil),  // 6: core.iam.DeletePolicyRequest
	(*AttachPolicyRequest)(nil),  // 7: core.iam.AttachPolicyRequest
	(*AttachPolicyResponse)(nil), // 8: core.iam.AttachPolicyResponse
	(*DetachPolicyRequest)(nil),  // 9: core.iam.DetachPolicyRequest
	(*Policy)(nil),               // 10: core.iam.Policy
	(*PolicyForUpsert)(nil),      // 11: core.iam.PolicyForUpsert
	(*RolePolicy)(nil),           // 12: core.iam.RolePolicy
}
var file_iam_policy_proto_depIdxs = []int32{
	10, // 0: core.iam.GetPolicyResponse.policy:type_name -> core.iam.Policy
	11, // 1: core.iam.PutPolicyRequest.policy:type_name -> core.iam.PolicyForUpsert
	10, // 2: core.iam.PutPolicyResponse.policy:type_name -> core.iam.Policy
	12, // 3: core.iam.AttachPolicyResponse.role_policy:type_name -> core.iam.RolePolicy
	4,  // [4:4] is the sub-list for method output_type
	4,  // [4:4] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_iam_policy_proto_init() }
func file_iam_policy_proto_init() {
	if File_iam_policy_proto != nil {
		return
	}
	file_iam_entity_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_iam_policy_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListPolicyRequest); i {
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
		file_iam_policy_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListPolicyResponse); i {
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
		file_iam_policy_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPolicyRequest); i {
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
		file_iam_policy_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPolicyResponse); i {
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
		file_iam_policy_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutPolicyRequest); i {
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
		file_iam_policy_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutPolicyResponse); i {
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
		file_iam_policy_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeletePolicyRequest); i {
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
		file_iam_policy_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AttachPolicyRequest); i {
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
		file_iam_policy_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AttachPolicyResponse); i {
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
		file_iam_policy_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DetachPolicyRequest); i {
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
			RawDescriptor: file_iam_policy_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_iam_policy_proto_goTypes,
		DependencyIndexes: file_iam_policy_proto_depIdxs,
		MessageInfos:      file_iam_policy_proto_msgTypes,
	}.Build()
	File_iam_policy_proto = out.File
	file_iam_policy_proto_rawDesc = nil
	file_iam_policy_proto_goTypes = nil
	file_iam_policy_proto_depIdxs = nil
}
