// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: service_config.proto

package types

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type MISTServiceConfig struct {
	// appname
	AppName string `protobuf:"bytes,1,opt,name=app_name,json=appName,proto3" json:"app_name,omitempty"`
	// mist master flag
	Enable bool `protobuf:"varint,2,opt,name=enable,proto3" json:"enable,omitempty"`
	// version
	Version string `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	// log level
	LogLevel int32 `protobuf:"varint,4,opt,name=log_level,json=logLevel,proto3" json:"log_level,omitempty"`
	// log path
	LogPath string `protobuf:"bytes,5,opt,name=log_path,json=logPath,proto3" json:"log_path,omitempty"`
	// biz service enable flag
	BizEnable bool `protobuf:"varint,6,opt,name=biz_enable,json=bizEnable,proto3" json:"biz_enable,omitempty"`
	// admin service enable flag
	AdminEnable bool `protobuf:"varint,7,opt,name=admin_enable,json=adminEnable,proto3" json:"admin_enable,omitempty"`
	// mist plugin stop
	PluginEnable bool `protobuf:"varint,8,opt,name=plugin_enable,json=pluginEnable,proto3" json:"plugin_enable,omitempty"`
	// go max procs
	MaxProcs             int32    `protobuf:"varint,9,opt,name=max_procs,json=maxProcs,proto3" json:"max_procs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MISTServiceConfig) Reset()         { *m = MISTServiceConfig{} }
func (m *MISTServiceConfig) String() string { return proto.CompactTextString(m) }
func (*MISTServiceConfig) ProtoMessage()    {}
func (*MISTServiceConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_452b382cbc0cfd24, []int{0}
}
func (m *MISTServiceConfig) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MISTServiceConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MISTServiceConfig.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MISTServiceConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MISTServiceConfig.Merge(m, src)
}
func (m *MISTServiceConfig) XXX_Size() int {
	return m.Size()
}
func (m *MISTServiceConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_MISTServiceConfig.DiscardUnknown(m)
}

var xxx_messageInfo_MISTServiceConfig proto.InternalMessageInfo

func (m *MISTServiceConfig) GetAppName() string {
	if m != nil {
		return m.AppName
	}
	return ""
}

func (m *MISTServiceConfig) GetEnable() bool {
	if m != nil {
		return m.Enable
	}
	return false
}

func (m *MISTServiceConfig) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *MISTServiceConfig) GetLogLevel() int32 {
	if m != nil {
		return m.LogLevel
	}
	return 0
}

func (m *MISTServiceConfig) GetLogPath() string {
	if m != nil {
		return m.LogPath
	}
	return ""
}

func (m *MISTServiceConfig) GetBizEnable() bool {
	if m != nil {
		return m.BizEnable
	}
	return false
}

func (m *MISTServiceConfig) GetAdminEnable() bool {
	if m != nil {
		return m.AdminEnable
	}
	return false
}

func (m *MISTServiceConfig) GetPluginEnable() bool {
	if m != nil {
		return m.PluginEnable
	}
	return false
}

func (m *MISTServiceConfig) GetMaxProcs() int32 {
	if m != nil {
		return m.MaxProcs
	}
	return 0
}

func init() {
	proto.RegisterType((*MISTServiceConfig)(nil), "meshapi.security.mist.v1.MISTServiceConfig")
}

func init() { proto.RegisterFile("service_config.proto", fileDescriptor_452b382cbc0cfd24) }

var fileDescriptor_452b382cbc0cfd24 = []byte{
	// 288 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x90, 0x41, 0x4a, 0xf4, 0x30,
	0x18, 0x86, 0xc9, 0xfc, 0xff, 0xcc, 0xb4, 0x71, 0x5c, 0x18, 0x45, 0x22, 0x83, 0xa5, 0xea, 0xa6,
	0xab, 0x82, 0x78, 0x03, 0xc5, 0x85, 0xa0, 0x32, 0x74, 0x5c, 0xb9, 0x29, 0x69, 0xfd, 0x6c, 0x03,
	0x49, 0x13, 0x9a, 0x4c, 0x99, 0x99, 0xb3, 0x78, 0x20, 0x97, 0x1e, 0x41, 0x7a, 0x12, 0x69, 0x52,
	0x71, 0xf9, 0xbe, 0x79, 0x5e, 0x9e, 0xf0, 0xe1, 0x13, 0x03, 0x6d, 0xc7, 0x4b, 0xc8, 0x4b, 0xd5,
	0xbc, 0xf3, 0x2a, 0xd5, 0xad, 0xb2, 0x8a, 0x50, 0x09, 0xa6, 0x66, 0x9a, 0xa7, 0x06, 0xca, 0x4d,
	0xcb, 0xed, 0x2e, 0x95, 0xdc, 0xd8, 0xb4, 0xbb, 0xbe, 0xfc, 0x98, 0xe0, 0xa3, 0xa7, 0x87, 0xf5,
	0xcb, 0xda, 0xcf, 0xee, 0xdc, 0x8a, 0x9c, 0xe1, 0x80, 0x69, 0x9d, 0x37, 0x4c, 0x02, 0x45, 0x31,
	0x4a, 0xc2, 0x6c, 0xce, 0xb4, 0x7e, 0x66, 0x12, 0xc8, 0x29, 0x9e, 0x41, 0xc3, 0x0a, 0x01, 0x74,
	0x12, 0xa3, 0x24, 0xc8, 0xc6, 0x44, 0x28, 0x9e, 0x77, 0xd0, 0x1a, 0xae, 0x1a, 0xfa, 0xcf, 0x2f,
	0xc6, 0x48, 0x96, 0x38, 0x14, 0xaa, 0xca, 0x05, 0x74, 0x20, 0xe8, 0xff, 0x18, 0x25, 0xd3, 0x2c,
	0x10, 0xaa, 0x7a, 0x1c, 0xf2, 0x60, 0x1a, 0x1e, 0x35, 0xb3, 0x35, 0x9d, 0xfa, 0x9d, 0x50, 0xd5,
	0x8a, 0xd9, 0x9a, 0x9c, 0x63, 0x5c, 0xf0, 0x7d, 0x3e, 0xda, 0x66, 0xce, 0x16, 0x16, 0x7c, 0x7f,
	0xef, 0x85, 0x17, 0x78, 0xc1, 0xde, 0x24, 0x6f, 0x7e, 0x81, 0xb9, 0x03, 0x0e, 0x5c, 0x37, 0x22,
	0x57, 0xf8, 0x50, 0x8b, 0x4d, 0xf5, 0xc7, 0x04, 0x8e, 0x59, 0xf8, 0x72, 0x84, 0x96, 0x38, 0x94,
	0x6c, 0x9b, 0xeb, 0x56, 0x95, 0x86, 0x86, 0xfe, 0x7b, 0x92, 0x6d, 0x57, 0x43, 0xbe, 0x3d, 0xfe,
	0xec, 0x23, 0xf4, 0xd5, 0x47, 0xe8, 0xbb, 0x8f, 0xd0, 0xeb, 0xd4, 0xee, 0x34, 0x98, 0x62, 0xe6,
	0x8e, 0x7a, 0xf3, 0x13, 0x00, 0x00, 0xff, 0xff, 0x41, 0xdf, 0xed, 0xa5, 0x6c, 0x01, 0x00, 0x00,
}

