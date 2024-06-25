// Code generated by protoc-gen-deepcopy. DO NOT EDIT.
package v1alpha1

import (
	proto "github.com/golang/protobuf/proto"
)

// DeepCopyInto supports using CreateDataVolumeRequest within kubernetes types, where deepcopy-gen is used.
func (in *CreateDataVolumeRequest) DeepCopyInto(out *CreateDataVolumeRequest) {
	p := proto.Clone(in).(*CreateDataVolumeRequest)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CreateDataVolumeRequest. Required by controller-gen.
func (in *CreateDataVolumeRequest) DeepCopy() *CreateDataVolumeRequest {
	if in == nil {
		return nil
	}
	out := new(CreateDataVolumeRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new CreateDataVolumeRequest. Required by controller-gen.
func (in *CreateDataVolumeRequest) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using DeleteDataVolumeRequest within kubernetes types, where deepcopy-gen is used.
func (in *DeleteDataVolumeRequest) DeepCopyInto(out *DeleteDataVolumeRequest) {
	p := proto.Clone(in).(*DeleteDataVolumeRequest)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeleteDataVolumeRequest. Required by controller-gen.
func (in *DeleteDataVolumeRequest) DeepCopy() *DeleteDataVolumeRequest {
	if in == nil {
		return nil
	}
	out := new(DeleteDataVolumeRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new DeleteDataVolumeRequest. Required by controller-gen.
func (in *DeleteDataVolumeRequest) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using DataVolumeConfig within kubernetes types, where deepcopy-gen is used.
func (in *DataVolumeConfig) DeepCopyInto(out *DataVolumeConfig) {
	p := proto.Clone(in).(*DataVolumeConfig)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataVolumeConfig. Required by controller-gen.
func (in *DataVolumeConfig) DeepCopy() *DataVolumeConfig {
	if in == nil {
		return nil
	}
	out := new(DataVolumeConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new DataVolumeConfig. Required by controller-gen.
func (in *DataVolumeConfig) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using DataVolumeConfig_DataSource within kubernetes types, where deepcopy-gen is used.
func (in *DataVolumeConfig_DataSource) DeepCopyInto(out *DataVolumeConfig_DataSource) {
	p := proto.Clone(in).(*DataVolumeConfig_DataSource)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataVolumeConfig_DataSource. Required by controller-gen.
func (in *DataVolumeConfig_DataSource) DeepCopy() *DataVolumeConfig_DataSource {
	if in == nil {
		return nil
	}
	out := new(DataVolumeConfig_DataSource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new DataVolumeConfig_DataSource. Required by controller-gen.
func (in *DataVolumeConfig_DataSource) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using DataVolumeConfig_DataSource_Blank within kubernetes types, where deepcopy-gen is used.
func (in *DataVolumeConfig_DataSource_Blank) DeepCopyInto(out *DataVolumeConfig_DataSource_Blank) {
	p := proto.Clone(in).(*DataVolumeConfig_DataSource_Blank)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataVolumeConfig_DataSource_Blank. Required by controller-gen.
func (in *DataVolumeConfig_DataSource_Blank) DeepCopy() *DataVolumeConfig_DataSource_Blank {
	if in == nil {
		return nil
	}
	out := new(DataVolumeConfig_DataSource_Blank)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new DataVolumeConfig_DataSource_Blank. Required by controller-gen.
func (in *DataVolumeConfig_DataSource_Blank) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using DataVolumeConfig_DataSource_Upload within kubernetes types, where deepcopy-gen is used.
func (in *DataVolumeConfig_DataSource_Upload) DeepCopyInto(out *DataVolumeConfig_DataSource_Upload) {
	p := proto.Clone(in).(*DataVolumeConfig_DataSource_Upload)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataVolumeConfig_DataSource_Upload. Required by controller-gen.
func (in *DataVolumeConfig_DataSource_Upload) DeepCopy() *DataVolumeConfig_DataSource_Upload {
	if in == nil {
		return nil
	}
	out := new(DataVolumeConfig_DataSource_Upload)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new DataVolumeConfig_DataSource_Upload. Required by controller-gen.
func (in *DataVolumeConfig_DataSource_Upload) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using DataVolumeConfig_DataSource_Http within kubernetes types, where deepcopy-gen is used.
func (in *DataVolumeConfig_DataSource_Http) DeepCopyInto(out *DataVolumeConfig_DataSource_Http) {
	p := proto.Clone(in).(*DataVolumeConfig_DataSource_Http)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataVolumeConfig_DataSource_Http. Required by controller-gen.
func (in *DataVolumeConfig_DataSource_Http) DeepCopy() *DataVolumeConfig_DataSource_Http {
	if in == nil {
		return nil
	}
	out := new(DataVolumeConfig_DataSource_Http)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new DataVolumeConfig_DataSource_Http. Required by controller-gen.
func (in *DataVolumeConfig_DataSource_Http) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using DataVolumeConfig_DataSource_Registry within kubernetes types, where deepcopy-gen is used.
func (in *DataVolumeConfig_DataSource_Registry) DeepCopyInto(out *DataVolumeConfig_DataSource_Registry) {
	p := proto.Clone(in).(*DataVolumeConfig_DataSource_Registry)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataVolumeConfig_DataSource_Registry. Required by controller-gen.
func (in *DataVolumeConfig_DataSource_Registry) DeepCopy() *DataVolumeConfig_DataSource_Registry {
	if in == nil {
		return nil
	}
	out := new(DataVolumeConfig_DataSource_Registry)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new DataVolumeConfig_DataSource_Registry. Required by controller-gen.
func (in *DataVolumeConfig_DataSource_Registry) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using DataVolumeConfig_DataSource_S3 within kubernetes types, where deepcopy-gen is used.
func (in *DataVolumeConfig_DataSource_S3) DeepCopyInto(out *DataVolumeConfig_DataSource_S3) {
	p := proto.Clone(in).(*DataVolumeConfig_DataSource_S3)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataVolumeConfig_DataSource_S3. Required by controller-gen.
func (in *DataVolumeConfig_DataSource_S3) DeepCopy() *DataVolumeConfig_DataSource_S3 {
	if in == nil {
		return nil
	}
	out := new(DataVolumeConfig_DataSource_S3)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new DataVolumeConfig_DataSource_S3. Required by controller-gen.
func (in *DataVolumeConfig_DataSource_S3) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using DataVolumeConfig_BoundPVC within kubernetes types, where deepcopy-gen is used.
func (in *DataVolumeConfig_BoundPVC) DeepCopyInto(out *DataVolumeConfig_BoundPVC) {
	p := proto.Clone(in).(*DataVolumeConfig_BoundPVC)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataVolumeConfig_BoundPVC. Required by controller-gen.
func (in *DataVolumeConfig_BoundPVC) DeepCopy() *DataVolumeConfig_BoundPVC {
	if in == nil {
		return nil
	}
	out := new(DataVolumeConfig_BoundPVC)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new DataVolumeConfig_BoundPVC. Required by controller-gen.
func (in *DataVolumeConfig_BoundPVC) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using DataVolumeConfig_OperatingSystem within kubernetes types, where deepcopy-gen is used.
func (in *DataVolumeConfig_OperatingSystem) DeepCopyInto(out *DataVolumeConfig_OperatingSystem) {
	p := proto.Clone(in).(*DataVolumeConfig_OperatingSystem)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataVolumeConfig_OperatingSystem. Required by controller-gen.
func (in *DataVolumeConfig_OperatingSystem) DeepCopy() *DataVolumeConfig_OperatingSystem {
	if in == nil {
		return nil
	}
	out := new(DataVolumeConfig_OperatingSystem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new DataVolumeConfig_OperatingSystem. Required by controller-gen.
func (in *DataVolumeConfig_OperatingSystem) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using DataVolume within kubernetes types, where deepcopy-gen is used.
func (in *DataVolume) DeepCopyInto(out *DataVolume) {
	p := proto.Clone(in).(*DataVolume)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DataVolume. Required by controller-gen.
func (in *DataVolume) DeepCopy() *DataVolume {
	if in == nil {
		return nil
	}
	out := new(DataVolume)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new DataVolume. Required by controller-gen.
func (in *DataVolume) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using DeleteDataVolumeResponse within kubernetes types, where deepcopy-gen is used.
func (in *DeleteDataVolumeResponse) DeepCopyInto(out *DeleteDataVolumeResponse) {
	p := proto.Clone(in).(*DeleteDataVolumeResponse)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeleteDataVolumeResponse. Required by controller-gen.
func (in *DeleteDataVolumeResponse) DeepCopy() *DeleteDataVolumeResponse {
	if in == nil {
		return nil
	}
	out := new(DeleteDataVolumeResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new DeleteDataVolumeResponse. Required by controller-gen.
func (in *DeleteDataVolumeResponse) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ListDataVolumesRequest within kubernetes types, where deepcopy-gen is used.
func (in *ListDataVolumesRequest) DeepCopyInto(out *ListDataVolumesRequest) {
	p := proto.Clone(in).(*ListDataVolumesRequest)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ListDataVolumesRequest. Required by controller-gen.
func (in *ListDataVolumesRequest) DeepCopy() *ListDataVolumesRequest {
	if in == nil {
		return nil
	}
	out := new(ListDataVolumesRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ListDataVolumesRequest. Required by controller-gen.
func (in *ListDataVolumesRequest) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ListDataVolumesResponse within kubernetes types, where deepcopy-gen is used.
func (in *ListDataVolumesResponse) DeepCopyInto(out *ListDataVolumesResponse) {
	p := proto.Clone(in).(*ListDataVolumesResponse)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ListDataVolumesResponse. Required by controller-gen.
func (in *ListDataVolumesResponse) DeepCopy() *ListDataVolumesResponse {
	if in == nil {
		return nil
	}
	out := new(ListDataVolumesResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ListDataVolumesResponse. Required by controller-gen.
func (in *ListDataVolumesResponse) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}
