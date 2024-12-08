// Code generated by protoc-gen-deepcopy. DO NOT EDIT.
package types

import (
	proto "github.com/golang/protobuf/proto"
)

// DeepCopyInto supports using NamespaceName within kubernetes types, where deepcopy-gen is used.
func (in *NamespaceName) DeepCopyInto(out *NamespaceName) {
	p := proto.Clone(in).(*NamespaceName)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespaceName. Required by controller-gen.
func (in *NamespaceName) DeepCopy() *NamespaceName {
	if in == nil {
		return nil
	}
	out := new(NamespaceName)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new NamespaceName. Required by controller-gen.
func (in *NamespaceName) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using FieldSelector within kubernetes types, where deepcopy-gen is used.
func (in *FieldSelector) DeepCopyInto(out *FieldSelector) {
	p := proto.Clone(in).(*FieldSelector)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FieldSelector. Required by controller-gen.
func (in *FieldSelector) DeepCopy() *FieldSelector {
	if in == nil {
		return nil
	}
	out := new(FieldSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new FieldSelector. Required by controller-gen.
func (in *FieldSelector) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}

// DeepCopyInto supports using FieldSelectorGroup within kubernetes types, where deepcopy-gen is used.
func (in *FieldSelectorGroup) DeepCopyInto(out *FieldSelectorGroup) {
	p := proto.Clone(in).(*FieldSelectorGroup)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FieldSelectorGroup. Required by controller-gen.
func (in *FieldSelectorGroup) DeepCopy() *FieldSelectorGroup {
	if in == nil {
		return nil
	}
	out := new(FieldSelectorGroup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new FieldSelectorGroup. Required by controller-gen.
func (in *FieldSelectorGroup) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}
