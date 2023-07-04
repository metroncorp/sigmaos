// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: socialnetwork/proto/post.proto

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

type POST_TYPE int32

const (
	POST_TYPE_UNKNOWN POST_TYPE = 0
	POST_TYPE_POST    POST_TYPE = 1
	POST_TYPE_REPOST  POST_TYPE = 2
	POST_TYPE_REPLY   POST_TYPE = 3
	POST_TYPE_DM      POST_TYPE = 4
)

// Enum value maps for POST_TYPE.
var (
	POST_TYPE_name = map[int32]string{
		0: "UNKNOWN",
		1: "POST",
		2: "REPOST",
		3: "REPLY",
		4: "DM",
	}
	POST_TYPE_value = map[string]int32{
		"UNKNOWN": 0,
		"POST":    1,
		"REPOST":  2,
		"REPLY":   3,
		"DM":      4,
	}
)

func (x POST_TYPE) Enum() *POST_TYPE {
	p := new(POST_TYPE)
	*p = x
	return p
}

func (x POST_TYPE) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (POST_TYPE) Descriptor() protoreflect.EnumDescriptor {
	return file_socialnetwork_proto_post_proto_enumTypes[0].Descriptor()
}

func (POST_TYPE) Type() protoreflect.EnumType {
	return &file_socialnetwork_proto_post_proto_enumTypes[0]
}

func (x POST_TYPE) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use POST_TYPE.Descriptor instead.
func (POST_TYPE) EnumDescriptor() ([]byte, []int) {
	return file_socialnetwork_proto_post_proto_rawDescGZIP(), []int{0}
}

type StorePostRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Post *Post `protobuf:"bytes,1,opt,name=post,proto3" json:"post,omitempty"`
}

func (x *StorePostRequest) Reset() {
	*x = StorePostRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socialnetwork_proto_post_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StorePostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StorePostRequest) ProtoMessage() {}

