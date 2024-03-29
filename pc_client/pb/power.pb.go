// Code generated by protoc-gen-go. DO NOT EDIT.
// source: power.proto

package pb

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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Power struct {
	CurrentRms           float32  `protobuf:"fixed32,1,opt,name=CurrentRms,json=currentRms,proto3" json:"CurrentRms,omitempty"`
	RealPower            float32  `protobuf:"fixed32,2,opt,name=RealPower,json=realPower,proto3" json:"RealPower,omitempty"`
	ApparentPower        float32  `protobuf:"fixed32,3,opt,name=ApparentPower,json=apparentPower,proto3" json:"ApparentPower,omitempty"`
	PowerFactor          float32  `protobuf:"fixed32,4,opt,name=PowerFactor,json=powerFactor,proto3" json:"PowerFactor,omitempty"`
	DC                   float32  `protobuf:"fixed32,5,opt,name=DC,json=dC,proto3" json:"DC,omitempty"`
	Current              []uint32 `protobuf:"varint,6,rep,packed,name=Current,json=current,proto3" json:"Current,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Power) Reset()         { *m = Power{} }
func (m *Power) String() string { return proto.CompactTextString(m) }
func (*Power) ProtoMessage()    {}
func (*Power) Descriptor() ([]byte, []int) {
	return fileDescriptor_a4fab2da8ea5416b, []int{0}
}

func (m *Power) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Power.Unmarshal(m, b)
}
func (m *Power) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Power.Marshal(b, m, deterministic)
}
func (m *Power) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Power.Merge(m, src)
}
func (m *Power) XXX_Size() int {
	return xxx_messageInfo_Power.Size(m)
}
func (m *Power) XXX_DiscardUnknown() {
	xxx_messageInfo_Power.DiscardUnknown(m)
}

var xxx_messageInfo_Power proto.InternalMessageInfo

func (m *Power) GetCurrentRms() float32 {
	if m != nil {
		return m.CurrentRms
	}
	return 0
}

func (m *Power) GetRealPower() float32 {
	if m != nil {
		return m.RealPower
	}
	return 0
}

func (m *Power) GetApparentPower() float32 {
	if m != nil {
		return m.ApparentPower
	}
	return 0
}

func (m *Power) GetPowerFactor() float32 {
	if m != nil {
		return m.PowerFactor
	}
	return 0
}

func (m *Power) GetDC() float32 {
	if m != nil {
		return m.DC
	}
	return 0
}

func (m *Power) GetCurrent() []uint32 {
	if m != nil {
		return m.Current
	}
	return nil
}

type PowerMessage struct {
	Powers               []*Power `protobuf:"bytes,1,rep,name=Powers,json=powers,proto3" json:"Powers,omitempty"`
	VoltageRms           float32  `protobuf:"fixed32,2,opt,name=VoltageRms,json=voltageRms,proto3" json:"VoltageRms,omitempty"`
	DC                   float32  `protobuf:"fixed32,3,opt,name=DC,json=dC,proto3" json:"DC,omitempty"`
	Voltage              []uint32 `protobuf:"varint,4,rep,packed,name=Voltage,json=voltage,proto3" json:"Voltage,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PowerMessage) Reset()         { *m = PowerMessage{} }
func (m *PowerMessage) String() string { return proto.CompactTextString(m) }
func (*PowerMessage) ProtoMessage()    {}
func (*PowerMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_a4fab2da8ea5416b, []int{1}
}

func (m *PowerMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PowerMessage.Unmarshal(m, b)
}
func (m *PowerMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PowerMessage.Marshal(b, m, deterministic)
}
func (m *PowerMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PowerMessage.Merge(m, src)
}
func (m *PowerMessage) XXX_Size() int {
	return xxx_messageInfo_PowerMessage.Size(m)
}
func (m *PowerMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_PowerMessage.DiscardUnknown(m)
}

var xxx_messageInfo_PowerMessage proto.InternalMessageInfo

func (m *PowerMessage) GetPowers() []*Power {
	if m != nil {
		return m.Powers
	}
	return nil
}

func (m *PowerMessage) GetVoltageRms() float32 {
	if m != nil {
		return m.VoltageRms
	}
	return 0
}

func (m *PowerMessage) GetDC() float32 {
	if m != nil {
		return m.DC
	}
	return 0
}

func (m *PowerMessage) GetVoltage() []uint32 {
	if m != nil {
		return m.Voltage
	}
	return nil
}

func init() {
	proto.RegisterType((*Power)(nil), "Power")
	proto.RegisterType((*PowerMessage)(nil), "PowerMessage")
}

func init() { proto.RegisterFile("power.proto", fileDescriptor_a4fab2da8ea5416b) }

var fileDescriptor_a4fab2da8ea5416b = []byte{
	// 263 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0x31, 0x4e, 0xc3, 0x30,
	0x14, 0x86, 0x65, 0x3b, 0x49, 0xd5, 0x97, 0x04, 0x55, 0x9e, 0x2c, 0x84, 0x50, 0x54, 0x01, 0xca,
	0xd4, 0x01, 0x0e, 0x10, 0x91, 0x20, 0x36, 0x24, 0xe4, 0x81, 0x81, 0xcd, 0x29, 0x56, 0x97, 0x12,
	0x5b, 0x4e, 0x80, 0x95, 0x0b, 0xb0, 0xe4, 0x24, 0x5c, 0x82, 0x7b, 0x21, 0x3f, 0x9b, 0xaa, 0x03,
	0x5b, 0xde, 0xff, 0xbe, 0xd8, 0xdf, 0x6f, 0xc8, 0xad, 0xf9, 0xd0, 0x6e, 0x63, 0x9d, 0x99, 0xcc,
	0x69, 0x31, 0xa8, 0xc1, 0xd8, 0x3e, 0x4c, 0xeb, 0x1f, 0x02, 0xe9, 0xa3, 0xdf, 0xf2, 0x73, 0x80,
	0xee, 0xcd, 0x39, 0x3d, 0x4c, 0xf2, 0x75, 0x14, 0xa4, 0x22, 0x35, 0x95, 0xb0, 0x3d, 0x24, 0xfc,
	0x0c, 0x96, 0x52, 0xab, 0x3d, 0xc2, 0x82, 0xe2, 0x7a, 0xe9, 0xfe, 0x02, 0x7e, 0x01, 0xe5, 0xad,
	0xb5, 0xca, 0xc3, 0x81, 0x60, 0x48, 0x94, 0xea, 0x38, 0xe4, 0x15, 0xe4, 0xf8, 0x71, 0xaf, 0xb6,
	0x93, 0x71, 0x22, 0x41, 0x26, 0xd8, 0x85, 0x88, 0x9f, 0x00, 0xbd, 0xeb, 0x44, 0x8a, 0x0b, 0xfa,
	0xd2, 0xf1, 0x2b, 0x58, 0x44, 0x2b, 0x91, 0x55, 0xac, 0x2e, 0xdb, 0x62, 0x6e, 0xd8, 0xea, 0x9b,
	0xcc, 0x0d, 0xfb, 0x24, 0x44, 0x2e, 0xa2, 0xe0, 0xfa, 0x8b, 0x40, 0x81, 0x47, 0x3f, 0xe8, 0x71,
	0x54, 0x3b, 0xcd, 0x2f, 0x21, 0xc3, 0xd9, 0x57, 0x61, 0x75, 0x7e, 0x9d, 0x6d, 0x70, 0x6c, 0xd3,
	0xb9, 0xa1, 0x2b, 0x26, 0x33, 0xbc, 0x74, 0xf4, 0xad, 0x9f, 0xcc, 0x7e, 0x52, 0x3b, 0xed, 0x5b,
	0x87, 0x5a, 0xf0, 0x7e, 0x48, 0xa2, 0x0f, 0x3b, 0xf6, 0x89, 0xbc, 0x48, 0xfe, 0xf3, 0x89, 0xbf,
	0xb6, 0xc9, 0x33, 0xb5, 0x7d, 0x9f, 0xe1, 0x23, 0xdf, 0xfc, 0x06, 0x00, 0x00, 0xff, 0xff, 0x96,
	0xb7, 0xff, 0x55, 0x81, 0x01, 0x00, 0x00,
}
