package types

import "time"

// cosnt
const (
	UserStatusUnknown Status = iota
	UserStatusActive
	UserStatusBlocked
	UserStatausDeleted
)

// Status represents the status of user
type Status int

func (s Status) String() string {
	switch s {
	case UserStatusActive:
		return "active"
	case UserStatusBlocked:
		return "blocked"
	case UserStatausDeleted:
		return "deleted"
	default:
		return "unknown"
	}
}

// User represents the user model
// swagger:model
type User struct {
	Base
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`

	Password  string     `json:"-" gorm:"not null"`
	LastLogin *time.Time `json:"last_login,omitempty"`

	Phone           string     `json:"phone" gorm:"uniqueIndex:uix_users_phone"`
	PhoneVerifiedAt *time.Time `json:"phone_verified_at,omitempty"`
	OTP             *string    `json:"-" gorm:"type:varchar(10)"`
	OTPSentAt       *time.Time `json:"-"`
	Email           string     `json:"email" gorm:"uniqueIndex:uix_users_email"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`

	Status string `json:"status" gorm:"type:varchar(20);default:active"` // active || blocked || deleted

	Profile *Profile `json:"profile,omitempty" gorm:"foreignkey:UserID"`
}
