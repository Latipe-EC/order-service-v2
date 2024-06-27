package persistence

import (
	"github.com/google/wire"
	"latipe-order-service-v2/internal/infrastructure/excel"
	"latipe-order-service-v2/internal/infrastructure/persistence/commission"
	"latipe-order-service-v2/internal/infrastructure/persistence/db"
	"latipe-order-service-v2/internal/infrastructure/persistence/order"
	"latipe-order-service-v2/internal/infrastructure/persistence/transaction"
)

var Set = wire.NewSet(
	db.NewMySQLConnection,
	order.NewGormRepository,
	commission.NewCommissionRepository,
	transaction.NewGormRepository,
	excel.NewExcelExportClient,
)
