// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: report/service.proto

package report

import (
	context "context"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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

type GetReportFindingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId  uint32   `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	FromDate   string   `protobuf:"bytes,2,opt,name=from_date,json=fromDate,proto3" json:"from_date,omitempty"`
	ToDate     string   `protobuf:"bytes,3,opt,name=to_date,json=toDate,proto3" json:"to_date,omitempty"`
	Score      float32  `protobuf:"fixed32,4,opt,name=score,proto3" json:"score,omitempty"`
	DataSource []string `protobuf:"bytes,5,rep,name=data_source,json=dataSource,proto3" json:"data_source,omitempty"`
}

func (x *GetReportFindingRequest) Reset() {
	*x = GetReportFindingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_report_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetReportFindingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetReportFindingRequest) ProtoMessage() {}

func (x *GetReportFindingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_report_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetReportFindingRequest.ProtoReflect.Descriptor instead.
func (*GetReportFindingRequest) Descriptor() ([]byte, []int) {
	return file_report_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetReportFindingRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *GetReportFindingRequest) GetFromDate() string {
	if x != nil {
		return x.FromDate
	}
	return ""
}

func (x *GetReportFindingRequest) GetToDate() string {
	if x != nil {
		return x.ToDate
	}
	return ""
}

func (x *GetReportFindingRequest) GetScore() float32 {
	if x != nil {
		return x.Score
	}
	return 0
}

func (x *GetReportFindingRequest) GetDataSource() []string {
	if x != nil {
		return x.DataSource
	}
	return nil
}

type GetReportFindingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReportFinding []*ReportFinding `protobuf:"bytes,1,rep,name=report_finding,json=reportFinding,proto3" json:"report_finding,omitempty"`
}

func (x *GetReportFindingResponse) Reset() {
	*x = GetReportFindingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_report_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetReportFindingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetReportFindingResponse) ProtoMessage() {}

func (x *GetReportFindingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_report_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetReportFindingResponse.ProtoReflect.Descriptor instead.
func (*GetReportFindingResponse) Descriptor() ([]byte, []int) {
	return file_report_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetReportFindingResponse) GetReportFinding() []*ReportFinding {
	if x != nil {
		return x.ReportFinding
	}
	return nil
}

type GetReportFindingAllRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProjectId  uint32   `protobuf:"varint,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	FromDate   string   `protobuf:"bytes,2,opt,name=from_date,json=fromDate,proto3" json:"from_date,omitempty"`
	ToDate     string   `protobuf:"bytes,3,opt,name=to_date,json=toDate,proto3" json:"to_date,omitempty"`
	Score      float32  `protobuf:"fixed32,4,opt,name=score,proto3" json:"score,omitempty"`
	DataSource []string `protobuf:"bytes,5,rep,name=data_source,json=dataSource,proto3" json:"data_source,omitempty"`
}

func (x *GetReportFindingAllRequest) Reset() {
	*x = GetReportFindingAllRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_report_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetReportFindingAllRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetReportFindingAllRequest) ProtoMessage() {}

func (x *GetReportFindingAllRequest) ProtoReflect() protoreflect.Message {
	mi := &file_report_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetReportFindingAllRequest.ProtoReflect.Descriptor instead.
func (*GetReportFindingAllRequest) Descriptor() ([]byte, []int) {
	return file_report_service_proto_rawDescGZIP(), []int{2}
}

func (x *GetReportFindingAllRequest) GetProjectId() uint32 {
	if x != nil {
		return x.ProjectId
	}
	return 0
}

func (x *GetReportFindingAllRequest) GetFromDate() string {
	if x != nil {
		return x.FromDate
	}
	return ""
}

func (x *GetReportFindingAllRequest) GetToDate() string {
	if x != nil {
		return x.ToDate
	}
	return ""
}

func (x *GetReportFindingAllRequest) GetScore() float32 {
	if x != nil {
		return x.Score
	}
	return 0
}

func (x *GetReportFindingAllRequest) GetDataSource() []string {
	if x != nil {
		return x.DataSource
	}
	return nil
}

type GetReportFindingAllResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReportFinding []*ReportFinding `protobuf:"bytes,1,rep,name=report_finding,json=reportFinding,proto3" json:"report_finding,omitempty"`
}

