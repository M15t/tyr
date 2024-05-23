package user

import (
	"tyr/internal/rbac"
	"tyr/internal/types"

	"github.com/M15t/gram/pkg/server"

	contextutil "tyr/internal/api/context"

	structutil "github.com/M15t/gram/pkg/util/struct"
)

// Create creates new user
func (s *User) Create(c contextutil.Context, data CreateUserReq) (*types.User, error) {
	if err := s.enforce(c, rbac.ActionCreateAll); err != nil {
		return nil, err
	}

	if existed, err := s.repo.User.Existed(c.GetContext(), map[string]interface{}{"email": data.Email}); err != nil || existed {
		return nil, ErrEmailExisted.SetInternal(err)
	}

	rec := &types.User{
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Phone:     data.Phone,
		Password:  s.cr.HashPassword(data.Password),
		Role:      data.Role,
	}

	if err := s.repo.User.Create(c.GetContext(), rec); err != nil {
		return nil, server.NewHTTPInternalError("error creating user").SetInternal(err)
	}

	return rec, nil
}

// Read returns single user by id
func (s *User) Read(c contextutil.Context, id string) (*types.User, error) {
	if err := s.enforce(c, rbac.ActionReadAll); err != nil {
		return nil, err
	}

	rec := &types.User{}
	if err := s.repo.User.ReadByID(c.GetContext(), rec, id); err != nil {
		return nil, server.NewHTTPInternalError("error reading user").SetInternal(err)
	}

	return rec, nil
}

// List returns the list of users
func (s *User) List(c contextutil.Context, req ListUserReq) (*ListUsersResp, error) {
	if err := s.enforce(c, rbac.ActionReadAll); err != nil {
		return nil, err
	}

	var count int64 = 0
	data := []*types.User{}
	if err := s.repo.User.List(c.GetContext(), &data, &count, req.ToListCond()); err != nil {
		return nil, server.NewHTTPInternalError("Error listing user").SetInternal(err)
	}

	return &ListUsersResp{
		Data:       data,
		TotalCount: count,
	}, nil
}

// Update updates user information
func (s *User) Update(c contextutil.Context, id string, data UpdateUserReq) (*types.User, error) {
	if err := s.enforce(c, rbac.ActionUpdateAll); err != nil {
		return nil, err
	}

	if err := s.repo.User.Update(c.GetContext(), structutil.ToMap(data), id); err != nil {
		return nil, server.NewHTTPInternalError("error reading user").SetInternal(err)
	}

	return s.Read(c, id)
}

// Delete deletes user by id
func (s *User) Delete(c contextutil.Context, id string) error {
	if err := s.enforce(c, rbac.ActionDeleteAll); err != nil {
		return err
	}

	if existed, err := s.repo.User.Existed(c.GetContext(), id); err != nil || !existed {
		return ErrUserNotFound.SetInternal(err)
	}

	return s.repo.User.Delete(c.GetContext(), id)
}

// Me returns current authenticated user
func (s *User) Me(c contextutil.Context) (*types.User, error) {
	rec := &types.User{}
	if err := s.repo.User.ReadByID(c.GetContext(), rec, c.AuthUser().ID); err != nil {
		return nil, server.NewHTTPInternalError("error reading user").SetInternal(err)
	}

	return rec, nil
}

// ChangePassword changes user password
func (s *User) ChangePassword(c contextutil.Context, data ChangePasswordReq) error {
	rec, err := s.Me(c)
	if err != nil {
		return err
	}

	if !s.cr.CompareHashAndPassword(rec.Password, data.OldPassword) {
		return ErrIncorrectPassword
	}

	return s.repo.User.Update(c.GetContext(), &types.User{
		Password: s.cr.HashPassword(data.NewPassword),
	}, rec.ID)
}

// enforce checks user permission to perform the action
func (s *User) enforce(c contextutil.Context, action string) error {
	au := c.AuthUser()
	if au == nil || !s.rbac.Enforce(au.Role, rbac.ObjectUser, action) {
		return rbac.ErrForbiddenAction
	}
	return nil
}
