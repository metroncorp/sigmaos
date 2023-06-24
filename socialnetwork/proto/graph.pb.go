// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.20.0
// 	protoc        v3.12.4
// source: socialnetwork/proto/graph.proto

package proto

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type GetFollowersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Followeeid int64 `protobuf:"varint,1,opt,name=followeeid,proto3" json:"followeeid,omitempty"`
}

func (x *GetFollowersRequest) Reset() {
	*x = GetFollowersRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socialnetwork_proto_graph_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFollowersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFollowersRequest) ProtoMessage() {}

func (x *GetFollowersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_socialnetwork_proto_graph_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFollowersRequest.ProtoReflect.Descriptor instead.
func (*GetFollowersRequest) Descriptor() ([]byte, []int) {
	return file_socialnetwork_proto_graph_proto_rawDescGZIP(), []int{0}
}

func (x *GetFollowersRequest) GetFolloweeid() int64 {
	if x != nil {
		return x.Followeeid
	}
	return 0
}

type GetFolloweesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Followerid int64 `protobuf:"varint,1,opt,name=followerid,proto3" json:"followerid,omitempty"`
}

func (x *GetFolloweesRequest) Reset() {
	*x = GetFolloweesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socialnetwork_proto_graph_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFolloweesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFolloweesRequest) ProtoMessage() {}

func (x *GetFolloweesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_socialnetwork_proto_graph_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFolloweesRequest.ProtoReflect.Descriptor instead.
func (*GetFolloweesRequest) Descriptor() ([]byte, []int) {
	return file_socialnetwork_proto_graph_proto_rawDescGZIP(), []int{1}
}

func (x *GetFolloweesRequest) GetFollowerid() int64 {
	if x != nil {
		return x.Followerid
	}
	return 0
}

type GraphGetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok      string  `protobuf:"bytes,1,opt,name=ok,proto3" json:"ok,omitempty"`
	Userids []int64 `protobuf:"varint,2,rep,packed,name=userids,proto3" json:"userids,omitempty"`
}

func (x *GraphGetResponse) Reset() {
	*x = GraphGetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socialnetwork_proto_graph_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GraphGetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GraphGetResponse) ProtoMessage() {}

func (x *GraphGetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_socialnetwork_proto_graph_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GraphGetResponse.ProtoReflect.Descriptor instead.
func (*GraphGetResponse) Descriptor() ([]byte, []int) {
	return file_socialnetwork_proto_graph_proto_rawDescGZIP(), []int{2}
}

func (x *GraphGetResponse) GetOk() string {
	if x != nil {
		return x.Ok
	}
	return ""
}

func (x *GraphGetResponse) GetUserids() []int64 {
	if x != nil {
		return x.Userids
	}
	return nil
}

type FollowRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Followerid int64 `protobuf:"varint,1,opt,name=followerid,proto3" json:"followerid,omitempty"`
	Followeeid int64 `protobuf:"varint,2,opt,name=followeeid,proto3" json:"followeeid,omitempty"`
}

func (x *FollowRequest) Reset() {
	*x = FollowRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socialnetwork_proto_graph_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FollowRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FollowRequest) ProtoMessage() {}

func (x *FollowRequest) ProtoReflect() protoreflect.Message {
	mi := &file_socialnetwork_proto_graph_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FollowRequest.ProtoReflect.Descriptor instead.
func (*FollowRequest) Descriptor() ([]byte, []int) {
	return file_socialnetwork_proto_graph_proto_rawDescGZIP(), []int{3}
}

func (x *FollowRequest) GetFollowerid() int64 {
	if x != nil {
		return x.Followerid
	}
	return 0
}

func (x *FollowRequest) GetFolloweeid() int64 {
	if x != nil {
		return x.Followeeid
	}
	return 0
}

type UnfollowRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Followerid int64 `protobuf:"varint,1,opt,name=followerid,proto3" json:"followerid,omitempty"`
	Followeeid int64 `protobuf:"varint,2,opt,name=followeeid,proto3" json:"followeeid,omitempty"`
}

