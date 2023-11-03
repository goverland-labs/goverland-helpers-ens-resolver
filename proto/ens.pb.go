// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.3
// source: ens.proto

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

type ResolveAddressesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Domains []string `protobuf:"bytes,1,rep,name=domains,proto3" json:"domains,omitempty"`
}

func (x *ResolveAddressesRequest) Reset() {
	*x = ResolveAddressesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ens_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResolveAddressesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResolveAddressesRequest) ProtoMessage() {}

func (x *ResolveAddressesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ens_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResolveAddressesRequest.ProtoReflect.Descriptor instead.
func (*ResolveAddressesRequest) Descriptor() ([]byte, []int) {
	return file_ens_proto_rawDescGZIP(), []int{0}
}

func (x *ResolveAddressesRequest) GetDomains() []string {
	if x != nil {
		return x.Domains
	}
	return nil
}

type ResolveDomainsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addresses []string `protobuf:"bytes,1,rep,name=addresses,proto3" json:"addresses,omitempty"`
}

func (x *ResolveDomainsRequest) Reset() {
	*x = ResolveDomainsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ens_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResolveDomainsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResolveDomainsRequest) ProtoMessage() {}

func (x *ResolveDomainsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ens_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResolveDomainsRequest.ProtoReflect.Descriptor instead.
func (*ResolveDomainsRequest) Descriptor() ([]byte, []int) {
	return file_ens_proto_rawDescGZIP(), []int{1}
}

func (x *ResolveDomainsRequest) GetAddresses() []string {
	if x != nil {
		return x.Addresses
	}
	return nil
}

type Address struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EnsName string `protobuf:"bytes,1,opt,name=ens_name,json=ensName,proto3" json:"ens_name,omitempty"`
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *Address) Reset() {
	*x = Address{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ens_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Address) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Address) ProtoMessage() {}

func (x *Address) ProtoReflect() protoreflect.Message {
	mi := &file_ens_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Address.ProtoReflect.Descriptor instead.
func (*Address) Descriptor() ([]byte, []int) {
	return file_ens_proto_rawDescGZIP(), []int{2}
}

func (x *Address) GetEnsName() string {
	if x != nil {
		return x.EnsName
	}
	return ""
}

func (x *Address) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

type ResolveResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addresses []*Address `protobuf:"bytes,1,rep,name=addresses,proto3" json:"addresses,omitempty"`
}

func (x *ResolveResponse) Reset() {
	*x = ResolveResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ens_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResolveResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResolveResponse) ProtoMessage() {}

func (x *ResolveResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ens_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResolveResponse.ProtoReflect.Descriptor instead.
func (*ResolveResponse) Descriptor() ([]byte, []int) {
	return file_ens_proto_rawDescGZIP(), []int{3}
}

func (x *ResolveResponse) GetAddresses() []*Address {
	if x != nil {
		return x.Addresses
	}
	return nil
}

var File_ens_proto protoreflect.FileDescriptor

var file_ens_proto_rawDesc = []byte{
	0x0a, 0x09, 0x65, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x33, 0x0a, 0x17, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a,
	0x07, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07,
	0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x22, 0x35, 0x0a, 0x15, 0x52, 0x65, 0x73, 0x6f, 0x6c,
	0x76, 0x65, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1c, 0x0a, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x22, 0x3e,
	0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x6e, 0x73,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x6e, 0x73,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22, 0x3f,
	0x0a, 0x0f, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x2c, 0x0a, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x52, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x32,
	0x99, 0x01, 0x0a, 0x03, 0x45, 0x6e, 0x73, 0x12, 0x4a, 0x0a, 0x10, 0x52, 0x65, 0x73, 0x6f, 0x6c,
	0x76, 0x65, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x12, 0x1e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x41, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a, 0x0e, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x44, 0x6f,
	0x6d, 0x61, 0x69, 0x6e, 0x73, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65,
	0x73, 0x6f, 0x6c, 0x76, 0x65, 0x44, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x6f,
	0x6c, 0x76, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0a, 0x5a, 0x08, 0x2e,
	0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ens_proto_rawDescOnce sync.Once
	file_ens_proto_rawDescData = file_ens_proto_rawDesc
)

func file_ens_proto_rawDescGZIP() []byte {
	file_ens_proto_rawDescOnce.Do(func() {
		file_ens_proto_rawDescData = protoimpl.X.CompressGZIP(file_ens_proto_rawDescData)
	})
	return file_ens_proto_rawDescData
}

var file_ens_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_ens_proto_goTypes = []interface{}{
	(*ResolveAddressesRequest)(nil), // 0: proto.ResolveAddressesRequest
	(*ResolveDomainsRequest)(nil),   // 1: proto.ResolveDomainsRequest
	(*Address)(nil),                 // 2: proto.Address
	(*ResolveResponse)(nil),         // 3: proto.ResolveResponse
}
var file_ens_proto_depIdxs = []int32{
	2, // 0: proto.ResolveResponse.addresses:type_name -> proto.Address
	0, // 1: proto.Ens.ResolveAddresses:input_type -> proto.ResolveAddressesRequest
	1, // 2: proto.Ens.ResolveDomains:input_type -> proto.ResolveDomainsRequest
	3, // 3: proto.Ens.ResolveAddresses:output_type -> proto.ResolveResponse
	3, // 4: proto.Ens.ResolveDomains:output_type -> proto.ResolveResponse
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_ens_proto_init() }
func file_ens_proto_init() {
	if File_ens_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ens_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResolveAddressesRequest); i {
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
		file_ens_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResolveDomainsRequest); i {
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
		file_ens_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Address); i {
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
		file_ens_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResolveResponse); i {
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
			RawDescriptor: file_ens_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ens_proto_goTypes,
		DependencyIndexes: file_ens_proto_depIdxs,
		MessageInfos:      file_ens_proto_msgTypes,
	}.Build()
	File_ens_proto = out.File
	file_ens_proto_rawDesc = nil
	file_ens_proto_goTypes = nil
	file_ens_proto_depIdxs = nil
}