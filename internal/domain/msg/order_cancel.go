package msg

import "time"

type OrderCancelMessage struct {
	OrderID       string    `json:"OrderID"`
	OrderStatus   int       `json:"orderStatus"`
	PaymentMethod int       `json:"paymentMethod"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
