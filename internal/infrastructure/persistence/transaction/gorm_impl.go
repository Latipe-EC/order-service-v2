package transaction

import (
	"context"
	gormF "gorm.io/gorm"
	entity "latipe-order-service-v2/internal/domain/entities/order"
	"latipe-order-service-v2/internal/domain/entities/transaction"
	"latipe-order-service-v2/pkg/db/gorm"
)

type transactionRepository struct {
	client gorm.Gorm
}

func NewGormRepository(client gorm.Gorm) transaction.TransactionRepository {
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
	return &transactionRepository{
		client: client,
	}
}

func (t transactionRepository) StartTransaction(ctx context.Context, commandFunc ...func(*gormF.DB) error) error {
	result := t.client.Transaction(func(tx *gormF.DB) error {
		for _, cmd := range commandFunc {
			if err := cmd(tx); err != nil {
				return err
			}
		}
		return nil
	})
	if result != nil {
		return result
	}
	return nil
}
