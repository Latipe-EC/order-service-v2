package ordercommand

import (
	"context"
	orderDTO "latipe-order-service-v2/internal/domain/dto/order"
	"latipe-order-service-v2/internal/domain/dto/order/delivery"
	"latipe-order-service-v2/internal/domain/dto/order/store"
	"latipe-order-service-v2/internal/domain/msgDTO"
)

type OrderCommandUsecase interface {
	UpdateStatusOrder(ctx context.Context, dto *orderDTO.UpdateOrderStatusRequest) error
	UpdateOrder(ctx context.Context, dto *orderDTO.UpdateOrderRequest) error
	AdminCancelOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error
	UpdateOrderStatusByEvent(ctx context.Context, dto *msgDTO.OrderStatusMessage) error

	//user
	CreateOrder(ctx context.Context, dto *orderDTO.CreateOrderRequest) (*orderDTO.CreateOrderResponse, error)
	CancelOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error
	UserRefundOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error
	//store
	CancelOrderItem(ctx context.Context, dto *store.UpdateOrderItemRequest) (*store.UpdateOrderItemResponse, error)
	UpdateOrderItem(ctx context.Context, dto *store.UpdateOrderItemRequest) (*store.UpdateOrderItemResponse, error)
	//deli
	DeliveryUpdateStatusOrder(ctx context.Context, dto delivery.UpdateOrderStatusRequest) (*delivery.UpdateOrderStatusResponse, error)
}
