package order

import (
	"context"
	"github.com/stretchr/testify/mock"
	"latipe-order-service-v2/internal/domain/dto/custom_entity"
	"latipe-order-service-v2/internal/domain/entities/order"
	"latipe-order-service-v2/pkg/util/pagable"
)

type OrderMockRepos struct {
	mock.Mock
}

func (o OrderMockRepos) FindByItemId(ctx context.Context, itemId string) (*order.OrderItem, error) {
	args := o.Called(ctx, itemId)
	return args.Get(0).(*order.OrderItem), args.Error(1)
}

func (o OrderMockRepos) FindByID(ctx context.Context, Id string) (*order.Order, error) {
	args := o.Called(ctx, Id)
	return args.Get(0).(*order.Order), args.Error(1)
}

func (o OrderMockRepos) FindOrderByStoreID(ctx context.Context, storeId string, query *pagable.Query, keyword string) ([]order.Order, error) {
	args := o.Called(ctx, storeId, query, keyword)
	return args.Get(0).([]order.Order), args.Error(1)
}

func (o OrderMockRepos) FindOrderByDelivery(ctx context.Context, deliID string, keyword string, query *pagable.Query) ([]order.Order, error) {
	args := o.Called(ctx, deliID, query, keyword)
	return args.Get(0).([]order.Order), args.Error(1)
}

func (o OrderMockRepos) FindAll(ctx context.Context, query *pagable.Query) ([]order.Order, error) {
	args := o.Called(ctx, query)
	return args.Get(0).([]order.Order), args.Error(1)
}

func (o OrderMockRepos) FindByUserId(ctx context.Context, userId string, query *pagable.Query) ([]order.Order, error) {
	args := o.Called(ctx, userId, query)
	return args.Get(0).([]order.Order), args.Error(1)
}

func (o OrderMockRepos) SearchOrderByStoreID(ctx context.Context, storeId string, keyword string, query *pagable.Query) ([]order.Order, error) {
	args := o.Called(ctx, storeId, query, keyword)
	return args.Get(0).([]order.Order), args.Error(1)
}

func (o OrderMockRepos) FindOrderLogByOrderId(ctx context.Context, orderId string) ([]order.OrderStatusLog, error) {
	args := o.Called(ctx, orderId)
	return args.Get(0).([]order.OrderStatusLog), args.Error(1)
}

func (o OrderMockRepos) FindOrderByUserAndProduct(ctx context.Context, userId string, productId string) ([]order.Order, error) {
	args := o.Called(ctx, userId, productId)
	return args.Get(0).([]order.Order), args.Error(1)
}

func (o OrderMockRepos) GetOrderAmountOfStore(ctx context.Context, orderId string) ([]custom_entity.AmountItemOfStoreInOrder, error) {
	args := o.Called(ctx, orderId)
	return args.Get(0).([]custom_entity.AmountItemOfStoreInOrder), args.Error(1)
}

func (o OrderMockRepos) Save(ctx context.Context, order *order.Order) error {
	args := o.Called(ctx, order)
	return args.Error(0)
}

func (o OrderMockRepos) Update(ctx context.Context, order order.Order) error {
	args := o.Called(ctx, order)
	return args.Error(0)
}

func (o OrderMockRepos) UpdateStatus(ctx context.Context, orderId string, status int, message ...string) error {
	args := o.Called(ctx, orderId, status, message)
	return args.Error(0)
}

func (o OrderMockRepos) UpdateOrderItem(ctx context.Context, orderItem string, status int) error {
	args := o.Called(ctx, orderItem, status)
	return args.Error(0)
}

func (o OrderMockRepos) Total(ctx context.Context, query *pagable.Query) (int, error) {
	args := o.Called(ctx, query)
	return args.Get(0).(int), args.Error(1)
}

func (o OrderMockRepos) UserQueryTotal(ctx context.Context, userId string, query *pagable.Query) (int, error) {
	args := o.Called(ctx, userId, query)
	return args.Get(0).(int), args.Error(1)
}

func (o OrderMockRepos) TotalStoreOrder(ctx context.Context, storeId string, query *pagable.Query, keyword string) (int, error) {
	args := o.Called(ctx, storeId, query, keyword)
	return args.Get(0).(int), args.Error(1)
}

