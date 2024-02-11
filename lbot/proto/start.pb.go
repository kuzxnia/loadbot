// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.2
// source: start/start.proto

package proto

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type StartRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Watch bool `protobuf:"varint,1,opt,name=watch,proto3" json:"watch,omitempty"`
}

func (x *StartRequest) Reset() {
	*x = StartRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_start_start_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartRequest) ProtoMessage() {}

func (x *StartRequest) ProtoReflect() protoreflect.Message {
	mi := &file_start_start_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartRequest.ProtoReflect.Descriptor instead.
func (*StartRequest) Descriptor() ([]byte, []int) {
	return file_start_start_proto_rawDescGZIP(), []int{0}
}

func (x *StartRequest) GetWatch() bool {
	if x != nil {
		return x.Watch
	}
	return false
}

type StartResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StartResponse) Reset() {
	*x = StartResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_start_start_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartResponse) ProtoMessage() {}

func (x *StartResponse) ProtoReflect() protoreflect.Message {
	mi := &file_start_start_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartResponse.ProtoReflect.Descriptor instead.
func (*StartResponse) Descriptor() ([]byte, []int) {
	return file_start_start_proto_rawDescGZIP(), []int{1}
}

var File_start_start_proto protoreflect.FileDescriptor

var file_start_start_proto_rawDesc = []byte{
	0x0a, 0x11, 0x73, 0x74, 0x61, 0x72, 0x74, 0x2f, 0x73, 0x74, 0x61, 0x72, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x22, 0x24, 0x0a, 0x0c, 0x53, 0x74,
	0x61, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x77, 0x61,
	0x74, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x77, 0x61, 0x74, 0x63, 0x68,
	0x22, 0x0f, 0x0a, 0x0d, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x32, 0x42, 0x0a, 0x0c, 0x53, 0x74, 0x61, 0x72, 0x74, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73,
	0x73, 0x12, 0x32, 0x0a, 0x03, 0x52, 0x75, 0x6e, 0x12, 0x13, 0x2e, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x08, 0x5a, 0x06, 0x73, 0x74, 0x61, 0x72, 0x74, 0x2f, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_start_start_proto_rawDescOnce sync.Once
	file_start_start_proto_rawDescData = file_start_start_proto_rawDesc
)

func file_start_start_proto_rawDescGZIP() []byte {
	file_start_start_proto_rawDescOnce.Do(func() {
		file_start_start_proto_rawDescData = protoimpl.X.CompressGZIP(file_start_start_proto_rawDescData)
	})
	return file_start_start_proto_rawDescData
}

var (
	file_start_start_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
	file_start_start_proto_goTypes  = []interface{}{
		(*StartRequest)(nil),  // 0: start.StartRequest
		(*StartResponse)(nil), // 1: start.StartResponse
	}
)
var file_start_start_proto_depIdxs = []int32{
	0, // 0: start.StartProcess.Run:input_type -> start.StartRequest
	1, // 1: start.StartProcess.Run:output_type -> start.StartResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_start_start_proto_init() }
func file_start_start_proto_init() {
	if File_start_start_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_start_start_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartRequest); i {
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
		file_start_start_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartResponse); i {
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
			RawDescriptor: file_start_start_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_start_start_proto_goTypes,
		DependencyIndexes: file_start_start_proto_depIdxs,
		MessageInfos:      file_start_start_proto_msgTypes,
	}.Build()
	File_start_start_proto = out.File
	file_start_start_proto_rawDesc = nil
	file_start_start_proto_goTypes = nil
	file_start_start_proto_depIdxs = nil
}
