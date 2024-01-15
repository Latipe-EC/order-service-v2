package app

import (
	"github.com/google/wire"
	"latipe-order-service-v2/internal/app/orders"
)

var Set = wire.NewSet(
	orders.NewOrderService,
)
