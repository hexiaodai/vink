// Code generated by protoc-gen-deepcopy. DO NOT EDIT.
package common

import (
	proto "github.com/golang/protobuf/proto"
)

// DeepCopyInto supports using RPCError within kubernetes types, where deepcopy-gen is used.
func (in *RPCError) DeepCopyInto(out *RPCError) {
	p := proto.Clone(in).(*RPCError)
	*out = *p
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RPCError. Required by controller-gen.
func (in *RPCError) DeepCopy() *RPCError {
	if in == nil {
		return nil
	}
	out := new(RPCError)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInterface is an autogenerated deepcopy function, copying the receiver, creating a new RPCError. Required by controller-gen.
func (in *RPCError) DeepCopyInterface() interface{} {
	return in.DeepCopy()
}
