package order

import (
	"context"
	"latipe-order-service-v2/internal/domain/dto/custom_entity"
	"latipe-order-service-v2/pkg/util/pagable"
)

type OrderRepository interface {
	FindByItemId(ctx context.Context, itemId string) (*OrderItem, error)
	FindByID(ctx context.Context, orderId string) (*Order, error)
	FindByIdForUpdate(ctx context.Context, orderId string) (*Order, error)
	FindByIdSingleObject(ctx context.Context, orderId string) (*Order, error)
	FindOrderByStoreID(ctx context.Context, storeId string, query *pagable.Query, keyword string) ([]Order, error)
	FindOrderByDelivery(ctx context.Context, deliID string, keyword string, query *pagable.Query) ([]Order, error)
	FindAll(ctx context.Context, query *pagable.Query) ([]Order, error)
	FindByUserId(ctx context.Context, userId string, query *pagable.Query) ([]Order, error)
	SearchOrderByStoreID(ctx context.Context, storeId string, keyword string, query *pagable.Query) ([]Order, error)
	FindOrderLogByOrderId(ctx context.Context, orderId string) ([]OrderStatusLog, error)
	FindOrderByUserAndProduct(ctx context.Context, userId string, productId string) ([]Order, error)
	GetOrderAmountOfStore(ctx context.Context, orderId string) ([]custom_entity.AmountItemOfStoreInOrder, error)
	Save(ctx context.Context, order *Order) error
	SaveMultiple(ctx context.Context, dao []*Order) error

	Update(ctx context.Context, order Order) error
	UpdateStatus(ctx context.Context, orderID string, status int, message ...string) error
	UpdateOrderRating(ctx context.Context, itemId string, ratingId string) error
	UpdateOrderItem(ctx context.Context, orderItem string, status int) error
	Total(ctx context.Context, query *pagable.Query) (int, error)
	UserQueryTotal(ctx context.Context, userId string, query *pagable.Query) (int, error)
	TotalStoreOrder(ctx context.Context, storeId string, query *pagable.Query, keyword string) (int, error)
	TotalOrdersOfDelivery(ctx context.Context, deliveryId string, keyword string, query *pagable.Query) (int, error)
	TotalSearchOrderByStoreID(ctx context.Context, storeId string, keyword string) (int, error)
	//custom_entity - admin
	GetTotalOrderInSystemInDay(ctx context.Context, date string) ([]custom_entity.TotalOrderInSystemInHours, error)
	GetTotalOrderInSystemInMonth(ctx context.Context, date string) ([]custom_entity.TotalOrderInSystemInDay, error)
	GetTotalOrderInSystemInYear(ctx context.Context, year int) ([]custom_entity.TotalOrderInSystemInMonth, error)
	GetTotalCommissionOrderInYear(ctx context.Context, date string) ([]custom_entity.SystemOrderCommissionDetail, error)
	TopOfProductSold(ctx context.Context, date string, count int) ([]custom_entity.TopOfProductSold, error)

	//custom_entity - store
	GetTotalOrderInSystemInMonthOfStore(ctx context.Context, date string, storeId string) ([]custom_entity.TotalOrderInSystemInDay, error)
	GetTotalOrderInSystemInYearOfStore(ctx context.Context, year int, storeId string) ([]custom_entity.TotalOrderInSystemInMonth, error)
	GetTotalCommissionOrderInYearOfStore(ctx context.Context, date string, storeId string) ([]custom_entity.OrderCommissionDetail, error)
	TopOfProductSoldOfStore(ctx context.Context, date string, count int, storeId string) ([]custom_entity.TopOfProductSold, error)

	UserCountingOrder(ctx context.Context, userId string) (int, error)
	StoreCountingOrder(ctx context.Context, storeId string) (int, error)
	DeliveryCountingOrder(ctx context.Context, deliveryId string) (int, error)
	AdminCountingOrder(ctx context.Context) (int, error)

	GetStoreRevenueDistributionInMonth(ctx context.Context, storeId string, date string) (*custom_entity.StoreRevenuePer, error)
	AdminRevenueDistributionInMonth(ctx context.Context, date string) (*custom_entity.AdminRevenuePer, error)
	GetAllOrderDataRecordByAdmin(ctx context.Context, date string) ([]custom_entity.StatisticOrderRecordData, error)
	GetAllOrderDataRecordByStore(ctx context.Context, storeId string, date string) ([]custom_entity.StatisticOrderRecordData, error)
}
