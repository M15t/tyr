package document

import (
	"io"
	contextutil "tyr/internal/api/context"
	"tyr/internal/repo"
	"tyr/third_party/azure"

	"github.com/M15t/gram/pkg/rbac"
)

// New creates new document application service
func New(repo *repo.Service, rbac rbac.Intf, cr Crypter, azure Azure) *Document {
	return &Document{repo: repo, rbac: rbac, cr: cr, azure: azure}
}

// Document represents document application service
type Document struct {
	repo  *repo.Service
	rbac  rbac.Intf
	cr    Crypter
	azure Azure
}

// Azure represents azure interface
type Azure interface {
	AnalyzeDocument(c contextutil.Context, modelID, apiVersion string, payload io.Reader) (*azure.ResponseHeaders, error)
	GetAnalyzeDocument(c contextutil.Context, url string) (*azure.ResultAnalyzeResponse, error)
}

// Crypter represents security interface
type Crypter interface {
}
