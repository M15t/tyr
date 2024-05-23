package session

import (
	"tyr/internal/types"

	requestutil "github.com/M15t/gram/pkg/util/request"
)

// UpdateSessionReq contains request data to update existing session
// swagger:model
type UpdateSessionReq struct {
	IsBlocked *bool `json:"is_blocked,omitempty"`
}

// ListSessionReq contains request data to get list of sessions
// swagger:parameters sessionsList
type ListSessionReq struct {
	requestutil.ListQueryRequest
}

// ListSessionsResp contains list of paginated users and total numbers after filtered
// swagger:model
type ListSessionsResp struct {
	Data       []*types.Session `json:"data"`
	TotalCount int64            `json:"total_count"`
}
