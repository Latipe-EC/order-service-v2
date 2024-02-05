package msgDTO

type CheckoutMessage struct {
	CheckoutID    string      `json:"checkout_id"`
	UserID        string      `json:"user_id"`
	OrderData     []OrderData `json:"order_data"`
	TotalAmount   uint        `json:"total_amount"`
	PaymentMethod int         `json:"payment_method"`
}

type OrderData struct {
	OrderID string `json:"order_id"`
	Amount  uint   `json:"amount"`
}
