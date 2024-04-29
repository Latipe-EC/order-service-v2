package order

import (
	"context"
	"fmt"
	gormF "gorm.io/gorm"
	entity "latipe-order-service-v2/internal/domain/entities/order"
	"latipe-order-service-v2/pkg/db/gorm"
	"latipe-order-service-v2/pkg/util/pagable"
	"strings"
)

type GormRepository struct {
	client gorm.Gorm
}

func NewGormRepository(client gorm.Gorm) entity.OrderRepository {
	// auto migrate
	err := client.DB().AutoMigrate(
		&entity.Order{},
		&entity.OrderItem{},
		&entity.OrderStatusLog{},
		&entity.DeliveryOrder{},
		&entity.OrderCommission{},
	)
	if err != nil {
		panic(err)
	}
	return &GormRepository{
		client: client,
	}
}

func (g GormRepository) FindByItemId(ctx context.Context, itemId string) (*entity.OrderItem, error) {
	item := entity.OrderItem{}

	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Model(&entity.OrderItem{}).
			First(&item, "id = ?", itemId).Error
	}, ctx)
	if result != nil {
		return nil, result
	}

	return &item, nil
}

func (g GormRepository) FindByID(ctx context.Context, orderId string) (*entity.Order, error) {
	order := entity.Order{}

	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Model(&entity.Order{}).
			Preload("OrderItem").
			Preload("Delivery").
			Preload("OrderStatusLog").
			Preload("OrderCommission").
			First(&order, "order_id = ?", orderId).Error
	}, ctx)
	if result != nil {
		return nil, result
	}

	return &order, nil
}

func (g GormRepository) FindByIdForUpdate(ctx context.Context, orderId string) (*entity.Order, error) {
	order := entity.Order{}

	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Model(&entity.Order{}).
			Preload("Delivery").
			First(&order, "order_id = ?", orderId).Error
	}, ctx)
	if result != nil {
		return nil, result
	}

	return &order, nil
}

func (g GormRepository) FindByIdSingleObject(ctx context.Context, orderId string) (*entity.Order, error) {
	order := entity.Order{}

	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Model(&entity.Order{}).
			First(&order, "order_id = ?", orderId).Error
	}, ctx)
	if result != nil {
		return nil, result
	}

	return &order, nil
}

func (g GormRepository) FindAll(ctx context.Context, query *pagable.Query) ([]entity.Order, error) {
	var orders []entity.Order
	whereState := query.ORMConditions().(string)

	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Model(&entity.Order{}).
			Preload("OrderItem").
			Preload("Delivery").
			Where(whereState).
			Limit(query.GetLimit()).Offset(query.GetOffset()).
			Find(&orders).Error
	}, ctx)

	if result != nil {
		return nil, result
	}

	return orders, nil
}

func (g GormRepository) FindByUserId(ctx context.Context, userId string, query *pagable.Query) ([]entity.Order, error) {
	var orders []entity.Order

	whereState := query.UserORMConditions().(string)

	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Model(&entity.Order{}).
			Preload("Delivery").
			Where(whereState).
			Where("orders.user_id", userId).
			Order("created_at desc").
			Limit(query.GetLimit()).Offset(query.GetOffset()).
			Find(&orders).Error
	}, ctx)

	if result != nil {
		return nil, result
	}

	return orders, nil
}

func (g GormRepository) FindOrderByStoreID(ctx context.Context, storeId string, query *pagable.Query, keyword string) ([]entity.Order, error) {
	var orders []entity.Order
	whereState := query.UserORMConditions().(string)
	var likeState string

	if strings.Contains(whereState, "status") {
		whereState = strings.Replace(whereState, "status", "orders.status", 1)
	}

	if len(keyword) > 2 {
		likeState = fmt.Sprintf("orders.order_id like %v", fmt.Sprintf("'%%%v%%'", keyword))
	}

	err := g.client.DB().Model(&entity.Order{}).
		Preload("Delivery").
		Where("orders.store_id=?", storeId).
		Where(whereState).
		Where(likeState).
		Limit(query.GetLimit()).Offset(query.GetOffset()).
		Find(&orders).Error

	if err != nil {
		return nil, err
	}

	return orders, err
}

func (g GormRepository) SearchOrderByStoreID(ctx context.Context, storeId string, keyword string, query *pagable.Query) ([]entity.Order, error) {
	var orders []entity.Order
	err := g.client.DB().Model(&entity.Order{}).
		Preload("Delivery").
		Where("orders.store_id=?", storeId).
		Where("orders.order_id like ?", fmt.Sprintf("'%%%v%%'", keyword)).
		Limit(query.GetLimit()).Offset(query.GetOffset()).
		Find(&orders).Error

	if err != nil {
		return nil, err
	}

	return orders, err
}

func (g GormRepository) TotalSearchOrderByStoreID(ctx context.Context, storeId string, keyword string) (int, error) {
	var count int64

	err := g.client.DB().Select("*").Model(&entity.Order{}).
		Where("orders.store_id=?", storeId).
		Where("orders.order_id like ?", fmt.Sprintf("'%%%v%%'", keyword)).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return int(count), err
}

func (g GormRepository) FindOrderByDelivery(ctx context.Context, deliID string, keyword string, query *pagable.Query) ([]entity.Order, error) {
	var orders []entity.Order
	var searchState string

	whereState := query.UserORMConditions().(string)
	if strings.Contains(whereState, "status") {
		whereState = strings.Replace(whereState, "status", "orders.status", 1)
	}

	if len(keyword) > 2 {
		searchState = fmt.Sprintf("orders.order_id like %v", fmt.Sprintf("'%%%v%%'", keyword))
	}

	err := g.client.DB().Model(&entity.Order{}).Preload("Delivery").
		Joins("inner join delivery_orders ON orders.order_id = delivery_orders.order_id").
		Where("delivery_orders.delivery_id=?", deliID).
		Order("orders.created_at DESC").
		Where(searchState).
		Where(whereState).
		Limit(query.GetLimit()).Offset(query.GetOffset()).
		Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, err
}

