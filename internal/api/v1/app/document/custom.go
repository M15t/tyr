package document

import (
	"net/http"

	"github.com/M15t/gram/pkg/server"
)

// Custom errors
var (
	ErrDocumentIsEmpty      = server.NewHTTPError(http.StatusBadRequest, "DOCUMENT_EMPTY", "Azure returns empty document")
	ErrDocumentNotFound     = server.NewHTTPError(http.StatusBadRequest, "DOCUMENT_NOTFOUND", "Document not found")
	ErrCreateTransferIntent = server.NewHTTPError(http.StatusBadRequest, "PLAID_CREATE_TRANSFER_INTENT_FAILED", "Create transfer intent failed")
)
