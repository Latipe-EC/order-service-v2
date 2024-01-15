package store

type UpdateOrderItemRequest struct {
	OrderID string `params:"id" validate:"required"`
	ItemID  string `json:"item_id"`
	StoreId string
}

type UpdateOrderItemResponse struct {
	OrderID string `json:"order_id"`
	ItemID  string `json:"item_id"`
	Status  int    `json:"status"`
}
