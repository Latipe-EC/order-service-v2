package authserv

import (
	"context"
	"github.com/stretchr/testify/mock"
	"latipe-order-service-v2/internal/infrastructure/adapter/authserv/dto"
)

type AuthServiceMock struct {
	mock.Mock
}

func (a *AuthServiceMock) Authorization(ctx context.Context, req *dto.AuthorizationRequest) (*dto.AuthorizationResponse, error) {
	args := a.Called(ctx, req)
	return args.Get(0).(*dto.AuthorizationResponse), args.Error(1)
}