func (x *UnfollowRequest) Reset() {
	*x = UnfollowRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socialnetwork_proto_graph_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnfollowRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnfollowRequest) ProtoMessage() {}

func (x *UnfollowRequest) ProtoReflect() protoreflect.Message {
	mi := &file_socialnetwork_proto_graph_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnfollowRequest.ProtoReflect.Descriptor instead.
func (*UnfollowRequest) Descriptor() ([]byte, []int) {
	return file_socialnetwork_proto_graph_proto_rawDescGZIP(), []int{4}
}

func (x *UnfollowRequest) GetFollowerid() int64 {
	if x != nil {
		return x.Followerid
	}
	return 0
}

func (x *UnfollowRequest) GetFolloweeid() int64 {
	if x != nil {
		return x.Followeeid
	}
	return 0
}

type FollowWithUnameRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Followeruname string `protobuf:"bytes,1,opt,name=followeruname,proto3" json:"followeruname,omitempty"`
	Followeeuname string `protobuf:"bytes,2,opt,name=followeeuname,proto3" json:"followeeuname,omitempty"`
}

func (x *FollowWithUnameRequest) Reset() {
	*x = FollowWithUnameRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socialnetwork_proto_graph_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FollowWithUnameRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FollowWithUnameRequest) ProtoMessage() {}

func (x *FollowWithUnameRequest) ProtoReflect() protoreflect.Message {
	mi := &file_socialnetwork_proto_graph_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FollowWithUnameRequest.ProtoReflect.Descriptor instead.
func (*FollowWithUnameRequest) Descriptor() ([]byte, []int) {
	return file_socialnetwork_proto_graph_proto_rawDescGZIP(), []int{5}
}

func (x *FollowWithUnameRequest) GetFolloweruname() string {
	if x != nil {
		return x.Followeruname
	}
	return ""
}

func (x *FollowWithUnameRequest) GetFolloweeuname() string {
	if x != nil {
		return x.Followeeuname
	}
	return ""
}

type UnfollowWithUnameRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Followeruname string `protobuf:"bytes,1,opt,name=followeruname,proto3" json:"followeruname,omitempty"`
	Followeeuname string `protobuf:"bytes,2,opt,name=followeeuname,proto3" json:"followeeuname,omitempty"`
}

func (x *UnfollowWithUnameRequest) Reset() {
	*x = UnfollowWithUnameRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socialnetwork_proto_graph_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnfollowWithUnameRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnfollowWithUnameRequest) ProtoMessage() {}

