// Code generated by protoc-gen-go. DO NOT EDIT.
// source: envoy/config/filter/network/server_forward/v2/server_forward.proto

package v2alpha

import (
	fmt "fmt"
	_ "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/wrappers"
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

type SofaForward_FallbackPolicy int32

const (
	SofaForward_NONE SofaForward_FallbackPolicy = 0
	SofaForward_NEXT SofaForward_FallbackPolicy = 1
)

var SofaForward_FallbackPolicy_name = map[int32]string{
	0: "NONE",
	1: "NEXT",
}

var SofaForward_FallbackPolicy_value = map[string]int32{
	"NONE": 0,
	"NEXT": 1,
}

func (x SofaForward_FallbackPolicy) String() string {
	return proto.EnumName(SofaForward_FallbackPolicy_name, int32(x))
}

func (SofaForward_FallbackPolicy) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_032736f017b2b8b3, []int{2, 0}
}

type Action_ActionType int32

const (
	Action_PASS    Action_ActionType = 0
	Action_FORWARD Action_ActionType = 1
)

var Action_ActionType_name = map[int32]string{
	0: "PASS",
	1: "FORWARD",
}

var Action_ActionType_value = map[string]int32{
	"PASS":    0,
	"FORWARD": 1,
}

func (x Action_ActionType) String() string {
	return proto.EnumName(Action_ActionType_name, int32(x))
}

func (Action_ActionType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_032736f017b2b8b3, []int{4, 0}
}

type ServerForward struct {
	Configs              []*Config `protobuf:"bytes,1,rep,name=configs,proto3" json:"configs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ServerForward) Reset()         { *m = ServerForward{} }
func (m *ServerForward) String() string { return proto.CompactTextString(m) }
func (*ServerForward) ProtoMessage()    {}
func (*ServerForward) Descriptor() ([]byte, []int) {
	return fileDescriptor_032736f017b2b8b3, []int{0}
}

func (m *ServerForward) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServerForward.Unmarshal(m, b)
}
func (m *ServerForward) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServerForward.Marshal(b, m, deterministic)
}
func (m *ServerForward) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServerForward.Merge(m, src)
}
func (m *ServerForward) XXX_Size() int {
	return xxx_messageInfo_ServerForward.Size(m)
}
func (m *ServerForward) XXX_DiscardUnknown() {
	xxx_messageInfo_ServerForward.DiscardUnknown(m)
}

var xxx_messageInfo_ServerForward proto.InternalMessageInfo

func (m *ServerForward) GetConfigs() []*Config {
	if m != nil {
		return m.Configs
	}
	return nil
}

type Config struct {
	Domains              []string       `protobuf:"bytes,1,rep,name=domains,proto3" json:"domains,omitempty"`
	Sofas                []*SofaForward `protobuf:"bytes,2,rep,name=sofas,proto3" json:"sofas,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_032736f017b2b8b3, []int{1}
}

