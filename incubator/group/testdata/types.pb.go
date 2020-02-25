// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: testdata/types.proto

package testdata

import (
	fmt "fmt"
	github_com_cosmos_modules_incubator_group "github.com/cosmos/modules/incubator/group"
	group "github.com/cosmos/modules/incubator/group"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/regen-network/cosmos-proto"
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

type MyAppProposal struct {
	// Types that are valid to be assigned to Sum:
	//	*MyAppProposal_A
	//	*MyAppProposal_B
	Sum isMyAppProposal_Sum `protobuf_oneof:"sum"`
}

func (m *MyAppProposal) Reset()         { *m = MyAppProposal{} }
func (m *MyAppProposal) String() string { return proto.CompactTextString(m) }
func (*MyAppProposal) ProtoMessage()    {}
func (*MyAppProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_2447ab8d7bf628b8, []int{0}
}
func (m *MyAppProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MyAppProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MyAppProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MyAppProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MyAppProposal.Merge(m, src)
}
func (m *MyAppProposal) XXX_Size() int {
	return m.Size()
}
func (m *MyAppProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_MyAppProposal.DiscardUnknown(m)
}

var xxx_messageInfo_MyAppProposal proto.InternalMessageInfo

type isMyAppProposal_Sum interface {
	isMyAppProposal_Sum()
	MarshalTo([]byte) (int, error)
	Size() int
}

type MyAppProposal_A struct {
	A *AMyAppProposal `protobuf:"bytes,1,opt,name=A,proto3,oneof" json:"A,omitempty"`
}
type MyAppProposal_B struct {
	B *BMyAppProposal `protobuf:"bytes,2,opt,name=B,proto3,oneof" json:"B,omitempty"`
}

func (*MyAppProposal_A) isMyAppProposal_Sum() {}
func (*MyAppProposal_B) isMyAppProposal_Sum() {}

func (m *MyAppProposal) GetSum() isMyAppProposal_Sum {
	if m != nil {
		return m.Sum
	}
	return nil
}

func (m *MyAppProposal) GetA() *AMyAppProposal {
	if x, ok := m.GetSum().(*MyAppProposal_A); ok {
		return x.A
	}
	return nil
}

func (m *MyAppProposal) GetB() *BMyAppProposal {
	if x, ok := m.GetSum().(*MyAppProposal_B); ok {
		return x.B
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*MyAppProposal) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*MyAppProposal_A)(nil),
		(*MyAppProposal_B)(nil),
	}
}

type AMyAppProposal struct {
	Base group.ProposalBase `protobuf:"bytes,1,opt,name=base,proto3" json:"base"`
}

func (m *AMyAppProposal) Reset()         { *m = AMyAppProposal{} }
func (m *AMyAppProposal) String() string { return proto.CompactTextString(m) }
func (*AMyAppProposal) ProtoMessage()    {}
func (*AMyAppProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_2447ab8d7bf628b8, []int{1}
}
func (m *AMyAppProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AMyAppProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AMyAppProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AMyAppProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AMyAppProposal.Merge(m, src)
}
func (m *AMyAppProposal) XXX_Size() int {
	return m.Size()
}
func (m *AMyAppProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_AMyAppProposal.DiscardUnknown(m)
}

var xxx_messageInfo_AMyAppProposal proto.InternalMessageInfo

func (m *AMyAppProposal) GetBase() group.ProposalBase {
	if m != nil {
		return m.Base
	}
	return group.ProposalBase{}
}

type BMyAppProposal struct {
	Base group.ProposalBase `protobuf:"bytes,1,opt,name=base,proto3" json:"base"`
}

func (m *BMyAppProposal) Reset()         { *m = BMyAppProposal{} }
func (m *BMyAppProposal) String() string { return proto.CompactTextString(m) }
func (*BMyAppProposal) ProtoMessage()    {}
func (*BMyAppProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_2447ab8d7bf628b8, []int{2}
}
func (m *BMyAppProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BMyAppProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BMyAppProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BMyAppProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BMyAppProposal.Merge(m, src)
}
func (m *BMyAppProposal) XXX_Size() int {
	return m.Size()
}
func (m *BMyAppProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_BMyAppProposal.DiscardUnknown(m)
}

var xxx_messageInfo_BMyAppProposal proto.InternalMessageInfo

func (m *BMyAppProposal) GetBase() group.ProposalBase {
	if m != nil {
		return m.Base
	}
	return group.ProposalBase{}
}

type MsgProposeA struct {
	Base group.MsgProposeBase `protobuf:"bytes,1,opt,name=base,proto3" json:"base"`
}