func (m *MISTServiceConfig) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MISTServiceConfig) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MISTServiceConfig) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.MaxProcs != 0 {
		i = encodeVarintServiceConfig(dAtA, i, uint64(m.MaxProcs))
		i--
		dAtA[i] = 0x48
	}
	if m.PluginEnable {
		i--
		if m.PluginEnable {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x40
	}
	if m.AdminEnable {
		i--
		if m.AdminEnable {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x38
	}
	if m.BizEnable {
		i--
		if m.BizEnable {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x30
	}
	if len(m.LogPath) > 0 {
		i -= len(m.LogPath)
		copy(dAtA[i:], m.LogPath)
		i = encodeVarintServiceConfig(dAtA, i, uint64(len(m.LogPath)))
		i--
		dAtA[i] = 0x2a
	}
	if m.LogLevel != 0 {
		i = encodeVarintServiceConfig(dAtA, i, uint64(m.LogLevel))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Version) > 0 {
		i -= len(m.Version)
		copy(dAtA[i:], m.Version)
		i = encodeVarintServiceConfig(dAtA, i, uint64(len(m.Version)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Enable {
		i--
		if m.Enable {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if len(m.AppName) > 0 {
		i -= len(m.AppName)
		copy(dAtA[i:], m.AppName)
		i = encodeVarintServiceConfig(dAtA, i, uint64(len(m.AppName)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintServiceConfig(dAtA []byte, offset int, v uint64) int {
	offset -= sovServiceConfig(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MISTServiceConfig) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.AppName)
	if l > 0 {
		n += 1 + l + sovServiceConfig(uint64(l))
	}
	if m.Enable {
		n += 2
	}
	l = len(m.Version)
	if l > 0 {
		n += 1 + l + sovServiceConfig(uint64(l))
	}
	if m.LogLevel != 0 {
		n += 1 + sovServiceConfig(uint64(m.LogLevel))
	}
	l = len(m.LogPath)
	if l > 0 {
		n += 1 + l + sovServiceConfig(uint64(l))
	}
	if m.BizEnable {
		n += 2
	}
	if m.AdminEnable {
		n += 2
	}
	if m.PluginEnable {
		n += 2
	}
	if m.MaxProcs != 0 {
		n += 1 + sovServiceConfig(uint64(m.MaxProcs))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovServiceConfig(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozServiceConfig(x uint64) (n int) {
	return sovServiceConfig(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MISTServiceConfig) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowServiceConfig
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MISTServiceConfig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MISTServiceConfig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AppName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServiceConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServiceConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthServiceConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AppName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Enable", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServiceConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Enable = bool(v != 0)
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServiceConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServiceConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthServiceConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Version = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LogLevel", wireType)
			}
			m.LogLevel = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServiceConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LogLevel |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LogPath", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServiceConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthServiceConfig
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthServiceConfig
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LogPath = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BizEnable", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServiceConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.BizEnable = bool(v != 0)
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AdminEnable", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServiceConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.AdminEnable = bool(v != 0)
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PluginEnable", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServiceConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.PluginEnable = bool(v != 0)
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxProcs", wireType)
			}
			m.MaxProcs = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowServiceConfig
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxProcs |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipServiceConfig(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthServiceConfig
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthServiceConfig
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipServiceConfig(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowServiceConfig
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowServiceConfig
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowServiceConfig
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthServiceConfig
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupServiceConfig
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthServiceConfig
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthServiceConfig        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowServiceConfig          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupServiceConfig = fmt.Errorf("proto: unexpected end of group")
)