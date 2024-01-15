package storeserv

import (
	"context"
	"github.com/stretchr/testify/mock"
	storeDTO "latipe-order-service-v2/internal/infrastructure/adapter/storeserv/dto"
)

type StoreServiceMock struct {
	mock.Mock
}

func (s *StoreServiceMock) GetStoreByUserId(ctx context.Context, req *storeDTO.GetStoreIdByUserRequest) (*storeDTO.GetStoreIdByUserResponse, error) {
	args := s.Called(ctx, req)
	return args.Get(0).(*storeDTO.GetStoreIdByUserResponse), args.Error(1)
}
