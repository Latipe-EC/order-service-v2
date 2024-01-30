package ordercommand

import (
	"github.com/gofiber/fiber/v2/log"
	"latipe-order-service-v2/internal/domain/entities/order"
	"latipe-order-service-v2/internal/domain/msgDTO"
	"latipe-order-service-v2/pkg/util/mapper"
)

func MappingDataToMessage(dao *order.Order, cartIds []string) *msgDTO.OrderMessage {
	orderMsg := msgDTO.OrderMessage{}

	if err := mapper.BindingStruct(dao, &orderMsg); err != nil {
		log.Error(err)
		return nil
	}

	orderMsg.Address.AddressId = dao.Delivery.AddressId
	orderMsg.Delivery.DeliveryId = dao.Delivery.DeliveryId
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
