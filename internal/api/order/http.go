package order

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	_ "latipe-order-service-v2/docs"
	"latipe-order-service-v2/internal/common/errors"
	"latipe-order-service-v2/internal/common/responses"
	dto "latipe-order-service-v2/internal/domain/dto/order"
	"latipe-order-service-v2/internal/domain/dto/order/delivery"
	internalDTO "latipe-order-service-v2/internal/domain/dto/order/internal-service"
	"latipe-order-service-v2/internal/domain/dto/order/store"
	"latipe-order-service-v2/internal/middleware/auth"
	"latipe-order-service-v2/internal/services/commands/orderCmd"
	"latipe-order-service-v2/internal/services/queries/orderQuery"
	"latipe-order-service-v2/pkg/util/pagable"
	"latipe-order-service-v2/pkg/util/valid"
	"strings"
)

// @title API Documentation
// @version 2.0
// @description This is a server for Latipe Order Service.
// @host localhost:5000
// @BasePath /api/v2/order
// @schemes http
type orderApiHandler struct {
	orderCommandServ orderCmd.OrderCommandUsecase
	orderQueryServ   orderQuery.OrderQueryUsecase
}

func NewOrderHandler(orderCommandServ orderCmd.OrderCommandUsecase, orderQueryServ orderQuery.OrderQueryUsecase) OrderApiHandler {
	return orderApiHandler{
		orderCommandServ: orderCommandServ,
		orderQueryServ:   orderQueryServ,
	}
}

// @Summary Create Order
// @Description Create Order
// @Tags Order
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header"
// @Param CreateOrderRequest body CreateOrderRequest true "Create Order Request"
// @Router /user/order [post]
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
		log.Error(err)
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

// @Summary Get My Order
// @Description Get My Order
// @Tags Order
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Param order query string false "Order"
// @Router /user/order [get]
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

// @Summary User Cancel Order
// @Description User Cancel Order
// @Tags Order
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header"
// @Param CancelOrderRequest body CancelOrderRequest true "Cancel Order Request"
// @Router /user/order/cancel [post]
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

// @Summary User Refund Order
// @Description User Refund Order
// @Tags Order
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header"
// @Param CancelOrderRequest body CancelOrderRequest true "Cancel Order Request"
// @Router /user/order/refund [post]
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

// @Summary Admin Cancel Order
// @Description Admin Cancel Order
// @Tags Order
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header"
// @Param CancelOrderRequest body CancelOrderRequest true "Cancel Order Request"
// @Router /admin/order/cancel [post]
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

// @Summary Update Order Status
// @Description Update Order Status
// @Tags Order
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header"
// @Param UpdateOrderStatusRequest body UpdateOrderStatusRequest true "Update Order Status Request"
// @Router /role/order/status [put]
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

// @Summary List Of Order
// @Description List Of Order
// @Tags Order
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Param order query string false "Order"
// @Router /role/order [get]
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

// @Summary Get By Order ID
// @Description Get By Order ID
// @Tags Order
// @Accept json
// @Produce json
// @Param orderId path string true "Order ID"
// @Router /role/order/{orderId} [get]
func (o orderApiHandler) GetOrderDetailByAdmin(ctx *fiber.Ctx) error {
	context := ctx.Context()
	req := new(dto.GetOrderByIDRequest)

	if err := ctx.ParamsParser(req); err != nil {
		return errors.ErrInternalServer
	}

	result, err := o.orderQueryServ.GetOrderByIdofAdmin(context, req)
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

// @Summary User Get Order By ID
// @Description User Get Order By ID
// @Tags Order
// @Accept json
// @Produce json
// @Param orderId path string true "Order ID"
// @Router /user/order/{orderId} [get]
func (o orderApiHandler) GetOrderDetailOfUser(ctx *fiber.Ctx) error {
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

	result, err := o.orderQueryServ.GetOrderByIdOfUser(context, req)
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

// @Summary Delivery Get Order By ID
// @Description Delivery Get Order By ID
// @Tags Order
// @Accept json
// @Produce json
// @Param orderId path string true "Order ID"
// @Router /delivery/order/{orderId} [get]
func (o orderApiHandler) GetOrderDetailByDelivery(ctx *fiber.Ctx) error {
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

	result, err := o.orderQueryServ.GetOrderDetailOfDelivery(context, req)
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

// @Summary Internal Get Order By Order ID
// @Description Internal Get Order By Order ID
// @Tags Order
// @Accept json
// @Produce json
// @Param orderId path string true "Order ID"
// @Router /internal/order/{orderId} [get]
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

// @Summary Get My Store Order
// @Description Get My Store Order
// @Tags Order
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Param order query string false "Order"
// @Router /store/order [get]
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

// @Summary Get Store Order Detail
// @Description Get Store Order Detail
// @Tags Order
// @Accept json
// @Produce json
// @Param orderId path string true "Order ID"
// @Router /store/order/{orderId} [get]
func (o orderApiHandler) GetStoreOrderDetail(ctx *fiber.Ctx) error {
	context := ctx.Context()

	var req store.GetOrderOfStoreByIDRequest

	if err := ctx.ParamsParser(&req); err != nil {
		return errors.ErrBadRequest.WithInternalError(err)
	}

	storeID := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.StoreID = storeID

	result, err := o.orderQueryServ.GetDetailOrderOfStore(context, &req)
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

// @Summary Update Order Status By Store
// @Description Update Order Status By Store
// @Tags Order
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header"
// @Param UpdateOrderStatusRequest body UpdateOrderStatusRequest true "Update Order Status Request"
// @Router /store/order/status [put]
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

// @Summary Update Order Status By Delivery
// @Description Update Order Status By Delivery
// @Tags Order
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header"
// @Param UpdateOrderStatusRequest body UpdateOrderStatusRequest true "Update Order Status Request"
// @Router /delivery/order/status [put]
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
		log.Error(err)
		return errors.ErrBadRequest
	}

	err := o.orderCommandServ.DeliveryUpdateOrderStatus(context, req)
	if err != nil {
		return err
	}

	data := responses.DefaultSuccess

	return data.JSON(ctx)
}

// @Summary Get Orders By Delivery
// @Description Get Orders By Delivery
// @Tags Order
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Param order query string false "Order"
// @Router /delivery/order [get]
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

// @Summary Search Order ID By Keyword
// @Description Search Order ID By Keyword
// @Tags Order
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Param order query string false "Order"
// @Param keyword query string true "Keyword"
// @Router /order/store/search [get]
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
		log.Error(err)
		return errors.ErrBadRequest
	}

	storeID := fmt.Sprintf("%v", ctx.Locals(auth.STORE_ID))
	req.StoreID = storeID

	result, err := o.orderQueryServ.SearchStoreOrderID(context, &req)
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

// @Summary Admin Counting Order
// @Description Admin Counting Order
// @Tags Order
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Param order query string false "Order"
// @Router /admin/order/count [get]
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

// @Summary User Counting Order
// @Description User Counting Order
// @Tags Order
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Param order query string false "Order"
// @Router /user/order/count [get]
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

// @Summary Store Counting Order
// @Description Store Counting Order
// @Tags Order
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Param order query string false "Order"
// @Router /store/order/count [get]
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

// @Summary Delivery Counting Order
// @Description Delivery Counting Order
// @Tags Order
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param sort query string false "Sort"
// @Param order query string false "Order"
// @Router /delivery/order/count [get]
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
