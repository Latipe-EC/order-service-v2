package order

import (
	"latipe-order-service-v2/internal/domain/msgDTO"
)

type CreateOrderRequest struct {
	Header        BaseHeader
	UserRequest   UserRequest
	PaymentMethod int            `json:"payment_method" validate:"required,min=0,max=3"`
	Address       OrderAddress   `json:"address" validate:"required"`
	StoreOrders   []StoreOrder   `json:"store_orders"`
	PromotionData *PromotionData `json:"promotion_data"`
}

type StoreOrder struct {
	StoreID  string       `json:"store_id" validate:"required"`
	Delivery Delivery     `json:"delivery" validate:"required"`
	Items    []OrderItems `json:"order_items" validate:"required"`
	CartIds  []string     `json:"cart_ids"`
}

type CreateOrderResponse struct {
	msgDTO.CheckoutMessage
	FailedOrder FailedOrder `json:"failed_order,omitempty"`
}

type FailedOrder struct {
	StoreID string `json:"store_id,omitempty"`
	Message string `json:"message,omitempty"`
}

type UserRequest struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type OrderItems struct {
	CartId    string `json:"cart_id,omitempty"`
	ProductId string `json:"product_id" validate:"required"`
	OptionId  string `json:"option_id"`
	Quantity  int    `json:"quantity" validate:"required"`
	Price     int    `json:"price" validate:"required"`
}

type OrderAddress struct {
	AddressId string `json:"address_id" validate:"required"`
}

type Delivery struct {
	DeliveryId string `json:"delivery_id" validate:"required"`
}

type PromotionData struct {
	FreeShippingVoucherInfo *FreeShippingVoucherInfo `json:"free_shipping_voucher"`
	PaymentVoucherInfo      *PaymentVoucherInfo      `json:"payment_voucher"`
	ShopVoucherInfo         []ShopVoucherInfo        `json:"shop_vouchers"`
}

type FreeShippingVoucherInfo struct {
	StoreIds    []string `json:"store_ids"`
	VoucherCode string   `json:"voucher_code"`
}

type PaymentVoucherInfo struct {
	VoucherCode string `json:"voucher_code"`
}

type ShopVoucherInfo struct {
	StoreId     string `json:"store_id" validate:"required"`
	VoucherCode string `json:"voucher_code" validate:"required"`
}
