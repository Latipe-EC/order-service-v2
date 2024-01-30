package store

type StoreUpdateOrderStatusRequest struct {
	OrderID string `params:"id" validate:"required"`
	Message string `json:"message"`
	Status  int    `json:"status" validate:"required"`
	StoreId string
}
