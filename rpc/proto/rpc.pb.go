// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.3
// source: rpc/proto/rpc.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sigmap "sigmaos/sigmap"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Method string `protobuf:"bytes,1,opt,name=method,proto3" json:"method,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_rpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_rpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rpc_proto_rawDescGZIP(), []int{0}
}

func (x *Request) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

type Reply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Err *sigmap.Rerror `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *Reply) Reset() {
	*x = Reply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_rpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Reply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Reply) ProtoMessage() {}

func (x *Reply) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_rpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Reply.ProtoReflect.Descriptor instead.
func (*Reply) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rpc_proto_rawDescGZIP(), []int{1}
}

func (x *Reply) GetErr() *sigmap.Rerror {
	if x != nil {
		return x.Err
	}
	return nil
}

// Users of rpc package can use Blob to pass data directly through to
// the transport without the rpc package marshaling it.
type Blob struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Iov [][]byte `protobuf:"bytes,1,rep,name=iov,proto3" json:"iov,omitempty"`
}

func (x *Blob) Reset() {
	*x = Blob{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_rpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Blob) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Blob) ProtoMessage() {}

func (x *Blob) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_rpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Blob.ProtoReflect.Descriptor instead.
func (*Blob) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rpc_proto_rawDescGZIP(), []int{2}
}

func (x *Blob) GetIov() [][]byte {
	if x != nil {
		return x.Iov
	}
	return nil
}

var File_rpc_proto_rpc_proto protoreflect.FileDescriptor

var file_rpc_proto_rpc_proto_rawDesc = []byte{
	0x0a, 0x13, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x70, 0x63, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x73, 0x69, 0x67, 0x6d, 0x61, 0x70, 0x2f, 0x73, 0x69,
	0x67, 0x6d, 0x61, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x21, 0x0a, 0x07, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x22, 0x22, 0x0a,
	0x05, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x19, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x07, 0x2e, 0x52, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x03, 0x65, 0x72,
	0x72, 0x22, 0x18, 0x0a, 0x04, 0x42, 0x6c, 0x6f, 0x62, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x6f, 0x76,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x03, 0x69, 0x6f, 0x76, 0x42, 0x13, 0x5a, 0x11, 0x73,
	0x69, 0x67, 0x6d, 0x61, 0x6f, 0x73, 0x2f, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_proto_rpc_proto_rawDescOnce sync.Once
	file_rpc_proto_rpc_proto_rawDescData = file_rpc_proto_rpc_proto_rawDesc
)

func file_rpc_proto_rpc_proto_rawDescGZIP() []byte {
	file_rpc_proto_rpc_proto_rawDescOnce.Do(func() {
		file_rpc_proto_rpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_proto_rpc_proto_rawDescData)
	})
	return file_rpc_proto_rpc_proto_rawDescData
}

var file_rpc_proto_rpc_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_rpc_proto_rpc_proto_goTypes = []interface{}{
	(*Request)(nil),       // 0: Request
	(*Reply)(nil),         // 1: Reply
	(*Blob)(nil),          // 2: Blob
	(*sigmap.Rerror)(nil), // 3: Rerror
}
var file_rpc_proto_rpc_proto_depIdxs = []int32{
	3, // 0: Reply.err:type_name -> Rerror
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rpc_proto_rpc_proto_init() }
func file_rpc_proto_rpc_proto_init() {
	if File_rpc_proto_rpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_proto_rpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
		file_rpc_proto_rpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Reply); i {
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
		file_rpc_proto_rpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Blob); i {
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
			RawDescriptor: file_rpc_proto_rpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rpc_proto_rpc_proto_goTypes,
		DependencyIndexes: file_rpc_proto_rpc_proto_depIdxs,
		MessageInfos:      file_rpc_proto_rpc_proto_msgTypes,
	}.Build()
	File_rpc_proto_rpc_proto = out.File
	file_rpc_proto_rpc_proto_rawDesc = nil
	file_rpc_proto_rpc_proto_goTypes = nil
	file_rpc_proto_rpc_proto_depIdxs = nil
}
