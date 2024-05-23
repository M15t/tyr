package types

// Profile model
// swagger:model
type Profile struct {
	Base
	UserID           string `json:"user_id"`
	PlaidAccessToken string `json:"plaid_access_token"`
}
