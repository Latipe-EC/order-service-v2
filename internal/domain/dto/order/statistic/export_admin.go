package statistic

import (
	"bytes"
)

type ExportOrderDataForAdminRequest struct {
	Date   string `query:"date" validate:"required"` //yyyy-mm
	UserID string
}

type ExportOrderDataForAdminResponse struct {
	QueryDate      string `json:"query_date"`
	FileAttachment *bytes.Reader
	FileName       string `json:"file_name"`
}
