package api

import (
	"github.com/google/wire"
	"latipe-order-service-v2/internal/api/order"
)

var Set = wire.NewSet(
	order.NewOrderHandler,
	order.NewStatisticHandler,
)
