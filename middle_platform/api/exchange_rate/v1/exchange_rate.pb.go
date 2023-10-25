// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.19.4
// source: exchange_rate/v1/exchange_rate.proto

package v1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type RateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RateRequest) Reset() {
	*x = RateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_exchange_rate_v1_exchange_rate_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateRequest) ProtoMessage() {}

func (x *RateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_exchange_rate_v1_exchange_rate_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RateRequest.ProtoReflect.Descriptor instead.
func (*RateRequest) Descriptor() ([]byte, []int) {
	return file_exchange_rate_v1_exchange_rate_proto_rawDescGZIP(), []int{0}
}

type RateReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Currencies []string `protobuf:"bytes,1,rep,name=currencies,proto3" json:"currencies,omitempty"`
}

func (x *RateReply) Reset() {
	*x = RateReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_exchange_rate_v1_exchange_rate_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RateReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RateReply) ProtoMessage() {}

func (x *RateReply) ProtoReflect() protoreflect.Message {
	mi := &file_exchange_rate_v1_exchange_rate_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RateReply.ProtoReflect.Descriptor instead.
func (*RateReply) Descriptor() ([]byte, []int) {
	return file_exchange_rate_v1_exchange_rate_proto_rawDescGZIP(), []int{1}
}

func (x *RateReply) GetCurrencies() []string {
	if x != nil {
		return x.Currencies
	}
	return nil
}

type BaseCurrencyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Base string `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
}

func (x *BaseCurrencyRequest) Reset() {
	*x = BaseCurrencyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_exchange_rate_v1_exchange_rate_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BaseCurrencyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BaseCurrencyRequest) ProtoMessage() {}

func (x *BaseCurrencyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_exchange_rate_v1_exchange_rate_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BaseCurrencyRequest.ProtoReflect.Descriptor instead.
func (*BaseCurrencyRequest) Descriptor() ([]byte, []int) {
	return file_exchange_rate_v1_exchange_rate_proto_rawDescGZIP(), []int{2}
}

func (x *BaseCurrencyRequest) GetBase() string {
	if x != nil {
		return x.Base
	}
	return ""
}

type BaseCurrencyReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp int32              `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Base      string             `protobuf:"bytes,2,opt,name=base,proto3" json:"base,omitempty"`
	Rates     map[string]float64 `protobuf:"bytes,3,rep,name=rates,proto3" json:"rates,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed64,2,opt,name=value,proto3"`
}

func (x *BaseCurrencyReply) Reset() {
	*x = BaseCurrencyReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_exchange_rate_v1_exchange_rate_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BaseCurrencyReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BaseCurrencyReply) ProtoMessage() {}

func (x *BaseCurrencyReply) ProtoReflect() protoreflect.Message {
	mi := &file_exchange_rate_v1_exchange_rate_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BaseCurrencyReply.ProtoReflect.Descriptor instead.
func (*BaseCurrencyReply) Descriptor() ([]byte, []int) {
	return file_exchange_rate_v1_exchange_rate_proto_rawDescGZIP(), []int{3}
}

func (x *BaseCurrencyReply) GetTimestamp() int32 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *BaseCurrencyReply) GetBase() string {
	if x != nil {
		return x.Base
	}
	return ""
}

func (x *BaseCurrencyReply) GetRates() map[string]float64 {
	if x != nil {
		return x.Rates
	}
	return nil
}

var File_exchange_rate_v1_exchange_rate_proto protoreflect.FileDescriptor