func (m *Config) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config.Unmarshal(m, b)
}
func (m *Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config.Marshal(b, m, deterministic)
}
func (m *Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config.Merge(m, src)
}
func (m *Config) XXX_Size() int {
	return xxx_messageInfo_Config.Size(m)
}
func (m *Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Config proto.InternalMessageInfo

func (m *Config) GetDomains() []string {
	if m != nil {
		return m.Domains
	}
	return nil
}

func (m *Config) GetSofas() []*SofaForward {
	if m != nil {
		return m.Sofas
	}
	return nil
}

type SofaForward struct {
	Priority             int32                      `protobuf:"varint,1,opt,name=priority,proto3" json:"priority,omitempty"`
	RuleFallbackPolicy   SofaForward_FallbackPolicy `protobuf:"varint,2,opt,name=ruleFallbackPolicy,proto3,enum=envoy.config.filter.server_forward.v2alpha.SofaForward_FallbackPolicy" json:"ruleFallbackPolicy,omitempty"`
	ActionFallbackPolicy SofaForward_FallbackPolicy `protobuf:"varint,3,opt,name=actionFallbackPolicy,proto3,enum=envoy.config.filter.server_forward.v2alpha.SofaForward_FallbackPolicy" json:"actionFallbackPolicy,omitempty"`
	Match                []*MatchRequirement        `protobuf:"bytes,4,rep,name=match,proto3" json:"match,omitempty"`
	Actions              []*Action                  `protobuf:"bytes,5,rep,name=actions,proto3" json:"actions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *SofaForward) Reset()         { *m = SofaForward{} }
func (m *SofaForward) String() string { return proto.CompactTextString(m) }
func (*SofaForward) ProtoMessage()    {}
func (*SofaForward) Descriptor() ([]byte, []int) {
	return fileDescriptor_032736f017b2b8b3, []int{2}
}

func (m *SofaForward) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SofaForward.Unmarshal(m, b)
}
func (m *SofaForward) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SofaForward.Marshal(b, m, deterministic)
}
func (m *SofaForward) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SofaForward.Merge(m, src)
}
func (m *SofaForward) XXX_Size() int {
	return xxx_messageInfo_SofaForward.Size(m)
}
func (m *SofaForward) XXX_DiscardUnknown() {
	xxx_messageInfo_SofaForward.DiscardUnknown(m)
}

var xxx_messageInfo_SofaForward proto.InternalMessageInfo

func (m *SofaForward) GetPriority() int32 {
	if m != nil {
		return m.Priority
	}
	return 0
}

func (m *SofaForward) GetRuleFallbackPolicy() SofaForward_FallbackPolicy {
	if m != nil {
		return m.RuleFallbackPolicy
	}
	return SofaForward_NONE
}

func (m *SofaForward) GetActionFallbackPolicy() SofaForward_FallbackPolicy {
	if m != nil {
		return m.ActionFallbackPolicy
	}
	return SofaForward_NONE
}

func (m *SofaForward) GetMatch() []*MatchRequirement {
	if m != nil {
		return m.Match
	}
	return nil
}

func (m *SofaForward) GetActions() []*Action {
	if m != nil {
		return m.Actions
	}
	return nil
}

type MatchRequirement struct {
	Ref                  string   `protobuf:"bytes,1,opt,name=ref,proto3" json:"ref,omitempty"`
	Key                  string   `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Operator             string   `protobuf:"bytes,3,opt,name=operator,proto3" json:"operator,omitempty"`
	Values               []string `protobuf:"bytes,4,rep,name=values,proto3" json:"values,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MatchRequirement) Reset()         { *m = MatchRequirement{} }
func (m *MatchRequirement) String() string { return proto.CompactTextString(m) }
func (*MatchRequirement) ProtoMessage()    {}
func (*MatchRequirement) Descriptor() ([]byte, []int) {
	return fileDescriptor_032736f017b2b8b3, []int{3}
}

func (m *MatchRequirement) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MatchRequirement.Unmarshal(m, b)
}
func (m *MatchRequirement) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MatchRequirement.Marshal(b, m, deterministic)
}
func (m *MatchRequirement) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MatchRequirement.Merge(m, src)
}
func (m *MatchRequirement) XXX_Size() int {
	return xxx_messageInfo_MatchRequirement.Size(m)
}
func (m *MatchRequirement) XXX_DiscardUnknown() {
	xxx_messageInfo_MatchRequirement.DiscardUnknown(m)
}

var xxx_messageInfo_MatchRequirement proto.InternalMessageInfo

func (m *MatchRequirement) GetRef() string {
	if m != nil {
		return m.Ref
	}
	return ""
}

func (m *MatchRequirement) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *MatchRequirement) GetOperator() string {
	if m != nil {
		return m.Operator
	}
	return ""
}

func (m *MatchRequirement) GetValues() []string {
	if m != nil {
		return m.Values
	}
	return nil
}

type Action struct {
	Weight               int32             `protobuf:"varint,1,opt,name=weight,proto3" json:"weight,omitempty"`
	Type                 Action_ActionType `protobuf:"varint,2,opt,name=type,proto3,enum=envoy.config.filter.server_forward.v2alpha.Action_ActionType" json:"type,omitempty"`
	Hosts                []string          `protobuf:"bytes,3,rep,name=hosts,proto3" json:"hosts,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Action) Reset()         { *m = Action{} }
func (m *Action) String() string { return proto.CompactTextString(m) }
func (*Action) ProtoMessage()    {}
func (*Action) Descriptor() ([]byte, []int) {
	return fileDescriptor_032736f017b2b8b3, []int{4}
}

func (m *Action) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Action.Unmarshal(m, b)
}
func (m *Action) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Action.Marshal(b, m, deterministic)
}
func (m *Action) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Action.Merge(m, src)
}
func (m *Action) XXX_Size() int {
	return xxx_messageInfo_Action.Size(m)
}
func (m *Action) XXX_DiscardUnknown() {
	xxx_messageInfo_Action.DiscardUnknown(m)
}

