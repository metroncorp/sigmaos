// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.3
// source: proc/proc.proto

package proc

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
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

type ProcEnvProto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PidStr          string                 `protobuf:"bytes,1,opt,name=pidStr,proto3" json:"pidStr,omitempty"`
	Program         string                 `protobuf:"bytes,2,opt,name=program,proto3" json:"program,omitempty"`
	RealmStr        string                 `protobuf:"bytes,3,opt,name=realmStr,proto3" json:"realmStr,omitempty"`
	UnameStr        string                 `protobuf:"bytes,4,opt,name=unameStr,proto3" json:"unameStr,omitempty"`
	ProcDir         string                 `protobuf:"bytes,5,opt,name=procDir,proto3" json:"procDir,omitempty"`
	ParentDir       string                 `protobuf:"bytes,6,opt,name=parentDir,proto3" json:"parentDir,omitempty"`
	EtcdIP          string                 `protobuf:"bytes,7,opt,name=etcdIP,proto3" json:"etcdIP,omitempty"`
	LocalIPStr      string                 `protobuf:"bytes,8,opt,name=localIPStr,proto3" json:"localIPStr,omitempty"`
	KernelID        string                 `protobuf:"bytes,9,opt,name=kernelID,proto3" json:"kernelID,omitempty"`
	BuildTag        string                 `protobuf:"bytes,10,opt,name=buildTag,proto3" json:"buildTag,omitempty"`
	Net             string                 `protobuf:"bytes,11,opt,name=net,proto3" json:"net,omitempty"`
	Perf            string                 `protobuf:"bytes,12,opt,name=perf,proto3" json:"perf,omitempty"`
	Debug           string                 `protobuf:"bytes,13,opt,name=debug,proto3" json:"debug,omitempty"`
	UprocdPIDStr    string                 `protobuf:"bytes,14,opt,name=uprocdPIDStr,proto3" json:"uprocdPIDStr,omitempty"`
	Privileged      bool                   `protobuf:"varint,15,opt,name=privileged,proto3" json:"privileged,omitempty"`
	Overlays        bool                   `protobuf:"varint,16,opt,name=overlays,proto3" json:"overlays,omitempty"`
	NetFail         int64                  `protobuf:"varint,17,opt,name=netFail,proto3" json:"netFail,omitempty"`
	Crash           int64                  `protobuf:"varint,18,opt,name=crash,proto3" json:"crash,omitempty"`
	Partition       int64                  `protobuf:"varint,19,opt,name=partition,proto3" json:"partition,omitempty"`
	HowInt          int32                  `protobuf:"varint,20,opt,name=howInt,proto3" json:"howInt,omitempty"`
	SpawnTimePB     *timestamppb.Timestamp `protobuf:"bytes,21,opt,name=spawnTimePB,proto3" json:"spawnTimePB,omitempty"`
	Strace          string                 `protobuf:"bytes,22,opt,name=strace,proto3" json:"strace,omitempty"`
	ScheddAddr      *sigmap.Taddr          `protobuf:"bytes,23,opt,name=scheddAddr,proto3" json:"scheddAddr,omitempty"`
	NamedMountProto *sigmap.TmountProto    `protobuf:"bytes,24,opt,name=namedMountProto,proto3" json:"namedMountProto,omitempty"`
	UseSigmaclntd   bool                   `protobuf:"varint,25,opt,name=useSigmaclntd,proto3" json:"useSigmaclntd,omitempty"`
}

func (x *ProcEnvProto) Reset() {
	*x = ProcEnvProto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proc_proc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcEnvProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcEnvProto) ProtoMessage() {}