func (x *UnfollowWithUnameRequest) ProtoReflect() protoreflect.Message {
	mi := &file_socialnetwork_proto_graph_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnfollowWithUnameRequest.ProtoReflect.Descriptor instead.
func (*UnfollowWithUnameRequest) Descriptor() ([]byte, []int) {
	return file_socialnetwork_proto_graph_proto_rawDescGZIP(), []int{6}
}

func (x *UnfollowWithUnameRequest) GetFolloweruname() string {
	if x != nil {
		return x.Followeruname
	}
	return ""
}

func (x *UnfollowWithUnameRequest) GetFolloweeuname() string {
	if x != nil {
		return x.Followeeuname
	}
	return ""
}

type GraphUpdateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok string `protobuf:"bytes,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *GraphUpdateResponse) Reset() {
	*x = GraphUpdateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socialnetwork_proto_graph_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GraphUpdateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GraphUpdateResponse) ProtoMessage() {}

func (x *GraphUpdateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_socialnetwork_proto_graph_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GraphUpdateResponse.ProtoReflect.Descriptor instead.
func (*GraphUpdateResponse) Descriptor() ([]byte, []int) {
	return file_socialnetwork_proto_graph_proto_rawDescGZIP(), []int{7}
}

func (x *GraphUpdateResponse) GetOk() string {
	if x != nil {
		return x.Ok
	}
	return ""
}

var File_socialnetwork_proto_graph_proto protoreflect.FileDescriptor

var file_socialnetwork_proto_graph_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x6c, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x35, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x66, 0x6f, 0x6c, 0x6c,
	0x6f, 0x77, 0x65, 0x65, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x66, 0x6f,
	0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x69, 0x64, 0x22, 0x35, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x46,
	0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1e, 0x0a, 0x0a, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0a, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x69, 0x64, 0x22,
	0x3c, 0x0a, 0x10, 0x47, 0x72, 0x61, 0x70, 0x68, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x6f, 0x6b, 0x12, 0x18, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x69, 0x64, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x03, 0x52, 0x07, 0x75, 0x73, 0x65, 0x72, 0x69, 0x64, 0x73, 0x22, 0x4f, 0x0a,
	0x0d, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e,
	0x0a, 0x0a, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0a, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x69, 0x64, 0x12, 0x1e,
	0x0a, 0x0a, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0a, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x69, 0x64, 0x22, 0x51,
	0x0a, 0x0f, 0x55, 0x6e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x69,
	0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x69,
	0x64, 0x22, 0x64, 0x0a, 0x16, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x57, 0x69, 0x74, 0x68, 0x55,
	0x6e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x66,
	0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x75, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x75, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x24, 0x0a, 0x0d, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x75, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77,
	0x65, 0x65, 0x75, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x66, 0x0a, 0x18, 0x55, 0x6e, 0x66, 0x6f, 0x6c,
	0x6c, 0x6f, 0x77, 0x57, 0x69, 0x74, 0x68, 0x55, 0x6e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x75,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x66, 0x6f, 0x6c, 0x6c,
	0x6f, 0x77, 0x65, 0x72, 0x75, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x66, 0x6f, 0x6c,
	0x6c, 0x6f, 0x77, 0x65, 0x65, 0x75, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x75, 0x6e, 0x61, 0x6d, 0x65, 0x22,
	0x25, 0x0a, 0x13, 0x47, 0x72, 0x61, 0x70, 0x68, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x6f, 0x6b, 0x32, 0xec, 0x02, 0x0a, 0x0c, 0x47, 0x72, 0x61, 0x70, 0x68,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x37, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x46, 0x6f,
	0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x72, 0x73, 0x12, 0x14, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x6c,
	0x6c, 0x6f, 0x77, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e,
	0x47, 0x72, 0x61, 0x70, 0x68, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x37, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x73,
	0x12, 0x14, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x65, 0x65, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x47, 0x72, 0x61, 0x70, 0x68, 0x47, 0x65,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x06, 0x46, 0x6f, 0x6c,
	0x6c, 0x6f, 0x77, 0x12, 0x0e, 0x2e, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x47, 0x72, 0x61, 0x70, 0x68, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x08, 0x55, 0x6e, 0x66,
	0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x12, 0x10, 0x2e, 0x55, 0x6e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x47, 0x72, 0x61, 0x70, 0x68, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a,
	0x0f, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x57, 0x69, 0x74, 0x68, 0x55, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x17, 0x2e, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x57, 0x69, 0x74, 0x68, 0x55, 0x6e, 0x61,
	0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x47, 0x72, 0x61, 0x70,
	0x68, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x44, 0x0a, 0x11, 0x55, 0x6e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x57, 0x69, 0x74, 0x68, 0x55,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x19, 0x2e, 0x55, 0x6e, 0x66, 0x6f, 0x6c, 0x6c, 0x6f, 0x77, 0x57,
	0x69, 0x74, 0x68, 0x55, 0x6e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x14, 0x2e, 0x47, 0x72, 0x61, 0x70, 0x68, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x1d, 0x5a, 0x1b, 0x73, 0x69, 0x67, 0x6d, 0x61, 0x6f, 0x73,
	0x2f, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x6c, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_socialnetwork_proto_graph_proto_rawDescOnce sync.Once
	file_socialnetwork_proto_graph_proto_rawDescData = file_socialnetwork_proto_graph_proto_rawDesc
)

func file_socialnetwork_proto_graph_proto_rawDescGZIP() []byte {
	file_socialnetwork_proto_graph_proto_rawDescOnce.Do(func() {
		file_socialnetwork_proto_graph_proto_rawDescData = protoimpl.X.CompressGZIP(file_socialnetwork_proto_graph_proto_rawDescData)
	})
	return file_socialnetwork_proto_graph_proto_rawDescData
}

var file_socialnetwork_proto_graph_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_socialnetwork_proto_graph_proto_goTypes = []interface{}{
	(*GetFollowersRequest)(nil),      // 0: GetFollowersRequest
	(*GetFolloweesRequest)(nil),      // 1: GetFolloweesRequest
	(*GraphGetResponse)(nil),         // 2: GraphGetResponse
	(*FollowRequest)(nil),            // 3: FollowRequest
	(*UnfollowRequest)(nil),          // 4: UnfollowRequest
	(*FollowWithUnameRequest)(nil),   // 5: FollowWithUnameRequest
	(*UnfollowWithUnameRequest)(nil), // 6: UnfollowWithUnameRequest
	(*GraphUpdateResponse)(nil),      // 7: GraphUpdateResponse
}
var file_socialnetwork_proto_graph_proto_depIdxs = []int32{
	0, // 0: GraphService.GetFollowers:input_type -> GetFollowersRequest
	1, // 1: GraphService.GetFollowees:input_type -> GetFolloweesRequest
	3, // 2: GraphService.Follow:input_type -> FollowRequest
	4, // 3: GraphService.Unfollow:input_type -> UnfollowRequest
	5, // 4: GraphService.FollowWithUname:input_type -> FollowWithUnameRequest
	6, // 5: GraphService.UnfollowWithUname:input_type -> UnfollowWithUnameRequest
	2, // 6: GraphService.GetFollowers:output_type -> GraphGetResponse
	2, // 7: GraphService.GetFollowees:output_type -> GraphGetResponse
	7, // 8: GraphService.Follow:output_type -> GraphUpdateResponse
	7, // 9: GraphService.Unfollow:output_type -> GraphUpdateResponse
	7, // 10: GraphService.FollowWithUname:output_type -> GraphUpdateResponse
	7, // 11: GraphService.UnfollowWithUname:output_type -> GraphUpdateResponse
	6, // [6:12] is the sub-list for method output_type
	0, // [0:6] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_socialnetwork_proto_graph_proto_init() }
func file_socialnetwork_proto_graph_proto_init() {
	if File_socialnetwork_proto_graph_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_socialnetwork_proto_graph_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFollowersRequest); i {
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
		file_socialnetwork_proto_graph_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFolloweesRequest); i {
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
		file_socialnetwork_proto_graph_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GraphGetResponse); i {
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
		file_socialnetwork_proto_graph_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FollowRequest); i {
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
		file_socialnetwork_proto_graph_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnfollowRequest); i {
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
		file_socialnetwork_proto_graph_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FollowWithUnameRequest); i {
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
		file_socialnetwork_proto_graph_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnfollowWithUnameRequest); i {
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
		file_socialnetwork_proto_graph_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GraphUpdateResponse); i {
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
			RawDescriptor: file_socialnetwork_proto_graph_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_socialnetwork_proto_graph_proto_goTypes,
		DependencyIndexes: file_socialnetwork_proto_graph_proto_depIdxs,
		MessageInfos:      file_socialnetwork_proto_graph_proto_msgTypes,
	}.Build()
	File_socialnetwork_proto_graph_proto = out.File
	file_socialnetwork_proto_graph_proto_rawDesc = nil
	file_socialnetwork_proto_graph_proto_goTypes = nil
	file_socialnetwork_proto_graph_proto_depIdxs = nil
}
