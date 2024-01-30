package order

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"latipe-order-service-v2/internal/app/commands/ordercommand"
	"latipe-order-service-v2/internal/app/queries/orderquery"
	"latipe-order-service-v2/internal/common/errors"
	dto "latipe-order-service-v2/internal/domain/dto/order"
	"latipe-order-service-v2/internal/domain/dto/order/delivery"
	internalDTO "latipe-order-service-v2/internal/domain/dto/order/internal-service"
	"latipe-order-service-v2/internal/domain/dto/order/store"
	"latipe-order-service-v2/internal/middleware/auth"
	"latipe-order-service-v2/internal/responses"
	"latipe-order-service-v2/pkg/util/pagable"
	"latipe-order-service-v2/pkg/util/valid"
	"strings"
)

type orderApiHandler struct {
	orderCommandServ ordercommand.OrderCommandUsecase
	orderQueryServ   orderquery.OrderQueryUsecase
}

func NewOrderHandler(orderCommandServ ordercommand.OrderCommandUsecase, orderQueryServ orderquery.OrderQueryUsecase) OrderApiHandler {
	return orderApiHandler{
		orderCommandServ: orderCommandServ,
		orderQueryServ:   orderQueryServ,
	}
}

func (o orderApiHandler) CreateOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	bodyReq := dto.CreateOrderRequest{}

	userId := fmt.Sprintf("%v", ctx.Locals(auth.USER_ID))
	if userId == "" {
		return errors.ErrUnauthenticated
	}

	username := fmt.Sprintf("%v", ctx.Locals(auth.USERNAME))
	if username == "" {
		return errors.ErrUnauthenticated
	}

	bearerToken := fmt.Sprintf("%v", ctx.Locals(auth.BEARER_TOKEN))
	if bearerToken == "" {
		return errors.ErrUnauthenticated
	}

	if err := ctx.BodyParser(&bodyReq); err != nil {
		return errors.ErrInternalServer.WithInternalError(err)
	}

	if err := valid.GetValidator().Validate(bodyReq); err != nil {
		return errors.ErrBadRequest
	}

	bodyReq.Header.BearerToken = bearerToken
	bodyReq.UserRequest.UserId = userId
	bodyReq.UserRequest.Username = username

	dataResp, err := o.orderCommandServ.CreateOrder(context, &bodyReq)
	if err != nil {
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = dataResp
	return resp.JSON(ctx)
}

func (o orderApiHandler) GetMyOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	query, err := pagable.GetQueryFromFiberCtx(ctx)
	if err != nil {
		return errors.ErrBadRequest.WithInternalError(err)
	}

	req := new(dto.GetByUserIdRequest)
	req.Query = query

	userId := fmt.Sprintf("%v", ctx.Locals(auth.USER_ID))
	if userId == "" {
		return errors.ErrUnauthenticated
	}

	req.UserId = userId

	result, err := o.orderQueryServ.GetOrderByUserId(context, req)
	if err != nil {
		return errors.ErrInternalServer
	}

	resp := responses.DefaultSuccess
	resp.Data = result

	return resp.JSON(ctx)
}

func (o orderApiHandler) UserCancelOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := new(dto.CancelOrderRequest)

	if err := ctx.BodyParser(&req); err != nil {
		return errors.ErrInternalServer.WithInternalError(err)
	}

	userId := fmt.Sprintf("%v", ctx.Locals(auth.USER_ID))
	if userId == "" {
		return errors.ErrUnauthenticated
	}

	req.UserId = userId

	err := o.orderCommandServ.UserCancelOrder(context, req)
	if err != nil {
		return err
	}

	resp := responses.DefaultSuccess
	return resp.JSON(ctx)
}

func (o orderApiHandler) UserRefundOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := new(dto.CancelOrderRequest)

	if err := ctx.BodyParser(&req); err != nil {
		return errors.ErrInternalServer.WithInternalError(err)
	}

	userId := fmt.Sprintf("%v", ctx.Locals(auth.USER_ID))
	if userId == "" {
		return errors.ErrUnauthenticated
	}

	req.UserId = userId

	err := o.orderCommandServ.UserRefundOrder(context, req)
	if err != nil {
		return err
	}

	resp := responses.DefaultSuccess
	return resp.JSON(ctx)
}

