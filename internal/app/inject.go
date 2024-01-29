package app

import (
	"github.com/google/wire"
	"latipe-order-service-v2/internal/app/commands/ordercommand"
	"latipe-order-service-v2/internal/app/queries/orderquery"
	"latipe-order-service-v2/internal/app/queries/orderstatistic"
)

var Set = wire.NewSet(
	ordercommand.NewOrderCommmandService,
	orderquery.NewOrderQueryService,
	orderstatistic.NewOrderStatisicService,
)
