// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.22.0
// source: refactor.proto

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

type RefactorArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InspirationFiles []string `protobuf:"bytes,1,rep,name=inspirationFiles,proto3" json:"inspirationFiles,omitempty"`
	FileIn           string   `protobuf:"bytes,2,opt,name=fileIn,proto3" json:"fileIn,omitempty"`
	FileOut          string   `protobuf:"bytes,3,opt,name=fileOut,proto3" json:"fileOut,omitempty"`
	Mission          string   `protobuf:"bytes,4,opt,name=mission,proto3" json:"mission,omitempty"`
	Task             string   `protobuf:"bytes,5,opt,name=task,proto3" json:"task,omitempty"`
	Plan             string   `protobuf:"bytes,6,opt,name=plan,proto3" json:"plan,omitempty"`
	OutputFormat     string   `protobuf:"bytes,7,opt,name=outputFormat,proto3" json:"outputFormat,omitempty"`
	Example          string   `protobuf:"bytes,8,opt,name=example,proto3" json:"example,omitempty"`
}

func (x *RefactorArgs) Reset() {
	*x = RefactorArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_refactor_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RefactorArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefactorArgs) ProtoMessage() {}

func (x *RefactorArgs) ProtoReflect() protoreflect.Message {
	mi := &file_refactor_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefactorArgs.ProtoReflect.Descriptor instead.
func (*RefactorArgs) Descriptor() ([]byte, []int) {
	return file_refactor_proto_rawDescGZIP(), []int{0}
}

func (x *RefactorArgs) GetInspirationFiles() []string {
	if x != nil {
		return x.InspirationFiles
	}
	return nil
}

func (x *RefactorArgs) GetFileIn() string {
	if x != nil {
		return x.FileIn
	}
	return ""
}

func (x *RefactorArgs) GetFileOut() string {
	if x != nil {
		return x.FileOut
	}
	return ""
}

func (x *RefactorArgs) GetMission() string {
	if x != nil {
		return x.Mission
	}
	return ""
}

func (x *RefactorArgs) GetTask() string {
	if x != nil {
		return x.Task
	}
	return ""
}

func (x *RefactorArgs) GetPlan() string {
	if x != nil {
		return x.Plan
	}
	return ""
}

func (x *RefactorArgs) GetOutputFormat() string {
	if x != nil {
		return x.OutputFormat
	}
	return ""
}

func (x *RefactorArgs) GetExample() string {
	if x != nil {
		return x.Example
	}
	return ""
}

type RefactorRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RefactorArgs *RefactorArgs `protobuf:"bytes,1,opt,name=refactorArgs,proto3" json:"refactorArgs,omitempty"`
}

func (x *RefactorRequest) Reset() {
	*x = RefactorRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_refactor_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RefactorRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefactorRequest) ProtoMessage() {}

func (x *RefactorRequest) ProtoReflect() protoreflect.Message {
	mi := &file_refactor_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefactorRequest.ProtoReflect.Descriptor instead.
func (*RefactorRequest) Descriptor() ([]byte, []int) {
	return file_refactor_proto_rawDescGZIP(), []int{1}
}

func (x *RefactorRequest) GetRefactorArgs() *RefactorArgs {
	if x != nil {
		return x.RefactorArgs
	}
	return nil
}

type RefactorResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileWrites []*FileWrite `protobuf:"bytes,1,rep,name=fileWrites,proto3" json:"fileWrites,omitempty"`
}

func (x *RefactorResponse) Reset() {
	*x = RefactorResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_refactor_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RefactorResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefactorResponse) ProtoMessage() {}

