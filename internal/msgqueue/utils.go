package msgqueue

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2/log"
	"latipe-order-service-v2/internal/domain/entities/order"
	"latipe-order-service-v2/internal/domain/msgDTO"
	"latipe-order-service-v2/pkg/util/mapper"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func ParseOrderToByte(request interface{}) ([]byte, error) {
	jsonObj, err := json.Marshal(&request)
	if err != nil {
		return nil, err
	}
	return jsonObj, err
}

func MappingDataToMessage(dao *order.Order) *msgDTO.OrderMessage {
	orderMsg := msgDTO.OrderMessage{}

	if err := mapper.BindingStruct(dao, &orderMsg); err != nil {
		log.Error(err)
		return nil
	}

	orderMsg.Address.AddressId = dao.Delivery.AddressId
	orderMsg.Delivery.DeliveryId = dao.Delivery.DeliveryId

	//order detail
	var orderItems []msgDTO.OrderItemsMessage
	for _, i := range dao.OrderItem {
		item := msgDTO.OrderItemsMessage{
			ProductItem: msgDTO.ProductItem{
				ProductID:   i.ProductID,
				ProductName: i.ProductName,
				StoreID:     i.StoreID,
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

	return &orderMsg
}
