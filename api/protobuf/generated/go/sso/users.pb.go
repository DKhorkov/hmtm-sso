// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: sso/users.proto

package sso

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetMeIn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestID   string `protobuf:"bytes,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	AccessToken string `protobuf:"bytes,2,opt,name=accessToken,proto3" json:"accessToken,omitempty"`
}

func (x *GetMeIn) Reset() {
	*x = GetMeIn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sso_users_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMeIn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMeIn) ProtoMessage() {}

func (x *GetMeIn) ProtoReflect() protoreflect.Message {
	mi := &file_sso_users_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMeIn.ProtoReflect.Descriptor instead.
func (*GetMeIn) Descriptor() ([]byte, []int) {
	return file_sso_users_proto_rawDescGZIP(), []int{0}
}

func (x *GetMeIn) GetRequestID() string {
	if x != nil {
		return x.RequestID
	}
	return ""
}

func (x *GetMeIn) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

type GetUserIn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestID string `protobuf:"bytes,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	ID        uint64 `protobuf:"varint,2,opt,name=ID,proto3" json:"ID,omitempty"`
}

func (x *GetUserIn) Reset() {
	*x = GetUserIn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sso_users_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserIn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserIn) ProtoMessage() {}

func (x *GetUserIn) ProtoReflect() protoreflect.Message {
	mi := &file_sso_users_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserIn.ProtoReflect.Descriptor instead.
func (*GetUserIn) Descriptor() ([]byte, []int) {
	return file_sso_users_proto_rawDescGZIP(), []int{1}
}

func (x *GetUserIn) GetRequestID() string {
	if x != nil {
		return x.RequestID
	}
	return ""
}

func (x *GetUserIn) GetID() uint64 {
	if x != nil {
		return x.ID
	}
	return 0
}

type GetUserOut struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID        uint64                 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Email     string                 `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`
}

func (x *GetUserOut) Reset() {
	*x = GetUserOut{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sso_users_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserOut) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserOut) ProtoMessage() {}

func (x *GetUserOut) ProtoReflect() protoreflect.Message {
	mi := &file_sso_users_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserOut.ProtoReflect.Descriptor instead.
func (*GetUserOut) Descriptor() ([]byte, []int) {
	return file_sso_users_proto_rawDescGZIP(), []int{2}
}

func (x *GetUserOut) GetID() uint64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *GetUserOut) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *GetUserOut) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *GetUserOut) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type GetUsersIn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestID string `protobuf:"bytes,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
}

func (x *GetUsersIn) Reset() {
	*x = GetUsersIn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sso_users_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUsersIn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUsersIn) ProtoMessage() {}

func (x *GetUsersIn) ProtoReflect() protoreflect.Message {
	mi := &file_sso_users_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUsersIn.ProtoReflect.Descriptor instead.
func (*GetUsersIn) Descriptor() ([]byte, []int) {
	return file_sso_users_proto_rawDescGZIP(), []int{3}
}

func (x *GetUsersIn) GetRequestID() string {
	if x != nil {
		return x.RequestID
	}
	return ""
}

type GetUsersOut struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Users []*GetUserOut `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
}

func (x *GetUsersOut) Reset() {
	*x = GetUsersOut{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sso_users_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUsersOut) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUsersOut) ProtoMessage() {}

