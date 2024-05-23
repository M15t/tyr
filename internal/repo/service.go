package repo

import "gorm.io/gorm"

// Service provides all databases
type Service struct {
	User         *User
	Session      *Session
	ActivityLog  *ActivityLog
	Document     *Document
	DocumentItem *DocumentItem
	Profile      *Profile
}

// New creates db service
func New(db *gorm.DB) *Service {
	return &Service{
		User:         NewUser(db),
		Session:      NewSession(db),
		ActivityLog:  NewActivityLog(db),
		Document:     NewDocument(db),
		DocumentItem: NewDocumentItem(db),
		Profile:      NewProfile(db),
	}
}