func (o OrderMockRepos) TotalOrdersOfDelivery(ctx context.Context, deliveryId string, keyword string, query *pagable.Query) (int, error) {
	args := o.Called(ctx, deliveryId, query, keyword)
	return args.Get(0).(int), args.Error(1)
}

func (o OrderMockRepos) TotalSearchOrderByStoreID(ctx context.Context, storeId string, keyword string) (int, error) {
	args := o.Called(ctx, storeId, keyword)
	return args.Get(0).(int), args.Error(1)
}

func (o OrderMockRepos) GetTotalOrderInSystemInDay(ctx context.Context, date string) ([]custom_entity.TotalOrderInSystemInHours, error) {
	args := o.Called(ctx, date)
	return args.Get(0).([]custom_entity.TotalOrderInSystemInHours), args.Error(1)
}

func (o OrderMockRepos) GetTotalOrderInSystemInMonth(ctx context.Context, date string) ([]custom_entity.TotalOrderInSystemInDay, error) {
	args := o.Called(ctx, date)
	return args.Get(0).([]custom_entity.TotalOrderInSystemInDay), args.Error(1)
}

func (o OrderMockRepos) GetTotalOrderInSystemInYear(ctx context.Context, year int) ([]custom_entity.TotalOrderInSystemInMonth, error) {
	args := o.Called(ctx, year)
	return args.Get(0).([]custom_entity.TotalOrderInSystemInMonth), args.Error(1)
}

func (o OrderMockRepos) GetTotalCommissionOrderInYear(ctx context.Context, date string) ([]custom_entity.SystemOrderCommissionDetail, error) {
	args := o.Called(ctx, date)
	return args.Get(0).([]custom_entity.SystemOrderCommissionDetail), args.Error(1)
}

func (o OrderMockRepos) TopOfProductSold(ctx context.Context, date string, count int) ([]custom_entity.TopOfProductSold, error) {
	args := o.Called(ctx, date, count)
	return args.Get(0).([]custom_entity.TopOfProductSold), args.Error(1)
}

func (o OrderMockRepos) GetTotalOrderInSystemInMonthOfStore(ctx context.Context, date string, storeId string) ([]custom_entity.TotalOrderInSystemInDay, error) {
	args := o.Called(ctx, date, storeId)
	return args.Get(0).([]custom_entity.TotalOrderInSystemInDay), args.Error(1)
}

func (o OrderMockRepos) GetTotalOrderInSystemInYearOfStore(ctx context.Context, year int, storeId string) ([]custom_entity.TotalOrderInSystemInMonth, error) {
	args := o.Called(ctx, year, storeId)
	return args.Get(0).([]custom_entity.TotalOrderInSystemInMonth), args.Error(1)
}

func (o OrderMockRepos) GetTotalCommissionOrderInYearOfStore(ctx context.Context, date string, storeId string) ([]custom_entity.OrderCommissionDetail, error) {
	args := o.Called(ctx, date, storeId)
	return args.Get(0).([]custom_entity.OrderCommissionDetail), args.Error(1)
}

func (o OrderMockRepos) TopOfProductSoldOfStore(ctx context.Context, date string, count int, storeId string) ([]custom_entity.TopOfProductSold, error) {
	args := o.Called(ctx, date, count, storeId)
	return args.Get(0).([]custom_entity.TopOfProductSold), args.Error(1)
}

func (o OrderMockRepos) UserCountingOrder(ctx context.Context, userId string) (int, error) {
	args := o.Called(ctx, userId)
	return args.Get(0).(int), args.Error(1)
}

func (o OrderMockRepos) StoreCountingOrder(ctx context.Context, storeId string) (int, error) {
	args := o.Called(ctx, storeId)
	return args.Get(0).(int), args.Error(1)
}

func (o OrderMockRepos) DeliveryCountingOrder(ctx context.Context, deliveryId string) (int, error) {
	args := o.Called(ctx, deliveryId)
	return args.Get(0).(int), args.Error(1)
}

func (o OrderMockRepos) AdminCountingOrder(ctx context.Context) (int, error) {
	args := o.Called(ctx)
	return args.Get(0).(int), args.Error(1)
}
