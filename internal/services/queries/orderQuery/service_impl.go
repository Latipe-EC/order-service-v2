package orderQuery

import (
	"context"
	"latipe-order-service-v2/internal/common/errors"
	orderDTO "latipe-order-service-v2/internal/domain/dto/order"
	"latipe-order-service-v2/internal/domain/dto/order/delivery"
	internalDTO "latipe-order-service-v2/internal/domain/dto/order/internal-service"
	"latipe-order-service-v2/internal/domain/dto/order/store"
	"latipe-order-service-v2/internal/domain/entities/order"
	"latipe-order-service-v2/internal/middleware/auth"
	"latipe-order-service-v2/pkg/util/mapper"
)

type orderQueryService struct {
	orderRepo order.OrderRepository
}

func NewOrderQueryService(orderRepo order.OrderRepository) OrderQueryUsecase {
	return &orderQueryService{
		orderRepo: orderRepo,
	}
}

func (o orderQueryService) GetOrderByIdofAdmin(ctx context.Context, dto *orderDTO.GetOrderByIDRequest) (*orderDTO.AdminOrderResponse, error) {
	orderResp := orderDTO.AdminOrderResponse{}

	orderDAO, err := o.orderRepo.FindByID(ctx, dto.OrderId)
	if err != nil {
		return nil, err
	}

	if err = mapper.BindingStruct(orderDAO, &orderResp); err != nil {
		return nil, err
	}

	if err = mapper.BindingStruct(orderDAO.OrderCommission, &orderResp.CommissionDetail); err != nil {
		return nil, err
	}

	return &orderResp, err
}

func (o orderQueryService) AdminCountingOrderAmount(ctx context.Context, dto *orderDTO.CountingOrderAmountRequest) (*orderDTO.CountingOrderAmountResponse, error) {
	count, err := o.orderRepo.AdminCountingOrder(ctx)
	if err != nil {
		return nil, err
	}

	dataResp := orderDTO.CountingOrderAmountResponse{Count: count}
	return &dataResp, nil
}

func (o orderQueryService) UserCountingOrder(ctx context.Context, dto *orderDTO.CountingOrderAmountRequest) (*orderDTO.CountingOrderAmountResponse, error) {
	count, err := o.orderRepo.UserCountingOrder(ctx, dto.OwnerID)
	if err != nil {
		return nil, err
	}

	dataResp := orderDTO.CountingOrderAmountResponse{Count: count}
	return &dataResp, nil
}

func (o orderQueryService) StoreCountingOrder(ctx context.Context, dto *orderDTO.CountingOrderAmountRequest) (*orderDTO.CountingOrderAmountResponse, error) {
	count, err := o.orderRepo.StoreCountingOrder(ctx, dto.OwnerID)
	if err != nil {
		return nil, err
	}

	dataResp := orderDTO.CountingOrderAmountResponse{Count: count}
	return &dataResp, nil
}

func (o orderQueryService) DeliveryCountingOrder(ctx context.Context, dto *orderDTO.CountingOrderAmountRequest) (*orderDTO.CountingOrderAmountResponse, error) {
	count, err := o.orderRepo.DeliveryCountingOrder(ctx, dto.OwnerID)
	if err != nil {
		return nil, err
	}

	dataResp := orderDTO.CountingOrderAmountResponse{Count: count}
	return &dataResp, nil
}

func (o orderQueryService) InternalGetRatingID(ctx context.Context, dto *internalDTO.GetOrderRatingItemRequest) (*internalDTO.GetOrderRatingItemResponse, error) {
	resp := internalDTO.GetOrderRatingItemResponse{}
	orderDAO, err := o.orderRepo.FindByItemId(ctx, dto.ItemID)
	if err != nil {
		return nil, err
	}

	if err = mapper.BindingStruct(orderDAO, &resp); err != nil {
		return nil, err
	}

	return &resp, err
}

