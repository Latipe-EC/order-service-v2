package publisher

import (
	"github.com/google/wire"
)

var Set = wire.NewSet(NewTransactionProducer,
	NewNotificationMessagePublisher)