func (o orderApiHandler) AdminCancelOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := new(dto.CancelOrderRequest)

	if err := ctx.BodyParser(&req); err != nil {
		return errors.ErrInternalServer.WithInternalError(err)
	}

	userId := fmt.Sprintf("%v", ctx.Locals(auth.USER_ID))
	if userId == "" {
		return errors.ErrUnauthenticated
	}

	req.UserId = userId

	err := o.orderCommandServ.AdminCancelOrder(context, req)
	if err != nil {
		return err
	}

	resp := responses.DefaultSuccess
	return resp.JSON(ctx)
}

func (o orderApiHandler) UpdateOrderStatus(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := new(dto.UpdateOrderStatusRequest)

	if err := ctx.BodyParser(&req); err != nil {
		return errors.ErrInternalServer.WithInternalError(err)
	}

	userId := fmt.Sprintf("%v", ctx.Locals(auth.USER_ID))
	if userId == "" {
		return errors.ErrUnauthenticated
	}

	role := fmt.Sprintf("%v", ctx.Locals(auth.ROLE))
	if role == "" {
		return errors.ErrPermissionDenied
	}

	req.UserId = userId
	req.Role = role

	err := o.orderCommandServ.UpdateStatusOrder(context, req)
	if err != nil {
		return errors.ErrInternalServer
	}

	return responses.DefaultSuccess.JSON(ctx)
}

