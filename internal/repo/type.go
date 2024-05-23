package repo

// * definition of custom filters
type (
	// UsersFilter represents the filter type for listing and filtering users
	UsersFilter struct {
		Search string
	}

	// DocumentsFilter represents the filter type for listing and filtering documents
	DocumentsFilter struct {
		UserID string
		Search string
	}
)
