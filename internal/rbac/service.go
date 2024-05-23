package rbac

import (
	"github.com/M15t/gram/pkg/rbac"
)

// New returns new RBAC service
func New(enableLog bool) *rbac.RBAC {
	r := rbac.NewWithConfig(rbac.Config{EnableLog: enableLog})

	// Add permission for user role
	r.AddPolicy(RoleUser, ObjectUser, ActionReadAll)

	r.AddPolicy(RoleUser, ObjectDocument, ActionCreate)
	r.AddPolicy(RoleUser, ObjectDocument, ActionRead)
	r.AddPolicy(RoleUser, ObjectDocument, ActionUpdate)
	r.AddPolicy(RoleUser, ObjectDocument, ActionDelete)

	r.AddPolicy(RoleUser, ObjectPlaid, ActionCreate)

	// Add permission for admin role
	r.AddPolicy(RoleAdmin, ObjectUser, ActionAny)
	r.AddPolicy(RoleAdmin, ObjectSession, ActionAny)
	r.AddPolicy(RoleAdmin, ObjectDocument, ActionAny)

	// Add permission for superadmin role
	r.AddPolicy(RoleSuperAdmin, ObjectAny, ActionAny)

	// Roles inheritance
	r.AddGroupingPolicy(RoleAdmin, RoleUser)
	r.AddGroupingPolicy(RoleSuperAdmin, RoleAdmin)

	r.GetModel().PrintPolicy()

	return r
}
