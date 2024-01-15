package persistence

import (
	"github.com/google/wire"
	"latipe-order-service-v2/internal/infrastructure/persistence/db"
	"latipe-order-service-v2/internal/infrastructure/persistence/order"
)

var Set = wire.NewSet(
	db.NewMySQLConnection,
	order.NewGormRepository,
)
