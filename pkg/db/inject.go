package db

import (
	"github.com/google/wire"
	"latipe-order-service-v2/pkg/db/gorm"
)

var Set = wire.NewSet(gorm.NewCacheGormPlugin)
