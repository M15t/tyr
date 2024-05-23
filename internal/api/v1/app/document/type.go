package document

import (
	"mime/multipart"
	"tyr/internal/repo"
	"tyr/internal/types"

	requestutil "github.com/M15t/gram/pkg/util/request"
)

// AnalyzeDocumentReq struct
// swagger:model
type AnalyzeDocumentReq struct {
	Document *multipart.FileHeader `form:"document"`
}

// AnalyzeDocumentRes struct
// swagger:model
type AnalyzeDocumentRes struct {
	APIMRequestID string `json:"apim_request_id"`
}

// UpdateDocumentReq contains request data to update existing document
// swagger:model
type UpdateDocumentReq struct {
	// Vendor
	VendorName             *string `json:"vendor_name,omitempty"`
	VendorAddressRecipient *string `json:"vendor_address_recipient,omitempty"`
	VendorAddress          *string `json:"vendor_address,omitempty"`

	// Customer
	CustomerID               *string `json:"customer_id,omitempty"`
	CustomerName             *string `json:"customer_name,omitempty"`
	CustomerAddressRecipient *string `json:"customer_address_recipient,omitempty"`
	CustomerAddress          *string `json:"customer_address,omitempty"`

	// Billing
	BillingAddressRecipient *string `json:"billing_address_recipient,omitempty"`
	BillingAddress          *string `json:"billing_address,omitempty"`

	// Shipping
	ShippingAddressRecipient *string `json:"shipping_address_recipient,omitempty"`
	ShippingAddress          *string `json:"shipping_address,omitempty"`

	// Service
	ServiceAddressRecipient *string `json:"service_address_recipient,omitempty"`
	ServiceAddress          *string `json:"service_address,omitempty"`

	// Remit
	RemittanceAddressRecipient *string `json:"remittance_address_recipient,omitempty"`
	RemittanceAddress          *string `json:"remittance_address,omitempty"`

	PurchaseOrder *string `json:"purchase_order,omitempty"`

	// Invoice
	InvoiceID   *string `json:"invoice_id,omitempty"`
	InvoiceDate *string `json:"invoice_date,omitempty"`
	DueDate     *string `json:"due_date,omitempty"`
	PaymentTerm *string `json:"payment_term,omitempty"`

	Currency *string `json:"currency,omitempty"`

	// $$$
	SubTotal              **float64 `json:"sub_total,omitempty"`
	InvoiceTotal          *float64  `json:"invoice_total,omitempty"`
	TotalTax              *float64  `json:"total_tax,omitempty"`
	PreviousUnpaidBalance *float64  `json:"previous_unpaid_balance,omitempty"`
	AmountDue             *float64  `json:"amount_due,omitempty"`
}

// ListDocumentReq contains request data to get list of documents
// swagger:parameters documentsList
type ListDocumentReq struct {
	requestutil.ListQueryRequest
	// Search for document(s) by?
	// TODO: TBD
	Search string `json:"search,omitempty" query:"search"`
}

// ToListCond transforms the service request to repo conditions
func (lq *ListDocumentReq) ToListCond() *requestutil.ListCondition[repo.DocumentsFilter] {
	return &requestutil.ListCondition[repo.DocumentsFilter]{
		Page:    lq.Page,
		PerPage: lq.PerPage,
		Sort:    lq.Sort,
		Count:   true,
		Filter: repo.DocumentsFilter{
			Search: lq.Search,
		},
	}
}

// ListDocumentsResp contains list of paginated documents and total numbers after filtered
// swagger:model
type ListDocumentsResp struct {
	Data       []*types.Document `json:"data"`
	TotalCount int64             `json:"total_count"`
}
