package userserv

import (
	"context"
	"latipe-order-service-v2/internal/infrastructure/adapter/userserv/dto"
)

type Service interface {
	GetAddressDetails(ctx context.Context, req *dto.GetDetailAddressRequest) (*dto.GetDetailAddressResponse, error)
}
