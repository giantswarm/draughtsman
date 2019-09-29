// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v2/enums/media_type.proto

package enums

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// The type of media.
type MediaTypeEnum_MediaType int32

const (
	// The media type has not been specified.
	MediaTypeEnum_UNSPECIFIED MediaTypeEnum_MediaType = 0
	// The received value is not known in this version.
	//
	// This is a response-only value.
	MediaTypeEnum_UNKNOWN MediaTypeEnum_MediaType = 1
	// Static image, used for image ad.
	MediaTypeEnum_IMAGE MediaTypeEnum_MediaType = 2
	// Small image, used for map ad.
	MediaTypeEnum_ICON MediaTypeEnum_MediaType = 3
	// ZIP file, used in fields of template ads.
	MediaTypeEnum_MEDIA_BUNDLE MediaTypeEnum_MediaType = 4
	// Audio file.
	MediaTypeEnum_AUDIO MediaTypeEnum_MediaType = 5
	// Video file.
	MediaTypeEnum_VIDEO MediaTypeEnum_MediaType = 6
	// Animated image, such as animated GIF.
	MediaTypeEnum_DYNAMIC_IMAGE MediaTypeEnum_MediaType = 7
)

var MediaTypeEnum_MediaType_name = map[int32]string{
	0: "UNSPECIFIED",
	1: "UNKNOWN",
	2: "IMAGE",
	3: "ICON",
	4: "MEDIA_BUNDLE",
	5: "AUDIO",
	6: "VIDEO",
	7: "DYNAMIC_IMAGE",
}

var MediaTypeEnum_MediaType_value = map[string]int32{
	"UNSPECIFIED":   0,
	"UNKNOWN":       1,
	"IMAGE":         2,
	"ICON":          3,
	"MEDIA_BUNDLE":  4,
	"AUDIO":         5,
	"VIDEO":         6,
	"DYNAMIC_IMAGE": 7,
}

func (x MediaTypeEnum_MediaType) String() string {
	return proto.EnumName(MediaTypeEnum_MediaType_name, int32(x))
}

func (MediaTypeEnum_MediaType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e1c9415301a4e420, []int{0, 0}
}

// Container for enum describing the types of media.
type MediaTypeEnum struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MediaTypeEnum) Reset()         { *m = MediaTypeEnum{} }
func (m *MediaTypeEnum) String() string { return proto.CompactTextString(m) }
func (*MediaTypeEnum) ProtoMessage()    {}
func (*MediaTypeEnum) Descriptor() ([]byte, []int) {
	return fileDescriptor_e1c9415301a4e420, []int{0}
}

func (m *MediaTypeEnum) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MediaTypeEnum.Unmarshal(m, b)
}
func (m *MediaTypeEnum) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MediaTypeEnum.Marshal(b, m, deterministic)
}
func (m *MediaTypeEnum) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MediaTypeEnum.Merge(m, src)
}
func (m *MediaTypeEnum) XXX_Size() int {
	return xxx_messageInfo_MediaTypeEnum.Size(m)
}
func (m *MediaTypeEnum) XXX_DiscardUnknown() {
	xxx_messageInfo_MediaTypeEnum.DiscardUnknown(m)
}

var xxx_messageInfo_MediaTypeEnum proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("google.ads.googleads.v2.enums.MediaTypeEnum_MediaType", MediaTypeEnum_MediaType_name, MediaTypeEnum_MediaType_value)
	proto.RegisterType((*MediaTypeEnum)(nil), "google.ads.googleads.v2.enums.MediaTypeEnum")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v2/enums/media_type.proto", fileDescriptor_e1c9415301a4e420)
}