func (x *StorePostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_socialnetwork_proto_post_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StorePostRequest.ProtoReflect.Descriptor instead.
func (*StorePostRequest) Descriptor() ([]byte, []int) {
	return file_socialnetwork_proto_post_proto_rawDescGZIP(), []int{0}
}

func (x *StorePostRequest) GetPost() *Post {
	if x != nil {
		return x.Post
	}
	return nil
}

type StorePostResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok string `protobuf:"bytes,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *StorePostResponse) Reset() {
	*x = StorePostResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socialnetwork_proto_post_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StorePostResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StorePostResponse) ProtoMessage() {}

func (x *StorePostResponse) ProtoReflect() protoreflect.Message {
	mi := &file_socialnetwork_proto_post_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StorePostResponse.ProtoReflect.Descriptor instead.
func (*StorePostResponse) Descriptor() ([]byte, []int) {
	return file_socialnetwork_proto_post_proto_rawDescGZIP(), []int{1}
}

func (x *StorePostResponse) GetOk() string {
	if x != nil {
		return x.Ok
	}
	return ""
}

type ReadPostsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Postids []int64 `protobuf:"varint,1,rep,packed,name=postids,proto3" json:"postids,omitempty"`
}

func (x *ReadPostsRequest) Reset() {
	*x = ReadPostsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socialnetwork_proto_post_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadPostsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadPostsRequest) ProtoMessage() {}

func (x *ReadPostsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_socialnetwork_proto_post_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadPostsRequest.ProtoReflect.Descriptor instead.
func (*ReadPostsRequest) Descriptor() ([]byte, []int) {
	return file_socialnetwork_proto_post_proto_rawDescGZIP(), []int{2}
}

func (x *ReadPostsRequest) GetPostids() []int64 {
	if x != nil {
		return x.Postids
	}
	return nil
}

type ReadPostsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok    string  `protobuf:"bytes,1,opt,name=ok,proto3" json:"ok,omitempty"`
	Posts []*Post `protobuf:"bytes,2,rep,name=posts,proto3" json:"posts,omitempty"`
}

func (x *ReadPostsResponse) Reset() {
	*x = ReadPostsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socialnetwork_proto_post_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadPostsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadPostsResponse) ProtoMessage() {}

func (x *ReadPostsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_socialnetwork_proto_post_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadPostsResponse.ProtoReflect.Descriptor instead.
func (*ReadPostsResponse) Descriptor() ([]byte, []int) {
	return file_socialnetwork_proto_post_proto_rawDescGZIP(), []int{3}
}

func (x *ReadPostsResponse) GetOk() string {
	if x != nil {
		return x.Ok
	}
	return ""
}

func (x *ReadPostsResponse) GetPosts() []*Post {
	if x != nil {
		return x.Posts
	}
	return nil
}

type StoreMediaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Mediatype string `protobuf:"bytes,1,opt,name=mediatype,proto3" json:"mediatype,omitempty"`
	Mediadata []byte `protobuf:"bytes,2,opt,name=mediadata,proto3" json:"mediadata,omitempty"`
}

func (x *StoreMediaRequest) Reset() {
	*x = StoreMediaRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socialnetwork_proto_post_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreMediaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreMediaRequest) ProtoMessage() {}

func (x *StoreMediaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_socialnetwork_proto_post_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreMediaRequest.ProtoReflect.Descriptor instead.
func (*StoreMediaRequest) Descriptor() ([]byte, []int) {
	return file_socialnetwork_proto_post_proto_rawDescGZIP(), []int{4}
}

func (x *StoreMediaRequest) GetMediatype() string {
	if x != nil {
		return x.Mediatype
	}
	return ""
}

func (x *StoreMediaRequest) GetMediadata() []byte {
	if x != nil {
		return x.Mediadata
	}
	return nil
}

type StoreMediaResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok      string `protobuf:"bytes,1,opt,name=ok,proto3" json:"ok,omitempty"`
	Mediaid int64  `protobuf:"varint,2,opt,name=mediaid,proto3" json:"mediaid,omitempty"`
}

func (x *StoreMediaResponse) Reset() {
	*x = StoreMediaResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socialnetwork_proto_post_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreMediaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreMediaResponse) ProtoMessage() {}

func (x *StoreMediaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_socialnetwork_proto_post_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreMediaResponse.ProtoReflect.Descriptor instead.
func (*StoreMediaResponse) Descriptor() ([]byte, []int) {
	return file_socialnetwork_proto_post_proto_rawDescGZIP(), []int{5}
}

func (x *StoreMediaResponse) GetOk() string {
	if x != nil {
		return x.Ok
	}
	return ""
}

func (x *StoreMediaResponse) GetMediaid() int64 {
	if x != nil {
		return x.Mediaid
	}
	return 0
}

type ReadMediaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Mediaids []int64 `protobuf:"varint,1,rep,packed,name=mediaids,proto3" json:"mediaids,omitempty"`
}

func (x *ReadMediaRequest) Reset() {
	*x = ReadMediaRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socialnetwork_proto_post_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadMediaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadMediaRequest) ProtoMessage() {}

func (x *ReadMediaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_socialnetwork_proto_post_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadMediaRequest.ProtoReflect.Descriptor instead.
func (*ReadMediaRequest) Descriptor() ([]byte, []int) {
	return file_socialnetwork_proto_post_proto_rawDescGZIP(), []int{6}
}

func (x *ReadMediaRequest) GetMediaids() []int64 {
	if x != nil {
		return x.Mediaids
	}
	return nil
}

type ReadMediaResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok         string   `protobuf:"bytes,1,opt,name=ok,proto3" json:"ok,omitempty"`
	Mediatypes []string `protobuf:"bytes,2,rep,name=mediatypes,proto3" json:"mediatypes,omitempty"`
	Mediadatas [][]byte `protobuf:"bytes,3,rep,name=mediadatas,proto3" json:"mediadatas,omitempty"`
}

func (x *ReadMediaResponse) Reset() {
	*x = ReadMediaResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socialnetwork_proto_post_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReadMediaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadMediaResponse) ProtoMessage() {}

func (x *ReadMediaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_socialnetwork_proto_post_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadMediaResponse.ProtoReflect.Descriptor instead.
func (*ReadMediaResponse) Descriptor() ([]byte, []int) {
	return file_socialnetwork_proto_post_proto_rawDescGZIP(), []int{7}
}

func (x *ReadMediaResponse) GetOk() string {
	if x != nil {
		return x.Ok
	}
	return ""
}

func (x *ReadMediaResponse) GetMediatypes() []string {
	if x != nil {
		return x.Mediatypes
	}
	return nil
}

func (x *ReadMediaResponse) GetMediadatas() [][]byte {
	if x != nil {
		return x.Mediadatas
	}
	return nil
}

type ComposePostRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string    `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Userid   int64     `protobuf:"varint,2,opt,name=userid,proto3" json:"userid,omitempty"`
	Text     string    `protobuf:"bytes,3,opt,name=text,proto3" json:"text,omitempty"`
	Posttype POST_TYPE `protobuf:"varint,4,opt,name=posttype,proto3,enum=POST_TYPE" json:"posttype,omitempty"`
	Mediaids []int64   `protobuf:"varint,5,rep,packed,name=mediaids,proto3" json:"mediaids,omitempty"`
}

