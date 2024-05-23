package document

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	contextutil "tyr/internal/api/context"
	"tyr/internal/rbac"
	"tyr/internal/types"
	"tyr/third_party/azure"

	"github.com/M15t/gram/pkg/server"
	structutil "github.com/M15t/gram/pkg/util/struct"
	"github.com/iancoleman/strcase"
	"gorm.io/datatypes"
)

// Analyze sends a document to Azure for analysis.
// It encodes the file content as base64, creates a JSON payload, and sends it to Azure for analysis.
// It then creates a new document entry in the repository with the analysis results.
// Returns the APIM request ID of the analysis.
func (s *Document) Analyze(c contextutil.Context, req AnalyzeDocumentReq) (*AnalyzeDocumentRes, error) {
	if err := s.enforce(c, rbac.ActionCreate); err != nil {
		return nil, err
	}

	// ! move it to configuration on PROD
	apiVersion := "2023-07-31"
	modelID := "prebuilt-receipt"

	// Open file from multipart.FileHeader
	file, err := req.Document.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Encode the file content as base64
	base64Str := base64.StdEncoding.EncodeToString(fileContent)

	reqPayload := map[string]interface{}{
		"base64Source": base64Str,
		// "urlSource":    "https://m15t-public-bucket.s3.amazonaws.com/11832804300.pdf", // * support url btw
	}

	// Encode the map as JSON
	jsonData, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("error encoding JSON: %s", err)
	}

	// Create a bytes.Reader from the JSON-encoded byte slice
	payload := bytes.NewReader(jsonData)

	resHeaders, err := s.azure.AnalyzeDocument(c, modelID, apiVersion, payload)
	if err != nil {
		return nil, err
	}

	newDocument := types.Document{
		FileName:          req.Document.Filename,
		FilePath:          req.Document.Header.Get("Content-Disposition"),
		OriginalFileName:  req.Document.Filename,
		APIMRequestID:     resHeaders.APIMRequestID[0],
		OperationLocation: resHeaders.OperationLocation[0],
		ModelID:           modelID,
		APIVersion:        apiVersion,
		DocumentItem: &types.DocumentItem{
			Data: datatypes.JSON([]byte{}),
		},
	}

	if err := s.repo.Document.Create(c.GetContext(), &newDocument); err != nil {
		return nil, err
	}

	return &AnalyzeDocumentRes{
		APIMRequestID: resHeaders.APIMRequestID[0],
	}, nil
}

// Get retrieves the document information by the given APIM request ID.
// It fetches the document from the repository based on the APIM request ID.
// If the document is not found in the activity logs, it requests Azure to get the result document.
// It then updates the document details including vendor information, invoice details, payment information, and items.
// Finally, it updates the document item and returns the updated document.
func (s *Document) Get(c contextutil.Context, apimReqID string) (*types.Document, error) {
	// get document by apimReqID
	document, err := s.repo.Document.FindByAPIMRequestID(c.GetContext(), apimReqID)
	if err != nil {
		return nil, err
	}

	// check in activity logs first
	var resRawDocument *azure.ResultAnalyzeResponse
	activityLog, err := s.repo.ActivityLog.FindByAPIMRequestID(c.GetContext(), apimReqID)
	if err == nil && activityLog != nil {
		responseBodyBytes, err := json.Marshal(activityLog.ResponseBody)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(responseBodyBytes, &resRawDocument); err != nil {
			return nil, err
		}
	} else {
		// request Azure to get result document if not existed record in activity logs
		resRawDocument, err = s.azure.GetAnalyzeDocument(c, document.OperationLocation)
		if err != nil {
			return nil, err
		}
	}

	if len(resRawDocument.AnalyzeResult.Documents) == 0 {
		return nil, ErrDocumentIsEmpty
	}

	// update document
	totalPage := len(resRawDocument.AnalyzeResult.Pages)
	page := 0 // in case proccess multiple pages

	merchantName := resRawDocument.AnalyzeResult.Documents[page].Fields.MerchantName.Content
	merchantAddress := resRawDocument.AnalyzeResult.Documents[page].Fields.MerchantAddress.Content
	merchantPhoneNumber := resRawDocument.AnalyzeResult.Documents[page].Fields.MerchantPhoneNumber.Content

	transactionDate := parseStringToDate(resRawDocument.AnalyzeResult.Documents[page].Fields.TransactionDate.Content)
	transactionTime := resRawDocument.AnalyzeResult.Documents[page].Fields.TransactionTime.Content

	currency := resRawDocument.AnalyzeResult.Documents[page].Fields.TaxDetails.ValueArray[0].ValueObject.Amount.ValueCurrency.CurrencyCode
	taxDetails := resRawDocument.AnalyzeResult.Documents[page].Fields.TaxDetails.ValueArray[0].Content
	totalTax := resRawDocument.AnalyzeResult.Documents[page].Fields.TotalTax.ValueNumber
	total := resRawDocument.AnalyzeResult.Documents[page].Fields.Total.ValueNumber

	items := resRawDocument.AnalyzeResult.Documents[page].Fields.Items

	// update document item
	mappedItems := make([]map[string]interface{}, 0)
	for _, item := range items.ValueArray {
		newItem := make(map[string]interface{})
		for fieldName, fieldValue := range item.ValueObject {
			if fieldValueMap, ok := fieldValue.(map[string]interface{}); ok {
				if valueString, ok := fieldValueMap["content"].(string); ok {
					newItem[strcase.ToSnake(fieldName)] = valueString
				}
			}
		}
		mappedItems = append(mappedItems, newItem)
	}

	jsonData, err := json.Marshal(mappedItems)
	if err != nil {
		return nil, err
	}

	if err := s.repo.DocumentItem.Update(c.GetContext(), &types.DocumentItem{
		Data: jsonData,
	}, "document_id", document.ID); err != nil {
		return nil, err
	}

	if err := s.repo.Document.Update(c.GetContext(), &types.Document{
		UserID:              c.AuthUser().ID,
		TotalPage:           totalPage,
		MerchantName:        merchantName,
		MerchantAddress:     merchantAddress,
		MerchantPhoneNumber: merchantPhoneNumber,
		Currency:            currency,
		TotalTax:            totalTax,
		Total:               total,
		TaxDetails:          taxDetails,
		TransactionDate:     transactionDate,
		TransactionTime:     transactionTime,
	}, "id = ?", document.ID); err != nil {
		return nil, err
	}

	return s.repo.Document.FindByAPIMRequestID(c.GetContext(), apimReqID)
}

