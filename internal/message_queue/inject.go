package message_queue

import "github.com/google/wire"

var Set = wire.NewSet(NewTransactionProducer)
