// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.22.0
// source: prompt.proto

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

type PromptArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InspirationFiles []string      `protobuf:"bytes,1,rep,name=inspirationFiles,proto3" json:"inspirationFiles,omitempty"`
	ExcludeFiles     []string      `protobuf:"bytes,2,rep,name=excludeFiles,proto3" json:"excludeFiles,omitempty"`
	SearchArgss      []*SearchArgs `protobuf:"bytes,4,rep,name=searchArgss,proto3" json:"searchArgss,omitempty"`
	Thread           []*Message    `protobuf:"bytes,6,rep,name=thread,proto3" json:"thread,omitempty"`
}

func (x *PromptArgs) Reset() {
	*x = PromptArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_prompt_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PromptArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PromptArgs) ProtoMessage() {}

func (x *PromptArgs) ProtoReflect() protoreflect.Message {
	mi := &file_prompt_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PromptArgs.ProtoReflect.Descriptor instead.
func (*PromptArgs) Descriptor() ([]byte, []int) {
	return file_prompt_proto_rawDescGZIP(), []int{0}
}

func (x *PromptArgs) GetInspirationFiles() []string {
	if x != nil {
		return x.InspirationFiles
	}
	return nil
}

func (x *PromptArgs) GetExcludeFiles() []string {
	if x != nil {
		return x.ExcludeFiles
	}
	return nil
}

func (x *PromptArgs) GetSearchArgss() []*SearchArgs {
	if x != nil {
		return x.SearchArgss
	}
	return nil
}

func (x *PromptArgs) GetThread() []*Message {
	if x != nil {
		return x.Thread
	}
	return nil
}

type PromptRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PromptArgs *PromptArgs `protobuf:"bytes,1,opt,name=promptArgs,proto3" json:"promptArgs,omitempty"`
}

func (x *PromptRequest) Reset() {
	*x = PromptRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_prompt_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PromptRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PromptRequest) ProtoMessage() {}

func (x *PromptRequest) ProtoReflect() protoreflect.Message {
	mi := &file_prompt_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PromptRequest.ProtoReflect.Descriptor instead.
func (*PromptRequest) Descriptor() ([]byte, []int) {
	return file_prompt_proto_rawDescGZIP(), []int{1}
}

func (x *PromptRequest) GetPromptArgs() *PromptArgs {
	if x != nil {
		return x.PromptArgs
	}
	return nil
}

type PromptResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Thread []*Message `protobuf:"bytes,1,rep,name=thread,proto3" json:"thread,omitempty"`
}

func (x *PromptResponse) Reset() {
	*x = PromptResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_prompt_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PromptResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PromptResponse) ProtoMessage() {}

func (x *PromptResponse) ProtoReflect() protoreflect.Message {
	mi := &file_prompt_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PromptResponse.ProtoReflect.Descriptor instead.
func (*PromptResponse) Descriptor() ([]byte, []int) {
	return file_prompt_proto_rawDescGZIP(), []int{2}
}

func (x *PromptResponse) GetThread() []*Message {
	if x != nil {
		return x.Thread
	}
	return nil
}

var File_prompt_proto protoreflect.FileDescriptor

var file_prompt_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x1a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbf, 0x01, 0x0a, 0x0a, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x41,
	0x72, 0x67, 0x73, 0x12, 0x2a, 0x0a, 0x10, 0x69, 0x6e, 0x73, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x10, 0x69,
	0x6e, 0x73, 0x70, 0x69, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x12,
	0x22, 0x0a, 0x0c, 0x65, 0x78, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x78, 0x63, 0x6c, 0x75, 0x64, 0x65, 0x46, 0x69,
	0x6c, 0x65, 0x73, 0x12, 0x36, 0x0a, 0x0b, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x41, 0x72, 0x67,
	0x73, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x41, 0x72, 0x67, 0x73, 0x52, 0x0b,
	0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x41, 0x72, 0x67, 0x73, 0x73, 0x12, 0x29, 0x0a, 0x06, 0x74,
	0x68, 0x72, 0x65, 0x61, 0x64, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x06,
	0x74, 0x68, 0x72, 0x65, 0x61, 0x64, 0x22, 0x45, 0x0a, 0x0d, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x34, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6d, 0x70,
	0x74, 0x41, 0x72, 0x67, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2e, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x41, 0x72, 0x67,
	0x73, 0x52, 0x0a, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x41, 0x72, 0x67, 0x73, 0x22, 0x3b, 0x0a,
	0x0e, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x29, 0x0a, 0x06, 0x74, 0x68, 0x72, 0x65, 0x61, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x11, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x06, 0x74, 0x68, 0x72, 0x65, 0x61, 0x64, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_prompt_proto_rawDescOnce sync.Once
	file_prompt_proto_rawDescData = file_prompt_proto_rawDesc
)

func file_prompt_proto_rawDescGZIP() []byte {
	file_prompt_proto_rawDescOnce.Do(func() {
		file_prompt_proto_rawDescData = protoimpl.X.CompressGZIP(file_prompt_proto_rawDescData)
	})
	return file_prompt_proto_rawDescData
}

var file_prompt_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_prompt_proto_goTypes = []interface{}{
	(*PromptArgs)(nil),     // 0: actionpb.PromptArgs
	(*PromptRequest)(nil),  // 1: actionpb.PromptRequest
	(*PromptResponse)(nil), // 2: actionpb.PromptResponse
	(*SearchArgs)(nil),     // 3: actionpb.SearchArgs
	(*Message)(nil),        // 4: actionpb.Message
}
var file_prompt_proto_depIdxs = []int32{
	3, // 0: actionpb.PromptArgs.searchArgss:type_name -> actionpb.SearchArgs
	4, // 1: actionpb.PromptArgs.thread:type_name -> actionpb.Message
	0, // 2: actionpb.PromptRequest.promptArgs:type_name -> actionpb.PromptArgs
	4, // 3: actionpb.PromptResponse.thread:type_name -> actionpb.Message
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_prompt_proto_init() }
func file_prompt_proto_init() {
	if File_prompt_proto != nil {
		return
	}
	file_common_proto_init()
	file_search_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_prompt_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PromptArgs); i {
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
		file_prompt_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PromptRequest); i {
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
		file_prompt_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PromptResponse); i {
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
			RawDescriptor: file_prompt_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_prompt_proto_goTypes,
		DependencyIndexes: file_prompt_proto_depIdxs,
		MessageInfos:      file_prompt_proto_msgTypes,
	}.Build()
	File_prompt_proto = out.File
	file_prompt_proto_rawDesc = nil
	file_prompt_proto_goTypes = nil
	file_prompt_proto_depIdxs = nil
}