// Read returns single user by id
func (s *Document) Read(c contextutil.Context, id string) (*types.Document, error) {
	if err := s.enforce(c, rbac.ActionRead); err != nil {
		return nil, err
	}

	if existed, err := s.repo.Document.Existed(c.GetContext(), `user_id`, c.AuthUser().ID); !existed && err == nil {
		return nil, ErrDocumentNotFound
	}

	return s.repo.Document.ReadByID(c.GetContext(), id)
}

// List returns the list of users
func (s *Document) List(c contextutil.Context, req ListDocumentReq) (*ListDocumentsResp, error) {
	if err := s.enforce(c, rbac.ActionRead); err != nil {
		return nil, err
	}

	var count int64 = 0
	data := []*types.Document{}
	lc := req.ToListCond()
	lc.Filter.UserID = c.AuthUser().ID
	preloadConds := []string{"DocumentItem"}
	if err := s.repo.Document.List(c.GetContext(), &data, &count, lc, preloadConds); err != nil {
		return nil, server.NewHTTPInternalError("Error listing user").SetInternal(err)
	}

	return &ListDocumentsResp{
		Data:       data,
		TotalCount: count,
	}, nil
}

// Update updates document information
func (s *Document) Update(c contextutil.Context, id string, data UpdateDocumentReq) (*types.Document, error) {
	if err := s.enforce(c, rbac.ActionUpdate); err != nil {
		return nil, err
	}

	if existed, err := s.repo.Document.Existed(c.GetContext(), id); err != nil || !existed {
		return nil, ErrDocumentNotFound.SetInternal(err)
	}

	if err := s.repo.Document.Update(c.GetContext(), structutil.ToMap(data), id); err != nil {
		return nil, server.NewHTTPInternalError("error updating document").SetInternal(err)
	}

	return s.Read(c, id)
}

// Delete deletes document by id
func (s *Document) Delete(c contextutil.Context, id string) error {
	if err := s.enforce(c, rbac.ActionDelete); err != nil {
		return err
	}

	if existed, err := s.repo.Document.Existed(c.GetContext(), id); err != nil || !existed {
		return ErrDocumentNotFound.SetInternal(err)
	}

	if err := s.repo.DocumentItem.Delete(c.GetContext(), `document_id`, id); err != nil {
		return ErrDocumentNotFound.SetInternal(err)
	}

	return s.repo.Document.Delete(c.GetContext(), id)
}

// enforce checks document permission to perform the action
func (s *Document) enforce(c contextutil.Context, action string) error {
	au := c.AuthUser()
	if au == nil || !s.rbac.Enforce(au.Role, rbac.ObjectDocument, action) {
		return rbac.ErrForbiddenAction
	}
	return nil
}
