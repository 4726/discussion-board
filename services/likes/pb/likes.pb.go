// Code generated by protoc-gen-go. DO NOT EDIT.
// source: likes.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type HealthCheckResponse_ServingStatus int32

const (
	HealthCheckResponse_UNKNOWN     HealthCheckResponse_ServingStatus = 0
	HealthCheckResponse_SERVING     HealthCheckResponse_ServingStatus = 1
	HealthCheckResponse_NOT_SERVING HealthCheckResponse_ServingStatus = 2
)

var HealthCheckResponse_ServingStatus_name = map[int32]string{
	0: "UNKNOWN",
	1: "SERVING",
	2: "NOT_SERVING",
}

var HealthCheckResponse_ServingStatus_value = map[string]int32{
	"UNKNOWN":     0,
	"SERVING":     1,
	"NOT_SERVING": 2,
}

func (x HealthCheckResponse_ServingStatus) Enum() *HealthCheckResponse_ServingStatus {
	p := new(HealthCheckResponse_ServingStatus)
	*p = x
	return p
}

func (x HealthCheckResponse_ServingStatus) String() string {
	return proto.EnumName(HealthCheckResponse_ServingStatus_name, int32(x))
}

func (x *HealthCheckResponse_ServingStatus) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(HealthCheckResponse_ServingStatus_value, data, "HealthCheckResponse_ServingStatus")
	if err != nil {
		return err
	}
	*x = HealthCheckResponse_ServingStatus(value)
	return nil
}

func (HealthCheckResponse_ServingStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_cff81f36f81c8d8e, []int{9, 0}
}

type IDUserID struct {
	Id                   *uint64  `protobuf:"varint,1,req,name=id" json:"id,omitempty"`
	UserId               *uint64  `protobuf:"varint,2,req,name=user_id,json=userId" json:"user_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IDUserID) Reset()         { *m = IDUserID{} }
func (m *IDUserID) String() string { return proto.CompactTextString(m) }
func (*IDUserID) ProtoMessage()    {}
func (*IDUserID) Descriptor() ([]byte, []int) {
	return fileDescriptor_cff81f36f81c8d8e, []int{0}
}

func (m *IDUserID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IDUserID.Unmarshal(m, b)
}
func (m *IDUserID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IDUserID.Marshal(b, m, deterministic)
}
func (m *IDUserID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IDUserID.Merge(m, src)
}
func (m *IDUserID) XXX_Size() int {
	return xxx_messageInfo_IDUserID.Size(m)
}
func (m *IDUserID) XXX_DiscardUnknown() {
	xxx_messageInfo_IDUserID.DiscardUnknown(m)
}

var xxx_messageInfo_IDUserID proto.InternalMessageInfo

func (m *IDUserID) GetId() uint64 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *IDUserID) GetUserId() uint64 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

type Total struct {
	Total                *uint64  `protobuf:"varint,1,req,name=total" json:"total,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Total) Reset()         { *m = Total{} }
func (m *Total) String() string { return proto.CompactTextString(m) }
func (*Total) ProtoMessage()    {}
func (*Total) Descriptor() ([]byte, []int) {
	return fileDescriptor_cff81f36f81c8d8e, []int{1}
}

func (m *Total) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Total.Unmarshal(m, b)
}
func (m *Total) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Total.Marshal(b, m, deterministic)
}
func (m *Total) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Total.Merge(m, src)
}
func (m *Total) XXX_Size() int {
	return xxx_messageInfo_Total.Size(m)
}
func (m *Total) XXX_DiscardUnknown() {
	xxx_messageInfo_Total.DiscardUnknown(m)
}

var xxx_messageInfo_Total proto.InternalMessageInfo

func (m *Total) GetTotal() uint64 {
	if m != nil && m.Total != nil {
		return *m.Total
	}
	return 0
}

