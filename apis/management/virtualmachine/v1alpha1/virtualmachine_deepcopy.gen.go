// Code generated by protoc-gen-deepcopy. DO NOT EDIT.
package v1alpha1

import (
	proto "github.com/golang/protobuf/proto"
)

// DeepCopyInto supports using VirtualMachine within kubernetes types, where deepcopy-gen is used.
func (in *VirtualMachine) DeepCopyInto(out *VirtualMachine) {
	p := proto.Clone(in).(*VirtualMachine)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachine. Required by controller-gen.
func (in *VirtualMachine) DeepCopy() *VirtualMachine {
	if in == nil {
		return nil
	}
	out := new(VirtualMachine)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachine. Required by controller-gen.
func (in *VirtualMachine) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using VirtualMachine_DataVolume within kubernetes types, where deepcopy-gen is used.
func (in *VirtualMachine_DataVolume) DeepCopyInto(out *VirtualMachine_DataVolume) {
	p := proto.Clone(in).(*VirtualMachine_DataVolume)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachine_DataVolume. Required by controller-gen.
func (in *VirtualMachine_DataVolume) DeepCopy() *VirtualMachine_DataVolume {
	if in == nil {
		return nil
	}
	out := new(VirtualMachine_DataVolume)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachine_DataVolume. Required by controller-gen.
func (in *VirtualMachine_DataVolume) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using VirtualMachineConfig within kubernetes types, where deepcopy-gen is used.
func (in *VirtualMachineConfig) DeepCopyInto(out *VirtualMachineConfig) {
	p := proto.Clone(in).(*VirtualMachineConfig)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineConfig. Required by controller-gen.
func (in *VirtualMachineConfig) DeepCopy() *VirtualMachineConfig {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineConfig. Required by controller-gen.
func (in *VirtualMachineConfig) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using VirtualMachineConfig_Storage within kubernetes types, where deepcopy-gen is used.
func (in *VirtualMachineConfig_Storage) DeepCopyInto(out *VirtualMachineConfig_Storage) {
	p := proto.Clone(in).(*VirtualMachineConfig_Storage)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineConfig_Storage. Required by controller-gen.
func (in *VirtualMachineConfig_Storage) DeepCopy() *VirtualMachineConfig_Storage {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineConfig_Storage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineConfig_Storage. Required by controller-gen.
func (in *VirtualMachineConfig_Storage) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using VirtualMachineConfig_Storage_DataVolume within kubernetes types, where deepcopy-gen is used.
func (in *VirtualMachineConfig_Storage_DataVolume) DeepCopyInto(out *VirtualMachineConfig_Storage_DataVolume) {
	p := proto.Clone(in).(*VirtualMachineConfig_Storage_DataVolume)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineConfig_Storage_DataVolume. Required by controller-gen.
func (in *VirtualMachineConfig_Storage_DataVolume) DeepCopy() *VirtualMachineConfig_Storage_DataVolume {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineConfig_Storage_DataVolume)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineConfig_Storage_DataVolume. Required by controller-gen.
func (in *VirtualMachineConfig_Storage_DataVolume) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using VirtualMachineConfig_Network within kubernetes types, where deepcopy-gen is used.
func (in *VirtualMachineConfig_Network) DeepCopyInto(out *VirtualMachineConfig_Network) {
	p := proto.Clone(in).(*VirtualMachineConfig_Network)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineConfig_Network. Required by controller-gen.
func (in *VirtualMachineConfig_Network) DeepCopy() *VirtualMachineConfig_Network {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineConfig_Network)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineConfig_Network. Required by controller-gen.
func (in *VirtualMachineConfig_Network) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using VirtualMachineConfig_Compute within kubernetes types, where deepcopy-gen is used.
func (in *VirtualMachineConfig_Compute) DeepCopyInto(out *VirtualMachineConfig_Compute) {
	p := proto.Clone(in).(*VirtualMachineConfig_Compute)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineConfig_Compute. Required by controller-gen.
func (in *VirtualMachineConfig_Compute) DeepCopy() *VirtualMachineConfig_Compute {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineConfig_Compute)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineConfig_Compute. Required by controller-gen.
func (in *VirtualMachineConfig_Compute) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using VirtualMachineConfig_UserConfig within kubernetes types, where deepcopy-gen is used.
func (in *VirtualMachineConfig_UserConfig) DeepCopyInto(out *VirtualMachineConfig_UserConfig) {
	p := proto.Clone(in).(*VirtualMachineConfig_UserConfig)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineConfig_UserConfig. Required by controller-gen.
func (in *VirtualMachineConfig_UserConfig) DeepCopy() *VirtualMachineConfig_UserConfig {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineConfig_UserConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineConfig_UserConfig. Required by controller-gen.
func (in *VirtualMachineConfig_UserConfig) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using CreateVirtualMachineRequest within kubernetes types, where deepcopy-gen is used.
func (in *CreateVirtualMachineRequest) DeepCopyInto(out *CreateVirtualMachineRequest) {
	p := proto.Clone(in).(*CreateVirtualMachineRequest)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CreateVirtualMachineRequest. Required by controller-gen.
func (in *CreateVirtualMachineRequest) DeepCopy() *CreateVirtualMachineRequest {
	if in == nil {
		return nil
	}
	out := new(CreateVirtualMachineRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new CreateVirtualMachineRequest. Required by controller-gen.
func (in *CreateVirtualMachineRequest) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using DeleteVirtualMachineRequest within kubernetes types, where deepcopy-gen is used.
func (in *DeleteVirtualMachineRequest) DeepCopyInto(out *DeleteVirtualMachineRequest) {
	p := proto.Clone(in).(*DeleteVirtualMachineRequest)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeleteVirtualMachineRequest. Required by controller-gen.
func (in *DeleteVirtualMachineRequest) DeepCopy() *DeleteVirtualMachineRequest {
	if in == nil {
		return nil
	}
	out := new(DeleteVirtualMachineRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new DeleteVirtualMachineRequest. Required by controller-gen.
func (in *DeleteVirtualMachineRequest) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using DeleteVirtualMachineResponse within kubernetes types, where deepcopy-gen is used.
func (in *DeleteVirtualMachineResponse) DeepCopyInto(out *DeleteVirtualMachineResponse) {
	p := proto.Clone(in).(*DeleteVirtualMachineResponse)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeleteVirtualMachineResponse. Required by controller-gen.
func (in *DeleteVirtualMachineResponse) DeepCopy() *DeleteVirtualMachineResponse {
	if in == nil {
		return nil
	}
	out := new(DeleteVirtualMachineResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new DeleteVirtualMachineResponse. Required by controller-gen.
func (in *DeleteVirtualMachineResponse) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ListVirtualMachinesRequest within kubernetes types, where deepcopy-gen is used.
func (in *ListVirtualMachinesRequest) DeepCopyInto(out *ListVirtualMachinesRequest) {
	p := proto.Clone(in).(*ListVirtualMachinesRequest)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ListVirtualMachinesRequest. Required by controller-gen.
func (in *ListVirtualMachinesRequest) DeepCopy() *ListVirtualMachinesRequest {
	if in == nil {
		return nil
	}
	out := new(ListVirtualMachinesRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ListVirtualMachinesRequest. Required by controller-gen.
func (in *ListVirtualMachinesRequest) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ListVirtualMachinesResponse within kubernetes types, where deepcopy-gen is used.
func (in *ListVirtualMachinesResponse) DeepCopyInto(out *ListVirtualMachinesResponse) {
	p := proto.Clone(in).(*ListVirtualMachinesResponse)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ListVirtualMachinesResponse. Required by controller-gen.
func (in *ListVirtualMachinesResponse) DeepCopy() *ListVirtualMachinesResponse {
	if in == nil {
		return nil
	}
	out := new(ListVirtualMachinesResponse)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ListVirtualMachinesResponse. Required by controller-gen.
func (in *ListVirtualMachinesResponse) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using ManageVirtualMachinePowerStateRequest within kubernetes types, where deepcopy-gen is used.
func (in *ManageVirtualMachinePowerStateRequest) DeepCopyInto(out *ManageVirtualMachinePowerStateRequest) {
	p := proto.Clone(in).(*ManageVirtualMachinePowerStateRequest)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManageVirtualMachinePowerStateRequest. Required by controller-gen.
func (in *ManageVirtualMachinePowerStateRequest) DeepCopy() *ManageVirtualMachinePowerStateRequest {
	if in == nil {
		return nil
	}
	out := new(ManageVirtualMachinePowerStateRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new ManageVirtualMachinePowerStateRequest. Required by controller-gen.
func (in *ManageVirtualMachinePowerStateRequest) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}
