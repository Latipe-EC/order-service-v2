package statistic

import "latipe-order-service-v2/internal/domain/dto/custom_entity"

type GetTotalStoreOrderInMonthRequest struct {
	Date    string `query:"date"`
	StoreId string
}

type GetTotalOrderInMonthResponse struct {
	FilterDate string                                  `json:"filter_date,omitempty"`
	Items      []custom_entity.TotalOrderInSystemInDay `json:"items"`
}

type GetTotalOrderInYearOfStoreRequest struct {
	Year    int `query:"year"`
	StoreID string
}

type GetTotalOrderInYearOfStoreResponse struct {
	FilterDate string                                    `json:"filter_date,omitempty"`
	Items      []custom_entity.TotalOrderInSystemInMonth `json:"items"`
}
