// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ScopeEnumPb.proto

package proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type ScopeEnumPb int32

const (
	ScopeEnumPb_zone       ScopeEnumPb = 0
	ScopeEnumPb_dataCenter ScopeEnumPb = 1
	ScopeEnumPb_global     ScopeEnumPb = 2
)

var ScopeEnumPb_name = map[int32]string{
	0: "zone",
	1: "dataCenter",
	2: "global",
}

var ScopeEnumPb_value = map[string]int32{
	"zone":       0,
	"dataCenter": 1,
	"global":     2,
}

func (x ScopeEnumPb) String() string {
	return proto.EnumName(ScopeEnumPb_name, int32(x))
}
func (ScopeEnumPb) EnumDescriptor() ([]byte, []int) { return fileDescriptor8, []int{0} }

func init() {
	proto.RegisterEnum("ScopeEnumPb", ScopeEnumPb_name, ScopeEnumPb_value)
}

func init() { proto.RegisterFile("ScopeEnumPb.proto", fileDescriptor8) }

var fileDescriptor8 = []byte{
	// 141 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0xcd, 0xb1, 0xca, 0xc2, 0x30,
	0x10, 0x00, 0xe0, 0xbf, 0xe5, 0xa7, 0xc8, 0x09, 0x12, 0xf3, 0x14, 0xe2, 0x70, 0x08, 0x7d, 0x03,
	0xc5, 0xbd, 0xe0, 0xe6, 0x76, 0x69, 0xcf, 0x52, 0x48, 0x72, 0xc7, 0x25, 0x0a, 0xf5, 0xe9, 0x05,
	0x27, 0xd7, 0x6f, 0xf9, 0x60, 0x7f, 0x1b, 0x45, 0xf9, 0x9a, 0x9f, 0x69, 0x08, 0xa8, 0x26, 0x55,
	0x8e, 0x3d, 0x6c, 0x7f, 0xd0, 0x6f, 0xe0, 0xff, 0x2d, 0x99, 0xdd, 0x9f, 0xdf, 0x01, 0x4c, 0x54,
	0xe9, 0xc2, 0xb9, 0xb2, 0xb9, 0xc6, 0x03, 0x74, 0x73, 0x94, 0x40, 0xd1, 0xb5, 0xe7, 0x13, 0x1c,
	0x46, 0x49, 0x48, 0x71, 0x51, 0x5a, 0xb1, 0xc8, 0x83, 0xd0, 0x78, 0x5e, 0x4a, 0xb5, 0x15, 0x0b,
	0xdb, 0x8b, 0x0d, 0x93, 0x4c, 0x1c, 0x51, 0xc3, 0xd0, 0xdc, 0x5b, 0x0d, 0xa1, 0xfb, 0x6e, 0xfd,
	0x27, 0x00, 0x00, 0xff, 0xff, 0x92, 0xa9, 0x6c, 0x9c, 0x82, 0x00, 0x00, 0x00,
}