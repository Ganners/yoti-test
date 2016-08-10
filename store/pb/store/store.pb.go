// Code generated by protoc-gen-go.
// source: store.proto
// DO NOT EDIT!

/*
Package store is a generated protocol buffer package.

It is generated from these files:
	store.proto

It has these top-level messages:
	StoreState
	ReadRequest
	ReadResponse
	WriteRequest
*/
package store

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// We want to let our siblings know our full state upon joining. It's fairly
// straightforward here
type StoreState struct {
	Data map[string][]byte `protobuf:"bytes,1,rep,name=data" json:"data,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *StoreState) Reset()                    { *m = StoreState{} }
func (m *StoreState) String() string            { return proto.CompactTextString(m) }
func (*StoreState) ProtoMessage()               {}
func (*StoreState) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *StoreState) GetData() map[string][]byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type ReadRequest struct {
	Key []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (m *ReadRequest) Reset()                    { *m = ReadRequest{} }
func (m *ReadRequest) String() string            { return proto.CompactTextString(m) }
func (*ReadRequest) ProtoMessage()               {}
func (*ReadRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type ReadResponse struct {
	Value []byte `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *ReadResponse) Reset()                    { *m = ReadResponse{} }
func (m *ReadResponse) String() string            { return proto.CompactTextString(m) }
func (*ReadResponse) ProtoMessage()               {}
func (*ReadResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type WriteRequest struct {
	Key       []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value     []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Overwrite bool   `protobuf:"varint,3,opt,name=overwrite" json:"overwrite,omitempty"`
}

func (m *WriteRequest) Reset()                    { *m = WriteRequest{} }
func (m *WriteRequest) String() string            { return proto.CompactTextString(m) }
func (*WriteRequest) ProtoMessage()               {}
func (*WriteRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func init() {
	proto.RegisterType((*StoreState)(nil), "StoreState")
	proto.RegisterType((*ReadRequest)(nil), "ReadRequest")
	proto.RegisterType((*ReadResponse)(nil), "ReadResponse")
	proto.RegisterType((*WriteRequest)(nil), "WriteRequest")
}

func init() { proto.RegisterFile("store.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 196 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x2e, 0x2e, 0xc9, 0x2f,
	0x4a, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x2a, 0xe0, 0xe2, 0x0a, 0x06, 0x71, 0x83, 0x4b,
	0x12, 0x4b, 0x52, 0x85, 0x34, 0xb9, 0x58, 0x52, 0x12, 0x4b, 0x12, 0x25, 0x18, 0x15, 0x98, 0x35,
	0xb8, 0x8d, 0x44, 0xf5, 0x10, 0x52, 0x7a, 0x2e, 0x40, 0x71, 0xd7, 0xbc, 0x92, 0xa2, 0xca, 0x20,
	0xb0, 0x12, 0x29, 0x73, 0x2e, 0x4e, 0xb8, 0x90, 0x90, 0x00, 0x17, 0x73, 0x76, 0x6a, 0x25, 0x50,
	0x1b, 0xa3, 0x06, 0x67, 0x10, 0x88, 0x29, 0x24, 0xc2, 0xc5, 0x5a, 0x96, 0x98, 0x53, 0x9a, 0x2a,
	0xc1, 0x04, 0x14, 0xe3, 0x09, 0x82, 0x70, 0xac, 0x98, 0x2c, 0x18, 0x95, 0xe4, 0xb9, 0xb8, 0x83,
	0x52, 0x13, 0x53, 0x82, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x90, 0xb5, 0xf2, 0x80, 0xb5, 0x2a,
	0xa9, 0x70, 0xf1, 0x40, 0x14, 0x14, 0x17, 0xe4, 0xe7, 0x15, 0xa7, 0x22, 0x8c, 0x62, 0x44, 0x32,
	0x4a, 0x29, 0x84, 0x8b, 0x27, 0xbc, 0x28, 0xb3, 0x24, 0x15, 0xa7, 0x39, 0xd8, 0x9d, 0x20, 0x24,
	0xc3, 0xc5, 0x99, 0x5f, 0x96, 0x5a, 0x54, 0x0e, 0xd2, 0x2b, 0xc1, 0x0c, 0x94, 0xe1, 0x08, 0x42,
	0x08, 0x24, 0xb1, 0x81, 0x43, 0xc5, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x3b, 0x5f, 0xcb, 0xa9,
	0x24, 0x01, 0x00, 0x00,
}
