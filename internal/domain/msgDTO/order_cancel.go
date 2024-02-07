package msgDTO

type OrderCancelMessage struct {
	OrderID      string `json:"order_id"`
	CancelStatus int    `json:"cancel_status"`
}
