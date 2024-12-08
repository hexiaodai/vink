// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: types/types.proto

package types

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

type ResourceType int32

const (
	ResourceType_UNSPECIFIED              ResourceType = 0
	ResourceType_VIRTUAL_MACHINE          ResourceType = 1
	ResourceType_VIRTUAL_MACHINE_INSTANCE ResourceType = 2
	ResourceType_DATA_VOLUME              ResourceType = 3
	ResourceType_NODE                     ResourceType = 4
	ResourceType_NAMESPACE                ResourceType = 5
	ResourceType_MULTUS                   ResourceType = 6
	ResourceType_SUBNET                   ResourceType = 7
	ResourceType_VPC                      ResourceType = 8
	ResourceType_IPPOOL                   ResourceType = 9
	ResourceType_STORAGE_CLASS            ResourceType = 10
	ResourceType_IPS                      ResourceType = 11
	ResourceType_VIRTUAL_MACHINE_SUMMARY  ResourceType = 12
	ResourceType_EVENT                    ResourceType = 13
)

// Enum value maps for ResourceType.
var (
	ResourceType_name = map[int32]string{
		0:  "UNSPECIFIED",
		1:  "VIRTUAL_MACHINE",
		2:  "VIRTUAL_MACHINE_INSTANCE",
		3:  "DATA_VOLUME",
		4:  "NODE",
		5:  "NAMESPACE",
		6:  "MULTUS",
		7:  "SUBNET",
		8:  "VPC",
		9:  "IPPOOL",
		10: "STORAGE_CLASS",
		11: "IPS",
		12: "VIRTUAL_MACHINE_SUMMARY",
		13: "EVENT",
	}
	ResourceType_value = map[string]int32{
		"UNSPECIFIED":              0,
		"VIRTUAL_MACHINE":          1,
		"VIRTUAL_MACHINE_INSTANCE": 2,
		"DATA_VOLUME":              3,
		"NODE":                     4,
		"NAMESPACE":                5,
		"MULTUS":                   6,
		"SUBNET":                   7,
		"VPC":                      8,
		"IPPOOL":                   9,
		"STORAGE_CLASS":            10,
		"IPS":                      11,
		"VIRTUAL_MACHINE_SUMMARY":  12,
		"EVENT":                    13,
	}
)

func (x ResourceType) Enum() *ResourceType {
	p := new(ResourceType)
	*p = x
	return p
}

func (x ResourceType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ResourceType) Descriptor() protoreflect.EnumDescriptor {
	return file_types_types_proto_enumTypes[0].Descriptor()
}

func (ResourceType) Type() protoreflect.EnumType {
	return &file_types_types_proto_enumTypes[0]
}

func (x ResourceType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ResourceType.Descriptor instead.
func (ResourceType) EnumDescriptor() ([]byte, []int) {
	return file_types_types_proto_rawDescGZIP(), []int{0}
}

type NamespaceName struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Namespace string `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Name      string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *NamespaceName) Reset() {
	*x = NamespaceName{}
	if protoimpl.UnsafeEnabled {
		mi := &file_types_types_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NamespaceName) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NamespaceName) ProtoMessage() {}

func (x *NamespaceName) ProtoReflect() protoreflect.Message {
	mi := &file_types_types_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NamespaceName.ProtoReflect.Descriptor instead.
func (*NamespaceName) Descriptor() ([]byte, []int) {
	return file_types_types_proto_rawDescGZIP(), []int{0}
}

func (x *NamespaceName) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *NamespaceName) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type FieldSelector struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FieldPath string   `protobuf:"bytes,1,opt,name=field_path,json=fieldPath,proto3" json:"field_path,omitempty"`
	Operator  string   `protobuf:"bytes,2,opt,name=operator,proto3" json:"operator,omitempty"`
	Values    []string `protobuf:"bytes,3,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *FieldSelector) Reset() {
	*x = FieldSelector{}
	if protoimpl.UnsafeEnabled {
		mi := &file_types_types_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FieldSelector) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FieldSelector) ProtoMessage() {}

func (x *FieldSelector) ProtoReflect() protoreflect.Message {
	mi := &file_types_types_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FieldSelector.ProtoReflect.Descriptor instead.
func (*FieldSelector) Descriptor() ([]byte, []int) {
	return file_types_types_proto_rawDescGZIP(), []int{1}
}

func (x *FieldSelector) GetFieldPath() string {
	if x != nil {
		return x.FieldPath
	}
	return ""
}

func (x *FieldSelector) GetOperator() string {
	if x != nil {
		return x.Operator
	}
	return ""
}

func (x *FieldSelector) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

type FieldSelectorGroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Operator       string           `protobuf:"bytes,1,opt,name=operator,proto3" json:"operator,omitempty"`
	FieldSelectors []*FieldSelector `protobuf:"bytes,2,rep,name=field_selectors,json=fieldSelectors,proto3" json:"field_selectors,omitempty"`
}

func (x *FieldSelectorGroup) Reset() {
	*x = FieldSelectorGroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_types_types_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FieldSelectorGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FieldSelectorGroup) ProtoMessage() {}

