package storeserv

import (
	"context"
	"latipe-order-service-v2/internal/infrastructure/adapter/storeserv/dto"
)

type Service interface {
	GetStoreByUserId(ctx context.Context, req *dto.GetStoreIdByUserRequest) (*dto.GetStoreIdByUserResponse, error)
}
