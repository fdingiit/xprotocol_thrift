// Code generated by protoc-gen-go. DO NOT EDIT.
// source: sofadrm/model/proto/BaseInfoPb.proto

package model

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

type BaseInfoPb struct {
	Zone                 string            `protobuf:"bytes,1,opt,name=zone,proto3" json:"zone,omitempty"`
	DataId               string            `protobuf:"bytes,2,opt,name=dataId,proto3" json:"dataId,omitempty"`
	Uuid                 string            `protobuf:"bytes,3,opt,name=uuid,proto3" json:"uuid,omitempty"`
	InstanceId           string            `protobuf:"bytes,4,opt,name=instanceId,proto3" json:"instanceId,omitempty"`
	Attributes           map[string]string `protobuf:"bytes,5,rep,name=attributes,proto3" json:"attributes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Profile              string            `protobuf:"bytes,6,opt,name=profile,proto3" json:"profile,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *BaseInfoPb) Reset()         { *m = BaseInfoPb{} }
func (m *BaseInfoPb) String() string { return proto.CompactTextString(m) }
func (*BaseInfoPb) ProtoMessage()    {}
func (*BaseInfoPb) Descriptor() ([]byte, []int) {
	return fileDescriptor_BaseInfoPb_23162b433a290552, []int{0}
}

func (m *BaseInfoPb) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BaseInfoPb.Unmarshal(m, b)
}

func (m *BaseInfoPb) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BaseInfoPb.Marshal(b, m, deterministic)
}

func (dst *BaseInfoPb) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BaseInfoPb.Merge(dst, src)
}

func (m *BaseInfoPb) XXX_Size() int {
	return xxx_messageInfo_BaseInfoPb.Size(m)
}

func (m *BaseInfoPb) XXX_DiscardUnknown() {
	xxx_messageInfo_BaseInfoPb.DiscardUnknown(m)
}

var xxx_messageInfo_BaseInfoPb proto.InternalMessageInfo

func (m *BaseInfoPb) GetZone() string {
	if m != nil {
		return m.Zone
	}
	return ""
}

func (m *BaseInfoPb) GetDataId() string {
	if m != nil {
		return m.DataId
	}
	return ""
}

func (m *BaseInfoPb) GetUuid() string {
	if m != nil {
		return m.Uuid
	}
	return ""
}

func (m *BaseInfoPb) GetInstanceId() string {
	if m != nil {
		return m.InstanceId
	}
	return ""
}

func (m *BaseInfoPb) GetAttributes() map[string]string {
	if m != nil {
		return m.Attributes
	}
	return nil
}

func (m *BaseInfoPb) GetProfile() string {
	if m != nil {
		return m.Profile
	}
	return ""
}

func init() {
	proto.RegisterType((*BaseInfoPb)(nil), "BaseInfoPb")
	proto.RegisterMapType((map[string]string)(nil), "BaseInfoPb.AttributesEntry")
}

func init() {
	proto.RegisterFile("sofadrm/model/proto/BaseInfoPb.proto", fileDescriptor_BaseInfoPb_23162b433a290552)
}

var fileDescriptor_BaseInfoPb_23162b433a290552 = []byte{
	// 254 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0x41, 0x4b, 0xc4, 0x30,
	0x10, 0x85, 0xe9, 0x76, 0x5b, 0x71, 0x3c, 0x28, 0x41, 0x24, 0x28, 0xc8, 0xb2, 0x78, 0xd8, 0x8b,
	0x29, 0xe8, 0x45, 0x14, 0x0f, 0x2e, 0x78, 0xe8, 0x6d, 0xd9, 0xa3, 0xb7, 0x69, 0x33, 0x85, 0x60,
	0x9a, 0x84, 0x34, 0x15, 0xea, 0xd1, 0x5f, 0x2e, 0xcd, 0x56, 0x2d, 0xde, 0xde, 0xfb, 0xe6, 0x25,
	0xcc, 0x3c, 0xb8, 0xe9, 0x6c, 0x83, 0xd2, 0xb7, 0x45, 0x6b, 0x25, 0xe9, 0xc2, 0x79, 0x1b, 0x6c,
	0xb1, 0xc5, 0x8e, 0x4a, 0xd3, 0xd8, 0x5d, 0x25, 0x22, 0x58, 0x7f, 0x2d, 0x00, 0xfe, 0x20, 0x63,
	0xb0, 0xfc, 0xb4, 0x86, 0x78, 0xb2, 0x4a, 0x36, 0xc7, 0xfb, 0xa8, 0xd9, 0x05, 0xe4, 0x12, 0x03,
	0x96, 0x92, 0x2f, 0x22, 0x9d, 0xdc, 0x98, 0xed, 0x7b, 0x25, 0x79, 0x7a, 0xc8, 0x8e, 0x9a, 0x5d,
	0x03, 0x28, 0xd3, 0x05, 0x34, 0x35, 0x95, 0x92, 0x2f, 0xe3, 0x64, 0x46, 0xd8, 0x13, 0x00, 0x86,
	0xe0, 0x55, 0xd5, 0x07, 0xea, 0x78, 0xb6, 0x4a, 0x37, 0x27, 0x77, 0x57, 0x62, 0xb6, 0xd5, 0xcb,
	0xef, 0xf4, 0xd5, 0x04, 0x3f, 0xec, 0x67, 0x71, 0xc6, 0xe1, 0xc8, 0x79, 0xdb, 0x28, 0x4d, 0x3c,
	0x8f, 0x3f, 0xff, 0xd8, 0xcb, 0x67, 0x38, 0xfd, 0xf7, 0x90, 0x9d, 0x41, 0xfa, 0x4e, 0xc3, 0x74,
	0xc8, 0x28, 0xd9, 0x39, 0x64, 0x1f, 0xa8, 0x7b, 0x9a, 0xce, 0x38, 0x98, 0xc7, 0xc5, 0x43, 0xb2,
	0xbd, 0x85, 0x75, 0x6d, 0x5b, 0x81, 0x5a, 0x39, 0x1c, 0x84, 0xf4, 0xad, 0xa8, 0xb5, 0x22, 0x13,
	0x04, 0x3a, 0x25, 0x62, 0x7d, 0xc2, 0x55, 0xbb, 0xe4, 0x2d, 0x8b, 0xba, 0xca, 0x63, 0x75, 0xf7,
	0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xe5, 0xf1, 0x57, 0xa9, 0x62, 0x01, 0x00, 0x00,
}