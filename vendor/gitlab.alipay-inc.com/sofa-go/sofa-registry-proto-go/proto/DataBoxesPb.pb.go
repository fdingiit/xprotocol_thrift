// Code generated by protoc-gen-go. DO NOT EDIT.
// source: DataBoxesPb.proto

package proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type DataBoxesPb struct {
	Data []*DataBoxPb `protobuf:"bytes,1,rep,name=data" json:"data,omitempty"`
}

func (m *DataBoxesPb) Reset()                    { *m = DataBoxesPb{} }
func (m *DataBoxesPb) String() string            { return proto.CompactTextString(m) }
func (*DataBoxesPb) ProtoMessage()               {}
func (*DataBoxesPb) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *DataBoxesPb) GetData() []*DataBoxPb {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*DataBoxesPb)(nil), "DataBoxesPb")
}

func init() { proto.RegisterFile("DataBoxesPb.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 132 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x74, 0x49, 0x2c, 0x49,
	0x74, 0xca, 0xaf, 0x48, 0x2d, 0x0e, 0x48, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x97, 0xe2, 0x87,
	0x0a, 0xc1, 0x04, 0x94, 0x74, 0xb9, 0xb8, 0x91, 0x54, 0x09, 0xc9, 0x71, 0xb1, 0xa4, 0x24, 0x96,
	0x24, 0x4a, 0x30, 0x2a, 0x30, 0x6b, 0x70, 0x1b, 0x71, 0xe9, 0xc1, 0x95, 0x07, 0x81, 0xc5, 0x9d,
	0x0c, 0xb8, 0x34, 0x92, 0xf3, 0x73, 0xf5, 0x12, 0x73, 0x32, 0x0b, 0x12, 0x2b, 0xf5, 0x8a, 0xf3,
	0xd3, 0x12, 0xf5, 0x8a, 0x52, 0xd3, 0x33, 0x8b, 0x4b, 0x8a, 0x2a, 0xf5, 0x8a, 0x53, 0x8b, 0xca,
	0x52, 0x8b, 0xf4, 0x72, 0xf3, 0x53, 0x52, 0x73, 0xf4, 0x0a, 0x92, 0x02, 0x18, 0xa3, 0x98, 0x0a,
	0x92, 0x92, 0xd8, 0xc0, 0xf6, 0x18, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x0c, 0x88, 0xe7, 0x7d,
	0x8d, 0x00, 0x00, 0x00,
}