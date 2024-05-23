package azure

import (
	"tyr/config"
	"tyr/internal/repo"
)

// Service represents azure service
type Service struct {
	cfg  config.Azure
	repo *repo.Service
}

// New returns azure service
func New(cfg config.Azure, repo *repo.Service) *Service {
	return &Service{cfg, repo}
}
