package repo

import (
	"context"
	"strings"
	"tyr/internal/types"

	repoutil "github.com/M15t/gram/pkg/util/repo"
	requestutil "github.com/M15t/gram/pkg/util/request"
	"gorm.io/gorm"
)

// Document represents the client for document table
type Document struct {
	*repoutil.Repo[types.Document]
}

// NewDocument returns a new document database instance
func NewDocument(gdb *gorm.DB) *Document {
	return &Document{repoutil.NewRepo[types.Document](gdb)}
}

// FindByAPIMRequestID finds a document by the given apimrequestID
func (r *Document) FindByAPIMRequestID(ctx context.Context, apimReqID string) (*types.Document, error) {
	rec := &types.Document{}
	if err := r.GDB.WithContext(ctx).Preload("DocumentItem").Where(`apim_request_id = ?`, apimReqID).Take(rec).Error; err != nil {
		return nil, err
	}

	return rec, nil
}

// ReadByID read a document by given id
func (r *Document) ReadByID(ctx context.Context, documentID string) (*types.Document, error) {
	rec := &types.Document{}
	if err := r.GDB.WithContext(ctx).Preload("DocumentItem").Where(`id = ?`, documentID).Take(rec).Error; err != nil {
		return nil, err
	}

	return rec, nil
}

// List reads all documents by given conditions
func (r *Document) List(ctx context.Context, output interface{}, count *int64, lc *requestutil.ListCondition[DocumentsFilter], preloadConds []string) error {
	conds := []string{}
	vars := []any{}
	if lc.Filter.Search != "" {
		conds = append(conds, "(vendor_name like ? OR customer_name like ?)")
		sVal := strings.ReplaceAll(lc.Filter.Search, "%", "")
		sVal = strings.ReplaceAll(sVal, "?", "")
		sVal += "%"
		vars = append(vars, sVal, sVal, sVal, sVal)
	}

	if lc.Filter.UserID != "" {
		conds = append(conds, "user_id = ?")
		vars = append(vars, lc.Filter.UserID)
	}

	return r.ReadAllByCondition(ctx, output, count, &requestutil.ListQueryCondition{
		Page:    lc.Page,
		PerPage: lc.PerPage,
		Sort:    lc.Sort,
		Count:   lc.Count,
		Filter:  append([]any{strings.Join(conds, " AND ")}, vars...),
	}, preloadConds...)

}