func (x *GetReportFindingAllResponse) Reset() {
	*x = GetReportFindingAllResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_report_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetReportFindingAllResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetReportFindingAllResponse) ProtoMessage() {}

func (x *GetReportFindingAllResponse) ProtoReflect() protoreflect.Message {
	mi := &file_report_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetReportFindingAllResponse.ProtoReflect.Descriptor instead.
func (*GetReportFindingAllResponse) Descriptor() ([]byte, []int) {
	return file_report_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetReportFindingAllResponse) GetReportFinding() []*ReportFinding {
	if x != nil {
		return x.ReportFinding
	}
	return nil
}

var File_report_service_proto protoreflect.FileDescriptor

var file_report_service_proto_rawDesc = []byte{
	0x0a, 0x14, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x72, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x13, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8f,
	0x02, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x46, 0x69, 0x6e, 0x64,
	0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x26, 0x0a, 0x0a, 0x70, 0x72,
	0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x07,
	0xfa, 0x42, 0x04, 0x2a, 0x02, 0x28, 0x01, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x49, 0x64, 0x12, 0x43, 0x0a, 0x09, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x26, 0xfa, 0x42, 0x23, 0x72, 0x21, 0x32, 0x1f, 0x5e, 0x28,
	0x7c, 0x5b, 0x30, 0x2d, 0x39, 0x5d, 0x7b, 0x34, 0x7d, 0x2d, 0x5b, 0x30, 0x2d, 0x39, 0x5d, 0x7b,
	0x32, 0x7d, 0x2d, 0x5b, 0x30, 0x2d, 0x39, 0x5d, 0x7b, 0x32, 0x7d, 0x29, 0x24, 0x52, 0x08, 0x66,
	0x72, 0x6f, 0x6d, 0x44, 0x61, 0x74, 0x65, 0x12, 0x3f, 0x0a, 0x07, 0x74, 0x6f, 0x5f, 0x64, 0x61,
	0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x26, 0xfa, 0x42, 0x23, 0x72, 0x21, 0x32,
	0x1f, 0x5e, 0x28, 0x7c, 0x5b, 0x30, 0x2d, 0x39, 0x5d, 0x7b, 0x34, 0x7d, 0x2d, 0x5b, 0x30, 0x2d,
	0x39, 0x5d, 0x7b, 0x32, 0x7d, 0x2d, 0x5b, 0x30, 0x2d, 0x39, 0x5d, 0x7b, 0x32, 0x7d, 0x29, 0x24,
	0x52, 0x06, 0x74, 0x6f, 0x44, 0x61, 0x74, 0x65, 0x12, 0x25, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x72,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x42, 0x0f, 0xfa, 0x42, 0x0c, 0x0a, 0x0a, 0x1d, 0x00,
	0x00, 0x80, 0x3f, 0x2d, 0x00, 0x00, 0x00, 0x00, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x12,
	0x1f, 0x0a, 0x0b, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x05,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x22, 0x5d, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x46, 0x69, 0x6e,
	0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x0e,
	0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x66, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x72, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x46, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67,
	0x52, 0x0d, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x46, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x22,
	0x89, 0x02, 0x0a, 0x1a, 0x47, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x46, 0x69, 0x6e,
	0x64, 0x69, 0x6e, 0x67, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d,
	0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49, 0x64, 0x12, 0x43, 0x0a,
	0x09, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x26, 0xfa, 0x42, 0x23, 0x72, 0x21, 0x32, 0x1f, 0x5e, 0x28, 0x7c, 0x5b, 0x30, 0x2d, 0x39,
	0x5d, 0x7b, 0x34, 0x7d, 0x2d, 0x5b, 0x30, 0x2d, 0x39, 0x5d, 0x7b, 0x32, 0x7d, 0x2d, 0x5b, 0x30,
	0x2d, 0x39, 0x5d, 0x7b, 0x32, 0x7d, 0x29, 0x24, 0x52, 0x08, 0x66, 0x72, 0x6f, 0x6d, 0x44, 0x61,
	0x74, 0x65, 0x12, 0x3f, 0x0a, 0x07, 0x74, 0x6f, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x26, 0xfa, 0x42, 0x23, 0x72, 0x21, 0x32, 0x1f, 0x5e, 0x28, 0x7c, 0x5b,
	0x30, 0x2d, 0x39, 0x5d, 0x7b, 0x34, 0x7d, 0x2d, 0x5b, 0x30, 0x2d, 0x39, 0x5d, 0x7b, 0x32, 0x7d,
	0x2d, 0x5b, 0x30, 0x2d, 0x39, 0x5d, 0x7b, 0x32, 0x7d, 0x29, 0x24, 0x52, 0x06, 0x74, 0x6f, 0x44,
	0x61, 0x74, 0x65, 0x12, 0x25, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x02, 0x42, 0x0f, 0xfa, 0x42, 0x0c, 0x0a, 0x0a, 0x1d, 0x00, 0x00, 0x80, 0x3f, 0x2d, 0x00,
	0x00, 0x00, 0x00, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x61,
	0x74, 0x61, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x0a, 0x64, 0x61, 0x74, 0x61, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0x60, 0x0a, 0x1b, 0x47,
	0x65, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x46, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x41,
	0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x0e, 0x72, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x5f, 0x66, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74,
	0x2e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x46, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x0d,
	0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x46, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x32, 0xa2, 0x02,
	0x0a, 0x0d, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x5f, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x46, 0x69, 0x6e, 0x64,
	0x69, 0x6e, 0x67, 0x12, 0x24, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x72, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x46, 0x69, 0x6e, 0x64, 0x69,
	0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x72,
	0x74, 0x46, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x68, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x46, 0x69, 0x6e,
	0x64, 0x69, 0x6e, 0x67, 0x41, 0x6c, 0x6c, 0x12, 0x27, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x72,
	0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x46,
	0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x28, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x47,
	0x65, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x46, 0x69, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x41,
	0x6c, 0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a, 0x14, 0x43, 0x6f,
	0x6c, 0x6c, 0x65, 0x63, 0x74, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x46, 0x69, 0x6e, 0x64, 0x69,
	0x6e, 0x67, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x43, 0x79, 0x62, 0x65, 0x72, 0x41, 0x67, 0x65, 0x6e, 0x74, 0x2f, 0x6d, 0x69, 0x6d, 0x6f,
	0x73, 0x61, 0x2d, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_report_service_proto_rawDescOnce sync.Once
	file_report_service_proto_rawDescData = file_report_service_proto_rawDesc
)

func file_report_service_proto_rawDescGZIP() []byte {
	file_report_service_proto_rawDescOnce.Do(func() {
		file_report_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_report_service_proto_rawDescData)
	})
	return file_report_service_proto_rawDescData
}

