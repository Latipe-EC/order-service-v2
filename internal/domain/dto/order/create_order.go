package order

import "time"

type CreateOrderRequest struct {
	Header      BaseHeader
	UserRequest UserRequest
	Address     OrderAddress `json:"address" validate:"required"`
	StoreOrders []StoreOrder `json:"store_orders"`
	CartItems   []string     `json:"cart_items"`
}

type StoreOrder struct {
	StoreID       string       `json:"store_id" validate:"required"`
	PaymentMethod int          `json:"payment_method" validate:"required"`
	VoucherCode   []string     `json:"vouchers"`
	Delivery      Delivery     `json:"delivery" validate:"required"`
	Items         []OrderItems `json:"items" validate:"required"`
}

type CreateOrderResponse struct {
	UserOrder UserRequest `json:"user_order"`
	OrderKeys []string    `json:"order_keys"`
	CreatedAt time.Time   `json:"created_at"`
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
