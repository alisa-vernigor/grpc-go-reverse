// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: chat.proto

package chat

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

type ClientMessageType int32

const (
	ClientMessageType_CONNECT         ClientMessageType = 0
	ClientMessageType_DISCONNECT      ClientMessageType = 1
	ClientMessageType_FINISH_THE_DAY  ClientMessageType = 2
	ClientMessageType_MESSAGE_TO_CHAT ClientMessageType = 3
	ClientMessageType_VOTE_FOR        ClientMessageType = 4
	ClientMessageType_PUBLIC          ClientMessageType = 5
)

// Enum value maps for ClientMessageType.
var (
	ClientMessageType_name = map[int32]string{
		0: "CONNECT",
		1: "DISCONNECT",
		2: "FINISH_THE_DAY",
		3: "MESSAGE_TO_CHAT",
		4: "VOTE_FOR",
		5: "PUBLIC",
	}
	ClientMessageType_value = map[string]int32{
		"CONNECT":         0,
		"DISCONNECT":      1,
		"FINISH_THE_DAY":  2,
		"MESSAGE_TO_CHAT": 3,
		"VOTE_FOR":        4,
		"PUBLIC":          5,
	}
)

func (x ClientMessageType) Enum() *ClientMessageType {
	p := new(ClientMessageType)
	*p = x
	return p
}

func (x ClientMessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ClientMessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_chat_proto_enumTypes[0].Descriptor()
}

func (ClientMessageType) Type() protoreflect.EnumType {
	return &file_chat_proto_enumTypes[0]
}

func (x ClientMessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ClientMessageType.Descriptor instead.
func (ClientMessageType) EnumDescriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{0}
}

type ServerMessageType int32

const (
	ServerMessageType_CONNECTED         ServerMessageType = 0
	ServerMessageType_DISCONNECTED      ServerMessageType = 1
	ServerMessageType_GAME_STARTED      ServerMessageType = 2
	ServerMessageType_ROLE              ServerMessageType = 3
	ServerMessageType_GAME_FINISHED     ServerMessageType = 4
	ServerMessageType_DAY               ServerMessageType = 5
	ServerMessageType_NIGHT             ServerMessageType = 6
	ServerMessageType_MESSAGE_FROM_CHAT ServerMessageType = 7
	ServerMessageType_KILLED            ServerMessageType = 8
	ServerMessageType_DRAW              ServerMessageType = 9
	ServerMessageType_DENIED            ServerMessageType = 10
	ServerMessageType_VOTED_OUT         ServerMessageType = 11
	ServerMessageType_PUBLIC_FROM_COM   ServerMessageType = 12
	ServerMessageType_OK                ServerMessageType = 13
)

// Enum value maps for ServerMessageType.
var (
	ServerMessageType_name = map[int32]string{
		0:  "CONNECTED",
		1:  "DISCONNECTED",
		2:  "GAME_STARTED",
		3:  "ROLE",
		4:  "GAME_FINISHED",
		5:  "DAY",
		6:  "NIGHT",
		7:  "MESSAGE_FROM_CHAT",
		8:  "KILLED",
		9:  "DRAW",
		10: "DENIED",
		11: "VOTED_OUT",
		12: "PUBLIC_FROM_COM",
		13: "OK",
	}
	ServerMessageType_value = map[string]int32{
		"CONNECTED":         0,
		"DISCONNECTED":      1,
		"GAME_STARTED":      2,
		"ROLE":              3,
		"GAME_FINISHED":     4,
		"DAY":               5,
		"NIGHT":             6,
		"MESSAGE_FROM_CHAT": 7,
		"KILLED":            8,
		"DRAW":              9,
		"DENIED":            10,
		"VOTED_OUT":         11,
		"PUBLIC_FROM_COM":   12,
		"OK":                13,
	}
)

func (x ServerMessageType) Enum() *ServerMessageType {
	p := new(ServerMessageType)
	*p = x
	return p
}

func (x ServerMessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ServerMessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_chat_proto_enumTypes[1].Descriptor()
}

func (ServerMessageType) Type() protoreflect.EnumType {
	return &file_chat_proto_enumTypes[1]
}

func (x ServerMessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ServerMessageType.Descriptor instead.
func (ServerMessageType) EnumDescriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{1}
}

type ClientMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type     ClientMessageType `protobuf:"varint,1,opt,name=type,proto3,enum=chat.ClientMessageType" json:"type,omitempty"`
	Body     string            `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
	Nickname string            `protobuf:"bytes,3,opt,name=nickname,proto3" json:"nickname,omitempty"`
	Role     string            `protobuf:"bytes,4,opt,name=role,proto3" json:"role,omitempty"`
}

func (x *ClientMessage) Reset() {
	*x = ClientMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientMessage) ProtoMessage() {}

func (x *ClientMessage) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientMessage.ProtoReflect.Descriptor instead.
func (*ClientMessage) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{0}
}

func (x *ClientMessage) GetType() ClientMessageType {
	if x != nil {
		return x.Type
	}
	return ClientMessageType_CONNECT
}

func (x *ClientMessage) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

func (x *ClientMessage) GetNickname() string {
	if x != nil {
		return x.Nickname
	}
	return ""
}

func (x *ClientMessage) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

type ServerMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type     ServerMessageType `protobuf:"varint,1,opt,name=type,proto3,enum=chat.ServerMessageType" json:"type,omitempty"`
	Body     string            `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
	Nickname string            `protobuf:"bytes,3,opt,name=nickname,proto3" json:"nickname,omitempty"`
	Role     string            `protobuf:"bytes,4,opt,name=role,proto3" json:"role,omitempty"`
	Day      int64             `protobuf:"varint,5,opt,name=day,proto3" json:"day,omitempty"`
}

func (x *ServerMessage) Reset() {
	*x = ServerMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_chat_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerMessage) ProtoMessage() {}

func (x *ServerMessage) ProtoReflect() protoreflect.Message {
	mi := &file_chat_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerMessage.ProtoReflect.Descriptor instead.
func (*ServerMessage) Descriptor() ([]byte, []int) {
	return file_chat_proto_rawDescGZIP(), []int{1}
}

func (x *ServerMessage) GetType() ServerMessageType {
	if x != nil {
		return x.Type
	}
	return ServerMessageType_CONNECTED
}

func (x *ServerMessage) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

func (x *ServerMessage) GetNickname() string {
	if x != nil {
		return x.Nickname
	}
	return ""
}

func (x *ServerMessage) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *ServerMessage) GetDay() int64 {
	if x != nil {
		return x.Day
	}
	return 0
}

var File_chat_proto protoreflect.FileDescriptor

