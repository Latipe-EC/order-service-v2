package order

import "time"

type BaseHeader struct {
	BearerToken string `reqHeader:"Authorization"`
}

type OrderResponse struct {
	OrderID          string            `json:"order_id"`
	Amount           int               `json:"amount"`
	ShippingDiscount int               `json:"shipping_discount"`
	ItemDiscount     int               `json:"item_discount"`
	SubTotal         int               `json:"sub_total"`
	Status           int               `json:"status"`
	PaymentMethod    int               `json:"payment_method"`
	VoucherCode      string            `json:"vouchers"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	Delivery         DeliveryResp      `json:"delivery"`
	OrderItems       []OrderItemsResp  `json:"order_items,omitempty"`
	OrderStatus      []OrderStatusResp `json:"order_status,omitempty"`
}

type DeliveryResp struct {
	DeliveryId      string    `json:"delivery_id"`
	DeliveryName    string    `json:"delivery_name"`
	Cost            int       `json:"cost"`
	ReceivingDate   time.Time `json:"receiving_date"`
	AddressId       string    `json:"address_id"`
	ShippingName    string    `json:"shipping_name" `
	ShippingPhone   string    `json:"shipping_phone" `
	ShippingAddress string    `json:"shipping_address" `
}

type OrderItemsResp struct {
	ItemId      string `json:"item_id"`
	ProductId   string `json:"product_id" `
	SubTotal    int    `json:"sub_total"`
	OptionId    string `json:"option_id"`
	Quantity    int    `json:"quantity"`
	NetPrice    int    `json:"net_price"`
	RatingID    string `json:"rating_id,omitempty"`
	ProductName string `json:"product_name"`
	ProdImg     string `json:"image"`
	StoreID     string `json:"store_id"`
	Price       int    `json:"price" `
}

type OrderStatusResp struct {
	Message      string    `json:"message_queue"`
	StatusChange int       `json:"status_change"`
	CreatedAt    time.Time `json:"created_at"`
}
