// Code generated by protoc-gen-go. DO NOT EDIT.
// source: envoy/config/filter/network/headers_config/v2/headers_config.proto

package v2alpha

import (
	fmt "fmt"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
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

type HeadersConfig struct {
	Configs              []*Config `protobuf:"bytes,1,rep,name=configs,proto3" json:"configs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *HeadersConfig) Reset()         { *m = HeadersConfig{} }
func (m *HeadersConfig) String() string { return proto.CompactTextString(m) }
func (*HeadersConfig) ProtoMessage()    {}
func (*HeadersConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a0fbaebebee81f7, []int{0}
}

func (m *HeadersConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeadersConfig.Unmarshal(m, b)
}
func (m *HeadersConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeadersConfig.Marshal(b, m, deterministic)
}
func (m *HeadersConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeadersConfig.Merge(m, src)
}
func (m *HeadersConfig) XXX_Size() int {
	return xxx_messageInfo_HeadersConfig.Size(m)
}
func (m *HeadersConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_HeadersConfig.DiscardUnknown(m)
}

var xxx_messageInfo_HeadersConfig proto.InternalMessageInfo

func (m *HeadersConfig) GetConfigs() []*Config {
	if m != nil {
		return m.Configs
	}
	return nil
}

type Config struct {
	Domains              []string           `protobuf:"bytes,1,rep,name=domains,proto3" json:"domains,omitempty"`
	Sofas                []*SofaHeadersRule `protobuf:"bytes,2,rep,name=sofas,proto3" json:"sofas,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a0fbaebebee81f7, []int{1}
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

func (m *Config) GetSofas() []*SofaHeadersRule {
	if m != nil {
		return m.Sofas
	}
	return nil
}

type SofaHeadersRule struct {
	Match                []*MatchRequirement `protobuf:"bytes,1,rep,name=match,proto3" json:"match,omitempty"`
	Actions              []*Action           `protobuf:"bytes,2,rep,name=actions,proto3" json:"actions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *SofaHeadersRule) Reset()         { *m = SofaHeadersRule{} }
func (m *SofaHeadersRule) String() string { return proto.CompactTextString(m) }
func (*SofaHeadersRule) ProtoMessage()    {}
func (*SofaHeadersRule) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a0fbaebebee81f7, []int{2}
}

func (m *SofaHeadersRule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SofaHeadersRule.Unmarshal(m, b)
}
func (m *SofaHeadersRule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SofaHeadersRule.Marshal(b, m, deterministic)
}
func (m *SofaHeadersRule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SofaHeadersRule.Merge(m, src)
}
func (m *SofaHeadersRule) XXX_Size() int {
	return xxx_messageInfo_SofaHeadersRule.Size(m)
}
func (m *SofaHeadersRule) XXX_DiscardUnknown() {
	xxx_messageInfo_SofaHeadersRule.DiscardUnknown(m)
}

var xxx_messageInfo_SofaHeadersRule proto.InternalMessageInfo

func (m *SofaHeadersRule) GetMatch() []*MatchRequirement {
	if m != nil {
		return m.Match
	}
	return nil
}

func (m *SofaHeadersRule) GetActions() []*Action {
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
	return fileDescriptor_8a0fbaebebee81f7, []int{3}
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
	Weight                  *wrappers.UInt32Value     `protobuf:"bytes,1,opt,name=weight,proto3" json:"weight,omitempty"`
	RequestHeadersToAdd     []*core.HeaderValueOption `protobuf:"bytes,2,rep,name=request_headers_to_add,json=requestHeadersToAdd,proto3" json:"request_headers_to_add,omitempty"`
	RequestHeadersToRemove  []string                  `protobuf:"bytes,3,rep,name=request_headers_to_remove,json=requestHeadersToRemove,proto3" json:"request_headers_to_remove,omitempty"`
	ResponseHeadersToAdd    []*core.HeaderValueOption `protobuf:"bytes,4,rep,name=response_headers_to_add,json=responseHeadersToAdd,proto3" json:"response_headers_to_add,omitempty"`
	ResponseHeadersToRemove []string                  `protobuf:"bytes,5,rep,name=response_headers_to_remove,json=responseHeadersToRemove,proto3" json:"response_headers_to_remove,omitempty"`
	XXX_NoUnkeyedLiteral    struct{}                  `json:"-"`
	XXX_unrecognized        []byte                    `json:"-"`
	XXX_sizecache           int32                     `json:"-"`
}

func (m *Action) Reset()         { *m = Action{} }
func (m *Action) String() string { return proto.CompactTextString(m) }
func (*Action) ProtoMessage()    {}
func (*Action) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a0fbaebebee81f7, []int{4}
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

func (m *Action) GetWeight() *wrappers.UInt32Value {
	if m != nil {
		return m.Weight
	}
	return nil
}

func (m *Action) GetRequestHeadersToAdd() []*core.HeaderValueOption {
	if m != nil {
		return m.RequestHeadersToAdd
	}
	return nil
}

func (m *Action) GetRequestHeadersToRemove() []string {
	if m != nil {
		return m.RequestHeadersToRemove
	}
	return nil
}

func (m *Action) GetResponseHeadersToAdd() []*core.HeaderValueOption {
	if m != nil {
		return m.ResponseHeadersToAdd
	}
	return nil
}

func (m *Action) GetResponseHeadersToRemove() []string {
	if m != nil {
		return m.ResponseHeadersToRemove
	}
	return nil
}

func init() {
	proto.RegisterType((*HeadersConfig)(nil), "envoy.config.filter.headers_config.v2alpha.HeadersConfig")
	proto.RegisterType((*Config)(nil), "envoy.config.filter.headers_config.v2alpha.Config")
	proto.RegisterType((*SofaHeadersRule)(nil), "envoy.config.filter.headers_config.v2alpha.SofaHeadersRule")
	proto.RegisterType((*MatchRequirement)(nil), "envoy.config.filter.headers_config.v2alpha.MatchRequirement")
	proto.RegisterType((*Action)(nil), "envoy.config.filter.headers_config.v2alpha.Action")
}

func init() {
	proto.RegisterFile("envoy/config/filter/network/headers_config/v2/headers_config.proto", fileDescriptor_8a0fbaebebee81f7)
}

var fileDescriptor_8a0fbaebebee81f7 = []byte{
	// 511 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0xcf, 0x8e, 0xd3, 0x30,
	0x10, 0xc6, 0xd5, 0xed, 0xb6, 0x65, 0xbd, 0x42, 0xac, 0x0c, 0x6a, 0x43, 0xb5, 0x42, 0xab, 0x88,
	0xc3, 0x8a, 0x83, 0x83, 0xb2, 0x5c, 0xd0, 0x72, 0xd9, 0x72, 0x01, 0x69, 0x11, 0x60, 0xfe, 0x48,
	0x80, 0x50, 0xe5, 0x36, 0x93, 0x36, 0x34, 0xcd, 0x78, 0x6d, 0x27, 0xa5, 0xcf, 0xc5, 0x7b, 0xf1,
	0x0c, 0x28, 0xb6, 0x83, 0xd4, 0xb0, 0x07, 0x7a, 0x8b, 0x67, 0x3c, 0xbf, 0xef, 0x9b, 0xf1, 0x84,
	0x4c, 0xa0, 0xa8, 0x70, 0x1b, 0xcd, 0xb1, 0x48, 0xb3, 0x45, 0x94, 0x66, 0xb9, 0x01, 0x15, 0x15,
	0x60, 0x36, 0xa8, 0x56, 0xd1, 0x12, 0x44, 0x02, 0x4a, 0x4f, 0x7d, 0xb6, 0x8a, 0x5b, 0x11, 0x26,
	0x15, 0x1a, 0xa4, 0x4f, 0x2c, 0x83, 0xf9, 0x98, 0x63, 0xb0, 0xd6, 0xcd, 0x2a, 0x16, 0xb9, 0x5c,
	0x8a, 0xf1, 0xa8, 0x12, 0x79, 0x96, 0x08, 0x03, 0x51, 0xf3, 0xe1, 0x20, 0xe3, 0x53, 0x67, 0x44,
	0xc8, 0xac, 0xd6, 0x99, 0xa3, 0x82, 0x68, 0x26, 0x74, 0x93, 0x7d, 0xb4, 0x40, 0x5c, 0xe4, 0x10,
	0xd9, 0xd3, 0xac, 0x4c, 0xa3, 0x8d, 0x12, 0x52, 0x82, 0xd2, 0x2e, 0x1f, 0x7e, 0x27, 0x77, 0x5f,
	0x39, 0xc1, 0x97, 0x56, 0x8f, 0x5e, 0x93, 0x81, 0x53, 0xd6, 0x41, 0xe7, 0xac, 0x7b, 0x7e, 0x1c,
	0xc7, 0xec, 0xff, 0x5d, 0x32, 0x07, 0xe1, 0x0d, 0x22, 0x2c, 0x49, 0xdf, 0x73, 0x03, 0x32, 0x48,
	0x70, 0x2d, 0xb2, 0xc2, 0x71, 0x8f, 0x78, 0x73, 0xa4, 0xef, 0x49, 0x4f, 0x63, 0x2a, 0x74, 0x70,
	0x60, 0xf5, 0x2e, 0xf7, 0xd1, 0xfb, 0x80, 0xa9, 0xf0, 0xfe, 0x79, 0x99, 0x03, 0x77, 0xa4, 0xf0,
	0x57, 0x87, 0xdc, 0x6b, 0xa5, 0x28, 0x27, 0xbd, 0xb5, 0x30, 0xf3, 0xa5, 0x6f, 0xeb, 0xc5, 0x3e,
	0x32, 0x6f, 0xea, 0x42, 0x0e, 0x37, 0x65, 0xa6, 0x60, 0x0d, 0x85, 0xe1, 0x0e, 0x55, 0x0f, 0x4b,
	0xcc, 0x4d, 0x86, 0x45, 0x63, 0x7e, 0xaf, 0x61, 0x5d, 0xd9, 0x52, 0xde, 0x20, 0xc2, 0x1f, 0xe4,
	0xa4, 0x2d, 0x44, 0x4f, 0x48, 0x57, 0x41, 0x1a, 0x74, 0xce, 0x3a, 0xe7, 0x47, 0xbc, 0xfe, 0xac,
	0x23, 0x2b, 0xd8, 0x06, 0x07, 0x2e, 0xb2, 0x82, 0x2d, 0x1d, 0x93, 0x3b, 0x28, 0x41, 0x09, 0x83,
	0x2a, 0xe8, 0xda, 0xf0, 0xdf, 0x33, 0x1d, 0x92, 0x7e, 0x25, 0xf2, 0x12, 0x74, 0x70, 0x68, 0xa7,
	0xee, 0x4f, 0xe1, 0xef, 0x03, 0xd2, 0x77, 0xfa, 0xf4, 0x19, 0xe9, 0x6f, 0x20, 0x5b, 0x2c, 0x8d,
	0x55, 0x39, 0x8e, 0x4f, 0x99, 0xdb, 0x19, 0xd6, 0xec, 0x0c, 0xfb, 0xf4, 0xba, 0x30, 0x17, 0xf1,
	0xe7, 0xba, 0x92, 0xfb, 0xbb, 0xf4, 0x0b, 0x19, 0x2a, 0xb8, 0x29, 0x41, 0x9b, 0x69, 0xd3, 0x9e,
	0xc1, 0xa9, 0x48, 0x12, 0x3f, 0x89, 0xc7, 0x7e, 0x12, 0x42, 0x66, 0xac, 0x8a, 0x59, 0xbd, 0x97,
	0xcc, 0x3d, 0x87, 0xe5, 0xbc, 0x95, 0xb6, 0xf7, 0xfb, 0x9e, 0xe1, 0x1f, 0xea, 0x23, 0x5e, 0x25,
	0x09, 0x7d, 0x4e, 0x1e, 0xde, 0x82, 0x56, 0xb0, 0xc6, 0x0a, 0x82, 0xae, 0x6d, 0x63, 0xd8, 0xae,
	0xe3, 0x36, 0x4b, 0xbf, 0x91, 0x91, 0x02, 0x2d, 0xb1, 0xd0, 0xd0, 0xb6, 0x75, 0xb8, 0x87, 0xad,
	0x07, 0x0d, 0x64, 0xc7, 0xd7, 0x25, 0x19, 0xdf, 0x06, 0xf7, 0xc6, 0x7a, 0xd6, 0xd8, 0xe8, 0x9f,
	0x4a, 0xe7, 0x6c, 0x72, 0x4d, 0x9e, 0x66, 0xe8, 0xc4, 0xa5, 0xc2, 0x9f, 0xdb, 0xdd, 0x45, 0xf1,
	0xb8, 0xdd, 0x05, 0x99, 0xec, 0xfe, 0x9a, 0xef, 0x3a, 0x5f, 0x07, 0x3e, 0x33, 0xeb, 0xdb, 0xb7,
	0xb9, 0xf8, 0x13, 0x00, 0x00, 0xff, 0xff, 0xe6, 0x14, 0x5d, 0xa0, 0x86, 0x04, 0x00, 0x00,
}
