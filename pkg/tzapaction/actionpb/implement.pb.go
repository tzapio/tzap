// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.22.0
// source: implement.proto

package actionpb

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

type ImplementRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ImplementArgs *ImplementArgs `protobuf:"bytes,1,opt,name=ImplementArgs,proto3" json:"ImplementArgs,omitempty"`
}

func (x *ImplementRequest) Reset() {
	*x = ImplementRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_implement_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImplementRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImplementRequest) ProtoMessage() {}

func (x *ImplementRequest) ProtoReflect() protoreflect.Message {
	mi := &file_implement_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImplementRequest.ProtoReflect.Descriptor instead.
func (*ImplementRequest) Descriptor() ([]byte, []int) {
	return file_implement_proto_rawDescGZIP(), []int{0}
}

func (x *ImplementRequest) GetImplementArgs() *ImplementArgs {
	if x != nil {
		return x.ImplementArgs
	}
	return nil
}

type ImplementArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Mission string  `protobuf:"bytes,1,opt,name=mission,proto3" json:"mission,omitempty"`
	Plan    string  `protobuf:"bytes,2,opt,name=plan,proto3" json:"plan,omitempty"`
	Tasks   []*Task `protobuf:"bytes,3,rep,name=tasks,proto3" json:"tasks,omitempty"`
}

func (x *ImplementArgs) Reset() {
	*x = ImplementArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_implement_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImplementArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImplementArgs) ProtoMessage() {}

func (x *ImplementArgs) ProtoReflect() protoreflect.Message {
	mi := &file_implement_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImplementArgs.ProtoReflect.Descriptor instead.
func (*ImplementArgs) Descriptor() ([]byte, []int) {
	return file_implement_proto_rawDescGZIP(), []int{1}
}

func (x *ImplementArgs) GetMission() string {
	if x != nil {
		return x.Mission
	}
	return ""
}

func (x *ImplementArgs) GetPlan() string {
	if x != nil {
		return x.Plan
	}
	return ""
}

func (x *ImplementArgs) GetTasks() []*Task {
	if x != nil {
		return x.Tasks
	}
	return nil
}

type Task struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Task             string   `protobuf:"bytes,1,opt,name=task,proto3" json:"task,omitempty"`
	FileIn           string   `protobuf:"bytes,2,opt,name=fileIn,proto3" json:"fileIn,omitempty"`
	FileOut          string   `protobuf:"bytes,3,opt,name=fileOut,proto3" json:"fileOut,omitempty"`
	Code             string   `protobuf:"bytes,4,opt,name=code,proto3" json:"code,omitempty"`
	InspirationFiles []string `protobuf:"bytes,5,rep,name=inspirationFiles,proto3" json:"inspirationFiles,omitempty"`
}

func (x *Task) Reset() {
	*x = Task{}
	if protoimpl.UnsafeEnabled {
		mi := &file_implement_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Task) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Task) ProtoMessage() {}

func (x *Task) ProtoReflect() protoreflect.Message {
	mi := &file_implement_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Task.ProtoReflect.Descriptor instead.
func (*Task) Descriptor() ([]byte, []int) {
	return file_implement_proto_rawDescGZIP(), []int{2}
}

func (x *Task) GetTask() string {
	if x != nil {
		return x.Task
	}
	return ""
}

func (x *Task) GetFileIn() string {
	if x != nil {
		return x.FileIn
	}
	return ""
}

func (x *Task) GetFileOut() string {
	if x != nil {
		return x.FileOut
	}
	return ""
}

func (x *Task) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *Task) GetInspirationFiles() []string {
	if x != nil {
		return x.InspirationFiles
	}
	return nil
}

type ImplementResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileWrites []*FileWrite `protobuf:"bytes,1,rep,name=fileWrites,proto3" json:"fileWrites,omitempty"`
}

func (x *ImplementResponse) Reset() {
	*x = ImplementResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_implement_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImplementResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImplementResponse) ProtoMessage() {}