func (o orderQueryService) GetOrderByIdOfUser(ctx context.Context, dto *orderDTO.GetOrderByIDRequest) (*orderDTO.GetOrderResponse, error) {
	orderResp := orderDTO.OrderResponse{}

	orderDAO, err := o.orderRepo.FindByID(ctx, dto.OrderId)
	if err != nil {
		return nil, err
	}

	if orderDAO.UserId != dto.OwnerId {
		return nil, errors.ErrNotFound
	}

	if err = mapper.BindingStruct(orderDAO, &orderResp); err != nil {
		return nil, err
	}

	resp := orderDTO.GetOrderResponse{Order: orderResp}

	return &resp, err
}

func (o orderQueryService) GetOrderDetailOfDelivery(ctx context.Context, dto *orderDTO.GetOrderByIDRequest) (*orderDTO.GetOrderResponse, error) {
	orderResp := orderDTO.OrderResponse{}

	orderDAO, err := o.orderRepo.FindByID(ctx, dto.OrderId)
	if err != nil {
		return nil, err
	}

	switch dto.Role {
	case auth.ROLE_USER:
		if orderDAO.UserId != dto.OwnerId {
			return nil, errors.ErrNotFound
		}
	case auth.ROLE_STORE:
		if orderDAO.StoreId != dto.OwnerId {
			return nil, errors.ErrNotFound
		}
	case auth.ROLE_DELIVERY:
		if orderDAO.Delivery.DeliveryId != dto.OwnerId {
			return nil, errors.ErrNotFound
		}
	}

	if err = mapper.BindingStruct(orderDAO, &orderResp); err != nil {
		return nil, err
	}

	resp := orderDTO.GetOrderResponse{Order: orderResp}

	return &resp, err
}

func (o orderQueryService) GetOrderList(ctx context.Context, dto *orderDTO.GetOrderListRequest) (*orderDTO.GetOrderListResponse, error) {
	var dataResp []orderDTO.OrderResponse

	orders, err := o.orderRepo.FindAll(ctx, dto.Query)
	if err != nil {
		return nil, err
	}

	total, err := o.orderRepo.Total(ctx, dto.Query)
	if err != nil {
		return nil, err
	}

	if err = mapper.BindingStruct(orders, &dataResp); err != nil {
		return nil, err
	}

	resp := orderDTO.GetOrderListResponse{}
	resp.Items = dataResp
	resp.Size = dto.Query.Size
	resp.Page = dto.Query.Page
	resp.Total = dto.Query.GetTotalPages(total)
	resp.HasMore = dto.Query.GetHasMore(total)

	return &resp, err
}

func (o orderQueryService) GetOrderByUserId(ctx context.Context, dto *orderDTO.GetByUserIdRequest) (*orderDTO.GetByUserIdResponse, error) {
	var dataResp []orderDTO.OrderResponse

	orders, err := o.orderRepo.FindByUserId(ctx, dto.UserId, dto.Query)
	if err != nil {
		return nil, err
	}

	total, err := o.orderRepo.UserQueryTotal(ctx, dto.UserId, dto.Query)
	if err != nil {
		return nil, err
	}

	if err = mapper.BindingStruct(orders, &dataResp); err != nil {
		return nil, err
	}

	resp := orderDTO.GetByUserIdResponse{}
	resp.Items = dataResp
	resp.Size = dto.Query.Size
	resp.Page = dto.Query.Page
	resp.Total = dto.Query.GetTotalPages(total)
	resp.HasMore = dto.Query.GetHasMore(total)

	return &resp, err
}

func (o orderQueryService) SearchStoreOrderID(ctx context.Context, dto *store.FindStoreOrderRequest) (*orderDTO.GetOrderListResponse, error) {
	var dataResp []store.StoreOrderResponse

	orders, err := o.orderRepo.SearchOrderByStoreID(ctx, dto.StoreID, dto.Keyword, dto.Query)
	if err != nil {
		return nil, err
	}

	total, err := o.orderRepo.TotalSearchOrderByStoreID(ctx, dto.StoreID, dto.Keyword)
	if err != nil {
		return nil, err
	}

	if err = mapper.BindingStruct(orders, &dataResp); err != nil {
		return nil, err
	}

	resp := orderDTO.GetOrderListResponse{}
	resp.Items = dataResp
	resp.Size = dto.Query.Size
	resp.Page = dto.Query.Page
	resp.Total = dto.Query.GetTotalPages(total)
	resp.HasMore = dto.Query.GetHasMore(total)

	return &resp, err
}