func (x *ComposePostRequest) Reset() {
	*x = ComposePostRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socialnetwork_proto_post_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ComposePostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ComposePostRequest) ProtoMessage() {}

func (x *ComposePostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_socialnetwork_proto_post_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ComposePostRequest.ProtoReflect.Descriptor instead.
func (*ComposePostRequest) Descriptor() ([]byte, []int) {
	return file_socialnetwork_proto_post_proto_rawDescGZIP(), []int{8}
}

func (x *ComposePostRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *ComposePostRequest) GetUserid() int64 {
	if x != nil {
		return x.Userid
	}
	return 0
}

func (x *ComposePostRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *ComposePostRequest) GetPosttype() POST_TYPE {
	if x != nil {
		return x.Posttype
	}
	return POST_TYPE_UNKNOWN
}

func (x *ComposePostRequest) GetMediaids() []int64 {
	if x != nil {
		return x.Mediaids
	}
	return nil
}

type ComposePostResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok string `protobuf:"bytes,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *ComposePostResponse) Reset() {
	*x = ComposePostResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socialnetwork_proto_post_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ComposePostResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ComposePostResponse) ProtoMessage() {}

func (x *ComposePostResponse) ProtoReflect() protoreflect.Message {
	mi := &file_socialnetwork_proto_post_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ComposePostResponse.ProtoReflect.Descriptor instead.
func (*ComposePostResponse) Descriptor() ([]byte, []int) {
	return file_socialnetwork_proto_post_proto_rawDescGZIP(), []int{9}
}

func (x *ComposePostResponse) GetOk() string {
	if x != nil {
		return x.Ok
	}
	return ""
}

type Post struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Postid       int64     `protobuf:"varint,1,opt,name=postid,proto3" json:"postid,omitempty"`
	Posttype     POST_TYPE `protobuf:"varint,2,opt,name=posttype,proto3,enum=POST_TYPE" json:"posttype,omitempty"`
	Timestamp    int64     `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Creator      int64     `protobuf:"varint,4,opt,name=creator,proto3" json:"creator,omitempty"`
	Creatoruname string    `protobuf:"bytes,5,opt,name=creatoruname,proto3" json:"creatoruname,omitempty"`
	Text         string    `protobuf:"bytes,6,opt,name=text,proto3" json:"text,omitempty"`
	Usermentions []int64   `protobuf:"varint,7,rep,packed,name=usermentions,proto3" json:"usermentions,omitempty"`
	Medias       []int64   `protobuf:"varint,8,rep,packed,name=medias,proto3" json:"medias,omitempty"`
	Urls         []string  `protobuf:"bytes,9,rep,name=urls,proto3" json:"urls,omitempty"`
}

func (x *Post) Reset() {
	*x = Post{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socialnetwork_proto_post_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Post) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Post) ProtoMessage() {}

