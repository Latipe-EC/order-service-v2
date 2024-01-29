// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.25.2
// source: internal/infrastructure/grpc/productServ/product_service.proto

package productgrpc

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

type GetPurchaseProductRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StoreId string                    `protobuf:"bytes,1,opt,name=storeId,proto3" json:"storeId,omitempty"`
	Items   []*GetPurchaseItemRequest `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *GetPurchaseProductRequest) Reset() {
	*x = GetPurchaseProductRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_infrastructure_grpc_productServ_product_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPurchaseProductRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPurchaseProductRequest) ProtoMessage() {}

func (x *GetPurchaseProductRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_infrastructure_grpc_productServ_product_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPurchaseProductRequest.ProtoReflect.Descriptor instead.
func (*GetPurchaseProductRequest) Descriptor() ([]byte, []int) {
	return file_internal_infrastructure_grpc_productServ_product_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetPurchaseProductRequest) GetStoreId() string {
	if x != nil {
		return x.StoreId
	}
	return ""
}

func (x *GetPurchaseProductRequest) GetItems() []*GetPurchaseItemRequest {
	if x != nil {
		return x.Items
	}
	return nil
}

type GetPurchaseItemRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductId string `protobuf:"bytes,1,opt,name=productId,proto3" json:"productId,omitempty"`
	OptionId  string `protobuf:"bytes,2,opt,name=optionId,proto3" json:"optionId,omitempty"`
	Quantity  int32  `protobuf:"varint,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
}

func (x *GetPurchaseItemRequest) Reset() {
	*x = GetPurchaseItemRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_infrastructure_grpc_productServ_product_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPurchaseItemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPurchaseItemRequest) ProtoMessage() {}

func (x *GetPurchaseItemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_infrastructure_grpc_productServ_product_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPurchaseItemRequest.ProtoReflect.Descriptor instead.
func (*GetPurchaseItemRequest) Descriptor() ([]byte, []int) {
	return file_internal_infrastructure_grpc_productServ_product_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetPurchaseItemRequest) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}

func (x *GetPurchaseItemRequest) GetOptionId() string {
	if x != nil {
		return x.OptionId
	}
	return ""
}

func (x *GetPurchaseItemRequest) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

type ItemResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductId        string  `protobuf:"bytes,1,opt,name=productId,proto3" json:"productId,omitempty"`
	Name             string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Quantity         int32   `protobuf:"varint,3,opt,name=quantity,proto3" json:"quantity,omitempty"`
	Image            string  `protobuf:"bytes,4,opt,name=image,proto3" json:"image,omitempty"`
	Price            float32 `protobuf:"fixed32,5,opt,name=price,proto3" json:"price,omitempty"`
	PromotionalPrice float32 `protobuf:"fixed32,6,opt,name=promotionalPrice,proto3" json:"promotionalPrice,omitempty"`
	OptionId         string  `protobuf:"bytes,7,opt,name=optionId,proto3" json:"optionId,omitempty"`
	NameOption       string  `protobuf:"bytes,8,opt,name=nameOption,proto3" json:"nameOption,omitempty"`
	StoreId          string  `protobuf:"bytes,9,opt,name=storeId,proto3" json:"storeId,omitempty"`
	TotalPrice       float32 `protobuf:"fixed32,10,opt,name=totalPrice,proto3" json:"totalPrice,omitempty"`
}

func (x *ItemResponse) Reset() {
	*x = ItemResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_infrastructure_grpc_productServ_product_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ItemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ItemResponse) ProtoMessage() {}

func (x *ItemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_infrastructure_grpc_productServ_product_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ItemResponse.ProtoReflect.Descriptor instead.
func (*ItemResponse) Descriptor() ([]byte, []int) {
	return file_internal_infrastructure_grpc_productServ_product_service_proto_rawDescGZIP(), []int{2}
}

func (x *ItemResponse) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}

func (x *ItemResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ItemResponse) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

func (x *ItemResponse) GetImage() string {
	if x != nil {
		return x.Image
	}
	return ""
}

func (x *ItemResponse) GetPrice() float32 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *ItemResponse) GetPromotionalPrice() float32 {
	if x != nil {
		return x.PromotionalPrice
	}
	return 0
}

func (x *ItemResponse) GetOptionId() string {
	if x != nil {
		return x.OptionId
	}
	return ""
}

func (x *ItemResponse) GetNameOption() string {
	if x != nil {
		return x.NameOption
	}
	return ""
}

func (x *ItemResponse) GetStoreId() string {
	if x != nil {
		return x.StoreId
	}
	return ""
}

func (x *ItemResponse) GetTotalPrice() float32 {
	if x != nil {
		return x.TotalPrice
	}
	return 0
}

type GetPurchaseItemResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StoreId      string          `protobuf:"bytes,1,opt,name=storeId,proto3" json:"storeId,omitempty"`
	ProvinceCode string          `protobuf:"bytes,2,opt,name=provinceCode,proto3" json:"provinceCode,omitempty"`
	TotalPrice   int64           `protobuf:"varint,3,opt,name=totalPrice,proto3" json:"totalPrice,omitempty"`
	Items        []*ItemResponse `protobuf:"bytes,4,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *GetPurchaseItemResponse) Reset() {
	*x = GetPurchaseItemResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_infrastructure_grpc_productServ_product_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPurchaseItemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPurchaseItemResponse) ProtoMessage() {}

func (x *GetPurchaseItemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_infrastructure_grpc_productServ_product_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPurchaseItemResponse.ProtoReflect.Descriptor instead.
func (*GetPurchaseItemResponse) Descriptor() ([]byte, []int) {
	return file_internal_infrastructure_grpc_productServ_product_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetPurchaseItemResponse) GetStoreId() string {
	if x != nil {
		return x.StoreId
	}
	return ""
}

func (x *GetPurchaseItemResponse) GetProvinceCode() string {
	if x != nil {
		return x.ProvinceCode
	}
	return ""
}

func (x *GetPurchaseItemResponse) GetTotalPrice() int64 {
	if x != nil {
		return x.TotalPrice
	}
	return 0
}

func (x *GetPurchaseItemResponse) GetItems() []*ItemResponse {
	if x != nil {
		return x.Items
	}
	return nil
}

type UpdateProductQuantityRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StoreId string                    `protobuf:"bytes,1,opt,name=storeId,proto3" json:"storeId,omitempty"`
	Items   []*GetPurchaseItemRequest `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *UpdateProductQuantityRequest) Reset() {
	*x = UpdateProductQuantityRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_infrastructure_grpc_productServ_product_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateProductQuantityRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateProductQuantityRequest) ProtoMessage() {}

func (x *UpdateProductQuantityRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_infrastructure_grpc_productServ_product_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateProductQuantityRequest.ProtoReflect.Descriptor instead.
func (*UpdateProductQuantityRequest) Descriptor() ([]byte, []int) {
	return file_internal_infrastructure_grpc_productServ_product_service_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateProductQuantityRequest) GetStoreId() string {
	if x != nil {
		return x.StoreId
	}
	return ""
}

func (x *UpdateProductQuantityRequest) GetItems() []*GetPurchaseItemRequest {
	if x != nil {
		return x.Items
	}
	return nil
}

type UpdateProductQuantityResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsSuccess bool `protobuf:"varint,1,opt,name=isSuccess,proto3" json:"isSuccess,omitempty"`
}

func (x *UpdateProductQuantityResponse) Reset() {
	*x = UpdateProductQuantityResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_infrastructure_grpc_productServ_product_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateProductQuantityResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateProductQuantityResponse) ProtoMessage() {}

func (x *UpdateProductQuantityResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_infrastructure_grpc_productServ_product_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateProductQuantityResponse.ProtoReflect.Descriptor instead.
func (*UpdateProductQuantityResponse) Descriptor() ([]byte, []int) {
	return file_internal_infrastructure_grpc_productServ_product_service_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateProductQuantityResponse) GetIsSuccess() bool {
	if x != nil {
		return x.IsSuccess
	}
	return false
}

var File_internal_infrastructure_grpc_productServ_product_service_proto protoreflect.FileDescriptor

var file_internal_infrastructure_grpc_productServ_product_service_proto_rawDesc = []byte{
	0x0a, 0x3e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x69, 0x6e, 0x66, 0x72, 0x61,
	0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x53, 0x65, 0x72, 0x76, 0x2f, 0x70, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x22, 0x6d, 0x0a, 0x19, 0x47, 0x65,
	0x74, 0x50, 0x75, 0x72, 0x63, 0x68, 0x61, 0x73, 0x65, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x49,
	0x64, 0x12, 0x36, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x47, 0x65, 0x74, 0x50,
	0x75, 0x72, 0x63, 0x68, 0x61, 0x73, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x6e, 0x0a, 0x16, 0x47, 0x65, 0x74,
	0x50, 0x75, 0x72, 0x63, 0x68, 0x61, 0x73, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49,
	0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1a, 0x0a,
	0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x22, 0xaa, 0x02, 0x0a, 0x0c, 0x49, 0x74,
	0x65, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08,
	0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x70,
	0x72, 0x69, 0x63, 0x65, 0x12, 0x2a, 0x0a, 0x10, 0x70, 0x72, 0x6f, 0x6d, 0x6f, 0x74, 0x69, 0x6f,
	0x6e, 0x61, 0x6c, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x02, 0x52, 0x10,
	0x70, 0x72, 0x6f, 0x6d, 0x6f, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x50, 0x72, 0x69, 0x63, 0x65,
	0x12, 0x1a, 0x0a, 0x08, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0a,
	0x6e, 0x61, 0x6d, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x6e, 0x61, 0x6d, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x50,
	0x72, 0x69, 0x63, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61,
	0x6c, 0x50, 0x72, 0x69, 0x63, 0x65, 0x22, 0xa5, 0x01, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x50, 0x75,
	0x72, 0x63, 0x68, 0x61, 0x73, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x64, 0x12, 0x22, 0x0a, 0x0c,
	0x70, 0x72, 0x6f, 0x76, 0x69, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x50, 0x72, 0x69, 0x63, 0x65,
	0x12, 0x2c, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x49, 0x74, 0x65, 0x6d, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x70,
	0x0a, 0x1c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x51,
	0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18,
	0x0a, 0x07, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x49, 0x64, 0x12, 0x36, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x75, 0x72, 0x63, 0x68, 0x61, 0x73, 0x65, 0x49, 0x74,
	0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73,
	0x22, 0x3d, 0x0a, 0x1d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x73, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x32,
	0xd3, 0x01, 0x0a, 0x12, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x47, 0x52, 0x50, 0x43, 0x12, 0x58, 0x0a, 0x0c, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x49,
	0x6e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x12, 0x23, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x75, 0x72, 0x63, 0x68, 0x61, 0x73, 0x65, 0x50, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x75, 0x72, 0x63, 0x68, 0x61,
	0x73, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x63, 0x0a, 0x0e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x12, 0x26, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x51, 0x75, 0x61, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x51, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0e, 0x5a, 0x0c, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x67, 0x72, 0x70, 0x63, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_infrastructure_grpc_productServ_product_service_proto_rawDescOnce sync.Once
	file_internal_infrastructure_grpc_productServ_product_service_proto_rawDescData = file_internal_infrastructure_grpc_productServ_product_service_proto_rawDesc
)

func file_internal_infrastructure_grpc_productServ_product_service_proto_rawDescGZIP() []byte {
	file_internal_infrastructure_grpc_productServ_product_service_proto_rawDescOnce.Do(func() {
		file_internal_infrastructure_grpc_productServ_product_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_infrastructure_grpc_productServ_product_service_proto_rawDescData)
	})
	return file_internal_infrastructure_grpc_productServ_product_service_proto_rawDescData
}

