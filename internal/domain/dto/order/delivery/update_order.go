package delivery

type UpdateOrderStatusRequest struct {
	OrderID    string `params:"id" validate:"required"`
	Status     int    `json:"status"  validate:"required"`
	DeliveryID string
}

type UpdateOrderStatusResponse struct {
}