type IDs struct {
	Id                   []uint64 `protobuf:"varint,1,rep,packed,name=id" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IDs) Reset()         { *m = IDs{} }
func (m *IDs) String() string { return proto.CompactTextString(m) }
func (*IDs) ProtoMessage()    {}
func (*IDs) Descriptor() ([]byte, []int) {
	return fileDescriptor_cff81f36f81c8d8e, []int{2}
}

func (m *IDs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IDs.Unmarshal(m, b)
}
func (m *IDs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IDs.Marshal(b, m, deterministic)
}
func (m *IDs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IDs.Merge(m, src)
}
func (m *IDs) XXX_Size() int {
	return xxx_messageInfo_IDs.Size(m)
}
func (m *IDs) XXX_DiscardUnknown() {
	xxx_messageInfo_IDs.DiscardUnknown(m)
}

var xxx_messageInfo_IDs proto.InternalMessageInfo

func (m *IDs) GetId() []uint64 {
	if m != nil {
		return m.Id
	}
	return nil
}

type TotalLikes struct {
	IdLikes              []*TotalLikes_IDLikes `protobuf:"bytes,1,rep,name=id_likes,json=idLikes" json:"id_likes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *TotalLikes) Reset()         { *m = TotalLikes{} }
func (m *TotalLikes) String() string { return proto.CompactTextString(m) }
func (*TotalLikes) ProtoMessage()    {}
func (*TotalLikes) Descriptor() ([]byte, []int) {
	return fileDescriptor_cff81f36f81c8d8e, []int{3}
}

func (m *TotalLikes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TotalLikes.Unmarshal(m, b)
}
func (m *TotalLikes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TotalLikes.Marshal(b, m, deterministic)
}
func (m *TotalLikes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TotalLikes.Merge(m, src)
}
func (m *TotalLikes) XXX_Size() int {
	return xxx_messageInfo_TotalLikes.Size(m)
}
func (m *TotalLikes) XXX_DiscardUnknown() {
	xxx_messageInfo_TotalLikes.DiscardUnknown(m)
}

var xxx_messageInfo_TotalLikes proto.InternalMessageInfo

func (m *TotalLikes) GetIdLikes() []*TotalLikes_IDLikes {
	if m != nil {
		return m.IdLikes
	}
	return nil
}

type TotalLikes_IDLikes struct {
	Id                   *uint64  `protobuf:"varint,1,req,name=id" json:"id,omitempty"`
	Total                *uint64  `protobuf:"varint,2,req,name=total" json:"total,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TotalLikes_IDLikes) Reset()         { *m = TotalLikes_IDLikes{} }
func (m *TotalLikes_IDLikes) String() string { return proto.CompactTextString(m) }
func (*TotalLikes_IDLikes) ProtoMessage()    {}
func (*TotalLikes_IDLikes) Descriptor() ([]byte, []int) {
	return fileDescriptor_cff81f36f81c8d8e, []int{3, 0}
}

func (m *TotalLikes_IDLikes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TotalLikes_IDLikes.Unmarshal(m, b)
}
func (m *TotalLikes_IDLikes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TotalLikes_IDLikes.Marshal(b, m, deterministic)
}
func (m *TotalLikes_IDLikes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TotalLikes_IDLikes.Merge(m, src)
}
func (m *TotalLikes_IDLikes) XXX_Size() int {
	return xxx_messageInfo_TotalLikes_IDLikes.Size(m)
}
func (m *TotalLikes_IDLikes) XXX_DiscardUnknown() {
	xxx_messageInfo_TotalLikes_IDLikes.DiscardUnknown(m)
}

var xxx_messageInfo_TotalLikes_IDLikes proto.InternalMessageInfo

func (m *TotalLikes_IDLikes) GetId() uint64 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *TotalLikes_IDLikes) GetTotal() uint64 {
	if m != nil && m.Total != nil {
		return *m.Total
	}
	return 0
}

type IDsUserID struct {
	Id                   []uint64 `protobuf:"varint,1,rep,name=id" json:"id,omitempty"`
	UserId               *uint64  `protobuf:"varint,2,req,name=user_id,json=userId" json:"user_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IDsUserID) Reset()         { *m = IDsUserID{} }
func (m *IDsUserID) String() string { return proto.CompactTextString(m) }
func (*IDsUserID) ProtoMessage()    {}
func (*IDsUserID) Descriptor() ([]byte, []int) {
	return fileDescriptor_cff81f36f81c8d8e, []int{4}
}

func (m *IDsUserID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IDsUserID.Unmarshal(m, b)
}
func (m *IDsUserID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IDsUserID.Marshal(b, m, deterministic)
}
func (m *IDsUserID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IDsUserID.Merge(m, src)
}
func (m *IDsUserID) XXX_Size() int {
	return xxx_messageInfo_IDsUserID.Size(m)
}
func (m *IDsUserID) XXX_DiscardUnknown() {
	xxx_messageInfo_IDsUserID.DiscardUnknown(m)
}

var xxx_messageInfo_IDsUserID proto.InternalMessageInfo

func (m *IDsUserID) GetId() []uint64 {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *IDsUserID) GetUserId() uint64 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

type HaveLikes struct {
	HaveLikes            []*HaveLikes_HaveLike `protobuf:"bytes,1,rep,name=have_likes,json=haveLikes" json:"have_likes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *HaveLikes) Reset()         { *m = HaveLikes{} }
func (m *HaveLikes) String() string { return proto.CompactTextString(m) }
func (*HaveLikes) ProtoMessage()    {}
func (*HaveLikes) Descriptor() ([]byte, []int) {
	return fileDescriptor_cff81f36f81c8d8e, []int{5}
}

func (m *HaveLikes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HaveLikes.Unmarshal(m, b)
}
func (m *HaveLikes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HaveLikes.Marshal(b, m, deterministic)
}
func (m *HaveLikes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HaveLikes.Merge(m, src)
}
func (m *HaveLikes) XXX_Size() int {
	return xxx_messageInfo_HaveLikes.Size(m)
}
func (m *HaveLikes) XXX_DiscardUnknown() {
	xxx_messageInfo_HaveLikes.DiscardUnknown(m)
}

var xxx_messageInfo_HaveLikes proto.InternalMessageInfo

func (m *HaveLikes) GetHaveLikes() []*HaveLikes_HaveLike {
	if m != nil {
		return m.HaveLikes
	}
	return nil
}

type HaveLikes_HaveLike struct {
	Id                   *uint64  `protobuf:"varint,1,req,name=id" json:"id,omitempty"`
	HasLike              *bool    `protobuf:"varint,2,req,name=has_like,json=hasLike" json:"has_like,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HaveLikes_HaveLike) Reset()         { *m = HaveLikes_HaveLike{} }
func (m *HaveLikes_HaveLike) String() string { return proto.CompactTextString(m) }
func (*HaveLikes_HaveLike) ProtoMessage()    {}
func (*HaveLikes_HaveLike) Descriptor() ([]byte, []int) {
	return fileDescriptor_cff81f36f81c8d8e, []int{5, 0}
}

func (m *HaveLikes_HaveLike) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HaveLikes_HaveLike.Unmarshal(m, b)
}
func (m *HaveLikes_HaveLike) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HaveLikes_HaveLike.Marshal(b, m, deterministic)
}
func (m *HaveLikes_HaveLike) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HaveLikes_HaveLike.Merge(m, src)
}
func (m *HaveLikes_HaveLike) XXX_Size() int {
	return xxx_messageInfo_HaveLikes_HaveLike.Size(m)
}
func (m *HaveLikes_HaveLike) XXX_DiscardUnknown() {
	xxx_messageInfo_HaveLikes_HaveLike.DiscardUnknown(m)
}