func (x *ImplementResponse) ProtoReflect() protoreflect.Message {
	mi := &file_implement_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImplementResponse.ProtoReflect.Descriptor instead.
func (*ImplementResponse) Descriptor() ([]byte, []int) {
	return file_implement_proto_rawDescGZIP(), []int{3}
}

func (x *ImplementResponse) GetFileWrites() []*FileWrite {
	if x != nil {
		return x.FileWrites
	}
	return nil
}

var File_implement_proto protoreflect.FileDescriptor

var file_implement_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x08, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x1a, 0x0c, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x51, 0x0a, 0x10, 0x49, 0x6d, 0x70,
	0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3d, 0x0a,
	0x0d, 0x49, 0x6d, 0x70, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x41, 0x72, 0x67, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2e,
	0x49, 0x6d, 0x70, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x41, 0x72, 0x67, 0x73, 0x52, 0x0d, 0x49,
	0x6d, 0x70, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x41, 0x72, 0x67, 0x73, 0x22, 0x63, 0x0a, 0x0d,
	0x49, 0x6d, 0x70, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x41, 0x72, 0x67, 0x73, 0x12, 0x18, 0x0a,
	0x07, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6c, 0x61, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x6c, 0x61, 0x6e, 0x12, 0x24, 0x0a, 0x05, 0x74,
	0x61, 0x73, 0x6b, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x05, 0x74, 0x61, 0x73, 0x6b,
	0x73, 0x22, 0x8c, 0x01, 0x0a, 0x04, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61,
	0x73, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x61, 0x73, 0x6b, 0x12, 0x16,
	0x0a, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x66, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x4f, 0x75,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x66, 0x69, 0x6c, 0x65, 0x4f, 0x75, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x12, 0x2a, 0x0a, 0x10, 0x69, 0x6e, 0x73, 0x70, 0x69, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x10,
	0x69, 0x6e, 0x73, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6c, 0x65, 0x73,
	0x22, 0x48, 0x0a, 0x11, 0x49, 0x6d, 0x70, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x33, 0x0a, 0x0a, 0x66, 0x69, 0x6c, 0x65, 0x57, 0x72, 0x69,
	0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x70, 0x62, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x0a,
	0x66, 0x69, 0x6c, 0x65, 0x57, 0x72, 0x69, 0x74, 0x65, 0x73, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_implement_proto_rawDescOnce sync.Once
	file_implement_proto_rawDescData = file_implement_proto_rawDesc
)

func file_implement_proto_rawDescGZIP() []byte {
	file_implement_proto_rawDescOnce.Do(func() {
		file_implement_proto_rawDescData = protoimpl.X.CompressGZIP(file_implement_proto_rawDescData)
	})
	return file_implement_proto_rawDescData
}

var file_implement_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_implement_proto_goTypes = []interface{}{
	(*ImplementRequest)(nil),  // 0: actionpb.ImplementRequest
	(*ImplementArgs)(nil),     // 1: actionpb.ImplementArgs
	(*Task)(nil),              // 2: actionpb.Task
	(*ImplementResponse)(nil), // 3: actionpb.ImplementResponse
	(*FileWrite)(nil),         // 4: actionpb.FileWrite
}
var file_implement_proto_depIdxs = []int32{
	1, // 0: actionpb.ImplementRequest.ImplementArgs:type_name -> actionpb.ImplementArgs
	2, // 1: actionpb.ImplementArgs.tasks:type_name -> actionpb.Task
	4, // 2: actionpb.ImplementResponse.fileWrites:type_name -> actionpb.FileWrite
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_implement_proto_init() }
func file_implement_proto_init() {
	if File_implement_proto != nil {
		return
	}
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_implement_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImplementRequest); i {
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
		file_implement_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImplementArgs); i {
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
		file_implement_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Task); i {
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
		file_implement_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImplementResponse); i {
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
			RawDescriptor: file_implement_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_implement_proto_goTypes,
		DependencyIndexes: file_implement_proto_depIdxs,
		MessageInfos:      file_implement_proto_msgTypes,
	}.Build()
	File_implement_proto = out.File
	file_implement_proto_rawDesc = nil
	file_implement_proto_goTypes = nil
	file_implement_proto_depIdxs = nil
}
