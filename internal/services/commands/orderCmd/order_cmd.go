package orderCmd

import (
	"context"
	orderDTO "latipe-order-service-v2/internal/domain/dto/order"
	"latipe-order-service-v2/internal/domain/dto/order/delivery"
	"latipe-order-service-v2/internal/domain/dto/order/store"
	"latipe-order-service-v2/internal/domain/msgDTO"
)

type OrderCommandUsecase interface {
	CreateOrder(ctx context.Context, dto *orderDTO.CreateOrderRequest) (*orderDTO.CreateOrderResponse, error)

	UpdateStatusOrder(ctx context.Context, dto *orderDTO.UpdateOrderStatusRequest) error
	UpdateOrder(ctx context.Context, dto *orderDTO.UpdateOrderRequest) error
	UpdateOrderStatusByReplyMessage(ctx context.Context, dto *msgDTO.OrderStatusMessage) error
	StoreUpdateOrderStatus(ctx context.Context, dto *store.StoreUpdateOrderStatusRequest) error
	DeliveryUpdateOrderStatus(ctx context.Context, dto delivery.UpdateOrderStatusRequest) error

	UserCancelOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error
	UserRefundOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error
	AdminCancelOrder(ctx context.Context, dto *orderDTO.CancelOrderRequest) error

	UpdateRatingItem(ctx context.Context, data *msgDTO.RatingMessage) error
}
