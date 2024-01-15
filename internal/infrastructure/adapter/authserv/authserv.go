package authserv

import (
	"context"
	"latipe-order-service-v2/internal/infrastructure/adapter/authserv/dto"
)

type Service interface {
	Authorization(ctx context.Context, req *dto.AuthorizationRequest) (*dto.AuthorizationResponse, error)
}
