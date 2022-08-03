// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.2
// source: idl.proto

package idl

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

type ID int32

const (
	ID_NONE     ID = 0
	ID_REQ_ECHO ID = 1
	ID_RES_ECHO ID = 2
)

// Enum value maps for ID.
var (
	ID_name = map[int32]string{
		0: "NONE",
		1: "REQ_ECHO",
		2: "RES_ECHO",
	}
	ID_value = map[string]int32{
		"NONE":     0,
		"REQ_ECHO": 1,
		"RES_ECHO": 2,
	}
)

func (x ID) Enum() *ID {
	p := new(ID)
	*p = x
	return p
}

func (x ID) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ID) Descriptor() protoreflect.EnumDescriptor {
	return file_idl_proto_enumTypes[0].Descriptor()
}

func (ID) Type() protoreflect.EnumType {
	return &file_idl_proto_enumTypes[0]
}

func (x ID) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ID.Descriptor instead.
func (ID) EnumDescriptor() ([]byte, []int) {
	return file_idl_proto_rawDescGZIP(), []int{0}
}

type ReqEcho struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	From    string `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	Message string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ReqEcho) Reset() {
	*x = ReqEcho{}
	if protoimpl.UnsafeEnabled {
		mi := &file_idl_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqEcho) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqEcho) ProtoMessage() {}

func (x *ReqEcho) ProtoReflect() protoreflect.Message {
	mi := &file_idl_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqEcho.ProtoReflect.Descriptor instead.
func (*ReqEcho) Descriptor() ([]byte, []int) {
	return file_idl_proto_rawDescGZIP(), []int{0}
}

func (x *ReqEcho) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *ReqEcho) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type ResEcho struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	To      string `protobuf:"bytes,1,opt,name=to,proto3" json:"to,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ResEcho) Reset() {
	*x = ResEcho{}
	if protoimpl.UnsafeEnabled {
		mi := &file_idl_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResEcho) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResEcho) ProtoMessage() {}

func (x *ResEcho) ProtoReflect() protoreflect.Message {
	mi := &file_idl_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResEcho.ProtoReflect.Descriptor instead.
func (*ResEcho) Descriptor() ([]byte, []int) {
	return file_idl_proto_rawDescGZIP(), []int{1}
}

func (x *ResEcho) GetTo() string {
	if x != nil {
		return x.To
	}
	return ""
}

func (x *ResEcho) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_idl_proto protoreflect.FileDescriptor

var file_idl_proto_rawDesc = []byte{
	0x0a, 0x09, 0x69, 0x64, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x37, 0x0a, 0x07, 0x52,
	0x65, 0x71, 0x45, 0x63, 0x68, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0x33, 0x0a, 0x07, 0x52, 0x65, 0x73, 0x45, 0x63, 0x68, 0x6f, 0x12,
	0x0e, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x74, 0x6f, 0x12,
	0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2a, 0x2a, 0x0a, 0x02, 0x49, 0x44, 0x12,
	0x08, 0x0a, 0x04, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x45, 0x51,
	0x5f, 0x45, 0x43, 0x48, 0x4f, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x45, 0x53, 0x5f, 0x45,
	0x43, 0x48, 0x4f, 0x10, 0x02, 0x42, 0x08, 0x5a, 0x06, 0x2e, 0x2e, 0x2f, 0x69, 0x64, 0x6c, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_idl_proto_rawDescOnce sync.Once
	file_idl_proto_rawDescData = file_idl_proto_rawDesc
)

func file_idl_proto_rawDescGZIP() []byte {
	file_idl_proto_rawDescOnce.Do(func() {
		file_idl_proto_rawDescData = protoimpl.X.CompressGZIP(file_idl_proto_rawDescData)
	})
	return file_idl_proto_rawDescData
}

var file_idl_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_idl_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_idl_proto_goTypes = []interface{}{
	(ID)(0),         // 0: ID
	(*ReqEcho)(nil), // 1: ReqEcho
	(*ResEcho)(nil), // 2: ResEcho
}
var file_idl_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_idl_proto_init() }
func file_idl_proto_init() {
	if File_idl_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_idl_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqEcho); i {
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
		file_idl_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResEcho); i {
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
			RawDescriptor: file_idl_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_idl_proto_goTypes,
		DependencyIndexes: file_idl_proto_depIdxs,
		EnumInfos:         file_idl_proto_enumTypes,
		MessageInfos:      file_idl_proto_msgTypes,
	}.Build()
	File_idl_proto = out.File
	file_idl_proto_rawDesc = nil
	file_idl_proto_goTypes = nil
	file_idl_proto_depIdxs = nil
}
