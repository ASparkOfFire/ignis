// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v3.21.12
// source: proto/types.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type FDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Method           string                   `protobuf:"bytes,1,opt,name=Method,proto3" json:"Method,omitempty"`
	Header           map[string]*HeaderFields `protobuf:"bytes,2,rep,name=Header,proto3" json:"Header,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Body             []byte                   `protobuf:"bytes,3,opt,name=Body,proto3" json:"Body,omitempty"`
	ContentLength    int64                    `protobuf:"varint,4,opt,name=ContentLength,proto3" json:"ContentLength,omitempty"`
	TransferEncoding *StringSlice             `protobuf:"bytes,5,opt,name=TransferEncoding,proto3" json:"TransferEncoding,omitempty"`
	Host             string                   `protobuf:"bytes,6,opt,name=Host,proto3" json:"Host,omitempty"`
	RemoteAddr       string                   `protobuf:"bytes,7,opt,name=RemoteAddr,proto3" json:"RemoteAddr,omitempty"`
	RequestURI       string                   `protobuf:"bytes,8,opt,name=RequestURI,proto3" json:"RequestURI,omitempty"`
	Pattern          string                   `protobuf:"bytes,9,opt,name=Pattern,proto3" json:"Pattern,omitempty"`
}

func (x *FDRequest) Reset() {
	*x = FDRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_types_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FDRequest) ProtoMessage() {}

func (x *FDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_types_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FDRequest.ProtoReflect.Descriptor instead.
func (*FDRequest) Descriptor() ([]byte, []int) {
	return file_proto_types_proto_rawDescGZIP(), []int{0}
}

func (x *FDRequest) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *FDRequest) GetHeader() map[string]*HeaderFields {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *FDRequest) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *FDRequest) GetContentLength() int64 {
	if x != nil {
		return x.ContentLength
	}
	return 0
}

func (x *FDRequest) GetTransferEncoding() *StringSlice {
	if x != nil {
		return x.TransferEncoding
	}
	return nil
}

func (x *FDRequest) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *FDRequest) GetRemoteAddr() string {
	if x != nil {
		return x.RemoteAddr
	}
	return ""
}

func (x *FDRequest) GetRequestURI() string {
	if x != nil {
		return x.RequestURI
	}
	return ""
}

func (x *FDRequest) GetPattern() string {
	if x != nil {
		return x.Pattern
	}
	return ""
}

type FDResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Body       []byte                   `protobuf:"bytes,1,opt,name=Body,proto3" json:"Body,omitempty"`
	StatusCode int32                    `protobuf:"varint,2,opt,name=StatusCode,proto3" json:"StatusCode,omitempty"`
	Length     int32                    `protobuf:"varint,3,opt,name=Length,proto3" json:"Length,omitempty"`
	Header     map[string]*HeaderFields `protobuf:"bytes,4,rep,name=Header,proto3" json:"Header,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *FDResponse) Reset() {
	*x = FDResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_types_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FDResponse) ProtoMessage() {}

func (x *FDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_types_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FDResponse.ProtoReflect.Descriptor instead.
func (*FDResponse) Descriptor() ([]byte, []int) {
	return file_proto_types_proto_rawDescGZIP(), []int{1}
}

func (x *FDResponse) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *FDResponse) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *FDResponse) GetLength() int32 {
	if x != nil {
		return x.Length
	}
	return 0
}

func (x *FDResponse) GetHeader() map[string]*HeaderFields {
	if x != nil {
		return x.Header
	}
	return nil
}

type HeaderFields struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Fields []string `protobuf:"bytes,1,rep,name=fields,proto3" json:"fields,omitempty"`
}

func (x *HeaderFields) Reset() {
	*x = HeaderFields{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_types_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeaderFields) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeaderFields) ProtoMessage() {}

func (x *HeaderFields) ProtoReflect() protoreflect.Message {
	mi := &file_proto_types_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeaderFields.ProtoReflect.Descriptor instead.
func (*HeaderFields) Descriptor() ([]byte, []int) {
	return file_proto_types_proto_rawDescGZIP(), []int{2}
}

func (x *HeaderFields) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

type StringSlice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Fields []string `protobuf:"bytes,1,rep,name=fields,proto3" json:"fields,omitempty"`
}

func (x *StringSlice) Reset() {
	*x = StringSlice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_types_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StringSlice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StringSlice) ProtoMessage() {}

