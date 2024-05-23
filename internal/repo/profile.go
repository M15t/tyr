package repo

import (
	"tyr/internal/types"

	repoutil "github.com/M15t/gram/pkg/util/repo"
	"gorm.io/gorm"
)

// Profile represents the client for profile table
type Profile struct {
	*repoutil.Repo[types.Profile]
}

// NewProfile returns a new profile database instance
func NewProfile(gdb *gorm.DB) *Profile {
	return &Profile{repoutil.NewRepo[types.Profile](gdb)}
}