func (x *FieldSelectorGroup) ProtoReflect() protoreflect.Message {
	mi := &file_types_types_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FieldSelectorGroup.ProtoReflect.Descriptor instead.
func (*FieldSelectorGroup) Descriptor() ([]byte, []int) {
	return file_types_types_proto_rawDescGZIP(), []int{2}
}

func (x *FieldSelectorGroup) GetOperator() string {
	if x != nil {
		return x.Operator
	}
	return ""
}

func (x *FieldSelectorGroup) GetFieldSelectors() []*FieldSelector {
	if x != nil {
		return x.FieldSelectors
	}
	return nil
}

var File_types_types_proto protoreflect.FileDescriptor

var file_types_types_proto_rawDesc = []byte{
	0x0a, 0x11, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x19, 0x76, 0x69, 0x6e, 0x6b, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x76, 0x6d,
	0x2e, 0x69, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x22, 0x41,
	0x0a, 0x0d, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x22, 0x62, 0x0a, 0x0d, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x70, 0x61, 0x74, 0x68,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x50, 0x61, 0x74,
	0x68, 0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x16, 0x0a,
	0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0x83, 0x01, 0x0a, 0x12, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x53,
	0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x1a, 0x0a, 0x08,
	0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x51, 0x0a, 0x0f, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x5f, 0x73, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x28, 0x2e, 0x76, 0x69, 0x6e, 0x6b, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x76, 0x6d, 0x2e,
	0x69, 0x6f, 0x2e, 0x61, 0x70, 0x69, 0x73, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x0e, 0x66, 0x69, 0x65,
	0x6c, 0x64, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x2a, 0xed, 0x01, 0x0a, 0x0c,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0f, 0x0a, 0x0b,
	0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x13, 0x0a,
	0x0f, 0x56, 0x49, 0x52, 0x54, 0x55, 0x41, 0x4c, 0x5f, 0x4d, 0x41, 0x43, 0x48, 0x49, 0x4e, 0x45,
	0x10, 0x01, 0x12, 0x1c, 0x0a, 0x18, 0x56, 0x49, 0x52, 0x54, 0x55, 0x41, 0x4c, 0x5f, 0x4d, 0x41,
	0x43, 0x48, 0x49, 0x4e, 0x45, 0x5f, 0x49, 0x4e, 0x53, 0x54, 0x41, 0x4e, 0x43, 0x45, 0x10, 0x02,
	0x12, 0x0f, 0x0a, 0x0b, 0x44, 0x41, 0x54, 0x41, 0x5f, 0x56, 0x4f, 0x4c, 0x55, 0x4d, 0x45, 0x10,
	0x03, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x4f, 0x44, 0x45, 0x10, 0x04, 0x12, 0x0d, 0x0a, 0x09, 0x4e,
	0x41, 0x4d, 0x45, 0x53, 0x50, 0x41, 0x43, 0x45, 0x10, 0x05, 0x12, 0x0a, 0x0a, 0x06, 0x4d, 0x55,
	0x4c, 0x54, 0x55, 0x53, 0x10, 0x06, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x55, 0x42, 0x4e, 0x45, 0x54,
	0x10, 0x07, 0x12, 0x07, 0x0a, 0x03, 0x56, 0x50, 0x43, 0x10, 0x08, 0x12, 0x0a, 0x0a, 0x06, 0x49,
	0x50, 0x50, 0x4f, 0x4f, 0x4c, 0x10, 0x09, 0x12, 0x11, 0x0a, 0x0d, 0x53, 0x54, 0x4f, 0x52, 0x41,
	0x47, 0x45, 0x5f, 0x43, 0x4c, 0x41, 0x53, 0x53, 0x10, 0x0a, 0x12, 0x07, 0x0a, 0x03, 0x49, 0x50,
	0x53, 0x10, 0x0b, 0x12, 0x1b, 0x0a, 0x17, 0x56, 0x49, 0x52, 0x54, 0x55, 0x41, 0x4c, 0x5f, 0x4d,
	0x41, 0x43, 0x48, 0x49, 0x4e, 0x45, 0x5f, 0x53, 0x55, 0x4d, 0x4d, 0x41, 0x52, 0x59, 0x10, 0x0c,
	0x12, 0x09, 0x0a, 0x05, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x10, 0x0d, 0x42, 0x26, 0x5a, 0x24, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x75, 0x62, 0x65, 0x76, 0x6d,
	0x2e, 0x69, 0x6f, 0x2f, 0x76, 0x69, 0x6e, 0x6b, 0x2f, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_types_types_proto_rawDescOnce sync.Once
	file_types_types_proto_rawDescData = file_types_types_proto_rawDesc
)

func file_types_types_proto_rawDescGZIP() []byte {
	file_types_types_proto_rawDescOnce.Do(func() {
		file_types_types_proto_rawDescData = protoimpl.X.CompressGZIP(file_types_types_proto_rawDescData)
	})
	return file_types_types_proto_rawDescData
}

var file_types_types_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_types_types_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_types_types_proto_goTypes = []interface{}{
	(ResourceType)(0),          // 0: vink.kubevm.io.apis.types.ResourceType
	(*NamespaceName)(nil),      // 1: vink.kubevm.io.apis.types.NamespaceName
	(*FieldSelector)(nil),      // 2: vink.kubevm.io.apis.types.FieldSelector
	(*FieldSelectorGroup)(nil), // 3: vink.kubevm.io.apis.types.FieldSelectorGroup
}
var file_types_types_proto_depIdxs = []int32{
	2, // 0: vink.kubevm.io.apis.types.FieldSelectorGroup.field_selectors:type_name -> vink.kubevm.io.apis.types.FieldSelector
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_types_types_proto_init() }
func file_types_types_proto_init() {
	if File_types_types_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_types_types_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NamespaceName); i {
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
		file_types_types_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FieldSelector); i {
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
		file_types_types_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FieldSelectorGroup); i {
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
			RawDescriptor: file_types_types_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_types_types_proto_goTypes,
		DependencyIndexes: file_types_types_proto_depIdxs,
		EnumInfos:         file_types_types_proto_enumTypes,
		MessageInfos:      file_types_types_proto_msgTypes,
	}.Build()
	File_types_types_proto = out.File
	file_types_types_proto_rawDesc = nil
	file_types_types_proto_goTypes = nil
	file_types_types_proto_depIdxs = nil
}
