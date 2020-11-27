// Code generated by protoc-gen-go. DO NOT EDIT.
// source: envoy/api/v2/discovery.proto

package envoy_api_v2

import (
	fmt "fmt"
	_ "github.com/cncf/udpa/go/udpa/annotations"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	status "google.golang.org/genproto/googleapis/rpc/status"
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

type DiscoveryRequest struct {
	VersionInfo          string         `protobuf:"bytes,1,opt,name=version_info,json=versionInfo,proto3" json:"version_info,omitempty"`
	Node                 *core.Node     `protobuf:"bytes,2,opt,name=node,proto3" json:"node,omitempty"`
	ResourceNames        []string       `protobuf:"bytes,3,rep,name=resource_names,json=resourceNames,proto3" json:"resource_names,omitempty"`
	TypeUrl              string         `protobuf:"bytes,4,opt,name=type_url,json=typeUrl,proto3" json:"type_url,omitempty"`
	ResponseNonce        string         `protobuf:"bytes,5,opt,name=response_nonce,json=responseNonce,proto3" json:"response_nonce,omitempty"`
	ErrorDetail          *status.Status `protobuf:"bytes,6,opt,name=error_detail,json=errorDetail,proto3" json:"error_detail,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *DiscoveryRequest) Reset()         { *m = DiscoveryRequest{} }
func (m *DiscoveryRequest) String() string { return proto.CompactTextString(m) }
func (*DiscoveryRequest) ProtoMessage()    {}
func (*DiscoveryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2c7365e287e5c035, []int{0}
}

func (m *DiscoveryRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DiscoveryRequest.Unmarshal(m, b)
}
func (m *DiscoveryRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DiscoveryRequest.Marshal(b, m, deterministic)
}
func (m *DiscoveryRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DiscoveryRequest.Merge(m, src)
}
func (m *DiscoveryRequest) XXX_Size() int {
	return xxx_messageInfo_DiscoveryRequest.Size(m)
}
func (m *DiscoveryRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DiscoveryRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DiscoveryRequest proto.InternalMessageInfo

func (m *DiscoveryRequest) GetVersionInfo() string {
	if m != nil {
		return m.VersionInfo
	}
	return ""
}

func (m *DiscoveryRequest) GetNode() *core.Node {
	if m != nil {
		return m.Node
	}
	return nil
}

func (m *DiscoveryRequest) GetResourceNames() []string {
	if m != nil {
		return m.ResourceNames
	}
	return nil
}

func (m *DiscoveryRequest) GetTypeUrl() string {
	if m != nil {
		return m.TypeUrl
	}
	return ""
}

func (m *DiscoveryRequest) GetResponseNonce() string {
	if m != nil {
		return m.ResponseNonce
	}
	return ""
}

func (m *DiscoveryRequest) GetErrorDetail() *status.Status {
	if m != nil {
		return m.ErrorDetail
	}
	return nil
}

type DiscoveryResponse struct {
	VersionInfo          string                   `protobuf:"bytes,1,opt,name=version_info,json=versionInfo,proto3" json:"version_info,omitempty"`
	Resources            []*any.Any               `protobuf:"bytes,2,rep,name=resources,proto3" json:"resources,omitempty"`
	Canary               bool                     `protobuf:"varint,3,opt,name=canary,proto3" json:"canary,omitempty"`
	TypeUrl              string                   `protobuf:"bytes,4,opt,name=type_url,json=typeUrl,proto3" json:"type_url,omitempty"`
	Nonce                string                   `protobuf:"bytes,5,opt,name=nonce,proto3" json:"nonce,omitempty"`
	ControlPlane         *core.ControlPlane       `protobuf:"bytes,6,opt,name=control_plane,json=controlPlane,proto3" json:"control_plane,omitempty"`
	ErrorDetail          *status.Status           `protobuf:"bytes,80,opt,name=error_detail,json=errorDetail,proto3" json:"error_detail,omitempty"`
	ExpandField          *DiscoveryResponseExpand `protobuf:"bytes,90,opt,name=expand_field,json=expandField,proto3" json:"expand_field,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *DiscoveryResponse) Reset()         { *m = DiscoveryResponse{} }
func (m *DiscoveryResponse) String() string { return proto.CompactTextString(m) }
func (*DiscoveryResponse) ProtoMessage()    {}
func (*DiscoveryResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2c7365e287e5c035, []int{1}
}

func (m *DiscoveryResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DiscoveryResponse.Unmarshal(m, b)
}
func (m *DiscoveryResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DiscoveryResponse.Marshal(b, m, deterministic)
}
func (m *DiscoveryResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DiscoveryResponse.Merge(m, src)
}
func (m *DiscoveryResponse) XXX_Size() int {
	return xxx_messageInfo_DiscoveryResponse.Size(m)
}
func (m *DiscoveryResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DiscoveryResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DiscoveryResponse proto.InternalMessageInfo

func (m *DiscoveryResponse) GetVersionInfo() string {
	if m != nil {
		return m.VersionInfo
	}
	return ""
}

func (m *DiscoveryResponse) GetResources() []*any.Any {
	if m != nil {
		return m.Resources
	}
	return nil
}

func (m *DiscoveryResponse) GetCanary() bool {
	if m != nil {
		return m.Canary
	}
	return false
}

func (m *DiscoveryResponse) GetTypeUrl() string {
	if m != nil {
		return m.TypeUrl
	}
	return ""
}

func (m *DiscoveryResponse) GetNonce() string {
	if m != nil {
		return m.Nonce
	}
	return ""
}

func (m *DiscoveryResponse) GetControlPlane() *core.ControlPlane {
	if m != nil {
		return m.ControlPlane
	}
	return nil
}

func (m *DiscoveryResponse) GetErrorDetail() *status.Status {
	if m != nil {
		return m.ErrorDetail
	}
	return nil
}

func (m *DiscoveryResponse) GetExpandField() *DiscoveryResponseExpand {
	if m != nil {
		return m.ExpandField
	}
	return nil
}

type DiscoveryResponseExpand struct {
	CheckMessage         *CheckMessage `protobuf:"bytes,1,opt,name=check_message,json=checkMessage,proto3" json:"check_message,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *DiscoveryResponseExpand) Reset()         { *m = DiscoveryResponseExpand{} }
func (m *DiscoveryResponseExpand) String() string { return proto.CompactTextString(m) }
func (*DiscoveryResponseExpand) ProtoMessage()    {}
func (*DiscoveryResponseExpand) Descriptor() ([]byte, []int) {
	return fileDescriptor_2c7365e287e5c035, []int{2}
}

func (m *DiscoveryResponseExpand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DiscoveryResponseExpand.Unmarshal(m, b)
}
func (m *DiscoveryResponseExpand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DiscoveryResponseExpand.Marshal(b, m, deterministic)
}
func (m *DiscoveryResponseExpand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DiscoveryResponseExpand.Merge(m, src)
}
func (m *DiscoveryResponseExpand) XXX_Size() int {
	return xxx_messageInfo_DiscoveryResponseExpand.Size(m)
}
func (m *DiscoveryResponseExpand) XXX_DiscardUnknown() {
	xxx_messageInfo_DiscoveryResponseExpand.DiscardUnknown(m)
}

var xxx_messageInfo_DiscoveryResponseExpand proto.InternalMessageInfo

func (m *DiscoveryResponseExpand) GetCheckMessage() *CheckMessage {
	if m != nil {
		return m.CheckMessage
	}
	return nil
}

type CheckMessage struct {
	CrSummarys           []*CRSummary `protobuf:"bytes,1,rep,name=cr_summarys,json=crSummarys,proto3" json:"cr_summarys,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *CheckMessage) Reset()         { *m = CheckMessage{} }
func (m *CheckMessage) String() string { return proto.CompactTextString(m) }
func (*CheckMessage) ProtoMessage()    {}
func (*CheckMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_2c7365e287e5c035, []int{3}
}

func (m *CheckMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckMessage.Unmarshal(m, b)
}
func (m *CheckMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckMessage.Marshal(b, m, deterministic)
}
func (m *CheckMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckMessage.Merge(m, src)
}
func (m *CheckMessage) XXX_Size() int {
	return xxx_messageInfo_CheckMessage.Size(m)
}
func (m *CheckMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckMessage.DiscardUnknown(m)
}

var xxx_messageInfo_CheckMessage proto.InternalMessageInfo

func (m *CheckMessage) GetCrSummarys() []*CRSummary {
	if m != nil {
		return m.CrSummarys
	}
	return nil
}

type CRSummary struct {
	CrName               string   `protobuf:"bytes,1,opt,name=cr_name,json=crName,proto3" json:"cr_name,omitempty"`
	CrKind               string   `protobuf:"bytes,2,opt,name=cr_kind,json=crKind,proto3" json:"cr_kind,omitempty"`
	CrVersion            string   `protobuf:"bytes,3,opt,name=cr_version,json=crVersion,proto3" json:"cr_version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CRSummary) Reset()         { *m = CRSummary{} }
func (m *CRSummary) String() string { return proto.CompactTextString(m) }
func (*CRSummary) ProtoMessage()    {}
func (*CRSummary) Descriptor() ([]byte, []int) {
	return fileDescriptor_2c7365e287e5c035, []int{4}
}

func (m *CRSummary) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CRSummary.Unmarshal(m, b)
}
func (m *CRSummary) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CRSummary.Marshal(b, m, deterministic)
}
func (m *CRSummary) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CRSummary.Merge(m, src)
}
func (m *CRSummary) XXX_Size() int {
	return xxx_messageInfo_CRSummary.Size(m)
}
func (m *CRSummary) XXX_DiscardUnknown() {
	xxx_messageInfo_CRSummary.DiscardUnknown(m)
}

var xxx_messageInfo_CRSummary proto.InternalMessageInfo

func (m *CRSummary) GetCrName() string {
	if m != nil {
		return m.CrName
	}
	return ""
}

func (m *CRSummary) GetCrKind() string {
	if m != nil {
		return m.CrKind
	}
	return ""
}

func (m *CRSummary) GetCrVersion() string {
	if m != nil {
		return m.CrVersion
	}
	return ""
}

type DeltaDiscoveryRequest struct {
	Node                     *core.Node        `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
	TypeUrl                  string            `protobuf:"bytes,2,opt,name=type_url,json=typeUrl,proto3" json:"type_url,omitempty"`
	ResourceNamesSubscribe   []string          `protobuf:"bytes,3,rep,name=resource_names_subscribe,json=resourceNamesSubscribe,proto3" json:"resource_names_subscribe,omitempty"`
	ResourceNamesUnsubscribe []string          `protobuf:"bytes,4,rep,name=resource_names_unsubscribe,json=resourceNamesUnsubscribe,proto3" json:"resource_names_unsubscribe,omitempty"`
	InitialResourceVersions  map[string]string `protobuf:"bytes,5,rep,name=initial_resource_versions,json=initialResourceVersions,proto3" json:"initial_resource_versions,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	ResponseNonce            string            `protobuf:"bytes,6,opt,name=response_nonce,json=responseNonce,proto3" json:"response_nonce,omitempty"`
	ErrorDetail              *status.Status    `protobuf:"bytes,7,opt,name=error_detail,json=errorDetail,proto3" json:"error_detail,omitempty"`
	XXX_NoUnkeyedLiteral     struct{}          `json:"-"`
	XXX_unrecognized         []byte            `json:"-"`
	XXX_sizecache            int32             `json:"-"`
}

func (m *DeltaDiscoveryRequest) Reset()         { *m = DeltaDiscoveryRequest{} }
func (m *DeltaDiscoveryRequest) String() string { return proto.CompactTextString(m) }
func (*DeltaDiscoveryRequest) ProtoMessage()    {}
func (*DeltaDiscoveryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2c7365e287e5c035, []int{5}
}

func (m *DeltaDiscoveryRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeltaDiscoveryRequest.Unmarshal(m, b)
}
func (m *DeltaDiscoveryRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeltaDiscoveryRequest.Marshal(b, m, deterministic)
}
func (m *DeltaDiscoveryRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeltaDiscoveryRequest.Merge(m, src)
}
func (m *DeltaDiscoveryRequest) XXX_Size() int {
	return xxx_messageInfo_DeltaDiscoveryRequest.Size(m)
}
func (m *DeltaDiscoveryRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeltaDiscoveryRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeltaDiscoveryRequest proto.InternalMessageInfo

func (m *DeltaDiscoveryRequest) GetNode() *core.Node {
	if m != nil {
		return m.Node
	}
	return nil
}

func (m *DeltaDiscoveryRequest) GetTypeUrl() string {
	if m != nil {
		return m.TypeUrl
	}
	return ""
}

func (m *DeltaDiscoveryRequest) GetResourceNamesSubscribe() []string {
	if m != nil {
		return m.ResourceNamesSubscribe
	}
	return nil
}

func (m *DeltaDiscoveryRequest) GetResourceNamesUnsubscribe() []string {
	if m != nil {
		return m.ResourceNamesUnsubscribe
	}
	return nil
}

func (m *DeltaDiscoveryRequest) GetInitialResourceVersions() map[string]string {
	if m != nil {
		return m.InitialResourceVersions
	}
	return nil
}

func (m *DeltaDiscoveryRequest) GetResponseNonce() string {
	if m != nil {
		return m.ResponseNonce
	}
	return ""
}

func (m *DeltaDiscoveryRequest) GetErrorDetail() *status.Status {
	if m != nil {
		return m.ErrorDetail
	}
	return nil
}

type DeltaDiscoveryResponse struct {
	SystemVersionInfo    string      `protobuf:"bytes,1,opt,name=system_version_info,json=systemVersionInfo,proto3" json:"system_version_info,omitempty"`
	Resources            []*Resource `protobuf:"bytes,2,rep,name=resources,proto3" json:"resources,omitempty"`
	TypeUrl              string      `protobuf:"bytes,4,opt,name=type_url,json=typeUrl,proto3" json:"type_url,omitempty"`
	RemovedResources     []string    `protobuf:"bytes,6,rep,name=removed_resources,json=removedResources,proto3" json:"removed_resources,omitempty"`
	Nonce                string      `protobuf:"bytes,5,opt,name=nonce,proto3" json:"nonce,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *DeltaDiscoveryResponse) Reset()         { *m = DeltaDiscoveryResponse{} }
func (m *DeltaDiscoveryResponse) String() string { return proto.CompactTextString(m) }
func (*DeltaDiscoveryResponse) ProtoMessage()    {}
func (*DeltaDiscoveryResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2c7365e287e5c035, []int{6}
}

func (m *DeltaDiscoveryResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeltaDiscoveryResponse.Unmarshal(m, b)
}
func (m *DeltaDiscoveryResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeltaDiscoveryResponse.Marshal(b, m, deterministic)
}
func (m *DeltaDiscoveryResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeltaDiscoveryResponse.Merge(m, src)
}
func (m *DeltaDiscoveryResponse) XXX_Size() int {
	return xxx_messageInfo_DeltaDiscoveryResponse.Size(m)
}
func (m *DeltaDiscoveryResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeltaDiscoveryResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeltaDiscoveryResponse proto.InternalMessageInfo

func (m *DeltaDiscoveryResponse) GetSystemVersionInfo() string {
	if m != nil {
		return m.SystemVersionInfo
	}
	return ""
}

func (m *DeltaDiscoveryResponse) GetResources() []*Resource {
	if m != nil {
		return m.Resources
	}
	return nil
}

func (m *DeltaDiscoveryResponse) GetTypeUrl() string {
	if m != nil {
		return m.TypeUrl
	}
	return ""
}

func (m *DeltaDiscoveryResponse) GetRemovedResources() []string {
	if m != nil {
		return m.RemovedResources
	}
	return nil
}

func (m *DeltaDiscoveryResponse) GetNonce() string {
	if m != nil {
		return m.Nonce
	}
	return ""
}

type Resource struct {
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Aliases              []string `protobuf:"bytes,4,rep,name=aliases,proto3" json:"aliases,omitempty"`
	Version              string   `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	Resource             *any.Any `protobuf:"bytes,2,opt,name=resource,proto3" json:"resource,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Resource) Reset()         { *m = Resource{} }
func (m *Resource) String() string { return proto.CompactTextString(m) }
func (*Resource) ProtoMessage()    {}
func (*Resource) Descriptor() ([]byte, []int) {
	return fileDescriptor_2c7365e287e5c035, []int{7}
}

func (m *Resource) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Resource.Unmarshal(m, b)
}
func (m *Resource) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Resource.Marshal(b, m, deterministic)
}
func (m *Resource) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Resource.Merge(m, src)
}
func (m *Resource) XXX_Size() int {
	return xxx_messageInfo_Resource.Size(m)
}
func (m *Resource) XXX_DiscardUnknown() {
	xxx_messageInfo_Resource.DiscardUnknown(m)
}

var xxx_messageInfo_Resource proto.InternalMessageInfo

func (m *Resource) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Resource) GetAliases() []string {
	if m != nil {
		return m.Aliases
	}
	return nil
}

func (m *Resource) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *Resource) GetResource() *any.Any {
	if m != nil {
		return m.Resource
	}
	return nil
}

func init() {
	proto.RegisterType((*DiscoveryRequest)(nil), "envoy.api.v2.DiscoveryRequest")
	proto.RegisterType((*DiscoveryResponse)(nil), "envoy.api.v2.DiscoveryResponse")
	proto.RegisterType((*DiscoveryResponseExpand)(nil), "envoy.api.v2.DiscoveryResponseExpand")
	proto.RegisterType((*CheckMessage)(nil), "envoy.api.v2.CheckMessage")
	proto.RegisterType((*CRSummary)(nil), "envoy.api.v2.CRSummary")
	proto.RegisterType((*DeltaDiscoveryRequest)(nil), "envoy.api.v2.DeltaDiscoveryRequest")
	proto.RegisterMapType((map[string]string)(nil), "envoy.api.v2.DeltaDiscoveryRequest.InitialResourceVersionsEntry")
	proto.RegisterType((*DeltaDiscoveryResponse)(nil), "envoy.api.v2.DeltaDiscoveryResponse")
	proto.RegisterType((*Resource)(nil), "envoy.api.v2.Resource")
}

func init() { proto.RegisterFile("envoy/api/v2/discovery.proto", fileDescriptor_2c7365e287e5c035) }

var fileDescriptor_2c7365e287e5c035 = []byte{
	// 856 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x55, 0xcf, 0x6b, 0xe3, 0x46,
	0x14, 0x46, 0x76, 0xe2, 0xc4, 0xcf, 0xce, 0x92, 0x4c, 0xb7, 0xb1, 0x62, 0xd2, 0xd6, 0x35, 0x2c,
	0x18, 0x16, 0xe4, 0xe2, 0x6d, 0x21, 0x94, 0x42, 0xdb, 0x4d, 0xb6, 0xec, 0xb6, 0x34, 0x18, 0x85,
	0xdd, 0xc3, 0x52, 0x10, 0xe3, 0xd1, 0x4b, 0x3a, 0x44, 0x9e, 0x51, 0x67, 0x24, 0xb1, 0x82, 0x9e,
	0x4a, 0xef, 0xbd, 0xf6, 0x8f, 0xeb, 0x5f, 0xd1, 0x63, 0x0f, 0x6d, 0xd1, 0x68, 0x64, 0x4b, 0x89,
	0x37, 0xf8, 0xa6, 0x37, 0xdf, 0xf7, 0xde, 0xcc, 0xfb, 0xde, 0x0f, 0xc1, 0x29, 0x8a, 0x4c, 0xe6,
	0x53, 0x1a, 0xf3, 0x69, 0x36, 0x9b, 0x86, 0x5c, 0x33, 0x99, 0xa1, 0xca, 0xbd, 0x58, 0xc9, 0x44,
	0x92, 0xbe, 0x41, 0x3d, 0x1a, 0x73, 0x2f, 0x9b, 0x0d, 0x9b, 0x5c, 0x26, 0x15, 0x4e, 0x17, 0x54,
	0x63, 0xc9, 0x1d, 0x9e, 0xdc, 0x48, 0x79, 0x13, 0xe1, 0xd4, 0x58, 0x8b, 0xf4, 0x7a, 0x4a, 0x85,
	0x0d, 0x33, 0x1c, 0x58, 0x48, 0xc5, 0x6c, 0xaa, 0x13, 0x9a, 0xa4, 0xda, 0x02, 0x1f, 0xa7, 0x61,
	0x4c, 0xa7, 0x54, 0x08, 0x99, 0xd0, 0x84, 0x4b, 0xa1, 0xa7, 0x4b, 0x7e, 0xa3, 0x68, 0x62, 0x63,
	0x8e, 0x7f, 0x6b, 0xc1, 0xe1, 0x45, 0xf5, 0x26, 0x1f, 0x7f, 0x49, 0x51, 0x27, 0xe4, 0x53, 0xe8,
	0x67, 0xa8, 0x34, 0x97, 0x22, 0xe0, 0xe2, 0x5a, 0xba, 0xce, 0xc8, 0x99, 0x74, 0xfd, 0x9e, 0x3d,
	0x7b, 0x25, 0xae, 0x25, 0x79, 0x0a, 0x3b, 0x42, 0x86, 0xe8, 0xb6, 0x46, 0xce, 0xa4, 0x37, 0x1b,
	0x78, 0xf5, 0x34, 0xbc, 0xe2, 0xe1, 0xde, 0xa5, 0x0c, 0xd1, 0x37, 0x24, 0xf2, 0x04, 0x1e, 0x29,
	0xd4, 0x32, 0x55, 0x0c, 0x03, 0x41, 0x97, 0xa8, 0xdd, 0xf6, 0xa8, 0x3d, 0xe9, 0xfa, 0x07, 0xd5,
	0xe9, 0x65, 0x71, 0x48, 0x4e, 0x60, 0x3f, 0xc9, 0x63, 0x0c, 0x52, 0x15, 0xb9, 0x3b, 0xe6, 0xca,
	0xbd, 0xc2, 0x7e, 0xad, 0x22, 0x1b, 0x21, 0x96, 0x42, 0x63, 0x20, 0xa4, 0x60, 0xe8, 0xee, 0x1a,
	0xc2, 0x41, 0x75, 0x7a, 0x59, 0x1c, 0x92, 0x2f, 0xa0, 0x8f, 0x4a, 0x49, 0x15, 0x84, 0x98, 0x50,
	0x1e, 0xb9, 0x1d, 0xf3, 0x3a, 0xe2, 0x95, 0xea, 0x78, 0x2a, 0x66, 0xde, 0x95, 0x51, 0xc7, 0xef,
	0x19, 0xde, 0x85, 0xa1, 0x8d, 0xff, 0x69, 0xc1, 0x51, 0x4d, 0x84, 0x32, 0xe2, 0x36, 0x2a, 0xcc,
	0xa0, 0x5b, 0xa5, 0xa0, 0xdd, 0xd6, 0xa8, 0x3d, 0xe9, 0xcd, 0x1e, 0x57, 0x97, 0x55, 0x55, 0xf2,
	0xbe, 0x15, 0xb9, 0xbf, 0xa6, 0x91, 0x63, 0xe8, 0x30, 0x2a, 0xa8, 0xca, 0xdd, 0xf6, 0xc8, 0x99,
	0xec, 0xfb, 0xd6, 0x7a, 0x28, 0xfb, 0xc7, 0xb0, 0x5b, 0x4f, 0xba, 0x34, 0xc8, 0x05, 0x1c, 0x30,
	0x29, 0x12, 0x25, 0xa3, 0x20, 0x8e, 0xa8, 0x40, 0x9b, 0xed, 0x27, 0x1b, 0x6a, 0x71, 0x5e, 0xf2,
	0xe6, 0x05, 0xcd, 0xef, 0xb3, 0x9a, 0x75, 0x4f, 0xb2, 0xf9, 0x56, 0x92, 0x91, 0x97, 0xd0, 0xc7,
	0x77, 0x31, 0x15, 0x61, 0x70, 0xcd, 0x31, 0x0a, 0xdd, 0xb7, 0xc6, 0xed, 0x49, 0xf3, 0xee, 0x7b,
	0x9a, 0xbe, 0x30, 0x2e, 0x7e, 0xaf, 0x74, 0xfd, 0xae, 0xf0, 0x1c, 0xbf, 0x85, 0xc1, 0x7b, 0x78,
	0xe4, 0x6b, 0x38, 0x60, 0x3f, 0x23, 0xbb, 0x0d, 0x96, 0xa8, 0x35, 0xbd, 0x41, 0x53, 0x82, 0xde,
	0x6c, 0xd8, 0xbc, 0xe5, 0xbc, 0xa0, 0xfc, 0x58, 0x32, 0xfc, 0x3e, 0xab, 0x59, 0xe3, 0x97, 0xd0,
	0xaf, 0xa3, 0xe4, 0x0c, 0x7a, 0x4c, 0x05, 0x3a, 0x5d, 0x2e, 0xa9, 0xca, 0xb5, 0xeb, 0x98, 0x8a,
	0xdd, 0x69, 0xde, 0x73, 0xff, 0xaa, 0xc4, 0x7d, 0x60, 0xca, 0x7e, 0xea, 0xf1, 0x4f, 0xd0, 0x5d,
	0x01, 0x64, 0x00, 0x7b, 0x4c, 0x99, 0x4e, 0xb6, 0x4d, 0xd1, 0x61, 0xaa, 0x68, 0x61, 0x0b, 0xdc,
	0x72, 0x11, 0x9a, 0xc1, 0x30, 0xc0, 0x0f, 0x5c, 0x84, 0xe4, 0x23, 0x00, 0xa6, 0x02, 0xdb, 0x3a,
	0xa6, 0xf0, 0x5d, 0xbf, 0xcb, 0xd4, 0x9b, 0xf2, 0x60, 0xfc, 0x5f, 0x1b, 0x3e, 0xbc, 0xc0, 0x28,
	0xa1, 0xf7, 0x46, 0xb1, 0x9a, 0x33, 0x67, 0x9b, 0x39, 0xab, 0xb7, 0x50, 0xab, 0xd9, 0x42, 0x67,
	0xe0, 0x36, 0x47, 0x30, 0xd0, 0xe9, 0x42, 0x33, 0xc5, 0x17, 0x68, 0x87, 0xf1, 0xb8, 0x31, 0x8c,
	0x57, 0x15, 0x4a, 0xbe, 0x82, 0xe1, 0x1d, 0xcf, 0x54, 0xac, 0x7d, 0x77, 0x8c, 0xaf, 0xdb, 0xf0,
	0x7d, 0xbd, 0xc6, 0xc9, 0xaf, 0x70, 0xc2, 0x05, 0x4f, 0x38, 0x8d, 0x82, 0x55, 0x14, 0x2b, 0x83,
	0x76, 0x77, 0x8d, 0xfe, 0xdf, 0xdc, 0x69, 0x9a, 0x4d, 0x3a, 0x78, 0xaf, 0xca, 0x20, 0xbe, 0x8d,
	0x61, 0x85, 0xd3, 0x2f, 0x44, 0xa2, 0x72, 0x7f, 0xc0, 0x37, 0xa3, 0x1b, 0xd6, 0x46, 0x67, 0x9b,
	0xb5, 0xb1, 0xb7, 0xd5, 0x0c, 0x0c, 0xbf, 0x87, 0xd3, 0x87, 0x9e, 0x45, 0x0e, 0xa1, 0x7d, 0x8b,
	0xb9, 0x6d, 0x91, 0xe2, 0xb3, 0x18, 0xe4, 0x8c, 0x46, 0x29, 0xda, 0xea, 0x94, 0xc6, 0x97, 0xad,
	0x33, 0x67, 0xfc, 0x97, 0x03, 0xc7, 0x77, 0x33, 0xb7, 0x7b, 0xc8, 0x83, 0x0f, 0x74, 0xae, 0x13,
	0x5c, 0x06, 0x1b, 0xd6, 0xd1, 0x51, 0x09, 0xbd, 0xa9, 0x2d, 0xa5, 0xcf, 0xef, 0x2f, 0xa5, 0xe3,
	0xa6, 0xc4, 0xd5, 0x73, 0xeb, 0x6b, 0xe9, 0x81, 0xf5, 0xf3, 0x14, 0x8e, 0x14, 0x2e, 0x65, 0x86,
	0x61, 0xb0, 0x0e, 0xdc, 0x31, 0x85, 0x3f, 0xb4, 0x80, 0xbf, 0x8a, 0xb3, 0x71, 0x57, 0x8d, 0x7f,
	0x77, 0x60, 0xbf, 0xe2, 0x10, 0x02, 0x3b, 0x66, 0x76, 0xca, 0x31, 0x30, 0xdf, 0xc4, 0x85, 0x3d,
	0x1a, 0x71, 0xaa, 0x51, 0xdb, 0x96, 0xaa, 0xcc, 0x02, 0xa9, 0xe6, 0xa6, 0x4c, 0xb9, 0x32, 0xc9,
	0x67, 0xb0, 0x5f, 0xbd, 0xc7, 0xfe, 0x87, 0x36, 0x2f, 0xdf, 0x15, 0xeb, 0xf9, 0xfc, 0xef, 0x3f,
	0xff, 0xfd, 0x63, 0xf7, 0x94, 0x0c, 0x4b, 0x39, 0x34, 0xaa, 0x8c, 0x33, 0xf4, 0xd6, 0x3f, 0xe5,
	0xec, 0x19, 0x0c, 0xb9, 0x2c, 0xd5, 0x8a, 0x95, 0x7c, 0x97, 0x37, 0x84, 0x7b, 0xfe, 0x68, 0x55,
	0x9d, 0x79, 0x71, 0xc1, 0xdc, 0x59, 0x74, 0xcc, 0x4d, 0xcf, 0xfe, 0x0f, 0x00, 0x00, 0xff, 0xff,
	0xa3, 0xea, 0x50, 0x3c, 0xe6, 0x07, 0x00, 0x00,
}
