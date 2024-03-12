package store

import (
	"latipe-order-service-v2/internal/domain/dto/order"
	"time"
)

type GetOrderOfStoreByIDRequest struct {
	BaseHeader order.BaseHeader
	OrderID    string `json:"order_id" params:"id"`
	StoreID    string
}

type GetOrderOfStoreByIDResponse struct {
	StoreOrderResponse
}

type StoreOrderResponse struct {
	OrderID          string             `json:"order_id"`
	StoreOrderAmount int                `json:"store_order_amount,omitempty"`
	Status           int                `json:"status"`
	StoreID          string             `json:"store_id"`
	PaymentMethod    int                `json:"payment_method"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
	Delivery         order.DeliveryResp `json:"delivery"`
	CommissionDetail *CommissionDetail  `json:"commission_detail,omitempty"`
	OrderItems       []OrderStoreItem   `json:"order_items,omitempty"`
}

type DeliveryOrderResponse struct {
	OrderID          string             `json:"order_id"`
	StoreOrderAmount int                `json:"store_order_amount,omitempty"`
	Status           int                `json:"status"`
	PaymentMethod    int                `json:"payment_method"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
	Delivery         order.DeliveryResp `json:"delivery"`
	OrderItems       []OrderStoreItem   `json:"order_items,omitempty"`
}

type CommissionDetail struct {
	AmountReceived    int `json:"amount_received"`
	SystemFee         int `json:"system_fee"`
	DiscountFromStore int ` json:"discount_from_store"`
}

type OrderStoreItem struct {
	Id          string `json:"item_id,omitempty"`
	ProductId   string `json:"product_id" `
	OptionId    string `json:"option_id"`
	Quantity    int    `json:"quantity" `
	Price       int    `json:"price"`
	NetPrice    int    `gorm:"not null;type:bigint" json:"net_price"`
	NameOption  string `json:"name_option"`
	Status      int    `json:"is_prepared"`
	SubTotal    int    `json:"sub_total"`
	ProductName string `json:"product_name"`
	ProdImg     string `json:"image"`
}
