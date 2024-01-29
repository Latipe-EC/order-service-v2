// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.2
// source: internal/infrastructure/grpc/deliveryServ/delivery_service.proto

package deliverygrpc

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

type GetShippingCostRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SrcCode    string `protobuf:"bytes,1,opt,name=src_code,json=srcCode,proto3" json:"src_code,omitempty"`
	DestCode   string `protobuf:"bytes,2,opt,name=dest_code,json=destCode,proto3" json:"dest_code,omitempty"`
	DeliveryId string `protobuf:"bytes,3,opt,name=delivery_id,json=deliveryId,proto3" json:"delivery_id,omitempty"`
}

func (x *GetShippingCostRequest) Reset() {
	*x = GetShippingCostRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetShippingCostRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetShippingCostRequest) ProtoMessage() {}

func (x *GetShippingCostRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetShippingCostRequest.ProtoReflect.Descriptor instead.
func (*GetShippingCostRequest) Descriptor() ([]byte, []int) {
	return file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetShippingCostRequest) GetSrcCode() string {
	if x != nil {
		return x.SrcCode
	}
	return ""
}

func (x *GetShippingCostRequest) GetDestCode() string {
	if x != nil {
		return x.DestCode
	}
	return ""
}

func (x *GetShippingCostRequest) GetDeliveryId() string {
	if x != nil {
		return x.DeliveryId
	}
	return ""
}

type GetShippingCostResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReceiveDate  string `protobuf:"bytes,1,opt,name=receive_date,json=receiveDate,proto3" json:"receive_date,omitempty"`
	DeliveryId   string `protobuf:"bytes,2,opt,name=delivery_id,json=deliveryId,proto3" json:"delivery_id,omitempty"`
	DeliveryName string `protobuf:"bytes,3,opt,name=delivery_name,json=deliveryName,proto3" json:"delivery_name,omitempty"`
	Cost         int64  `protobuf:"varint,4,opt,name=cost,proto3" json:"cost,omitempty"`
}

func (x *GetShippingCostResponse) Reset() {
	*x = GetShippingCostResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetShippingCostResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetShippingCostResponse) ProtoMessage() {}

func (x *GetShippingCostResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetShippingCostResponse.ProtoReflect.Descriptor instead.
func (*GetShippingCostResponse) Descriptor() ([]byte, []int) {
	return file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetShippingCostResponse) GetReceiveDate() string {
	if x != nil {
		return x.ReceiveDate
	}
	return ""
}

func (x *GetShippingCostResponse) GetDeliveryId() string {
	if x != nil {
		return x.DeliveryId
	}
	return ""
}

func (x *GetShippingCostResponse) GetDeliveryName() string {
	if x != nil {
		return x.DeliveryName
	}
	return ""
}

func (x *GetShippingCostResponse) GetCost() int64 {
	if x != nil {
		return x.Cost
	}
	return 0
}

var File_internal_infrastructure_grpc_deliveryServ_delivery_service_proto protoreflect.FileDescriptor

var file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_rawDesc = []byte{
	0x0a, 0x40, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x69, 0x6e, 0x66, 0x72, 0x61,
	0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x64,
	0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x53, 0x65, 0x72, 0x76, 0x2f, 0x64, 0x65, 0x6c, 0x69,
	0x76, 0x65, 0x72, 0x79, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x71, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x53, 0x68, 0x69, 0x70, 0x70, 0x69, 0x6e,
	0x67, 0x43, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08,
	0x73, 0x72, 0x63, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x73, 0x72, 0x63, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x65, 0x73, 0x74, 0x5f,
	0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x65, 0x73, 0x74,
	0x43, 0x6f, 0x64, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79,
	0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x65, 0x6c, 0x69, 0x76,
	0x65, 0x72, 0x79, 0x49, 0x64, 0x22, 0x96, 0x01, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x53, 0x68, 0x69,
	0x70, 0x70, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x5f, 0x64, 0x61, 0x74,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65,
	0x44, 0x61, 0x74, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x65, 0x6c, 0x69, 0x76,
	0x65, 0x72, 0x79, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72,
	0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x64, 0x65,
	0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f,
	0x73, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x63, 0x6f, 0x73, 0x74, 0x32, 0x63,
	0x0a, 0x13, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x47, 0x52, 0x50, 0x43, 0x12, 0x4c, 0x0a, 0x15, 0x43, 0x61, 0x6c, 0x63, 0x75, 0x6c, 0x61,
	0x74, 0x65, 0x53, 0x68, 0x69, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x73, 0x74, 0x12, 0x17,
	0x2e, 0x47, 0x65, 0x74, 0x53, 0x68, 0x69, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x73, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x68, 0x69,
	0x70, 0x70, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x0f, 0x5a, 0x0d, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x79, 0x67,
	0x72, 0x70, 0x63, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_rawDescOnce sync.Once
	file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_rawDescData = file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_rawDesc
)

func file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_rawDescGZIP() []byte {
	file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_rawDescOnce.Do(func() {
		file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_rawDescData)
	})
	return file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_rawDescData
}

var file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_goTypes = []interface{}{
	(*GetShippingCostRequest)(nil),  // 0: GetShippingCostRequest
	(*GetShippingCostResponse)(nil), // 1: GetShippingCostResponse
}
var file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_depIdxs = []int32{
	0, // 0: DeliveryServiceGRPC.CalculateShippingCost:input_type -> GetShippingCostRequest
	1, // 1: DeliveryServiceGRPC.CalculateShippingCost:output_type -> GetShippingCostResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_init() }
func file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_init() {
	if File_internal_infrastructure_grpc_deliveryServ_delivery_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetShippingCostRequest); i {
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
		file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetShippingCostResponse); i {
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
			RawDescriptor: file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_goTypes,
		DependencyIndexes: file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_depIdxs,
		MessageInfos:      file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_msgTypes,
	}.Build()
	File_internal_infrastructure_grpc_deliveryServ_delivery_service_proto = out.File
	file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_rawDesc = nil
	file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_goTypes = nil
	file_internal_infrastructure_grpc_deliveryServ_delivery_service_proto_depIdxs = nil
}
