package delivery

type UpdateOrderStatusRequest struct {
	OrderID    string `params:"id" validate:"required"`
	Status     int    `json:"status"  validate:"required"`
	Message    string `json:"message"`
	DeliveryID string
}
