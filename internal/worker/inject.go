package worker

import "github.com/google/wire"

var Set = wire.NewSet(
	NewPurchaseReplySubscriber,
)
