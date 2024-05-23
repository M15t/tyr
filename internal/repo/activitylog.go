package repo

import (
	"context"
	"tyr/internal/types"

	repoutil "github.com/M15t/gram/pkg/util/repo"

	"gorm.io/gorm"
)

// ActivityLog represents the client for activity log table
type ActivityLog struct {
	*repoutil.Repo[types.ActivityLog]
}

// NewActivityLog returns a new activity log database instance
func NewActivityLog(gdb *gorm.DB) *ActivityLog {
	return &ActivityLog{repoutil.NewRepo[types.ActivityLog](gdb)}
}

// FindByAPIMRequestID finds a record by apim_request_id
func (r *ActivityLog) FindByAPIMRequestID(ctx context.Context, APIMRequestID string) (*types.ActivityLog, error) {
	rec := &types.ActivityLog{}
	if err := r.GDB.WithContext(ctx).Where(`apim_request_id = ?`, APIMRequestID).Take(rec).Error; err != nil {
		return nil, err
	}

	return rec, nil
}