func (x *GetUsersOut) ProtoReflect() protoreflect.Message {
	mi := &file_sso_users_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUsersOut.ProtoReflect.Descriptor instead.
func (*GetUsersOut) Descriptor() ([]byte, []int) {
	return file_sso_users_proto_rawDescGZIP(), []int{4}
}

func (x *GetUsersOut) GetUsers() []*GetUserOut {
	if x != nil {
		return x.Users
	}
	return nil
}

var File_sso_users_proto protoreflect.FileDescriptor

var file_sso_users_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x73, 0x73, 0x6f, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x49, 0x0a, 0x07, 0x47, 0x65, 0x74,
	0x4d, 0x65, 0x49, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49,
	0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x49, 0x44, 0x12, 0x20, 0x0a, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x39, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49,
	0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x44, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x44, 0x12,
	0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x49, 0x44, 0x22,
	0xa6, 0x01, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x75, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x49, 0x44, 0x12, 0x14,
	0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x12, 0x38, 0x0a, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x38,
	0x0a, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x2a, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x55,
	0x73, 0x65, 0x72, 0x73, 0x49, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x49, 0x44, 0x22, 0x36, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73,
	0x4f, 0x75, 0x74, 0x12, 0x27, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x11, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x4f, 0x75, 0x74, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x32, 0xa3, 0x01, 0x0a,
	0x0c, 0x55, 0x73, 0x65, 0x72, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x30, 0x0a,
	0x07, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x12, 0x10, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x73,
	0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x1a, 0x11, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x75, 0x74, 0x22, 0x00, 0x12,
	0x33, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73, 0x12, 0x11, 0x2e, 0x75, 0x73,
	0x65, 0x72, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73, 0x49, 0x6e, 0x1a, 0x12,
	0x2e, 0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x73, 0x4f,
	0x75, 0x74, 0x22, 0x00, 0x12, 0x2c, 0x0a, 0x05, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x12, 0x0e, 0x2e,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x49, 0x6e, 0x1a, 0x11, 0x2e,
	0x75, 0x73, 0x65, 0x72, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x4f, 0x75, 0x74,
	0x22, 0x00, 0x42, 0x33, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x44, 0x4b, 0x68, 0x6f, 0x72, 0x6b, 0x6f, 0x76, 0x2f, 0x68, 0x6d, 0x74, 0x6d, 0x2d, 0x73,
	0x73, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x73, 0x73, 0x6f, 0x3b, 0x73, 0x73, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_sso_users_proto_rawDescOnce sync.Once
	file_sso_users_proto_rawDescData = file_sso_users_proto_rawDesc
)

func file_sso_users_proto_rawDescGZIP() []byte {
	file_sso_users_proto_rawDescOnce.Do(func() {
		file_sso_users_proto_rawDescData = protoimpl.X.CompressGZIP(file_sso_users_proto_rawDescData)
	})
	return file_sso_users_proto_rawDescData
}

var file_sso_users_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_sso_users_proto_goTypes = []interface{}{
	(*GetMeIn)(nil),               // 0: users.GetMeIn
	(*GetUserIn)(nil),             // 1: users.GetUserIn
	(*GetUserOut)(nil),            // 2: users.GetUserOut
	(*GetUsersIn)(nil),            // 3: users.GetUsersIn
	(*GetUsersOut)(nil),           // 4: users.GetUsersOut
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
}
var file_sso_users_proto_depIdxs = []int32{
	5, // 0: users.GetUserOut.createdAt:type_name -> google.protobuf.Timestamp
	5, // 1: users.GetUserOut.updatedAt:type_name -> google.protobuf.Timestamp
	2, // 2: users.GetUsersOut.users:type_name -> users.GetUserOut
	1, // 3: users.UsersService.GetUser:input_type -> users.GetUserIn
	3, // 4: users.UsersService.GetUsers:input_type -> users.GetUsersIn
	0, // 5: users.UsersService.GetMe:input_type -> users.GetMeIn
	2, // 6: users.UsersService.GetUser:output_type -> users.GetUserOut
	4, // 7: users.UsersService.GetUsers:output_type -> users.GetUsersOut
	2, // 8: users.UsersService.GetMe:output_type -> users.GetUserOut
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_sso_users_proto_init() }
func file_sso_users_proto_init() {
	if File_sso_users_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_sso_users_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMeIn); i {
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
		file_sso_users_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserIn); i {
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
		file_sso_users_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserOut); i {
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
		file_sso_users_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUsersIn); i {
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
		file_sso_users_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUsersOut); i {
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
			RawDescriptor: file_sso_users_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_sso_users_proto_goTypes,
		DependencyIndexes: file_sso_users_proto_depIdxs,
		MessageInfos:      file_sso_users_proto_msgTypes,
	}.Build()
	File_sso_users_proto = out.File
	file_sso_users_proto_rawDesc = nil
	file_sso_users_proto_goTypes = nil
	file_sso_users_proto_depIdxs = nil
}
