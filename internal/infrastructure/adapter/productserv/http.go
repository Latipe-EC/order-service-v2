package productserv

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/wire"
	"latipe-order-service-v2/config"
	"latipe-order-service-v2/internal/common/errors"
	productDTO "latipe-order-service-v2/internal/infrastructure/adapter/productserv/dto"
	"latipe-order-service-v2/pkg/internal_http"
	"latipe-order-service-v2/pkg/util/mapper"
)

var Set = wire.NewSet(
	NewProductServAdapter,
)

type httpAdapter struct {
	client http.Client
}

func NewProductServAdapter(config *config.Config) Service {
	restyClient := http.New()
	restyClient.SetRestyClient(
		restyClient.
			Resty().
			SetBaseURL(config.AdapterService.ProductService.BaseURL).
			SetHeader("X-INTERNAL-SERVICE", config.AdapterService.ProductService.InternalKey))
	return httpAdapter{
		client: restyClient,
	}
}

func (h httpAdapter) GetProductOrderInfo(ctx context.Context, req *productDTO.OrderProductRequest) (*productDTO.OrderProductResponse, error) {
	resp, err := h.client.MakeRequest().
		SetBody(req.StoreOrders).
		SetContext(ctx).
		Post(req.URL())

	if err != nil {
		log.Errorf("[Get product]: %s", err)
		return nil, err
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[Get product]: %s", resp.Body())
		return nil, errors.ErrInternalServer
	}

	if resp.StatusCode() >= 400 {
		log.Errorf("[Get product]: %s", resp.Body())
		return nil, errors.ErrBadRequest
	}

	var rawResp productDTO.BaseResponse
	if err := json.Unmarshal(resp.Body(), &rawResp.Data); err != nil {
		log.Errorf("[Get product]: %s", err)
		return nil, errors.ErrInternalServer
	}

	var regResp *productDTO.OrderProductResponse

	err = mapper.BindingStruct(rawResp.Data, &regResp)
	if err != nil {
		log.Errorf(" [Get product]: %s", err)
		return nil, errors.ErrInternalServer
	}

	return regResp, nil
}

func (h httpAdapter) ReduceProductQuantity(ctx context.Context, req *productDTO.ReduceProductRequest) (*productDTO.ReduceProductResponse, error) {
	resp, err := h.client.MakeRequest().
		SetBody(req.Items).
		SetContext(ctx).
		Patch(req.URL())

	if err != nil {
		log.Errorf("[Reduce Quantity]: %s", err)
		return nil, err
	}

	if resp.StatusCode() >= 400 {
		log.Errorf("[Get product]: %s", resp.Body())
		return nil, errors.ErrBadRequest
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[Reduce Quantity]: %s", resp.Body())
		return nil, err
	}

	var rawResp productDTO.BaseResponse
	if err := json.Unmarshal(resp.Body(), &rawResp); err != nil {
		log.Errorf("[Reduce Quantity]: %s", err)
		return nil, errors.ErrInternalServer
	}

	/*	if rawResp.Code != 0 && resp.StatusCode() != 200 {
		return nil, errors.ErrorMapping(baseResp.Code)
	}*/
	var regResp *productDTO.ReduceProductResponse
	err = mapper.BindingStruct(rawResp.Data, &regResp)
	if err != nil {
		log.Errorf("[%s] [Reduce Quantity]: %s", "ERROR", err)
		return nil, errors.ErrInternalServer
	}

	return regResp, nil
}

func (h httpAdapter) RollBackQuantityOrder(ctx context.Context, req *productDTO.RollbackQuantityRequest) (*productDTO.RollbackQuantityResponse, error) {
	resp, err := h.client.MakeRequest().
		SetBody(req).
		SetContext(ctx).
		Patch(req.URL())

	if err != nil {
		log.Errorf("[%s] [Reduce Quantity]: %s", "ERROR", err)
		return nil, err
	}

	if resp.StatusCode() >= 500 {
		log.Errorf("[%s] [Reduce Quantity]: %s", "ERROR", resp.Body())
		return nil, err
	}

	var rawResp productDTO.BaseResponse
	if err := json.Unmarshal(resp.Body(), &rawResp); err != nil {
		log.Errorf("[%s] [Reduce Quantity]: %s", "ERROR", err)
		return nil, errors.ErrInternalServer
	}

	/*	if rawResp.Code != 0 && resp.StatusCode() != 200 {
		return nil, errors.ErrorMapping(baseResp.Code)
	}*/
	var regResp *productDTO.RollbackQuantityResponse
	err = mapper.BindingStruct(rawResp.Data, &regResp)
	if err != nil {
		log.Errorf("[%s] [Reduce Quantity]: %s", "ERROR", err)
		return nil, errors.ErrInternalServer
	}

	return regResp, nil
}