var xxx_messageInfo_Action proto.InternalMessageInfo

func (m *Action) GetWeight() int32 {
	if m != nil {
		return m.Weight
	}
	return 0
}

func (m *Action) GetType() Action_ActionType {
	if m != nil {
		return m.Type
	}
	return Action_PASS
}

func (m *Action) GetHosts() []string {
	if m != nil {
		return m.Hosts
	}
	return nil
}

func init() {
	proto.RegisterEnum("envoy.config.filter.server_forward.v2alpha.SofaForward_FallbackPolicy", SofaForward_FallbackPolicy_name, SofaForward_FallbackPolicy_value)
	proto.RegisterEnum("envoy.config.filter.server_forward.v2alpha.Action_ActionType", Action_ActionType_name, Action_ActionType_value)
	proto.RegisterType((*ServerForward)(nil), "envoy.config.filter.server_forward.v2alpha.ServerForward")
	proto.RegisterType((*Config)(nil), "envoy.config.filter.server_forward.v2alpha.Config")
	proto.RegisterType((*SofaForward)(nil), "envoy.config.filter.server_forward.v2alpha.SofaForward")
	proto.RegisterType((*MatchRequirement)(nil), "envoy.config.filter.server_forward.v2alpha.MatchRequirement")
	proto.RegisterType((*Action)(nil), "envoy.config.filter.server_forward.v2alpha.Action")
}

func init() {
	proto.RegisterFile("envoy/config/filter/network/server_forward/v2/server_forward.proto", fileDescriptor_032736f017b2b8b3)
}