var file_report_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_report_service_proto_goTypes = []interface{}{
	(*GetReportFindingRequest)(nil),     // 0: core.report.GetReportFindingRequest
	(*GetReportFindingResponse)(nil),    // 1: core.report.GetReportFindingResponse
	(*GetReportFindingAllRequest)(nil),  // 2: core.report.GetReportFindingAllRequest
	(*GetReportFindingAllResponse)(nil), // 3: core.report.GetReportFindingAllResponse
	(*ReportFinding)(nil),               // 4: core.report.ReportFinding
	(*emptypb.Empty)(nil),               // 5: google.protobuf.Empty
}
var file_report_service_proto_depIdxs = []int32{
	4, // 0: core.report.GetReportFindingResponse.report_finding:type_name -> core.report.ReportFinding
	4, // 1: core.report.GetReportFindingAllResponse.report_finding:type_name -> core.report.ReportFinding
	0, // 2: core.report.ReportService.GetReportFinding:input_type -> core.report.GetReportFindingRequest
	2, // 3: core.report.ReportService.GetReportFindingAll:input_type -> core.report.GetReportFindingAllRequest
	5, // 4: core.report.ReportService.CollectReportFinding:input_type -> google.protobuf.Empty
	1, // 5: core.report.ReportService.GetReportFinding:output_type -> core.report.GetReportFindingResponse
	3, // 6: core.report.ReportService.GetReportFindingAll:output_type -> core.report.GetReportFindingAllResponse
	5, // 7: core.report.ReportService.CollectReportFinding:output_type -> google.protobuf.Empty
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_report_service_proto_init() }
func file_report_service_proto_init() {
	if File_report_service_proto != nil {
		return
	}
	file_report_entity_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_report_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetReportFindingRequest); i {
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
		file_report_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetReportFindingResponse); i {
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
		file_report_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetReportFindingAllRequest); i {
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
		file_report_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetReportFindingAllResponse); i {
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
			RawDescriptor: file_report_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_report_service_proto_goTypes,
		DependencyIndexes: file_report_service_proto_depIdxs,
		MessageInfos:      file_report_service_proto_msgTypes,
	}.Build()
	File_report_service_proto = out.File
	file_report_service_proto_rawDesc = nil
	file_report_service_proto_goTypes = nil
	file_report_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// ReportServiceClient is the client API for ReportService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ReportServiceClient interface {
	// report
	GetReportFinding(ctx context.Context, in *GetReportFindingRequest, opts ...grpc.CallOption) (*GetReportFindingResponse, error)
	GetReportFindingAll(ctx context.Context, in *GetReportFindingAllRequest, opts ...grpc.CallOption) (*GetReportFindingAllResponse, error)
	CollectReportFinding(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type reportServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewReportServiceClient(cc grpc.ClientConnInterface) ReportServiceClient {
	return &reportServiceClient{cc}
}

func (c *reportServiceClient) GetReportFinding(ctx context.Context, in *GetReportFindingRequest, opts ...grpc.CallOption) (*GetReportFindingResponse, error) {
	out := new(GetReportFindingResponse)
	err := c.cc.Invoke(ctx, "/core.report.ReportService/GetReportFinding", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reportServiceClient) GetReportFindingAll(ctx context.Context, in *GetReportFindingAllRequest, opts ...grpc.CallOption) (*GetReportFindingAllResponse, error) {
	out := new(GetReportFindingAllResponse)
	err := c.cc.Invoke(ctx, "/core.report.ReportService/GetReportFindingAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reportServiceClient) CollectReportFinding(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/core.report.ReportService/CollectReportFinding", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReportServiceServer is the server API for ReportService service.
type ReportServiceServer interface {
	// report
	GetReportFinding(context.Context, *GetReportFindingRequest) (*GetReportFindingResponse, error)
	GetReportFindingAll(context.Context, *GetReportFindingAllRequest) (*GetReportFindingAllResponse, error)
	CollectReportFinding(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
}

// UnimplementedReportServiceServer can be embedded to have forward compatible implementations.
type UnimplementedReportServiceServer struct {
}

func (*UnimplementedReportServiceServer) GetReportFinding(context.Context, *GetReportFindingRequest) (*GetReportFindingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReportFinding not implemented")
}
func (*UnimplementedReportServiceServer) GetReportFindingAll(context.Context, *GetReportFindingAllRequest) (*GetReportFindingAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReportFindingAll not implemented")
}
func (*UnimplementedReportServiceServer) CollectReportFinding(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CollectReportFinding not implemented")
}

func RegisterReportServiceServer(s *grpc.Server, srv ReportServiceServer) {
	s.RegisterService(&_ReportService_serviceDesc, srv)
}

func _ReportService_GetReportFinding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReportFindingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReportServiceServer).GetReportFinding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.report.ReportService/GetReportFinding",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReportServiceServer).GetReportFinding(ctx, req.(*GetReportFindingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReportService_GetReportFindingAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReportFindingAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReportServiceServer).GetReportFindingAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.report.ReportService/GetReportFindingAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReportServiceServer).GetReportFindingAll(ctx, req.(*GetReportFindingAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReportService_CollectReportFinding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReportServiceServer).CollectReportFinding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/core.report.ReportService/CollectReportFinding",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReportServiceServer).CollectReportFinding(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _ReportService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "core.report.ReportService",
	HandlerType: (*ReportServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetReportFinding",
			Handler:    _ReportService_GetReportFinding_Handler,
		},
		{
			MethodName: "GetReportFindingAll",
			Handler:    _ReportService_GetReportFindingAll_Handler,
		},
		{
			MethodName: "CollectReportFinding",
			Handler:    _ReportService_CollectReportFinding_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "report/service.proto",
}
