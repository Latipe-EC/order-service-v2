package router

import (
	"github.com/google/wire"
	adminRouter "latipe-order-service-v2/internal/router/admin"
	deliveryRouter "latipe-order-service-v2/internal/router/delivery"
	"latipe-order-service-v2/internal/router/internalRouter"

	statisticRouter "latipe-order-service-v2/internal/router/statistic"
	storeRouter "latipe-order-service-v2/internal/router/store"
	userRouter "latipe-order-service-v2/internal/router/user"
)

var Set = wire.NewSet(
	adminRouter.NewAdminOrderRouter,
	userRouter.NewUserOrderRouter,
	storeRouter.NewStoreOrderRouter,
	deliveryRouter.NewDeliveryOrderRouter,
	statisticRouter.NewStatisticOrderRouter,
	internalRouter.NewInternalOrderRouter,
)
