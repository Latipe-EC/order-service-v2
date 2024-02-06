package order

type CommissionRepository interface {
	//cms
	UpdateOrderCommission(order *Order, ocms *OrderCommission, log *OrderStatusLog) error
	CreateOrderCommission(ocms *OrderCommission) error
	FindCommissionByOrderId(orderId int) (*OrderCommission, error)
	UpdateCommission(Id string) (*Order, error)
}
