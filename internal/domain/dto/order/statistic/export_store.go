package statistic

import "bytes"

type ExportOrderDataForStoreRequest struct {
	Date     string `query:"date" validate:"required"` //yyyy-mm
	Username string
	StoreID  string
}

type ExportOrderDataForStoreResponse struct {
	QueryDate      string `json:"query_date"`
	FileAttachment *bytes.Reader
	FileName       string `json:"file_name"`
}