var file_chat_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x63, 0x68,
	0x61, 0x74, 0x22, 0x80, 0x01, 0x0a, 0x0d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x17, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x62, 0x6f, 0x64, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x72, 0x6f, 0x6c, 0x65, 0x22, 0x92, 0x01, 0x0a, 0x0d, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x53, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x69, 0x63, 0x6b,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6e, 0x69, 0x63, 0x6b,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x64, 0x61, 0x79, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x64, 0x61, 0x79, 0x2a, 0x73, 0x0a, 0x11, 0x43, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x0b, 0x0a, 0x07, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a,
	0x44, 0x49, 0x53, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x10, 0x01, 0x12, 0x12, 0x0a, 0x0e,
	0x46, 0x49, 0x4e, 0x49, 0x53, 0x48, 0x5f, 0x54, 0x48, 0x45, 0x5f, 0x44, 0x41, 0x59, 0x10, 0x02,
	0x12, 0x13, 0x0a, 0x0f, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x5f, 0x54, 0x4f, 0x5f, 0x43,
	0x48, 0x41, 0x54, 0x10, 0x03, 0x12, 0x0c, 0x0a, 0x08, 0x56, 0x4f, 0x54, 0x45, 0x5f, 0x46, 0x4f,
	0x52, 0x10, 0x04, 0x12, 0x0a, 0x0a, 0x06, 0x50, 0x55, 0x42, 0x4c, 0x49, 0x43, 0x10, 0x05, 0x2a,
	0xdc, 0x01, 0x0a, 0x11, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0d, 0x0a, 0x09, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54,
	0x45, 0x44, 0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c, 0x44, 0x49, 0x53, 0x43, 0x4f, 0x4e, 0x4e, 0x45,
	0x43, 0x54, 0x45, 0x44, 0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x47, 0x41, 0x4d, 0x45, 0x5f, 0x53,
	0x54, 0x41, 0x52, 0x54, 0x45, 0x44, 0x10, 0x02, 0x12, 0x08, 0x0a, 0x04, 0x52, 0x4f, 0x4c, 0x45,
	0x10, 0x03, 0x12, 0x11, 0x0a, 0x0d, 0x47, 0x41, 0x4d, 0x45, 0x5f, 0x46, 0x49, 0x4e, 0x49, 0x53,
	0x48, 0x45, 0x44, 0x10, 0x04, 0x12, 0x07, 0x0a, 0x03, 0x44, 0x41, 0x59, 0x10, 0x05, 0x12, 0x09,
	0x0a, 0x05, 0x4e, 0x49, 0x47, 0x48, 0x54, 0x10, 0x06, 0x12, 0x15, 0x0a, 0x11, 0x4d, 0x45, 0x53,
	0x53, 0x41, 0x47, 0x45, 0x5f, 0x46, 0x52, 0x4f, 0x4d, 0x5f, 0x43, 0x48, 0x41, 0x54, 0x10, 0x07,
	0x12, 0x0a, 0x0a, 0x06, 0x4b, 0x49, 0x4c, 0x4c, 0x45, 0x44, 0x10, 0x08, 0x12, 0x08, 0x0a, 0x04,
	0x44, 0x52, 0x41, 0x57, 0x10, 0x09, 0x12, 0x0a, 0x0a, 0x06, 0x44, 0x45, 0x4e, 0x49, 0x45, 0x44,
	0x10, 0x0a, 0x12, 0x0d, 0x0a, 0x09, 0x56, 0x4f, 0x54, 0x45, 0x44, 0x5f, 0x4f, 0x55, 0x54, 0x10,
	0x0b, 0x12, 0x13, 0x0a, 0x0f, 0x50, 0x55, 0x42, 0x4c, 0x49, 0x43, 0x5f, 0x46, 0x52, 0x4f, 0x4d,
	0x5f, 0x43, 0x4f, 0x4d, 0x10, 0x0c, 0x12, 0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10, 0x0d, 0x32, 0x42,
	0x0a, 0x08, 0x43, 0x68, 0x61, 0x74, 0x52, 0x6f, 0x6f, 0x6d, 0x12, 0x36, 0x0a, 0x04, 0x43, 0x68,
	0x61, 0x74, 0x12, 0x13, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x13, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x2e, 0x53,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x00, 0x28, 0x01,
	0x30, 0x01, 0x42, 0x3a, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x61, 0x6c, 0x69, 0x73, 0x61, 0x2d, 0x76, 0x65, 0x72, 0x6e, 0x69, 0x67, 0x6f, 0x72, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x2d, 0x67, 0x6f, 0x2d, 0x72, 0x65, 0x76, 0x65, 0x72, 0x73, 0x65, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x68, 0x61, 0x74, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_chat_proto_rawDescOnce sync.Once
	file_chat_proto_rawDescData = file_chat_proto_rawDesc
)

func file_chat_proto_rawDescGZIP() []byte {
	file_chat_proto_rawDescOnce.Do(func() {
		file_chat_proto_rawDescData = protoimpl.X.CompressGZIP(file_chat_proto_rawDescData)
	})
	return file_chat_proto_rawDescData
}

var file_chat_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_chat_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_chat_proto_goTypes = []interface{}{
	(ClientMessageType)(0), // 0: chat.ClientMessageType
	(ServerMessageType)(0), // 1: chat.ServerMessageType
	(*ClientMessage)(nil),  // 2: chat.ClientMessage
	(*ServerMessage)(nil),  // 3: chat.ServerMessage
}
var file_chat_proto_depIdxs = []int32{
	0, // 0: chat.ClientMessage.type:type_name -> chat.ClientMessageType
	1, // 1: chat.ServerMessage.type:type_name -> chat.ServerMessageType
	2, // 2: chat.ChatRoom.Chat:input_type -> chat.ClientMessage
	3, // 3: chat.ChatRoom.Chat:output_type -> chat.ServerMessage
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_chat_proto_init() }
func file_chat_proto_init() {
	if File_chat_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_chat_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientMessage); i {
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
		file_chat_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerMessage); i {
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
			RawDescriptor: file_chat_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_chat_proto_goTypes,
		DependencyIndexes: file_chat_proto_depIdxs,
		EnumInfos:         file_chat_proto_enumTypes,
		MessageInfos:      file_chat_proto_msgTypes,
	}.Build()
	File_chat_proto = out.File
	file_chat_proto_rawDesc = nil
	file_chat_proto_goTypes = nil
	file_chat_proto_depIdxs = nil
}
