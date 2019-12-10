// Code generated by protoc-gen-go. DO NOT EDIT.
// source: search.proto

package pb

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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
	return fileDescriptor_453745cff914010e, []int{11, 0}
}

type Post struct {
	Title                *string  `protobuf:"bytes,1,req,name=title" json:"title,omitempty"`
	Body                 *string  `protobuf:"bytes,2,req,name=body" json:"body,omitempty"`
	UserId               *uint64  `protobuf:"varint,3,req,name=user_id,json=userId" json:"user_id,omitempty"`
	Id                   *uint64  `protobuf:"varint,4,req,name=id" json:"id,omitempty"`
	Timestamp            *int64   `protobuf:"varint,5,opt,name=timestamp" json:"timestamp,omitempty"`
	Likes                *int64   `protobuf:"varint,6,opt,name=likes" json:"likes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Post) Reset()         { *m = Post{} }
func (m *Post) String() string { return proto.CompactTextString(m) }
func (*Post) ProtoMessage()    {}
func (*Post) Descriptor() ([]byte, []int) {
	return fileDescriptor_453745cff914010e, []int{0}
}

func (m *Post) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Post.Unmarshal(m, b)
}
func (m *Post) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Post.Marshal(b, m, deterministic)
}
func (m *Post) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Post.Merge(m, src)
}
func (m *Post) XXX_Size() int {
	return xxx_messageInfo_Post.Size(m)
}
func (m *Post) XXX_DiscardUnknown() {
	xxx_messageInfo_Post.DiscardUnknown(m)
}

var xxx_messageInfo_Post proto.InternalMessageInfo

func (m *Post) GetTitle() string {
	if m != nil && m.Title != nil {
		return *m.Title
	}
	return ""
}

func (m *Post) GetBody() string {
	if m != nil && m.Body != nil {
		return *m.Body
	}
	return ""
}

func (m *Post) GetUserId() uint64 {
	if m != nil && m.UserId != nil {
		return *m.UserId
	}
	return 0
}

func (m *Post) GetId() uint64 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Post) GetTimestamp() int64 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

func (m *Post) GetLikes() int64 {
	if m != nil && m.Likes != nil {
		return *m.Likes
	}
	return 0
}

type IndexResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IndexResponse) Reset()         { *m = IndexResponse{} }
func (m *IndexResponse) String() string { return proto.CompactTextString(m) }
func (*IndexResponse) ProtoMessage()    {}
func (*IndexResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_453745cff914010e, []int{1}
}

func (m *IndexResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IndexResponse.Unmarshal(m, b)
}
func (m *IndexResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IndexResponse.Marshal(b, m, deterministic)
}
func (m *IndexResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IndexResponse.Merge(m, src)
}
func (m *IndexResponse) XXX_Size() int {
	return xxx_messageInfo_IndexResponse.Size(m)
}
func (m *IndexResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_IndexResponse.DiscardUnknown(m)
}

var xxx_messageInfo_IndexResponse proto.InternalMessageInfo

type SearchQuery struct {
	Term                 *string  `protobuf:"bytes,1,req,name=term" json:"term,omitempty"`
	Total                *uint64  `protobuf:"varint,2,req,name=total" json:"total,omitempty"`
	From                 *uint64  `protobuf:"varint,3,opt,name=from" json:"from,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SearchQuery) Reset()         { *m = SearchQuery{} }
func (m *SearchQuery) String() string { return proto.CompactTextString(m) }
func (*SearchQuery) ProtoMessage()    {}
func (*SearchQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_453745cff914010e, []int{2}
}

func (m *SearchQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SearchQuery.Unmarshal(m, b)
}
func (m *SearchQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SearchQuery.Marshal(b, m, deterministic)
}
func (m *SearchQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SearchQuery.Merge(m, src)
}
func (m *SearchQuery) XXX_Size() int {
	return xxx_messageInfo_SearchQuery.Size(m)
}
func (m *SearchQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_SearchQuery.DiscardUnknown(m)
}

var xxx_messageInfo_SearchQuery proto.InternalMessageInfo

func (m *SearchQuery) GetTerm() string {
	if m != nil && m.Term != nil {
		return *m.Term
	}
	return ""
}

func (m *SearchQuery) GetTotal() uint64 {
	if m != nil && m.Total != nil {
		return *m.Total
	}
	return 0
}

func (m *SearchQuery) GetFrom() uint64 {
	if m != nil && m.From != nil {
		return *m.From
	}
	return 0
}

type SearchResult struct {
	Id                   []uint64 `protobuf:"varint,1,rep,name=id" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SearchResult) Reset()         { *m = SearchResult{} }
func (m *SearchResult) String() string { return proto.CompactTextString(m) }
func (*SearchResult) ProtoMessage()    {}
func (*SearchResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_453745cff914010e, []int{3}
}

func (m *SearchResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SearchResult.Unmarshal(m, b)
}
func (m *SearchResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SearchResult.Marshal(b, m, deterministic)
}
func (m *SearchResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SearchResult.Merge(m, src)
}
func (m *SearchResult) XXX_Size() int {
	return xxx_messageInfo_SearchResult.Size(m)
}
func (m *SearchResult) XXX_DiscardUnknown() {
	xxx_messageInfo_SearchResult.DiscardUnknown(m)
}

var xxx_messageInfo_SearchResult proto.InternalMessageInfo

func (m *SearchResult) GetId() []uint64 {
	if m != nil {
		return m.Id
	}
	return nil
}

type Likes struct {
	Id                   *uint64  `protobuf:"varint,1,req,name=id" json:"id,omitempty"`
	Likes                *int64   `protobuf:"varint,2,req,name=likes" json:"likes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Likes) Reset()         { *m = Likes{} }
func (m *Likes) String() string { return proto.CompactTextString(m) }
func (*Likes) ProtoMessage()    {}
func (*Likes) Descriptor() ([]byte, []int) {
	return fileDescriptor_453745cff914010e, []int{4}
}

func (m *Likes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Likes.Unmarshal(m, b)
}
func (m *Likes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Likes.Marshal(b, m, deterministic)
}
func (m *Likes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Likes.Merge(m, src)
}
func (m *Likes) XXX_Size() int {
	return xxx_messageInfo_Likes.Size(m)
}
func (m *Likes) XXX_DiscardUnknown() {
	xxx_messageInfo_Likes.DiscardUnknown(m)
}

var xxx_messageInfo_Likes proto.InternalMessageInfo

func (m *Likes) GetId() uint64 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Likes) GetLikes() int64 {
	if m != nil && m.Likes != nil {
		return *m.Likes
	}
	return 0
}

type LikesResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LikesResponse) Reset()         { *m = LikesResponse{} }
func (m *LikesResponse) String() string { return proto.CompactTextString(m) }
func (*LikesResponse) ProtoMessage()    {}
func (*LikesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_453745cff914010e, []int{5}
}