func (x *RefactorResponse) ProtoReflect() protoreflect.Message {
	mi := &file_refactor_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefactorResponse.ProtoReflect.Descriptor instead.
func (*RefactorResponse) Descriptor() ([]byte, []int) {
	return file_refactor_proto_rawDescGZIP(), []int{2}
}

func (x *RefactorResponse) GetFileWrites() []*FileWrite {
	if x != nil {
		return x.FileWrites
	}
	return nil
}

var File_refactor_proto protoreflect.FileDescriptor

var file_refactor_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x72, 0x65, 0x66, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x1a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xec, 0x01, 0x0a, 0x0c, 0x52, 0x65, 0x66,
	0x61, 0x63, 0x74, 0x6f, 0x72, 0x41, 0x72, 0x67, 0x73, 0x12, 0x2a, 0x0a, 0x10, 0x69, 0x6e, 0x73,
	0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x10, 0x69, 0x6e, 0x73, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x46, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x12, 0x18, 0x0a,
	0x07, 0x66, 0x69, 0x6c, 0x65, 0x4f, 0x75, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x66, 0x69, 0x6c, 0x65, 0x4f, 0x75, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x69, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x73, 0x6b, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x74, 0x61, 0x73, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6c, 0x61, 0x6e, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x6c, 0x61, 0x6e, 0x12, 0x22, 0x0a, 0x0c, 0x6f, 0x75, 0x74,
	0x70, 0x75, 0x74, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x18, 0x0a,
	0x07, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x22, 0x4d, 0x0a, 0x0f, 0x52, 0x65, 0x66, 0x61, 0x63,
	0x74, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3a, 0x0a, 0x0c, 0x72, 0x65,
	0x66, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x41, 0x72, 0x67, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x16, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x66, 0x61,
	0x63, 0x74, 0x6f, 0x72, 0x41, 0x72, 0x67, 0x73, 0x52, 0x0c, 0x72, 0x65, 0x66, 0x61, 0x63, 0x74,
	0x6f, 0x72, 0x41, 0x72, 0x67, 0x73, 0x22, 0x47, 0x0a, 0x10, 0x52, 0x65, 0x66, 0x61, 0x63, 0x74,
	0x6f, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x33, 0x0a, 0x0a, 0x66, 0x69,
	0x6c, 0x65, 0x57, 0x72, 0x69, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13,
	0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x57, 0x72,
	0x69, 0x74, 0x65, 0x52, 0x0a, 0x66, 0x69, 0x6c, 0x65, 0x57, 0x72, 0x69, 0x74, 0x65, 0x73, 0x42,
	0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_refactor_proto_rawDescOnce sync.Once
	file_refactor_proto_rawDescData = file_refactor_proto_rawDesc
)

func file_refactor_proto_rawDescGZIP() []byte {
	file_refactor_proto_rawDescOnce.Do(func() {
		file_refactor_proto_rawDescData = protoimpl.X.CompressGZIP(file_refactor_proto_rawDescData)
	})
	return file_refactor_proto_rawDescData
}

var file_refactor_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_refactor_proto_goTypes = []interface{}{
	(*RefactorArgs)(nil),     // 0: actionpb.RefactorArgs
	(*RefactorRequest)(nil),  // 1: actionpb.RefactorRequest
	(*RefactorResponse)(nil), // 2: actionpb.RefactorResponse
	(*FileWrite)(nil),        // 3: actionpb.FileWrite
}
var file_refactor_proto_depIdxs = []int32{
	0, // 0: actionpb.RefactorRequest.refactorArgs:type_name -> actionpb.RefactorArgs
	3, // 1: actionpb.RefactorResponse.fileWrites:type_name -> actionpb.FileWrite
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_refactor_proto_init() }
func file_refactor_proto_init() {
	if File_refactor_proto != nil {
		return
	}
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_refactor_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RefactorArgs); i {
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
		file_refactor_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RefactorRequest); i {
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
		file_refactor_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RefactorResponse); i {
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
			RawDescriptor: file_refactor_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_refactor_proto_goTypes,
		DependencyIndexes: file_refactor_proto_depIdxs,
		MessageInfos:      file_refactor_proto_msgTypes,
	}.Build()
	File_refactor_proto = out.File
	file_refactor_proto_rawDesc = nil
	file_refactor_proto_goTypes = nil
	file_refactor_proto_depIdxs = nil
}