var fileDescriptor_032736f017b2b8b3 = []byte{
	// 533 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0x4d, 0x8f, 0xd3, 0x30,
	0x10, 0xdd, 0x6c, 0xbf, 0xb6, 0x53, 0xb1, 0x8a, 0xac, 0x15, 0x44, 0x15, 0x42, 0x55, 0xe0, 0x50,
	0x71, 0x48, 0x50, 0x39, 0x70, 0x81, 0x43, 0x0b, 0xdb, 0xd3, 0x7e, 0x14, 0x77, 0x25, 0x10, 0x12,
	0x42, 0x6e, 0xea, 0xb4, 0xa6, 0x69, 0x26, 0xeb, 0xb8, 0x29, 0xe1, 0x47, 0xf1, 0xaf, 0xf8, 0x1f,
	0x28, 0x76, 0xba, 0xd0, 0x68, 0x0f, 0x14, 0x71, 0xca, 0xbc, 0x99, 0xf1, 0x7b, 0x2f, 0xe3, 0x91,
	0x61, 0xc4, 0xe3, 0x0c, 0x73, 0x3f, 0xc0, 0x38, 0x14, 0x0b, 0x3f, 0x14, 0x91, 0xe2, 0xd2, 0x8f,
	0xb9, 0xda, 0xa2, 0x5c, 0xf9, 0x29, 0x97, 0x19, 0x97, 0x5f, 0x42, 0x94, 0x5b, 0x26, 0xe7, 0x7e,
	0x36, 0xa8, 0x64, 0xbc, 0x44, 0xa2, 0x42, 0xf2, 0x5c, 0x73, 0x78, 0x86, 0xc3, 0x33, 0x1c, 0x5e,
	0xa5, 0x33, 0x1b, 0xb0, 0x28, 0x59, 0xb2, 0xee, 0xa3, 0x8c, 0x45, 0x62, 0xce, 0x14, 0xf7, 0x77,
	0x81, 0x21, 0xe9, 0x3e, 0x36, 0x46, 0x58, 0x22, 0x0a, 0x9d, 0x00, 0x25, 0xf7, 0x67, 0x2c, 0xdd,
	0x55, 0x9f, 0x2c, 0x10, 0x17, 0x11, 0xf7, 0x35, 0x9a, 0x6d, 0x42, 0x7f, 0x2b, 0x59, 0x92, 0x70,
	0x99, 0x9a, 0xba, 0xfb, 0x19, 0x1e, 0x4c, 0xb5, 0xe0, 0xd8, 0xe8, 0x91, 0x0b, 0x68, 0x19, 0x3f,
	0xa9, 0x63, 0xf5, 0x6a, 0xfd, 0xce, 0x60, 0xe0, 0xfd, 0xbd, 0x4b, 0xef, 0xad, 0x6e, 0xa2, 0x3b,
	0x0a, 0xf7, 0x16, 0x9a, 0x26, 0x45, 0x1c, 0x68, 0xcd, 0x71, 0xcd, 0x44, 0x6c, 0x78, 0xdb, 0x74,
	0x07, 0xc9, 0x25, 0x34, 0x52, 0x0c, 0x59, 0xea, 0x1c, 0x6b, 0xbd, 0x57, 0x87, 0xe8, 0x4d, 0x31,
	0x64, 0xa5, 0x73, 0x6a, 0x58, 0xdc, 0x9f, 0x35, 0xe8, 0xfc, 0x91, 0x26, 0x5d, 0x38, 0x49, 0xa4,
	0x40, 0x29, 0x54, 0xee, 0x58, 0x3d, 0xab, 0xdf, 0xa0, 0x77, 0x98, 0x64, 0x40, 0xe4, 0x26, 0xe2,
	0x63, 0x16, 0x45, 0x33, 0x16, 0xac, 0x26, 0x18, 0x89, 0x20, 0x77, 0x8e, 0x7b, 0x56, 0xff, 0x74,
	0x30, 0xfe, 0x47, 0x1f, 0xde, 0x3e, 0x1b, 0xbd, 0x47, 0x81, 0x7c, 0x87, 0x33, 0x16, 0x28, 0x81,
	0x71, 0x45, 0xb9, 0xf6, 0x5f, 0x95, 0xef, 0xd5, 0x20, 0x14, 0x1a, 0x6b, 0xa6, 0x82, 0xa5, 0x53,
	0xd7, 0xe3, 0x7e, 0x7d, 0x88, 0xd8, 0x65, 0x71, 0x90, 0xf2, 0xdb, 0x8d, 0x90, 0x7c, 0xcd, 0x63,
	0x45, 0x0d, 0x55, 0xb1, 0x34, 0x46, 0x2b, 0x75, 0x1a, 0x87, 0x2f, 0xcd, 0x50, 0x1f, 0xa5, 0x3b,
	0x0a, 0xf7, 0x19, 0x9c, 0x56, 0x3c, 0x9f, 0x40, 0xfd, 0xea, 0xfa, 0xea, 0xdc, 0x3e, 0xd2, 0xd1,
	0xf9, 0xc7, 0x1b, 0xdb, 0x72, 0xbf, 0x82, 0x5d, 0xb5, 0x43, 0x6c, 0xa8, 0x49, 0x1e, 0xea, 0x6b,
	0x6e, 0xd3, 0x22, 0x2c, 0x32, 0x2b, 0x6e, 0xae, 0xb4, 0x4d, 0x8b, 0xb0, 0xd8, 0x07, 0x4c, 0xb8,
	0x64, 0x0a, 0xa5, 0x9e, 0x77, 0x9b, 0xde, 0x61, 0xf2, 0x10, 0x9a, 0x19, 0x8b, 0x36, 0x3c, 0xd5,
	0xc3, 0x69, 0xd3, 0x12, 0xb9, 0x3f, 0x2c, 0x68, 0x1a, 0x97, 0x45, 0xcb, 0x96, 0x8b, 0xc5, 0x52,
	0x95, 0xcb, 0x54, 0x22, 0xf2, 0x1e, 0xea, 0x2a, 0x4f, 0x78, 0xb9, 0x3c, 0x6f, 0x0e, 0xff, 0xff,
	0xf2, 0x73, 0x93, 0x27, 0x9c, 0x6a, 0x2a, 0x72, 0x06, 0x8d, 0x25, 0xa6, 0x2a, 0x75, 0x6a, 0xda,
	0x8c, 0x01, 0xee, 0x53, 0x80, 0xdf, 0x9d, 0xc5, 0x3c, 0x26, 0xc3, 0xe9, 0xd4, 0x3e, 0x22, 0x1d,
	0x68, 0x8d, 0xaf, 0xe9, 0x87, 0x21, 0x7d, 0x67, 0x5b, 0xa3, 0x0b, 0x78, 0x21, 0xd0, 0x78, 0x48,
	0x24, 0x7e, 0xcb, 0xf7, 0xed, 0x18, 0x1f, 0x15, 0x1b, 0xa3, 0xfd, 0x87, 0x60, 0x62, 0x7d, 0x6a,
	0x95, 0x95, 0x59, 0x53, 0xbf, 0x15, 0x2f, 0x7f, 0x05, 0x00, 0x00, 0xff, 0xff, 0xf6, 0xeb, 0x47,
	0x0a, 0xf4, 0x04, 0x00, 0x00,
}
