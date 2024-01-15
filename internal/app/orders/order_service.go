package orders

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"latipe-order-service-v2/config"
	"latipe-order-service-v2/internal/common/errors"
	orderDTO "latipe-order-service-v2/internal/domain/dto/order"
	"latipe-order-service-v2/internal/domain/dto/order/delivery"
	internalDTO "latipe-order-service-v2/internal/domain/dto/order/internal-service"
	"latipe-order-service-v2/internal/domain/dto/order/store"
	"latipe-order-service-v2/internal/domain/entities/order"
	"latipe-order-service-v2/internal/infrastructure/adapter/deliveryserv"
	"latipe-order-service-v2/internal/infrastructure/adapter/productserv"
	"latipe-order-service-v2/internal/infrastructure/adapter/userserv"
	voucherserv "latipe-order-service-v2/internal/infrastructure/adapter/vouchersev"
	publishMsg "latipe-order-service-v2/internal/message_queue"
	"latipe-order-service-v2/internal/middleware/auth"
	"latipe-order-service-v2/pkg/cache/redis"
	"latipe-order-service-v2/pkg/util/mapper"
	"strings"
)

type orderService struct {
	orderRepo   order.Repository
	cacheEngine *redis.CacheEngine
	productServ productserv.Service
	userServ    userserv.Service
	deliServ    deliveryserv.Service
	voucherSer  voucherserv.Service
	publisher   *publishMsg.PublisherTransactionMessage
	cfg         *config.Config
}

func NewOrderService(cfg *config.Config, orderRepo order.Repository, productServ productserv.Service,
	cacheEngine *redis.CacheEngine, userServ userserv.Service, deliServ deliveryserv.Service,
	voucherServ voucherserv.Service, publisher *publishMsg.PublisherTransactionMessage) Usecase {
	return orderService{
		orderRepo:   orderRepo,
		cacheEngine: cacheEngine,
		productServ: productServ,
		userServ:    userServ,
		deliServ:    deliServ,
		voucherSer:  voucherServ,
		publisher:   publisher,
		cfg:         cfg,
	}
}

func (o orderService) InternalGetRatingID(ctx context.Context, dto *internalDTO.GetOrderRatingItemRequest) (*internalDTO.GetOrderRatingItemResponse, error) {
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

func (o orderService) GetOrderById(ctx context.Context, dto *orderDTO.GetOrderByIDRequest) (*orderDTO.GetOrderResponse, error) {
	orderResp := orderDTO.OrderResponse{}

	orderDAO, err := o.orderRepo.FindById(ctx, dto.OrderId)
	if err != nil {
		return nil, err
	}

	switch dto.Role {
	case auth.ROLE_USER:
		if orderDAO.UserId != dto.OwnerId {
			return nil, errors.ErrNotFound
		}
	case auth.ROLE_STORE:
		if !CheckStoreHaveOrder(*orderDAO, dto.OwnerId) {
			return nil, errors.ErrNotFound
		}
	}

	if err = mapper.BindingStruct(orderDAO, &orderResp); err != nil {
		return nil, err
	}

	resp := orderDTO.GetOrderResponse{Order: orderResp}

	return &resp, err
}

func (o orderService) GetOrderByUUID(ctx context.Context, dto *orderDTO.GetOrderByUUIDRequest) (*orderDTO.GetOrderResponse, error) {
	orderResp := orderDTO.OrderResponse{}

	orderDAO, err := o.orderRepo.FindByUUID(ctx, dto.OrderId)
	if err != nil {
		return nil, err
	}

	switch dto.Role {
	case auth.ROLE_USER:
		if orderDAO.UserId != dto.OwnerId {
			return nil, errors.ErrNotFound
		}
	case auth.ROLE_STORE:
		if !CheckStoreHaveOrder(*orderDAO, dto.OwnerId) {
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

func (o orderService) GetOrderList(ctx context.Context, dto *orderDTO.GetOrderListRequest) (*orderDTO.GetOrderListResponse, error) {
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

func (o orderService) GetOrderByUserId(ctx context.Context, dto *orderDTO.GetByUserIdRequest) (*orderDTO.GetByUserIdResponse, error) {
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

func (o orderService) SearchStoreOrderId(ctx context.Context, dto *store.FindStoreOrderRequest) (*orderDTO.GetOrderListResponse, error) {
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

func (o orderService) GetOrdersOfStore(ctx context.Context, dto *store.GetStoreOrderRequest) (*orderDTO.GetOrderListResponse, error) {
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

func (o orderService) GetOrdersOfDelivery(ctx context.Context, dto *delivery.GetOrderListRequest) (*delivery.GetOrderListResponse, error) {
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

func (o orderService) ViewDetailStoreOrder(ctx context.Context, dto *store.GetOrderOfStoreByIDRequest) (*store.GetOrderOfStoreByIDResponse, error) {
	orderResp := store.GetOrderOfStoreByIDResponse{}

	orderDAO, err := o.orderRepo.FindByUUID(ctx, dto.OrderUUID)
	if err != nil {
		return nil, err
	}

	if !CheckStoreHaveOrder(*orderDAO, dto.StoreID) {
		return nil, errors.ErrNotFound
	}

	if err = mapper.BindingStruct(orderDAO, &orderResp); err != nil {
		return nil, err
	}

	storeAmount := 0
	var items []store.OrderStoreItem
	for _, o := range orderDAO.OrderItem {
		if o.StoreID == dto.StoreID {
			i := store.OrderStoreItem{
				ProductId:   o.ProductID,
				OptionId:    o.OptionID,
				Quantity:    o.Quantity,
				Price:       o.Price,
				Status:      o.Status,
				Id:          o.Id,
				SubTotal:    o.SubTotal,
				ProdImg:     o.ProdImg,
				ProductName: o.ProductName,
				NetPrice:    o.NetPrice,
			}
			items = append(items, i)
			storeAmount += o.SubTotal
		}
	}

	if len(items) < 1 {
		return nil, errors.ErrNotFoundRecord
	}

	orderResp.CommissionDetail.SystemFee = orderDAO.OrderCommission.SystemFee
	orderResp.CommissionDetail.AmountReceived = orderDAO.OrderCommission.AmountReceived
	orderResp.StoreOrderAmount = storeAmount
	orderResp.OrderItems = items

	return &orderResp, err
}

func (o orderService) CheckProductPurchased(ctx context.Context, dto *orderDTO.CheckUserOrderRequest) (*orderDTO.CheckUserOrderResponse, error) {
	orders, err := o.orderRepo.FindOrderByUserAndProduct(ctx, dto.UserId, dto.ProductId)
	if err != nil {
		return nil, err
	}
	data := orderDTO.CheckUserOrderResponse{}

	if len(orders) > 0 {
		var ordersKeys []string
		for _, i := range orders {
			ordersKeys = append(ordersKeys, i.OrderUUID)
		}
		data.IsPurchased = true
		data.Orders = ordersKeys
	}

	return &data, err
}

func (o orderService) genOrderKey(userId string) string {
	keyGen := strings.ReplaceAll(uuid.NewString(), "-", "")[:10]
	key := fmt.Sprintf("%v%v%v", o.cfg.Server.KeyID, userId[:4], keyGen)

	return key
}
