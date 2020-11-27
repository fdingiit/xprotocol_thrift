// Code generated by protoc-gen-go. DO NOT EDIT.
// source: SyncConfigRequestPb.proto

package proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type SyncConfigRequestPb struct {
	DataCenter string `protobuf:"bytes,1,opt,name=dataCenter" json:"dataCenter,omitempty"`
	Zone       string `protobuf:"bytes,2,opt,name=zone" json:"zone,omitempty"`
}

func (m *SyncConfigRequestPb) Reset()                    { *m = SyncConfigRequestPb{} }
func (m *SyncConfigRequestPb) String() string            { return proto.CompactTextString(m) }
func (*SyncConfigRequestPb) ProtoMessage()               {}
func (*SyncConfigRequestPb) Descriptor() ([]byte, []int) { return fileDescriptor10, []int{0} }

func (m *SyncConfigRequestPb) GetDataCenter() string {
	if m != nil {
		return m.DataCenter
	}
	return ""
}

func (m *SyncConfigRequestPb) GetZone() string {
	if m != nil {
		return m.Zone
	}
	return ""
}

func init() {
	proto.RegisterType((*SyncConfigRequestPb)(nil), "SyncConfigRequestPb")
}

func init() { proto.RegisterFile("SyncConfigRequestPb.proto", fileDescriptor10) }

var fileDescriptor10 = []byte{
	// 149 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x0c, 0xae, 0xcc, 0x4b,
	0x76, 0xce, 0xcf, 0x4b, 0xcb, 0x4c, 0x0f, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x09, 0x48, 0xd2,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0xf2, 0xe4, 0x12, 0xc6, 0x22, 0x29, 0x24, 0xc7, 0xc5, 0x95,
	0x92, 0x58, 0x92, 0xe8, 0x9c, 0x9a, 0x57, 0x92, 0x5a, 0x24, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x19,
	0x84, 0x24, 0x22, 0x24, 0xc4, 0xc5, 0x52, 0x95, 0x9f, 0x97, 0x2a, 0xc1, 0x04, 0x96, 0x01, 0xb3,
	0x9d, 0x0c, 0xb8, 0x34, 0x92, 0xf3, 0x73, 0xf5, 0x12, 0x73, 0x32, 0x0b, 0x12, 0x2b, 0xf5, 0x8a,
	0xf3, 0xd3, 0x12, 0xf5, 0x8a, 0x52, 0xd3, 0x33, 0x8b, 0x4b, 0x8a, 0x2a, 0xf5, 0x8a, 0x53, 0x8b,
	0xca, 0x52, 0x8b, 0xf4, 0x72, 0xf3, 0x53, 0x52, 0x73, 0xf4, 0x0a, 0x92, 0x02, 0x18, 0xa3, 0x98,
	0x0a, 0x92, 0x92, 0xd8, 0xc0, 0x6e, 0x30, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x60, 0xad, 0xb2,
	0x86, 0xa0, 0x00, 0x00, 0x00,
}
