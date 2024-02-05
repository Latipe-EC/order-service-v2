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

func MappingToProductRequest(dto reqDTO.StoreOrder) (*productgrpc.GetPurchaseProductRequest, map[string]int) {

	itemMap := make(map[string]int)
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
		if i.OptionId != "" {
			itemMap[i.OptionId] = i.Quantity
		} else {
			itemMap[i.ProductId] = i.Quantity
		}

		items = append(items, &item)
	}

	req.Items = items

	return &req, itemMap
}

func GetQuantityItems(productId string, optionId string, item map[string]int) int {
	if optionId != "" {
		return item[optionId]
	} else {
		return item[productId]
	}
}
func MappingDataToMessage(daos []*order.Order, cartMap map[string][]string) *msgDTO.CreateOrderMessage {
	orderMsg := msgDTO.CreateOrderMessage{}

	orderMsg.UserRequest.UserId = daos[0].UserId
	orderMsg.UserRequest.Username = daos[0].Username
	orderMsg.Address.AddressId = daos[0].Delivery.AddressId
	orderMsg.Address.AddressDetail = daos[0].Delivery.ShippingAddress
	orderMsg.Address.Name = daos[0].Delivery.ShippingName
	orderMsg.Address.Phone = daos[0].Delivery.ShippingPhone

	for _, i := range daos {
		orderDetail := msgDTO.OrderDetail{}
		if err := mapper.Copy(i, orderDetail); err != nil {
			log.Error(err)
			return nil
		}

		//assign delivery data
		orderDetail.Delivery.DeliveryId = i.Delivery.DeliveryId
		orderDetail.Delivery.Cost = i.Delivery.Cost
		orderDetail.Delivery.ReceivingDate = i.Delivery.ReceivingDate
		orderDetail.Delivery.Name = i.Delivery.DeliveryName

		//assign order store data
		orderDetail.StoreID = i.StoreId
		//order detail
		var orderItems []msgDTO.OrderItemsMessage
		for _, j := range i.OrderItem {
			item := msgDTO.OrderItemsMessage{
				ProductItem: msgDTO.ProductItem{
					ProductID:   j.ProductID,
					ProductName: j.ProductName,
					NameOption:  j.NameOption,
					OptionID:    j.OptionID,
					Quantity:    j.Quantity,
					Price:       j.Price,
					NetPrice:    j.NetPrice,
					Image:       j.ProdImg,
				},
			}
			orderItems = append(orderItems, item)
		}

		orderDetail.OrderItems = orderItems
		orderDetail.CartIds = cartMap[i.StoreId]
	}

	return &orderMsg
}

func GenKeyOrder(userId string) string {
	keyGen := strings.ReplaceAll(uuid.NewString(), "-", "")[:8]
	return fmt.Sprintf("orvn%v%v", userId[:4], keyGen)
}