func (m *MsgProposeA) Reset()         { *m = MsgProposeA{} }
func (m *MsgProposeA) String() string { return proto.CompactTextString(m) }
func (*MsgProposeA) ProtoMessage()    {}
func (*MsgProposeA) Descriptor() ([]byte, []int) {
	return fileDescriptor_2447ab8d7bf628b8, []int{3}
}
func (m *MsgProposeA) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgProposeA) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgProposeA.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgProposeA) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgProposeA.Merge(m, src)
}
func (m *MsgProposeA) XXX_Size() int {
	return m.Size()
}
func (m *MsgProposeA) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgProposeA.DiscardUnknown(m)
}

var xxx_messageInfo_MsgProposeA proto.InternalMessageInfo

func (m *MsgProposeA) GetBase() group.MsgProposeBase {
	if m != nil {
		return m.Base
	}
	return group.MsgProposeBase{}
}

type MsgProposeB struct {
	Base group.MsgProposeBase `protobuf:"bytes,1,opt,name=base,proto3" json:"base"`
}

func (m *MsgProposeB) Reset()         { *m = MsgProposeB{} }
func (m *MsgProposeB) String() string { return proto.CompactTextString(m) }
func (*MsgProposeB) ProtoMessage()    {}
func (*MsgProposeB) Descriptor() ([]byte, []int) {
	return fileDescriptor_2447ab8d7bf628b8, []int{4}
}
func (m *MsgProposeB) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgProposeB) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgProposeB.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgProposeB) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgProposeB.Merge(m, src)
}
func (m *MsgProposeB) XXX_Size() int {
	return m.Size()
}
func (m *MsgProposeB) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgProposeB.DiscardUnknown(m)
}

var xxx_messageInfo_MsgProposeB proto.InternalMessageInfo

func (m *MsgProposeB) GetBase() group.MsgProposeBase {
	if m != nil {
		return m.Base
	}
	return group.MsgProposeBase{}
}

func init() {
	proto.RegisterType((*MyAppProposal)(nil), "cosmos_modules.incubator.group.v1_alpha.testdata.MyAppProposal")
	proto.RegisterType((*AMyAppProposal)(nil), "cosmos_modules.incubator.group.v1_alpha.testdata.AMyAppProposal")
	proto.RegisterType((*BMyAppProposal)(nil), "cosmos_modules.incubator.group.v1_alpha.testdata.BMyAppProposal")
	proto.RegisterType((*MsgProposeA)(nil), "cosmos_modules.incubator.group.v1_alpha.testdata.MsgProposeA")
	proto.RegisterType((*MsgProposeB)(nil), "cosmos_modules.incubator.group.v1_alpha.testdata.MsgProposeB")
}

func init() { proto.RegisterFile("testdata/types.proto", fileDescriptor_2447ab8d7bf628b8) }

var fileDescriptor_2447ab8d7bf628b8 = []byte{
	// 335 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x29, 0x49, 0x2d, 0x2e,
	0x49, 0x49, 0x2c, 0x49, 0xd4, 0x2f, 0xa9, 0x2c, 0x48, 0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0x32, 0x48, 0xce, 0x2f, 0xce, 0xcd, 0x2f, 0x8e, 0xcf, 0xcd, 0x4f, 0x29, 0xcd, 0x49, 0x2d,
	0xd6, 0xcb, 0xcc, 0x4b, 0x2e, 0x4d, 0x4a, 0x2c, 0xc9, 0x2f, 0xd2, 0x4b, 0x2f, 0xca, 0x2f, 0x2d,
	0xd0, 0x2b, 0x33, 0x8c, 0x4f, 0xcc, 0x29, 0xc8, 0x48, 0xd4, 0x83, 0xe9, 0x96, 0xd2, 0x2e, 0xc9,
	0xc8, 0x2c, 0x4a, 0x89, 0x2f, 0x48, 0x2c, 0x2a, 0xa9, 0xd4, 0x07, 0x1b, 0xa2, 0x0f, 0x31, 0x43,
	0x17, 0x99, 0x03, 0x31, 0x5e, 0x4a, 0x24, 0x3d, 0x3f, 0x3d, 0x1f, 0x22, 0x0e, 0x62, 0x41, 0x45,
	0x05, 0xc1, 0x66, 0x23, 0xbb, 0x43, 0xe9, 0x0b, 0x23, 0x17, 0xaf, 0x6f, 0xa5, 0x63, 0x41, 0x41,
	0x40, 0x51, 0x7e, 0x41, 0x7e, 0x71, 0x62, 0x8e, 0x50, 0x00, 0x17, 0xa3, 0xa3, 0x04, 0xa3, 0x02,
	0xa3, 0x06, 0xb7, 0x91, 0x83, 0x1e, 0xa9, 0xae, 0xd4, 0x73, 0x44, 0x31, 0xcc, 0x83, 0x21, 0x88,
	0xd1, 0x11, 0x64, 0xa2, 0x93, 0x04, 0x13, 0xb9, 0x26, 0x3a, 0x61, 0x98, 0xe8, 0x64, 0x65, 0x71,
	0x6a, 0x8b, 0xae, 0x89, 0x56, 0x7a, 0x66, 0x49, 0x46, 0x69, 0x92, 0x5e, 0x72, 0x7e, 0x2e, 0xd4,
	0xf3, 0xfa, 0x50, 0x53, 0xf5, 0xe1, 0xa6, 0xea, 0x43, 0x4c, 0x85, 0xe9, 0xf6, 0x74, 0x62, 0xe5,
	0x62, 0x2e, 0x2e, 0xcd, 0x55, 0x4a, 0xe4, 0xe2, 0x43, 0x75, 0xa9, 0x90, 0x3f, 0x17, 0x4b, 0x52,
	0x62, 0x71, 0x2a, 0xd4, 0xe7, 0xa6, 0x44, 0xbb, 0x13, 0x66, 0x80, 0x53, 0x62, 0x71, 0xaa, 0x13,
	0xcb, 0x89, 0x7b, 0xf2, 0x0c, 0x41, 0x60, 0x83, 0x40, 0x56, 0x38, 0xd1, 0xd8, 0x8a, 0x04, 0x2e,
	0x6e, 0xdf, 0xe2, 0x74, 0x88, 0x74, 0xaa, 0xa3, 0x50, 0x20, 0x8a, 0xf9, 0xe6, 0x44, 0x9b, 0x8f,
	0x30, 0x03, 0xbf, 0x0d, 0x4e, 0x34, 0xb0, 0xc1, 0xc9, 0xe7, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f,
	0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0x9c, 0xf0, 0x58, 0x8e, 0xe1, 0xc2, 0x63, 0x39, 0x86, 0x1b,
	0x8f, 0xe5, 0x18, 0xa2, 0x8c, 0x88, 0x8e, 0x5f, 0x7d, 0x58, 0x62, 0x49, 0x62, 0x03, 0xa7, 0x6a,
	0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xd0, 0xdb, 0x8f, 0x0f, 0x75, 0x03, 0x00, 0x00,
}

