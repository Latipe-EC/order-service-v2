package services

import (
	"github.com/google/wire"
	"latipe-order-service-v2/internal/services/commands/orderCmd"
	"latipe-order-service-v2/internal/services/queries/orderQuery"
	"latipe-order-service-v2/internal/services/queries/statisticQuery"
)

var Set = wire.NewSet(
	orderCmd.NewOrderCommmandService,
	orderQuery.NewOrderQueryService,
	statisticQuery.NewOrderStatisicService,
)
