// Code generated by protoc-gen-go. DO NOT EDIT.
// source: envoy/api/v2/nds.proto

package envoy_api_v2

import (
	fmt "fmt"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/any"
	_ "github.com/golang/protobuf/ptypes/duration"
	_ "github.com/golang/protobuf/ptypes/struct"
	_ "github.com/golang/protobuf/ptypes/wrappers"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type NodeSentry struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Namespace            string   `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	InstanceIp           string   `protobuf:"bytes,3,opt,name=instance_ip,json=instanceIp,proto3" json:"instance_ip,omitempty"`
	NodeName             string   `protobuf:"bytes,4,opt,name=node_name,json=nodeName,proto3" json:"node_name,omitempty"`
	NodeIp               string   `protobuf:"bytes,5,opt,name=node_ip,json=nodeIp,proto3" json:"node_ip,omitempty"`
	ClusterName          string   `protobuf:"bytes,6,opt,name=cluster_name,json=clusterName,proto3" json:"cluster_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeSentry) Reset()         { *m = NodeSentry{} }
func (m *NodeSentry) String() string { return proto.CompactTextString(m) }
func (*NodeSentry) ProtoMessage()    {}
func (*NodeSentry) Descriptor() ([]byte, []int) {
	return fileDescriptor_19b4a8d1fb15240f, []int{0}
}

func (m *NodeSentry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeSentry.Unmarshal(m, b)
}
func (m *NodeSentry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeSentry.Marshal(b, m, deterministic)
}
func (m *NodeSentry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeSentry.Merge(m, src)
}
func (m *NodeSentry) XXX_Size() int {
	return xxx_messageInfo_NodeSentry.Size(m)
}
func (m *NodeSentry) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeSentry.DiscardUnknown(m)
}

var xxx_messageInfo_NodeSentry proto.InternalMessageInfo

func (m *NodeSentry) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *NodeSentry) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *NodeSentry) GetInstanceIp() string {
	if m != nil {
		return m.InstanceIp
	}
	return ""
}

func (m *NodeSentry) GetNodeName() string {
	if m != nil {
		return m.NodeName
	}
	return ""
}

func (m *NodeSentry) GetNodeIp() string {
	if m != nil {
		return m.NodeIp
	}
	return ""
}

func (m *NodeSentry) GetClusterName() string {
	if m != nil {
		return m.ClusterName
	}
	return ""
}

func init() {
	proto.RegisterType((*NodeSentry)(nil), "envoy.api.v2.NodeSentry")
}

func init() { proto.RegisterFile("envoy/api/v2/nds.proto", fileDescriptor_19b4a8d1fb15240f) }

var fileDescriptor_19b4a8d1fb15240f = []byte{
	// 286 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0x4f, 0x4e, 0xf3, 0x30,
	0x10, 0xc5, 0xe5, 0xef, 0x2b, 0xa5, 0x71, 0xbb, 0xf2, 0x82, 0x86, 0x10, 0xf1, 0x6f, 0x85, 0x84,
	0x94, 0x48, 0xe5, 0x06, 0xdd, 0x75, 0x13, 0x55, 0x70, 0x80, 0xca, 0x8d, 0x4d, 0x65, 0x29, 0x78,
	0x46, 0xb6, 0x13, 0xc8, 0x0d, 0x38, 0x12, 0xc7, 0x43, 0x9e, 0x26, 0x02, 0x55, 0x62, 0xe5, 0x99,
	0xf7, 0x9b, 0x37, 0xd6, 0xb3, 0xf9, 0x85, 0xb6, 0x1d, 0xf4, 0xa5, 0x44, 0x53, 0x76, 0xab, 0xd2,
	0x2a, 0x5f, 0xa0, 0x83, 0x00, 0x62, 0x41, 0x7a, 0x21, 0xd1, 0x14, 0xdd, 0x2a, 0xcb, 0x0f, 0x00,
	0x87, 0x46, 0xd3, 0x98, 0xb4, 0x16, 0x82, 0x0c, 0x06, 0xec, 0x30, 0x9b, 0x5d, 0x0e, 0x94, 0xba,
	0x7d, 0xfb, 0x5a, 0x4a, 0xdb, 0x0f, 0xe8, 0xfa, 0x14, 0xa9, 0xd6, 0x91, 0x77, 0xe0, 0xf9, 0x29,
	0xf7, 0xc1, 0xb5, 0x75, 0xf8, 0xcb, 0xfd, 0xee, 0x24, 0xa2, 0x76, 0xe3, 0xc5, 0xcb, 0x4e, 0x36,
	0x46, 0xc9, 0xa0, 0xcb, 0xb1, 0x38, 0x82, 0xfb, 0x2f, 0xc6, 0x79, 0x05, 0x4a, 0xbf, 0x68, 0x1b,
	0x5c, 0x2f, 0x04, 0x9f, 0x58, 0xf9, 0xa6, 0x53, 0x76, 0xcb, 0x1e, 0x92, 0x67, 0xaa, 0x45, 0xce,
	0x93, 0x78, 0x7a, 0x94, 0xb5, 0x4e, 0xff, 0x11, 0xf8, 0x11, 0xc4, 0x0d, 0x9f, 0x1b, 0xeb, 0x83,
	0xb4, 0xb5, 0xde, 0x19, 0x4c, 0xff, 0x13, 0xe7, 0xa3, 0xb4, 0x41, 0x71, 0xc5, 0x13, 0x0b, 0x4a,
	0xef, 0x68, 0xef, 0x84, 0xf0, 0x2c, 0x0a, 0x55, 0xdc, 0xbd, 0xe4, 0xe7, 0x04, 0x0d, 0xa6, 0x67,
	0x84, 0xa6, 0xb1, 0xdd, 0xa0, 0xb8, 0xe3, 0x8b, 0xba, 0x69, 0x7d, 0xd0, 0xee, 0x68, 0x9c, 0x12,
	0x9d, 0x0f, 0x5a, 0xf4, 0xae, 0x1f, 0x79, 0x66, 0xa0, 0xa0, 0xd7, 0x47, 0x07, 0x1f, 0x7d, 0xf1,
	0xfb, 0x23, 0xd6, 0xb3, 0x4a, 0xf9, 0x6d, 0x8c, 0xb8, 0x65, 0x9f, 0x8c, 0xed, 0xa7, 0x14, 0xf7,
	0xe9, 0x3b, 0x00, 0x00, 0xff, 0xff, 0x92, 0x10, 0x3a, 0xfb, 0xc6, 0x01, 0x00, 0x00,
}