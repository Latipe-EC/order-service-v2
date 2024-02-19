package subscriber

import (
	"github.com/google/wire"
	"latipe-order-service-v2/internal/subscriber/purchase"
	"latipe-order-service-v2/internal/subscriber/rating"
)

var Set = wire.NewSet(
	purchase.NewPurchaseReplySubscriber,
	rating.NewRatingItemSubscriber,
)