func (x *ProcEnvProto) ProtoReflect() protoreflect.Message {
	mi := &file_proc_proc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcEnvProto.ProtoReflect.Descriptor instead.
func (*ProcEnvProto) Descriptor() ([]byte, []int) {
	return file_proc_proc_proto_rawDescGZIP(), []int{0}
}

func (x *ProcEnvProto) GetPidStr() string {
	if x != nil {
		return x.PidStr
	}
	return ""
}

func (x *ProcEnvProto) GetProgram() string {
	if x != nil {
		return x.Program
	}
	return ""
}

func (x *ProcEnvProto) GetRealmStr() string {
	if x != nil {
		return x.RealmStr
	}
	return ""
}

func (x *ProcEnvProto) GetUnameStr() string {
	if x != nil {
		return x.UnameStr
	}
	return ""
}

func (x *ProcEnvProto) GetProcDir() string {
	if x != nil {
		return x.ProcDir
	}
	return ""
}

func (x *ProcEnvProto) GetParentDir() string {
	if x != nil {
		return x.ParentDir
	}
	return ""
}

func (x *ProcEnvProto) GetEtcdIP() string {
	if x != nil {
		return x.EtcdIP
	}
	return ""
}

func (x *ProcEnvProto) GetLocalIPStr() string {
	if x != nil {
		return x.LocalIPStr
	}
	return ""
}

func (x *ProcEnvProto) GetKernelID() string {
	if x != nil {
		return x.KernelID
	}
	return ""
}

func (x *ProcEnvProto) GetBuildTag() string {
	if x != nil {
		return x.BuildTag
	}
	return ""
}

func (x *ProcEnvProto) GetNet() string {
	if x != nil {
		return x.Net
	}
	return ""
}

func (x *ProcEnvProto) GetPerf() string {
	if x != nil {
		return x.Perf
	}
	return ""
}

func (x *ProcEnvProto) GetDebug() string {
	if x != nil {
		return x.Debug
	}
	return ""
}

func (x *ProcEnvProto) GetUprocdPIDStr() string {
	if x != nil {
		return x.UprocdPIDStr
	}
	return ""
}

func (x *ProcEnvProto) GetPrivileged() bool {
	if x != nil {
		return x.Privileged
	}
	return false
}

func (x *ProcEnvProto) GetOverlays() bool {
	if x != nil {
		return x.Overlays
	}
	return false
}

func (x *ProcEnvProto) GetNetFail() int64 {
	if x != nil {
		return x.NetFail
	}
	return 0
}

func (x *ProcEnvProto) GetCrash() int64 {
	if x != nil {
		return x.Crash
	}
	return 0
}

func (x *ProcEnvProto) GetPartition() int64 {
	if x != nil {
		return x.Partition
	}
	return 0
}

func (x *ProcEnvProto) GetHowInt() int32 {
	if x != nil {
		return x.HowInt
	}
	return 0
}

func (x *ProcEnvProto) GetSpawnTimePB() *timestamppb.Timestamp {
	if x != nil {
		return x.SpawnTimePB
	}
	return nil
}

func (x *ProcEnvProto) GetStrace() string {
	if x != nil {
		return x.Strace
	}
	return ""
}

func (x *ProcEnvProto) GetScheddAddr() *sigmap.Taddr {
	if x != nil {
		return x.ScheddAddr
	}
	return nil
}

func (x *ProcEnvProto) GetNamedMountProto() *sigmap.TmountProto {
	if x != nil {
		return x.NamedMountProto
	}
	return nil
}

func (x *ProcEnvProto) GetUseSigmaclntd() bool {
	if x != nil {
		return x.UseSigmaclntd
	}
	return false
}

type ProcProto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProcEnvProto *ProcEnvProto     `protobuf:"bytes,1,opt,name=procEnvProto,proto3" json:"procEnvProto,omitempty"`
	Args         []string          `protobuf:"bytes,6,rep,name=args,proto3" json:"args,omitempty"`
	Env          map[string]string `protobuf:"bytes,7,rep,name=env,proto3" json:"env,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	TypeInt      uint32            `protobuf:"varint,8,opt,name=typeInt,proto3" json:"typeInt,omitempty"`
	McpuInt      uint32            `protobuf:"varint,9,opt,name=mcpuInt,proto3" json:"mcpuInt,omitempty"`
	MemInt       uint32            `protobuf:"varint,10,opt,name=memInt,proto3" json:"memInt,omitempty"`
	SharedTarget string            `protobuf:"bytes,13,opt,name=sharedTarget,proto3" json:"sharedTarget,omitempty"`
}

func (x *ProcProto) Reset() {
	*x = ProcProto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proc_proc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcProto) ProtoMessage() {}

func (x *ProcProto) ProtoReflect() protoreflect.Message {
	mi := &file_proc_proc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcProto.ProtoReflect.Descriptor instead.
func (*ProcProto) Descriptor() ([]byte, []int) {
	return file_proc_proc_proto_rawDescGZIP(), []int{1}
}

func (x *ProcProto) GetProcEnvProto() *ProcEnvProto {
	if x != nil {
		return x.ProcEnvProto
	}
	return nil
}

func (x *ProcProto) GetArgs() []string {
	if x != nil {
		return x.Args
	}
	return nil
}

func (x *ProcProto) GetEnv() map[string]string {
	if x != nil {
		return x.Env
	}
	return nil
}

func (x *ProcProto) GetTypeInt() uint32 {
	if x != nil {
		return x.TypeInt
	}
	return 0
}

func (x *ProcProto) GetMcpuInt() uint32 {
	if x != nil {
		return x.McpuInt
	}
	return 0
}

func (x *ProcProto) GetMemInt() uint32 {
	if x != nil {
		return x.MemInt
	}
	return 0
}

func (x *ProcProto) GetSharedTarget() string {
	if x != nil {
		return x.SharedTarget
	}
	return ""
}

var File_proc_proc_proto protoreflect.FileDescriptor

var file_proc_proc_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x70, 0x72, 0x6f, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x13, 0x73, 0x69, 0x67, 0x6d, 0x61, 0x70, 0x2f, 0x73, 0x69, 0x67, 0x6d, 0x61,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xfe, 0x05, 0x0a, 0x0c, 0x50, 0x72, 0x6f, 0x63,
	0x45, 0x6e, 0x76, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x69, 0x64, 0x53,
	0x74, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x69, 0x64, 0x53, 0x74, 0x72,
	0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65,
	0x61, 0x6c, 0x6d, 0x53, 0x74, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65,
	0x61, 0x6c, 0x6d, 0x53, 0x74, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x6e, 0x61, 0x6d, 0x65, 0x53,
	0x74, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x6e, 0x61, 0x6d, 0x65, 0x53,
	0x74, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x63, 0x44, 0x69, 0x72, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x63, 0x44, 0x69, 0x72, 0x12, 0x1c, 0x0a, 0x09,
	0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x44, 0x69, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x44, 0x69, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x65, 0x74,
	0x63, 0x64, 0x49, 0x50, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x65, 0x74, 0x63, 0x64,
	0x49, 0x50, 0x12, 0x1e, 0x0a, 0x0a, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x49, 0x50, 0x53, 0x74, 0x72,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x49, 0x50, 0x53,
	0x74, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x49, 0x44, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6b, 0x65, 0x72, 0x6e, 0x65, 0x6c, 0x49, 0x44, 0x12, 0x1a,
	0x0a, 0x08, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x54, 0x61, 0x67, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x54, 0x61, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x6e, 0x65,
	0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6e, 0x65, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x70, 0x65, 0x72, 0x66, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x65, 0x72, 0x66,
	0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x62, 0x75, 0x67, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x64, 0x65, 0x62, 0x75, 0x67, 0x12, 0x22, 0x0a, 0x0c, 0x75, 0x70, 0x72, 0x6f, 0x63, 0x64,
	0x50, 0x49, 0x44, 0x53, 0x74, 0x72, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x75, 0x70,
	0x72, 0x6f, 0x63, 0x64, 0x50, 0x49, 0x44, 0x53, 0x74, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x72,
	0x69, 0x76, 0x69, 0x6c, 0x65, 0x67, 0x65, 0x64, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a,
	0x70, 0x72, 0x69, 0x76, 0x69, 0x6c, 0x65, 0x67, 0x65, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x76,
	0x65, 0x72, 0x6c, 0x61, 0x79, 0x73, 0x18, 0x10, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x6f, 0x76,
	0x65, 0x72, 0x6c, 0x61, 0x79, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6e, 0x65, 0x74, 0x46, 0x61, 0x69,
	0x6c, 0x18, 0x11, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x6e, 0x65, 0x74, 0x46, 0x61, 0x69, 0x6c,
	0x12, 0x14, 0x0a, 0x05, 0x63, 0x72, 0x61, 0x73, 0x68, 0x18, 0x12, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x05, 0x63, 0x72, 0x61, 0x73, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x61, 0x72, 0x74, 0x69, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x13, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x70, 0x61, 0x72, 0x74, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x6f, 0x77, 0x49, 0x6e, 0x74, 0x18, 0x14,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x68, 0x6f, 0x77, 0x49, 0x6e, 0x74, 0x12, 0x3c, 0x0a, 0x0b,
	0x73, 0x70, 0x61, 0x77, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x50, 0x42, 0x18, 0x15, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0b, 0x73,
	0x70, 0x61, 0x77, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x50, 0x42, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74,
	0x72, 0x61, 0x63, 0x65, 0x18, 0x16, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x72, 0x61,
	0x63, 0x65, 0x12, 0x26, 0x0a, 0x0a, 0x73, 0x63, 0x68, 0x65, 0x64, 0x64, 0x41, 0x64, 0x64, 0x72,
	0x18, 0x17, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x54, 0x61, 0x64, 0x64, 0x72, 0x52, 0x0a,
	0x73, 0x63, 0x68, 0x65, 0x64, 0x64, 0x41, 0x64, 0x64, 0x72, 0x12, 0x36, 0x0a, 0x0f, 0x6e, 0x61,
	0x6d, 0x65, 0x64, 0x4d, 0x6f, 0x75, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x18, 0x18, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x54, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x52, 0x0f, 0x6e, 0x61, 0x6d, 0x65, 0x64, 0x4d, 0x6f, 0x75, 0x6e, 0x74, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x24, 0x0a, 0x0d, 0x75, 0x73, 0x65, 0x53, 0x69, 0x67, 0x6d, 0x61, 0x63, 0x6c,
	0x6e, 0x74, 0x64, 0x18, 0x19, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x75, 0x73, 0x65, 0x53, 0x69,
	0x67, 0x6d, 0x61, 0x63, 0x6c, 0x6e, 0x74, 0x64, 0x22, 0xa1, 0x02, 0x0a, 0x09, 0x50, 0x72, 0x6f,
	0x63, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x31, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x63, 0x45, 0x6e,
	0x76, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x50,
	0x72, 0x6f, 0x63, 0x45, 0x6e, 0x76, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52, 0x0c, 0x70, 0x72, 0x6f,
	0x63, 0x45, 0x6e, 0x76, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x61, 0x72, 0x67,
	0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x61, 0x72, 0x67, 0x73, 0x12, 0x25, 0x0a,
	0x03, 0x65, 0x6e, 0x76, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x50, 0x72, 0x6f,
	0x63, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x45, 0x6e, 0x76, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x03, 0x65, 0x6e, 0x76, 0x12, 0x18, 0x0a, 0x07, 0x74, 0x79, 0x70, 0x65, 0x49, 0x6e, 0x74, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x74, 0x79, 0x70, 0x65, 0x49, 0x6e, 0x74, 0x12, 0x18,
	0x0a, 0x07, 0x6d, 0x63, 0x70, 0x75, 0x49, 0x6e, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x07, 0x6d, 0x63, 0x70, 0x75, 0x49, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x6d, 0x49,
	0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x6d, 0x65, 0x6d, 0x49, 0x6e, 0x74,
	0x12, 0x22, 0x0a, 0x0c, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x54, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x54, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x1a, 0x36, 0x0a, 0x08, 0x45, 0x6e, 0x76, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x0e, 0x5a, 0x0c,
	0x73, 0x69, 0x67, 0x6d, 0x61, 0x6f, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x63, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proc_proc_proto_rawDescOnce sync.Once
	file_proc_proc_proto_rawDescData = file_proc_proc_proto_rawDesc
)

func file_proc_proc_proto_rawDescGZIP() []byte {
	file_proc_proc_proto_rawDescOnce.Do(func() {
		file_proc_proc_proto_rawDescData = protoimpl.X.CompressGZIP(file_proc_proc_proto_rawDescData)
	})
	return file_proc_proc_proto_rawDescData
}

var file_proc_proc_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proc_proc_proto_goTypes = []interface{}{
	(*ProcEnvProto)(nil),          // 0: ProcEnvProto
	(*ProcProto)(nil),             // 1: ProcProto
	nil,                           // 2: ProcProto.EnvEntry
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
	(*sigmap.Taddr)(nil),          // 4: Taddr
	(*sigmap.TmountProto)(nil),    // 5: TmountProto
}
var file_proc_proc_proto_depIdxs = []int32{
	3, // 0: ProcEnvProto.spawnTimePB:type_name -> google.protobuf.Timestamp
	4, // 1: ProcEnvProto.scheddAddr:type_name -> Taddr
	5, // 2: ProcEnvProto.namedMountProto:type_name -> TmountProto
	0, // 3: ProcProto.procEnvProto:type_name -> ProcEnvProto
	2, // 4: ProcProto.env:type_name -> ProcProto.EnvEntry
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_proc_proc_proto_init() }
func file_proc_proc_proto_init() {
	if File_proc_proc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proc_proc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcEnvProto); i {
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
		file_proc_proc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcProto); i {
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
			RawDescriptor: file_proc_proc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proc_proc_proto_goTypes,
		DependencyIndexes: file_proc_proc_proto_depIdxs,
		MessageInfos:      file_proc_proc_proto_msgTypes,
	}.Build()
	File_proc_proc_proto = out.File
	file_proc_proc_proto_rawDesc = nil
	file_proc_proc_proto_goTypes = nil
	file_proc_proc_proto_depIdxs = nil
}
