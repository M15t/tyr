package repo

import (
	"tyr/internal/types"

	repoutil "github.com/M15t/gram/pkg/util/repo"
	"gorm.io/gorm"
)

// DocumentItem represents the client for document table
type DocumentItem struct {
	*repoutil.Repo[types.DocumentItem]
}

// NewDocumentItem returns a new document item database instance
func NewDocumentItem(gdb *gorm.DB) *DocumentItem {
	return &DocumentItem{repoutil.NewRepo[types.DocumentItem](gdb)}
}
