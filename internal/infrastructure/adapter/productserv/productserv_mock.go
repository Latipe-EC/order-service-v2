package productserv

import (
	"context"
	"github.com/stretchr/testify/mock"
	"latipe-order-service-v2/internal/infrastructure/adapter/productserv/dto"
)

type ProductServiceMock struct {
	mock.Mock
}

func (d *ProductServiceMock) GetProductOrderInfo(ctx context.Context, req *dto.OrderProductRequest) (*dto.OrderProductResponse, error) {
	args := d.Called(ctx, req)
	return args.Get(0).(*dto.OrderProductResponse), args.Error(1)
}

func (d *ProductServiceMock) ReduceProductQuantity(ctx context.Context, req *dto.ReduceProductRequest) (*dto.ReduceProductResponse, error) {
	args := d.Called(ctx, req)
	return args.Get(0).(*dto.ReduceProductResponse), args.Error(1)
}

func (d *ProductServiceMock) RollBackQuantityOrder(ctx context.Context, req *dto.RollbackQuantityRequest) (*dto.RollbackQuantityResponse, error) {
	args := d.Called(ctx, req)
	return args.Get(0).(*dto.RollbackQuantityResponse), args.Error(1)
}
