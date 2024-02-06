package transaction

import (
	"context"
	gormF "gorm.io/gorm"
)

type TransactionRepository interface {
	StartTransaction(ctx context.Context, commandFunc ...func(*gormF.DB) error) error
}
