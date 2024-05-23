package session

import (
	"tyr/internal/rbac"
	"tyr/internal/types"

	"github.com/M15t/gram/pkg/server"

	contextutil "tyr/internal/api/context"

	structutil "github.com/M15t/gram/pkg/util/struct"
)

// Read returns single session by id
func (s *Session) Read(c contextutil.Context, id string) (*types.Session, error) {
	if err := s.enforce(c, rbac.ActionReadAll); err != nil {
		return nil, err
	}

	rec := &types.Session{}
	if err := s.repo.Session.ReadByID(c.GetContext(), rec, id); err != nil {
		return nil, server.NewHTTPInternalError("error reading session").SetInternal(err)
	}

	return rec, nil
}

// List returns the list of users
func (s *Session) List(c contextutil.Context, req ListSessionReq) (*ListSessionsResp, error) {
	if err := s.enforce(c, rbac.ActionReadAll); err != nil {
		return nil, err
	}

	filter := map[string]any{}
	lqc := req.ToListQueryCond([]any{filter})

	var count int64 = 0
	data := []*types.Session{}
	if err := s.repo.Session.ReadAllByCondition(c.GetContext(), &data, &count, lqc); err != nil {
		return nil, server.NewHTTPInternalError("Error listing session").SetInternal(err)
	}

	return &ListSessionsResp{
		Data:       data,
		TotalCount: count,
	}, nil
}

// Update updates session information
func (s *Session) Update(c contextutil.Context, id string, data UpdateSessionReq) (*types.Session, error) {
	if err := s.enforce(c, rbac.ActionUpdateAll); err != nil {
		return nil, err
	}

	if err := s.repo.Session.Update(c.GetContext(), structutil.ToMap(data), id); err != nil {
		return nil, server.NewHTTPInternalError("error reading session").SetInternal(err)
	}

	return s.Read(c, id)
}

// Delete deletes session by id
func (s *Session) Delete(c contextutil.Context, id string) error {
	if err := s.enforce(c, rbac.ActionDeleteAll); err != nil {
		return err
	}

	if existed, err := s.repo.Session.Existed(c.GetContext(), id); err != nil || !existed {
		return ErrSessionNotFound.SetInternal(err)
	}

	return s.repo.Session.Delete(c.GetContext(), id)
}

// enforce checks user permission to perform the action
func (s *Session) enforce(c contextutil.Context, action string) error {
	au := c.AuthUser()
	if au == nil || !s.rbac.Enforce(au.Role, rbac.ObjectSession, action) {
		return rbac.ErrForbiddenAction
	}
	return nil
}