func (o orderQueryService) GetOrdersOfStore(ctx context.Context, dto *store.GetStoreOrderRequest) (*orderDTO.GetOrderListResponse, error) {
	var dataResp []store.StoreOrderResponse

	orders, err := o.orderRepo.FindOrderByStoreID(ctx, dto.StoreID, dto.Query, dto.Keyword)
	if err != nil {
		return nil, err
	}

	total, err := o.orderRepo.TotalStoreOrder(ctx, dto.StoreID, dto.Query, dto.Keyword)
	if err != nil {
		return nil, err
	}

	if err = mapper.BindingStruct(orders, &dataResp); err != nil {
		return nil, err
	}

	resp := orderDTO.GetOrderListResponse{}
	resp.Items = dataResp
	resp.Size = dto.Query.Size
	resp.Page = dto.Query.Page
	resp.Total = dto.Query.GetTotalPages(total)
	resp.HasMore = dto.Query.GetHasMore(total)

	return &resp, err
}

func (o orderQueryService) GetOrdersOfDelivery(ctx context.Context, dto *delivery.GetOrderListRequest) (*delivery.GetOrderListResponse, error) {
	var dataResp []store.DeliveryOrderResponse

	orders, err := o.orderRepo.FindOrderByDelivery(ctx, dto.DeliveryID, dto.Keyword, dto.Query)
	if err != nil {
		return nil, err
	}

	total, err := o.orderRepo.TotalOrdersOfDelivery(ctx, dto.DeliveryID, dto.Keyword, dto.Query)
	if err != nil {
		return nil, err
	}

	if err = mapper.BindingStruct(orders, &dataResp); err != nil {
		return nil, err
	}

	resp := delivery.GetOrderListResponse{}
	resp.Items = dataResp
	resp.Size = dto.Query.Size
	resp.Page = dto.Query.Page
	resp.Total = dto.Query.GetTotalPages(total)
	resp.HasMore = dto.Query.GetHasMore(total)

	return &resp, err
}

func (o orderQueryService) GetDetailOrderOfStore(ctx context.Context, dto *store.GetOrderOfStoreByIDRequest) (*store.GetOrderOfStoreByIDResponse, error) {
	orderResp := store.GetOrderOfStoreByIDResponse{}

	orderDAO, err := o.orderRepo.FindByID(ctx, dto.OrderID)
	if err != nil {
		return nil, err
	}

	if orderDAO.StoreId != dto.StoreID {
		return nil, errors.ErrNotFound
	}

	if err = mapper.BindingStruct(orderDAO, &orderResp); err != nil {
		return nil, err
	}

	if err = mapper.BindingStruct(orderDAO.OrderCommission, &orderResp.CommissionDetail); err != nil {
		return nil, err
	}

	storeAmount := 0
	var items []store.OrderStoreItem
	for _, o := range orderDAO.OrderItem {
		i := store.OrderStoreItem{
			ProductId:   o.ProductID,
			OptionId:    o.OptionID,
			Quantity:    o.Quantity,
			Price:       o.Price,
			Id:          o.Id,
			SubTotal:    o.SubTotal,
			NameOption:  o.NameOption,
			ProdImg:     o.ProdImg,
			ProductName: o.ProductName,
			NetPrice:    o.NetPrice,
		}
		items = append(items, i)
		storeAmount += o.SubTotal
	}

	if len(items) < 1 {
		return nil, errors.ErrNotFoundRecord
	}

	orderResp.StoreOrderAmount = storeAmount
	orderResp.OrderItems = items

	return &orderResp, err
}

func (o orderQueryService) CheckProductPurchased(ctx context.Context, dto *orderDTO.CheckUserOrderRequest) (*orderDTO.CheckUserOrderResponse, error) {
	orders, err := o.orderRepo.FindOrderByUserAndProduct(ctx, dto.UserId, dto.ProductId)
	if err != nil {
		return nil, err
	}
	data := orderDTO.CheckUserOrderResponse{}

	if len(orders) > 0 {
		var ordersKeys []string
		for _, i := range orders {
			ordersKeys = append(ordersKeys, i.OrderID)
		}
		data.IsPurchased = true
		data.Orders = ordersKeys
	}

	return &data, err
}
