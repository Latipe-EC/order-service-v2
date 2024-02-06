package commission

import (
	gormF "gorm.io/gorm"
	entity "latipe-order-service-v2/internal/domain/entities/order"
	"latipe-order-service-v2/pkg/db/gorm"
)

type commissionRepository struct {
	client gorm.Gorm
}

func NewCommissionRepository(client gorm.Gorm) entity.CommissionRepository {
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
	return &commissionRepository{
		client: client,
	}
}

func (g *commissionRepository) UpdateOrderCommission(dao *entity.Order, ocms *entity.OrderCommission, log *entity.OrderStatusLog) error {
	err := g.client.DB().Transaction(func(tx *gormF.DB) error {
		if err := tx.Model(&entity.OrderCommission{}).Where("id=?", ocms.Id).
			Update("status", ocms.Status).Error; err != nil {
			return err
		}

		if err := tx.Model(&entity.Order{}).Where("id=?", dao.OrderID).Update("status", dao.Status).Error; err != nil {
			return err
		}

		if err := tx.Create(&log).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (g *commissionRepository) CreateOrderCommission(ocms *entity.OrderCommission) error {
	err := g.client.DB().Create(&ocms).Error
	if err != nil {
		return err
	}

	return nil
}

func (g *commissionRepository) FindCommissionByOrderId(orderId int) (*entity.OrderCommission, error) {
	var data entity.OrderCommission

	err := g.client.DB().Model(&entity.OrderCommission{}).Where("orders_commission.order_id =?", orderId).
		Find(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (g commissionRepository) UpdateCommission(Id string) (*entity.Order, error) {
	//TODO implement me
	panic("implement me")
}
