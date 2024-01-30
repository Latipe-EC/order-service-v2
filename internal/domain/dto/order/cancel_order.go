package order

type CancelOrderRequest struct {
	Header  BaseHeader
	OrderID string `json:"order_id" validate:"required"`
	Message string `json:"message" validate:"required"`
	UserId  string
}

type CancelOrderResponse struct {
}