var xxx_messageInfo_HaveLikes_HaveLike proto.InternalMessageInfo

func (m *HaveLikes_HaveLike) GetId() uint64 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *HaveLikes_HaveLike) GetHasLike() bool {
	if m != nil && m.HasLike != nil {
		return *m.HasLike
	}
	return false
}

type Id struct {
	Id                   *uint64  `protobuf:"varint,1,req,name=id" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Id) Reset()         { *m = Id{} }
func (m *Id) String() string { return proto.CompactTextString(m) }
func (*Id) ProtoMessage()    {}
func (*Id) Descriptor() ([]byte, []int) {
	return fileDescriptor_cff81f36f81c8d8e, []int{6}
}

func (m *Id) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Id.Unmarshal(m, b)
}
func (m *Id) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Id.Marshal(b, m, deterministic)
}
func (m *Id) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Id.Merge(m, src)
}
func (m *Id) XXX_Size() int {
	return xxx_messageInfo_Id.Size(m)
}
func (m *Id) XXX_DiscardUnknown() {
	xxx_messageInfo_Id.DiscardUnknown(m)
}

var xxx_messageInfo_Id proto.InternalMessageInfo

func (m *Id) GetId() uint64 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

type DeletePostResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeletePostResponse) Reset()         { *m = DeletePostResponse{} }
func (m *DeletePostResponse) String() string { return proto.CompactTextString(m) }
func (*DeletePostResponse) ProtoMessage()    {}
func (*DeletePostResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cff81f36f81c8d8e, []int{7}
}

func (m *DeletePostResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeletePostResponse.Unmarshal(m, b)
}
func (m *DeletePostResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeletePostResponse.Marshal(b, m, deterministic)
}
func (m *DeletePostResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeletePostResponse.Merge(m, src)
}
func (m *DeletePostResponse) XXX_Size() int {
	return xxx_messageInfo_DeletePostResponse.Size(m)
}
func (m *DeletePostResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeletePostResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeletePostResponse proto.InternalMessageInfo

type HealthCheckRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HealthCheckRequest) Reset()         { *m = HealthCheckRequest{} }
func (m *HealthCheckRequest) String() string { return proto.CompactTextString(m) }
func (*HealthCheckRequest) ProtoMessage()    {}
func (*HealthCheckRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cff81f36f81c8d8e, []int{8}
}

func (m *HealthCheckRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HealthCheckRequest.Unmarshal(m, b)
}
func (m *HealthCheckRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HealthCheckRequest.Marshal(b, m, deterministic)
}
func (m *HealthCheckRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HealthCheckRequest.Merge(m, src)
}
func (m *HealthCheckRequest) XXX_Size() int {
	return xxx_messageInfo_HealthCheckRequest.Size(m)
}
func (m *HealthCheckRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HealthCheckRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HealthCheckRequest proto.InternalMessageInfo

type HealthCheckResponse struct {
	Status               *HealthCheckResponse_ServingStatus `protobuf:"varint,1,req,name=status,enum=pb.HealthCheckResponse_ServingStatus" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                           `json:"-"`
	XXX_unrecognized     []byte                             `json:"-"`
	XXX_sizecache        int32                              `json:"-"`
}