func (this *MyAppProposal) GetProposalI() github_com_cosmos_modules_incubator_group.ProposalI {
	if x := this.GetA(); x != nil {
		return x
	}
	if x := this.GetB(); x != nil {
		return x
	}
	return nil
}

func (this *MyAppProposal) SetProposalI(value github_com_cosmos_modules_incubator_group.ProposalI) error {
	if value == nil {
		this.Sum = nil
		return nil
	}
	switch vt := value.(type) {
	case *AMyAppProposal:
		this.Sum = &MyAppProposal_A{vt}
		return nil
	case *BMyAppProposal:
		this.Sum = &MyAppProposal_B{vt}
		return nil
	}
	return fmt.Errorf("can't encode value of type %T as message MyAppProposal", value)
}

func (m *MyAppProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MyAppProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MyAppProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Sum != nil {
		{
			size := m.Sum.Size()
			i -= size
			if _, err := m.Sum.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	return len(dAtA) - i, nil
}

func (m *MyAppProposal_A) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MyAppProposal_A) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.A != nil {
		{
			size, err := m.A.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTypes(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}
func (m *MyAppProposal_B) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MyAppProposal_B) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.B != nil {
		{
			size, err := m.B.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTypes(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	return len(dAtA) - i, nil
}
func (m *AMyAppProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AMyAppProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AMyAppProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Base.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *BMyAppProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BMyAppProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BMyAppProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Base.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *MsgProposeA) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgProposeA) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgProposeA) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Base.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *MsgProposeB) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgProposeB) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgProposeB) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Base.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MyAppProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Sum != nil {
		n += m.Sum.Size()
	}
	return n
}

func (m *MyAppProposal_A) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.A != nil {
		l = m.A.Size()
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}
func (m *MyAppProposal_B) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.B != nil {
		l = m.B.Size()
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}
func (m *AMyAppProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Base.Size()
	n += 1 + l + sovTypes(uint64(l))
	return n
}

func (m *BMyAppProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Base.Size()
	n += 1 + l + sovTypes(uint64(l))
	return n
}

func (m *MsgProposeA) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Base.Size()
	n += 1 + l + sovTypes(uint64(l))
	return n
}

func (m *MsgProposeB) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Base.Size()
	n += 1 + l + sovTypes(uint64(l))
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MyAppProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: MyAppProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MyAppProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field A", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &AMyAppProposal{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Sum = &MyAppProposal_A{v}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field B", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &BMyAppProposal{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Sum = &MyAppProposal_B{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *AMyAppProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: AMyAppProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AMyAppProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Base", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Base.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *BMyAppProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: BMyAppProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BMyAppProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Base", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Base.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgProposeA) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: MsgProposeA: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgProposeA: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Base", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Base.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgProposeB) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: MsgProposeB: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgProposeB: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Base", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Base.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTypes
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)
