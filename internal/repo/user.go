package repo

import (
	"context"
	"strings"

	"tyr/internal/types"

	repoutil "github.com/M15t/gram/pkg/util/repo"
	requestutil "github.com/M15t/gram/pkg/util/request"

	"gorm.io/gorm"
)

// User represents the client for user table
type User struct {
	*repoutil.Repo[types.User]
}

// NewUser returns a new user database instance
func NewUser(gdb *gorm.DB) *User {
	return &User{repoutil.NewRepo[types.User](gdb)}
}

// List reads all users by given conditions
func (r *User) List(ctx context.Context, output interface{}, count *int64, lc *requestutil.ListCondition[UsersFilter]) error {
	conds := []string{}
	vars := []any{}
	if lc.Filter.Search != "" {
		conds = append(conds, "(first_name like ? OR last_name like ? OR email like ?)")
		sVal := strings.ReplaceAll(lc.Filter.Search, "%", "")
		sVal = strings.ReplaceAll(sVal, "?", "")
		sVal += "%"
		vars = append(vars, sVal, sVal, sVal)
	}

	return r.ReadAllByCondition(ctx, output, count, &requestutil.ListQueryCondition{
		Page:    lc.Page,
		PerPage: lc.PerPage,
		Sort:    lc.Sort,
		Count:   lc.Count,
		Filter:  append([]any{strings.Join(conds, " AND ")}, vars...),
	})

}

// FindByEmail finds a user by the given email
func (r *User) FindByEmail(ctx context.Context, email string) (*types.User, error) {
	rec := &types.User{}
	if err := r.GDB.WithContext(ctx).Where(`email = ?`, email).Take(rec).Error; err != nil {
		return nil, err
	}

	return rec, nil
}

// UpdateRefreshToken updates the refresh token of the given user
func (r *User) UpdateRefreshToken(ctx context.Context, userID, refreshToken string) error {
	return r.GDB.WithContext(ctx).Model(&types.User{}).Where(`id = ?`, userID).Update(`refresh_token`, refreshToken).Error
}

// Create is a method of the User struct that creates a new user in the database.
/**
 * Parameters:
 * - ctx: The context.Context object for the request.
 * - user: A pointer to the types.User object representing the user to be created.
 *
 * Returns:
 * - error: An error object if there was an error creating the user, otherwise nil.
 */
func (r *User) Create(ctx context.Context, user *types.User) error {
	txErr := r.GDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(user).Error
	})

	return txErr
}
