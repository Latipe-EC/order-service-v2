package dto

const orderProductUrl = "/api/v1/products/check-in-stock"

type OrderProductRequest struct {
	StoreOrders []StoreOrderRequest
}
type StoreOrderRequest struct {
	StoreID string `json:"store_id"`
	Items   []ValidateItems
}

type ValidateItems struct {
	ProductId string `json:"productId"`
	OptionId  string `json:"optionId"`
	Quantity  int    `json:"quantity"`
}

type OrderProductResponse struct {
	Items []ShopOrders `json:"items"`
}

type ShopOrders struct {
	StoreID      string    `json:"store_id"`
	ProvinceCode string    `json:"province_code"`
	Items        []Product `json:"items"`
	TotalPrice   int       `json:"totalPrice"`
}

type Product struct {
	ProductId        string  `json:"productId"`
	Name             string  `json:"name"`
	Quantity         int     `json:"quantity"`
	Image            string  `json:"image"`
	Price            float64 `json:"price"`
	PromotionalPrice float64 `json:"promotionalPrice"`
	OptionId         string  `json:"optionId"`
	NameOption       string  `json:"nameOption"`
	StoreId          string  `json:"storeId"`
	TotalPrice       float64 `json:"totalPrice"`
}

func (OrderProductRequest) URL() string {
	return orderProductUrl
}