func (x *Post) ProtoReflect() protoreflect.Message {
	mi := &file_socialnetwork_proto_post_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Post.ProtoReflect.Descriptor instead.
func (*Post) Descriptor() ([]byte, []int) {
	return file_socialnetwork_proto_post_proto_rawDescGZIP(), []int{10}
}

func (x *Post) GetPostid() int64 {
	if x != nil {
		return x.Postid
	}
	return 0
}

func (x *Post) GetPosttype() POST_TYPE {
	if x != nil {
		return x.Posttype
	}
	return POST_TYPE_UNKNOWN
}

func (x *Post) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *Post) GetCreator() int64 {
	if x != nil {
		return x.Creator
	}
	return 0
}

func (x *Post) GetCreatoruname() string {
	if x != nil {
		return x.Creatoruname
	}
	return ""
}

func (x *Post) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *Post) GetUsermentions() []int64 {
	if x != nil {
		return x.Usermentions
	}
	return nil
}

func (x *Post) GetMedias() []int64 {
	if x != nil {
		return x.Medias
	}
	return nil
}

func (x *Post) GetUrls() []string {
	if x != nil {
		return x.Urls
	}
	return nil
}

var File_socialnetwork_proto_post_proto protoreflect.FileDescriptor

var file_socialnetwork_proto_post_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x6c, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6f, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x2d, 0x0a, 0x10, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x04, 0x70, 0x6f, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x05, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x04, 0x70, 0x6f, 0x73, 0x74, 0x22,
	0x23, 0x0a, 0x11, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x6f, 0x6b, 0x22, 0x2c, 0x0a, 0x10, 0x52, 0x65, 0x61, 0x64, 0x50, 0x6f, 0x73, 0x74,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x6f, 0x73, 0x74,
	0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x03, 0x52, 0x07, 0x70, 0x6f, 0x73, 0x74, 0x69,
	0x64, 0x73, 0x22, 0x40, 0x0a, 0x11, 0x52, 0x65, 0x61, 0x64, 0x50, 0x6f, 0x73, 0x74, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x6f, 0x6b, 0x12, 0x1b, 0x0a, 0x05, 0x70, 0x6f, 0x73, 0x74, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x05, 0x70,
	0x6f, 0x73, 0x74, 0x73, 0x22, 0x4f, 0x0a, 0x11, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x4d, 0x65, 0x64,
	0x69, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x6d, 0x65, 0x64,
	0x69, 0x61, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x65,
	0x64, 0x69, 0x61, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6d, 0x65, 0x64, 0x69, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x6d, 0x65, 0x64, 0x69,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x22, 0x3e, 0x0a, 0x12, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x4d, 0x65,
	0x64, 0x69, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f,
	0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x6f, 0x6b, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x64, 0x69, 0x61, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x6d, 0x65,
	0x64, 0x69, 0x61, 0x69, 0x64, 0x22, 0x2e, 0x0a, 0x10, 0x52, 0x65, 0x61, 0x64, 0x4d, 0x65, 0x64,
	0x69, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x64,
	0x69, 0x61, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x03, 0x52, 0x08, 0x6d, 0x65, 0x64,
	0x69, 0x61, 0x69, 0x64, 0x73, 0x22, 0x63, 0x0a, 0x11, 0x52, 0x65, 0x61, 0x64, 0x4d, 0x65, 0x64,
	0x69, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x6f, 0x6b, 0x12, 0x1e, 0x0a, 0x0a, 0x6d, 0x65,
	0x64, 0x69, 0x61, 0x74, 0x79, 0x70, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a,
	0x6d, 0x65, 0x64, 0x69, 0x61, 0x74, 0x79, 0x70, 0x65, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x6d, 0x65,
	0x64, 0x69, 0x61, 0x64, 0x61, 0x74, 0x61, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x0a,
	0x6d, 0x65, 0x64, 0x69, 0x61, 0x64, 0x61, 0x74, 0x61, 0x73, 0x22, 0xa0, 0x01, 0x0a, 0x12, 0x43,
	0x6f, 0x6d, 0x70, 0x6f, 0x73, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x26, 0x0a, 0x08, 0x70, 0x6f, 0x73,
	0x74, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0a, 0x2e, 0x50, 0x4f,
	0x53, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x52, 0x08, 0x70, 0x6f, 0x73, 0x74, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x69, 0x64, 0x73, 0x18, 0x05, 0x20,
	0x03, 0x28, 0x03, 0x52, 0x08, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x69, 0x64, 0x73, 0x22, 0x25, 0x0a,
	0x13, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x73, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x6f, 0x6b, 0x22, 0x86, 0x02, 0x0a, 0x04, 0x50, 0x6f, 0x73, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x70, 0x6f, 0x73, 0x74, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x70,
	0x6f, 0x73, 0x74, 0x69, 0x64, 0x12, 0x26, 0x0a, 0x08, 0x70, 0x6f, 0x73, 0x74, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0a, 0x2e, 0x50, 0x4f, 0x53, 0x54, 0x5f, 0x54,
	0x59, 0x50, 0x45, 0x52, 0x08, 0x70, 0x6f, 0x73, 0x74, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x22, 0x0a, 0x0c, 0x63, 0x72, 0x65, 0x61, 0x74, 0x6f, 0x72,
	0x75, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x6f, 0x72, 0x75, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78,
	0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x22, 0x0a,
	0x0c, 0x75, 0x73, 0x65, 0x72, 0x6d, 0x65, 0x6e, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x07, 0x20,
	0x03, 0x28, 0x03, 0x52, 0x0c, 0x75, 0x73, 0x65, 0x72, 0x6d, 0x65, 0x6e, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28,
	0x03, 0x52, 0x06, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x72, 0x6c,
	0x73, 0x18, 0x09, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x75, 0x72, 0x6c, 0x73, 0x2a, 0x41, 0x0a,
	0x09, 0x50, 0x4f, 0x53, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e,
	0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x50, 0x4f, 0x53, 0x54, 0x10,
	0x01, 0x12, 0x0a, 0x0a, 0x06, 0x52, 0x45, 0x50, 0x4f, 0x53, 0x54, 0x10, 0x02, 0x12, 0x09, 0x0a,
	0x05, 0x52, 0x45, 0x50, 0x4c, 0x59, 0x10, 0x03, 0x12, 0x06, 0x0a, 0x02, 0x44, 0x4d, 0x10, 0x04,
	0x32, 0x75, 0x0a, 0x0b, 0x50, 0x6f, 0x73, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x32, 0x0a, 0x09, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x12, 0x11, 0x2e, 0x53,
	0x74, 0x6f, 0x72, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x12, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x09, 0x52, 0x65, 0x61, 0x64, 0x50, 0x6f, 0x73, 0x74, 0x73,
	0x12, 0x11, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x50, 0x6f, 0x73, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x50, 0x6f, 0x73, 0x74, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x4a, 0x0a, 0x0e, 0x43, 0x6f, 0x6d, 0x70, 0x6f,
	0x73, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x38, 0x0a, 0x0b, 0x43, 0x6f, 0x6d,
	0x70, 0x6f, 0x73, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x12, 0x13, 0x2e, 0x43, 0x6f, 0x6d, 0x70, 0x6f,
	0x73, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e,
	0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x73, 0x65, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x32, 0x79, 0x0a, 0x0c, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x35, 0x0a, 0x0a, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x4d, 0x65, 0x64, 0x69,
	0x61, 0x12, 0x12, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x4d, 0x65, 0x64,
	0x69, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x09, 0x52, 0x65,
	0x61, 0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x12, 0x11, 0x2e, 0x52, 0x65, 0x61, 0x64, 0x4d, 0x65,
	0x64, 0x69, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x52, 0x65, 0x61,
	0x64, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x1d,
	0x5a, 0x1b, 0x73, 0x69, 0x67, 0x6d, 0x61, 0x6f, 0x73, 0x2f, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x6c,
	0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_socialnetwork_proto_post_proto_rawDescOnce sync.Once
	file_socialnetwork_proto_post_proto_rawDescData = file_socialnetwork_proto_post_proto_rawDesc
)

