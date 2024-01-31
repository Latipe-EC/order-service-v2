package ordercommand

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	reqDTO "latipe-order-service-v2/internal/domain/dto/order"
	"latipe-order-service-v2/internal/domain/entities/order"
	"latipe-order-service-v2/internal/domain/msgDTO"
	productgrpc "latipe-order-service-v2/internal/infrastructure/grpc/productServ"
	"latipe-order-service-v2/pkg/util/mapper"
	"strings"
)

func MappingToProductRequest(dto reqDTO.StoreOrder) *productgrpc.GetPurchaseProductRequest {
	req := productgrpc.GetPurchaseProductRequest{
		StoreId: dto.StoreID,
	}
	var items []*productgrpc.GetPurchaseItemRequest

	for _, i := range dto.Items {
		item := productgrpc.GetPurchaseItemRequest{
			ProductId: i.ProductId,
			OptionId:  i.OptionId,
			Quantity:  int32(i.Quantity),
		}

		items = append(items, &item)
	}

	req.Items = items

	return &req
}

func MappingDataToMessage(dao *order.Order, cartIds []string) *msgDTO.OrderMessage {
	orderMsg := msgDTO.OrderMessage{}
	dao.OrderStatusLog = nil
	if err := mapper.Copy(&orderMsg, dao); err != nil {
		log.Error(err)
		return nil
	}
	//assign user request
	orderMsg.UserRequest.UserId = dao.UserId
	orderMsg.UserRequest.Username = dao.Username

	//assign address data
	orderMsg.Address.AddressId = dao.Delivery.AddressId
	orderMsg.Address.AddressDetail = dao.Delivery.ShippingAddress
	orderMsg.Address.Name = dao.Delivery.ShippingName
	orderMsg.Address.Phone = dao.Delivery.ShippingPhone

	//assign delivery data
	orderMsg.Delivery.DeliveryId = dao.Delivery.DeliveryId
	orderMsg.Delivery.Cost = dao.Delivery.Cost
	orderMsg.Delivery.ReceivingDate = dao.Delivery.ReceivingDate
	orderMsg.Delivery.Name = dao.Delivery.DeliveryName

	//assign order store data
	orderMsg.StoreID = dao.StoreId
	//order detail
	var orderItems []msgDTO.OrderItemsMessage
	for _, i := range dao.OrderItem {
		item := msgDTO.OrderItemsMessage{
			ProductItem: msgDTO.ProductItem{
				ProductID:   i.ProductID,
				ProductName: i.ProductName,
				NameOption:  i.NameOption,
				OptionID:    i.OptionID,
				Quantity:    i.Quantity,
				Price:       i.Price,
				NetPrice:    i.NetPrice,
				Image:       i.ProdImg,
			},
		}
		orderItems = append(orderItems, item)
	}
	orderMsg.OrderItems = orderItems
	orderMsg.CartIds = cartIds

	return &orderMsg
}

func GenKeyOrder(userId string) string {
	keyGen := strings.ReplaceAll(uuid.NewString(), "-", "")[:8]
	return fmt.Sprintf("orvn%v%v", userId[:4], keyGen)
}