func (g GormRepository) FindOrderByUserAndProduct(ctx context.Context, userId string, productId string) ([]entity.Order, error) {
	var orders []entity.Order
	err := g.client.DB().
		Raw("select * from orders "+
			"inner join order_items on orders.order_id = order_items.order_id"+
			"where orders.user_id= ? and order_items.product_id = ?", userId, productId).Scan(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, err
}

func (g GormRepository) FindOrderLogByOrderId(ctx context.Context, orderId string) ([]entity.OrderStatusLog, error) {
	var orderStatus []entity.OrderStatusLog
	result := g.client.DB().Model(&entity.OrderStatusLog{}).
		Where("order_status_logs.order_id", orderId).
		Order("created_at desc").
		Find(&orderStatus).Error
	if result != nil {
		return nil, result
	}

	return orderStatus, nil
}

func (g GormRepository) Save(ctx context.Context, dao *entity.Order) error {
	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Create(&dao).Error
	}, ctx)

	return result
}

func (g GormRepository) SaveMultiple(ctx context.Context, dao []*entity.Order) error {
	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Save(&dao).Error
	}, ctx)

	return result
}

func (g GormRepository) UpdateStatus(ctx context.Context, orderID string, status int, message ...string) error {

	updateLog := entity.OrderStatusLog{
		OrderID:      orderID,
		OrderType:    "orders",
		StatusChange: status,
	}

	if len(message) == 0 {
		switch status {
		case entity.ORDER_PREPARED:
			updateLog.Message = "Đơn hàng đã chuẩn bị bởi nhà bán hàng"
		case entity.ORDER_DELIVERY:
			updateLog.Message = "Đơn hàng đang được vận chuyển"
		case entity.ORDER_SHIPPING_FINISH:
			updateLog.Message = "Đơn hàng được giao thành công"
		}
	} else {
		updateLog.Message = message[0]
	}

	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Transaction(func(tx *gormF.DB) error {
			if err := tx.Model(&entity.Order{}).
				Where("order_id = ?", orderID).Update("status", status).Error; err != nil {
				return err
			}

			if err := tx.Create(&updateLog).Error; err != nil {
				return err
			}
			return nil
		})
	}, ctx)

	if result != nil {
		return result
	}

	return nil
}

func (g GormRepository) Update(ctx context.Context, order entity.Order) error {
	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Updates(order).Error
	}, ctx)

	if result != nil {
		return result
	}
	return nil
}

func (g GormRepository) TotalOrdersOfDelivery(ctx context.Context, deliveryId string, keyword string, query *pagable.Query) (int, error) {
	var count int64
	var searchState string

	whereState := query.UserORMConditions().(string)
	if strings.Contains(whereState, "status") {
		whereState = strings.Replace(whereState, "status", "orders.status", 1)
	}

	if len(keyword) > 2 {
		searchState = fmt.Sprintf("orders.order_id like %v", fmt.Sprintf("'%%%v%%'", keyword))
	}

	err := g.client.DB().Select("*").Model(&entity.Order{}).
		Joins("inner join delivery_orders ON orders.order_id = delivery_orders.order_id").
		Where(searchState).
		Where(whereState).
		Where("delivery_orders.delivery_id=?", deliveryId).
		Order("orders.created_at DESC").
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (g GormRepository) Total(ctx context.Context, query *pagable.Query) (int, error) {
	var count int64
	whereState := query.ORMConditions().(string)
	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Select("*").Table(entity.Order{}.TableName()).
			Where(whereState).
			Count(&count).Error
	}, ctx)

	return int(count), result
}

func (g GormRepository) UserQueryTotal(ctx context.Context, userId string, query *pagable.Query) (int, error) {
	var count int64
	whereState := query.UserORMConditions().(string)
	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Select("*").Table(entity.Order{}.TableName()).
			Where("orders.user_id", userId).
			Where(whereState).
			Count(&count).Error
	}, ctx)

	return int(count), result
}

func (g GormRepository) TotalStoreOrder(ctx context.Context, storeId string, query *pagable.Query, keyword string) (int, error) {
	var count int64
	var likeState string

	whereState := query.UserORMConditions().(string)
	if strings.Contains(whereState, "status") {
		whereState = strings.Replace(whereState, "status", "orders.status", 1)
	}

	if len(keyword) > 2 {
		likeState = fmt.Sprintf("orders.order_id like %v", fmt.Sprintf("'%%%v%%'", keyword))
	}

	err := g.client.DB().Select("*").Model(&entity.Order{}).
		Where("orders.store_id=?", storeId).
		Where(whereState).
		Where(likeState).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (g GormRepository) UpdateOrderItem(ctx context.Context, orderItemID string, status int) error {
	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Model(&entity.OrderItem{}).
			Where("id = ?", orderItemID).Update("status", status).Error
	}, ctx)

	if result != nil {
		return result
	}
	return nil
}

func (g GormRepository) UpdateOrderRating(ctx context.Context, itemId string, ratingId string) error {
	result := g.client.Exec(func(tx *gormF.DB) error {
		return tx.Model(&entity.OrderItem{}).
			Where("id = ?", itemId).Update("rating_id", ratingId).Error
	}, ctx)

	if result != nil {
		return result
	}

	return nil
}