func (o orderApiHandler) ListOfOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	query, err := pagable.GetQueryFromFiberCtx(ctx)
	if err != nil {
		return errors.ErrBadRequest.WithInternalError(err)
	}

	req := new(dto.GetOrderListRequest)
	req.Query = query

	result, err := o.orderQueryServ.GetOrderList(context, req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

func (o orderApiHandler) GetByOrderId(ctx *fiber.Ctx) error {
	context := ctx.Context()
	req := new(dto.GetOrderByIDRequest)

	if err := ctx.ParamsParser(req); err != nil {
		return errors.ErrInternalServer
	}

	result, err := o.orderQueryServ.GetOrderByID(context, req)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return errors.ErrNotFound
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

func (o orderApiHandler) UserGetOrderByID(ctx *fiber.Ctx) error {
	context := ctx.Context()
	req := new(dto.GetOrderByIDRequest)

	if err := ctx.ParamsParser(req); err != nil {
		return errors.ErrInternalServer
	}

	role := fmt.Sprintf("%v", ctx.Locals(auth.ROLE))
	if role == "" {
		return errors.ErrPermissionDenied
	}

	userId := fmt.Sprintf("%v", ctx.Locals(auth.USER_ID))
	if userId == "" {
		return errors.ErrUnauthenticated
	}

	req.OwnerId = userId
	req.Role = role

	result, err := o.orderQueryServ.GetOrderByID(context, req)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return errors.ErrNotFound
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

func (o orderApiHandler) DeliveryGetOrderByID(ctx *fiber.Ctx) error {
	context := ctx.Context()
	req := new(dto.GetOrderByIDRequest)

	if err := ctx.ParamsParser(req); err != nil {
		return errors.ErrInternalServer
	}

	role := fmt.Sprintf("%v", ctx.Locals(auth.ROLE))
	if role == "" {
		return errors.ErrPermissionDenied
	}

	deliID := fmt.Sprintf("%v", ctx.Locals(auth.DELIVERY_ID))
	if deliID == "" {
		return errors.ErrUnauthenticated
	}

	req.OwnerId = deliID
	req.Role = role

	result, err := o.orderQueryServ.GetOrderById(context, req)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return errors.ErrNotFound
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

func (o orderApiHandler) InternalGetOrderByOrderID(ctx *fiber.Ctx) error {
	context := ctx.Context()
	req := internalDTO.GetOrderRatingItemRequest{}

	if err := ctx.ParamsParser(&req); err != nil {
		return errors.ErrInternalServer
	}

	result, err := o.orderQueryServ.InternalGetRatingID(context, &req)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return errors.ErrNotFound
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

func (o orderApiHandler) GetMyStoreOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	query, err := pagable.GetQueryFromFiberCtx(ctx)
	if err != nil {
		return errors.ErrBadRequest.WithInternalError(err)
	}

	req := store.GetStoreOrderRequest{}
	req.Query = query

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	storeID := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.StoreID = storeID

	result, err := o.orderQueryServ.GetOrdersOfStore(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

func (o orderApiHandler) GetStoreOrderDetail(ctx *fiber.Ctx) error {
	context := ctx.Context()

	var req store.GetOrderOfStoreByIDRequest

	if err := ctx.ParamsParser(&req); err != nil {
		return errors.ErrBadRequest.WithInternalError(err)
	}

	storeID := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.StoreID = storeID

	result, err := o.orderQueryServ.ViewDetailStoreOrder(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		case errors.Is(err, gorm.ErrRecordNotFound):
			return errors.ErrNotFoundRecord
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

func (o orderApiHandler) UpdateOrderStatusByStore(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := store.StoreUpdateOrderStatusRequest{}

	if err := ctx.BodyParser(&req); err != nil {
		return errors.ErrInternalServer.WithInternalError(err)
	}

	if err := ctx.ParamsParser(&req); err != nil {
		return errors.ErrInternalServer.WithInternalError(err)
	}

	storeId := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	if storeId == "" {
		return errors.ErrUnauthenticated
	}

	req.StoreId = storeId

	err := o.orderCommandServ.StoreUpdateOrderStatus(context, &req)
	if err != nil {
		return err
	}

	data := responses.DefaultSuccess

	return data.JSON(ctx)
}

func (o orderApiHandler) UpdateOrderStatusByDelivery(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := delivery.UpdateOrderStatusRequest{}

	if err := ctx.BodyParser(&req); err != nil {
		return errors.ErrInternalServer.WithInternalError(err)
	}

	if err := ctx.ParamsParser(&req); err != nil {
		return errors.ErrBadRequest.WithInternalError(err)
	}

	deli := fmt.Sprintf("%v", ctx.Locals(auth.DELIVERY_ID))
	if deli == "" {
		return errors.ErrUnauthenticated
	}

	req.DeliveryID = deli

	if err := valid.GetValidator().Validate(req); err != nil {
		return errors.ErrBadRequest
	}

	err := o.orderCommandServ.DeliveryUpdateOrderStatus(context, req)
	if err != nil {
		return err
	}

	data := responses.DefaultSuccess

	return data.JSON(ctx)
}

func (o orderApiHandler) GetOrdersByDelivery(ctx *fiber.Ctx) error {
	context := ctx.Context()

	query, err := pagable.GetQueryFromFiberCtx(ctx)
	if err != nil {
		return errors.ErrBadRequest.WithInternalError(err)
	}

	req := delivery.GetOrderListRequest{}
	req.Query = query

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	deliId := fmt.Sprintf("%v", ctx.Locals(auth.DELIVERY_ID))
	req.DeliveryID = deliId

	result, err := o.orderQueryServ.GetOrdersOfDelivery(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

func (o orderApiHandler) SearchOrderIdByKeyword(ctx *fiber.Ctx) error {
	context := ctx.Context()

	query, err := pagable.GetQueryFromFiberCtx(ctx)
	if err != nil {
		return errors.ErrBadRequest.WithInternalError(err)
	}

	req := store.FindStoreOrderRequest{}
	req.Query = query

	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrInvalidParameters
	}

	if err := valid.GetValidator().Validate(&req); err != nil {
		return errors.ErrBadRequest
	}

	storeID := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.StoreID = storeID

	result, err := o.orderQueryServ.SearchStoreOrderId(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

func (o orderApiHandler) AdminCountingOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := dto.CountingOrderAmountRequest{}
	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrBadRequest
	}

	result, err := o.orderQueryServ.AdminCountingOrderAmount(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}

func (o orderApiHandler) UserCountingOrder(ctx *fiber.Ctx) error {

	context := ctx.Context()

	req := dto.CountingOrderAmountRequest{}
	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrBadRequest
	}

	userId := fmt.Sprintf("%v", ctx.Locals(auth.USER_ID))
	req.OwnerID = userId

	result, err := o.orderQueryServ.UserCountingOrder(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)

}

func (o orderApiHandler) StoreCountingOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := dto.CountingOrderAmountRequest{}
	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrBadRequest
	}

	storeId := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.OwnerID = storeId

	result, err := o.orderQueryServ.StoreCountingOrder(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)

}

func (o orderApiHandler) DeliveryCountingOrder(ctx *fiber.Ctx) error {
	context := ctx.Context()

	req := dto.CountingOrderAmountRequest{}
	if err := ctx.QueryParser(&req); err != nil {
		return errors.ErrBadRequest
	}

	deli := fmt.Sprintf("%v", ctx.Locals(auth.DELIVERY_ID))
	req.OwnerID = deli

	result, err := o.orderQueryServ.DeliveryCountingOrder(context, &req)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Unknown column"):
			return errors.ErrBadRequest.WithInternalError(err)
		}
		return err
	}

	resp := responses.DefaultSuccess
	resp.Data = result
	return resp.JSON(ctx)
}