func (m *HealthCheckResponse) Reset()         { *m = HealthCheckResponse{} }
func (m *HealthCheckResponse) String() string { return proto.CompactTextString(m) }
func (*HealthCheckResponse) ProtoMessage()    {}
func (*HealthCheckResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cff81f36f81c8d8e, []int{9}
}

func (m *HealthCheckResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HealthCheckResponse.Unmarshal(m, b)
}
func (m *HealthCheckResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HealthCheckResponse.Marshal(b, m, deterministic)
}
func (m *HealthCheckResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HealthCheckResponse.Merge(m, src)
}
func (m *HealthCheckResponse) XXX_Size() int {
	return xxx_messageInfo_HealthCheckResponse.Size(m)
}
func (m *HealthCheckResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HealthCheckResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HealthCheckResponse proto.InternalMessageInfo

func (m *HealthCheckResponse) GetStatus() HealthCheckResponse_ServingStatus {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return HealthCheckResponse_UNKNOWN
}

func init() {
	proto.RegisterEnum("pb.HealthCheckResponse_ServingStatus", HealthCheckResponse_ServingStatus_name, HealthCheckResponse_ServingStatus_value)
	proto.RegisterType((*IDUserID)(nil), "pb.IDUserID")
	proto.RegisterType((*Total)(nil), "pb.Total")
	proto.RegisterType((*IDs)(nil), "pb.IDs")
	proto.RegisterType((*TotalLikes)(nil), "pb.TotalLikes")
	proto.RegisterType((*TotalLikes_IDLikes)(nil), "pb.TotalLikes.IDLikes")
	proto.RegisterType((*IDsUserID)(nil), "pb.IDsUserID")
	proto.RegisterType((*HaveLikes)(nil), "pb.HaveLikes")
	proto.RegisterType((*HaveLikes_HaveLike)(nil), "pb.HaveLikes.HaveLike")
	proto.RegisterType((*Id)(nil), "pb.Id")
	proto.RegisterType((*DeletePostResponse)(nil), "pb.DeletePostResponse")
	proto.RegisterType((*HealthCheckRequest)(nil), "pb.HealthCheckRequest")
	proto.RegisterType((*HealthCheckResponse)(nil), "pb.HealthCheckResponse")
}

func init() { proto.RegisterFile("likes.proto", fileDescriptor_cff81f36f81c8d8e) }

var fileDescriptor_cff81f36f81c8d8e = []byte{
	// 491 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x6f, 0x6b, 0xd3, 0x50,
	0x14, 0xc6, 0x6d, 0xba, 0x36, 0xe9, 0xd3, 0xa5, 0x2b, 0xc7, 0xb2, 0x3f, 0x01, 0x41, 0x22, 0xd3,
	0x22, 0x1a, 0xb5, 0x3a, 0x10, 0xc1, 0x37, 0x2e, 0xd2, 0x05, 0x25, 0x93, 0x74, 0xd5, 0x97, 0x25,
	0x33, 0x17, 0x13, 0x96, 0x35, 0xb1, 0xf7, 0xb6, 0xe0, 0xb7, 0xf0, 0x3b, 0xf9, 0xc5, 0xe4, 0xde,
	0xa4, 0x99, 0x5d, 0xd0, 0xfa, 0xaa, 0xf7, 0x9e, 0xfe, 0xce, 0x73, 0x9e, 0x93, 0xfb, 0xa0, 0x9b,
	0x26, 0x57, 0x8c, 0x3b, 0xf9, 0x22, 0x13, 0x19, 0x69, 0xf9, 0xa5, 0xfd, 0x12, 0x86, 0xe7, 0x4e,
	0x39, 0x5b, 0x78, 0x2e, 0xf5, 0xa0, 0x25, 0xd1, 0x61, 0xe3, 0xbe, 0x36, 0xdc, 0x09, 0xb4, 0x24,
	0xa2, 0x03, 0xe8, 0x4b, 0xce, 0x16, 0xb3, 0x24, 0x3a, 0xd4, 0x54, 0xb1, 0x2d, 0xaf, 0x5e, 0x64,
	0xdf, 0x43, 0xeb, 0x22, 0x13, 0x61, 0x4a, 0x03, 0xb4, 0x84, 0x3c, 0x94, 0x4d, 0xc5, 0xc5, 0x3e,
	0x42, 0xd3, 0x73, 0x39, 0x51, 0x29, 0xd7, 0x1c, 0xee, 0xbc, 0xd3, 0xfa, 0x0d, 0x29, 0x69, 0xe7,
	0x80, 0xea, 0xfc, 0x28, 0x6d, 0xd0, 0x0b, 0x18, 0x49, 0x34, 0x53, 0x96, 0x14, 0xd7, 0x1d, 0xed,
	0x3b, 0xf9, 0xa5, 0x73, 0x43, 0x38, 0x9e, 0xab, 0x7e, 0x03, 0x3d, 0x89, 0xd4, 0xc1, 0x7a, 0x06,
	0xbd, 0xac, 0xd5, 0xec, 0x56, 0x66, 0xb4, 0x3f, 0xcd, 0xbc, 0x42, 0xc7, 0x73, 0xf9, 0xad, 0x0d,
	0x9b, 0xdb, 0x36, 0xfc, 0x81, 0xce, 0x59, 0xb8, 0x62, 0xc5, 0xa0, 0x13, 0x20, 0x0e, 0x57, 0xac,
	0x6e, 0xb4, 0x42, 0xaa, 0x53, 0xd0, 0x89, 0xd7, 0x35, 0xeb, 0x04, 0xc6, 0xba, 0x5c, 0xf3, 0x7a,
	0x04, 0x23, 0x0e, 0xb9, 0x52, 0x54, 0x93, 0x8d, 0x40, 0x8f, 0x43, 0x2e, 0x51, 0x7b, 0x00, 0xcd,
	0x8b, 0x6e, 0x37, 0xd8, 0x03, 0x90, 0xcb, 0x52, 0x26, 0xd8, 0xa7, 0x8c, 0x8b, 0x80, 0xf1, 0x3c,
	0x9b, 0x73, 0xc9, 0xd2, 0x19, 0x0b, 0x53, 0x11, 0x9f, 0xc6, 0xec, 0xeb, 0x55, 0xc0, 0xbe, 0x2f,
	0x19, 0x17, 0xf6, 0xcf, 0x06, 0xee, 0x6e, 0x94, 0x0b, 0x9a, 0xde, 0xa2, 0xcd, 0x45, 0x28, 0x96,
	0x5c, 0xe9, 0xf6, 0x46, 0xc7, 0x6a, 0x87, 0x3a, 0xe8, 0x4c, 0xd8, 0x62, 0x95, 0xcc, 0xbf, 0x4d,
	0x14, 0x1c, 0x94, 0x4d, 0xf6, 0x1b, 0x98, 0x1b, 0x7f, 0x50, 0x17, 0xfa, 0xd4, 0xff, 0xe0, 0x9f,
	0x7f, 0xf1, 0xfb, 0x77, 0xe4, 0x65, 0xf2, 0x3e, 0xf8, 0xec, 0xf9, 0xe3, 0x7e, 0x83, 0xf6, 0xd0,
	0xf5, 0xcf, 0x2f, 0x66, 0xeb, 0x82, 0x36, 0xfa, 0xd5, 0x44, 0xab, 0xf8, 0x98, 0x0f, 0x60, 0xc8,
	0x83, 0x5c, 0x83, 0x76, 0xa5, 0x81, 0x75, 0xfc, 0xac, 0x4e, 0xf5, 0xf6, 0x74, 0x0c, 0x4c, 0xe7,
	0xe9, 0x56, 0xec, 0x21, 0xba, 0x52, 0xeb, 0x34, 0xbb, 0xbe, 0x66, 0xf3, 0x7f, 0x70, 0x43, 0x98,
	0x85, 0xdc, 0x56, 0xf2, 0x11, 0x76, 0xc7, 0x4c, 0xc8, 0xa9, 0x85, 0x5b, 0xbd, 0x00, 0xb9, 0xd5,
	0xdb, 0x0c, 0x26, 0x3d, 0xc6, 0xde, 0x98, 0x89, 0x52, 0x6f, 0x0b, 0xfb, 0x14, 0xa6, 0x54, 0xe4,
	0x55, 0x1a, 0xcc, 0x92, 0x2c, 0xe7, 0x9b, 0x1b, 0x59, 0xa2, 0xe7, 0xe8, 0x97, 0xba, 0xff, 0xdb,
	0xf1, 0x04, 0xb8, 0x09, 0x07, 0xb5, 0x15, 0x1b, 0x59, 0x2a, 0xa2, 0xf5, 0xd0, 0xd0, 0x6b, 0xb4,
	0xd4, 0x73, 0xd3, 0x7e, 0xed, 0xfd, 0x55, 0x7e, 0xac, 0x83, 0xbf, 0xe4, 0xe2, 0x77, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x1f, 0xd9, 0x29, 0xf1, 0x3e, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// LikesClient is the client API for Likes service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LikesClient interface {
	LikePost(ctx context.Context, in *IDUserID, opts ...grpc.CallOption) (*Total, error)
	UnlikePost(ctx context.Context, in *IDUserID, opts ...grpc.CallOption) (*Total, error)
	LikeComment(ctx context.Context, in *IDUserID, opts ...grpc.CallOption) (*Total, error)
	UnlikeComment(ctx context.Context, in *IDUserID, opts ...grpc.CallOption) (*Total, error)
	GetPostLikes(ctx context.Context, in *IDs, opts ...grpc.CallOption) (*TotalLikes, error)
	GetCommentLikes(ctx context.Context, in *IDs, opts ...grpc.CallOption) (*TotalLikes, error)
	PostsHaveLike(ctx context.Context, in *IDsUserID, opts ...grpc.CallOption) (*HaveLikes, error)
	CommentsHaveLike(ctx context.Context, in *IDsUserID, opts ...grpc.CallOption) (*HaveLikes, error)
	DeletePost(ctx context.Context, in *Id, opts ...grpc.CallOption) (*DeletePostResponse, error)
	Check(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error)
}

type likesClient struct {
	cc *grpc.ClientConn
}

func NewLikesClient(cc *grpc.ClientConn) LikesClient {
	return &likesClient{cc}
}

func (c *likesClient) LikePost(ctx context.Context, in *IDUserID, opts ...grpc.CallOption) (*Total, error) {
	out := new(Total)
	err := c.cc.Invoke(ctx, "/pb.Likes/LikePost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likesClient) UnlikePost(ctx context.Context, in *IDUserID, opts ...grpc.CallOption) (*Total, error) {
	out := new(Total)
	err := c.cc.Invoke(ctx, "/pb.Likes/UnlikePost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likesClient) LikeComment(ctx context.Context, in *IDUserID, opts ...grpc.CallOption) (*Total, error) {
	out := new(Total)
	err := c.cc.Invoke(ctx, "/pb.Likes/LikeComment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likesClient) UnlikeComment(ctx context.Context, in *IDUserID, opts ...grpc.CallOption) (*Total, error) {
	out := new(Total)
	err := c.cc.Invoke(ctx, "/pb.Likes/UnlikeComment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likesClient) GetPostLikes(ctx context.Context, in *IDs, opts ...grpc.CallOption) (*TotalLikes, error) {
	out := new(TotalLikes)
	err := c.cc.Invoke(ctx, "/pb.Likes/GetPostLikes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likesClient) GetCommentLikes(ctx context.Context, in *IDs, opts ...grpc.CallOption) (*TotalLikes, error) {
	out := new(TotalLikes)
	err := c.cc.Invoke(ctx, "/pb.Likes/GetCommentLikes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likesClient) PostsHaveLike(ctx context.Context, in *IDsUserID, opts ...grpc.CallOption) (*HaveLikes, error) {
	out := new(HaveLikes)
	err := c.cc.Invoke(ctx, "/pb.Likes/PostsHaveLike", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likesClient) CommentsHaveLike(ctx context.Context, in *IDsUserID, opts ...grpc.CallOption) (*HaveLikes, error) {
	out := new(HaveLikes)
	err := c.cc.Invoke(ctx, "/pb.Likes/CommentsHaveLike", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likesClient) DeletePost(ctx context.Context, in *Id, opts ...grpc.CallOption) (*DeletePostResponse, error) {
	out := new(DeletePostResponse)
	err := c.cc.Invoke(ctx, "/pb.Likes/DeletePost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *likesClient) Check(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error) {
	out := new(HealthCheckResponse)
	err := c.cc.Invoke(ctx, "/pb.Likes/Check", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LikesServer is the server API for Likes service.
type LikesServer interface {
	LikePost(context.Context, *IDUserID) (*Total, error)
	UnlikePost(context.Context, *IDUserID) (*Total, error)
	LikeComment(context.Context, *IDUserID) (*Total, error)
	UnlikeComment(context.Context, *IDUserID) (*Total, error)
	GetPostLikes(context.Context, *IDs) (*TotalLikes, error)
	GetCommentLikes(context.Context, *IDs) (*TotalLikes, error)
	PostsHaveLike(context.Context, *IDsUserID) (*HaveLikes, error)
	CommentsHaveLike(context.Context, *IDsUserID) (*HaveLikes, error)
	DeletePost(context.Context, *Id) (*DeletePostResponse, error)
	Check(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error)
}

// UnimplementedLikesServer can be embedded to have forward compatible implementations.
type UnimplementedLikesServer struct {
}

func (*UnimplementedLikesServer) LikePost(ctx context.Context, req *IDUserID) (*Total, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LikePost not implemented")
}
func (*UnimplementedLikesServer) UnlikePost(ctx context.Context, req *IDUserID) (*Total, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnlikePost not implemented")
}
func (*UnimplementedLikesServer) LikeComment(ctx context.Context, req *IDUserID) (*Total, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LikeComment not implemented")
}
func (*UnimplementedLikesServer) UnlikeComment(ctx context.Context, req *IDUserID) (*Total, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnlikeComment not implemented")
}
func (*UnimplementedLikesServer) GetPostLikes(ctx context.Context, req *IDs) (*TotalLikes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPostLikes not implemented")
}
func (*UnimplementedLikesServer) GetCommentLikes(ctx context.Context, req *IDs) (*TotalLikes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCommentLikes not implemented")
}
func (*UnimplementedLikesServer) PostsHaveLike(ctx context.Context, req *IDsUserID) (*HaveLikes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostsHaveLike not implemented")
}
func (*UnimplementedLikesServer) CommentsHaveLike(ctx context.Context, req *IDsUserID) (*HaveLikes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommentsHaveLike not implemented")
}
func (*UnimplementedLikesServer) DeletePost(ctx context.Context, req *Id) (*DeletePostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePost not implemented")
}
func (*UnimplementedLikesServer) Check(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}

func RegisterLikesServer(s *grpc.Server, srv LikesServer) {
	s.RegisterService(&_Likes_serviceDesc, srv)
}

func _Likes_LikePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDUserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikesServer).LikePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Likes/LikePost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikesServer).LikePost(ctx, req.(*IDUserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Likes_UnlikePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDUserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikesServer).UnlikePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Likes/UnlikePost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikesServer).UnlikePost(ctx, req.(*IDUserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Likes_LikeComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDUserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikesServer).LikeComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Likes/LikeComment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikesServer).LikeComment(ctx, req.(*IDUserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Likes_UnlikeComment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDUserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikesServer).UnlikeComment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Likes/UnlikeComment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikesServer).UnlikeComment(ctx, req.(*IDUserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Likes_GetPostLikes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikesServer).GetPostLikes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Likes/GetPostLikes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikesServer).GetPostLikes(ctx, req.(*IDs))
	}
	return interceptor(ctx, in, info, handler)
}

func _Likes_GetCommentLikes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikesServer).GetCommentLikes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Likes/GetCommentLikes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikesServer).GetCommentLikes(ctx, req.(*IDs))
	}
	return interceptor(ctx, in, info, handler)
}

func _Likes_PostsHaveLike_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDsUserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikesServer).PostsHaveLike(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Likes/PostsHaveLike",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikesServer).PostsHaveLike(ctx, req.(*IDsUserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Likes_CommentsHaveLike_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDsUserID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikesServer).CommentsHaveLike(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Likes/CommentsHaveLike",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikesServer).CommentsHaveLike(ctx, req.(*IDsUserID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Likes_DeletePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Id)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikesServer).DeletePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Likes/DeletePost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikesServer).DeletePost(ctx, req.(*Id))
	}
	return interceptor(ctx, in, info, handler)
}

func _Likes_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LikesServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Likes/Check",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LikesServer).Check(ctx, req.(*HealthCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Likes_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Likes",
	HandlerType: (*LikesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LikePost",
			Handler:    _Likes_LikePost_Handler,
		},
		{
			MethodName: "UnlikePost",
			Handler:    _Likes_UnlikePost_Handler,
		},
		{
			MethodName: "LikeComment",
			Handler:    _Likes_LikeComment_Handler,
		},
		{
			MethodName: "UnlikeComment",
			Handler:    _Likes_UnlikeComment_Handler,
		},
		{
			MethodName: "GetPostLikes",
			Handler:    _Likes_GetPostLikes_Handler,
		},
		{
			MethodName: "GetCommentLikes",
			Handler:    _Likes_GetCommentLikes_Handler,
		},
		{
			MethodName: "PostsHaveLike",
			Handler:    _Likes_PostsHaveLike_Handler,
		},
		{
			MethodName: "CommentsHaveLike",
			Handler:    _Likes_CommentsHaveLike_Handler,
		},
		{
			MethodName: "DeletePost",
			Handler:    _Likes_DeletePost_Handler,
		},
		{
			MethodName: "Check",
			Handler:    _Likes_Check_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "likes.proto",
}
