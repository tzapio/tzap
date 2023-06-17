// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.22.0
// source: tzap.proto

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

type TzapRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Role    string `protobuf:"bytes,1,opt,name=role,proto3" json:"role,omitempty"`
	Content string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *TzapRequest) Reset() {
	*x = TzapRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tzap_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TzapRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TzapRequest) ProtoMessage() {}

func (x *TzapRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tzap_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TzapRequest.ProtoReflect.Descriptor instead.
func (*TzapRequest) Descriptor() ([]byte, []int) {
	return file_tzap_proto_rawDescGZIP(), []int{0}
}

func (x *TzapRequest) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *TzapRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Role    string `protobuf:"bytes,1,opt,name=role,proto3" json:"role,omitempty"`
	Content string `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tzap_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_tzap_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_tzap_proto_rawDescGZIP(), []int{1}
}

func (x *Message) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *Message) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type TzapResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Thread []*Message `protobuf:"bytes,1,rep,name=thread,proto3" json:"thread,omitempty"`
}

func (x *TzapResponse) Reset() {
	*x = TzapResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tzap_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TzapResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TzapResponse) ProtoMessage() {}

func (x *TzapResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tzap_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TzapResponse.ProtoReflect.Descriptor instead.
func (*TzapResponse) Descriptor() ([]byte, []int) {
	return file_tzap_proto_rawDescGZIP(), []int{2}
}

func (x *TzapResponse) GetThread() []*Message {
	if x != nil {
		return x.Thread
	}
	return nil
}

var File_tzap_proto protoreflect.FileDescriptor

var file_tzap_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x74, 0x7a, 0x61, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x1a, 0x0c, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x70, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x3b, 0x0a, 0x0b, 0x54, 0x7a, 0x61, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22,
	0x37, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f,
	0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x39, 0x0a, 0x0c, 0x54, 0x7a, 0x61, 0x70,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x06, 0x74, 0x68, 0x72, 0x65,
	0x61, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x06, 0x74, 0x68, 0x72,
	0x65, 0x61, 0x64, 0x32, 0xbf, 0x01, 0x0a, 0x0b, 0x54, 0x7a, 0x61, 0x70, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x3b, 0x0a, 0x06, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x12, 0x17, 0x2e,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x70,
	0x62, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x39, 0x0a, 0x06, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x12, 0x17, 0x2e, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2e, 0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2e, 0x54,
	0x7a, 0x61, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x07, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x15, 0x2e, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x70,
	0x62, 0x2e, 0x54, 0x7a, 0x61, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x70, 0x62, 0x2e, 0x54, 0x7a, 0x61, 0x70, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tzap_proto_rawDescOnce sync.Once
	file_tzap_proto_rawDescData = file_tzap_proto_rawDesc
)

func file_tzap_proto_rawDescGZIP() []byte {
	file_tzap_proto_rawDescOnce.Do(func() {
		file_tzap_proto_rawDescData = protoimpl.X.CompressGZIP(file_tzap_proto_rawDescData)
	})
	return file_tzap_proto_rawDescData
}

var file_tzap_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_tzap_proto_goTypes = []interface{}{
	(*TzapRequest)(nil),    // 0: actionpb.TzapRequest
	(*Message)(nil),        // 1: actionpb.Message
	(*TzapResponse)(nil),   // 2: actionpb.TzapResponse
	(*SearchRequest)(nil),  // 3: actionpb.SearchRequest
	(*PromptRequest)(nil),  // 4: actionpb.PromptRequest
	(*SearchResponse)(nil), // 5: actionpb.SearchResponse
}
var file_tzap_proto_depIdxs = []int32{
	1, // 0: actionpb.TzapResponse.thread:type_name -> actionpb.Message
	3, // 1: actionpb.TzapService.Search:input_type -> actionpb.SearchRequest
	4, // 2: actionpb.TzapService.Prompt:input_type -> actionpb.PromptRequest
	0, // 3: actionpb.TzapService.Request:input_type -> actionpb.TzapRequest
	5, // 4: actionpb.TzapService.Search:output_type -> actionpb.SearchResponse
	2, // 5: actionpb.TzapService.Prompt:output_type -> actionpb.TzapResponse
	2, // 6: actionpb.TzapService.Request:output_type -> actionpb.TzapResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_tzap_proto_init() }
func file_tzap_proto_init() {
	if File_tzap_proto != nil {
		return
	}
	file_search_proto_init()
	file_prompt_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_tzap_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TzapRequest); i {
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
		file_tzap_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_tzap_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TzapResponse); i {
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
			RawDescriptor: file_tzap_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tzap_proto_goTypes,
		DependencyIndexes: file_tzap_proto_depIdxs,
		MessageInfos:      file_tzap_proto_msgTypes,
	}.Build()
	File_tzap_proto = out.File
	file_tzap_proto_rawDesc = nil
	file_tzap_proto_goTypes = nil
	file_tzap_proto_depIdxs = nil
}