func (m *LikesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LikesResponse.Unmarshal(m, b)
}
func (m *LikesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LikesResponse.Marshal(b, m, deterministic)
}
func (m *LikesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LikesResponse.Merge(m, src)
}
func (m *LikesResponse) XXX_Size() int {
	return xxx_messageInfo_LikesResponse.Size(m)
}
func (m *LikesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LikesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LikesResponse proto.InternalMessageInfo

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
	return fileDescriptor_453745cff914010e, []int{6}
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
	return fileDescriptor_453745cff914010e, []int{7}
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

type Timestamp struct {
	Id                   *uint64  `protobuf:"varint,1,req,name=id" json:"id,omitempty"`
	Timestamp            *int64   `protobuf:"varint,2,req,name=timestamp" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Timestamp) Reset()         { *m = Timestamp{} }
func (m *Timestamp) String() string { return proto.CompactTextString(m) }
func (*Timestamp) ProtoMessage()    {}
func (*Timestamp) Descriptor() ([]byte, []int) {
	return fileDescriptor_453745cff914010e, []int{8}
}

func (m *Timestamp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Timestamp.Unmarshal(m, b)
}
func (m *Timestamp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Timestamp.Marshal(b, m, deterministic)
}
func (m *Timestamp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Timestamp.Merge(m, src)
}
func (m *Timestamp) XXX_Size() int {
	return xxx_messageInfo_Timestamp.Size(m)
}
func (m *Timestamp) XXX_DiscardUnknown() {
	xxx_messageInfo_Timestamp.DiscardUnknown(m)
}

var xxx_messageInfo_Timestamp proto.InternalMessageInfo

func (m *Timestamp) GetId() uint64 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Timestamp) GetTimestamp() int64 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

type SetTimestampResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetTimestampResponse) Reset()         { *m = SetTimestampResponse{} }
func (m *SetTimestampResponse) String() string { return proto.CompactTextString(m) }
func (*SetTimestampResponse) ProtoMessage()    {}
func (*SetTimestampResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_453745cff914010e, []int{9}
}

func (m *SetTimestampResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetTimestampResponse.Unmarshal(m, b)
}
func (m *SetTimestampResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetTimestampResponse.Marshal(b, m, deterministic)
}
func (m *SetTimestampResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetTimestampResponse.Merge(m, src)
}
func (m *SetTimestampResponse) XXX_Size() int {
	return xxx_messageInfo_SetTimestampResponse.Size(m)
}
func (m *SetTimestampResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SetTimestampResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SetTimestampResponse proto.InternalMessageInfo

type HealthCheckRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HealthCheckRequest) Reset()         { *m = HealthCheckRequest{} }
func (m *HealthCheckRequest) String() string { return proto.CompactTextString(m) }
func (*HealthCheckRequest) ProtoMessage()    {}
func (*HealthCheckRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_453745cff914010e, []int{10}
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
	return fileDescriptor_453745cff914010e, []int{11}
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
	proto.RegisterType((*Post)(nil), "pb.Post")
	proto.RegisterType((*IndexResponse)(nil), "pb.IndexResponse")
	proto.RegisterType((*SearchQuery)(nil), "pb.SearchQuery")
	proto.RegisterType((*SearchResult)(nil), "pb.SearchResult")
	proto.RegisterType((*Likes)(nil), "pb.Likes")
	proto.RegisterType((*LikesResponse)(nil), "pb.LikesResponse")
	proto.RegisterType((*Id)(nil), "pb.Id")
	proto.RegisterType((*DeletePostResponse)(nil), "pb.DeletePostResponse")
	proto.RegisterType((*Timestamp)(nil), "pb.Timestamp")
	proto.RegisterType((*SetTimestampResponse)(nil), "pb.SetTimestampResponse")
	proto.RegisterType((*HealthCheckRequest)(nil), "pb.HealthCheckRequest")
	proto.RegisterType((*HealthCheckResponse)(nil), "pb.HealthCheckResponse")
}

func init() { proto.RegisterFile("search.proto", fileDescriptor_453745cff914010e) }

var fileDescriptor_453745cff914010e = []byte{
	// 487 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x52, 0x5d, 0x6f, 0xd3, 0x30,
	0x14, 0x25, 0x69, 0x9a, 0xad, 0x77, 0xeb, 0x5a, 0x2e, 0x55, 0x17, 0x45, 0x08, 0x45, 0x11, 0x48,
	0x91, 0x80, 0x3e, 0xec, 0x85, 0x0f, 0x89, 0x27, 0x40, 0x50, 0x0d, 0x65, 0xe0, 0x0c, 0x78, 0x9c,
	0xd2, 0xe5, 0x42, 0xa3, 0xa5, 0x4d, 0x88, 0x1d, 0xc4, 0x7e, 0x04, 0x12, 0xbf, 0x82, 0xdf, 0x89,
	0x6c, 0x77, 0x4e, 0x43, 0xe1, 0x29, 0xf7, 0xe3, 0xd8, 0xe7, 0xe4, 0x1c, 0xc3, 0x21, 0xa7, 0xb4,
	0xbe, 0x5c, 0xce, 0xaa, 0xba, 0x14, 0x25, 0xda, 0xd5, 0x22, 0xfc, 0x69, 0x81, 0xf3, 0xbe, 0xe4,
	0x02, 0x27, 0xd0, 0x17, 0xb9, 0x28, 0xc8, 0xb3, 0x02, 0x3b, 0x1a, 0x30, 0xdd, 0x20, 0x82, 0xb3,
	0x28, 0xb3, 0x6b, 0xcf, 0x56, 0x43, 0x55, 0xe3, 0x31, 0xec, 0x35, 0x9c, 0xea, 0x8b, 0x3c, 0xf3,
	0x7a, 0x81, 0x1d, 0x39, 0xcc, 0x95, 0xed, 0x3c, 0xc3, 0x23, 0xb0, 0xf3, 0xcc, 0x73, 0xd4, 0xcc,
	0xce, 0x33, 0xbc, 0x0b, 0x03, 0x91, 0xaf, 0x88, 0x8b, 0x74, 0x55, 0x79, 0xfd, 0xc0, 0x8a, 0x7a,
	0xac, 0x1d, 0x48, 0xc2, 0x22, 0xbf, 0x22, 0xee, 0xb9, 0x6a, 0xa3, 0x9b, 0x70, 0x04, 0xc3, 0xf9,
	0x3a, 0xa3, 0x1f, 0x8c, 0x78, 0x55, 0xae, 0x39, 0x85, 0xa7, 0x70, 0x90, 0x28, 0xd1, 0x1f, 0x1a,
	0xaa, 0xaf, 0xa5, 0x20, 0x41, 0xf5, 0x6a, 0xa3, 0x52, 0xd5, 0x4a, 0x7a, 0x29, 0xd2, 0x42, 0xa9,
	0x74, 0x98, 0x6e, 0x24, 0xf2, 0x4b, 0x5d, 0xae, 0xbc, 0x5e, 0x60, 0x45, 0x0e, 0x53, 0x75, 0x78,
	0x0f, 0x0e, 0xf5, 0x65, 0x8c, 0x78, 0x53, 0x88, 0x8d, 0x62, 0x2b, 0xe8, 0x69, 0xc5, 0xe1, 0x63,
	0xe8, 0xbf, 0x93, 0x32, 0xcc, 0xe2, 0xe6, 0x57, 0x8c, 0x58, 0x49, 0xb1, 0x2d, 0x56, 0xc1, 0x8d,
	0xd8, 0x09, 0xd8, 0xc6, 0x07, 0x73, 0x38, 0x9c, 0x00, 0xbe, 0xa2, 0x82, 0x04, 0x49, 0xa3, 0x0d,
	0xf6, 0x19, 0x0c, 0xce, 0x8d, 0x19, 0x7f, 0xf3, 0x75, 0xac, 0xd3, 0x9c, 0xed, 0x20, 0x9c, 0xc2,
	0x24, 0x21, 0x61, 0x4e, 0x6f, 0xd1, 0xe3, 0x5b, 0x4a, 0x0b, 0xb1, 0x7c, 0xb9, 0xa4, 0xcb, 0x2b,
	0x46, 0xdf, 0x1a, 0xe2, 0x22, 0xfc, 0x65, 0xc1, 0x9d, 0xce, 0x58, 0xa3, 0xf1, 0x05, 0xb8, 0x5c,
	0xa4, 0xa2, 0xe1, 0x8a, 0xf7, 0xe8, 0xe4, 0xc1, 0xac, 0x5a, 0xcc, 0xfe, 0x01, 0x9c, 0x25, 0x54,
	0x7f, 0xcf, 0xd7, 0x5f, 0x13, 0x05, 0x66, 0x9b, 0x43, 0xe1, 0x73, 0x18, 0x76, 0x16, 0x78, 0x00,
	0x7b, 0x1f, 0xe3, 0xd3, 0xf8, 0xec, 0x73, 0x3c, 0xbe, 0x25, 0x9b, 0xe4, 0x35, 0xfb, 0x34, 0x8f,
	0xdf, 0x8c, 0x2d, 0x1c, 0xc1, 0x41, 0x7c, 0x76, 0x7e, 0x71, 0x33, 0xb0, 0x4f, 0x7e, 0xdb, 0xe0,
	0xea, 0x20, 0xf0, 0x3e, 0xf4, 0x55, 0xe0, 0xb8, 0x2f, 0xe9, 0xa5, 0x43, 0xfe, 0x6d, 0x59, 0x75,
	0x5e, 0x01, 0x3e, 0x34, 0xf8, 0x91, 0x5c, 0x6e, 0xbd, 0x08, 0x7f, 0xdc, 0x0e, 0x36, 0xa9, 0x46,
	0xb0, 0x9f, 0x90, 0xd0, 0x41, 0x0e, 0xe4, 0x56, 0x95, 0xfa, 0xda, 0x4e, 0x5e, 0xf8, 0x08, 0xa0,
	0x4d, 0x06, 0x5d, 0xc5, 0x9b, 0xf9, 0x53, 0xf9, 0xdd, 0x4d, 0x0c, 0x9f, 0xc8, 0xd7, 0xd3, 0xda,
	0x8e, 0x43, 0x89, 0x33, 0xad, 0xef, 0x69, 0x21, 0xbb, 0xb9, 0xe0, 0x53, 0xe8, 0x2b, 0x47, 0x71,
	0xba, 0x63, 0xb1, 0x8a, 0xc8, 0x3f, 0xfe, 0x8f, 0xf5, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0x6d,
	0x05, 0x23, 0xde, 0xb1, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SearchClient is the client API for Search service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SearchClient interface {
	Index(ctx context.Context, in *Post, opts ...grpc.CallOption) (*IndexResponse, error)
	Search(ctx context.Context, in *SearchQuery, opts ...grpc.CallOption) (*SearchResult, error)
	SetLikes(ctx context.Context, in *Likes, opts ...grpc.CallOption) (*LikesResponse, error)
	DeletePost(ctx context.Context, in *Id, opts ...grpc.CallOption) (*DeletePostResponse, error)
	SetTimestamp(ctx context.Context, in *Timestamp, opts ...grpc.CallOption) (*SetTimestampResponse, error)
	Check(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error)
}

type searchClient struct {
	cc *grpc.ClientConn
}

func NewSearchClient(cc *grpc.ClientConn) SearchClient {
	return &searchClient{cc}
}

func (c *searchClient) Index(ctx context.Context, in *Post, opts ...grpc.CallOption) (*IndexResponse, error) {
	out := new(IndexResponse)
	err := c.cc.Invoke(ctx, "/pb.Search/Index", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) Search(ctx context.Context, in *SearchQuery, opts ...grpc.CallOption) (*SearchResult, error) {
	out := new(SearchResult)
	err := c.cc.Invoke(ctx, "/pb.Search/Search", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) SetLikes(ctx context.Context, in *Likes, opts ...grpc.CallOption) (*LikesResponse, error) {
	out := new(LikesResponse)
	err := c.cc.Invoke(ctx, "/pb.Search/SetLikes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) DeletePost(ctx context.Context, in *Id, opts ...grpc.CallOption) (*DeletePostResponse, error) {
	out := new(DeletePostResponse)
	err := c.cc.Invoke(ctx, "/pb.Search/DeletePost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) SetTimestamp(ctx context.Context, in *Timestamp, opts ...grpc.CallOption) (*SetTimestampResponse, error) {
	out := new(SetTimestampResponse)
	err := c.cc.Invoke(ctx, "/pb.Search/SetTimestamp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) Check(ctx context.Context, in *HealthCheckRequest, opts ...grpc.CallOption) (*HealthCheckResponse, error) {
	out := new(HealthCheckResponse)
	err := c.cc.Invoke(ctx, "/pb.Search/Check", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SearchServer is the server API for Search service.
type SearchServer interface {
	Index(context.Context, *Post) (*IndexResponse, error)
	Search(context.Context, *SearchQuery) (*SearchResult, error)
	SetLikes(context.Context, *Likes) (*LikesResponse, error)
	DeletePost(context.Context, *Id) (*DeletePostResponse, error)
	SetTimestamp(context.Context, *Timestamp) (*SetTimestampResponse, error)
	Check(context.Context, *HealthCheckRequest) (*HealthCheckResponse, error)
}

// UnimplementedSearchServer can be embedded to have forward compatible implementations.
type UnimplementedSearchServer struct {
}

func (*UnimplementedSearchServer) Index(ctx context.Context, req *Post) (*IndexResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Index not implemented")
}
func (*UnimplementedSearchServer) Search(ctx context.Context, req *SearchQuery) (*SearchResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}
func (*UnimplementedSearchServer) SetLikes(ctx context.Context, req *Likes) (*LikesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetLikes not implemented")
}
func (*UnimplementedSearchServer) DeletePost(ctx context.Context, req *Id) (*DeletePostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePost not implemented")
}
func (*UnimplementedSearchServer) SetTimestamp(ctx context.Context, req *Timestamp) (*SetTimestampResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetTimestamp not implemented")
}
func (*UnimplementedSearchServer) Check(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}

func RegisterSearchServer(s *grpc.Server, srv SearchServer) {
	s.RegisterService(&_Search_serviceDesc, srv)
}

func _Search_Index_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Post)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).Index(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Search/Index",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServer).Index(ctx, req.(*Post))
	}
	return interceptor(ctx, in, info, handler)
}

func _Search_Search_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Search/Search",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServer).Search(ctx, req.(*SearchQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _Search_SetLikes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Likes)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).SetLikes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Search/SetLikes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServer).SetLikes(ctx, req.(*Likes))
	}
	return interceptor(ctx, in, info, handler)
}

func _Search_DeletePost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Id)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).DeletePost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Search/DeletePost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServer).DeletePost(ctx, req.(*Id))
	}
	return interceptor(ctx, in, info, handler)
}

func _Search_SetTimestamp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Timestamp)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).SetTimestamp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Search/SetTimestamp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServer).SetTimestamp(ctx, req.(*Timestamp))
	}
	return interceptor(ctx, in, info, handler)
}

func _Search_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Search/Check",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SearchServer).Check(ctx, req.(*HealthCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Search_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Search",
	HandlerType: (*SearchServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Index",
			Handler:    _Search_Index_Handler,
		},
		{
			MethodName: "Search",
			Handler:    _Search_Search_Handler,
		},
		{
			MethodName: "SetLikes",
			Handler:    _Search_SetLikes_Handler,
		},
		{
			MethodName: "DeletePost",
			Handler:    _Search_DeletePost_Handler,
		},
		{
			MethodName: "SetTimestamp",
			Handler:    _Search_SetTimestamp_Handler,
		},
		{
			MethodName: "Check",
			Handler:    _Search_Check_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "search.proto",
}
