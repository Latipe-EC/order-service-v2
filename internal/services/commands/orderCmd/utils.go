package orderCmd

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
func MappingDataToMessage(daos []*order.Order, cartMap map[string][]string,
	checkout msgDTO.CheckoutMessage) *msgDTO.CreateOrderMessage {

	orderMsg := &msgDTO.CreateOrderMessage{}

	// Copy checkout message
	if err := mapper.Copy(&orderMsg.CheckoutMessage, checkout); err != nil {
		log.Error(err)
		return nil
	}

	// Extract user and address details from first DAO
	firstDao := daos[0]
	orderMsg.UserRequest = msgDTO.UserRequest{
		UserId:   firstDao.UserId,
		Username: firstDao.Username,
	}
	orderMsg.Address = msgDTO.OrderAddress{
		AddressId:     firstDao.Delivery.AddressId,
		AddressDetail: firstDao.Delivery.ShippingAddress,
		Name:          firstDao.Delivery.ShippingName,
		Phone:         firstDao.Delivery.ShippingPhone,
	}

	// Map order items
	for _, dao := range daos {
		orderDetail := msgDTO.OrderDetail{}

		// Copy data from DAO to order detail
		if err := mapper.Copy(&orderDetail, dao); err != nil {
			log.Error(err)
			return nil
		}

		// Assign delivery data
		orderDetail.Delivery = msgDTO.Delivery{
			DeliveryId:    dao.Delivery.DeliveryId,
			Cost:          dao.Delivery.Cost,
			ReceivingDate: dao.Delivery.ReceivingDate,
			Name:          dao.Delivery.DeliveryName,
		}

		// Map order store data
		orderDetail.Vouchers = dao.Vouchers
		orderDetail.StoreID = dao.StoreId

		// Map order items
		orderItems := make([]msgDTO.OrderItemsMessage, len(dao.OrderItem))
		for idx, item := range dao.OrderItem {
			orderItems[idx] = msgDTO.OrderItemsMessage{
				ProductItem: msgDTO.ProductItem{
					ProductID:   item.ProductID,
					ProductName: item.ProductName,
					NameOption:  item.NameOption,
					OptionID:    item.OptionID,
					Quantity:    item.Quantity,
					Price:       item.Price,
					NetPrice:    item.NetPrice,
					Image:       item.ProdImg,
				},
			}
		}
		orderDetail.OrderItems = orderItems
		orderDetail.CartIds = cartMap[dao.StoreId]

		orderMsg.OrderDetail = append(orderMsg.OrderDetail, orderDetail)
	}

	return orderMsg
}

func GenKeyOrder(userId string) string {
	keyGen := strings.ReplaceAll(uuid.NewString(), "-", "")[:8]
	return fmt.Sprintf("orvn%v%v", userId[:4], keyGen)
}

func isMatchPaymentMethod(required int, act int) bool {
	return required == 0 || required == act
}