var file_internal_infrastructure_grpc_productServ_product_service_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_internal_infrastructure_grpc_productServ_product_service_proto_goTypes = []interface{}{
	(*GetPurchaseProductRequest)(nil),     // 0: protobuf.GetPurchaseProductRequest
	(*GetPurchaseItemRequest)(nil),        // 1: protobuf.GetPurchaseItemRequest
	(*ItemResponse)(nil),                  // 2: protobuf.ItemResponse
	(*GetPurchaseItemResponse)(nil),       // 3: protobuf.GetPurchaseItemResponse
	(*UpdateProductQuantityRequest)(nil),  // 4: protobuf.UpdateProductQuantityRequest
	(*UpdateProductQuantityResponse)(nil), // 5: protobuf.UpdateProductQuantityResponse
}
var file_internal_infrastructure_grpc_productServ_product_service_proto_depIdxs = []int32{
	1, // 0: protobuf.GetPurchaseProductRequest.items:type_name -> protobuf.GetPurchaseItemRequest
	2, // 1: protobuf.GetPurchaseItemResponse.items:type_name -> protobuf.ItemResponse
	1, // 2: protobuf.UpdateProductQuantityRequest.items:type_name -> protobuf.GetPurchaseItemRequest
	0, // 3: protobuf.ProductServiceGRPC.CheckInStock:input_type -> protobuf.GetPurchaseProductRequest
	4, // 4: protobuf.ProductServiceGRPC.UpdateQuantity:input_type -> protobuf.UpdateProductQuantityRequest
	3, // 5: protobuf.ProductServiceGRPC.CheckInStock:output_type -> protobuf.GetPurchaseItemResponse
	5, // 6: protobuf.ProductServiceGRPC.UpdateQuantity:output_type -> protobuf.UpdateProductQuantityResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_internal_infrastructure_grpc_productServ_product_service_proto_init() }
func file_internal_infrastructure_grpc_productServ_product_service_proto_init() {
	if File_internal_infrastructure_grpc_productServ_product_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_infrastructure_grpc_productServ_product_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPurchaseProductRequest); i {
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
		file_internal_infrastructure_grpc_productServ_product_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPurchaseItemRequest); i {
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
		file_internal_infrastructure_grpc_productServ_product_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ItemResponse); i {
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
		file_internal_infrastructure_grpc_productServ_product_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPurchaseItemResponse); i {
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
		file_internal_infrastructure_grpc_productServ_product_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateProductQuantityRequest); i {
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
		file_internal_infrastructure_grpc_productServ_product_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateProductQuantityResponse); i {
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
			RawDescriptor: file_internal_infrastructure_grpc_productServ_product_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_infrastructure_grpc_productServ_product_service_proto_goTypes,
		DependencyIndexes: file_internal_infrastructure_grpc_productServ_product_service_proto_depIdxs,
		MessageInfos:      file_internal_infrastructure_grpc_productServ_product_service_proto_msgTypes,
	}.Build()
	File_internal_infrastructure_grpc_productServ_product_service_proto = out.File
	file_internal_infrastructure_grpc_productServ_product_service_proto_rawDesc = nil
	file_internal_infrastructure_grpc_productServ_product_service_proto_goTypes = nil
	file_internal_infrastructure_grpc_productServ_product_service_proto_depIdxs = nil
}
