package order

type UpdateOrderStatusRequest struct {
	Role    string
	OrderID string `json:"order_id"`
	UserId  string
	Status  int `json:"status"`
}
type UpdateOrderStatusResponse struct {
}

type UpdateOrderRequest struct {
	Header  BaseHeader
	OrderID string `json:"order_id"`
	UserId  string
	Status  int `json:"status"`
}
