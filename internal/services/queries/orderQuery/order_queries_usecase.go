package orderQuery

import (
	"context"
	orderDTO "latipe-order-service-v2/internal/domain/dto/order"
	"latipe-order-service-v2/internal/domain/dto/order/delivery"
	internalDTO "latipe-order-service-v2/internal/domain/dto/order/internal-service"
	"latipe-order-service-v2/internal/domain/dto/order/store"
)

type OrderQueryUsecase interface {
	//admin
	GetOrderById(ctx context.Context, dto *orderDTO.GetOrderByIDRequest) (*orderDTO.GetOrderResponse, error)
	AdminCountingOrderAmount(ctx context.Context, dto *orderDTO.CountingOrderAmountRequest) (*orderDTO.CountingOrderAmountResponse, error)
	GetOrderList(ctx context.Context, dto *orderDTO.GetOrderListRequest) (*orderDTO.GetOrderListResponse, error)
	CheckProductPurchased(ctx context.Context, dto *orderDTO.CheckUserOrderRequest) (*orderDTO.CheckUserOrderResponse, error)

	InternalGetRatingID(ctx context.Context, dto *internalDTO.GetOrderRatingItemRequest) (*internalDTO.GetOrderRatingItemResponse, error)

	//user
	GetOrderByID(ctx context.Context, dto *orderDTO.GetOrderByIDRequest) (*orderDTO.GetOrderResponse, error)
	GetOrderByUserId(ctx context.Context, dto *orderDTO.GetByUserIdRequest) (*orderDTO.GetByUserIdResponse, error)
	UserCountingOrder(ctx context.Context, dto *orderDTO.CountingOrderAmountRequest) (*orderDTO.CountingOrderAmountResponse, error)

	//store
	SearchStoreOrderId(ctx context.Context, dto *store.FindStoreOrderRequest) (*orderDTO.GetOrderListResponse, error)
	GetOrdersOfStore(ctx context.Context, dto *store.GetStoreOrderRequest) (*orderDTO.GetOrderListResponse, error)
	ViewDetailStoreOrder(ctx context.Context, dto *store.GetOrderOfStoreByIDRequest) (*store.GetOrderOfStoreByIDResponse, error)
	StoreCountingOrder(ctx context.Context, dto *orderDTO.CountingOrderAmountRequest) (*orderDTO.CountingOrderAmountResponse, error)

	//delivery
	GetOrdersOfDelivery(ctx context.Context, dto *delivery.GetOrderListRequest) (*delivery.GetOrderListResponse, error)
	DeliveryCountingOrder(ctx context.Context, dto *orderDTO.CountingOrderAmountRequest) (*orderDTO.CountingOrderAmountResponse, error)
}
