package types

import (
	"gorm.io/datatypes"
)

// ActivityLog represents the activity log model
// swagger:model
type ActivityLog struct {
	Base
	RequestURL      string         `json:"request_url"`
	RequestMethod   string         `json:"request_method"`
	RequestHeaders  datatypes.JSON `json:"request_headers"`
	RequestBody     datatypes.JSON `json:"request_body"`
	ResponseCode    int            `json:"response_code"`
	ResponseHeaders datatypes.JSON `json:"response_headers"`
	ResponseBody    datatypes.JSON `json:"response_body"`
	DurationMS      int64          `json:"duration_ms"`
	IPAddress       string         `json:"ip_address"`
	APIMRequestID   string         `json:"apim_request_id" gorm:"column:apim_request_id;type:varchar(36)"`
}