func file_socialnetwork_proto_post_proto_rawDescGZIP() []byte {
	file_socialnetwork_proto_post_proto_rawDescOnce.Do(func() {
		file_socialnetwork_proto_post_proto_rawDescData = protoimpl.X.CompressGZIP(file_socialnetwork_proto_post_proto_rawDescData)
	})
	return file_socialnetwork_proto_post_proto_rawDescData
}

var file_socialnetwork_proto_post_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_socialnetwork_proto_post_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_socialnetwork_proto_post_proto_goTypes = []interface{}{
	(POST_TYPE)(0),              // 0: POST_TYPE
	(*StorePostRequest)(nil),    // 1: StorePostRequest
	(*StorePostResponse)(nil),   // 2: StorePostResponse
	(*ReadPostsRequest)(nil),    // 3: ReadPostsRequest
	(*ReadPostsResponse)(nil),   // 4: ReadPostsResponse
	(*StoreMediaRequest)(nil),   // 5: StoreMediaRequest
	(*StoreMediaResponse)(nil),  // 6: StoreMediaResponse
	(*ReadMediaRequest)(nil),    // 7: ReadMediaRequest
	(*ReadMediaResponse)(nil),   // 8: ReadMediaResponse
	(*ComposePostRequest)(nil),  // 9: ComposePostRequest
	(*ComposePostResponse)(nil), // 10: ComposePostResponse
	(*Post)(nil),                // 11: Post
}
var file_socialnetwork_proto_post_proto_depIdxs = []int32{
	11, // 0: StorePostRequest.post:type_name -> Post
	11, // 1: ReadPostsResponse.posts:type_name -> Post
	0,  // 2: ComposePostRequest.posttype:type_name -> POST_TYPE
	0,  // 3: Post.posttype:type_name -> POST_TYPE
	1,  // 4: PostService.StorePost:input_type -> StorePostRequest
	3,  // 5: PostService.ReadPosts:input_type -> ReadPostsRequest
	9,  // 6: ComposeService.ComposePost:input_type -> ComposePostRequest
	5,  // 7: MediaService.StoreMedia:input_type -> StoreMediaRequest
	7,  // 8: MediaService.ReadMedia:input_type -> ReadMediaRequest
	2,  // 9: PostService.StorePost:output_type -> StorePostResponse
	4,  // 10: PostService.ReadPosts:output_type -> ReadPostsResponse
	10, // 11: ComposeService.ComposePost:output_type -> ComposePostResponse
	6,  // 12: MediaService.StoreMedia:output_type -> StoreMediaResponse
	8,  // 13: MediaService.ReadMedia:output_type -> ReadMediaResponse
	9,  // [9:14] is the sub-list for method output_type
	4,  // [4:9] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_socialnetwork_proto_post_proto_init() }
func file_socialnetwork_proto_post_proto_init() {
	if File_socialnetwork_proto_post_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_socialnetwork_proto_post_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StorePostRequest); i {
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
		file_socialnetwork_proto_post_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StorePostResponse); i {
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
		file_socialnetwork_proto_post_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadPostsRequest); i {
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
		file_socialnetwork_proto_post_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadPostsResponse); i {
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
		file_socialnetwork_proto_post_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreMediaRequest); i {
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
		file_socialnetwork_proto_post_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreMediaResponse); i {
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
		file_socialnetwork_proto_post_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadMediaRequest); i {
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
		file_socialnetwork_proto_post_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReadMediaResponse); i {
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
		file_socialnetwork_proto_post_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ComposePostRequest); i {
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
		file_socialnetwork_proto_post_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ComposePostResponse); i {
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
		file_socialnetwork_proto_post_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Post); i {
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
			RawDescriptor: file_socialnetwork_proto_post_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   3,
		},
		GoTypes:           file_socialnetwork_proto_post_proto_goTypes,
		DependencyIndexes: file_socialnetwork_proto_post_proto_depIdxs,
		EnumInfos:         file_socialnetwork_proto_post_proto_enumTypes,
		MessageInfos:      file_socialnetwork_proto_post_proto_msgTypes,
	}.Build()
	File_socialnetwork_proto_post_proto = out.File
	file_socialnetwork_proto_post_proto_rawDesc = nil
	file_socialnetwork_proto_post_proto_goTypes = nil
	file_socialnetwork_proto_post_proto_depIdxs = nil
}
