// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: management/resource/v1alpha1/resource.proto

package v1alpha1

import (
	_ "github.com/golang/protobuf/ptypes/struct"
	v1alpha1 "github.com/kubevm.io/vink/apis/apiextensions/v1alpha1"
	types "github.com/kubevm.io/vink/apis/types"
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

type GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupVersionResource *types.GroupVersionKind `protobuf:"bytes,1,opt,name=group_version_resource,json=groupVersionResource,proto3" json:"group_version_resource,omitempty"`
	NamespaceName        *types.NamespaceName    `protobuf:"bytes,2,opt,name=namespace_name,json=namespaceName,proto3" json:"namespace_name,omitempty"`
}

func (x *GetRequest) Reset() {
	*x = GetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_management_resource_v1alpha1_resource_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRequest) ProtoMessage() {}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_management_resource_v1alpha1_resource_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRequest.ProtoReflect.Descriptor instead.
func (*GetRequest) Descriptor() ([]byte, []int) {
	return file_management_resource_v1alpha1_resource_proto_rawDescGZIP(), []int{0}
}

func (x *GetRequest) GetGroupVersionResource() *types.GroupVersionKind {
	if x != nil {
		return x.GroupVersionResource
	}
	return nil
}

func (x *GetRequest) GetNamespaceName() *types.NamespaceName {
	if x != nil {
		return x.NamespaceName
	}
	return nil
}

type CreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupVersionResource *types.GroupVersionResourceIdentifier `protobuf:"bytes,1,opt,name=group_version_resource,json=groupVersionResource,proto3" json:"group_version_resource,omitempty"`
	Data                 string                                `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *CreateRequest) Reset() {
	*x = CreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_management_resource_v1alpha1_resource_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRequest) ProtoMessage() {}

func (x *CreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_management_resource_v1alpha1_resource_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRequest.ProtoReflect.Descriptor instead.
func (*CreateRequest) Descriptor() ([]byte, []int) {
	return file_management_resource_v1alpha1_resource_proto_rawDescGZIP(), []int{1}
}

func (x *CreateRequest) GetGroupVersionResource() *types.GroupVersionResourceIdentifier {
	if x != nil {
		return x.GroupVersionResource
	}
	return nil
}

func (x *CreateRequest) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

type UpdateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UpdateRequest) Reset() {
	*x = UpdateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_management_resource_v1alpha1_resource_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRequest) ProtoMessage() {}

func (x *UpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_management_resource_v1alpha1_resource_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRequest.ProtoReflect.Descriptor instead.
func (*UpdateRequest) Descriptor() ([]byte, []int) {
	return file_management_resource_v1alpha1_resource_proto_rawDescGZIP(), []int{2}
}

type DeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupVersionResource *types.GroupVersionResourceIdentifier `protobuf:"bytes,1,opt,name=group_version_resource,json=groupVersionResource,proto3" json:"group_version_resource,omitempty"`
	NamespaceName        *types.NamespaceName                  `protobuf:"bytes,2,opt,name=namespace_name,json=namespaceName,proto3" json:"namespace_name,omitempty"`
}

func (x *DeleteRequest) Reset() {
	*x = DeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_management_resource_v1alpha1_resource_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRequest) ProtoMessage() {}

func (x *DeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_management_resource_v1alpha1_resource_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRequest.ProtoReflect.Descriptor instead.
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return file_management_resource_v1alpha1_resource_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteRequest) GetGroupVersionResource() *types.GroupVersionResourceIdentifier {
	if x != nil {
		return x.GroupVersionResource
	}
	return nil
}

func (x *DeleteRequest) GetNamespaceName() *types.NamespaceName {
	if x != nil {
		return x.NamespaceName
	}
	return nil
}

var File_management_resource_v1alpha1_resource_proto protoreflect.FileDescriptor

var file_management_resource_v1alpha1_resource_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x30, 0x76,
	0x69, 0x6e, 0x6b, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x76, 0x6d, 0x2e, 0x69, 0x6f, 0x2e, 0x61, 0x70,
	0x69, 0x73, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x1a,
	0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x2f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x66,
	0x69, 0x65, 0x6c, 0x64, 0x5f, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x6e, 0x61, 0x6d, 0x65, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x37, 0x61, 0x70, 0x69, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2f, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x5f, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x64, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc0, 0x01, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x61, 0x0a, 0x16, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x76, 0x69, 0x6e, 0x6b, 0x2e, 0x6b, 0x75, 0x62,
	0x65, 0x76, 0x6d, 0x2e, 0x69, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x4b, 0x69,
	0x6e, 0x64, 0x52, 0x14, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x4f, 0x0a, 0x0e, 0x6e, 0x61, 0x6d, 0x65,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x28, 0x2e, 0x76, 0x69, 0x6e, 0x6b, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x76, 0x6d, 0x2e, 0x69,
	0x6f, 0x2e, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x4e, 0x61, 0x6d,
	0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x0d, 0x6e, 0x61, 0x6d, 0x65,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x94, 0x01, 0x0a, 0x0d, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x6f, 0x0a, 0x16, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x39, 0x2e, 0x76, 0x69,
	0x6e, 0x6b, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x76, 0x6d, 0x2e, 0x69, 0x6f, 0x2e, 0x61, 0x70, 0x69,
	0x73, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x52, 0x14, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x56, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x22, 0x0f, 0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0xd1, 0x01, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x6f, 0x0a, 0x16, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x39, 0x2e, 0x76, 0x69, 0x6e, 0x6b, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x76,
	0x6d, 0x2e, 0x69, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x52, 0x14,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x12, 0x4f, 0x0a, 0x0e, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x76,
	0x69, 0x6e, 0x6b, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x76, 0x6d, 0x2e, 0x69, 0x6f, 0x2e, 0x61, 0x70,
	0x69, 0x73, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x0d, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x32, 0xa7, 0x04, 0x0a, 0x12, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x89, 0x01, 0x0a,
	0x03, 0x47, 0x65, 0x74, 0x12, 0x3c, 0x2e, 0x76, 0x69, 0x6e, 0x6b, 0x2e, 0x6b, 0x75, 0x62, 0x65,
	0x76, 0x6d, 0x2e, 0x69, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x76,
	0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x44, 0x2e, 0x76, 0x69, 0x6e, 0x6b, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x76, 0x6d,
	0x2e, 0x69, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x61, 0x70, 0x69, 0x65, 0x78, 0x74, 0x65,
	0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e,
	0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x44, 0x65,
	0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x8f, 0x01, 0x0a, 0x06, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x12, 0x3f, 0x2e, 0x76, 0x69, 0x6e, 0x6b, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x76,
	0x6d, 0x2e, 0x69, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65,
	0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x76, 0x31,
	0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x44, 0x2e, 0x76, 0x69, 0x6e, 0x6b, 0x2e, 0x6b, 0x75, 0x62, 0x65,
	0x76, 0x6d, 0x2e, 0x69, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x61, 0x70, 0x69, 0x65, 0x78,
	0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61,
	0x31, 0x2e, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x8f, 0x01, 0x0a, 0x06, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x3f, 0x2e, 0x76, 0x69, 0x6e, 0x6b, 0x2e, 0x6b, 0x75, 0x62,
	0x65, 0x76, 0x6d, 0x2e, 0x69, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x6d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e,
	0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x44, 0x2e, 0x76, 0x69, 0x6e, 0x6b, 0x2e, 0x6b, 0x75,
	0x62, 0x65, 0x76, 0x6d, 0x2e, 0x69, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x61, 0x70, 0x69,
	0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0x2e, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x61, 0x0a, 0x06,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x3f, 0x2e, 0x76, 0x69, 0x6e, 0x6b, 0x2e, 0x6b, 0x75,
	0x62, 0x65, 0x76, 0x6d, 0x2e, 0x69, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x6d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42,
	0x3d, 0x5a, 0x3b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x75,
	0x62, 0x65, 0x76, 0x6d, 0x2e, 0x69, 0x6f, 0x2f, 0x76, 0x69, 0x6e, 0x6b, 0x2f, 0x61, 0x70, 0x69,
	0x73, 0x2f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_management_resource_v1alpha1_resource_proto_rawDescOnce sync.Once
	file_management_resource_v1alpha1_resource_proto_rawDescData = file_management_resource_v1alpha1_resource_proto_rawDesc
)

func file_management_resource_v1alpha1_resource_proto_rawDescGZIP() []byte {
	file_management_resource_v1alpha1_resource_proto_rawDescOnce.Do(func() {
		file_management_resource_v1alpha1_resource_proto_rawDescData = protoimpl.X.CompressGZIP(file_management_resource_v1alpha1_resource_proto_rawDescData)
	})
	return file_management_resource_v1alpha1_resource_proto_rawDescData
}

var file_management_resource_v1alpha1_resource_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_management_resource_v1alpha1_resource_proto_goTypes = []interface{}{
	(*GetRequest)(nil),                           // 0: vink.kubevm.io.apis.management.resource.v1alpha1.GetRequest
	(*CreateRequest)(nil),                        // 1: vink.kubevm.io.apis.management.resource.v1alpha1.CreateRequest
	(*UpdateRequest)(nil),                        // 2: vink.kubevm.io.apis.management.resource.v1alpha1.UpdateRequest
	(*DeleteRequest)(nil),                        // 3: vink.kubevm.io.apis.management.resource.v1alpha1.DeleteRequest
	(*types.GroupVersionKind)(nil),               // 4: vink.kubevm.io.apis.types.GroupVersionKind
	(*types.NamespaceName)(nil),                  // 5: vink.kubevm.io.apis.types.NamespaceName
	(*types.GroupVersionResourceIdentifier)(nil), // 6: vink.kubevm.io.apis.types.GroupVersionResourceIdentifier
	(*v1alpha1.CustomResourceDefinition)(nil),    // 7: vink.kubevm.io.apis.apiextensions.v1alpha1.CustomResourceDefinition
	(*emptypb.Empty)(nil),                        // 8: google.protobuf.Empty
}
var file_management_resource_v1alpha1_resource_proto_depIdxs = []int32{
	4, // 0: vink.kubevm.io.apis.management.resource.v1alpha1.GetRequest.group_version_resource:type_name -> vink.kubevm.io.apis.types.GroupVersionKind
	5, // 1: vink.kubevm.io.apis.management.resource.v1alpha1.GetRequest.namespace_name:type_name -> vink.kubevm.io.apis.types.NamespaceName
	6, // 2: vink.kubevm.io.apis.management.resource.v1alpha1.CreateRequest.group_version_resource:type_name -> vink.kubevm.io.apis.types.GroupVersionResourceIdentifier
	6, // 3: vink.kubevm.io.apis.management.resource.v1alpha1.DeleteRequest.group_version_resource:type_name -> vink.kubevm.io.apis.types.GroupVersionResourceIdentifier
	5, // 4: vink.kubevm.io.apis.management.resource.v1alpha1.DeleteRequest.namespace_name:type_name -> vink.kubevm.io.apis.types.NamespaceName
	0, // 5: vink.kubevm.io.apis.management.resource.v1alpha1.ResourceManagement.Get:input_type -> vink.kubevm.io.apis.management.resource.v1alpha1.GetRequest
	1, // 6: vink.kubevm.io.apis.management.resource.v1alpha1.ResourceManagement.Create:input_type -> vink.kubevm.io.apis.management.resource.v1alpha1.CreateRequest
	2, // 7: vink.kubevm.io.apis.management.resource.v1alpha1.ResourceManagement.Update:input_type -> vink.kubevm.io.apis.management.resource.v1alpha1.UpdateRequest
	3, // 8: vink.kubevm.io.apis.management.resource.v1alpha1.ResourceManagement.Delete:input_type -> vink.kubevm.io.apis.management.resource.v1alpha1.DeleteRequest
	7, // 9: vink.kubevm.io.apis.management.resource.v1alpha1.ResourceManagement.Get:output_type -> vink.kubevm.io.apis.apiextensions.v1alpha1.CustomResourceDefinition
	7, // 10: vink.kubevm.io.apis.management.resource.v1alpha1.ResourceManagement.Create:output_type -> vink.kubevm.io.apis.apiextensions.v1alpha1.CustomResourceDefinition
	7, // 11: vink.kubevm.io.apis.management.resource.v1alpha1.ResourceManagement.Update:output_type -> vink.kubevm.io.apis.apiextensions.v1alpha1.CustomResourceDefinition
	8, // 12: vink.kubevm.io.apis.management.resource.v1alpha1.ResourceManagement.Delete:output_type -> google.protobuf.Empty
	9, // [9:13] is the sub-list for method output_type
	5, // [5:9] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_management_resource_v1alpha1_resource_proto_init() }
func file_management_resource_v1alpha1_resource_proto_init() {
	if File_management_resource_v1alpha1_resource_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_management_resource_v1alpha1_resource_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRequest); i {
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
		file_management_resource_v1alpha1_resource_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRequest); i {
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
		file_management_resource_v1alpha1_resource_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateRequest); i {
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
		file_management_resource_v1alpha1_resource_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRequest); i {
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
			RawDescriptor: file_management_resource_v1alpha1_resource_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_management_resource_v1alpha1_resource_proto_goTypes,
		DependencyIndexes: file_management_resource_v1alpha1_resource_proto_depIdxs,
		MessageInfos:      file_management_resource_v1alpha1_resource_proto_msgTypes,
	}.Build()
	File_management_resource_v1alpha1_resource_proto = out.File
	file_management_resource_v1alpha1_resource_proto_rawDesc = nil
	file_management_resource_v1alpha1_resource_proto_goTypes = nil
	file_management_resource_v1alpha1_resource_proto_depIdxs = nil
}
