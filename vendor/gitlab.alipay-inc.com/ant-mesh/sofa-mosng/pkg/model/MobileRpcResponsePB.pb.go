// Code generated by protoc-gen-go. DO NOT EDIT.
// source: MobileRpcResponsePb.proto

package model

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type MobileRpcResponsePB struct {
	UniqueId             *string                    `protobuf:"bytes,1,opt,name=uniqueId" json:"uniqueId,omitempty"`
	ResultStatus         *int32                     `protobuf:"varint,2,opt,name=resultStatus,def=1001" json:"resultStatus,omitempty"`
	ResponseDataBytes    []byte                     `protobuf:"bytes,3,opt,name=responseDataBytes" json:"responseDataBytes,omitempty"`
	DataEncodingType     *int32                     `protobuf:"varint,4,opt,name=dataEncodingType" json:"dataEncodingType,omitempty"`
	Headers              []*ResponseDefaultMapEntry `protobuf:"bytes,5,rep,name=headers" json:"headers,omitempty"`
	Cookies              []*ResponseDefaultMapEntry `protobuf:"bytes,6,rep,name=cookies" json:"cookies,omitempty"`
	Ctx                  []*ResponseDefaultMapEntry `protobuf:"bytes,7,rep,name=ctx" json:"ctx,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *MobileRpcResponsePB) Reset()         { *m = MobileRpcResponsePB{} }
func (m *MobileRpcResponsePB) String() string { return proto.CompactTextString(m) }
func (*MobileRpcResponsePB) ProtoMessage()    {}
func (*MobileRpcResponsePB) Descriptor() ([]byte, []int) {
	return fileDescriptor_43465ba7af5ddcb3, []int{0}
}

func (m *MobileRpcResponsePB) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MobileRpcResponsePB.Unmarshal(m, b)
}
func (m *MobileRpcResponsePB) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MobileRpcResponsePB.Marshal(b, m, deterministic)
}
func (m *MobileRpcResponsePB) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MobileRpcResponsePB.Merge(m, src)
}
func (m *MobileRpcResponsePB) XXX_Size() int {
	return xxx_messageInfo_MobileRpcResponsePB.Size(m)
}
func (m *MobileRpcResponsePB) XXX_DiscardUnknown() {
	xxx_messageInfo_MobileRpcResponsePB.DiscardUnknown(m)
}

var xxx_messageInfo_MobileRpcResponsePB proto.InternalMessageInfo

const Default_MobileRpcResponsePB_ResultStatus int32 = 1001

func (m *MobileRpcResponsePB) GetUniqueId() string {
	if m != nil && m.UniqueId != nil {
		return *m.UniqueId
	}
	return ""
}

func (m *MobileRpcResponsePB) GetResultStatus() int32 {
	if m != nil && m.ResultStatus != nil {
		return *m.ResultStatus
	}
	return Default_MobileRpcResponsePB_ResultStatus
}

func (m *MobileRpcResponsePB) GetResponseDataBytes() []byte {
	if m != nil {
		return m.ResponseDataBytes
	}
	return nil
}

func (m *MobileRpcResponsePB) GetDataEncodingType() int32 {
	if m != nil && m.DataEncodingType != nil {
		return *m.DataEncodingType
	}
	return 0
}

func (m *MobileRpcResponsePB) GetHeaders() []*ResponseDefaultMapEntry {
	if m != nil {
		return m.Headers
	}
	return nil
}

func (m *MobileRpcResponsePB) GetCookies() []*ResponseDefaultMapEntry {
	if m != nil {
		return m.Cookies
	}
	return nil
}

func (m *MobileRpcResponsePB) GetCtx() []*ResponseDefaultMapEntry {
	if m != nil {
		return m.Ctx
	}
	return nil
}

type ResponseDefaultMapEntry struct {
	Key                  *string  `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Value                *string  `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ResponseDefaultMapEntry) Reset()         { *m = ResponseDefaultMapEntry{} }
func (m *ResponseDefaultMapEntry) String() string { return proto.CompactTextString(m) }
func (*ResponseDefaultMapEntry) ProtoMessage()    {}
func (*ResponseDefaultMapEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_43465ba7af5ddcb3, []int{1}
}

func (m *ResponseDefaultMapEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResponseDefaultMapEntry.Unmarshal(m, b)
}
func (m *ResponseDefaultMapEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResponseDefaultMapEntry.Marshal(b, m, deterministic)
}
func (m *ResponseDefaultMapEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResponseDefaultMapEntry.Merge(m, src)
}
func (m *ResponseDefaultMapEntry) XXX_Size() int {
	return xxx_messageInfo_ResponseDefaultMapEntry.Size(m)
}
func (m *ResponseDefaultMapEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_ResponseDefaultMapEntry.DiscardUnknown(m)
}

var xxx_messageInfo_ResponseDefaultMapEntry proto.InternalMessageInfo

func (m *ResponseDefaultMapEntry) GetKey() string {
	if m != nil && m.Key != nil {
		return *m.Key
	}
	return ""
}

func (m *ResponseDefaultMapEntry) GetValue() string {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return ""
}

func init() {
	proto.RegisterType((*MobileRpcResponsePB)(nil), "MobileRpcResponsePB")
	proto.RegisterType((*ResponseDefaultMapEntry)(nil), "ResponseDefaultMapEntry")
}

func init() { proto.RegisterFile("MobileRpcResponsePb.proto", fileDescriptor_43465ba7af5ddcb3) }

var fileDescriptor_43465ba7af5ddcb3 = []byte{
	// 266 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x90, 0x31, 0x4f, 0xc2, 0x40,
	0x14, 0x80, 0x53, 0x4b, 0xad, 0x3c, 0x30, 0xc1, 0x3a, 0x78, 0x38, 0x35, 0x4d, 0x4c, 0xea, 0x60,
	0x03, 0x26, 0x2e, 0x8e, 0x0d, 0x0c, 0x0e, 0x24, 0x06, 0x9d, 0xdc, 0x1e, 0xd7, 0x07, 0x5e, 0x28,
	0xbd, 0xf3, 0xee, 0x1d, 0xda, 0x1f, 0xee, 0x6e, 0x20, 0x38, 0x19, 0x75, 0x7e, 0xdf, 0xf7, 0xf2,
	0xbe, 0x07, 0xc3, 0x99, 0x5e, 0xa8, 0x9a, 0xe6, 0x46, 0xce, 0xc9, 0x19, 0xdd, 0x38, 0x7a, 0x5c,
	0x14, 0xc6, 0x6a, 0xd6, 0xd9, 0x67, 0x00, 0xe7, 0x3f, 0xa7, 0x65, 0x32, 0x80, 0x13, 0xdf, 0xa8,
	0x37, 0x4f, 0x0f, 0x95, 0x08, 0xd2, 0x20, 0xef, 0x26, 0x97, 0xd0, 0xb7, 0xe4, 0x7c, 0xcd, 0x4f,
	0x8c, 0xec, 0x9d, 0x38, 0x4a, 0x83, 0x3c, 0xba, 0xef, 0x8c, 0x47, 0xa3, 0x71, 0x32, 0x84, 0x33,
	0x7b, 0x70, 0x27, 0xc8, 0x58, 0xb6, 0x4c, 0x4e, 0x84, 0x69, 0x90, 0xf7, 0x13, 0x01, 0x83, 0x0a,
	0x19, 0xa7, 0x8d, 0xd4, 0x95, 0x6a, 0x56, 0xcf, 0xad, 0x21, 0xd1, 0xd9, 0xa9, 0xc9, 0x35, 0xc4,
	0xaf, 0x84, 0x15, 0x59, 0x27, 0xa2, 0x34, 0xcc, 0x7b, 0xb7, 0xa2, 0xf8, 0x3e, 0x60, 0x42, 0x4b,
	0xf4, 0x35, 0xcf, 0xd0, 0x4c, 0x1b, 0xb6, 0xed, 0x0e, 0x95, 0x5a, 0xaf, 0x15, 0x39, 0x71, 0xfc,
	0x0f, 0x7a, 0x05, 0xa1, 0xe4, 0x0f, 0x11, 0xff, 0x8d, 0x65, 0x77, 0x70, 0xf1, 0xdb, 0x86, 0x1e,
	0x84, 0x6b, 0x6a, 0x0f, 0xd5, 0xa7, 0x10, 0x6d, 0xb1, 0xf6, 0xb4, 0xcf, 0xed, 0x96, 0x37, 0x90,
	0x49, 0xbd, 0x29, 0xb0, 0x56, 0x06, 0xdb, 0x62, 0xb3, 0x7f, 0xdc, 0xea, 0xbd, 0xc0, 0x0a, 0x0d,
	0x93, 0x75, 0x64, 0xb7, 0x4a, 0xd2, 0x4b, 0xec, 0xf4, 0x12, 0xad, 0x91, 0x5f, 0x01, 0x00, 0x00,
	0xff, 0xff, 0xe0, 0x73, 0x3f, 0xaf, 0x79, 0x01, 0x00, 0x00,
}
