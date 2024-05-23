package types

import (
	"time"

	"github.com/M15t/gram/pkg/util/ulidutil"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Document represents for document model
// swagger:model
type Document struct {
	Base
	UserID            string `json:"user_id"`
	OriginalFileName  string `json:"-" gorm:"type:varchar(255)"`
	FileName          string `json:"file_name"`
	FilePath          string `json:"file_path"`
	APIMRequestID     string `json:"apim_request_id" gorm:"column:apim_request_id;type:varchar(36)"`
	OperationLocation string `json:"-"`
	ModelID           string `json:"-" gorm:"type:varchar(20)"`
	APIVersion        string `json:"-" gorm:"type:varchar(20)"`

	// Merchant
	MerchantName        string `json:"merchant_name"`
	MerchantAddress     string `json:"merchant_address"`
	MerchantPhoneNumber string `json:"merchant_phone_number" gorm:"type:varchar(20)"`

	// Receipt
	TransactionDate string `json:"transaction_date" gorm:"type:varchar(50)"`
	TransactionTime string `json:"transaction_time" gorm:"type:varchar(20)"`

	Currency string `json:"currency" gorm:"type:varchar(3)"`

	// $$$
	SubTotal   float64 `json:"sub_total"`
	Total      float64 `json:"total"`
	TotalTax   float64 `json:"total_tax"`
	TaxDetails string  `json:"tax_details"`

	TotalPage int `json:"total_page"`

	DocumentItem *DocumentItem `json:"document_item,omitempty"`
}

// DocumentItem for document item (list of items) model
// swagger:model
type DocumentItem struct {
	// ID of the record
	ID string `json:"-" gorm:"primaryKey"`
	// The time that record is created
	CreatedAt time.Time `json:"-"`
	// The latest time that record is updated
	UpdatedAt time.Time `json:"-"`
	// Enable soft-delete
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	DocumentID string         `json:"-"`
	Data       datatypes.JSON `json:"data"`
}

// BeforeCreate hook executed by gorm
func (di *DocumentItem) BeforeCreate(tx *gorm.DB) (err error) {
	if di.ID == "" {
		di.ID = ulidutil.NewString()
	}
	return
}