var file_exchange_rate_v1_exchange_rate_proto_rawDesc = []byte{
	0x0a, 0x24, 0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x2f,
	0x76, 0x31, 0x2f, 0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x72, 0x61, 0x74, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65,
	0x5f, 0x72, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x0d, 0x0a, 0x0b, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x2b, 0x0a, 0x09, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x69,
	0x65, 0x73, 0x22, 0x29, 0x0a, 0x13, 0x42, 0x61, 0x73, 0x65, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e,
	0x63, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x61, 0x73,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x62, 0x61, 0x73, 0x65, 0x22, 0xc5, 0x01,
	0x0a, 0x11, 0x42, 0x61, 0x73, 0x65, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x61, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x62, 0x61, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x05, 0x72, 0x61, 0x74, 0x65, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f,
	0x72, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x61, 0x73, 0x65, 0x43, 0x75, 0x72, 0x72,
	0x65, 0x6e, 0x63, 0x79, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x2e, 0x52, 0x61, 0x74, 0x65, 0x73, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x72, 0x61, 0x74, 0x65, 0x73, 0x1a, 0x38, 0x0a, 0x0a, 0x52,
	0x61, 0x74, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0xf5, 0x01, 0x0a, 0x0c, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x52, 0x61, 0x74, 0x65, 0x12, 0x70, 0x0a, 0x13, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72,
	0x74, 0x65, 0x64, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x12, 0x1d, 0x2e,
	0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x65,
	0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e,
	0x72, 0x61, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x17, 0x12, 0x15, 0x2f, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x2d, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x63, 0x69, 0x65, 0x73, 0x12, 0x73, 0x0a, 0x0c, 0x42, 0x61, 0x73, 0x65,
	0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x25, 0x2e, 0x65, 0x78, 0x63, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x61, 0x73, 0x65,
	0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x23, 0x2e, 0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x42, 0x61, 0x73, 0x65, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x63, 0x79, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0x17, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x12, 0x0f, 0x2f, 0x65,
	0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x2d, 0x72, 0x61, 0x74, 0x65, 0x73, 0x42, 0x61, 0x0a,
	0x1f, 0x64, 0x65, 0x76, 0x2e, 0x6b, 0x72, 0x61, 0x74, 0x6f, 0x73, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x65, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x2e, 0x76, 0x31,
	0x42, 0x13, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x61, 0x74, 0x65, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x56, 0x31, 0x50, 0x01, 0x5a, 0x27, 0x6d, 0x69, 0x64, 0x64, 0x6c, 0x65, 0x5f,
	0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x65, 0x78, 0x63,
	0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x76, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_exchange_rate_v1_exchange_rate_proto_rawDescOnce sync.Once
	file_exchange_rate_v1_exchange_rate_proto_rawDescData = file_exchange_rate_v1_exchange_rate_proto_rawDesc
)

func file_exchange_rate_v1_exchange_rate_proto_rawDescGZIP() []byte {
	file_exchange_rate_v1_exchange_rate_proto_rawDescOnce.Do(func() {
		file_exchange_rate_v1_exchange_rate_proto_rawDescData = protoimpl.X.CompressGZIP(file_exchange_rate_v1_exchange_rate_proto_rawDescData)
	})
	return file_exchange_rate_v1_exchange_rate_proto_rawDescData
}

var file_exchange_rate_v1_exchange_rate_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_exchange_rate_v1_exchange_rate_proto_goTypes = []interface{}{
	(*RateRequest)(nil),         // 0: exchange_rate.v1.rateRequest
	(*RateReply)(nil),           // 1: exchange_rate.v1.rateReply
	(*BaseCurrencyRequest)(nil), // 2: exchange_rate.v1.BaseCurrencyRequest
	(*BaseCurrencyReply)(nil),   // 3: exchange_rate.v1.BaseCurrencyReply
	nil,                         // 4: exchange_rate.v1.BaseCurrencyReply.RatesEntry
}
var file_exchange_rate_v1_exchange_rate_proto_depIdxs = []int32{
	4, // 0: exchange_rate.v1.BaseCurrencyReply.rates:type_name -> exchange_rate.v1.BaseCurrencyReply.RatesEntry
	0, // 1: exchange_rate.v1.ExchangeRate.SupportedCurrencies:input_type -> exchange_rate.v1.rateRequest
	2, // 2: exchange_rate.v1.ExchangeRate.BaseCurrency:input_type -> exchange_rate.v1.BaseCurrencyRequest
	1, // 3: exchange_rate.v1.ExchangeRate.SupportedCurrencies:output_type -> exchange_rate.v1.rateReply
	3, // 4: exchange_rate.v1.ExchangeRate.BaseCurrency:output_type -> exchange_rate.v1.BaseCurrencyReply
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_exchange_rate_v1_exchange_rate_proto_init() }
func file_exchange_rate_v1_exchange_rate_proto_init() {
	if File_exchange_rate_v1_exchange_rate_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_exchange_rate_v1_exchange_rate_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RateRequest); i {
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
		file_exchange_rate_v1_exchange_rate_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RateReply); i {
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
		file_exchange_rate_v1_exchange_rate_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BaseCurrencyRequest); i {
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
		file_exchange_rate_v1_exchange_rate_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BaseCurrencyReply); i {
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
			RawDescriptor: file_exchange_rate_v1_exchange_rate_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_exchange_rate_v1_exchange_rate_proto_goTypes,
		DependencyIndexes: file_exchange_rate_v1_exchange_rate_proto_depIdxs,
		MessageInfos:      file_exchange_rate_v1_exchange_rate_proto_msgTypes,
	}.Build()
	File_exchange_rate_v1_exchange_rate_proto = out.File
	file_exchange_rate_v1_exchange_rate_proto_rawDesc = nil
	file_exchange_rate_v1_exchange_rate_proto_goTypes = nil
	file_exchange_rate_v1_exchange_rate_proto_depIdxs = nil
}