func (x *StringSlice) ProtoReflect() protoreflect.Message {
	mi := &file_proto_types_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StringSlice.ProtoReflect.Descriptor instead.
func (*StringSlice) Descriptor() ([]byte, []int) {
	return file_proto_types_proto_rawDescGZIP(), []int{3}
}

func (x *StringSlice) GetFields() []string {
	if x != nil {
		return x.Fields
	}
	return nil
}

var File_proto_types_proto protoreflect.FileDescriptor

var file_proto_types_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x91, 0x03, 0x0a, 0x09, 0x46,
	0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x68,
	0x6f, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x12, 0x34, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06,
	0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x42, 0x6f, 0x64, 0x79, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x24, 0x0a, 0x0d, 0x43, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0d, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68,
	0x12, 0x3e, 0x0a, 0x10, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x45, 0x6e, 0x63, 0x6f,
	0x64, 0x69, 0x6e, 0x67, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x52, 0x10,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67,
	0x12, 0x12, 0x0a, 0x04, 0x48, 0x6f, 0x73, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x48, 0x6f, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x41, 0x64,
	0x64, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x52, 0x65, 0x6d, 0x6f, 0x74, 0x65,
	0x41, 0x64, 0x64, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x55,
	0x52, 0x49, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x55, 0x52, 0x49, 0x12, 0x18, 0x0a, 0x07, 0x50, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x50, 0x61, 0x74, 0x74, 0x65, 0x72, 0x6e, 0x1a, 0x4e,
	0x0a, 0x0b, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x29, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x73, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xdf,
	0x01, 0x0a, 0x0a, 0x46, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x42, 0x6f, 0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x42, 0x6f, 0x64,
	0x79, 0x12, 0x1e, 0x0a, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x06, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x35, 0x0a, 0x06, 0x48, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x46, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x48, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x1a, 0x4e, 0x0a, 0x0b, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x29, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x73, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x22, 0x26, 0x0a, 0x0c, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73,
	0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x22, 0x25, 0x0a, 0x0b, 0x53, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x53, 0x6c, 0x69, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x42,
	0x1f, 0x5a, 0x1d, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x41, 0x53, 0x70, 0x61, 0x72, 0x6b, 0x4f, 0x66,
	0x46, 0x69, 0x72, 0x65, 0x2f, 0x69, 0x67, 0x6e, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_types_proto_rawDescOnce sync.Once
	file_proto_types_proto_rawDescData = file_proto_types_proto_rawDesc
)

func file_proto_types_proto_rawDescGZIP() []byte {
	file_proto_types_proto_rawDescOnce.Do(func() {
		file_proto_types_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_types_proto_rawDescData)
	})
	return file_proto_types_proto_rawDescData
}

var file_proto_types_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_types_proto_goTypes = []interface{}{
	(*FDRequest)(nil),    // 0: proto.FDRequest
	(*FDResponse)(nil),   // 1: proto.FDResponse
	(*HeaderFields)(nil), // 2: proto.HeaderFields
	(*StringSlice)(nil),  // 3: proto.StringSlice
	nil,                  // 4: proto.FDRequest.HeaderEntry
	nil,                  // 5: proto.FDResponse.HeaderEntry
}
var file_proto_types_proto_depIdxs = []int32{
	4, // 0: proto.FDRequest.Header:type_name -> proto.FDRequest.HeaderEntry
	3, // 1: proto.FDRequest.TransferEncoding:type_name -> proto.StringSlice
	5, // 2: proto.FDResponse.Header:type_name -> proto.FDResponse.HeaderEntry
	2, // 3: proto.FDRequest.HeaderEntry.value:type_name -> proto.HeaderFields
	2, // 4: proto.FDResponse.HeaderEntry.value:type_name -> proto.HeaderFields
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_proto_types_proto_init() }
func file_proto_types_proto_init() {
	if File_proto_types_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_types_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FDRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_types_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FDResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_types_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HeaderFields); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_types_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StringSlice); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_types_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_types_proto_goTypes,
		DependencyIndexes: file_proto_types_proto_depIdxs,
		MessageInfos:      file_proto_types_proto_msgTypes,
	}.Build()
	File_proto_types_proto = out.File
	file_proto_types_proto_rawDesc = nil
	file_proto_types_proto_goTypes = nil
	file_proto_types_proto_depIdxs = nil
}
