package auth

import "tyr/internal/types"

// Credentials represents login request data
// swagger:model
type Credentials struct {
	// example: collector@tyr.io
	Email string `json:"email" form:"email" validate:"required_without=Username"`
	// example: user123!@#
	Password string `json:"password" form:"password" validate:"required"`

	// This is for SwaggerUI authentication which only support `username` field
	// swagger:ignore
	Username string `json:"username" form:"username"`
	// example: app
	GrantType string `json:"grant_type" form:"grant_type" validate:"required"`
}

// RefreshTokenData represents refresh token request data
// swagger:model
type RefreshTokenData struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// AuthenticateInput represents internal authenticate data
type AuthenticateInput struct {
	User    *types.User
	IsLogin bool
}

// SignupData represents signup request data
// swagger:model
type SignupData struct {
	// example: John
	FirstName string `json:"first_name" validate:"required"`
	// example: Doe
	LastName string `json:"last_name" validate:"required"`
	// example: john.doe@tyr.io
	Email string `json:"email" validate:"required,email"`
	// example: 5551234567
	Phone string `json:"phone" validate:"required,max=10"`
	// example: passisburden!@#
	Password string `json:"password" validate:"required,min=6"`
}
