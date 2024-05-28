package document

import (
	"net/http"
	contextutil "tyr/internal/api/context"
	"tyr/internal/types"

	"github.com/labstack/echo/v4"

	httputil "github.com/M15t/gram/pkg/util/http"
)

// HTTP represents document http service
type HTTP struct {
	contextutil.Context
	svc Service
}

// Service represents document application interface
type Service interface {
	Analyze(contextutil.Context, AnalyzeDocumentReq) (*AnalyzeDocumentRes, error)
	Get(contextutil.Context, string) (*types.Document, error)

	Read(contextutil.Context, string) (*types.Document, error)
	List(contextutil.Context, ListDocumentReq) (*ListDocumentsResp, error)
	Update(contextutil.Context, string, UpdateDocumentReq) (*types.Document, error)
	Delete(contextutil.Context, string) error
}

// NewHTTP attaches handlers to Echo routers under given group
func NewHTTP(svc Service, eg *echo.Group) {
	h := HTTP{svc: svc}

	// swagger:operation POST /v1/app/documents/analyze/upload app-documents-analyze appDocumentAnalyzeUpload
	// ---
	// summary: Analyzes new document, upload and send file to Azure for processing
	// consumes:
	// - multipart/form-data
	// parameters:
	// - name: document
	//   in: formData
	//   type: file
	//   description: The document to upload
	// responses:
	//   "200":
	//     description: The request id of document
	//     schema:
	//       "$ref": "#/definitions/AnalyzeDocumentRes"
	//   default:
	//     description: 'Possible errors: 400, 401, 403, 500'
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	eg.POST("/analyze/upload", h.analyzeUpload)

	// swagger:operation GET /v1/app/documents/analyze/get/{id} app-documents-analyze appDocumentAnalyzeGet
	// ---
	// summary: Get document from Azure by apim_request_id and analyze the result
	// parameters:
	// - name: id
	//   in: path
	//   description: apim_request_id of document
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     description: The document
	//     schema:
	//       "$ref": "#/definitions/Document"
	//   default:
	//     description: 'Possible errors: 400, 401, 403, 500'
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	eg.GET("/analyze/get/:id", h.analyzeGet)

	// swagger:operation GET /v1/app/documents/{id} app-documents documentsRead
	// ---
	// summary: Returns a single document
	// parameters:
	// - name: id
	//   in: path
	//   description: id of document
	//   type: string
	//   required: true
	// responses:
	//   "200":
	//     description: The document
	//     schema:
	//       "$ref": "#/definitions/Document"
	//   default:
	//     description: 'Possible errors: 400, 401, 403, 404, 500'
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	eg.GET("/:id", h.read)

	// swagger:operation GET /v1/app/documents app-documents documentsList
	// ---
	// summary: Returns list of documents
	// responses:
	//   "200":
	//     description: List of documents
	//     schema:
	//       "$ref": "#/definitions/ListDocumentsResp"
	//   default:
	//     description: 'Possible errors: 400, 401, 403, 500'
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	eg.GET("", h.list)

	// swagger:operation PATCH /v1/app/documents/{id} app-documents documentsUpdate
	// ---
	// summary: Updates document information
	// parameters:
	// - name: id
	//   in: path
	//   description: id of document
	//   type: string
	//   required: true
	// - name: request
	//   in: body
	//   description: Request body
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/UpdateDocumentReq"
	// responses:
	//   "200":
	//     description: The updated document
	//     schema:
	//       "$ref": "#/definitions/Document"
	//   default:
	//     description: 'Possible errors: 400, 401, 403, 404, 500'
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	eg.PATCH("/:id", h.update)

	// swagger:operation DELETE /v1/app/documents/{id} app-documents documentsDelete
	// ---
	// summary: Deletes an document
	// parameters:
	// - name: id
	//   in: path
	//   description: id of document
	//   type: string
	//   required: true
	// responses:
	//   "204":
	//     "$ref": "#/responses/ok"
	//   default:
	//     description: 'Possible errors: 401, 403, 404, 500'
	//     schema:
	//       "$ref": "#/definitions/ErrorResponse"
	eg.DELETE("/:id", h.delete)
}

func (h *HTTP) analyzeUpload(c echo.Context) error {
	r := AnalyzeDocumentReq{}
	if err := c.Bind(&r); err != nil {
		return err
	}
	// Multipart form
	// Aug 2023, until now there is no way to use c.Bind() with multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	documents := form.File["document"]

	if len(documents) > 0 {
		r.Document = documents[0]
	}

	resp, err := h.svc.Analyze(contextutil.NewContext(c), r)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) analyzeGet(c echo.Context) error {
	id := c.Param("id")

	resp, err := h.svc.Get(contextutil.NewContext(c), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) read(c echo.Context) error {
	id, err := httputil.ReqID(c)
	if err != nil {
		return err
	}
	resp, err := h.svc.Read(contextutil.NewContext(c), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) list(c echo.Context) error {
	req := ListDocumentReq{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	resp, err := h.svc.List(contextutil.NewContext(c), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) update(c echo.Context) error {
	id, err := httputil.ReqID(c)
	if err != nil {
		return err
	}
	r := UpdateDocumentReq{}
	if err := c.Bind(&r); err != nil {
		return err
	}

	resp, err := h.svc.Update(contextutil.NewContext(c), id, r)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) delete(c echo.Context) error {
	id, err := httputil.ReqID(c)
	if err != nil {
		return err
	}
	if err := h.svc.Delete(contextutil.NewContext(c), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