var fileDescriptor_e1c9415301a4e420 = []byte{
	// 338 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x90, 0x4f, 0x4e, 0x83, 0x40,
	0x18, 0xc5, 0x85, 0xfe, 0xb3, 0x53, 0xab, 0x23, 0x4b, 0x63, 0x17, 0xed, 0x01, 0x86, 0x04, 0x77,
	0xe3, 0x6a, 0x28, 0xd8, 0x4c, 0x14, 0x68, 0xa2, 0x60, 0x34, 0x24, 0x0d, 0x0a, 0x21, 0x24, 0x65,
	0x86, 0x74, 0x68, 0x93, 0x5e, 0xc1, 0x63, 0xb8, 0xf4, 0x28, 0x5e, 0xc3, 0x9d, 0xa7, 0x30, 0xc3,
	0x58, 0x76, 0xba, 0x99, 0xbc, 0xcc, 0xf7, 0xfb, 0x5e, 0xbe, 0xf7, 0x00, 0xca, 0x39, 0xcf, 0xd7,
	0x99, 0x99, 0xa4, 0xc2, 0x54, 0x52, 0xaa, 0x9d, 0x65, 0x66, 0x6c, 0x5b, 0x0a, 0xb3, 0xcc, 0xd2,
	0x22, 0x59, 0xd5, 0xfb, 0x2a, 0x43, 0xd5, 0x86, 0xd7, 0xdc, 0x98, 0x28, 0x08, 0x25, 0xa9, 0x40,
	0x2d, 0x8f, 0x76, 0x16, 0x6a, 0xf8, 0x8b, 0xcb, 0x83, 0x5d, 0x55, 0x98, 0x09, 0x63, 0xbc, 0x4e,
	0xea, 0x82, 0x33, 0xa1, 0x96, 0x67, 0x6f, 0x1a, 0x18, 0x7b, 0xd2, 0xf1, 0x61, 0x5f, 0x65, 0x2e,
	0xdb, 0x96, 0xb3, 0x3d, 0x18, 0xb6, 0x1f, 0xc6, 0x19, 0x18, 0x85, 0xfe, 0xfd, 0xd2, 0x9d, 0xd3,
	0x1b, 0xea, 0x3a, 0xf0, 0xc8, 0x18, 0x81, 0x41, 0xe8, 0xdf, 0xfa, 0xc1, 0xa3, 0x0f, 0x35, 0x63,
	0x08, 0x7a, 0xd4, 0x23, 0x0b, 0x17, 0xea, 0xc6, 0x31, 0xe8, 0xd2, 0x79, 0xe0, 0xc3, 0x8e, 0x01,
	0xc1, 0x89, 0xe7, 0x3a, 0x94, 0xac, 0xec, 0xd0, 0x77, 0xee, 0x5c, 0xd8, 0x95, 0x18, 0x09, 0x1d,
	0x1a, 0xc0, 0x9e, 0x94, 0x11, 0x75, 0xdc, 0x00, 0xf6, 0x8d, 0x73, 0x30, 0x76, 0x9e, 0x7c, 0xe2,
	0xd1, 0xf9, 0x4a, 0x99, 0x0c, 0xec, 0x2f, 0x0d, 0x4c, 0x5f, 0x79, 0x89, 0xfe, 0x0d, 0x64, 0x9f,
	0xb6, 0xe7, 0x2d, 0x65, 0x84, 0xa5, 0xf6, 0x6c, 0xff, 0x2e, 0xe4, 0x7c, 0x9d, 0xb0, 0x1c, 0xf1,
	0x4d, 0x6e, 0xe6, 0x19, 0x6b, 0x02, 0x1e, 0x1a, 0xac, 0x0a, 0xf1, 0x47, 0xa1, 0xd7, 0xcd, 0xfb,
	0xae, 0x77, 0x16, 0x84, 0x7c, 0xe8, 0x93, 0x85, 0xb2, 0x22, 0xa9, 0x40, 0x4a, 0x4a, 0x15, 0x59,
	0x48, 0x76, 0x23, 0x3e, 0x0f, 0xf3, 0x98, 0xa4, 0x22, 0x6e, 0xe7, 0x71, 0x64, 0xc5, 0xcd, 0xfc,
	0x5b, 0x9f, 0xaa, 0x4f, 0x8c, 0x49, 0x2a, 0x30, 0x6e, 0x09, 0x8c, 0x23, 0x0b, 0xe3, 0x86, 0x79,
	0xe9, 0x37, 0x87, 0x5d, 0xfd, 0x04, 0x00, 0x00, 0xff, 0xff, 0x42, 0x0d, 0xc8, 0x96, 0xe8, 0x01,
	0x00, 0x00,
}
